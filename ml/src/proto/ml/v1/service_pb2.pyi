from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import (
    ClassVar as _ClassVar,
    Iterable as _Iterable,
    Mapping as _Mapping,
    Optional as _Optional,
    Union as _Union,
)

DESCRIPTOR: _descriptor.FileDescriptor

class ScenarioType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    UNDEFINED_SCENARIO: _ClassVar[ScenarioType]
    SYSTEM_PROMPT: _ClassVar[ScenarioType]

UNDEFINED_SCENARIO: ScenarioType
SYSTEM_PROMPT: ScenarioType

class Scenario(_message.Message):
    __slots__ = ("customType", "content")
    CUSTOMTYPE_FIELD_NUMBER: _ClassVar[int]
    CONTENT_FIELD_NUMBER: _ClassVar[int]
    customType: ScenarioType
    content: str
    def __init__(
        self,
        customType: _Optional[_Union[ScenarioType, str]] = ...,
        content: _Optional[str] = ...,
    ) -> None: ...

class Query(_message.Message):
    __slots__ = ("userId", "content")
    USERID_FIELD_NUMBER: _ClassVar[int]
    CONTENT_FIELD_NUMBER: _ClassVar[int]
    userId: int
    content: str
    def __init__(
        self, userId: _Optional[int] = ..., content: _Optional[str] = ...
    ) -> None: ...

class ProcessQueryRequest(_message.Message):
    __slots__ = ("query", "scenario", "documentIds")
    QUERY_FIELD_NUMBER: _ClassVar[int]
    SCENARIO_FIELD_NUMBER: _ClassVar[int]
    DOCUMENTIDS_FIELD_NUMBER: _ClassVar[int]
    query: Query
    scenario: Scenario
    documentIds: _containers.RepeatedScalarFieldContainer[str]
    def __init__(
        self,
        query: _Optional[_Union[Query, _Mapping]] = ...,
        scenario: _Optional[_Union[Scenario, _Mapping]] = ...,
        documentIds: _Optional[_Iterable[str]] = ...,
    ) -> None: ...

class Chunk(_message.Message):
    __slots__ = ("content",)
    CONTENT_FIELD_NUMBER: _ClassVar[int]
    content: str
    def __init__(self, content: _Optional[str] = ...) -> None: ...

class ProcessQueryResponse(_message.Message):
    __slots__ = ("chunk", "documentIds")
    CHUNK_FIELD_NUMBER: _ClassVar[int]
    DOCUMENTIDS_FIELD_NUMBER: _ClassVar[int]
    chunk: Chunk
    documentIds: _containers.RepeatedScalarFieldContainer[str]
    def __init__(
        self,
        chunk: _Optional[_Union[Chunk, _Mapping]] = ...,
        documentIds: _Optional[_Iterable[str]] = ...,
    ) -> None: ...
