// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.4
// source: proto/mpc.proto

package Security_Miniproject2

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

// ExchangeClient is the client API for Exchange service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExchangeClient interface {
	SendShare(ctx context.Context, in *SendShareRequest, opts ...grpc.CallOption) (*SendShareResponse, error)
}

type exchangeClient struct {
	cc grpc.ClientConnInterface
}

func NewExchangeClient(cc grpc.ClientConnInterface) ExchangeClient {
	return &exchangeClient{cc}
}

func (c *exchangeClient) SendShare(ctx context.Context, in *SendShareRequest, opts ...grpc.CallOption) (*SendShareResponse, error) {
	out := new(SendShareResponse)
	err := c.cc.Invoke(ctx, "/proto.Exchange/sendShare", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExchangeServer is the server API for Exchange service.
// All implementations must embed UnimplementedExchangeServer
// for forward compatibility
type ExchangeServer interface {
	SendShare(context.Context, *SendShareRequest) (*SendShareResponse, error)
	mustEmbedUnimplementedExchangeServer()
}

// UnimplementedExchangeServer must be embedded to have forward compatible implementations.
type UnimplementedExchangeServer struct {
}

func (UnimplementedExchangeServer) SendShare(context.Context, *SendShareRequest) (*SendShareResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendShare not implemented")
}
func (UnimplementedExchangeServer) mustEmbedUnimplementedExchangeServer() {}

// UnsafeExchangeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExchangeServer will
// result in compilation errors.
type UnsafeExchangeServer interface {
	mustEmbedUnimplementedExchangeServer()
}

func RegisterExchangeServer(s grpc.ServiceRegistrar, srv ExchangeServer) {
	s.RegisterService(&Exchange_ServiceDesc, srv)
}

func _Exchange_SendShare_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendShareRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExchangeServer).SendShare(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Exchange/sendShare",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExchangeServer).SendShare(ctx, req.(*SendShareRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Exchange_ServiceDesc is the grpc.ServiceDesc for Exchange service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Exchange_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Exchange",
	HandlerType: (*ExchangeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "sendShare",
			Handler:    _Exchange_SendShare_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/mpc.proto",
}
