from copy import deepcopy
import logging
import re
import traceback

from api.kgateway.policy.ai import prompt_guard
from api.envoy.service.ext_proc.v3 import external_processor_pb2
from collections import deque
from dataclasses import dataclass
from ext_proc.provider import Provider, Tokens
from ext_proc.streamchunkdata import StreamChunkData, StreamChunkDataType


# from ext_proc.stream import Handler
from opentelemetry import trace
from opentelemetry.trace.propagation.tracecontext import TraceContextTextMapPropagator
from presidio_analyzer import EntityRecognizer
from presidio_anonymizer import AnonymizerEngine
from telemetry.tracing import OtelTracer
from typing import Deque
from typing import List
from typing import Callable
from typing import Dict
from typing import Any
from guardrails.regex import regex_transform
from guardrails.webhook import call_response_webhook
from util import sse

logger = logging.getLogger().getChild("kgateway-ai-ext.streamchunks")


def update_chunk(llm_provider: Provider, chunk: StreamChunkData):
    """
    update the content and raw_data inside the chunk base on the json_data inside the chunk
    """
    chunk.contents = llm_provider.extract_contents_from_resp_chunk(chunk.json_data)
    logger.debug(f"json_data before: {chunk.json_data}")
    logger.debug(f"raw_data before: {chunk.raw_data}")

    if chunk.json_data is None:
        return

    try:
        chunk.raw_data = sse.replace_json_data(
            raw_data=chunk.raw_data, json_data=chunk.json_data
        )
    except sse.SSEParsingException as e:
        logger.error(f"update_chunk: failed to replace json data: {e}")

    logger.debug(f"raw_data after: {chunk.raw_data}")


@dataclass
class StreamChunksContent:
    """
    StreamChunksContent hold the content as a str and the begin index (inclusive) and end index
    (exclusive) of the chunks in the streaming_fifo that make up the content string
    """

    content: str = ""
    begin_index: int = -1
    end_index: int = -1


@dataclass
class BoundaryMatchData:
    """
    BoundaryMatchData is a dataclass that store the captured pattern and start and end
    position from a regex match
    """

    choice_index: int = -1
    """
    The index into the contents array in SteamChunks when multi-choice response is used
    """

    capture: str = ""
    """
    The string matching the capture group pattern. An empty string here means there was no match
    """

    start_pos: int = -1
    """
    The start position where capture appears on the original string (inclusive)
    """

    end_pos: int = -1
    """
    The end position where capture ends on the original string (exclusive)
    """


