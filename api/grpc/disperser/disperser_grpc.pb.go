// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: disperser/disperser.proto

package disperser

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

const (
	Disperser_DisperseBlob_FullMethodName              = "/disperser.Disperser/DisperseBlob"
	Disperser_PaidDisperseBlob_FullMethodName          = "/disperser.Disperser/PaidDisperseBlob"
	Disperser_DisperseBlobAuthenticated_FullMethodName = "/disperser.Disperser/DisperseBlobAuthenticated"
	Disperser_GetBlobStatus_FullMethodName             = "/disperser.Disperser/GetBlobStatus"
	Disperser_RetrieveBlob_FullMethodName              = "/disperser.Disperser/RetrieveBlob"
	Disperser_GetChunk_FullMethodName                  = "/disperser.Disperser/GetChunk"
)

// DisperserClient is the client API for Disperser service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DisperserClient interface {
	// This API accepts blob to disperse from clients.
	// This executes the dispersal async, i.e. it returns once the request
	// is accepted. The client could use GetBlobStatus() API to poll the the
	// processing status of the blob.
	DisperseBlob(ctx context.Context, in *DisperseBlobRequest, opts ...grpc.CallOption) (*DisperseBlobReply, error)
	// This API require valid payments to accept blob to disperse from clients.
	// This executes the dispersal async, i.e. it returns once the request
	// is accepted. The client could use GetBlobStatus() API to poll the the
	// processing status of the blob.
	PaidDisperseBlob(ctx context.Context, in *PaidDisperseBlobRequest, opts ...grpc.CallOption) (*DisperseBlobReply, error)
	// DisperseBlobAuthenticated is similar to DisperseBlob, except that it requires the
	// client to authenticate itself via the AuthenticationData message. The protoco is as follows:
	//  1. The client sends a DisperseBlobAuthenticated request with the DisperseBlobRequest message
	//  2. The Disperser sends back a BlobAuthHeader message containing information for the client to
	//     verify and sign.
	//  3. The client verifies the BlobAuthHeader and sends back the signed BlobAuthHeader in an
	//     AuthenticationData message.
	//  4. The Disperser verifies the signature and returns a DisperseBlobReply message.
	DisperseBlobAuthenticated(ctx context.Context, opts ...grpc.CallOption) (Disperser_DisperseBlobAuthenticatedClient, error)
	// This API is meant to be polled for the blob status.
	GetBlobStatus(ctx context.Context, in *BlobStatusRequest, opts ...grpc.CallOption) (*BlobStatusReply, error)
	// This retrieves the requested blob from the Disperser's backend.
	// This is a more efficient way to retrieve blobs than directly retrieving
	// from the DA Nodes (see detail about this approach in
	// api/proto/retriever/retriever.proto).
	// The blob should have been initially dispersed via this Disperser service
	// for this API to work.
	RetrieveBlob(ctx context.Context, in *RetrieveBlobRequest, opts ...grpc.CallOption) (*RetrieveBlobReply, error)
	// Retrieves the requested chunk from the Disperser's backend.
	GetChunk(ctx context.Context, in *GetChunkRequest, opts ...grpc.CallOption) (*GetChunkReply, error)
}

type disperserClient struct {
	cc grpc.ClientConnInterface
}

func NewDisperserClient(cc grpc.ClientConnInterface) DisperserClient {
	return &disperserClient{cc}
}

