import logging

from typing import List
from typing import Dict
from typing import Any
from enum import Enum

logger = logging.getLogger().getChild("kgateway-ai-ext.streamchunkdata")


class StreamChunkDataType(Enum):
    NORMAL_TEXT = 1  # a normal chunk with content
    NORMAL_BINARY = (
        2  # a chunk contains binary data like audio and have not text content
    )
    FINISH = 3  # a chunk with finish_reason that's not null.
    FINISH_NO_CONTENT = (
        4  # a chunk with finish_reason that's not null and no text content.
    )
    # For openai, this chunk normal does not have content and it's before the DONE chunk
    DONE = 5  # the chunk with `data: [DONE]` and has no json data
    INVALID = 6  # chunk that contains invalid json or empty SSE message


class StreamChunkData:
    __slot__ = ("raw_data", "json_data", "contents")

    def __init__(self, raw_data, json_data, contents, type):
        self.raw_data: bytes = raw_data
        self.json_data: Dict[str, Any] | None = json_data
        self.contents: List[bytes] | None = contents
        self.type: StreamChunkDataType = type
        # TODO(andy): there doesn't seems to be much value to provide role to the webhook in the response
        #             because there is only one role in response. Some API does not have a role in the response
        #             at all
        # self.roles: List[bytes] | None = roles
        # TODO(andy): add usage (tokens) data here as well?

    def __repr__(self) -> str:
        return f"raw_data:\n{self.raw_data}\njson_data:\n{self.json_data}\ncontents:\n{self.contents}"

    def get_content_length(self, choice_index: int) -> int:
        if self.contents is None or len(self.contents) <= choice_index:
            return 0

        content: bytes = self.contents[choice_index]
        if any(b > 0x7F for b in content):
            # The bytes contains non-ascii character, mostly likely utf-8 encoded.
            # so need to return the decoded length instead
            return len(content.decode("utf-8"))

        return len(content)
