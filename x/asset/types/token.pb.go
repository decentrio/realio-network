// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: realionetwork/asset/v1/token.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
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

// Token represents an asset in the module
type Token struct {
	TokenId     string `protobuf:"bytes,1,opt,name=token_id,json=tokenId,proto3" json:"token_id,omitempty"`
	Name        string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Symbol      string `protobuf:"bytes,3,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Decimal     uint32 `protobuf:"varint,4,opt,name=decimal,proto3" json:"decimal,omitempty"`
	Description string `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
}

func (m *Token) Reset()         { *m = Token{} }
func (m *Token) String() string { return proto.CompactTextString(m) }
func (*Token) ProtoMessage()    {}
func (*Token) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f83138fc60a3176, []int{0}
}
func (m *Token) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Token) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Token.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Token) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Token.Merge(m, src)
}
func (m *Token) XXX_Size() int {
	return m.Size()
}
func (m *Token) XXX_DiscardUnknown() {
	xxx_messageInfo_Token.DiscardUnknown(m)
}

var xxx_messageInfo_Token proto.InternalMessageInfo

func (m *Token) GetTokenId() string {
	if m != nil {
		return m.TokenId
	}
	return ""
}

func (m *Token) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Token) GetSymbol() string {
	if m != nil {
		return m.Symbol
	}
	return ""
}

func (m *Token) GetDecimal() uint32 {
	if m != nil {
		return m.Decimal
	}
	return 0
}

func (m *Token) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

type TokenManagement struct {
	Manager            string   `protobuf:"bytes,1,opt,name=manager,proto3" json:"manager,omitempty"`
	AddNewPrivilege    bool     `protobuf:"varint,2,opt,name=add_new_privilege,json=addNewPrivilege,proto3" json:"add_new_privilege,omitempty"`
	ExcludedPrivileges []string `protobuf:"bytes,3,rep,name=excluded_privileges,json=excludedPrivileges,proto3" json:"excluded_privileges,omitempty"`
}

func (m *TokenManagement) Reset()         { *m = TokenManagement{} }
func (m *TokenManagement) String() string { return proto.CompactTextString(m) }
func (*TokenManagement) ProtoMessage()    {}
func (*TokenManagement) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f83138fc60a3176, []int{1}
}
func (m *TokenManagement) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TokenManagement) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TokenManagement.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TokenManagement) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenManagement.Merge(m, src)
}
func (m *TokenManagement) XXX_Size() int {
	return m.Size()
}
func (m *TokenManagement) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenManagement.DiscardUnknown(m)
}

var xxx_messageInfo_TokenManagement proto.InternalMessageInfo

func (m *TokenManagement) GetManager() string {
	if m != nil {
		return m.Manager
	}
	return ""
}

func (m *TokenManagement) GetAddNewPrivilege() bool {
	if m != nil {
		return m.AddNewPrivilege
	}
	return false
}

func (m *TokenManagement) GetExcludedPrivileges() []string {
	if m != nil {
		return m.ExcludedPrivileges
	}
	return nil
}

