// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: communication.proto

package communication

import (
	context "context"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type Message struct {
	Body                 string   `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_b64068f22c460ac1, []int{0}
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

func (m *Message) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

type Invoice struct {
	CompanyID            string   `protobuf:"bytes,1,opt,name=companyID,proto3" json:"companyID,omitempty"`
	Price                float64  `protobuf:"fixed64,2,opt,name=price,proto3" json:"price,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Invoice) Reset()         { *m = Invoice{} }
func (m *Invoice) String() string { return proto.CompactTextString(m) }
func (*Invoice) ProtoMessage()    {}
func (*Invoice) Descriptor() ([]byte, []int) {
	return fileDescriptor_b64068f22c460ac1, []int{1}
}
func (m *Invoice) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Invoice.Unmarshal(m, b)
}
func (m *Invoice) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Invoice.Marshal(b, m, deterministic)
}
func (m *Invoice) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Invoice.Merge(m, src)
}
func (m *Invoice) XXX_Size() int {
	return xxx_messageInfo_Invoice.Size(m)
}
func (m *Invoice) XXX_DiscardUnknown() {
	xxx_messageInfo_Invoice.DiscardUnknown(m)
}

var xxx_messageInfo_Invoice proto.InternalMessageInfo

func (m *Invoice) GetCompanyID() string {
	if m != nil {
		return m.CompanyID
	}
	return ""
}

func (m *Invoice) GetPrice() float64 {
	if m != nil {
		return m.Price
	}
	return 0
}

func init() {
	proto.RegisterType((*Message)(nil), "communication.Message")
	proto.RegisterType((*Invoice)(nil), "communication.Invoice")
}

func init() { proto.RegisterFile("communication.proto", fileDescriptor_b64068f22c460ac1) }

var fileDescriptor_b64068f22c460ac1 = []byte{
	// 164 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4e, 0xce, 0xcf, 0xcd,
	0x2d, 0xcd, 0xcb, 0x4c, 0x4e, 0x2c, 0xc9, 0xcc, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17,
	0xe2, 0x45, 0x11, 0x54, 0x92, 0xe5, 0x62, 0xf7, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0x15, 0x12,
	0xe2, 0x62, 0x49, 0xca, 0x4f, 0xa9, 0x94, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0xb3, 0x95,
	0x6c, 0xb9, 0xd8, 0x3d, 0xf3, 0xca, 0xf2, 0x33, 0x93, 0x53, 0x85, 0x64, 0xb8, 0x38, 0x93, 0xf3,
	0x73, 0x0b, 0x12, 0xf3, 0x2a, 0x3d, 0x5d, 0xa0, 0x6a, 0x10, 0x02, 0x42, 0x22, 0x5c, 0xac, 0x05,
	0x45, 0x99, 0xc9, 0xa9, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0x8c, 0x41, 0x10, 0x8e, 0x51, 0x38, 0x97,
	0x88, 0x33, 0xb2, 0x75, 0xc1, 0xa9, 0x45, 0x65, 0x20, 0xb3, 0xec, 0xb9, 0xb8, 0x83, 0x53, 0x73,
	0x72, 0x60, 0x46, 0x8b, 0xe9, 0xa1, 0xba, 0x14, 0x2a, 0x2e, 0x85, 0x2e, 0x0e, 0x75, 0xa9, 0x12,
	0x43, 0x12, 0x1b, 0xd8, 0x33, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x12, 0xa9, 0x24, 0x99,
	0xe3, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// CommunicationServiceClient is the client API for CommunicationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CommunicationServiceClient interface {
	SellInvoice(ctx context.Context, in *Invoice, opts ...grpc.CallOption) (*Message, error)
}

type communicationServiceClient struct {
	cc *grpc.ClientConn
}

func NewCommunicationServiceClient(cc *grpc.ClientConn) CommunicationServiceClient {
	return &communicationServiceClient{cc}
}

func (c *communicationServiceClient) SellInvoice(ctx context.Context, in *Invoice, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, "/communication.CommunicationService/SellInvoice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommunicationServiceServer is the server API for CommunicationService service.
type CommunicationServiceServer interface {
	SellInvoice(context.Context, *Invoice) (*Message, error)
}

// UnimplementedCommunicationServiceServer can be embedded to have forward compatible implementations.
type UnimplementedCommunicationServiceServer struct {
}

func (*UnimplementedCommunicationServiceServer) SellInvoice(ctx context.Context, req *Invoice) (*Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SellInvoice not implemented")
}

func RegisterCommunicationServiceServer(s *grpc.Server, srv CommunicationServiceServer) {
	s.RegisterService(&_CommunicationService_serviceDesc, srv)
}

func _CommunicationService_SellInvoice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Invoice)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunicationServiceServer).SellInvoice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/communication.CommunicationService/SellInvoice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunicationServiceServer).SellInvoice(ctx, req.(*Invoice))
	}
	return interceptor(ctx, in, info, handler)
}

var _CommunicationService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "communication.CommunicationService",
	HandlerType: (*CommunicationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SellInvoice",
			Handler:    _CommunicationService_SellInvoice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "communication.proto",
}
