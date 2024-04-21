// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.3
// source: internal/proto/api.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ClientClient is the client API for Client service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClientClient interface {
	Ping(ctx context.Context, in *Init, opts ...grpc.CallOption) (Client_PingClient, error)
}

type clientClient struct {
	cc grpc.ClientConnInterface
}

func NewClientClient(cc grpc.ClientConnInterface) ClientClient {
	return &clientClient{cc}
}

func (c *clientClient) Ping(ctx context.Context, in *Init, opts ...grpc.CallOption) (Client_PingClient, error) {
	stream, err := c.cc.NewStream(ctx, &Client_ServiceDesc.Streams[0], "/proto.Client/Ping", opts...)
	if err != nil {
		return nil, err
	}
	x := &clientPingClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Client_PingClient interface {
	Recv() (*Ping, error)
	grpc.ClientStream
}

type clientPingClient struct {
	grpc.ClientStream
}

func (x *clientPingClient) Recv() (*Ping, error) {
	m := new(Ping)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ClientServer is the server API for Client service.
// All implementations must embed UnimplementedClientServer
// for forward compatibility
type ClientServer interface {
	Ping(*Init, Client_PingServer) error
	mustEmbedUnimplementedClientServer()
}

// UnimplementedClientServer must be embedded to have forward compatible implementations.
type UnimplementedClientServer struct {
}

func (UnimplementedClientServer) Ping(*Init, Client_PingServer) error {
	return status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedClientServer) mustEmbedUnimplementedClientServer() {}

// UnsafeClientServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClientServer will
// result in compilation errors.
type UnsafeClientServer interface {
	mustEmbedUnimplementedClientServer()
}

func RegisterClientServer(s grpc.ServiceRegistrar, srv ClientServer) {
	s.RegisterService(&Client_ServiceDesc, srv)
}

func _Client_Ping_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Init)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ClientServer).Ping(m, &clientPingServer{stream})
}

type Client_PingServer interface {
	Send(*Ping) error
	grpc.ServerStream
}

type clientPingServer struct {
	grpc.ServerStream
}

func (x *clientPingServer) Send(m *Ping) error {
	return x.ServerStream.SendMsg(m)
}

// Client_ServiceDesc is the grpc.ServiceDesc for Client service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Client_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Client",
	HandlerType: (*ClientServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Ping",
			Handler:       _Client_Ping_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "internal/proto/api.proto",
}