from google.protobuf import timestamp_pb2 as _timestamp_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MultiQuery(_message.Message):
    __slots__ = ("useMultiquery", "nQueries", "queryModelName")
    USEMULTIQUERY_FIELD_NUMBER: _ClassVar[int]
    NQUERIES_FIELD_NUMBER: _ClassVar[int]
    QUERYMODELNAME_FIELD_NUMBER: _ClassVar[int]
    useMultiquery: bool
    nQueries: int
    queryModelName: str
    def __init__(self, useMultiquery: bool = ..., nQueries: _Optional[int] = ..., queryModelName: _Optional[str] = ...) -> None: ...

class Reranker(_message.Message):
    __slots__ = ("useRerank", "rerankerModel", "rerankerMaxLength", "topK")
    USERERANK_FIELD_NUMBER: _ClassVar[int]
    RERANKERMODEL_FIELD_NUMBER: _ClassVar[int]
    RERANKERMAXLENGTH_FIELD_NUMBER: _ClassVar[int]
    TOPK_FIELD_NUMBER: _ClassVar[int]
    useRerank: bool
    rerankerModel: str
    rerankerMaxLength: int
    topK: int
    def __init__(self, useRerank: bool = ..., rerankerModel: _Optional[str] = ..., rerankerMaxLength: _Optional[int] = ..., topK: _Optional[int] = ...) -> None: ...

class LlmModel(_message.Message):
    __slots__ = ("modelName", "temperature", "topK", "topP", "systemPrompt")
    MODELNAME_FIELD_NUMBER: _ClassVar[int]
    TEMPERATURE_FIELD_NUMBER: _ClassVar[int]
    TOPK_FIELD_NUMBER: _ClassVar[int]
    TOPP_FIELD_NUMBER: _ClassVar[int]
    SYSTEMPROMPT_FIELD_NUMBER: _ClassVar[int]
    modelName: str
    temperature: float
    topK: int
    topP: float
    systemPrompt: str
    def __init__(self, modelName: _Optional[str] = ..., temperature: _Optional[float] = ..., topK: _Optional[int] = ..., topP: _Optional[float] = ..., systemPrompt: _Optional[str] = ...) -> None: ...

class VectorSearch(_message.Message):
    __slots__ = ("topN", "threshold", "searchByQuery")
    TOPN_FIELD_NUMBER: _ClassVar[int]
    THRESHOLD_FIELD_NUMBER: _ClassVar[int]
    SEARCHBYQUERY_FIELD_NUMBER: _ClassVar[int]
    topN: int
    threshold: float
    searchByQuery: bool
    def __init__(self, topN: _Optional[int] = ..., threshold: _Optional[float] = ..., searchByQuery: bool = ...) -> None: ...

class Scenario(_message.Message):
    __slots__ = ("id", "multiQuery", "reranker", "vectorSearch", "model", "createdAt", "updatedAt", "title", "domainId", "contextSize")
    ID_FIELD_NUMBER: _ClassVar[int]
    MULTIQUERY_FIELD_NUMBER: _ClassVar[int]
    RERANKER_FIELD_NUMBER: _ClassVar[int]
    VECTORSEARCH_FIELD_NUMBER: _ClassVar[int]
    MODEL_FIELD_NUMBER: _ClassVar[int]
    CREATEDAT_FIELD_NUMBER: _ClassVar[int]
    UPDATEDAT_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DOMAINID_FIELD_NUMBER: _ClassVar[int]
    CONTEXTSIZE_FIELD_NUMBER: _ClassVar[int]
    id: int
    multiQuery: MultiQuery
    reranker: Reranker
    vectorSearch: VectorSearch
    model: LlmModel
    createdAt: _timestamp_pb2.Timestamp
    updatedAt: _timestamp_pb2.Timestamp
    title: str
    domainId: int
    contextSize: int
    def __init__(self, id: _Optional[int] = ..., multiQuery: _Optional[_Union[MultiQuery, _Mapping]] = ..., reranker: _Optional[_Union[Reranker, _Mapping]] = ..., vectorSearch: _Optional[_Union[VectorSearch, _Mapping]] = ..., model: _Optional[_Union[LlmModel, _Mapping]] = ..., createdAt: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., updatedAt: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., title: _Optional[str] = ..., domainId: _Optional[int] = ..., contextSize: _Optional[int] = ...) -> None: ...

class Query(_message.Message):
    __slots__ = ("id", "userId", "content")
    ID_FIELD_NUMBER: _ClassVar[int]
    USERID_FIELD_NUMBER: _ClassVar[int]
    CONTENT_FIELD_NUMBER: _ClassVar[int]
    id: int
    userId: int
    content: str
    def __init__(self, id: _Optional[int] = ..., userId: _Optional[int] = ..., content: _Optional[str] = ...) -> None: ...

class ProcessQueryRequest(_message.Message):
    __slots__ = ("query", "scenario", "sourceIds")
    QUERY_FIELD_NUMBER: _ClassVar[int]
    SCENARIO_FIELD_NUMBER: _ClassVar[int]
    SOURCEIDS_FIELD_NUMBER: _ClassVar[int]
    query: Query
    scenario: Scenario
    sourceIds: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, query: _Optional[_Union[Query, _Mapping]] = ..., scenario: _Optional[_Union[Scenario, _Mapping]] = ..., sourceIds: _Optional[_Iterable[str]] = ...) -> None: ...

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
    def __init__(self, chunk: _Optional[_Union[Chunk, _Mapping]] = ..., sourceIds: _Optional[_Iterable[str]] = ...) -> None: ...

class ModelParams(_message.Message):
    __slots__ = ("multiQuery", "reranker", "vectorSearch", "model")
    MULTIQUERY_FIELD_NUMBER: _ClassVar[int]
    RERANKER_FIELD_NUMBER: _ClassVar[int]
    VECTORSEARCH_FIELD_NUMBER: _ClassVar[int]
    MODEL_FIELD_NUMBER: _ClassVar[int]
    multiQuery: MultiQuery
    reranker: Reranker
    vectorSearch: VectorSearch
    model: LlmModel
    def __init__(self, multiQuery: _Optional[_Union[MultiQuery, _Mapping]] = ..., reranker: _Optional[_Union[Reranker, _Mapping]] = ..., vectorSearch: _Optional[_Union[VectorSearch, _Mapping]] = ..., model: _Optional[_Union[LlmModel, _Mapping]] = ...) -> None: ...

class GetOptimalParamsRequest(_message.Message):
    __slots__ = ("sourceIds",)
    SOURCEIDS_FIELD_NUMBER: _ClassVar[int]
    sourceIds: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, sourceIds: _Optional[_Iterable[str]] = ...) -> None: ...

class ProcessFirstQueryRequest(_message.Message):
    __slots__ = ("query",)
    QUERY_FIELD_NUMBER: _ClassVar[int]
    query: str
    def __init__(self, query: _Optional[str] = ...) -> None: ...

class ProcessFirstQueryResponse(_message.Message):
    __slots__ = ("query",)
    QUERY_FIELD_NUMBER: _ClassVar[int]
    query: str
    def __init__(self, query: _Optional[str] = ...) -> None: ...
