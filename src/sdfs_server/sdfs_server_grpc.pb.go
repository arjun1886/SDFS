// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.5
// source: sdfs_server.proto

package sdfs_server

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

// SdfsServerClient is the client API for SdfsServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SdfsServerClient interface {
	Put(ctx context.Context, opts ...grpc.CallOption) (SdfsServer_PutClient, error)
	Delete(ctx context.Context, in *DeleteInput, opts ...grpc.CallOption) (*DeleteOutput, error)
	Get(ctx context.Context, in *GetInput, opts ...grpc.CallOption) (SdfsServer_GetClient, error)
	Ls(ctx context.Context, in *LsInput, opts ...grpc.CallOption) (*LsOutput, error)
	Store(ctx context.Context, in *StoreInput, opts ...grpc.CallOption) (*StoreOutput, error)
	GetNumVersions(ctx context.Context, in *GetNumVersionsInput, opts ...grpc.CallOption) (SdfsServer_GetNumVersionsClient, error)
}

type sdfsServerClient struct {
	cc grpc.ClientConnInterface
}

func NewSdfsServerClient(cc grpc.ClientConnInterface) SdfsServerClient {
	return &sdfsServerClient{cc}
}

func (c *sdfsServerClient) Put(ctx context.Context, opts ...grpc.CallOption) (SdfsServer_PutClient, error) {
	stream, err := c.cc.NewStream(ctx, &SdfsServer_ServiceDesc.Streams[0], "/sdfs_server.SdfsServer/Put", opts...)
	if err != nil {
		return nil, err
	}
	x := &sdfsServerPutClient{stream}
	return x, nil
}

type SdfsServer_PutClient interface {
	Send(*PutInput) error
	CloseAndRecv() (*PutOutput, error)
	grpc.ClientStream
}

type sdfsServerPutClient struct {
	grpc.ClientStream
}

func (x *sdfsServerPutClient) Send(m *PutInput) error {
	return x.ClientStream.SendMsg(m)
}

