# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: proto/ml/v1/service.proto
# Protobuf Python Version: 5.29.0
"""Generated protocol buffer code."""

from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder

_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC, 5, 29, 0, "", "proto/ml/v1/service.proto"
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(
    b'\n\x19proto/ml/v1/service.proto\x12\x02ml"A\n\x08Scenario\x12$\n\ncustomType\x18\x01 \x01(\x0e\x32\x10.ml.ScenarioType\x12\x0f\n\x07\x63ontent\x18\x02 \x01(\t"(\n\x05Query\x12\x0e\n\x06userId\x18\x01 \x01(\x03\x12\x0f\n\x07\x63ontent\x18\x02 \x01(\t"d\n\x13ProcessQueryRequest\x12\x18\n\x05query\x18\x01 \x01(\x0b\x32\t.ml.Query\x12\x1e\n\x08scenario\x18\x02 \x01(\x0b\x32\x0c.ml.Scenario\x12\x13\n\x0b\x64ocumentIds\x18\x03 \x03(\t"\x18\n\x05\x43hunk\x12\x0f\n\x07\x63ontent\x18\x01 \x01(\t"E\n\x14ProcessQueryResponse\x12\x18\n\x05\x63hunk\x18\x01 \x01(\x0b\x32\t.ml.Chunk\x12\x13\n\x0b\x64ocumentIds\x18\x02 \x03(\t*9\n\x0cScenarioType\x12\x16\n\x12UNDEFINED_SCENARIO\x10\x00\x12\x11\n\rSYSTEM_PROMPT\x10\x01\x32R\n\tMLService\x12\x45\n\x0cProcessQuery\x12\x17.ml.ProcessQueryRequest\x1a\x18.ml.ProcessQueryResponse"\x00\x30\x01\x42\x10Z\x0einternal/ml/pbb\x06proto3'
)

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(
    DESCRIPTOR, "proto.ml.v1.service_pb2", _globals
)
if not _descriptor._USE_C_DESCRIPTORS:
    _globals["DESCRIPTOR"]._loaded_options = None
    _globals["DESCRIPTOR"]._serialized_options = b"Z\016internal/ml/pb"
    _globals["_SCENARIOTYPE"]._serialized_start = 341
    _globals["_SCENARIOTYPE"]._serialized_end = 398
    _globals["_SCENARIO"]._serialized_start = 33
    _globals["_SCENARIO"]._serialized_end = 98
    _globals["_QUERY"]._serialized_start = 100
    _globals["_QUERY"]._serialized_end = 140
    _globals["_PROCESSQUERYREQUEST"]._serialized_start = 142
    _globals["_PROCESSQUERYREQUEST"]._serialized_end = 242
    _globals["_CHUNK"]._serialized_start = 244
    _globals["_CHUNK"]._serialized_end = 268
    _globals["_PROCESSQUERYRESPONSE"]._serialized_start = 270
    _globals["_PROCESSQUERYRESPONSE"]._serialized_end = 339
    _globals["_MLSERVICE"]._serialized_start = 400
    _globals["_MLSERVICE"]._serialized_end = 482
# @@protoc_insertion_point(module_scope)
