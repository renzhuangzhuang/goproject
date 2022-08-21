// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.5
// source: MR.proto

package protos

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

//定义发送消息
type MrRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *MrRequest) Reset() {
	*x = MrRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_MR_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MrRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MrRequest) ProtoMessage() {}

func (x *MrRequest) ProtoReflect() protoreflect.Message {
	mi := &file_MR_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MrRequest.ProtoReflect.Descriptor instead.
func (*MrRequest) Descriptor() ([]byte, []int) {
	return file_MR_proto_rawDescGZIP(), []int{0}
}

func (x *MrRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

//定义接收消息
type MrResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *MrResponse) Reset() {
	*x = MrResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_MR_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MrResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MrResponse) ProtoMessage() {}

func (x *MrResponse) ProtoReflect() protoreflect.Message {
	mi := &file_MR_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MrResponse.ProtoReflect.Descriptor instead.
func (*MrResponse) Descriptor() ([]byte, []int) {
	return file_MR_proto_rawDescGZIP(), []int{1}
}

func (x *MrResponse) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_MR_proto protoreflect.FileDescriptor

var file_MR_proto_rawDesc = []byte{
	0x0a, 0x08, 0x4d, 0x52, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1f, 0x0a, 0x09, 0x4d, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x20, 0x0a, 0x0a, 0x4d,
	0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0x60, 0x0a,
	0x09, 0x4d, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x27, 0x0a, 0x0a, 0x47, 0x65,
	0x74, 0x53, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x0a, 0x2e, 0x4d, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x0b, 0x2e, 0x4d, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x30, 0x01, 0x12, 0x2a, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x42, 0x69, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x12, 0x0a, 0x2e, 0x4d, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0b,
	0x2e, 0x4d, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x42,
	0x12, 0x5a, 0x10, 0x6d, 0x61, 0x70, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_MR_proto_rawDescOnce sync.Once
	file_MR_proto_rawDescData = file_MR_proto_rawDesc
)

func file_MR_proto_rawDescGZIP() []byte {
	file_MR_proto_rawDescOnce.Do(func() {
		file_MR_proto_rawDescData = protoimpl.X.CompressGZIP(file_MR_proto_rawDescData)
	})
	return file_MR_proto_rawDescData
}

var file_MR_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_MR_proto_goTypes = []interface{}{
	(*MrRequest)(nil),  // 0: MrRequest
	(*MrResponse)(nil), // 1: MrResponse
}
var file_MR_proto_depIdxs = []int32{
	0, // 0: MrService.GetSStream:input_type -> MrRequest
	0, // 1: MrService.GetBiStream:input_type -> MrRequest
	1, // 2: MrService.GetSStream:output_type -> MrResponse
	1, // 3: MrService.GetBiStream:output_type -> MrResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_MR_proto_init() }
func file_MR_proto_init() {
	if File_MR_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_MR_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MrRequest); i {
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
		file_MR_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MrResponse); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_MR_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_MR_proto_goTypes,
		DependencyIndexes: file_MR_proto_depIdxs,
		MessageInfos:      file_MR_proto_msgTypes,
	}.Build()
	File_MR_proto = out.File
	file_MR_proto_rawDesc = nil
	file_MR_proto_goTypes = nil
	file_MR_proto_depIdxs = nil
}
