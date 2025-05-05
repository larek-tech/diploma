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

class MultiQuery(_message.Message):
    __slots__ = ("useMultiquery", "nQueryes", "queryModelName")
    USEMULTIQUERY_FIELD_NUMBER: _ClassVar[int]
    NQUERYES_FIELD_NUMBER: _ClassVar[int]
    QUERYMODELNAME_FIELD_NUMBER: _ClassVar[int]
    useMultiquery: bool
    nQueryes: int
    queryModelName: str
    def __init__(
        self,
        useMultiquery: bool = ...,
        nQueryes: _Optional[int] = ...,
        queryModelName: _Optional[str] = ...,
    ) -> None: ...

class Reranker(_message.Message):
    __slots__ = ("useRerank", "rerankerModel", "rerankerMaxLenght", "topK")
    USERERANK_FIELD_NUMBER: _ClassVar[int]
    RERANKERMODEL_FIELD_NUMBER: _ClassVar[int]
    RERANKERMAXLENGHT_FIELD_NUMBER: _ClassVar[int]
    TOPK_FIELD_NUMBER: _ClassVar[int]
    useRerank: bool
    rerankerModel: str
    rerankerMaxLenght: int
    topK: int
    def __init__(
        self,
        useRerank: bool = ...,
        rerankerModel: _Optional[str] = ...,
        rerankerMaxLenght: _Optional[int] = ...,
        topK: _Optional[int] = ...,
    ) -> None: ...

class LLMMdel(_message.Message):
    __slots__ = ("modelName", "temprature", "topK", "topP")
    MODELNAME_FIELD_NUMBER: _ClassVar[int]
    TEMPRATURE_FIELD_NUMBER: _ClassVar[int]
    TOPK_FIELD_NUMBER: _ClassVar[int]
    TOPP_FIELD_NUMBER: _ClassVar[int]
    modelName: str
    temprature: float
    topK: int
    topP: float
    def __init__(
        self,
        modelName: _Optional[str] = ...,
        temprature: _Optional[float] = ...,
        topK: _Optional[int] = ...,
        topP: _Optional[float] = ...,
    ) -> None: ...

class VectorSearch(_message.Message):
    __slots__ = ("topN", "threshold", "searchByQuery")
    TOPN_FIELD_NUMBER: _ClassVar[int]
    THRESHOLD_FIELD_NUMBER: _ClassVar[int]
    SEARCHBYQUERY_FIELD_NUMBER: _ClassVar[int]
    topN: int
    threshold: float
    searchByQuery: bool
    def __init__(
        self,
        topN: _Optional[int] = ...,
        threshold: _Optional[float] = ...,
        searchByQuery: bool = ...,
    ) -> None: ...

class Scenario(_message.Message):
    __slots__ = (
        "customType",
        "content",
        "multiQuery",
        "reranker",
        "vectorSearch",
        "model",
    )
    CUSTOMTYPE_FIELD_NUMBER: _ClassVar[int]
    CONTENT_FIELD_NUMBER: _ClassVar[int]
    MULTIQUERY_FIELD_NUMBER: _ClassVar[int]
    RERANKER_FIELD_NUMBER: _ClassVar[int]
    VECTORSEARCH_FIELD_NUMBER: _ClassVar[int]
    MODEL_FIELD_NUMBER: _ClassVar[int]
    customType: ScenarioType
    content: str
    multiQuery: MultiQuery
    reranker: Reranker
    vectorSearch: VectorSearch
    model: LLMMdel
    def __init__(
        self,
        customType: _Optional[_Union[ScenarioType, str]] = ...,
        content: _Optional[str] = ...,
        multiQuery: _Optional[_Union[MultiQuery, _Mapping]] = ...,
        reranker: _Optional[_Union[Reranker, _Mapping]] = ...,
        vectorSearch: _Optional[_Union[VectorSearch, _Mapping]] = ...,
        model: _Optional[_Union[LLMMdel, _Mapping]] = ...,
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
    __slots__ = ("query", "scenario", "sourceIds")
    QUERY_FIELD_NUMBER: _ClassVar[int]
    SCENARIO_FIELD_NUMBER: _ClassVar[int]
    SOURCEIDS_FIELD_NUMBER: _ClassVar[int]
    query: Query
    scenario: Scenario
    sourceIds: _containers.RepeatedScalarFieldContainer[str]
    def __init__(
        self,
        query: _Optional[_Union[Query, _Mapping]] = ...,
        scenario: _Optional[_Union[Scenario, _Mapping]] = ...,
        sourceIds: _Optional[_Iterable[str]] = ...,
    ) -> None: ...

class Chunk(_message.Message):
    __slots__ = ("content",)
    CONTENT_FIELD_NUMBER: _ClassVar[int]
    content: str
    def __init__(self, content: _Optional[str] = ...) -> None: ...

class ProcessQueryResponse(_message.Message):
    __slots__ = ("chunk", "sourceIds")
    CHUNK_FIELD_NUMBER: _ClassVar[int]
    SOURCEIDS_FIELD_NUMBER: _ClassVar[int]
    chunk: Chunk
    sourceIds: _containers.RepeatedScalarFieldContainer[str]
    def __init__(
        self,
        chunk: _Optional[_Union[Chunk, _Mapping]] = ...,
        sourceIds: _Optional[_Iterable[str]] = ...,
    ) -> None: ...
