from solo_io.solo_kit.api.v1 import ref_pb2 as _ref_pb2
from google.protobuf import struct_pb2 as _struct_pb2
from extproto import ext_pb2 as _ext_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class SingleAuthToken(_message.Message):
    __slots__ = ("inline", "secret_ref", "passthrough")
    class Passthrough(_message.Message):
        __slots__ = ()
        def __init__(self) -> None: ...
    INLINE_FIELD_NUMBER: _ClassVar[int]
    SECRET_REF_FIELD_NUMBER: _ClassVar[int]
    PASSTHROUGH_FIELD_NUMBER: _ClassVar[int]
    inline: str
    secret_ref: _ref_pb2.ResourceRef
    passthrough: SingleAuthToken.Passthrough
    def __init__(self, inline: _Optional[str] = ..., secret_ref: _Optional[_Union[_ref_pb2.ResourceRef, _Mapping]] = ..., passthrough: _Optional[_Union[SingleAuthToken.Passthrough, _Mapping]] = ...) -> None: ...

class UpstreamSpec(_message.Message):
    __slots__ = ("openai", "mistral", "anthropic", "azure_openai", "multi", "gemini", "vertex_ai")
    class CustomHost(_message.Message):
        __slots__ = ("host", "port")
        HOST_FIELD_NUMBER: _ClassVar[int]
        PORT_FIELD_NUMBER: _ClassVar[int]
        host: str
        port: int
        def __init__(self, host: _Optional[str] = ..., port: _Optional[int] = ...) -> None: ...
    class OpenAI(_message.Message):
        __slots__ = ("auth_token", "custom_host", "model")
        AUTH_TOKEN_FIELD_NUMBER: _ClassVar[int]
        CUSTOM_HOST_FIELD_NUMBER: _ClassVar[int]
        MODEL_FIELD_NUMBER: _ClassVar[int]
        auth_token: SingleAuthToken
        custom_host: UpstreamSpec.CustomHost
        model: str
        def __init__(self, auth_token: _Optional[_Union[SingleAuthToken, _Mapping]] = ..., custom_host: _Optional[_Union[UpstreamSpec.CustomHost, _Mapping]] = ..., model: _Optional[str] = ...) -> None: ...
    class AzureOpenAI(_message.Message):
        __slots__ = ("auth_token", "endpoint", "deployment_name", "api_version")
        AUTH_TOKEN_FIELD_NUMBER: _ClassVar[int]
        ENDPOINT_FIELD_NUMBER: _ClassVar[int]
        DEPLOYMENT_NAME_FIELD_NUMBER: _ClassVar[int]
        API_VERSION_FIELD_NUMBER: _ClassVar[int]
        auth_token: SingleAuthToken
        endpoint: str
        deployment_name: str
        api_version: str
        def __init__(self, auth_token: _Optional[_Union[SingleAuthToken, _Mapping]] = ..., endpoint: _Optional[str] = ..., deployment_name: _Optional[str] = ..., api_version: _Optional[str] = ...) -> None: ...
    class Gemini(_message.Message):
        __slots__ = ("auth_token", "model", "api_version")
        AUTH_TOKEN_FIELD_NUMBER: _ClassVar[int]
        MODEL_FIELD_NUMBER: _ClassVar[int]
        API_VERSION_FIELD_NUMBER: _ClassVar[int]
        auth_token: SingleAuthToken
        model: str
        api_version: str
        def __init__(self, auth_token: _Optional[_Union[SingleAuthToken, _Mapping]] = ..., model: _Optional[str] = ..., api_version: _Optional[str] = ...) -> None: ...
    class VertexAI(_message.Message):
        __slots__ = ("auth_token", "model", "api_version", "project_id", "location", "model_path", "publisher")
        class Publisher(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
            __slots__ = ()
            GOOGLE: _ClassVar[UpstreamSpec.VertexAI.Publisher]
        GOOGLE: UpstreamSpec.VertexAI.Publisher
        AUTH_TOKEN_FIELD_NUMBER: _ClassVar[int]
        MODEL_FIELD_NUMBER: _ClassVar[int]
        API_VERSION_FIELD_NUMBER: _ClassVar[int]
        PROJECT_ID_FIELD_NUMBER: _ClassVar[int]
        LOCATION_FIELD_NUMBER: _ClassVar[int]
        MODEL_PATH_FIELD_NUMBER: _ClassVar[int]
        PUBLISHER_FIELD_NUMBER: _ClassVar[int]
        auth_token: SingleAuthToken
        model: str
        api_version: str
        project_id: str
        location: str
        model_path: str
        publisher: UpstreamSpec.VertexAI.Publisher
        def __init__(self, auth_token: _Optional[_Union[SingleAuthToken, _Mapping]] = ..., model: _Optional[str] = ..., api_version: _Optional[str] = ..., project_id: _Optional[str] = ..., location: _Optional[str] = ..., model_path: _Optional[str] = ..., publisher: _Optional[_Union[UpstreamSpec.VertexAI.Publisher, str]] = ...) -> None: ...
    class Mistral(_message.Message):
        __slots__ = ("auth_token", "custom_host", "model")
        AUTH_TOKEN_FIELD_NUMBER: _ClassVar[int]
        CUSTOM_HOST_FIELD_NUMBER: _ClassVar[int]
        MODEL_FIELD_NUMBER: _ClassVar[int]
        auth_token: SingleAuthToken
        custom_host: UpstreamSpec.CustomHost
        model: str
        def __init__(self, auth_token: _Optional[_Union[SingleAuthToken, _Mapping]] = ..., custom_host: _Optional[_Union[UpstreamSpec.CustomHost, _Mapping]] = ..., model: _Optional[str] = ...) -> None: ...
    class Anthropic(_message.Message):
        __slots__ = ("auth_token", "custom_host", "version", "model")
        AUTH_TOKEN_FIELD_NUMBER: _ClassVar[int]
        CUSTOM_HOST_FIELD_NUMBER: _ClassVar[int]
        VERSION_FIELD_NUMBER: _ClassVar[int]
        MODEL_FIELD_NUMBER: _ClassVar[int]
        auth_token: SingleAuthToken
        custom_host: UpstreamSpec.CustomHost
        version: str
        model: str
        def __init__(self, auth_token: _Optional[_Union[SingleAuthToken, _Mapping]] = ..., custom_host: _Optional[_Union[UpstreamSpec.CustomHost, _Mapping]] = ..., version: _Optional[str] = ..., model: _Optional[str] = ...) -> None: ...
    class MultiPool(_message.Message):
        __slots__ = ("priorities",)
        class Backend(_message.Message):
            __slots__ = ("openai", "mistral", "anthropic", "azure_openai", "gemini", "vertex_ai")
            OPENAI_FIELD_NUMBER: _ClassVar[int]
            MISTRAL_FIELD_NUMBER: _ClassVar[int]
            ANTHROPIC_FIELD_NUMBER: _ClassVar[int]
            AZURE_OPENAI_FIELD_NUMBER: _ClassVar[int]
            GEMINI_FIELD_NUMBER: _ClassVar[int]
            VERTEX_AI_FIELD_NUMBER: _ClassVar[int]
            openai: UpstreamSpec.OpenAI
            mistral: UpstreamSpec.Mistral
            anthropic: UpstreamSpec.Anthropic
            azure_openai: UpstreamSpec.AzureOpenAI
            gemini: UpstreamSpec.Gemini
            vertex_ai: UpstreamSpec.VertexAI
            def __init__(self, openai: _Optional[_Union[UpstreamSpec.OpenAI, _Mapping]] = ..., mistral: _Optional[_Union[UpstreamSpec.Mistral, _Mapping]] = ..., anthropic: _Optional[_Union[UpstreamSpec.Anthropic, _Mapping]] = ..., azure_openai: _Optional[_Union[UpstreamSpec.AzureOpenAI, _Mapping]] = ..., gemini: _Optional[_Union[UpstreamSpec.Gemini, _Mapping]] = ..., vertex_ai: _Optional[_Union[UpstreamSpec.VertexAI, _Mapping]] = ...) -> None: ...
        class Priority(_message.Message):
            __slots__ = ("pool",)
            POOL_FIELD_NUMBER: _ClassVar[int]
            pool: _containers.RepeatedCompositeFieldContainer[UpstreamSpec.MultiPool.Backend]
            def __init__(self, pool: _Optional[_Iterable[_Union[UpstreamSpec.MultiPool.Backend, _Mapping]]] = ...) -> None: ...
        PRIORITIES_FIELD_NUMBER: _ClassVar[int]
        priorities: _containers.RepeatedCompositeFieldContainer[UpstreamSpec.MultiPool.Priority]
        def __init__(self, priorities: _Optional[_Iterable[_Union[UpstreamSpec.MultiPool.Priority, _Mapping]]] = ...) -> None: ...
    OPENAI_FIELD_NUMBER: _ClassVar[int]
    MISTRAL_FIELD_NUMBER: _ClassVar[int]
    ANTHROPIC_FIELD_NUMBER: _ClassVar[int]
    AZURE_OPENAI_FIELD_NUMBER: _ClassVar[int]
    MULTI_FIELD_NUMBER: _ClassVar[int]
    GEMINI_FIELD_NUMBER: _ClassVar[int]
    VERTEX_AI_FIELD_NUMBER: _ClassVar[int]
    openai: UpstreamSpec.OpenAI
    mistral: UpstreamSpec.Mistral
    anthropic: UpstreamSpec.Anthropic
    azure_openai: UpstreamSpec.AzureOpenAI
    multi: UpstreamSpec.MultiPool
    gemini: UpstreamSpec.Gemini
    vertex_ai: UpstreamSpec.VertexAI
    def __init__(self, openai: _Optional[_Union[UpstreamSpec.OpenAI, _Mapping]] = ..., mistral: _Optional[_Union[UpstreamSpec.Mistral, _Mapping]] = ..., anthropic: _Optional[_Union[UpstreamSpec.Anthropic, _Mapping]] = ..., azure_openai: _Optional[_Union[UpstreamSpec.AzureOpenAI, _Mapping]] = ..., multi: _Optional[_Union[UpstreamSpec.MultiPool, _Mapping]] = ..., gemini: _Optional[_Union[UpstreamSpec.Gemini, _Mapping]] = ..., vertex_ai: _Optional[_Union[UpstreamSpec.VertexAI, _Mapping]] = ...) -> None: ...

class RouteSettings(_message.Message):
    __slots__ = ("prompt_enrichment", "prompt_guard", "rag", "semantic_cache", "defaults", "route_type")
    class RouteType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        CHAT: _ClassVar[RouteSettings.RouteType]
        CHAT_STREAMING: _ClassVar[RouteSettings.RouteType]
    CHAT: RouteSettings.RouteType
    CHAT_STREAMING: RouteSettings.RouteType
    PROMPT_ENRICHMENT_FIELD_NUMBER: _ClassVar[int]
    PROMPT_GUARD_FIELD_NUMBER: _ClassVar[int]
    RAG_FIELD_NUMBER: _ClassVar[int]
    SEMANTIC_CACHE_FIELD_NUMBER: _ClassVar[int]
    DEFAULTS_FIELD_NUMBER: _ClassVar[int]
    ROUTE_TYPE_FIELD_NUMBER: _ClassVar[int]
    prompt_enrichment: AIPromptEnrichment
    prompt_guard: AIPromptGuard
    rag: RAG
    semantic_cache: SemanticCache
    defaults: _containers.RepeatedCompositeFieldContainer[FieldDefault]
    route_type: RouteSettings.RouteType
    def __init__(self, prompt_enrichment: _Optional[_Union[AIPromptEnrichment, _Mapping]] = ..., prompt_guard: _Optional[_Union[AIPromptGuard, _Mapping]] = ..., rag: _Optional[_Union[RAG, _Mapping]] = ..., semantic_cache: _Optional[_Union[SemanticCache, _Mapping]] = ..., defaults: _Optional[_Iterable[_Union[FieldDefault, _Mapping]]] = ..., route_type: _Optional[_Union[RouteSettings.RouteType, str]] = ...) -> None: ...

class FieldDefault(_message.Message):
    __slots__ = ("field", "value", "override")
    FIELD_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    OVERRIDE_FIELD_NUMBER: _ClassVar[int]
    field: str
    value: _struct_pb2.Value
    override: bool
    def __init__(self, field: _Optional[str] = ..., value: _Optional[_Union[_struct_pb2.Value, _Mapping]] = ..., override: bool = ...) -> None: ...

class Postgres(_message.Message):
    __slots__ = ("connection_string", "collection_name")
    CONNECTION_STRING_FIELD_NUMBER: _ClassVar[int]
    COLLECTION_NAME_FIELD_NUMBER: _ClassVar[int]
    connection_string: str
    collection_name: str
    def __init__(self, connection_string: _Optional[str] = ..., collection_name: _Optional[str] = ...) -> None: ...

class Embedding(_message.Message):
    __slots__ = ("openai", "azure_openai")
    class OpenAI(_message.Message):
        __slots__ = ("auth_token",)
        AUTH_TOKEN_FIELD_NUMBER: _ClassVar[int]
        auth_token: SingleAuthToken
        def __init__(self, auth_token: _Optional[_Union[SingleAuthToken, _Mapping]] = ...) -> None: ...
    class AzureOpenAI(_message.Message):
        __slots__ = ("auth_token", "api_version", "endpoint", "deployment_name")
        AUTH_TOKEN_FIELD_NUMBER: _ClassVar[int]
        API_VERSION_FIELD_NUMBER: _ClassVar[int]
        ENDPOINT_FIELD_NUMBER: _ClassVar[int]
        DEPLOYMENT_NAME_FIELD_NUMBER: _ClassVar[int]
        auth_token: SingleAuthToken
        api_version: str
        endpoint: str
        deployment_name: str
        def __init__(self, auth_token: _Optional[_Union[SingleAuthToken, _Mapping]] = ..., api_version: _Optional[str] = ..., endpoint: _Optional[str] = ..., deployment_name: _Optional[str] = ...) -> None: ...
    OPENAI_FIELD_NUMBER: _ClassVar[int]
    AZURE_OPENAI_FIELD_NUMBER: _ClassVar[int]
    openai: Embedding.OpenAI
    azure_openai: Embedding.AzureOpenAI
    def __init__(self, openai: _Optional[_Union[Embedding.OpenAI, _Mapping]] = ..., azure_openai: _Optional[_Union[Embedding.AzureOpenAI, _Mapping]] = ...) -> None: ...

class SemanticCache(_message.Message):
    __slots__ = ("datastore", "embedding", "ttl", "mode")
    class Mode(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        READ_WRITE: _ClassVar[SemanticCache.Mode]
        READ_ONLY: _ClassVar[SemanticCache.Mode]
    READ_WRITE: SemanticCache.Mode
    READ_ONLY: SemanticCache.Mode
    class Redis(_message.Message):
        __slots__ = ("connection_string", "score_threshold")
        CONNECTION_STRING_FIELD_NUMBER: _ClassVar[int]
        SCORE_THRESHOLD_FIELD_NUMBER: _ClassVar[int]
        connection_string: str
        score_threshold: float
        def __init__(self, connection_string: _Optional[str] = ..., score_threshold: _Optional[float] = ...) -> None: ...
    class Weaviate(_message.Message):
        __slots__ = ("host", "http_port", "grpc_port", "insecure")
        HOST_FIELD_NUMBER: _ClassVar[int]
        HTTP_PORT_FIELD_NUMBER: _ClassVar[int]
        GRPC_PORT_FIELD_NUMBER: _ClassVar[int]
        INSECURE_FIELD_NUMBER: _ClassVar[int]
        host: str
        http_port: int
        grpc_port: int
        insecure: bool
        def __init__(self, host: _Optional[str] = ..., http_port: _Optional[int] = ..., grpc_port: _Optional[int] = ..., insecure: bool = ...) -> None: ...
    class DataStore(_message.Message):
        __slots__ = ("redis", "weaviate")
        REDIS_FIELD_NUMBER: _ClassVar[int]
        WEAVIATE_FIELD_NUMBER: _ClassVar[int]
        redis: SemanticCache.Redis
        weaviate: SemanticCache.Weaviate
        def __init__(self, redis: _Optional[_Union[SemanticCache.Redis, _Mapping]] = ..., weaviate: _Optional[_Union[SemanticCache.Weaviate, _Mapping]] = ...) -> None: ...
    DATASTORE_FIELD_NUMBER: _ClassVar[int]
    EMBEDDING_FIELD_NUMBER: _ClassVar[int]
    TTL_FIELD_NUMBER: _ClassVar[int]
    MODE_FIELD_NUMBER: _ClassVar[int]
    datastore: SemanticCache.DataStore
    embedding: Embedding
    ttl: int
    mode: SemanticCache.Mode
    def __init__(self, datastore: _Optional[_Union[SemanticCache.DataStore, _Mapping]] = ..., embedding: _Optional[_Union[Embedding, _Mapping]] = ..., ttl: _Optional[int] = ..., mode: _Optional[_Union[SemanticCache.Mode, str]] = ...) -> None: ...

class RAG(_message.Message):
    __slots__ = ("datastore", "embedding", "prompt_template")
    class DataStore(_message.Message):
        __slots__ = ("postgres",)
        POSTGRES_FIELD_NUMBER: _ClassVar[int]
        postgres: Postgres
        def __init__(self, postgres: _Optional[_Union[Postgres, _Mapping]] = ...) -> None: ...
    DATASTORE_FIELD_NUMBER: _ClassVar[int]
    EMBEDDING_FIELD_NUMBER: _ClassVar[int]
    PROMPT_TEMPLATE_FIELD_NUMBER: _ClassVar[int]
    datastore: RAG.DataStore
    embedding: Embedding
    prompt_template: str
    def __init__(self, datastore: _Optional[_Union[RAG.DataStore, _Mapping]] = ..., embedding: _Optional[_Union[Embedding, _Mapping]] = ..., prompt_template: _Optional[str] = ...) -> None: ...

class AIPromptEnrichment(_message.Message):
    __slots__ = ("prepend", "append")
    class Message(_message.Message):
        __slots__ = ("role", "content")
        ROLE_FIELD_NUMBER: _ClassVar[int]
        CONTENT_FIELD_NUMBER: _ClassVar[int]
        role: str
        content: str
        def __init__(self, role: _Optional[str] = ..., content: _Optional[str] = ...) -> None: ...
    PREPEND_FIELD_NUMBER: _ClassVar[int]
    APPEND_FIELD_NUMBER: _ClassVar[int]
    prepend: _containers.RepeatedCompositeFieldContainer[AIPromptEnrichment.Message]
    append: _containers.RepeatedCompositeFieldContainer[AIPromptEnrichment.Message]
    def __init__(self, prepend: _Optional[_Iterable[_Union[AIPromptEnrichment.Message, _Mapping]]] = ..., append: _Optional[_Iterable[_Union[AIPromptEnrichment.Message, _Mapping]]] = ...) -> None: ...

class AIPromptGuard(_message.Message):
    __slots__ = ("request", "response")
    class Regex(_message.Message):
        __slots__ = ("matches", "builtins", "action")
        class BuiltIn(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
            __slots__ = ()
            SSN: _ClassVar[AIPromptGuard.Regex.BuiltIn]
            CREDIT_CARD: _ClassVar[AIPromptGuard.Regex.BuiltIn]
            PHONE_NUMBER: _ClassVar[AIPromptGuard.Regex.BuiltIn]
            EMAIL: _ClassVar[AIPromptGuard.Regex.BuiltIn]
        SSN: AIPromptGuard.Regex.BuiltIn
        CREDIT_CARD: AIPromptGuard.Regex.BuiltIn
        PHONE_NUMBER: AIPromptGuard.Regex.BuiltIn
        EMAIL: AIPromptGuard.Regex.BuiltIn
        class Action(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
            __slots__ = ()
            MASK: _ClassVar[AIPromptGuard.Regex.Action]
            REJECT: _ClassVar[AIPromptGuard.Regex.Action]
        MASK: AIPromptGuard.Regex.Action
        REJECT: AIPromptGuard.Regex.Action
        class RegexMatch(_message.Message):
            __slots__ = ("pattern", "name")
            PATTERN_FIELD_NUMBER: _ClassVar[int]
            NAME_FIELD_NUMBER: _ClassVar[int]
            pattern: str
            name: str
            def __init__(self, pattern: _Optional[str] = ..., name: _Optional[str] = ...) -> None: ...
        MATCHES_FIELD_NUMBER: _ClassVar[int]
        BUILTINS_FIELD_NUMBER: _ClassVar[int]
        ACTION_FIELD_NUMBER: _ClassVar[int]
        matches: _containers.RepeatedCompositeFieldContainer[AIPromptGuard.Regex.RegexMatch]
        builtins: _containers.RepeatedScalarFieldContainer[AIPromptGuard.Regex.BuiltIn]
        action: AIPromptGuard.Regex.Action
        def __init__(self, matches: _Optional[_Iterable[_Union[AIPromptGuard.Regex.RegexMatch, _Mapping]]] = ..., builtins: _Optional[_Iterable[_Union[AIPromptGuard.Regex.BuiltIn, str]]] = ..., action: _Optional[_Union[AIPromptGuard.Regex.Action, str]] = ...) -> None: ...
    class Webhook(_message.Message):
        __slots__ = ("host", "port", "forwardHeaders")
        class HeaderMatch(_message.Message):
            __slots__ = ("key", "match_type")
            class MatchType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
                __slots__ = ()
                EXACT: _ClassVar[AIPromptGuard.Webhook.HeaderMatch.MatchType]
                PREFIX: _ClassVar[AIPromptGuard.Webhook.HeaderMatch.MatchType]
                SUFFIX: _ClassVar[AIPromptGuard.Webhook.HeaderMatch.MatchType]
                CONTAINS: _ClassVar[AIPromptGuard.Webhook.HeaderMatch.MatchType]
                REGEX: _ClassVar[AIPromptGuard.Webhook.HeaderMatch.MatchType]
            EXACT: AIPromptGuard.Webhook.HeaderMatch.MatchType
            PREFIX: AIPromptGuard.Webhook.HeaderMatch.MatchType
            SUFFIX: AIPromptGuard.Webhook.HeaderMatch.MatchType
            CONTAINS: AIPromptGuard.Webhook.HeaderMatch.MatchType
            REGEX: AIPromptGuard.Webhook.HeaderMatch.MatchType
            KEY_FIELD_NUMBER: _ClassVar[int]
            MATCH_TYPE_FIELD_NUMBER: _ClassVar[int]
            key: str
            match_type: AIPromptGuard.Webhook.HeaderMatch.MatchType
            def __init__(self, key: _Optional[str] = ..., match_type: _Optional[_Union[AIPromptGuard.Webhook.HeaderMatch.MatchType, str]] = ...) -> None: ...
        HOST_FIELD_NUMBER: _ClassVar[int]
        PORT_FIELD_NUMBER: _ClassVar[int]
        FORWARDHEADERS_FIELD_NUMBER: _ClassVar[int]
        host: str
        port: int
        forwardHeaders: _containers.RepeatedCompositeFieldContainer[AIPromptGuard.Webhook.HeaderMatch]
        def __init__(self, host: _Optional[str] = ..., port: _Optional[int] = ..., forwardHeaders: _Optional[_Iterable[_Union[AIPromptGuard.Webhook.HeaderMatch, _Mapping]]] = ...) -> None: ...
    class Moderation(_message.Message):
        __slots__ = ("openai",)
        class OpenAI(_message.Message):
            __slots__ = ("model", "auth_token")
            MODEL_FIELD_NUMBER: _ClassVar[int]
            AUTH_TOKEN_FIELD_NUMBER: _ClassVar[int]
            model: str
            auth_token: SingleAuthToken
            def __init__(self, model: _Optional[str] = ..., auth_token: _Optional[_Union[SingleAuthToken, _Mapping]] = ...) -> None: ...
        OPENAI_FIELD_NUMBER: _ClassVar[int]
        openai: AIPromptGuard.Moderation.OpenAI
        def __init__(self, openai: _Optional[_Union[AIPromptGuard.Moderation.OpenAI, _Mapping]] = ...) -> None: ...
    class Request(_message.Message):
        __slots__ = ("custom_response", "regex", "webhook", "moderation")
        class CustomResponse(_message.Message):
            __slots__ = ("message", "status_code")
            MESSAGE_FIELD_NUMBER: _ClassVar[int]
            STATUS_CODE_FIELD_NUMBER: _ClassVar[int]
            message: str
            status_code: int
            def __init__(self, message: _Optional[str] = ..., status_code: _Optional[int] = ...) -> None: ...
        CUSTOM_RESPONSE_FIELD_NUMBER: _ClassVar[int]
        REGEX_FIELD_NUMBER: _ClassVar[int]
        WEBHOOK_FIELD_NUMBER: _ClassVar[int]
        MODERATION_FIELD_NUMBER: _ClassVar[int]
        custom_response: AIPromptGuard.Request.CustomResponse
        regex: AIPromptGuard.Regex
        webhook: AIPromptGuard.Webhook
        moderation: AIPromptGuard.Moderation
        def __init__(self, custom_response: _Optional[_Union[AIPromptGuard.Request.CustomResponse, _Mapping]] = ..., regex: _Optional[_Union[AIPromptGuard.Regex, _Mapping]] = ..., webhook: _Optional[_Union[AIPromptGuard.Webhook, _Mapping]] = ..., moderation: _Optional[_Union[AIPromptGuard.Moderation, _Mapping]] = ...) -> None: ...
    class Response(_message.Message):
        __slots__ = ("regex", "webhook")
        REGEX_FIELD_NUMBER: _ClassVar[int]
        WEBHOOK_FIELD_NUMBER: _ClassVar[int]
        regex: AIPromptGuard.Regex
        webhook: AIPromptGuard.Webhook
        def __init__(self, regex: _Optional[_Union[AIPromptGuard.Regex, _Mapping]] = ..., webhook: _Optional[_Union[AIPromptGuard.Webhook, _Mapping]] = ...) -> None: ...
    REQUEST_FIELD_NUMBER: _ClassVar[int]
    RESPONSE_FIELD_NUMBER: _ClassVar[int]
    request: AIPromptGuard.Request
    response: AIPromptGuard.Response
    def __init__(self, request: _Optional[_Union[AIPromptGuard.Request, _Mapping]] = ..., response: _Optional[_Union[AIPromptGuard.Response, _Mapping]] = ...) -> None: ...
