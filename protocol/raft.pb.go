// Code generated by protoc-gen-go. DO NOT EDIT.
// source: raft.proto

package protocol

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ErrCode int32

const (
	ErrCode_PartnerHasRegistered ErrCode = 0
	ErrCode_ConnectFail          ErrCode = 1
)

var ErrCode_name = map[int32]string{
	0: "PartnerHasRegistered",
	1: "ConnectFail",
}

var ErrCode_value = map[string]int32{
	"PartnerHasRegistered": 0,
	"ConnectFail":          1,
}

func (x ErrCode) String() string {
	return proto.EnumName(ErrCode_name, int32(x))
}

func (ErrCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{0}
}

// entity
type Role int32

const (
	Role_Candidate Role = 0
	Role_Follower  Role = 1
	Role_Leader    Role = 2
)

var Role_name = map[int32]string{
	0: "Candidate",
	1: "Follower",
	2: "Leader",
}

var Role_value = map[string]int32{
	"Candidate": 0,
	"Follower":  1,
	"Leader":    2,
}

func (x Role) String() string {
	return proto.EnumName(Role_name, int32(x))
}

func (Role) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{1}
}

type Pstatus int32

const (
	Pstatus_Ok         Pstatus = 0
	Pstatus_Connecting Pstatus = 1
	Pstatus_Lost       Pstatus = 2
)

var Pstatus_name = map[int32]string{
	0: "Ok",
	1: "Connecting",
	2: "Lost",
}

var Pstatus_value = map[string]int32{
	"Ok":         0,
	"Connecting": 1,
	"Lost":       2,
}

func (x Pstatus) String() string {
	return proto.EnumName(Pstatus_name, int32(x))
}

func (Pstatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{2}
}

// common msg
type Message struct {
	Mtype                uint32   `protobuf:"varint,1,opt,name=mtype,proto3" json:"mtype,omitempty"`
	Data                 []byte   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{0}
}

func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetMtype() uint32 {
	if m != nil {
		return m.Mtype
	}
	return 0
}

