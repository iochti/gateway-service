// Code generated by protoc-gen-go.
// source: user.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	user.proto

It has these top-level messages:
	UserRequest
	UserMessage
	UserID
	UserDeleted
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type UserRequest struct {
	Categ string `protobuf:"bytes,1,opt,name=categ" json:"categ,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (m *UserRequest) Reset()                    { *m = UserRequest{} }
func (m *UserRequest) String() string            { return proto1.CompactTextString(m) }
func (*UserRequest) ProtoMessage()               {}
func (*UserRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *UserRequest) GetCateg() string {
	if m != nil {
		return m.Categ
	}
	return ""
}

func (m *UserRequest) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type UserMessage struct {
	User []byte `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (m *UserMessage) Reset()                    { *m = UserMessage{} }
func (m *UserMessage) String() string            { return proto1.CompactTextString(m) }
func (*UserMessage) ProtoMessage()               {}
func (*UserMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *UserMessage) GetUser() []byte {
	if m != nil {
		return m.User
	}
	return nil
}

type UserID struct {
	Id int32 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
}

func (m *UserID) Reset()                    { *m = UserID{} }
func (m *UserID) String() string            { return proto1.CompactTextString(m) }
func (*UserID) ProtoMessage()               {}
func (*UserID) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *UserID) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

type UserDeleted struct {
	Deleted bool  `protobuf:"varint,1,opt,name=deleted" json:"deleted,omitempty"`
	Id      int32 `protobuf:"varint,2,opt,name=id" json:"id,omitempty"`
}

func (m *UserDeleted) Reset()                    { *m = UserDeleted{} }
func (m *UserDeleted) String() string            { return proto1.CompactTextString(m) }
func (*UserDeleted) ProtoMessage()               {}
func (*UserDeleted) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *UserDeleted) GetDeleted() bool {
	if m != nil {
		return m.Deleted
	}
	return false
}

func (m *UserDeleted) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func init() {
	proto1.RegisterType((*UserRequest)(nil), "proto.UserRequest")
	proto1.RegisterType((*UserMessage)(nil), "proto.UserMessage")
	proto1.RegisterType((*UserID)(nil), "proto.UserID")
	proto1.RegisterType((*UserDeleted)(nil), "proto.UserDeleted")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for UserSvc service

type UserSvcClient interface {
	GetUser(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserMessage, error)
	CreateUser(ctx context.Context, in *UserMessage, opts ...grpc.CallOption) (*UserMessage, error)
	DeleteUser(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*UserDeleted, error)
}

type userSvcClient struct {
	cc *grpc.ClientConn
}

func NewUserSvcClient(cc *grpc.ClientConn) UserSvcClient {
	return &userSvcClient{cc}
}

func (c *userSvcClient) GetUser(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserMessage, error) {
	out := new(UserMessage)
	err := grpc.Invoke(ctx, "/proto.UserSvc/GetUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userSvcClient) CreateUser(ctx context.Context, in *UserMessage, opts ...grpc.CallOption) (*UserMessage, error) {
	out := new(UserMessage)
	err := grpc.Invoke(ctx, "/proto.UserSvc/CreateUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userSvcClient) DeleteUser(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*UserDeleted, error) {
	out := new(UserDeleted)
	err := grpc.Invoke(ctx, "/proto.UserSvc/DeleteUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UserSvc service

type UserSvcServer interface {
	GetUser(context.Context, *UserRequest) (*UserMessage, error)
	CreateUser(context.Context, *UserMessage) (*UserMessage, error)
	DeleteUser(context.Context, *UserID) (*UserDeleted, error)
}

func RegisterUserSvcServer(s *grpc.Server, srv UserSvcServer) {
	s.RegisterService(&_UserSvc_serviceDesc, srv)
}

func _UserSvc_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserSvcServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserSvc/GetUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserSvcServer).GetUser(ctx, req.(*UserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserSvc_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserSvcServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserSvc/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserSvcServer).CreateUser(ctx, req.(*UserMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserSvc_DeleteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserSvcServer).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserSvc/DeleteUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserSvcServer).DeleteUser(ctx, req.(*UserID))
	}
	return interceptor(ctx, in, info, handler)
}

var _UserSvc_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.UserSvc",
	HandlerType: (*UserSvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUser",
			Handler:    _UserSvc_GetUser_Handler,
		},
		{
			MethodName: "CreateUser",
			Handler:    _UserSvc_CreateUser_Handler,
		},
		{
			MethodName: "DeleteUser",
			Handler:    _UserSvc_DeleteUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user.proto",
}

func init() { proto1.RegisterFile("user.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 226 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x8f, 0xcd, 0x4a, 0xc5, 0x30,
	0x10, 0x85, 0x6f, 0x83, 0xbd, 0xd5, 0xf1, 0x67, 0x31, 0xb8, 0x08, 0xae, 0x34, 0x2b, 0x57, 0x17,
	0xb4, 0xa0, 0xb8, 0xb6, 0x20, 0x5d, 0xb8, 0x89, 0xf8, 0x00, 0xb1, 0x19, 0x4a, 0xa1, 0x50, 0x4d,
	0xd2, 0xbe, 0x93, 0x6f, 0x29, 0x99, 0xa4, 0x50, 0x11, 0x57, 0x99, 0x73, 0x72, 0xbe, 0x61, 0x0e,
	0xc0, 0xec, 0xc9, 0x1d, 0x3e, 0xdd, 0x14, 0x26, 0x2c, 0xf9, 0x51, 0x4f, 0x70, 0xfa, 0xee, 0xc9,
	0x69, 0xfa, 0x9a, 0xc9, 0x07, 0xbc, 0x84, 0xb2, 0x33, 0x81, 0x7a, 0x59, 0x5c, 0x17, 0xb7, 0x27,
	0x3a, 0x89, 0xe8, 0x2e, 0x66, 0x9c, 0x49, 0x8a, 0xe4, 0xb2, 0x50, 0x37, 0x09, 0x7d, 0x25, 0xef,
	0x4d, 0x4f, 0x88, 0x70, 0x14, 0xd7, 0x33, 0x79, 0xa6, 0x79, 0x56, 0x12, 0xf6, 0x31, 0xd2, 0x36,
	0x78, 0x01, 0x62, 0xb0, 0xfc, 0x57, 0x6a, 0x31, 0x58, 0xf5, 0x98, 0xe0, 0x86, 0x46, 0x0a, 0x64,
	0x51, 0x42, 0x65, 0xd3, 0xc8, 0x99, 0x63, 0xbd, 0xca, 0x0c, 0x8a, 0x15, 0xbc, 0xff, 0x2e, 0xa0,
	0x8a, 0xe4, 0xdb, 0xd2, 0x61, 0x0d, 0xd5, 0x0b, 0x85, 0xa8, 0x10, 0x53, 0xad, 0xc3, 0xa6, 0xcc,
	0xd5, 0xd6, 0xcb, 0x57, 0xaa, 0x1d, 0x3e, 0x00, 0x3c, 0x3b, 0x32, 0x81, 0xfe, 0x70, 0x39, 0xf3,
	0x0f, 0x77, 0x07, 0x90, 0xae, 0x65, 0xee, 0x7c, 0x93, 0x69, 0x9b, 0x5f, 0x48, 0xee, 0xa4, 0x76,
	0x1f, 0x7b, 0x36, 0xeb, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x77, 0x97, 0xbf, 0x2f, 0x78, 0x01,
	0x00, 0x00,
}
