import json

from ext_proc.streamchunkdata import StreamChunkData, StreamChunkDataType
from typing import List

utf8_test_chunk_data: List[StreamChunkData] = [
    # 0
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"role":"assistant","content":"","refusal":null},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"role":"assistant","content":"","refusal":null},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b""],
    ),
    # 1
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":"Of"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":"Of"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b"Of"],
    ),
    # 2
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" course"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" course"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b" course"],
    ),
    # 3
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b","],
    ),
    # 4
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" here\'s"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" here\'s"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b" here's"],
    ),
    # 5
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" another"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" another"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b" another"],
    ),
    # 6
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" one"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" one"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b" one"],
    ),
    # 7
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" for"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" for"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b" for"],
    ),
    # 8
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" you"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" you"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b" you"],
    ),
    # 9
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":":\\n\\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":":\\n\\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b":\n\n"],
    ),
    # 10
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":"Why"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":"Why"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b"Why"],
    ),
    # 11
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" do"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" do"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b" do"],
    ),
    # 12
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" programmers"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" programmers"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b" programmers"],
    ),
    # 13
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" prefer"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" prefer"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b" prefer"],
    ),
    # 14
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" using"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" using"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b" using"],
    ),
    # 15
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" the"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" the"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b" the"],
    ),
    # 16
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" keyboard"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" keyboard"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b" keyboard"],
    ),
    # 17
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":"?\\n\\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":"?\\n\\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b"?\n\n"],
    ),
    # 18
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":"Because"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":"Because"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b"Because"],
    ),
    # 19
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" they"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" they"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b" they"],
    ),
    # 20
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" don"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" don"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b" don"],
    ),
    # 21
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":"\xe2\x80\x99t"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":"\xe2\x80\x99t"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b"\xe2\x80\x99t"],
    ),
    # 22
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" want"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" want"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b" want"],
    ),
    # 23
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" to"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" to"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b" to"],
    ),
    # 24
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" deal"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" deal"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b" deal"],
    ),
    # 25
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" with"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" with"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b" with"],
    ),
    # 26
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" a"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" a"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b" a"],
    ),
    # 27
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" **"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" **"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b" **"],
    ),
    # 28
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":"mouse"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":"mouse"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b"mouse \xf0\x9f\x96\xb1"],
    ),
    # 29
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" trap"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" trap"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b" trap"],
    ),
    # 30
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":"**"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":"**"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b"**"],
    ),
    # 31
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":"!"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":"!"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b"!  "],
    ),
    # 32
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" \xf0\x9f\x96\xb1"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":" \xf0\x9f\x96\xb1"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b"   \xf0\x9f\x96\xb1"],
    ),
    # 33
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":"\xef\xb8\x8f"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":"\xef\xb8\x8f"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b"\xef\xb8\x8f"],
    ),
    # 34
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":"\xf0\x9f\x98\x82"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{"content":"\xf0\x9f\x98\x82"},"logprobs":null,"finish_reason":null}]}'
        ),
        type=StreamChunkDataType.NORMAL_TEXT,
        contents=[b"\xf0\x9f\x98\x82"],
    ),
    # 35
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{},"logprobs":null,"finish_reason":"stop"}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AtIDcvSoOQYUP1FigEYZfNMXIL2Fo","object":"chat.completion.chunk","created":1737741436,"model":"chatgpt-4o-latest","service_tier":"default","system_fingerprint":"fp_60a3f2dc65","choices":[{"index":0,"delta":{},"logprobs":null,"finish_reason":"stop"}]}'
        ),
        type=StreamChunkDataType.FINISH_NO_CONTENT,
        contents=None,
    ),
    # 36
    StreamChunkData(
        raw_data=b"data: [DONE]\n\n",
        json_data=None,
        type=StreamChunkDataType.DONE,
        contents=None,
    ),
]