func (m *Message) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type NodeInfo struct {
	Id                   uint32   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Addr                 string   `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodeInfo) Reset()         { *m = NodeInfo{} }
func (m *NodeInfo) String() string { return proto.CompactTextString(m) }
func (*NodeInfo) ProtoMessage()    {}
func (*NodeInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{1}
}

func (m *NodeInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeInfo.Unmarshal(m, b)
}
func (m *NodeInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeInfo.Marshal(b, m, deterministic)
}
func (m *NodeInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeInfo.Merge(m, src)
}
func (m *NodeInfo) XXX_Size() int {
	return xxx_messageInfo_NodeInfo.Size(m)
}
func (m *NodeInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeInfo.DiscardUnknown(m)
}

var xxx_messageInfo_NodeInfo proto.InternalMessageInfo

func (m *NodeInfo) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *NodeInfo) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

// 心跳请求
type HeartBeatReq struct {
	Id                   uint64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Addr                 string   `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	Role                 int32    `protobuf:"varint,3,opt,name=role,proto3" json:"role,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HeartBeatReq) Reset()         { *m = HeartBeatReq{} }
func (m *HeartBeatReq) String() string { return proto.CompactTextString(m) }
func (*HeartBeatReq) ProtoMessage()    {}
func (*HeartBeatReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{2}
}

func (m *HeartBeatReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HeartBeatReq.Unmarshal(m, b)
}
func (m *HeartBeatReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HeartBeatReq.Marshal(b, m, deterministic)
}
func (m *HeartBeatReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeartBeatReq.Merge(m, src)
}
func (m *HeartBeatReq) XXX_Size() int {
	return xxx_messageInfo_HeartBeatReq.Size(m)
}
func (m *HeartBeatReq) XXX_DiscardUnknown() {
	xxx_messageInfo_HeartBeatReq.DiscardUnknown(m)
}

var xxx_messageInfo_HeartBeatReq proto.InternalMessageInfo

func (m *HeartBeatReq) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *HeartBeatReq) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *HeartBeatReq) GetRole() int32 {
	if m != nil {
		return m.Role
	}
	return 0
}

type GetNodeListReq struct {
	Id                   uint32   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetNodeListReq) Reset()         { *m = GetNodeListReq{} }
func (m *GetNodeListReq) String() string { return proto.CompactTextString(m) }
func (*GetNodeListReq) ProtoMessage()    {}
func (*GetNodeListReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{3}
}

func (m *GetNodeListReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetNodeListReq.Unmarshal(m, b)
}
func (m *GetNodeListReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetNodeListReq.Marshal(b, m, deterministic)
}
func (m *GetNodeListReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetNodeListReq.Merge(m, src)
}
func (m *GetNodeListReq) XXX_Size() int {
	return xxx_messageInfo_GetNodeListReq.Size(m)
}
func (m *GetNodeListReq) XXX_DiscardUnknown() {
	xxx_messageInfo_GetNodeListReq.DiscardUnknown(m)
}

var xxx_messageInfo_GetNodeListReq proto.InternalMessageInfo

func (m *GetNodeListReq) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

type GetNodeListRes struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetNodeListRes) Reset()         { *m = GetNodeListRes{} }
func (m *GetNodeListRes) String() string { return proto.CompactTextString(m) }
func (*GetNodeListRes) ProtoMessage()    {}
func (*GetNodeListRes) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{4}
}

func (m *GetNodeListRes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetNodeListRes.Unmarshal(m, b)
}
func (m *GetNodeListRes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetNodeListRes.Marshal(b, m, deterministic)
}
func (m *GetNodeListRes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetNodeListRes.Merge(m, src)
}
func (m *GetNodeListRes) XXX_Size() int {
	return xxx_messageInfo_GetNodeListRes.Size(m)
}
func (m *GetNodeListRes) XXX_DiscardUnknown() {
	xxx_messageInfo_GetNodeListRes.DiscardUnknown(m)
}

var xxx_messageInfo_GetNodeListRes proto.InternalMessageInfo

type RegisterReq struct {
	Id                   uint64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Addr                 string   `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterReq) Reset()         { *m = RegisterReq{} }
func (m *RegisterReq) String() string { return proto.CompactTextString(m) }
func (*RegisterReq) ProtoMessage()    {}
func (*RegisterReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{5}
}

func (m *RegisterReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterReq.Unmarshal(m, b)
}
func (m *RegisterReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterReq.Marshal(b, m, deterministic)
}
func (m *RegisterReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterReq.Merge(m, src)
}
func (m *RegisterReq) XXX_Size() int {
	return xxx_messageInfo_RegisterReq.Size(m)
}
func (m *RegisterReq) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterReq.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterReq proto.InternalMessageInfo

func (m *RegisterReq) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *RegisterReq) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

type RegisterRes struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterRes) Reset()         { *m = RegisterRes{} }
func (m *RegisterRes) String() string { return proto.CompactTextString(m) }
func (*RegisterRes) ProtoMessage()    {}
func (*RegisterRes) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{6}
}

func (m *RegisterRes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterRes.Unmarshal(m, b)
}
func (m *RegisterRes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterRes.Marshal(b, m, deterministic)
}
func (m *RegisterRes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterRes.Merge(m, src)
}
func (m *RegisterRes) XXX_Size() int {
	return xxx_messageInfo_RegisterRes.Size(m)
}
func (m *RegisterRes) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterRes.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterRes proto.InternalMessageInfo

type Ping struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Ping) Reset()         { *m = Ping{} }
func (m *Ping) String() string { return proto.CompactTextString(m) }
func (*Ping) ProtoMessage()    {}
func (*Ping) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{7}
}

func (m *Ping) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Ping.Unmarshal(m, b)
}
func (m *Ping) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Ping.Marshal(b, m, deterministic)
}
func (m *Ping) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ping.Merge(m, src)
}
func (m *Ping) XXX_Size() int {
	return xxx_messageInfo_Ping.Size(m)
}
func (m *Ping) XXX_DiscardUnknown() {
	xxx_messageInfo_Ping.DiscardUnknown(m)
}

