import os
import logging
from openai import OpenAI, AzureOpenAI
from mistralai import Mistral
from anthropic import Anthropic
from openai import NotGiven, NOT_GIVEN
from openai.types.chat.chat_completion import ChatCompletion
from openai import Stream, NotFoundError
from openai.types.chat.chat_completion_chunk import ChatCompletionChunk
import google.generativeai as genai
import vertexai
from vertexai.generative_models import GenerativeModel
from tenacity import retry, retry_if_exception_type, stop_after_attempt

from google.auth import credentials

logger = logging.getLogger(__name__)


class FakeCredentials(credentials.Credentials):
    def refresh(self, request):
        # Implement refresh if needed
        pass

    def before_request(self, request, method, url, headers):
        # Fake the before_request functionality
        headers["Authorization"] = "Bearer FAKE_TOKEN"


class PassthroughCredentials(credentials.Credentials):
    def refresh(self, request):
        # Implement refresh if needed
        pass

    def before_request(self, request, method, url, headers):
        # Passthrough the before_request functionality
        headers["Authorization"] = f"Bearer {os.environ.get('VERTEX_AI_API_KEY', '')}"


TEST_OPENAI_BASE_URL = os.environ.get("TEST_OPENAI_BASE_URL", "")
TEST_AZURE_OPENAI_BASE_URL = os.environ.get("TEST_AZURE_OPENAI_BASE_URL", "")
TEST_MISTRAL_BASE_URL = os.environ.get("TEST_MISTRAL_BASE_URL", "")
TEST_ANTHROPIC_BASE_URL = os.environ.get("TEST_ANTHROPIC_BASE_URL", "")
TEST_GEMINI_BASE_URL = os.environ.get("TEST_GEMINI_BASE_URL", "")
TEST_VERTEX_AI_BASE_URL = os.environ.get("TEST_VERTEX_AI_BASE_URL", "")

passthrough = os.environ.get("TEST_TOKEN_PASSTHROUGH", "false").lower() == "true"


class LLMClient:
    openai_client = OpenAI(
        api_key=os.environ.get("OPENAI_API_KEY", "") if passthrough else "FAKE",
        base_url=TEST_OPENAI_BASE_URL,
    )
    azure_openai_client = AzureOpenAI(
        api_key=(os.environ.get("AZURE_OPENAI_API_KEY", "") if passthrough else "FAKE"),
        base_url=TEST_AZURE_OPENAI_BASE_URL,
        max_retries=10,
        api_version="2024-02-15-preview",
    )
    mistral_client = Mistral(
        api_key=(os.environ.get("MISTRAL_API_KEY", "") if passthrough else "FAKE"),
        server_url=TEST_MISTRAL_BASE_URL,
    )
    anthropic_client = Anthropic(
        api_key=(os.environ.get("ANTHROPIC_API_KEY", "") if passthrough else "FAKE"),
        base_url=TEST_ANTHROPIC_BASE_URL,
    )

    genai.configure(
        api_key=(os.environ.get("GEMINI_API_KEY", "") if passthrough else "FAKE"),
        client_options={"api_endpoint": TEST_GEMINI_BASE_URL},
        transport="rest",
    )
    gemini_client = genai.GenerativeModel("gemini-1.5-flash-latest")

    vertexai.init(
        project="gloo-ee",
        location="us-central1",
        api_endpoint=TEST_VERTEX_AI_BASE_URL,
        api_transport="rest",
        api_key=(os.environ.get("VERTEX_AI_API_KEY", "") if passthrough else "FAKE"),
        credentials=PassthroughCredentials() if passthrough else FakeCredentials(),
    )
    vertex_ai_client = GenerativeModel("gemini-1.5-flash-001")

    @retry(
        retry=retry_if_exception_type(NotFoundError),
        stop=stop_after_attempt(3),
    )
    def openai_chat_completion(
        self,
        model: str,
        messages: list,
        tools: list | NotGiven = NOT_GIVEN,
    ) -> ChatCompletion:
        resp = self.openai_client.chat.completions.create(
            model=model, messages=messages, tools=tools
        )
        # Need to make sure it's not a direct response: 404
        assert resp is not None
        logger.debug(f"openai completion response:\n{resp}")
        return resp

    @retry(
        retry=retry_if_exception_type(NotFoundError),
        stop=stop_after_attempt(3),
    )
    def openai_chat_completion_stream(
        self,
        model: str,
        messages: list,
        tools: list | NotGiven = NOT_GIVEN,
    ) -> Stream[ChatCompletionChunk]:
        resp = self.openai_client.chat.completions.create(
            model=model,
            messages=messages,
            tools=tools,
            stream=True,
        )
        logger.debug(f"openai completion response:\n{resp}")
        # Need to make sure it's not a direct response: 404
        assert resp is not None
        return resp
