import json
from dataclasses import dataclass, field
from typing import Optional, List
from ..ai.authtoken import SingleAuthToken


@dataclass
class CustomResponse:
    message: Optional[str] = "The request was rejected due to inappropriate content"
    status_code: Optional[int] = 403

    @staticmethod
    def from_json(data: dict) -> "CustomResponse":
        return CustomResponse(
            message=data.get(
                "message", "The request was rejected due to inappropriate content"
            ),
            status_code=data.get("status_code", 403),
        )


@dataclass
class RegexMatch:
    pattern: Optional[str] = None
    name: Optional[str] = None

    @staticmethod
    def from_json(data: dict) -> "RegexMatch":
        return RegexMatch(pattern=data.get("pattern"), name=data.get("name"))


class BuiltIn:
    SSN = "SSN"
    CREDIT_CARD = "CREDIT_CARD"
    PHONE_NUMBER = "PHONE_NUMBER"
    EMAIL = "EMAIL"


class Action:
    MASK = "MASK"
    REJECT = "REJECT"


@dataclass
class Regex:
    matches: Optional[List[RegexMatch]] = field(default_factory=list)
    builtins: Optional[List[BuiltIn]] = field(default_factory=list)
    action: Optional[str] = Action.MASK  # Use Action class for default

    @staticmethod
    def from_json(data: dict) -> "Regex":
        matches = [RegexMatch.from_json(m) for m in data.get("matches", [])]
        builtins = [BuiltIn(b) for b in data.get("builtins", [])]
        return Regex(
            matches=matches, builtins=builtins, action=data.get("action", Action.MASK)
        )


@dataclass
class Host:
    host: str
    port: int

    @staticmethod
    def from_json(data: dict) -> "Host":
        return Host(host=data["host"], port=data["port"])


class Type:
    EXACT = "Exact"
    REGULAR_EXPRESSION = "RegularExpression"


@dataclass
class HTTPHeaderMatch:
    name: str
    value: str
    type: Optional[Type] = Type.EXACT

    @staticmethod
    def from_json(data: dict) -> "HTTPHeaderMatch":
        return HTTPHeaderMatch(
            type=data.get("type", "Exact"), name=data["name"], value=data["value"]
        )


@dataclass
class Webhook:
    host: Host
    forwardHeaders: Optional[List[HTTPHeaderMatch]] = field(default_factory=list)

    @staticmethod
    def from_json(data: dict) -> "Webhook":
        host = Host.from_json(data["host"])
        forward_headers = [
            HTTPHeaderMatch.from_json(h) for h in data.get("forwardHeaders", [])
        ]
        return Webhook(host=host, forwardHeaders=forward_headers)


@dataclass
class Moderation:
    model: Optional[str] = None
    auth_token: Optional[SingleAuthToken] = None

    @staticmethod
    def from_json(data: dict) -> "Moderation":
        return Moderation(
            model=data.get("model"),
            auth_token=SingleAuthToken(**data.get("auth_token", {})),
        )


@dataclass
class PromptguardRequest:
    customResponse: Optional[CustomResponse] = None
    regex: Optional[Regex] = None
    webhook: Optional[Webhook] = None
    moderation: Optional[Moderation] = None


@dataclass
class PromptguardResponse:
    regex: Optional[Regex] = None
    webhook: Optional[Webhook] = None


@dataclass
class AIPromptGuard:
    request: Optional[PromptguardRequest] = None
    response: Optional[PromptguardResponse] = None


def from_json(data: str) -> AIPromptGuard:
    json_data = json.loads(data)
    return AIPromptGuard(
        request=PromptguardRequest(
            customResponse=CustomResponse(
                **json_data.get("request", {}).get("customResponse", {})
            ),
            regex=Regex(
                matches=[
                    RegexMatch(**m)
                    for m in json_data.get("request", {})
                    .get("regex", {})
                    .get("matches", [])
                ],
                builtins=[
                    BuiltIn(b)
                    for b in json_data.get("request", {})
                    .get("regex", {})
                    .get("builtins", [])
                ],
                action=json_data.get("request", {})
                .get("regex", {})
                .get("action", Action.MASK),
            ),
            webhook=Webhook.from_json(json_data.get("request", {}).get("webhook", {})),
            moderation=Moderation(
                model=json_data.get("request", {}).get("moderation", {}).get("model"),
                auth_token=SingleAuthToken(
                    **json_data.get("request", {})
                    .get("moderation", {})
                    .get("auth_token", {})
                ),
            ),
        ),
        response=PromptguardResponse(
            regex=Regex(
                matches=[
                    RegexMatch(**m)
                    for m in json_data.get("response", {})
                    .get("regex", {})
                    .get("matches", [])
                ],
                builtins=[
                    BuiltIn(b)
                    for b in json_data.get("response", {})
                    .get("regex", {})
                    .get("builtins", [])
                ],
                action=json_data.get("response", {})
                .get("regex", {})
                .get("action", Action.MASK),
            ),
            webhook=Webhook.from_json(json_data.get("response", {}).get("webhook", {})),
        ),
    )