var xxx_messageInfo_Ping proto.InternalMessageInfo

func (m *Ping) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type Pong struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Pong) Reset()         { *m = Pong{} }
func (m *Pong) String() string { return proto.CompactTextString(m) }
func (*Pong) ProtoMessage()    {}
func (*Pong) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{8}
}

func (m *Pong) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Pong.Unmarshal(m, b)
}
func (m *Pong) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Pong.Marshal(b, m, deterministic)
}
func (m *Pong) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Pong.Merge(m, src)
}
func (m *Pong) XXX_Size() int {
	return xxx_messageInfo_Pong.Size(m)
}
func (m *Pong) XXX_DiscardUnknown() {
	xxx_messageInfo_Pong.DiscardUnknown(m)
}

var xxx_messageInfo_Pong proto.InternalMessageInfo

func (m *Pong) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func init() {
	proto.RegisterEnum("protocol.ErrCode", ErrCode_name, ErrCode_value)
	proto.RegisterEnum("protocol.Role", Role_name, Role_value)
	proto.RegisterEnum("protocol.Pstatus", Pstatus_name, Pstatus_value)
	proto.RegisterType((*Message)(nil), "protocol.Message")
	proto.RegisterType((*NodeInfo)(nil), "protocol.NodeInfo")
	proto.RegisterType((*HeartBeatReq)(nil), "protocol.HeartBeatReq")
	proto.RegisterType((*GetNodeListReq)(nil), "protocol.GetNodeListReq")
	proto.RegisterType((*GetNodeListRes)(nil), "protocol.GetNodeListRes")
	proto.RegisterType((*RegisterReq)(nil), "protocol.RegisterReq")
	proto.RegisterType((*RegisterRes)(nil), "protocol.RegisterRes")
	proto.RegisterType((*Ping)(nil), "protocol.Ping")
	proto.RegisterType((*Pong)(nil), "protocol.Pong")
}

func init() { proto.RegisterFile("raft.proto", fileDescriptor_b042552c306ae59b) }

