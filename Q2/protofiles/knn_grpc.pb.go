// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: protofiles/knn.proto

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
	KNNService_FindKNearestNeighbors_FullMethodName = "/knn.KNNService/FindKNearestNeighbors"
)

// KNNServiceClient is the client API for KNNService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KNNServiceClient interface {
	FindKNearestNeighbors(ctx context.Context, in *KNNRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Neighbor], error)
}

type kNNServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewKNNServiceClient(cc grpc.ClientConnInterface) KNNServiceClient {
	return &kNNServiceClient{cc}
}

func (c *kNNServiceClient) FindKNearestNeighbors(ctx context.Context, in *KNNRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Neighbor], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &KNNService_ServiceDesc.Streams[0], KNNService_FindKNearestNeighbors_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[KNNRequest, Neighbor]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type KNNService_FindKNearestNeighborsClient = grpc.ServerStreamingClient[Neighbor]

// KNNServiceServer is the server API for KNNService service.
// All implementations must embed UnimplementedKNNServiceServer
// for forward compatibility.
type KNNServiceServer interface {
	FindKNearestNeighbors(*KNNRequest, grpc.ServerStreamingServer[Neighbor]) error
	mustEmbedUnimplementedKNNServiceServer()
}

// UnimplementedKNNServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedKNNServiceServer struct{}

func (UnimplementedKNNServiceServer) FindKNearestNeighbors(*KNNRequest, grpc.ServerStreamingServer[Neighbor]) error {
	return status.Errorf(codes.Unimplemented, "method FindKNearestNeighbors not implemented")
}
func (UnimplementedKNNServiceServer) mustEmbedUnimplementedKNNServiceServer() {}
func (UnimplementedKNNServiceServer) testEmbeddedByValue()                    {}

// UnsafeKNNServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KNNServiceServer will
// result in compilation errors.
type UnsafeKNNServiceServer interface {
	mustEmbedUnimplementedKNNServiceServer()
}

func RegisterKNNServiceServer(s grpc.ServiceRegistrar, srv KNNServiceServer) {
	// If the following call pancis, it indicates UnimplementedKNNServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&KNNService_ServiceDesc, srv)
}

func _KNNService_FindKNearestNeighbors_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(KNNRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(KNNServiceServer).FindKNearestNeighbors(m, &grpc.GenericServerStream[KNNRequest, Neighbor]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type KNNService_FindKNearestNeighborsServer = grpc.ServerStreamingServer[Neighbor]

// KNNService_ServiceDesc is the grpc.ServiceDesc for KNNService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KNNService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "knn.KNNService",
	HandlerType: (*KNNServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "FindKNearestNeighbors",
			Handler:       _KNNService_FindKNearestNeighbors_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "protofiles/knn.proto",
}