func (c *disperserClient) DisperseBlob(ctx context.Context, in *DisperseBlobRequest, opts ...grpc.CallOption) (*DisperseBlobReply, error) {
	out := new(DisperseBlobReply)
	err := c.cc.Invoke(ctx, Disperser_DisperseBlob_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *disperserClient) PaidDisperseBlob(ctx context.Context, in *PaidDisperseBlobRequest, opts ...grpc.CallOption) (*DisperseBlobReply, error) {
	out := new(DisperseBlobReply)
	err := c.cc.Invoke(ctx, Disperser_PaidDisperseBlob_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *disperserClient) DisperseBlobAuthenticated(ctx context.Context, opts ...grpc.CallOption) (Disperser_DisperseBlobAuthenticatedClient, error) {
	stream, err := c.cc.NewStream(ctx, &Disperser_ServiceDesc.Streams[0], Disperser_DisperseBlobAuthenticated_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &disperserDisperseBlobAuthenticatedClient{stream}
	return x, nil
}

type Disperser_DisperseBlobAuthenticatedClient interface {
	Send(*AuthenticatedRequest) error
	Recv() (*AuthenticatedReply, error)
	grpc.ClientStream
}

type disperserDisperseBlobAuthenticatedClient struct {
	grpc.ClientStream
}

func (x *disperserDisperseBlobAuthenticatedClient) Send(m *AuthenticatedRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *disperserDisperseBlobAuthenticatedClient) Recv() (*AuthenticatedReply, error) {
	m := new(AuthenticatedReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *disperserClient) GetBlobStatus(ctx context.Context, in *BlobStatusRequest, opts ...grpc.CallOption) (*BlobStatusReply, error) {
	out := new(BlobStatusReply)
	err := c.cc.Invoke(ctx, Disperser_GetBlobStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *disperserClient) RetrieveBlob(ctx context.Context, in *RetrieveBlobRequest, opts ...grpc.CallOption) (*RetrieveBlobReply, error) {
	out := new(RetrieveBlobReply)
	err := c.cc.Invoke(ctx, Disperser_RetrieveBlob_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *disperserClient) GetChunk(ctx context.Context, in *GetChunkRequest, opts ...grpc.CallOption) (*GetChunkReply, error) {
	out := new(GetChunkReply)
	err := c.cc.Invoke(ctx, Disperser_GetChunk_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DisperserServer is the server API for Disperser service.
// All implementations must embed UnimplementedDisperserServer
// for forward compatibility
type DisperserServer interface {
	// This API accepts blob to disperse from clients.
	// This executes the dispersal async, i.e. it returns once the request
	// is accepted. The client could use GetBlobStatus() API to poll the the
	// processing status of the blob.
	DisperseBlob(context.Context, *DisperseBlobRequest) (*DisperseBlobReply, error)
	// This API require valid payments to accept blob to disperse from clients.
	// This executes the dispersal async, i.e. it returns once the request
	// is accepted. The client could use GetBlobStatus() API to poll the the
	// processing status of the blob.
	PaidDisperseBlob(context.Context, *PaidDisperseBlobRequest) (*DisperseBlobReply, error)
	// DisperseBlobAuthenticated is similar to DisperseBlob, except that it requires the
	// client to authenticate itself via the AuthenticationData message. The protoco is as follows:
	//  1. The client sends a DisperseBlobAuthenticated request with the DisperseBlobRequest message
	//  2. The Disperser sends back a BlobAuthHeader message containing information for the client to
	//     verify and sign.
	//  3. The client verifies the BlobAuthHeader and sends back the signed BlobAuthHeader in an
	//     AuthenticationData message.
	//  4. The Disperser verifies the signature and returns a DisperseBlobReply message.
	DisperseBlobAuthenticated(Disperser_DisperseBlobAuthenticatedServer) error
	// This API is meant to be polled for the blob status.
	GetBlobStatus(context.Context, *BlobStatusRequest) (*BlobStatusReply, error)
	// This retrieves the requested blob from the Disperser's backend.
	// This is a more efficient way to retrieve blobs than directly retrieving
	// from the DA Nodes (see detail about this approach in
	// api/proto/retriever/retriever.proto).
	// The blob should have been initially dispersed via this Disperser service
	// for this API to work.
	RetrieveBlob(context.Context, *RetrieveBlobRequest) (*RetrieveBlobReply, error)
	// Retrieves the requested chunk from the Disperser's backend.
	GetChunk(context.Context, *GetChunkRequest) (*GetChunkReply, error)
	mustEmbedUnimplementedDisperserServer()
}

// UnimplementedDisperserServer must be embedded to have forward compatible implementations.
type UnimplementedDisperserServer struct {
}

func (UnimplementedDisperserServer) DisperseBlob(context.Context, *DisperseBlobRequest) (*DisperseBlobReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DisperseBlob not implemented")
}
func (UnimplementedDisperserServer) PaidDisperseBlob(context.Context, *PaidDisperseBlobRequest) (*DisperseBlobReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PaidDisperseBlob not implemented")
}
func (UnimplementedDisperserServer) DisperseBlobAuthenticated(Disperser_DisperseBlobAuthenticatedServer) error {
	return status.Errorf(codes.Unimplemented, "method DisperseBlobAuthenticated not implemented")
}
func (UnimplementedDisperserServer) GetBlobStatus(context.Context, *BlobStatusRequest) (*BlobStatusReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlobStatus not implemented")
}
func (UnimplementedDisperserServer) RetrieveBlob(context.Context, *RetrieveBlobRequest) (*RetrieveBlobReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RetrieveBlob not implemented")
}
func (UnimplementedDisperserServer) GetChunk(context.Context, *GetChunkRequest) (*GetChunkReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChunk not implemented")
}
func (UnimplementedDisperserServer) mustEmbedUnimplementedDisperserServer() {}

// UnsafeDisperserServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DisperserServer will
// result in compilation errors.
type UnsafeDisperserServer interface {
	mustEmbedUnimplementedDisperserServer()
}

func RegisterDisperserServer(s grpc.ServiceRegistrar, srv DisperserServer) {
	s.RegisterService(&Disperser_ServiceDesc, srv)
}

func _Disperser_DisperseBlob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DisperseBlobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DisperserServer).DisperseBlob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Disperser_DisperseBlob_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DisperserServer).DisperseBlob(ctx, req.(*DisperseBlobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Disperser_PaidDisperseBlob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PaidDisperseBlobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DisperserServer).PaidDisperseBlob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Disperser_PaidDisperseBlob_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DisperserServer).PaidDisperseBlob(ctx, req.(*PaidDisperseBlobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Disperser_DisperseBlobAuthenticated_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DisperserServer).DisperseBlobAuthenticated(&disperserDisperseBlobAuthenticatedServer{stream})
}

type Disperser_DisperseBlobAuthenticatedServer interface {
	Send(*AuthenticatedReply) error
	Recv() (*AuthenticatedRequest, error)
	grpc.ServerStream
}

type disperserDisperseBlobAuthenticatedServer struct {
	grpc.ServerStream
}

func (x *disperserDisperseBlobAuthenticatedServer) Send(m *AuthenticatedReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *disperserDisperseBlobAuthenticatedServer) Recv() (*AuthenticatedRequest, error) {
	m := new(AuthenticatedRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Disperser_GetBlobStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BlobStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DisperserServer).GetBlobStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Disperser_GetBlobStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DisperserServer).GetBlobStatus(ctx, req.(*BlobStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Disperser_RetrieveBlob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RetrieveBlobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DisperserServer).RetrieveBlob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Disperser_RetrieveBlob_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DisperserServer).RetrieveBlob(ctx, req.(*RetrieveBlobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Disperser_GetChunk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChunkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DisperserServer).GetChunk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Disperser_GetChunk_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DisperserServer).GetChunk(ctx, req.(*GetChunkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Disperser_ServiceDesc is the grpc.ServiceDesc for Disperser service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Disperser_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "disperser.Disperser",
	HandlerType: (*DisperserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DisperseBlob",
			Handler:    _Disperser_DisperseBlob_Handler,
		},
		{
			MethodName: "PaidDisperseBlob",
			Handler:    _Disperser_PaidDisperseBlob_Handler,
		},
		{
			MethodName: "GetBlobStatus",
			Handler:    _Disperser_GetBlobStatus_Handler,
		},
		{
			MethodName: "RetrieveBlob",
			Handler:    _Disperser_RetrieveBlob_Handler,
		},
		{
			MethodName: "GetChunk",
			Handler:    _Disperser_GetChunk_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "DisperseBlobAuthenticated",
			Handler:       _Disperser_DisperseBlobAuthenticated_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "disperser/disperser.proto",
}
