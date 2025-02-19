from api.enterprise.options.ai import ai_pb2 as ai_pb2
from api.envoy.config.core.v3 import base_pb2 as base_pb2
from langchain_openai import OpenAIEmbeddings, AzureOpenAIEmbeddings
from langchain_core.embeddings import Embeddings as LangchainEmbeddings
from pydantic import SecretStr

from util.proto import get_auth_token

from util.env import (
    open_ai_token_env,
    azure_open_ai_token_env,
)


def create_embeddings(
    embeddings: ai_pb2.Embedding,
    headers: base_pb2.HeaderMap,
) -> LangchainEmbeddings:
    match embeddings.WhichOneof("embedding"):
        case "openai":
            token = get_auth_token(
                embeddings.openai.auth_token, headers, open_ai_token_env
            )
            return OpenAIEmbeddings(api_key=SecretStr(token))
        case "azure_openai":
            token = get_auth_token(
                embeddings.azure_openai.auth_token, headers, azure_open_ai_token_env
            )
            return AzureOpenAIEmbeddings(
                api_key=SecretStr(token),
                azure_endpoint=embeddings.azure_openai.endpoint,
                azure_deployment=embeddings.azure_openai.deployment_name,
                api_version=embeddings.azure_openai.api_version,
            )
        case _:
            raise ValueError("Unknown embedding type")
