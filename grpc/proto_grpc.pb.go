// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.19.6
// source: grpc/proto.proto

package proto

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
	ChittyChat_PublishMessage_FullMethodName = "/ChittyChat/PublishMessage"
	ChittyChat_Join_FullMethodName           = "/ChittyChat/Join"
	ChittyChat_Leave_FullMethodName          = "/ChittyChat/Leave"
)

// ChittyChatClient is the client API for ChittyChat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChittyChatClient interface {
	PublishMessage(ctx context.Context, in *Response, opts ...grpc.CallOption) (*Status, error)
	Join(ctx context.Context, in *Client, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Response], error)
	Leave(ctx context.Context, in *Client, opts ...grpc.CallOption) (*Empty, error)
}

type chittyChatClient struct {
	cc grpc.ClientConnInterface
}

func NewChittyChatClient(cc grpc.ClientConnInterface) ChittyChatClient {
	return &chittyChatClient{cc}
}

func (c *chittyChatClient) PublishMessage(ctx context.Context, in *Response, opts ...grpc.CallOption) (*Status, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Status)
	err := c.cc.Invoke(ctx, ChittyChat_PublishMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chittyChatClient) Join(ctx context.Context, in *Client, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Response], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &ChittyChat_ServiceDesc.Streams[0], ChittyChat_Join_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[Client, Response]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ChittyChat_JoinClient = grpc.ServerStreamingClient[Response]

func (c *chittyChatClient) Leave(ctx context.Context, in *Client, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, ChittyChat_Leave_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChittyChatServer is the server API for ChittyChat service.
// All implementations must embed UnimplementedChittyChatServer
// for forward compatibility.
type ChittyChatServer interface {
	PublishMessage(context.Context, *Response) (*Status, error)
	Join(*Client, grpc.ServerStreamingServer[Response]) error
	Leave(context.Context, *Client) (*Empty, error)
	mustEmbedUnimplementedChittyChatServer()
}

// UnimplementedChittyChatServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedChittyChatServer struct{}

func (UnimplementedChittyChatServer) PublishMessage(context.Context, *Response) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishMessage not implemented")
}
func (UnimplementedChittyChatServer) Join(*Client, grpc.ServerStreamingServer[Response]) error {
	return status.Errorf(codes.Unimplemented, "method Join not implemented")
}
func (UnimplementedChittyChatServer) Leave(context.Context, *Client) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Leave not implemented")
}
func (UnimplementedChittyChatServer) mustEmbedUnimplementedChittyChatServer() {}
func (UnimplementedChittyChatServer) testEmbeddedByValue()                    {}

// UnsafeChittyChatServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChittyChatServer will
// result in compilation errors.
type UnsafeChittyChatServer interface {
	mustEmbedUnimplementedChittyChatServer()
}

func RegisterChittyChatServer(s grpc.ServiceRegistrar, srv ChittyChatServer) {
	// If the following call pancis, it indicates UnimplementedChittyChatServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ChittyChat_ServiceDesc, srv)
}

func _ChittyChat_PublishMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Response)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChittyChatServer).PublishMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChittyChat_PublishMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChittyChatServer).PublishMessage(ctx, req.(*Response))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChittyChat_Join_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Client)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChittyChatServer).Join(m, &grpc.GenericServerStream[Client, Response]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ChittyChat_JoinServer = grpc.ServerStreamingServer[Response]

func _ChittyChat_Leave_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Client)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChittyChatServer).Leave(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChittyChat_Leave_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChittyChatServer).Leave(ctx, req.(*Client))
	}
	return interceptor(ctx, in, info, handler)
}

// ChittyChat_ServiceDesc is the grpc.ServiceDesc for ChittyChat service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChittyChat_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ChittyChat",
	HandlerType: (*ChittyChatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PublishMessage",
			Handler:    _ChittyChat_PublishMessage_Handler,
		},
		{
			MethodName: "Leave",
			Handler:    _ChittyChat_Leave_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Join",
			Handler:       _ChittyChat_Join_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "grpc/proto.proto",
}
