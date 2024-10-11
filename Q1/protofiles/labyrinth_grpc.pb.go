// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: protofiles/labyrinth.proto

package protofiles

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	LabyrinthService_GetLabyrinthInfo_FullMethodName = "/labyrinth.LabyrinthService/GetLabyrinthInfo"
	LabyrinthService_Revelio_FullMethodName          = "/labyrinth.LabyrinthService/Revelio"
	LabyrinthService_Bombarda_FullMethodName         = "/labyrinth.LabyrinthService/Bombarda"
)

// LabyrinthServiceClient is the client API for LabyrinthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LabyrinthServiceClient interface {
	GetLabyrinthInfo(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*LabyrinthInfoResponse, error)
	Revelio(ctx context.Context, in *RevelioRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Position], error)
	Bombarda(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[BombardaRequest, EmptyMessage], error)
}

type labyrinthServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLabyrinthServiceClient(cc grpc.ClientConnInterface) LabyrinthServiceClient {
	return &labyrinthServiceClient{cc}
}

func (c *labyrinthServiceClient) GetLabyrinthInfo(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*LabyrinthInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LabyrinthInfoResponse)
	err := c.cc.Invoke(ctx, LabyrinthService_GetLabyrinthInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *labyrinthServiceClient) Revelio(ctx context.Context, in *RevelioRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Position], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &LabyrinthService_ServiceDesc.Streams[0], LabyrinthService_Revelio_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[RevelioRequest, Position]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type LabyrinthService_RevelioClient = grpc.ServerStreamingClient[Position]

func (c *labyrinthServiceClient) Bombarda(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[BombardaRequest, EmptyMessage], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &LabyrinthService_ServiceDesc.Streams[1], LabyrinthService_Bombarda_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[BombardaRequest, EmptyMessage]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type LabyrinthService_BombardaClient = grpc.ClientStreamingClient[BombardaRequest, EmptyMessage]

// LabyrinthServiceServer is the server API for LabyrinthService service.
// All implementations must embed UnimplementedLabyrinthServiceServer
// for forward compatibility.
type LabyrinthServiceServer interface {
	GetLabyrinthInfo(context.Context, *EmptyMessage) (*LabyrinthInfoResponse, error)
	Revelio(*RevelioRequest, grpc.ServerStreamingServer[Position]) error
	Bombarda(grpc.ClientStreamingServer[BombardaRequest, EmptyMessage]) error
	mustEmbedUnimplementedLabyrinthServiceServer()
}

// UnimplementedLabyrinthServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedLabyrinthServiceServer struct{}

func (UnimplementedLabyrinthServiceServer) GetLabyrinthInfo(context.Context, *EmptyMessage) (*LabyrinthInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLabyrinthInfo not implemented")
}
func (UnimplementedLabyrinthServiceServer) Revelio(*RevelioRequest, grpc.ServerStreamingServer[Position]) error {
	return status.Errorf(codes.Unimplemented, "method Revelio not implemented")
}
func (UnimplementedLabyrinthServiceServer) Bombarda(grpc.ClientStreamingServer[BombardaRequest, EmptyMessage]) error {
	return status.Errorf(codes.Unimplemented, "method Bombarda not implemented")
}
func (UnimplementedLabyrinthServiceServer) mustEmbedUnimplementedLabyrinthServiceServer() {}
func (UnimplementedLabyrinthServiceServer) testEmbeddedByValue()                          {}

// UnsafeLabyrinthServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LabyrinthServiceServer will
// result in compilation errors.
type UnsafeLabyrinthServiceServer interface {
	mustEmbedUnimplementedLabyrinthServiceServer()
}

func RegisterLabyrinthServiceServer(s grpc.ServiceRegistrar, srv LabyrinthServiceServer) {
	// If the following call pancis, it indicates UnimplementedLabyrinthServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&LabyrinthService_ServiceDesc, srv)
}

func _LabyrinthService_GetLabyrinthInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LabyrinthServiceServer).GetLabyrinthInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LabyrinthService_GetLabyrinthInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LabyrinthServiceServer).GetLabyrinthInfo(ctx, req.(*EmptyMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _LabyrinthService_Revelio_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RevelioRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(LabyrinthServiceServer).Revelio(m, &grpc.GenericServerStream[RevelioRequest, Position]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type LabyrinthService_RevelioServer = grpc.ServerStreamingServer[Position]

func _LabyrinthService_Bombarda_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(LabyrinthServiceServer).Bombarda(&grpc.GenericServerStream[BombardaRequest, EmptyMessage]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type LabyrinthService_BombardaServer = grpc.ClientStreamingServer[BombardaRequest, EmptyMessage]

// LabyrinthService_ServiceDesc is the grpc.ServiceDesc for LabyrinthService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LabyrinthService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "labyrinth.LabyrinthService",
	HandlerType: (*LabyrinthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetLabyrinthInfo",
			Handler:    _LabyrinthService_GetLabyrinthInfo_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Revelio",
			Handler:       _LabyrinthService_Revelio_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Bombarda",
			Handler:       _LabyrinthService_Bombarda_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "protofiles/labyrinth.proto",
}