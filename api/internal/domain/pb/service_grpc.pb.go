// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: ml/v1/service.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	MLService_ProcessQuery_FullMethodName      = "/pb.ml.MLService/ProcessQuery"
	MLService_GetDefaultParams_FullMethodName  = "/pb.ml.MLService/GetDefaultParams"
	MLService_GetOptimalParams_FullMethodName  = "/pb.ml.MLService/GetOptimalParams"
	MLService_ProcessFirstQuery_FullMethodName = "/pb.ml.MLService/ProcessFirstQuery"
)

// MLServiceClient is the client API for MLService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MLServiceClient interface {
	ProcessQuery(ctx context.Context, in *ProcessQueryRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ProcessQueryResponse], error)
	GetDefaultParams(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ModelParams, error)
	GetOptimalParams(ctx context.Context, in *GetOptimalParamsRequest, opts ...grpc.CallOption) (*ModelParams, error)
	ProcessFirstQuery(ctx context.Context, in *ProcessFirstQueryRequest, opts ...grpc.CallOption) (*ProcessFirstQueryResponse, error)
}

type mLServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMLServiceClient(cc grpc.ClientConnInterface) MLServiceClient {
	return &mLServiceClient{cc}
}

func (c *mLServiceClient) ProcessQuery(ctx context.Context, in *ProcessQueryRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ProcessQueryResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &MLService_ServiceDesc.Streams[0], MLService_ProcessQuery_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[ProcessQueryRequest, ProcessQueryResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type MLService_ProcessQueryClient = grpc.ServerStreamingClient[ProcessQueryResponse]

func (c *mLServiceClient) GetDefaultParams(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ModelParams, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ModelParams)
	err := c.cc.Invoke(ctx, MLService_GetDefaultParams_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mLServiceClient) GetOptimalParams(ctx context.Context, in *GetOptimalParamsRequest, opts ...grpc.CallOption) (*ModelParams, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ModelParams)
	err := c.cc.Invoke(ctx, MLService_GetOptimalParams_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mLServiceClient) ProcessFirstQuery(ctx context.Context, in *ProcessFirstQueryRequest, opts ...grpc.CallOption) (*ProcessFirstQueryResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProcessFirstQueryResponse)
	err := c.cc.Invoke(ctx, MLService_ProcessFirstQuery_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MLServiceServer is the server API for MLService service.
// All implementations must embed UnimplementedMLServiceServer
// for forward compatibility.
type MLServiceServer interface {
	ProcessQuery(*ProcessQueryRequest, grpc.ServerStreamingServer[ProcessQueryResponse]) error
	GetDefaultParams(context.Context, *emptypb.Empty) (*ModelParams, error)
	GetOptimalParams(context.Context, *GetOptimalParamsRequest) (*ModelParams, error)
	ProcessFirstQuery(context.Context, *ProcessFirstQueryRequest) (*ProcessFirstQueryResponse, error)
	mustEmbedUnimplementedMLServiceServer()
}

// UnimplementedMLServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMLServiceServer struct{}

func (UnimplementedMLServiceServer) ProcessQuery(*ProcessQueryRequest, grpc.ServerStreamingServer[ProcessQueryResponse]) error {
	return status.Errorf(codes.Unimplemented, "method ProcessQuery not implemented")
}
func (UnimplementedMLServiceServer) GetDefaultParams(context.Context, *emptypb.Empty) (*ModelParams, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDefaultParams not implemented")
}
func (UnimplementedMLServiceServer) GetOptimalParams(context.Context, *GetOptimalParamsRequest) (*ModelParams, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOptimalParams not implemented")
}
func (UnimplementedMLServiceServer) ProcessFirstQuery(context.Context, *ProcessFirstQueryRequest) (*ProcessFirstQueryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessFirstQuery not implemented")
}
func (UnimplementedMLServiceServer) mustEmbedUnimplementedMLServiceServer() {}
func (UnimplementedMLServiceServer) testEmbeddedByValue()                   {}

// UnsafeMLServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MLServiceServer will
// result in compilation errors.
type UnsafeMLServiceServer interface {
	mustEmbedUnimplementedMLServiceServer()
}

func RegisterMLServiceServer(s grpc.ServiceRegistrar, srv MLServiceServer) {
	// If the following call pancis, it indicates UnimplementedMLServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MLService_ServiceDesc, srv)
}

func _MLService_ProcessQuery_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ProcessQueryRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MLServiceServer).ProcessQuery(m, &grpc.GenericServerStream[ProcessQueryRequest, ProcessQueryResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type MLService_ProcessQueryServer = grpc.ServerStreamingServer[ProcessQueryResponse]

func _MLService_GetDefaultParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MLServiceServer).GetDefaultParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MLService_GetDefaultParams_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MLServiceServer).GetDefaultParams(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _MLService_GetOptimalParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOptimalParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MLServiceServer).GetOptimalParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MLService_GetOptimalParams_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MLServiceServer).GetOptimalParams(ctx, req.(*GetOptimalParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MLService_ProcessFirstQuery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessFirstQueryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MLServiceServer).ProcessFirstQuery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MLService_ProcessFirstQuery_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MLServiceServer).ProcessFirstQuery(ctx, req.(*ProcessFirstQueryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MLService_ServiceDesc is the grpc.ServiceDesc for MLService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MLService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.ml.MLService",
	HandlerType: (*MLServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDefaultParams",
			Handler:    _MLService_GetDefaultParams_Handler,
		},
		{
			MethodName: "GetOptimalParams",
			Handler:    _MLService_GetOptimalParams_Handler,
		},
		{
			MethodName: "ProcessFirstQuery",
			Handler:    _MLService_ProcessFirstQuery_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ProcessQuery",
			Handler:       _MLService_ProcessQuery_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "ml/v1/service.proto",
}
