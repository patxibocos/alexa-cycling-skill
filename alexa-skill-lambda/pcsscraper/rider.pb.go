// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v4.22.0
// source: rider.proto

package pcsscraper

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Rider struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                 string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	FirstName          string                 `protobuf:"bytes,2,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName           string                 `protobuf:"bytes,3,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	Country            string                 `protobuf:"bytes,4,opt,name=country,proto3" json:"country,omitempty"`
	BirthDate          *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=birth_date,json=birthDate,proto3,oneof" json:"birth_date,omitempty"`
	Photo              string                 `protobuf:"bytes,6,opt,name=photo,proto3" json:"photo,omitempty"`
	Website            *string                `protobuf:"bytes,7,opt,name=website,proto3,oneof" json:"website,omitempty"`
	BirthPlace         *string                `protobuf:"bytes,8,opt,name=birth_place,json=birthPlace,proto3,oneof" json:"birth_place,omitempty"`
	Weight             *uint32                `protobuf:"varint,9,opt,name=weight,proto3,oneof" json:"weight,omitempty"`
	Height             *uint32                `protobuf:"varint,10,opt,name=height,proto3,oneof" json:"height,omitempty"`
	UciRankingPosition *uint32                `protobuf:"varint,11,opt,name=uciRankingPosition,proto3,oneof" json:"uciRankingPosition,omitempty"`
}

func (x *Rider) Reset() {
	*x = Rider{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rider_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Rider) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rider) ProtoMessage() {}

func (x *Rider) ProtoReflect() protoreflect.Message {
	mi := &file_rider_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rider.ProtoReflect.Descriptor instead.
func (*Rider) Descriptor() ([]byte, []int) {
	return file_rider_proto_rawDescGZIP(), []int{0}
}

func (x *Rider) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Rider) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *Rider) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *Rider) GetCountry() string {
	if x != nil {
		return x.Country
	}
	return ""
}

func (x *Rider) GetBirthDate() *timestamppb.Timestamp {
	if x != nil {
		return x.BirthDate
	}
	return nil
}

func (x *Rider) GetPhoto() string {
	if x != nil {
		return x.Photo
	}
	return ""
}

func (x *Rider) GetWebsite() string {
	if x != nil && x.Website != nil {
		return *x.Website
	}
	return ""
}

func (x *Rider) GetBirthPlace() string {
	if x != nil && x.BirthPlace != nil {
		return *x.BirthPlace
	}
	return ""
}

func (x *Rider) GetWeight() uint32 {
	if x != nil && x.Weight != nil {
		return *x.Weight
	}
	return 0
}

func (x *Rider) GetHeight() uint32 {
	if x != nil && x.Height != nil {
		return *x.Height
	}
	return 0
}

func (x *Rider) GetUciRankingPosition() uint32 {
	if x != nil && x.UciRankingPosition != nil {
		return *x.UciRankingPosition
	}
	return 0
}

var File_rider_proto protoreflect.FileDescriptor

var file_rider_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x72, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x28, 0x69,
	0x6f, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x70, 0x61, 0x74, 0x78, 0x69, 0x62, 0x6f,
	0x63, 0x6f, 0x73, 0x2e, 0x70, 0x63, 0x73, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xcf, 0x03, 0x0a, 0x05, 0x52, 0x69, 0x64,
	0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x69, 0x72, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x61, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x3e, 0x0a, 0x0a, 0x62, 0x69, 0x72, 0x74,
	0x68, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x00, 0x52, 0x09, 0x62, 0x69, 0x72, 0x74,
	0x68, 0x44, 0x61, 0x74, 0x65, 0x88, 0x01, 0x01, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x68, 0x6f, 0x74,
	0x6f, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x12, 0x1d,
	0x0a, 0x07, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x01, 0x52, 0x07, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x88, 0x01, 0x01, 0x12, 0x24, 0x0a,
	0x0b, 0x62, 0x69, 0x72, 0x74, 0x68, 0x5f, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x02, 0x52, 0x0a, 0x62, 0x69, 0x72, 0x74, 0x68, 0x50, 0x6c, 0x61, 0x63, 0x65,
	0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x0d, 0x48, 0x03, 0x52, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x88, 0x01, 0x01,
	0x12, 0x1b, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0d,
	0x48, 0x04, 0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x88, 0x01, 0x01, 0x12, 0x33, 0x0a,
	0x12, 0x75, 0x63, 0x69, 0x52, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x50, 0x6f, 0x73, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x05, 0x52, 0x12, 0x75, 0x63, 0x69,
	0x52, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x88,
	0x01, 0x01, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x62, 0x69, 0x72, 0x74, 0x68, 0x5f, 0x64, 0x61, 0x74,
	0x65, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x42, 0x0e, 0x0a,
	0x0c, 0x5f, 0x62, 0x69, 0x72, 0x74, 0x68, 0x5f, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x42, 0x09, 0x0a,
	0x07, 0x5f, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x68, 0x65, 0x69,
	0x67, 0x68, 0x74, 0x42, 0x15, 0x0a, 0x13, 0x5f, 0x75, 0x63, 0x69, 0x52, 0x61, 0x6e, 0x6b, 0x69,
	0x6e, 0x67, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x22, 0x5a, 0x20, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x61, 0x74, 0x78, 0x69, 0x62, 0x6f,
	0x63, 0x6f, 0x73, 0x2f, 0x70, 0x63, 0x73, 0x73, 0x63, 0x72, 0x61, 0x70, 0x65, 0x72, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rider_proto_rawDescOnce sync.Once
	file_rider_proto_rawDescData = file_rider_proto_rawDesc
)

func file_rider_proto_rawDescGZIP() []byte {
	file_rider_proto_rawDescOnce.Do(func() {
		file_rider_proto_rawDescData = protoimpl.X.CompressGZIP(file_rider_proto_rawDescData)
	})
	return file_rider_proto_rawDescData
}

var file_rider_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_rider_proto_goTypes = []interface{}{
	(*Rider)(nil),                 // 0: io.github.patxibocos.pcsscraper.protobuf.Rider
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_rider_proto_depIdxs = []int32{
	1, // 0: io.github.patxibocos.pcsscraper.protobuf.Rider.birth_date:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_rider_proto_init() }
func file_rider_proto_init() {
	if File_rider_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rider_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Rider); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_rider_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rider_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rider_proto_goTypes,
		DependencyIndexes: file_rider_proto_depIdxs,
		MessageInfos:      file_rider_proto_msgTypes,
	}.Build()
	File_rider_proto = out.File
	file_rider_proto_rawDesc = nil
	file_rider_proto_goTypes = nil
	file_rider_proto_depIdxs = nil
}
