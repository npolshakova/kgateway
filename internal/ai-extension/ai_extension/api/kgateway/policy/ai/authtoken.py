import json
from enum import Enum
from dataclasses import dataclass
from typing import Optional


class SingleAuthTokenKind(Enum):
    INLINE = "Inline"
    SECRET_REF = "SecretRef"
    PASSTHROUGH = "Passthrough"


@dataclass
class LocalObjectReference:
    name: Optional[str] = None

    @staticmethod
    def from_json(data: dict) -> "LocalObjectReference":
        return LocalObjectReference(**data)


@dataclass
class SingleAuthToken:
    def __init__(
        self,
        kind: SingleAuthTokenKind,
        inline: Optional[str] = None,
        secret_ref: Optional[LocalObjectReference] = None,
    ):
        self.kind = kind
        self.inline = inline
        self.secret_ref = secret_ref

    def __repr__(self):
        return f"SingleAuthToken(kind={self.kind}, inline={self.inline}, secret_ref={self.secret_ref})"


def from_json(data: str) -> SingleAuthToken:
    json_data = json.loads(data)
    return SingleAuthToken(
        kind=SingleAuthTokenKind(json_data["kind"]),
        inline=json_data.get("inline"),
        secret_ref=LocalObjectReference(**json_data.get("secret_ref", {}))
        if json_data.get("secret_ref")
        else None,
    )
