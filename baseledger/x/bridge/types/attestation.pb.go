// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: bridge/attestation.proto

package types

import (
	fmt "fmt"
	types "github.com/cosmos/cosmos-sdk/codec/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
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

type ClaimType int32

const (
	CLAIM_TYPE_UNSPECIFIED        ClaimType = 0
	CLAIM_UBT_DEPOSITED           ClaimType = 1
	CLAIM_VALIDATOR_POWER_CHANGED ClaimType = 2
)

var ClaimType_name = map[int32]string{
	0: "CLAIM_TYPE_UNSPECIFIED",
	1: "CLAIM_UBT_DEPOSITED",
	2: "CLAIM_VALIDATOR_POWER_CHANGED",
}

var ClaimType_value = map[string]int32{
	"CLAIM_TYPE_UNSPECIFIED":        0,
	"CLAIM_UBT_DEPOSITED":           1,
	"CLAIM_VALIDATOR_POWER_CHANGED": 2,
}

func (x ClaimType) String() string {
	return proto.EnumName(ClaimType_name, int32(x))
}

func (ClaimType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_2c79119916c99d41, []int{0}
}

type Attestation struct {
	Observed    bool                                   `protobuf:"varint,1,opt,name=observed,proto3" json:"observed,omitempty"`
	Votes       []string                               `protobuf:"bytes,2,rep,name=votes,proto3" json:"votes,omitempty"`
	Height      uint64                                 `protobuf:"varint,3,opt,name=height,proto3" json:"height,omitempty"`
	Claim       *types.Any                             `protobuf:"bytes,4,opt,name=claim,proto3" json:"claim,omitempty"`
	UbtPrices   []string                               `protobuf:"bytes,5,rep,name=ubtPrices,proto3" json:"ubtPrices,omitempty"`
	AvgUbtPrice github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,6,opt,name=avgUbtPrice,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"avgUbtPrice"`
}

func (m *Attestation) Reset()         { *m = Attestation{} }
func (m *Attestation) String() string { return proto.CompactTextString(m) }
func (*Attestation) ProtoMessage()    {}
func (*Attestation) Descriptor() ([]byte, []int) {
	return fileDescriptor_2c79119916c99d41, []int{0}
}
func (m *Attestation) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Attestation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Attestation.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Attestation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Attestation.Merge(m, src)
}
func (m *Attestation) XXX_Size() int {
	return m.Size()
}
func (m *Attestation) XXX_DiscardUnknown() {
	xxx_messageInfo_Attestation.DiscardUnknown(m)
}

var xxx_messageInfo_Attestation proto.InternalMessageInfo

func (m *Attestation) GetObserved() bool {
	if m != nil {
		return m.Observed
	}
	return false
}

func (m *Attestation) GetVotes() []string {
	if m != nil {
		return m.Votes
	}
	return nil
}

func (m *Attestation) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *Attestation) GetClaim() *types.Any {
	if m != nil {
		return m.Claim
	}
	return nil
}

func (m *Attestation) GetUbtPrices() []string {
	if m != nil {
		return m.UbtPrices
	}
	return nil
}

func init() {
	proto.RegisterEnum("Baseledger.baseledger.bridge.ClaimType", ClaimType_name, ClaimType_value)
	proto.RegisterType((*Attestation)(nil), "Baseledger.baseledger.bridge.Attestation")
}

func init() { proto.RegisterFile("bridge/attestation.proto", fileDescriptor_2c79119916c99d41) }

var fileDescriptor_2c79119916c99d41 = []byte{
	// 417 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x92, 0xdf, 0x6e, 0xd3, 0x30,
	0x14, 0xc6, 0xe3, 0xad, 0xad, 0x56, 0xf7, 0xa6, 0x32, 0xd5, 0x08, 0xd1, 0xc8, 0x02, 0x17, 0x28,
	0x9a, 0x84, 0x2d, 0xc1, 0x13, 0xa4, 0x49, 0x80, 0x48, 0x63, 0x8d, 0xb2, 0x14, 0x04, 0x37, 0x51,
	0xfe, 0x18, 0x37, 0xa2, 0x8d, 0xab, 0xd8, 0xad, 0xc8, 0x1b, 0x70, 0xc9, 0x3b, 0xf0, 0x32, 0xbb,
	0xdc, 0x25, 0xe2, 0x62, 0x42, 0xed, 0x23, 0xf0, 0x02, 0x68, 0x49, 0x58, 0x7a, 0xe5, 0xf3, 0xf3,
	0xf9, 0x7c, 0xce, 0x67, 0x1f, 0x43, 0x35, 0x29, 0xf3, 0x8c, 0x51, 0x12, 0x4b, 0x49, 0x85, 0x8c,
	0x65, 0xce, 0x0b, 0xbc, 0x2e, 0xb9, 0xe4, 0xe8, 0x6c, 0x1a, 0x0b, 0xba, 0xa4, 0x19, 0xa3, 0x25,
	0x4e, 0x0e, 0xc2, 0x5a, 0xaf, 0x4d, 0x18, 0x67, 0xbc, 0x16, 0x92, 0xfb, 0xa8, 0x39, 0xa3, 0x3d,
	0x61, 0x9c, 0xb3, 0x25, 0x25, 0x35, 0x25, 0x9b, 0x2f, 0x24, 0x2e, 0xaa, 0x26, 0xf5, 0xfc, 0x2f,
	0x80, 0x23, 0xab, 0x6b, 0x82, 0x34, 0x78, 0xc2, 0x13, 0x41, 0xcb, 0x2d, 0xcd, 0x54, 0x60, 0x00,
	0xf3, 0x24, 0x78, 0x60, 0x34, 0x81, 0xfd, 0x2d, 0x97, 0x54, 0xa8, 0x47, 0xc6, 0xb1, 0x39, 0x0c,
	0x1a, 0x40, 0xa7, 0x70, 0xb0, 0xa0, 0x39, 0x5b, 0x48, 0xf5, 0xd8, 0x00, 0x66, 0x2f, 0x68, 0x09,
	0x5d, 0xc0, 0x7e, 0xba, 0x8c, 0xf3, 0x95, 0xda, 0x33, 0x80, 0x39, 0x7a, 0x35, 0xc1, 0x8d, 0x09,
	0xfc, 0xdf, 0x04, 0xb6, 0x8a, 0x2a, 0x68, 0x24, 0xe8, 0x0c, 0x0e, 0x37, 0x89, 0xf4, 0xcb, 0x3c,
	0xa5, 0x42, 0xed, 0xd7, 0xd5, 0xbb, 0x0d, 0xe4, 0xc3, 0x51, 0xbc, 0x65, 0xf3, 0x96, 0xd5, 0x81,
	0x01, 0xcc, 0xe1, 0x14, 0xdf, 0xdc, 0x9d, 0x2b, 0xbf, 0xef, 0xce, 0x5f, 0xb0, 0x5c, 0x2e, 0x36,
	0x09, 0x4e, 0xf9, 0x8a, 0xa4, 0x5c, 0xac, 0xb8, 0x68, 0x97, 0x97, 0x22, 0xfb, 0x4a, 0x64, 0xb5,
	0xa6, 0x02, 0x7b, 0x85, 0x0c, 0x0e, 0x4b, 0x5c, 0xe4, 0x70, 0x68, 0xdf, 0x37, 0x0e, 0xab, 0x35,
	0x45, 0x1a, 0x3c, 0xb5, 0x2f, 0x2d, 0xef, 0x7d, 0x14, 0x7e, 0xf2, 0xdd, 0x68, 0x7e, 0x75, 0xed,
	0xbb, 0xb6, 0xf7, 0xc6, 0x73, 0x9d, 0xb1, 0x82, 0x1e, 0xc3, 0x47, 0x4d, 0x6e, 0x3e, 0x0d, 0x23,
	0xc7, 0xf5, 0x67, 0xd7, 0x5e, 0xe8, 0x3a, 0x63, 0x80, 0x9e, 0xc1, 0xa7, 0x4d, 0xe2, 0x83, 0x75,
	0xe9, 0x39, 0x56, 0x38, 0x0b, 0x22, 0x7f, 0xf6, 0xd1, 0x0d, 0x22, 0xfb, 0x9d, 0x75, 0xf5, 0xd6,
	0x75, 0xc6, 0x47, 0x5a, 0xef, 0xfb, 0x4f, 0x5d, 0x99, 0x7a, 0x37, 0x3b, 0x1d, 0xdc, 0xee, 0x74,
	0xf0, 0x67, 0xa7, 0x83, 0x1f, 0x7b, 0x5d, 0xb9, 0xdd, 0xeb, 0xca, 0xaf, 0xbd, 0xae, 0x7c, 0x26,
	0x07, 0xce, 0xbb, 0xa1, 0x92, 0x6e, 0xa8, 0xe4, 0x1b, 0x69, 0xbf, 0x41, 0x7d, 0x8d, 0x64, 0x50,
	0x3f, 0xdd, 0xeb, 0x7f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x8b, 0xfe, 0x8a, 0x25, 0x1d, 0x02, 0x00,
	0x00,
}

func (m *Attestation) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Attestation) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Attestation) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.AvgUbtPrice.Size()
		i -= size
		if _, err := m.AvgUbtPrice.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintAttestation(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	if len(m.UbtPrices) > 0 {
		for iNdEx := len(m.UbtPrices) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.UbtPrices[iNdEx])
			copy(dAtA[i:], m.UbtPrices[iNdEx])
			i = encodeVarintAttestation(dAtA, i, uint64(len(m.UbtPrices[iNdEx])))
			i--
			dAtA[i] = 0x2a
		}
	}
	if m.Claim != nil {
		{
			size, err := m.Claim.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintAttestation(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.Height != 0 {
		i = encodeVarintAttestation(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Votes) > 0 {
		for iNdEx := len(m.Votes) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Votes[iNdEx])
			copy(dAtA[i:], m.Votes[iNdEx])
			i = encodeVarintAttestation(dAtA, i, uint64(len(m.Votes[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Observed {
		i--
		if m.Observed {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintAttestation(dAtA []byte, offset int, v uint64) int {
	offset -= sovAttestation(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Attestation) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Observed {
		n += 2
	}
	if len(m.Votes) > 0 {
		for _, s := range m.Votes {
			l = len(s)
			n += 1 + l + sovAttestation(uint64(l))
		}
	}
	if m.Height != 0 {
		n += 1 + sovAttestation(uint64(m.Height))
	}
	if m.Claim != nil {
		l = m.Claim.Size()
		n += 1 + l + sovAttestation(uint64(l))
	}
	if len(m.UbtPrices) > 0 {
		for _, s := range m.UbtPrices {
			l = len(s)
			n += 1 + l + sovAttestation(uint64(l))
		}
	}
	l = m.AvgUbtPrice.Size()
	n += 1 + l + sovAttestation(uint64(l))
	return n
}

func sovAttestation(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozAttestation(x uint64) (n int) {
	return sovAttestation(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Attestation) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAttestation
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
			return fmt.Errorf("proto: Attestation: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Attestation: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Observed", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAttestation
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
			m.Observed = bool(v != 0)
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Votes", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAttestation
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
				return ErrInvalidLengthAttestation
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAttestation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Votes = append(m.Votes, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAttestation
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Claim", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAttestation
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthAttestation
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAttestation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Claim == nil {
				m.Claim = &types.Any{}
			}
			if err := m.Claim.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UbtPrices", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAttestation
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
				return ErrInvalidLengthAttestation
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAttestation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UbtPrices = append(m.UbtPrices, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AvgUbtPrice", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAttestation
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
				return ErrInvalidLengthAttestation
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAttestation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.AvgUbtPrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAttestation(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAttestation
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
func skipAttestation(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAttestation
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
					return 0, ErrIntOverflowAttestation
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
					return 0, ErrIntOverflowAttestation
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
				return 0, ErrInvalidLengthAttestation
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupAttestation
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthAttestation
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthAttestation        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAttestation          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupAttestation = fmt.Errorf("proto: unexpected end of group")
)