func (x *sdfsServerPutClient) CloseAndRecv() (*PutOutput, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(PutOutput)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *sdfsServerClient) Delete(ctx context.Context, in *DeleteInput, opts ...grpc.CallOption) (*DeleteOutput, error) {
	out := new(DeleteOutput)
	err := c.cc.Invoke(ctx, "/sdfs_server.SdfsServer/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sdfsServerClient) Get(ctx context.Context, in *GetInput, opts ...grpc.CallOption) (SdfsServer_GetClient, error) {
	stream, err := c.cc.NewStream(ctx, &SdfsServer_ServiceDesc.Streams[1], "/sdfs_server.SdfsServer/Get", opts...)
	if err != nil {
		return nil, err
	}
	x := &sdfsServerGetClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SdfsServer_GetClient interface {
	Recv() (*GetOutput, error)
	grpc.ClientStream
}

type sdfsServerGetClient struct {
	grpc.ClientStream
}

func (x *sdfsServerGetClient) Recv() (*GetOutput, error) {
	m := new(GetOutput)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *sdfsServerClient) Ls(ctx context.Context, in *LsInput, opts ...grpc.CallOption) (*LsOutput, error) {
	out := new(LsOutput)
	err := c.cc.Invoke(ctx, "/sdfs_server.SdfsServer/Ls", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sdfsServerClient) Store(ctx context.Context, in *StoreInput, opts ...grpc.CallOption) (*StoreOutput, error) {
	out := new(StoreOutput)
	err := c.cc.Invoke(ctx, "/sdfs_server.SdfsServer/Store", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sdfsServerClient) GetNumVersions(ctx context.Context, in *GetNumVersionsInput, opts ...grpc.CallOption) (SdfsServer_GetNumVersionsClient, error) {
	stream, err := c.cc.NewStream(ctx, &SdfsServer_ServiceDesc.Streams[2], "/sdfs_server.SdfsServer/GetNumVersions", opts...)
	if err != nil {
		return nil, err
	}
	x := &sdfsServerGetNumVersionsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SdfsServer_GetNumVersionsClient interface {
	Recv() (*GetNumVersionsOutput, error)
	grpc.ClientStream
}

type sdfsServerGetNumVersionsClient struct {
	grpc.ClientStream
}

func (x *sdfsServerGetNumVersionsClient) Recv() (*GetNumVersionsOutput, error) {
	m := new(GetNumVersionsOutput)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SdfsServerServer is the server API for SdfsServer service.
// All implementations must embed UnimplementedSdfsServerServer
// for forward compatibility
type SdfsServerServer interface {
	Put(SdfsServer_PutServer) error
	Delete(context.Context, *DeleteInput) (*DeleteOutput, error)
	Get(*GetInput, SdfsServer_GetServer) error
	Ls(context.Context, *LsInput) (*LsOutput, error)
	Store(context.Context, *StoreInput) (*StoreOutput, error)
	GetNumVersions(*GetNumVersionsInput, SdfsServer_GetNumVersionsServer) error
	mustEmbedUnimplementedSdfsServerServer()
}

// UnimplementedSdfsServerServer must be embedded to have forward compatible implementations.
type UnimplementedSdfsServerServer struct {
}

func (UnimplementedSdfsServerServer) Put(SdfsServer_PutServer) error {
	return status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (UnimplementedSdfsServerServer) Delete(context.Context, *DeleteInput) (*DeleteOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedSdfsServerServer) Get(*GetInput, SdfsServer_GetServer) error {
	return status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedSdfsServerServer) Ls(context.Context, *LsInput) (*LsOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ls not implemented")
}
func (UnimplementedSdfsServerServer) Store(context.Context, *StoreInput) (*StoreOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Store not implemented")
}
func (UnimplementedSdfsServerServer) GetNumVersions(*GetNumVersionsInput, SdfsServer_GetNumVersionsServer) error {
	return status.Errorf(codes.Unimplemented, "method GetNumVersions not implemented")
}
func (UnimplementedSdfsServerServer) mustEmbedUnimplementedSdfsServerServer() {}

// UnsafeSdfsServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SdfsServerServer will
// result in compilation errors.
type UnsafeSdfsServerServer interface {
	mustEmbedUnimplementedSdfsServerServer()
}

func RegisterSdfsServerServer(s grpc.ServiceRegistrar, srv SdfsServerServer) {
	s.RegisterService(&SdfsServer_ServiceDesc, srv)
}

func _SdfsServer_Put_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SdfsServerServer).Put(&sdfsServerPutServer{stream})
}

type SdfsServer_PutServer interface {
	SendAndClose(*PutOutput) error
	Recv() (*PutInput, error)
	grpc.ServerStream
}

type sdfsServerPutServer struct {
	grpc.ServerStream
}

func (x *sdfsServerPutServer) SendAndClose(m *PutOutput) error {
	return x.ServerStream.SendMsg(m)
}

func (x *sdfsServerPutServer) Recv() (*PutInput, error) {
	m := new(PutInput)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _SdfsServer_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SdfsServerServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdfs_server.SdfsServer/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SdfsServerServer).Delete(ctx, req.(*DeleteInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _SdfsServer_Get_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetInput)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SdfsServerServer).Get(m, &sdfsServerGetServer{stream})
}

type SdfsServer_GetServer interface {
	Send(*GetOutput) error
	grpc.ServerStream
}

type sdfsServerGetServer struct {
	grpc.ServerStream
}

func (x *sdfsServerGetServer) Send(m *GetOutput) error {
	return x.ServerStream.SendMsg(m)
}

func _SdfsServer_Ls_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LsInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SdfsServerServer).Ls(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdfs_server.SdfsServer/Ls",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SdfsServerServer).Ls(ctx, req.(*LsInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _SdfsServer_Store_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SdfsServerServer).Store(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdfs_server.SdfsServer/Store",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SdfsServerServer).Store(ctx, req.(*StoreInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _SdfsServer_GetNumVersions_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetNumVersionsInput)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SdfsServerServer).GetNumVersions(m, &sdfsServerGetNumVersionsServer{stream})
}

type SdfsServer_GetNumVersionsServer interface {
	Send(*GetNumVersionsOutput) error
	grpc.ServerStream
}

type sdfsServerGetNumVersionsServer struct {
	grpc.ServerStream
}

func (x *sdfsServerGetNumVersionsServer) Send(m *GetNumVersionsOutput) error {
	return x.ServerStream.SendMsg(m)
}

// SdfsServer_ServiceDesc is the grpc.ServiceDesc for SdfsServer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SdfsServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sdfs_server.SdfsServer",
	HandlerType: (*SdfsServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Delete",
			Handler:    _SdfsServer_Delete_Handler,
		},
		{
			MethodName: "Ls",
			Handler:    _SdfsServer_Ls_Handler,
		},
		{
			MethodName: "Store",
			Handler:    _SdfsServer_Store_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Put",
			Handler:       _SdfsServer_Put_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Get",
			Handler:       _SdfsServer_Get_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetNumVersions",
			Handler:       _SdfsServer_GetNumVersions_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "sdfs_server.proto",
}