// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: domain/v1/role_service.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_domain_v1_role_service_proto protoreflect.FileDescriptor

const file_domain_v1_role_service_proto_rawDesc = "" +
	"\n" +
	"\x1cdomain/v1/role_service.proto\x12\tdomain.v1\x1a\x1adomain/v1/role_model.proto\x1a\x1bgoogle/protobuf/empty.proto2\x9e\x03\n" +
	"\vRoleService\x12=\n" +
	"\n" +
	"CreateRole\x12\x1c.domain.v1.CreateRoleRequest\x1a\x0f.domain.v1.Role\"\x00\x127\n" +
	"\aGetRole\x12\x19.domain.v1.GetRoleRequest\x1a\x0f.domain.v1.Role\"\x00\x12D\n" +
	"\n" +
	"DeleteRole\x12\x1c.domain.v1.DeleteRoleRequest\x1a\x16.google.protobuf.Empty\"\x00\x12H\n" +
	"\tListRoles\x12\x1b.domain.v1.ListRolesRequest\x1a\x1c.domain.v1.ListRolesResponse\"\x00\x12A\n" +
	"\aSetRole\x12\x1c.domain.v1.UpdateRoleRequest\x1a\x16.google.protobuf.Empty\"\x00\x12D\n" +
	"\n" +
	"RemoveRole\x12\x1c.domain.v1.UpdateRoleRequest\x1a\x16.google.protobuf.Empty\"\x00B\x14Z\x12internal/domain/pbb\x06proto3"

var file_domain_v1_role_service_proto_goTypes = []any{
	(*CreateRoleRequest)(nil), // 0: domain.v1.CreateRoleRequest
	(*GetRoleRequest)(nil),    // 1: domain.v1.GetRoleRequest
	(*DeleteRoleRequest)(nil), // 2: domain.v1.DeleteRoleRequest
	(*ListRolesRequest)(nil),  // 3: domain.v1.ListRolesRequest
	(*UpdateRoleRequest)(nil), // 4: domain.v1.UpdateRoleRequest
	(*Role)(nil),              // 5: domain.v1.Role
	(*emptypb.Empty)(nil),     // 6: google.protobuf.Empty
	(*ListRolesResponse)(nil), // 7: domain.v1.ListRolesResponse
}
var file_domain_v1_role_service_proto_depIdxs = []int32{
	0, // 0: domain.v1.RoleService.CreateRole:input_type -> domain.v1.CreateRoleRequest
	1, // 1: domain.v1.RoleService.GetRole:input_type -> domain.v1.GetRoleRequest
	2, // 2: domain.v1.RoleService.DeleteRole:input_type -> domain.v1.DeleteRoleRequest
	3, // 3: domain.v1.RoleService.ListRoles:input_type -> domain.v1.ListRolesRequest
	4, // 4: domain.v1.RoleService.SetRole:input_type -> domain.v1.UpdateRoleRequest
	4, // 5: domain.v1.RoleService.RemoveRole:input_type -> domain.v1.UpdateRoleRequest
	5, // 6: domain.v1.RoleService.CreateRole:output_type -> domain.v1.Role
	5, // 7: domain.v1.RoleService.GetRole:output_type -> domain.v1.Role
	6, // 8: domain.v1.RoleService.DeleteRole:output_type -> google.protobuf.Empty
	7, // 9: domain.v1.RoleService.ListRoles:output_type -> domain.v1.ListRolesResponse
	6, // 10: domain.v1.RoleService.SetRole:output_type -> google.protobuf.Empty
	6, // 11: domain.v1.RoleService.RemoveRole:output_type -> google.protobuf.Empty
	6, // [6:12] is the sub-list for method output_type
	0, // [0:6] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_domain_v1_role_service_proto_init() }
func file_domain_v1_role_service_proto_init() {
	if File_domain_v1_role_service_proto != nil {
		return
	}
	file_domain_v1_role_model_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_domain_v1_role_service_proto_rawDesc), len(file_domain_v1_role_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_domain_v1_role_service_proto_goTypes,
		DependencyIndexes: file_domain_v1_role_service_proto_depIdxs,
	}.Build()
	File_domain_v1_role_service_proto = out.File
	file_domain_v1_role_service_proto_goTypes = nil
	file_domain_v1_role_service_proto_depIdxs = nil
}
