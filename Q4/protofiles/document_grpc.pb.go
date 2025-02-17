// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: protofiles/document.proto

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
	CollaborativeDocumentService_SyncDocumentChanges_FullMethodName = "/collaborative_document.CollaborativeDocumentService/SyncDocumentChanges"
	CollaborativeDocumentService_StreamDocumentLogs_FullMethodName  = "/collaborative_document.CollaborativeDocumentService/StreamDocumentLogs"
)

// CollaborativeDocumentServiceClient is the client API for CollaborativeDocumentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CollaborativeDocumentServiceClient interface {
	SyncDocumentChanges(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[DocumentChange, DocumentChange], error)
	StreamDocumentLogs(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DocumentChange], error)
}

type collaborativeDocumentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCollaborativeDocumentServiceClient(cc grpc.ClientConnInterface) CollaborativeDocumentServiceClient {
	return &collaborativeDocumentServiceClient{cc}
}

func (c *collaborativeDocumentServiceClient) SyncDocumentChanges(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[DocumentChange, DocumentChange], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &CollaborativeDocumentService_ServiceDesc.Streams[0], CollaborativeDocumentService_SyncDocumentChanges_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[DocumentChange, DocumentChange]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type CollaborativeDocumentService_SyncDocumentChangesClient = grpc.BidiStreamingClient[DocumentChange, DocumentChange]

func (c *collaborativeDocumentServiceClient) StreamDocumentLogs(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DocumentChange], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &CollaborativeDocumentService_ServiceDesc.Streams[1], CollaborativeDocumentService_StreamDocumentLogs_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[EmptyMessage, DocumentChange]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type CollaborativeDocumentService_StreamDocumentLogsClient = grpc.ServerStreamingClient[DocumentChange]

// CollaborativeDocumentServiceServer is the server API for CollaborativeDocumentService service.
// All implementations must embed UnimplementedCollaborativeDocumentServiceServer
// for forward compatibility.
type CollaborativeDocumentServiceServer interface {
	SyncDocumentChanges(grpc.BidiStreamingServer[DocumentChange, DocumentChange]) error
	StreamDocumentLogs(*EmptyMessage, grpc.ServerStreamingServer[DocumentChange]) error
	mustEmbedUnimplementedCollaborativeDocumentServiceServer()
}

// UnimplementedCollaborativeDocumentServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCollaborativeDocumentServiceServer struct{}

func (UnimplementedCollaborativeDocumentServiceServer) SyncDocumentChanges(grpc.BidiStreamingServer[DocumentChange, DocumentChange]) error {
	return status.Errorf(codes.Unimplemented, "method SyncDocumentChanges not implemented")
}
func (UnimplementedCollaborativeDocumentServiceServer) StreamDocumentLogs(*EmptyMessage, grpc.ServerStreamingServer[DocumentChange]) error {
	return status.Errorf(codes.Unimplemented, "method StreamDocumentLogs not implemented")
}
func (UnimplementedCollaborativeDocumentServiceServer) mustEmbedUnimplementedCollaborativeDocumentServiceServer() {
}
func (UnimplementedCollaborativeDocumentServiceServer) testEmbeddedByValue() {}

// UnsafeCollaborativeDocumentServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CollaborativeDocumentServiceServer will
// result in compilation errors.
type UnsafeCollaborativeDocumentServiceServer interface {
	mustEmbedUnimplementedCollaborativeDocumentServiceServer()
}

func RegisterCollaborativeDocumentServiceServer(s grpc.ServiceRegistrar, srv CollaborativeDocumentServiceServer) {
	// If the following call pancis, it indicates UnimplementedCollaborativeDocumentServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CollaborativeDocumentService_ServiceDesc, srv)
}

func _CollaborativeDocumentService_SyncDocumentChanges_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CollaborativeDocumentServiceServer).SyncDocumentChanges(&grpc.GenericServerStream[DocumentChange, DocumentChange]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type CollaborativeDocumentService_SyncDocumentChangesServer = grpc.BidiStreamingServer[DocumentChange, DocumentChange]

func _CollaborativeDocumentService_StreamDocumentLogs_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EmptyMessage)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CollaborativeDocumentServiceServer).StreamDocumentLogs(m, &grpc.GenericServerStream[EmptyMessage, DocumentChange]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type CollaborativeDocumentService_StreamDocumentLogsServer = grpc.ServerStreamingServer[DocumentChange]

// CollaborativeDocumentService_ServiceDesc is the grpc.ServiceDesc for CollaborativeDocumentService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CollaborativeDocumentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "collaborative_document.CollaborativeDocumentService",
	HandlerType: (*CollaborativeDocumentServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SyncDocumentChanges",
			Handler:       _CollaborativeDocumentService_SyncDocumentChanges_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamDocumentLogs",
			Handler:       _CollaborativeDocumentService_StreamDocumentLogs_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "protofiles/document.proto",
}