var fileDescriptor_b042552c306ae59b = []byte{
	// 430 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0x5d, 0x8b, 0xd3, 0x40,
	0x14, 0xed, 0xc4, 0x6c, 0x9b, 0xde, 0x7e, 0x38, 0x5e, 0xd6, 0x12, 0xe2, 0x4b, 0xc9, 0x53, 0xa9,
	0x90, 0xea, 0xae, 0x0f, 0x82, 0x4f, 0x5a, 0xb6, 0x56, 0xa8, 0x5a, 0x82, 0xe0, 0xf3, 0xec, 0xce,
	0x6d, 0x08, 0x66, 0x33, 0x75, 0x66, 0x96, 0xa5, 0x7f, 0xdb, 0x5f, 0x20, 0x99, 0x36, 0x36, 0xee,
	0x5a, 0xf0, 0x29, 0x37, 0x67, 0xce, 0xb9, 0x07, 0xee, 0x39, 0x00, 0x5a, 0x6c, 0x6c, 0xb2, 0xd5,
	0xca, 0x2a, 0x0c, 0xdc, 0xe7, 0x46, 0x15, 0xd1, 0x8b, 0x4c, 0xa9, 0xac, 0xa0, 0x99, 0x03, 0xae,
	0xef, 0x36, 0x33, 0xba, 0xdd, 0xda, 0xdd, 0x9e, 0x16, 0x5f, 0x42, 0xe7, 0x33, 0x19, 0x23, 0x32,
	0xc2, 0x73, 0x38, 0xbb, 0xb5, 0xbb, 0x2d, 0x85, 0x6c, 0xcc, 0x26, 0x83, 0x74, 0xff, 0x83, 0x08,
	0xbe, 0x14, 0x56, 0x84, 0xde, 0x98, 0x4d, 0xfa, 0xa9, 0x9b, 0xe3, 0x04, 0x82, 0x2f, 0x4a, 0xd2,
	0xa7, 0x72, 0xa3, 0x70, 0x08, 0x5e, 0x2e, 0x0f, 0x12, 0x2f, 0x97, 0x15, 0x5f, 0x48, 0xa9, 0x1d,
	0xbf, 0x9b, 0xba, 0x39, 0x5e, 0x40, 0x7f, 0x49, 0x42, 0xdb, 0x0f, 0x24, 0x6c, 0x4a, 0x3f, 0x1b,
	0x1a, 0xff, 0x94, 0xa6, 0xc2, 0xb4, 0x2a, 0x28, 0x7c, 0x32, 0x66, 0x93, 0xb3, 0xd4, 0xcd, 0xf1,
	0x18, 0x86, 0x1f, 0xc9, 0x56, 0xd6, 0xab, 0xdc, 0x3c, 0xd8, 0xe4, 0xdc, 0x63, 0xfe, 0x80, 0x61,
	0xe2, 0xd7, 0xd0, 0x4b, 0x29, 0xcb, 0x8d, 0x25, 0xfd, 0x9f, 0xd6, 0xf1, 0xa0, 0x29, 0x31, 0xf1,
	0x08, 0xfc, 0x75, 0x5e, 0x66, 0x0d, 0x69, 0xd7, 0x79, 0x55, 0xb8, 0x7a, 0x8c, 0x4f, 0xdf, 0x40,
	0xe7, 0x4a, 0xeb, 0xb9, 0x92, 0x84, 0x21, 0x9c, 0xaf, 0x85, 0xb6, 0x25, 0xe9, 0xa5, 0x30, 0xf5,
	0x4e, 0x92, 0xbc, 0x85, 0x4f, 0xa1, 0x37, 0x57, 0x65, 0x49, 0x37, 0x76, 0x21, 0xf2, 0x82, 0xb3,
	0xe9, 0x0c, 0xfc, 0x54, 0x15, 0x84, 0x03, 0xe8, 0xce, 0x45, 0x29, 0x73, 0x29, 0x2c, 0xf1, 0x16,
	0xf6, 0x21, 0x58, 0xa8, 0xa2, 0x50, 0xf7, 0xa4, 0x39, 0x43, 0x80, 0xf6, 0x8a, 0x84, 0x24, 0xcd,
	0xbd, 0xe9, 0x4b, 0xe8, 0xac, 0x8d, 0x15, 0xf6, 0xce, 0x60, 0x1b, 0xbc, 0xaf, 0x3f, 0x78, 0x0b,
	0x87, 0x00, 0x87, 0xa5, 0x79, 0x99, 0x71, 0x86, 0x01, 0xf8, 0x2b, 0x65, 0x2c, 0xf7, 0x2e, 0x7e,
	0x31, 0xf0, 0x96, 0x02, 0xdf, 0x41, 0xf7, 0x4f, 0x10, 0x38, 0x4a, 0xea, 0x8a, 0x24, 0xcd, 0x74,
	0xa2, 0x51, 0xb2, 0x2f, 0x4c, 0x52, 0x17, 0x26, 0xb9, 0xaa, 0x0a, 0x83, 0xef, 0xa1, 0xd7, 0xb8,
	0x2d, 0x86, 0x47, 0xf9, 0xdf, 0xa1, 0x44, 0xa7, 0x5e, 0x0c, 0xbe, 0x85, 0xa0, 0xbe, 0x02, 0x3e,
	0x3f, 0xb2, 0x1a, 0x01, 0x45, 0xff, 0x84, 0x0d, 0x5e, 0x40, 0xfb, 0xdb, 0xbd, 0xfa, 0x2e, 0x76,
	0xf8, 0xec, 0x48, 0x38, 0x34, 0x37, 0x7a, 0x0c, 0x4d, 0xd8, 0x2b, 0x76, 0xdd, 0x76, 0xe8, 0xe5,
	0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0xd4, 0xca, 0x13, 0x27, 0x17, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// HaClient is the client API for Ha service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type HaClient interface {
	// 心跳包
	HeartBeat(ctx context.Context, in *HeartBeatReq, opts ...grpc.CallOption) (*empty.Empty, error)
	// 获取节点信息
	GetNodeList(ctx context.Context, in *GetNodeListReq, opts ...grpc.CallOption) (*GetNodeListRes, error)
	// 注册节点
	Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterRes, error)
	// 通用双向消息
	TwoWay(ctx context.Context, opts ...grpc.CallOption) (Ha_TwoWayClient, error)
}

