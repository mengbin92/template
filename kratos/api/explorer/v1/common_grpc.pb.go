// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: explorer/v1/common.proto

package v1

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
	Basic_Ping_FullMethodName = "/api.explorer.v1.Basic/Ping"
)

// BasicClient is the client API for Basic service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BasicClient interface {
	Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PingReply, error)
}

type basicClient struct {
	cc grpc.ClientConnInterface
}

func NewBasicClient(cc grpc.ClientConnInterface) BasicClient {
	return &basicClient{cc}
}

func (c *basicClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PingReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PingReply)
	err := c.cc.Invoke(ctx, Basic_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BasicServer is the server API for Basic service.
// All implementations must embed UnimplementedBasicServer
// for forward compatibility.
type BasicServer interface {
	Ping(context.Context, *emptypb.Empty) (*PingReply, error)
	mustEmbedUnimplementedBasicServer()
}

// UnimplementedBasicServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedBasicServer struct{}

func (UnimplementedBasicServer) Ping(context.Context, *emptypb.Empty) (*PingReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedBasicServer) mustEmbedUnimplementedBasicServer() {}
func (UnimplementedBasicServer) testEmbeddedByValue()               {}

// UnsafeBasicServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BasicServer will
// result in compilation errors.
type UnsafeBasicServer interface {
	mustEmbedUnimplementedBasicServer()
}

func RegisterBasicServer(s grpc.ServiceRegistrar, srv BasicServer) {
	// If the following call pancis, it indicates UnimplementedBasicServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Basic_ServiceDesc, srv)
}

func _Basic_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BasicServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Basic_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BasicServer).Ping(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Basic_ServiceDesc is the grpc.ServiceDesc for Basic service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Basic_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.explorer.v1.Basic",
	HandlerType: (*BasicServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Basic_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "explorer/v1/common.proto",
}