type Balance struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Amount  uint64 `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (m *Balance) Reset()         { *m = Balance{} }
func (m *Balance) String() string { return proto.CompactTextString(m) }
func (*Balance) ProtoMessage()    {}
func (*Balance) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f83138fc60a3176, []int{2}
}
func (m *Balance) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Balance) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Balance.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Balance) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Balance.Merge(m, src)
}
func (m *Balance) XXX_Size() int {
	return m.Size()
}
func (m *Balance) XXX_DiscardUnknown() {
	xxx_messageInfo_Balance.DiscardUnknown(m)
}

var xxx_messageInfo_Balance proto.InternalMessageInfo

func (m *Balance) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Balance) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

type PrivilegeList struct {
	Privileges []string `protobuf:"bytes,1,rep,name=privileges,proto3" json:"privileges,omitempty"`
}

func (m *PrivilegeList) Reset()         { *m = PrivilegeList{} }
func (m *PrivilegeList) String() string { return proto.CompactTextString(m) }
func (*PrivilegeList) ProtoMessage()    {}
func (*PrivilegeList) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f83138fc60a3176, []int{3}
}
func (m *PrivilegeList) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PrivilegeList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PrivilegeList.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PrivilegeList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PrivilegeList.Merge(m, src)
}
func (m *PrivilegeList) XXX_Size() int {
	return m.Size()
}
func (m *PrivilegeList) XXX_DiscardUnknown() {
	xxx_messageInfo_PrivilegeList.DiscardUnknown(m)
}

var xxx_messageInfo_PrivilegeList proto.InternalMessageInfo

func (m *PrivilegeList) GetPrivileges() []string {
	if m != nil {
		return m.Privileges
	}
	return nil
}

func init() {
	proto.RegisterType((*Token)(nil), "realionetwork.asset.v1.Token")
	proto.RegisterType((*TokenManagement)(nil), "realionetwork.asset.v1.TokenManagement")
	proto.RegisterType((*Balance)(nil), "realionetwork.asset.v1.Balance")
	proto.RegisterType((*PrivilegeList)(nil), "realionetwork.asset.v1.PrivilegeList")
}

func init() {
	proto.RegisterFile("realionetwork/asset/v1/token.proto", fileDescriptor_2f83138fc60a3176)
}

var fileDescriptor_2f83138fc60a3176 = []byte{
	// 418 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xb1, 0x6e, 0xd4, 0x30,
	0x1c, 0xc6, 0xcf, 0xdc, 0xb5, 0x69, 0x8d, 0xaa, 0x0a, 0x53, 0x55, 0x6e, 0x87, 0x10, 0x65, 0x3a,
	0x21, 0x35, 0x51, 0xcb, 0xc6, 0xc6, 0x6d, 0x48, 0x05, 0xa1, 0x00, 0x0b, 0x4b, 0xe4, 0x8b, 0xff,
	0x4a, 0xad, 0xc6, 0x76, 0x64, 0xfb, 0xee, 0xda, 0xb7, 0xe8, 0x23, 0x20, 0x9e, 0x81, 0x87, 0x60,
	0xac, 0x98, 0x18, 0xd1, 0xdd, 0xc2, 0x63, 0xa0, 0xd8, 0x49, 0x39, 0xa6, 0x6e, 0xff, 0xef, 0xfb,
	0x7f, 0x89, 0x7f, 0xb6, 0x3e, 0x9c, 0x1a, 0x60, 0x8d, 0xd0, 0x0a, 0xdc, 0x4a, 0x9b, 0xeb, 0x9c,
	0x59, 0x0b, 0x2e, 0x5f, 0x9e, 0xe7, 0x4e, 0x5f, 0x83, 0xca, 0x5a, 0xa3, 0x9d, 0x26, 0xc7, 0xff,
	0x65, 0x32, 0x9f, 0xc9, 0x96, 0xe7, 0xa7, 0x47, 0xb5, 0xae, 0xb5, 0x8f, 0xe4, 0xdd, 0x14, 0xd2,
	0xa7, 0x27, 0x95, 0xb6, 0x52, 0xdb, 0x32, 0x2c, 0x82, 0x08, 0xab, 0xf4, 0x0e, 0xe1, 0x9d, 0x4f,
	0xdd, 0x8f, 0xc9, 0x09, 0xde, 0xf3, 0x27, 0x94, 0x82, 0x53, 0x94, 0xa0, 0xe9, 0x7e, 0x11, 0x79,
	0xfd, 0x96, 0x13, 0x82, 0x27, 0x8a, 0x49, 0xa0, 0x4f, 0xbc, 0xed, 0x67, 0x72, 0x8c, 0x77, 0xed,
	0xad, 0x9c, 0xeb, 0x86, 0x8e, 0xbd, 0xdb, 0x2b, 0x42, 0x71, 0xc4, 0xa1, 0x12, 0x92, 0x35, 0x74,
	0x92, 0xa0, 0xe9, 0x41, 0x31, 0x48, 0x92, 0xe0, 0xa7, 0x1c, 0x6c, 0x65, 0x44, 0xeb, 0x84, 0x56,
	0x74, 0xc7, 0x7f, 0xb6, 0x6d, 0xbd, 0x9e, 0xfc, 0xf9, 0xfa, 0x62, 0x94, 0x7e, 0x43, 0xf8, 0xd0,
	0x23, 0xbd, 0x63, 0x8a, 0xd5, 0x20, 0x41, 0x39, 0x72, 0x81, 0x23, 0xe9, 0x95, 0x09, 0x6c, 0x33,
	0xfa, 0xf3, 0xfb, 0xd9, 0x51, 0x7f, 0x93, 0x37, 0x9c, 0x1b, 0xb0, 0xf6, 0xa3, 0x33, 0x42, 0xd5,
	0xc5, 0x10, 0x24, 0x2f, 0xf1, 0x33, 0xc6, 0x79, 0xa9, 0x60, 0x55, 0xb6, 0x46, 0x2c, 0x45, 0x03,
	0x75, 0xb8, 0xc2, 0x5e, 0x71, 0xc8, 0x38, 0x7f, 0x0f, 0xab, 0x0f, 0x83, 0x4d, 0x72, 0xfc, 0x1c,
	0x6e, 0xaa, 0x66, 0xc1, 0x81, 0xff, 0x0b, 0x5b, 0x3a, 0x4e, 0xc6, 0xd3, 0xfd, 0x82, 0x0c, 0xab,
	0x87, 0xbc, 0x4d, 0x3f, 0xe3, 0x68, 0xc6, 0x1a, 0xa6, 0x2a, 0xe8, 0xd8, 0x58, 0x20, 0x78, 0x9c,
	0xad, 0x0f, 0x76, 0xaf, 0xc7, 0xa4, 0x5e, 0x28, 0xe7, 0x81, 0x26, 0x45, 0xaf, 0xd2, 0x1c, 0x1f,
	0x3c, 0x1c, 0x72, 0x29, 0xac, 0x23, 0x31, 0xc6, 0x5b, 0x3c, 0xc8, 0xf3, 0x6c, 0x39, 0xb3, 0xcb,
	0x1f, 0xeb, 0x18, 0xdd, 0xaf, 0x63, 0xf4, 0x7b, 0x1d, 0xa3, 0xbb, 0x4d, 0x3c, 0xba, 0xdf, 0xc4,
	0xa3, 0x5f, 0x9b, 0x78, 0xf4, 0xe5, 0xa2, 0x16, 0xee, 0x6a, 0x31, 0xcf, 0x2a, 0x2d, 0xf3, 0xd0,
	0x16, 0x07, 0xd5, 0x55, 0x3f, 0x9e, 0x0d, 0xed, 0xba, 0xe9, 0xfb, 0xe5, 0x6e, 0x5b, 0xb0, 0xf3,
	0x5d, 0x5f, 0x8a, 0x57, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0xff, 0x10, 0x32, 0xba, 0x83, 0x02,
	0x00, 0x00,
}

func (m *Token) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Token) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Token) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintToken(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x2a
	}
	if m.Decimal != 0 {
		i = encodeVarintToken(dAtA, i, uint64(m.Decimal))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Symbol) > 0 {
		i -= len(m.Symbol)
		copy(dAtA[i:], m.Symbol)
		i = encodeVarintToken(dAtA, i, uint64(len(m.Symbol)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintToken(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.TokenId) > 0 {
		i -= len(m.TokenId)
		copy(dAtA[i:], m.TokenId)
		i = encodeVarintToken(dAtA, i, uint64(len(m.TokenId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TokenManagement) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TokenManagement) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TokenManagement) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ExcludedPrivileges) > 0 {
		for iNdEx := len(m.ExcludedPrivileges) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.ExcludedPrivileges[iNdEx])
			copy(dAtA[i:], m.ExcludedPrivileges[iNdEx])
			i = encodeVarintToken(dAtA, i, uint64(len(m.ExcludedPrivileges[iNdEx])))
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.AddNewPrivilege {
		i--
		if m.AddNewPrivilege {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if len(m.Manager) > 0 {
		i -= len(m.Manager)
		copy(dAtA[i:], m.Manager)
		i = encodeVarintToken(dAtA, i, uint64(len(m.Manager)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Balance) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Balance) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Balance) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Amount != 0 {
		i = encodeVarintToken(dAtA, i, uint64(m.Amount))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintToken(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *PrivilegeList) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PrivilegeList) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PrivilegeList) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Privileges) > 0 {
		for iNdEx := len(m.Privileges) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Privileges[iNdEx])
			copy(dAtA[i:], m.Privileges[iNdEx])
			i = encodeVarintToken(dAtA, i, uint64(len(m.Privileges[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintToken(dAtA []byte, offset int, v uint64) int {
	offset -= sovToken(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Token) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.TokenId)
	if l > 0 {
		n += 1 + l + sovToken(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovToken(uint64(l))
	}
	l = len(m.Symbol)
	if l > 0 {
		n += 1 + l + sovToken(uint64(l))
	}
	if m.Decimal != 0 {
		n += 1 + sovToken(uint64(m.Decimal))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovToken(uint64(l))
	}
	return n
}

func (m *TokenManagement) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Manager)
	if l > 0 {
		n += 1 + l + sovToken(uint64(l))
	}
	if m.AddNewPrivilege {
		n += 2
	}
	if len(m.ExcludedPrivileges) > 0 {
		for _, s := range m.ExcludedPrivileges {
			l = len(s)
			n += 1 + l + sovToken(uint64(l))
		}
	}
	return n
}

func (m *Balance) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovToken(uint64(l))
	}
	if m.Amount != 0 {
		n += 1 + sovToken(uint64(m.Amount))
	}
	return n
}

func (m *PrivilegeList) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Privileges) > 0 {
		for _, s := range m.Privileges {
			l = len(s)
			n += 1 + l + sovToken(uint64(l))
		}
	}
	return n
}

func sovToken(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozToken(x uint64) (n int) {
	return sovToken(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Token) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowToken
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Token: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Token: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowToken
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthToken
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthToken
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TokenId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowToken
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthToken
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthToken
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Symbol", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowToken
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthToken
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthToken
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Symbol = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Decimal", wireType)
			}
			m.Decimal = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowToken
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Decimal |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowToken
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthToken
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthToken
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipToken(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthToken
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TokenManagement) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowToken
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: TokenManagement: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TokenManagement: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Manager", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowToken
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthToken
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthToken
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Manager = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AddNewPrivilege", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowToken
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.AddNewPrivilege = bool(v != 0)
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExcludedPrivileges", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowToken
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthToken
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthToken
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ExcludedPrivileges = append(m.ExcludedPrivileges, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipToken(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthToken
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Balance) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowToken
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Balance: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Balance: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowToken
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthToken
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthToken
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			m.Amount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowToken
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Amount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipToken(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthToken
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *PrivilegeList) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowToken
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PrivilegeList: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PrivilegeList: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Privileges", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowToken
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthToken
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthToken
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Privileges = append(m.Privileges, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipToken(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthToken
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipToken(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowToken
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowToken
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowToken
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthToken
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupToken
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthToken
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthToken        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowToken          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupToken = fmt.Errorf("proto: unexpected end of group")
)
