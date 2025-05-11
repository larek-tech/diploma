from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class VectorSearchRequest(_message.Message):
    __slots__ = ("query", "sourceIds", "topK", "threshold", "useQuestions")
    QUERY_FIELD_NUMBER: _ClassVar[int]
    SOURCEIDS_FIELD_NUMBER: _ClassVar[int]
    TOPK_FIELD_NUMBER: _ClassVar[int]
    THRESHOLD_FIELD_NUMBER: _ClassVar[int]
    USEQUESTIONS_FIELD_NUMBER: _ClassVar[int]
    query: str
    sourceIds: _containers.RepeatedScalarFieldContainer[str]
    topK: int
    threshold: float
    useQuestions: bool
    def __init__(self, query: _Optional[str] = ..., sourceIds: _Optional[_Iterable[str]] = ..., topK: _Optional[int] = ..., threshold: _Optional[float] = ..., useQuestions: bool = ...) -> None: ...

class DocumentChunk(_message.Message):
    __slots__ = ("id", "index", "content", "metadata", "similarity")
    ID_FIELD_NUMBER: _ClassVar[int]
    INDEX_FIELD_NUMBER: _ClassVar[int]
    CONTENT_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    SIMILARITY_FIELD_NUMBER: _ClassVar[int]
    id: str
    index: int
    content: str
    metadata: bytes
    similarity: float
    def __init__(self, id: _Optional[str] = ..., index: _Optional[int] = ..., content: _Optional[str] = ..., metadata: _Optional[bytes] = ..., similarity: _Optional[float] = ...) -> None: ...

class VectorSearchResponse(_message.Message):
    __slots__ = ("chunks",)
    CHUNKS_FIELD_NUMBER: _ClassVar[int]
    chunks: _containers.RepeatedCompositeFieldContainer[DocumentChunk]
    def __init__(
        self,
        chunks: _Optional[_Iterable[_Union[DocumentChunk, _Mapping]]] = ...,
    ) -> None: ...

class GetDocumentsIn(_message.Message):
    __slots__ = ("sourceId", "size", "page")
    SOURCEID_FIELD_NUMBER: _ClassVar[int]
    SIZE_FIELD_NUMBER: _ClassVar[int]
    PAGE_FIELD_NUMBER: _ClassVar[int]
    sourceId: str
    size: int
    page: int
    def __init__(
        self,
        sourceId: _Optional[str] = ...,
        size: _Optional[int] = ...,
        page: _Optional[int] = ...,
    ) -> None: ...

class Document(_message.Message):
    __slots__ = ("id", "sourceId", "name", "content", "metadata")
    ID_FIELD_NUMBER: _ClassVar[int]
    SOURCEID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    CONTENT_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    id: str
    sourceId: str
    name: str
    content: str
    metadata: str
    def __init__(
        self,
        id: _Optional[str] = ...,
        sourceId: _Optional[str] = ...,
        name: _Optional[str] = ...,
        content: _Optional[str] = ...,
        metadata: _Optional[str] = ...,
    ) -> None: ...

class GetDocumentsOut(_message.Message):
    __slots__ = ("size", "page", "total", "documents")
    SIZE_FIELD_NUMBER: _ClassVar[int]
    PAGE_FIELD_NUMBER: _ClassVar[int]
    TOTAL_FIELD_NUMBER: _ClassVar[int]
    DOCUMENTS_FIELD_NUMBER: _ClassVar[int]
    size: int
    page: int
    total: int
    documents: _containers.RepeatedCompositeFieldContainer[Document]
    def __init__(
        self,
        size: _Optional[int] = ...,
        page: _Optional[int] = ...,
        total: _Optional[int] = ...,
        documents: _Optional[_Iterable[_Union[Document, _Mapping]]] = ...,
    ) -> None: ...
