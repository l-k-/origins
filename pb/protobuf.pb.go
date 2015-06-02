// Code generated by protoc-gen-go.
// source: pb/protobuf.proto
// DO NOT EDIT!

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	pb/protobuf.proto

It has these top-level messages:
	Transaction
	Log
	Segment
	Fact
*/
package pb

import proto "github.com/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type Transaction struct {
	ID               *uint64 `protobuf:"varint,1,req" json:"ID,omitempty"`
	StartTime        *int64  `protobuf:"varint,2,req" json:"StartTime,omitempty"`
	EndTime          *int64  `protobuf:"varint,3,req" json:"EndTime,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Transaction) Reset()         { *m = Transaction{} }
func (m *Transaction) String() string { return proto.CompactTextString(m) }
func (*Transaction) ProtoMessage()    {}

func (m *Transaction) GetID() uint64 {
	if m != nil && m.ID != nil {
		return *m.ID
	}
	return 0
}

func (m *Transaction) GetStartTime() int64 {
	if m != nil && m.StartTime != nil {
		return *m.StartTime
	}
	return 0
}

func (m *Transaction) GetEndTime() int64 {
	if m != nil && m.EndTime != nil {
		return *m.EndTime
	}
	return 0
}

type Log struct {
	Head             *uint64 `protobuf:"varint,1,opt" json:"Head,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Log) Reset()         { *m = Log{} }
func (m *Log) String() string { return proto.CompactTextString(m) }
func (*Log) ProtoMessage()    {}

func (m *Log) GetHead() uint64 {
	if m != nil && m.Head != nil {
		return *m.Head
	}
	return 0
}

type Segment struct {
	ID               *uint64 `protobuf:"varint,1,req" json:"ID,omitempty"`
	Blocks           *int32  `protobuf:"varint,2,req" json:"Blocks,omitempty"`
	Count            *int32  `protobuf:"varint,3,req" json:"Count,omitempty"`
	Bytes            *int32  `protobuf:"varint,4,req" json:"Bytes,omitempty"`
	Next             *uint64 `protobuf:"varint,5,opt" json:"Next,omitempty"`
	Base             *uint64 `protobuf:"varint,6,opt" json:"Base,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Segment) Reset()         { *m = Segment{} }
func (m *Segment) String() string { return proto.CompactTextString(m) }
func (*Segment) ProtoMessage()    {}

func (m *Segment) GetID() uint64 {
	if m != nil && m.ID != nil {
		return *m.ID
	}
	return 0
}

func (m *Segment) GetBlocks() int32 {
	if m != nil && m.Blocks != nil {
		return *m.Blocks
	}
	return 0
}

func (m *Segment) GetCount() int32 {
	if m != nil && m.Count != nil {
		return *m.Count
	}
	return 0
}

func (m *Segment) GetBytes() int32 {
	if m != nil && m.Bytes != nil {
		return *m.Bytes
	}
	return 0
}

func (m *Segment) GetNext() uint64 {
	if m != nil && m.Next != nil {
		return *m.Next
	}
	return 0
}

func (m *Segment) GetBase() uint64 {
	if m != nil && m.Base != nil {
		return *m.Base
	}
	return 0
}

type Fact struct {
	Added            *bool   `protobuf:"varint,1,req" json:"Added,omitempty"`
	EntityDomain     *string `protobuf:"bytes,2,req" json:"EntityDomain,omitempty"`
	Entity           *string `protobuf:"bytes,3,req" json:"Entity,omitempty"`
	AttributeDomain  *string `protobuf:"bytes,4,req" json:"AttributeDomain,omitempty"`
	Attribute        *string `protobuf:"bytes,5,req" json:"Attribute,omitempty"`
	ValueDomain      *string `protobuf:"bytes,6,opt" json:"ValueDomain,omitempty"`
	Value            *string `protobuf:"bytes,7,req" json:"Value,omitempty"`
	Time             *int64  `protobuf:"varint,8,opt" json:"Time,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Fact) Reset()         { *m = Fact{} }
func (m *Fact) String() string { return proto.CompactTextString(m) }
func (*Fact) ProtoMessage()    {}

func (m *Fact) GetAdded() bool {
	if m != nil && m.Added != nil {
		return *m.Added
	}
	return false
}

func (m *Fact) GetEntityDomain() string {
	if m != nil && m.EntityDomain != nil {
		return *m.EntityDomain
	}
	return ""
}

func (m *Fact) GetEntity() string {
	if m != nil && m.Entity != nil {
		return *m.Entity
	}
	return ""
}

func (m *Fact) GetAttributeDomain() string {
	if m != nil && m.AttributeDomain != nil {
		return *m.AttributeDomain
	}
	return ""
}

func (m *Fact) GetAttribute() string {
	if m != nil && m.Attribute != nil {
		return *m.Attribute
	}
	return ""
}

func (m *Fact) GetValueDomain() string {
	if m != nil && m.ValueDomain != nil {
		return *m.ValueDomain
	}
	return ""
}

func (m *Fact) GetValue() string {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return ""
}

func (m *Fact) GetTime() int64 {
	if m != nil && m.Time != nil {
		return *m.Time
	}
	return 0
}

func init() {
}