class StreamChunks:
    __total_bytes_buffered: int = 0
    """
    This is the regex pattern to detect segment boundary in content message so we can send a complete
    segment to webhook or regex match.
    """

    def __init__(self):
        self.__streaming_fifo: Deque["StreamChunkData"] = deque()
        """
        streaming_fifo is used to buffer the content if handler.resp_webhook or handler.resp_regex
        is set. The data the come out of this FIFO will be append to body for caching later.
        This is not used for request.
        """

        self.__contents: List[bytearray] = list()
        """
        The content message in the streaming fifo
        """

        self.__leftover: bytes = bytes()
        """
        The leftover bytes of any incomplete chunks from the previous buffer() call
        """

        self.model: str = ""
        """
        The model from the streaming response
        """

        self.tokens: Tokens | None = None
        """
        Accumulated tokens for the response stream
        """

        self.is_function_calling: bool = False
        """
        Indicate if this streaming response is a function calling response
        """

        self.is_completed: bool = False
        """
        This indicated the stream is completed properly base on the provider specific indicator
        """

    def get_contents_with_chunk_indices(self) -> List[StreamChunksContent]:
        contents: List[StreamChunksContent] = list()
        for item in self.__contents:
            contents.append(
                StreamChunksContent(
                    content=item.decode("utf-8"),
                    begin_index=0,
                    end_index=len(self.__streaming_fifo),
                )
            )
        return contents

    def reconstruct_contents(self):
        for item in self.__contents:
            item.clear()

        for chunk in self.__streaming_fifo:
            if chunk.contents is None:
                continue
            for choice_index, content in enumerate(chunk.contents):
                self.__contents[choice_index].extend(content)

    def append_chunk(self, chunk: StreamChunkData):
        """
        append one chunk to the end of streaming_fifo, update contents and do the necessary accounting
        """
        self.__streaming_fifo.append(chunk)
        StreamChunks.__total_bytes_buffered += len(chunk.raw_data)
        isFirstChunkWithContent = False
        if len(self.__contents) == 0:
            isFirstChunkWithContent = True
        if chunk.contents is not None:
            if not isFirstChunkWithContent and len(self.__contents) != len(
                chunk.contents
            ):
                logger.critical(
                    f"contents list size not matching. self.__contents {len(self.__contents)} vs chunk.contents {len(chunk.contents)}!"
                )
            for i, content in enumerate(chunk.contents):
                if isFirstChunkWithContent:
                    logger.debug(f"type of content: {type(content)}")
                    print(f"type of content: {type(content)}")
                    self.__contents.append(bytearray(content))
                else:
                    self.__contents[i].extend(content)
                    i += 1

    def pop_chunk(self) -> StreamChunkData | None:
        """
        pop one chunk from the head of streaming_fifo, update contents and do the necessary accounting
        """
        if len(self.__streaming_fifo) == 0:
            return None

        chunk = self.__streaming_fifo.popleft()
        StreamChunks.__total_bytes_buffered -= len(chunk.raw_data)
        if chunk.contents is not None:
            for i, content in enumerate(chunk.contents):
                if self.__contents[i].startswith(content):
                    del self.__contents[i][: len(content)]
                else:
                    logger.critical(
                        f'"{self.__contents[i]}" does not start with "{content}"'
                    )

        return chunk

    def pop_all(self) -> bytes:
        """
        pop_all pops all the chunks in streaming_filo and return all the raw_data in a single bytes object
        """
        raw_data = bytearray()
        for i, chunk in enumerate(self.__streaming_fifo):
            StreamChunks.__total_bytes_buffered -= len(chunk.raw_data)
            logger.debug(
                f"pop_all: i: {i} total_bytes_buffered: {StreamChunks.__total_bytes_buffered} raw_data: {chunk.raw_data}"
            )
            raw_data.extend(chunk.raw_data)
        self.__contents.clear()
        self.__streaming_fifo.clear()
        return bytes(raw_data)

    def pop_chunks(self, n: int) -> bytes | None:
        """
        pop_chunks pops n chunks from streaming_fifo and returns the raw_data in those chunks in a single bytes object
        if n is greater or equal to the size of streaming_fifo, it would just call pop_all()
        """

        if n <= 0:
            return None

        if n >= len(self.__streaming_fifo):
            return self.pop_all()

        raw_data = bytearray()
        for i in range(n):
            chunk = self.pop_chunk()
            if chunk is not None:
                # because we already check size above, chunk should never be None
                raw_data.extend(chunk.raw_data)
                logger.debug(f"pop_chunks: n: {n} i: {i} raw_data: {chunk.raw_data}")
            else:
                logger.debug(f"pop_chunks: n: {n} i: {i} chunk is None")
        return bytes(raw_data)

    def has_min_chunks_with_contents(self, min: int) -> bool:
        count = 0
        for chunk in self.__streaming_fifo:
            if (
                chunk.type == StreamChunkDataType.NORMAL_TEXT
                or chunk.type == StreamChunkDataType.FINISH
            ):
                count += 1
                if count >= min:
                    return True

        return False

    def align_contents_for_guardrail(
        self,
        llm_provider: Provider,
        contents: List[StreamChunksContent],
        min_content_length: int,
    ) -> int:
        """
        find the segment boundary in each content and adjust the content text and chunk indices
        re-align the chunks in the fifo if necessary to make all the end indices on the same chunk
        returns how many chunks should be popped out to send back to envoy
        """
        # TODO(andy): This is option 2 in the design doc. When we have more than 1 option, need to get
        #             which option from x-resp-guardrails-config:

        if not self.has_min_chunks_with_contents(2):
            # We need at least 2 chunks that have contents to move things around
            # For the end_of_steam scenario where we may not have 2 chunks, this function
            # should not be called as we just send all contents to guardrails at that point
            return 0

        matchList = self.find_segment_boundary(contents)
        if len(matchList) != len(contents):
            # This means we don't have a match or we have multi-choice response and not
            # every content in every choice has a match.
            # We wait until we have a match in every content and collapse the chunks
            # before webhook.
            return 0

        logger.debug(f"align: 0 contents {contents} matchList {matchList}")
        for match in matchList:
            if match.end_pos < min_content_length:
                logger.debug(
                    f"align: 0.1 boundary too short: contents {contents} matchList {matchList}"
                )
                return 0

        for match in matchList:
            # each match corresponds to one choice in the response
            # find which chunk(s) contains the matched capture. Starts from the end of the fifo
            # when found, split the content in the chunk so everything before and including the boundary
            # indicator will be moved to the chunk before this chunk in the fifo
            contentData = contents[match.choice_index]
            bytes_count_from_end = len(contentData.content) - match.end_pos
            for chunk_reverse_index, chunk in enumerate(
                reversed(self.__streaming_fifo)
            ):
                if chunk.contents is None:
                    # chunk.contents can be None for chunks that has no content field. eg the last chunk
                    # sometimes can contain of the the stop reason but no content
                    continue
                logger.debug(
                    f"align: 1.1 chunk {chunk} bytes_count_from_end: {bytes_count_from_end}"
                )
                content_len = chunk.get_content_length(match.choice_index)
                bytes_count_from_end -= content_len
                logger.debug(
                    f"reverse_index: {chunk_reverse_index}, content_len: {content_len}, bytes_count_from_end: {bytes_count_from_end}"
                )
                if bytes_count_from_end < 0:
                    # we found the chunk that contains the end of the matched capture
                    # and we want to split the content to keep what's after the matched capture
                    # in this chunk but put everything else in the chunk that's before this chunk
                    # in the fifo
                    content = deepcopy(chunk.contents[match.choice_index])
                    chunk_index = len(self.__streaming_fifo) - 1 - chunk_reverse_index

                    logger.debug(
                        f"align: 1 chunk_index: {chunk_index} match: {match} contentData: {contentData}"
                    )
                    # Adjust the contentData to strip out everything after the boundary indicator
                    # as they will not be passed to webhook or regex match
                    contentData.end_index = chunk_index
                    contentData.content = contentData.content[: match.end_pos]

                    # matchedBytes is the detected boundary indicator base on the regex match pattern
                    matchedBytes = match.capture.encode("utf-8")
                    matchedBytesLen = len(matchedBytes)
                    logger.debug(
                        f"align: 1.7 matchedBytes: {matchedBytes} content: {content}"
                    )
                    if content.endswith(matchedBytes) or (
                        len(content) < matchedBytesLen
                        and matchedBytes.endswith(content)
                    ):
                        # This means the chunk is ending with the boundary indicator, so just include this chunk and
                        # nothing else needs to be done
                        contentData.end_index += 1
                        break

                    # This is where the matchedBytes are either split between 2 or more chunks or
                    # at the beginning of the current chunk. So, the current chunk should contain the
                    # whole matchedBytes or the last part of the matchedBytes.

                    # For the current chunk, we need to strip out any of the matchedBytes at the
                    # beginning.
                    matchedBytesIndex = 0
                    for matchedBytesIndex in range(matchedBytesLen):
                        if content.startswith(matchedBytes[matchedBytesIndex:]):
                            break
                    logger.debug(f"matchedBytesIndex: {matchedBytesIndex}")
                    llm_provider.update_stream_resp_contents(
                        json_data=chunk.json_data,
                        choice_index=match.choice_index,
                        content=content[matchedBytesLen - matchedBytesIndex :],
                    )
                    update_chunk(llm_provider, chunk)

                    # For the chunk in front of the current chunk, we append the part
                    # of the matchedBytes we stripped off from the current chunk
                    chunk_index -= 1  # get to the previous chunk in the fifo
                    logger.debug(
                        f"align: 2 chunk_index: {chunk_index} match: {match} contentData: {contentData}"
                    )
                    if chunk_index >= 0:
                        prev_chunk = self.__streaming_fifo[chunk_index]
                        if prev_chunk.contents is not None:
                            new_content = (
                                prev_chunk.contents[match.choice_index]
                                + matchedBytes[matchedBytesIndex:]
                            )
                            llm_provider.update_stream_resp_contents(
                                json_data=prev_chunk.json_data,
                                choice_index=match.choice_index,
                                content=new_content,
                            )
                            update_chunk(llm_provider, prev_chunk)
                        logger.debug(
                            f"align: 3 chunk_index: {chunk_index} prev_chunk: {prev_chunk}"
                        )
                    else:
                        logger.critical(
                            "align_contents_for_guardrail: chunk index reached 0 unexpectedly"
                        )
                        pass

                    break

        if len(contents) > 1:
            # TODO(andy): This logic was assuming the multi-choices are in the array in a single chunk but
            #             OpenAI streaming response does not do that. They send the choices in different chunk
            #             So, this logic needs to be re-do. Need to check other providers like gemini and see
            #             how they do multi-choices streaming response.

            # Handle more than 1 choice. Align the chunks if the boundary for different
            # choices are at different chunk_index
            highest_end_index = contents[0].end_index
            lowest_end_index = contents[0].end_index
            for item in contents:
                if item.end_index > highest_end_index:
                    highest_end_index = item.end_index
                if item.end_index < lowest_end_index:
                    lowest_end_index = item.end_index

            # TODO(andy): bounds check the indices?
            if lowest_end_index != highest_end_index:
                # end indices are not aligned, need to re-align the chunks
                # so all the end indices are on the same chunk
                for choice_index, content in enumerate(contents):
                    if content.end_index == highest_end_index:
                        # move all contents to the chunk before lowest_end_index
                        self.collapse_contents(
                            choice_index,
                            llm_provider.update_stream_resp_contents,
                            lowest_end_index,
                            content.end_index,
                            lowest_end_index - 1,
                        )
                    elif content.end_index > lowest_end_index:
                        # move contents before and including the boundary to the chunk before lowest_end_index
                        self.collapse_contents(
                            choice_index,
                            llm_provider.update_stream_resp_contents,
                            lowest_end_index,
                            content.end_index,
                            lowest_end_index - 1,
                        )
                        # move the rest contents to highest_end_index chunk
                        self.collapse_contents(
                            choice_index,
                            llm_provider.update_stream_resp_contents,
                            content.end_index,
                            highest_end_index,
                            highest_end_index,
                        )
                    content.end_index = lowest_end_index

            # reconstruct all the data in these 2 chunks that we move contents into
            update_chunk(llm_provider, self.__streaming_fifo[lowest_end_index - 1])
            update_chunk(llm_provider, self.__streaming_fifo[highest_end_index])

            # delete all the chunck that we move the contents away
            # TODO(andy): need to move away the token usage data as well
            # at the point, the self.__contents should still be valid as we have not
            # change any actual contents
            for i in range(lowest_end_index, highest_end_index):
                del self.__streaming_fifo[i]

        # self.reconstruct_contents()
        # at this point, all chunks content should be aligned, so returning the end_index
        # of the first choice as the number of chunks to pop out and return to envoy for delivery
        logger.debug(f"align: 4.5 contents {self.__contents[0]}")
        logger.debug(f"align: 5 returning  {contents}")
        return contents[0].end_index

    def get_usage_from_chunks(
        self,
        extract_usage_func: Callable[[Dict[str, Any]], Tokens],
        start_index: int,
        end_index: int,
    ) -> Tokens:
        logger.debug(
            f"get_usage_from_chunks: start_index {start_index} end_index {end_index}"
        )
        total = Tokens()
        for i in range(start_index, end_index):
            json_data = self.__streaming_fifo[i].json_data
            if json_data is None:
                continue

            chunk_tokens = extract_usage_func(json_data)

            # For gemini, the prompt token is repeated in every chunk, so we should not add them up
            if total.prompt == 0:
                total.prompt = chunk_tokens.prompt

            total.completion += chunk_tokens.completion

            logger.debug(f"get_usage_from_chunks: i: {i} total {total}")

        # For OpenAI, the usage token is null in the chunk, so the total prompt and completion token would be 0
        return total

    def collapse_chunks_with_new_content(
        self,
        llm_provider: Provider,
        original_contents: List[StreamChunksContent],
        new_contents: List[str],
    ) -> int:
        """
        delete the first n - 1 chunks we are collapsing and then update the 1st chunk with the new content
        return the number of chunks we should pop out. Usually it would be 1 but if we are at the end of stream,
        the last few chunks might be a FINISH_NO_CONTENT or DONE chunks that we want to preserve as is and don't
        collapse them.
        """

        # all contents should have the same begin and end index and begin_index of every content should be 0,
        # so just using the first one
        # TODO(andy): maybe double check?

        logger.debug(f"regex: collapse_chunks: {original_contents}")
        fifo_size = len(self.__streaming_fifo)
        logger.debug(
            f"regex: collapse_chunks fifo (size={fifo_size}): {self.__streaming_fifo}"
        )
        end_index = original_contents[0].end_index
        begin_index = original_contents[0].begin_index
        num_chunks_to_collapse = end_index - begin_index

        if num_chunks_to_collapse < 1:
            # This should not happen
            logger.critical("num chunks to collapse is less then 1")
            return 0

        if end_index > fifo_size:
            # This should not happen
            logger.critical(
                f"collapse_chunks_with_new_contents: end_index out of bound. end_index: {end_index} fifo_size: {fifo_size}"
            )
            return 0

        if begin_index != 0:
            # This should not happen
            logger.critical(
                f"collapse_chunks_with_new_contents: begin_index {begin_index} is not 0 fifo_size: {fifo_size}"
            )
            return 0

        chunks_to_pop = 1
        for i in range(end_index - 1, begin_index, -1):
            logger.debug(f"collapse_chunks_with_new_contents: i = {i}")
            match self.__streaming_fifo[i].type:
                case StreamChunkDataType.DONE:
                    num_chunks_to_collapse -= 1
                    chunks_to_pop += 1
                case StreamChunkDataType.FINISH_NO_CONTENT:
                    num_chunks_to_collapse -= 1
                    chunks_to_pop += 1
                case StreamChunkDataType.INVALID:
                    num_chunks_to_collapse -= 1
                    chunks_to_pop += 1
                case _:
                    # break as soon as we see a normal chunk
                    break

        # we are getting the total usage for the first num_chunks_to_collapse, so this will
        # include the usage of the chunk that we will keep for storing the new contents
        usages = self.get_usage_from_chunks(
            llm_provider.tokens, 0, num_chunks_to_collapse
        )
        # pop out the first n - 1 chunks
        self.pop_chunks(num_chunks_to_collapse - 1)

        if (
            self.__streaming_fifo[0].type != StreamChunkDataType.NORMAL_TEXT
            and self.__streaming_fifo[0].type != StreamChunkDataType.FINISH
        ):
            logger.critical(
                "collapse_chunks_with_new_content: no normal chunks to set content after collapsing"
            )
            return 0

        for choice_index, content in enumerate(new_contents):
            llm_provider.update_stream_resp_contents(
                self.__streaming_fifo[0].json_data,
                choice_index,
                content.encode("utf-8"),
            )

        if (
            usages.completion > 0
            and usages.prompt > 0
            and self.__streaming_fifo[0].json_data
            is not None  # This is just to make pyright happy as the Chunk type check above should have prevented this
        ):
            # only set the token if both are non-zero because for OpenAI, the usage object is null and
            # will appear as 0. We don't want to change the null to 0. For others, they should not be 0
            # but if they are, that means we don't need to change it because the chunk is already 0
            logger.debug(f"before usage update: {self.__streaming_fifo[0].json_data}")
            llm_provider.update_stream_resp_usage_token(
                self.__streaming_fifo[0].json_data, usages
            )
            logger.debug(f"after usage update: {self.__streaming_fifo[0].json_data}")

        # The first chunk now has the new contents and the collapsed total usage,
        # Update the raw_data and content of the chunk and reconstruct the contents to maintain consistency
        update_chunk(llm_provider, self.__streaming_fifo[0])
        self.reconstruct_contents()

        logger.debug(f"collapse_chunks_with_new_content: returning {chunks_to_pop}")
        return chunks_to_pop

    async def do_guardrails_check(
        self,
        llm_provider: Provider,
        resp_headers: dict[str, str],
        webhook: prompt_guard.Webhook | None,
        regex: list[EntityRecognizer] | None,
        anonymizer_engine: AnonymizerEngine,
        parent_span: trace.Span,
        final: bool = False,
    ) -> int:
        """
        return how many chunks we should pop out from the fifo. 0 means we are just buffering until we get enough.
        """
        # this contents is a copy and will be modified locally
        contents = self.get_contents_with_chunk_indices()
        logger.debug(
            f"webhook (fifo: {len(self.__streaming_fifo)} final={final}): {contents}"
        )

        # TODO(andy): This is deviated from the original design that use a minimum chunks (mainly for simplicity)
        #             but turns out using minimum chunks do not work well with gemini because it packs a lot of
        #             tokens in a single chunk where OpenAI packs only a few tokens at most per chunk.
        #             According to ChatGPT, the average sentence for a chat is around 25 to 75 characters long,
        #             so picking 50 here. Do we need this to be configurable and get this from x-resp-guardrails-config?
        min_content_length = 50
        should_do_guardrails_check = True
        for content_data in contents:
            if len(content_data.content) < min_content_length:
                should_do_guardrails_check = False
                break

        if not final and not should_do_guardrails_check:
            # buffer the minimum before we do any guardrails check unless this is the final check (end of stream)
            return 0

        chunks_to_pop_out = len(self.__streaming_fifo)
        if not final:
            chunks_to_pop_out = self.align_contents_for_guardrail(
                llm_provider, contents, min_content_length
            )
        else:
            # if it's final, send all contents as is to guardrail
            pass

        if chunks_to_pop_out == 0:
            # didn't find any boundary indicator
            return 0

        # The order is important here as the webhook_modified_content will be passed into regex match
        # if we are to change the order in the future, need to make sure the regex_modifed_content is
        # used to pass into webhook. Cannot just swap the 2 sections.
        webhook_modified = False
        webhook_modified_contents: List[str] | None = None
        if webhook:
            with OtelTracer.get().start_as_current_span(
                "webhook",
                context=trace.set_span_in_context(parent_span),
            ):
                headers = deepcopy(resp_headers)
                TraceContextTextMapPropagator().inject(headers)
                (
                    webhook_modified,
                    webhook_modified_contents,
                ) = await call_response_webhook(
                    webhook_host=webhook.host,
                    webhook_port=webhook.port,
                    headers=headers,
                    contents=(content_data.content for content_data in contents),
                )
                if webhook_modified and webhook_modified_contents is not None:
                    if len(webhook_modified_contents) != len(contents):
                        logger.error(
                            f"guardrail response webhook response does not contains all choices of the original content {len(contents)} vs {len(webhook_modified_contents)} "
                        )
                        # set this to None so it won't get used
                        webhook_modified_contents = None
                    elif regex:
                        # we only need to do this if regex is also enabled; otherwise, webhook_modified_contents will be used directly
                        for i, item in enumerate(contents):
                            item.content = webhook_modified_contents[i]
        regex_modified = False
        regex_modified_contents: List[str] = []
        if regex:
            # regex_transform can throw RegexRejection exception. Deliberately not catching
            # it here so it bubbles up all the way to server.Process() so it can construct an
            # immediate error response there.
            with OtelTracer.get().start_as_current_span(
                "regex",
                context=trace.set_span_in_context(parent_span),
            ):
                for i, item in enumerate(contents):
                    logger.debug(f"regex: choice_index: {i} content: {item.content}")
                    regex_modified_contents.append(
                        regex_transform("", item.content, regex, anonymizer_engine)
                    )
                    if item.content != regex_modified_contents[i]:
                        # as long as there is one choice that got modified, we need to collapse
                        # the chunks for all choices so they are aligned
                        regex_modified = True
                        logger.debug(
                            f"regex: choice_index: {i} modifed_content: {regex_modified_contents[i]}"
                        )

        if regex_modified:
            # if webhook has modified the contents, the modified contents would have already pass into regex
            # so, only use the webhook_modified_contents if regex didn't modify them
            return self.collapse_chunks_with_new_content(
                llm_provider, contents, regex_modified_contents
            )
        elif webhook_modified and webhook_modified_contents is not None:
            return self.collapse_chunks_with_new_content(
                llm_provider, contents, webhook_modified_contents
            )

        # no modification, so pop out the chunks that already went through guardrails
        return chunks_to_pop_out

    def collect_stream_info(self, llm_provider: Provider, chunk: StreamChunkData):
        """
        This function get called on each chunk of the streaming response and try to
        collect information about the streaming response for use later
        """
        if not self.is_completed:
            # allow chunk.json_data to be None for calling this because on some API
            # the "completion" is indicated by a SSE data tag ["DONE"] without any json data
            self.is_completed = llm_provider.is_streaming_response_completed(chunk)

        if chunk.json_data is None:
            return

        if self.tokens is None:
            self.tokens = llm_provider.tokens(chunk.json_data)
        else:
            self.tokens += llm_provider.tokens(chunk.json_data)

        if self.model == "":
            self.model = llm_provider.get_model_resp(chunk.json_data)

        if not self.is_function_calling:
            self.is_function_calling = llm_provider.has_function_call_finish_reason(
                chunk.json_data
            )

    async def buffer(
        self,
        llm_provider: Provider,
        resp_webhook: prompt_guard.Webhook | None,
        resp_regex: list[EntityRecognizer] | None,
        anonymizer_engine: AnonymizerEngine,
        resp_headers: dict[str, str],
        resp_body: external_processor_pb2.HttpBody,
        parent_span: trace.Span,
    ) -> bytes | None:
        """
        Buffer data for Guardrail. Returns the bytes when the data comes out of the Fifo
        """
        if resp_webhook is None and resp_regex is None:
            # Guardrail feature is not enabled, so no need to buffer
            return resp_body.body

        try:
            chunks, self.__leftover = sse.parse_sse_messages(
                llm_provider=llm_provider,
                data=resp_body.body,
                prev_leftover=self.__leftover,
            )
            for chunk in chunks:
                logger.debug("    StreamChunkData(")
                logger.debug(f"        raw_data = {chunk.raw_data}")
                logger.debug(f"        json_data = {chunk.json_data}")
                logger.debug(f"        type = {chunk.type.name}")
                logger.debug(f"        contents = {chunk.contents}")
                # TODO(andy): if chunk type is BINARY, flush all the chunk and stop buffering until we see a text chunk again
                self.collect_stream_info(llm_provider=llm_provider, chunk=chunk)

                self.append_chunk(chunk)
        except Exception as exc:
            logger.error(f"error parsing_stream_chunks, {exc}")
            print(traceback.format_exc())
            return resp_body.body

        if resp_body.end_of_stream and self.__leftover:
            logger.critical(
                f"reached end of stream but still has leftover data: {self.__leftover}"
            )
            self.append_chunk(
                StreamChunkData(
                    raw_data=self.__leftover,
                    json_data=None,
                    contents=None,
                    type=StreamChunkDataType.INVALID,
                )
            )

        number_messages_to_remove = await self.do_guardrails_check(
            final=resp_body.end_of_stream,
            llm_provider=llm_provider,
            resp_headers=resp_headers,
            regex=resp_regex,
            webhook=resp_webhook,
            anonymizer_engine=anonymizer_engine,
            parent_span=parent_span,
        )

        return self.pop_chunks(number_messages_to_remove)

    #    __boundaryRegex = re.compile(r'([.,?!:;] +\n*|\n+)')
    __boundaryRegex = re.compile(
        r"([.?!;] +\n*|\n+)"
    )  # TODO(andy): should this be configurable

    def find_segment_boundary(
        self, contents: List[StreamChunksContent]
    ) -> List[BoundaryMatchData]:
        """
        Find a boundaryPattern specified by __boundaryRegex from the current contents in self.__contents that
        indicates a semantic segment boundary and return a BoundaryMatchData for the content that has a match.
        The BoundaryMatchData object only contains the last match (closest to the end)
        A boundary indicator is pre-compiled in __boundaryRegex class variable and is one of the
        punctuation . , , , ? , ! , : , ; followed by any white space or a newline by itself.

        Note about the newline character:
        A newline from the llm in the json field is converted to "\n" (2 characters across the wire) when the actual
        newline character is assigned to the field. When assigned the field value to a str in python or encode() into bytes,
        it's converted back to a single character. So, when matching, we need to match the single newline character
        and not 2 characters. This is probably true for any special ascii characters.
        """

        result: List[BoundaryMatchData] = []
        for i, item in enumerate(contents):
            for match in reversed(
                list(StreamChunks.__boundaryRegex.finditer(item.content))
            ):
                logger.debug(
                    f"found boundary: choice_index: {i} group: {match.group()} groups: {match.groups()} [{match.pos}, {match.endpos}) content: {item.content}"
                )
                result.append(
                    BoundaryMatchData(
                        choice_index=i,
                        capture=match.group(),
                        start_pos=match.start(),
                        end_pos=match.end(),
                    )
                )
                break

        return result

    def collapse_contents(
        self,
        choice_index: int,
        json_content_update_func: Callable[[Dict[str, Any], int, bytes], None],
        src_start_index: int,
        src_end_index: int,
        dst_index: int,
    ):
        """
        collapse_contents move content for one choice specified by choice_index from the src chunks pointed to
        from src_start_index (inclusive) to src_end_index (exclusive) into the dst chunk
        This function only get used from the logic where we assume the multi-choices response are in the array in
        a single chunk so we need to align the detected boundary in every choices to the same chunk. Turns out OpenAI
        doesn't do that and need to re-think the logic to support that.
        """
        dst_chunk = self.__streaming_fifo[dst_index]
        if dst_chunk.contents is None or dst_chunk.json_data is None:
            logger.critical(
                f"collapse_contents: dst chunk has no data! indexes: src_start: {src_start_index} src_end: {src_end_index} dst: {dst_index}"
            )
            return

        new_content = bytearray()
        # TODO(andy): bounds check all the indices
        prepend_to_dst = False
        if dst_index >= src_end_index:
            # prepend contents to dst chunk
            prepend_to_dst = True
        elif dst_index < src_start_index:
            # append contents to dst chunk
            dst_contents = dst_chunk.contents
            if dst_contents is not None and len(dst_contents) > choice_index:
                new_content.extend(dst_contents[choice_index])
        else:
            logger.critical(
                f"collapse_contents: invalid indexes: src_start: {src_start_index} src_end: {src_end_index} dst: {dst_index}"
            )
            return

        for i in range(src_start_index, src_end_index):
            contents = self.__streaming_fifo[i].contents
            if contents is None or len(contents) <= choice_index:
                continue
            new_content.extend(contents[choice_index])

        if prepend_to_dst and dst_chunk.contents:
            new_content.extend(dst_chunk.contents[choice_index])

        json_content_update_func(dst_chunk.json_data, choice_index, bytes(new_content))
