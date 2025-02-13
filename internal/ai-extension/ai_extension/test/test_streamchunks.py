import unittest

from ext_proc.streamchunks import StreamChunks
from typing import List
from ext_proc.provider import OpenAI
from test.test_data.sample_chunk_data import test_chunk_data, utf8_test_chunk_data
from copy import deepcopy


class TestStreamChunks(unittest.TestCase):
    def test_find_segment_boundary(self):
        chunks = StreamChunks()
        # The content from chunk 0 to 26 should be:
        #     In the heart of the code, a dance unfolds,
        #     Where whispers of logic, in layers, are told.
        #     A
        for chunk in test_chunk_data[0:27]:
            chunks.append_chunk(deepcopy(chunk))

        contents = chunks.get_contents_with_chunk_indices()
        matchData = chunks.find_segment_boundary(contents)
        assert matchData is not None
        assert len(matchData) == 1
        assert matchData[0].capture == ".  \n"
        assert matchData[0].start_pos == 89, (
            f"contents: {contents} len: {len(contents[0].content)}"
        )
        assert matchData[0].end_pos == 93, (
            f"contents: {contents} len: {len(contents[0].content)}"
        )

    def test_collapse_chunks_with_new_content(self):
        chunks = StreamChunks()
        # The content from chunk 0 to 10 should be:
        #     'In the heart of the code, a dance' (without quote)
        for chunk in test_chunk_data[0:10]:
            chunks.append_chunk(deepcopy(chunk))

        contents = chunks.get_contents_with_chunk_indices()
        print(contents)
        # The content from chunk 10 to 20 should be:
        #     ' unfolds,
        #      Where whispers of logic, in layers' (without quote)
        for chunk in test_chunk_data[10:20]:
            chunks.append_chunk(chunk)

        new_contents: List[str] = []
        new_contents.append("test")
        provider = OpenAI()
        chunks.collapse_chunks_with_new_content(
            llm_provider=provider, original_contents=contents, new_contents=new_contents
        )

        contents = chunks.get_contents_with_chunk_indices()
        print(contents)

        assert contents[0].begin_index == 0
        assert contents[0].end_index == 11
        assert (
            contents[0].content == "test unfolds,  \nWhere whispers of logic, in layers"
        )

        chunk = chunks.pop_chunk()
        assert chunk is not None
        assert chunk.contents is not None
        assert chunk.contents[0] == b"test"

        contents = chunks.get_contents_with_chunk_indices()
        print(contents)

        assert contents[0].begin_index == 0
        assert contents[0].end_index == 10
        assert contents[0].content == " unfolds,  \nWhere whispers of logic, in layers"

    def test_align_contents_for_guardrail(self):
        chunks = StreamChunks()
        # The content from chunk 0 to 27 should be:
        #     In the heart of the code, a dance unfolds,
        #     Where whispers of logic, in layers, are told.
        #     A mystery
        # Testing the boundary indicator span 2 chunks (chunk #23 and #24)
        # but 2nd chunk ends with the indicator
        for chunk in test_chunk_data[0:27]:
            chunks.append_chunk(deepcopy(chunk))

        provider = OpenAI()
        contents = chunks.get_contents_with_chunk_indices()
        print(contents)
        # TODO(andy): test multi-choice response
        all_chunks_contents_before_alignment = contents[0].content
        boundary_indicator = ".  \n"
        assert not contents[0].content.endswith(boundary_indicator)
        chunks_to_pop = chunks.align_contents_for_guardrail(
            llm_provider=provider, contents=contents, min_content_length=20
        )
        print(contents)
        assert contents[0].content.endswith(boundary_indicator)
        assert chunks_to_pop == 25
        contents = chunks.get_contents_with_chunk_indices()
        assert contents[0].content == all_chunks_contents_before_alignment

        chunks.pop_chunks(chunks_to_pop)
        contents = chunks.get_contents_with_chunk_indices()
        print(contents)

        # Testing only double newline in a single chunk (chunk #34)
        chunks.pop_all()
        for chunk in test_chunk_data[0:37]:
            chunks.append_chunk(deepcopy(chunk))
        contents = chunks.get_contents_with_chunk_indices()
        print(contents)
        all_chunks_contents_before_alignment = contents[0].content
        boundary_indicator = "\n\n"
        assert not contents[0].content.endswith(boundary_indicator)
        chunks_to_pop = chunks.align_contents_for_guardrail(
            llm_provider=provider, contents=contents, min_content_length=20
        )
        print(contents)
        assert contents[0].content.endswith(boundary_indicator)
        contents = chunks.get_contents_with_chunk_indices()
        print(contents)
        assert contents[0].content == all_chunks_contents_before_alignment
        assert chunks_to_pop == 35

        # Testing the boundary indicator span 2 chunks (chunk #7 and #8)
        # but extra data in 2nd chunk
        chunks.pop_all()
        for chunk in test_chunk_data[0:10]:
            chunks.append_chunk(deepcopy(chunk))
        contents = chunks.get_contents_with_chunk_indices()
        print(contents)
        all_chunks_contents_before_alignment = contents[0].content
        boundary_indicator = ". "
        assert not contents[0].content.endswith(boundary_indicator)
        chunks_to_pop = chunks.align_contents_for_guardrail(
            llm_provider=provider, contents=contents, min_content_length=20
        )
        print(contents)
        assert contents[0].content.endswith(boundary_indicator)
        contents = chunks.get_contents_with_chunk_indices()
        print(contents)
        assert contents[0].content == all_chunks_contents_before_alignment
        assert chunks_to_pop == 8

        # TODO(andy): test boundary pos less than min_content_length case

        # TODO(andy): test aligning in the last chunk

        # TODO(andy): test multi-choices response
        # https://github.com/solo-io/solo-projects/issues/7439

    def test_align_contents_for_guardrail_utf8(self):
        chunks = StreamChunks()
        provider = OpenAI()
        for i, chunk in enumerate(utf8_test_chunk_data[0:34]):
            chunks.append_chunk(deepcopy(chunk))

        contents = chunks.get_contents_with_chunk_indices()
        chunks_to_pop = chunks.align_contents_for_guardrail(
            llm_provider=provider, contents=contents, min_content_length=50
        )
        assert chunks_to_pop == 32

        chunks.pop_chunks(chunks_to_pop - 1)
        next_chunk = chunks.pop_chunk()
        # This is chunk #31 originally has contents=[b"!  "],
        # The space from the next chunk should moved up to here, so should be b"!     " now
        assert next_chunk is not None and next_chunk.contents is not None
        assert next_chunk.contents[0] == b"!     ", next_chunk.contents[0]
        next_chunk = chunks.pop_chunk()

        # This is chunk #32 originally has contents=[b"   \xf0\x9f\x96\xb1"],
        # the space should have moved up to chunk #31, so should be b"\xf0\x9f\x96\xb1" now
        assert next_chunk is not None and next_chunk.contents is not None
        assert next_chunk.contents[0] == b"\xf0\x9f\x96\xb1", next_chunk.contents[0]
        chunks.pop_all()