type haClient struct {
	cc *grpc.ClientConn
}

func NewHaClient(cc *grpc.ClientConn) HaClient {
	return &haClient{cc}
}

func (c *haClient) HeartBeat(ctx context.Context, in *HeartBeatReq, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/protocol.Ha/HeartBeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *haClient) GetNodeList(ctx context.Context, in *GetNodeListReq, opts ...grpc.CallOption) (*GetNodeListRes, error) {
	out := new(GetNodeListRes)
	err := c.cc.Invoke(ctx, "/protocol.Ha/GetNodeList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *haClient) Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterRes, error) {
	out := new(RegisterRes)
	err := c.cc.Invoke(ctx, "/protocol.Ha/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *haClient) TwoWay(ctx context.Context, opts ...grpc.CallOption) (Ha_TwoWayClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Ha_serviceDesc.Streams[0], "/protocol.Ha/TwoWay", opts...)
	if err != nil {
		return nil, err
	}
	x := &haTwoWayClient{stream}
	return x, nil
}

type Ha_TwoWayClient interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ClientStream
}

type haTwoWayClient struct {
	grpc.ClientStream
}

func (x *haTwoWayClient) Send(m *Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *haTwoWayClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// HaServer is the server API for Ha service.
type HaServer interface {
	// 心跳包
	HeartBeat(context.Context, *HeartBeatReq) (*empty.Empty, error)
	// 获取节点信息
	GetNodeList(context.Context, *GetNodeListReq) (*GetNodeListRes, error)
	// 注册节点
	Register(context.Context, *RegisterReq) (*RegisterRes, error)
	// 通用双向消息
	TwoWay(Ha_TwoWayServer) error
}

// UnimplementedHaServer can be embedded to have forward compatible implementations.
type UnimplementedHaServer struct {
}

func (*UnimplementedHaServer) HeartBeat(ctx context.Context, req *HeartBeatReq) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HeartBeat not implemented")
}
func (*UnimplementedHaServer) GetNodeList(ctx context.Context, req *GetNodeListReq) (*GetNodeListRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNodeList not implemented")
}
func (*UnimplementedHaServer) Register(ctx context.Context, req *RegisterReq) (*RegisterRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (*UnimplementedHaServer) TwoWay(srv Ha_TwoWayServer) error {
	return status.Errorf(codes.Unimplemented, "method TwoWay not implemented")
}

func RegisterHaServer(s *grpc.Server, srv HaServer) {
	s.RegisterService(&_Ha_serviceDesc, srv)
}

func _Ha_HeartBeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartBeatReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HaServer).HeartBeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.Ha/HeartBeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HaServer).HeartBeat(ctx, req.(*HeartBeatReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ha_GetNodeList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNodeListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HaServer).GetNodeList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.Ha/GetNodeList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HaServer).GetNodeList(ctx, req.(*GetNodeListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ha_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HaServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.Ha/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HaServer).Register(ctx, req.(*RegisterReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ha_TwoWay_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(HaServer).TwoWay(&haTwoWayServer{stream})
}

type Ha_TwoWayServer interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ServerStream
}

type haTwoWayServer struct {
	grpc.ServerStream
}

func (x *haTwoWayServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func (x *haTwoWayServer) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Ha_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protocol.Ha",
	HandlerType: (*HaServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HeartBeat",
			Handler:    _Ha_HeartBeat_Handler,
		},
		{
			MethodName: "GetNodeList",
			Handler:    _Ha_GetNodeList_Handler,
		},
		{
			MethodName: "Register",
			Handler:    _Ha_Register_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "TwoWay",
			Handler:       _Ha_TwoWay_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "raft.proto",
}