test_chunk_data: List[StreamChunkData] = [
    # 0
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"role":"assistant","content":"","refusal":null},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"role":"assistant","content":"","refusal":null},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b""],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 1
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"In"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"In"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"In"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 2
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" the"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" the"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" the"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 3
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" heart"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" heart"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" heart"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 4
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" of"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" of"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" of"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 5
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" the"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" the"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" the"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 6
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" code"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" code"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" code"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 7
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"."],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 8
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" a"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" a"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" a"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 9
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" dance"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" dance"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" dance"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 10
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" unfolds"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" unfolds"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" unfolds"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 11
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 12
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"  \n"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 13
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Where"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Where"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"Where"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 14
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" whispers"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" whispers"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" whispers"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 15
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" of"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" of"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" of"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 16
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" logic"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" logic"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" logic"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 17
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 18
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" in"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" in"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" in"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 19
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" layers"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" layers"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" layers"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 20
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 21
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" are"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" are"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" are"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 22
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" told"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" told"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" told"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 23
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"."},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"."},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"."],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 24
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"  \n"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 25
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"A"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"A"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"A"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 26
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" mystery"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" mystery"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" mystery"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 27
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" wrapped"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" wrapped"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" wrapped"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 28
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" in"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" in"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" in"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 29
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" loops"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" loops"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" loops"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 30
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" that"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" that"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" that"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 31
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" entw"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" entw"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" entw"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 32
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"ine"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"ine"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"ine"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 33
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 34
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"\\n\\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"\\n\\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"\n\n"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 35
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Rec"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Rec"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"Rec"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 36
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"ursion"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"ursion"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"ursion"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 37
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 38
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" dear"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" dear"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" dear"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 39
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" friend"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" friend"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" friend"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 40
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 41
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" is"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" is"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" is"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 42
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" both"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" both"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" both"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 43
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" simple"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" simple"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" simple"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 44
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" and"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" and"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" and"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 45
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" fine"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" fine"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" fine"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 46
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"."},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"."},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"."],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 47
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n\\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n\\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"  \n\n"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 48
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Imagine"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Imagine"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"Imagine"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 49
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" a"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" a"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" a"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 50
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" mirror"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" mirror"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" mirror"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 51
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 52
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" reflecting"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" reflecting"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" reflecting"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 53
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" a"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" a"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" a"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 54
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" face"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" face"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" face"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 55
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 56
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"  \n"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 57
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Each"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Each"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"Each"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 58
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" glance"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" glance"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" glance"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 59
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" leads"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" leads"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" leads"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 60
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" to"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" to"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" to"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 61
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" more"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" more"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" more"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 62
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 63
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" in"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" in"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" in"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 64
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" an"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" an"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" an"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 65
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" infinite"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" infinite"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" infinite"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 66
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" space"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" space"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" space"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 67
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"."},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"."},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"."],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 68
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"  \n"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 69
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"\\""},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"\\""},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b'"'],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 70
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Call"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Call"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"Call"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 71
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" me"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" me"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" me"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    # 72
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" again"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" again"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" again"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":",\\""},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":",\\""},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b',"'],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" says"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" says"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" says"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" the"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" the"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" the"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" function"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" function"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" function"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" with"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" with"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" with"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" g"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" g"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" g"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"lee"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"lee"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"lee"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"  \n"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"\\"I"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"\\"I"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b'"I'],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"\'ll"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"\'ll"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"'ll"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" solve"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" solve"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" solve"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" this"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" this"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" this"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" small"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" small"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" small"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" problem"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" problem"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" problem"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" then"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" then"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" then"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" set"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" set"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" set"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" you"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" you"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" you"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" free"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" free"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" free"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":".\\""},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":".\\""},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b'."'],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n\\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n\\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"  \n\n"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"The"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"The"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"The"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" base"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" base"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" base"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" case"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" case"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" case"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" is"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" is"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" is"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" vital"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" vital"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" vital"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" a"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" a"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" a"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" stop"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" stop"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" stop"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" in"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" in"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" in"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" the"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" the"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" the"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" flow"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" flow"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" flow"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"  \n"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"A"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"A"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"A"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" moment"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" moment"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" moment"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" of"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" of"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" of"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" pause"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" pause"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" pause"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" where"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" where"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" where"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" we"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" we"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" we"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" finally"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" finally"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" finally"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" know"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" know"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" know"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"\xe2\x80\x94"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"\xe2\x80\x94"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"\xe2\x80\x94"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"  \n"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"When"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"When"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"When"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" should"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" should"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" should"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" we"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" we"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" we"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" halt"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" halt"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" halt"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" when"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" when"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" when"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" should"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" should"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" should"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" we"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" we"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" we"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" refrain"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" refrain"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" refrain"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"?"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"?"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"?"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"  \n"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Without"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Without"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"Without"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" it"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" it"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" it"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" we"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" we"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" we"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"\xe2\x80\x99d"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"\xe2\x80\x99d"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"\xe2\x80\x99d"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" spiral"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" spiral"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" spiral"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" in"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" in"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" in"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" an"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" an"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" an"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" endless"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" endless"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" endless"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" refrain"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" refrain"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" refrain"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"."},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"."},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"."],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n\\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n\\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"  \n\n"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"A"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"A"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"A"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" factorial"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" factorial"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" factorial"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"\xe2\x80\x99s"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"\xe2\x80\x99s"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"\xe2\x80\x99s"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" journey"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" journey"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" journey"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" from"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" from"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" from"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" five"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" five"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" five"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" down"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" down"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" down"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" to"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" to"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" to"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" one"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" one"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" one"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"  \n"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Each"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Each"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"Each"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" step"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" step"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" step"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" of"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" of"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" of"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" its"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" its"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" its"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" path"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" path"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" path"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" a"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" a"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" a"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" new"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" new"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" new"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" call"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" call"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" call"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" begun"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" begun"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" begun"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"."},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"."},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"."],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"  \n"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Five"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Five"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"Five"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" times"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" times"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" times"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" four"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" four"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" four"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" times"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" times"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" times"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" three"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" three"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" three"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"\xe2\x80\x94"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"\xe2\x80\x94"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"\xe2\x80\x94"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"oh"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"oh"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"oh"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" what"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" what"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" what"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" a"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" a"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" a"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" delight"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" delight"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" delight"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"!"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"!"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"!"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"  \n"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Each"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Each"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"Each"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" function"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" function"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" function"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" returns"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" returns"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" returns"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" like"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" like"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" like"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" stars"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" stars"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" stars"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" back"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" back"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" back"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" to"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" to"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" to"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" night"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" night"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" night"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"."},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"."},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"."],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n\\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n\\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"  \n\n"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"In"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"In"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"In"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" trees"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" trees"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" trees"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" of"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" of"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" of"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" decisions"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" decisions"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" decisions"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" recursion"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" recursion"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" recursion"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" can"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" can"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" can"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" bloom"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" bloom"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" bloom"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"  \n"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Branch"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Branch"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"Branch"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"ing"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"ing"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"ing"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" out"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" out"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" out"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" softly"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" softly"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" softly"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" disp"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" disp"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" disp"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"elling"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"elling"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"elling"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" the"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" the"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" the"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" gloom"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" gloom"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" gloom"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"."},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"."},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"."],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"  \n"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"With"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"With"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"With"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" each"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" each"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" each"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" nested"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" nested"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" nested"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" call"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" call"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" call"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" we"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" we"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" we"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" traverse"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" traverse"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" traverse"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" with"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" with"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" with"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" grace"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" grace"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" grace"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"  \n"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Re"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Re"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"Re"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"vis"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"vis"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"vis"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"iting"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"iting"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"iting"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" paths"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" paths"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" paths"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" in"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" in"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" in"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" this"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" this"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" this"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" elegant"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" elegant"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" elegant"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" space"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" space"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" space"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"."},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"."},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"."],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n\\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n\\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"  \n\n"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"So"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"So"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"So"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" fear"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" fear"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" fear"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" not"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" not"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" not"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" the"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" the"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" the"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" depth"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" depth"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" depth"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" nor"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" nor"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" nor"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" the"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" the"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" the"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" layers"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" layers"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" layers"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" so"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" so"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" so"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" vast"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" vast"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" vast"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"  \n"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"For"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"For"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"For"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" with"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" with"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" with"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" eager"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" eager"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" eager"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" recursion"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" recursion"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" recursion"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" we"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" we"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" we"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" solve"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" solve"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" solve"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" the"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" the"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" the"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" cast"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" cast"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" cast"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"."},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"."},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"."],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"  \n"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Em"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"Em"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"Em"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"brace"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"brace"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"brace"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" this"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" this"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" this"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" technique"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" technique"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" technique"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" let"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" let"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" let"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" it"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" it"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" it"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" guide"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" guide"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" guide"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" your"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" your"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" your"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" way"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" way"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" way"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  \\n"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"  \n"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"In"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"In"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"In"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" the"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" the"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" the"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" realm"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" realm"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" realm"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" of"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" of"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" of"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" the"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" the"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" the"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" code"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" code"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" code"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" it"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" it"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" it"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"\xe2\x80\x99s"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"\xe2\x80\x99s"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"\xe2\x80\x99s"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" a"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" a"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" a"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" bright"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" bright"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" bright"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":","},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b","],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" shining"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" shining"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" shining"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" day"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":" day"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b" day"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"!"},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"!"},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"!"],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  "},"logprobs":null,"finish_reason":null}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{"content":"  "},"logprobs":null,"finish_reason":null}]}'
        ),
        contents=[b"  "],
        type=StreamChunkDataType.NORMAL_TEXT,
    ),
    StreamChunkData(
        raw_data=b'data: {"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{},"logprobs":null,"finish_reason":"stop"}]}\n\n',
        json_data=json.loads(
            '{"id":"chatcmpl-AVPZexJ0frlpPaMIa2MsW9a99fA1z","object":"chat.completion.chunk","created":1732049838,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_3de1288069","choices":[{"index":0,"delta":{},"logprobs":null,"finish_reason":"stop"}]}'
        ),
        contents=None,
        type=StreamChunkDataType.FINISH_NO_CONTENT,
    ),
    StreamChunkData(
        raw_data=b"data: [DONE]\n\n",
        json_data=None,
        contents=None,
        type=StreamChunkDataType.DONE,
    ),
]
