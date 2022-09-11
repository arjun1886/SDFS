// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.5
// source: coordinator.proto

package coordinator

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

type CoordinatorInput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data string `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Flag string `protobuf:"bytes,2,opt,name=flag,proto3" json:"flag,omitempty"`
}

func (x *CoordinatorInput) Reset() {
	*x = CoordinatorInput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coordinator_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CoordinatorInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CoordinatorInput) ProtoMessage() {}

func (x *CoordinatorInput) ProtoReflect() protoreflect.Message {
	mi := &file_coordinator_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CoordinatorInput.ProtoReflect.Descriptor instead.
func (*CoordinatorInput) Descriptor() ([]byte, []int) {
	return file_coordinator_proto_rawDescGZIP(), []int{0}
}

func (x *CoordinatorInput) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

func (x *CoordinatorInput) GetFlag() string {
	if x != nil {
		return x.Flag
	}
	return ""
}

type CoordinatorOutput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileName        []string `protobuf:"bytes,1,rep,name=fileName,proto3" json:"fileName,omitempty"`
	Matches         []string `protobuf:"bytes,2,rep,name=matches,proto3" json:"matches,omitempty"`
	TotalMatchCount string   `protobuf:"bytes,3,opt,name=totalMatchCount,proto3" json:"totalMatchCount,omitempty"`
}

func (x *CoordinatorOutput) Reset() {
	*x = CoordinatorOutput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coordinator_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CoordinatorOutput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CoordinatorOutput) ProtoMessage() {}

func (x *CoordinatorOutput) ProtoReflect() protoreflect.Message {
	mi := &file_coordinator_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CoordinatorOutput.ProtoReflect.Descriptor instead.
func (*CoordinatorOutput) Descriptor() ([]byte, []int) {
	return file_coordinator_proto_rawDescGZIP(), []int{1}
}

func (x *CoordinatorOutput) GetFileName() []string {
	if x != nil {
		return x.FileName
	}
	return nil
}

func (x *CoordinatorOutput) GetMatches() []string {
	if x != nil {
		return x.Matches
	}
	return nil
}

func (x *CoordinatorOutput) GetTotalMatchCount() string {
	if x != nil {
		return x.TotalMatchCount
	}
	return ""
}

var File_coordinator_proto protoreflect.FileDescriptor

var file_coordinator_proto_rawDesc = []byte{
	0x0a, 0x11, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72,
	0x22, 0x3a, 0x0a, 0x10, 0x43, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x49,
	0x6e, 0x70, 0x75, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x6c, 0x61, 0x67,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x6c, 0x61, 0x67, 0x22, 0x73, 0x0a, 0x11,
	0x43, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x4f, 0x75, 0x74, 0x70, 0x75,
	0x74, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07,
	0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x12, 0x28, 0x0a, 0x0f, 0x74, 0x6f, 0x74, 0x61, 0x6c,
	0x4d, 0x61, 0x74, 0x63, 0x68, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0f, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x32, 0x6f, 0x0a, 0x12, 0x43, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x59, 0x0a, 0x16, 0x46, 0x65, 0x74, 0x63, 0x68,
	0x43, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x4f, 0x75, 0x74, 0x70, 0x75,
	0x74, 0x12, 0x1d, 0x2e, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x2e,
	0x43, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x49, 0x6e, 0x70, 0x75, 0x74,
	0x1a, 0x1e, 0x2e, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x43,
	0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74,
	0x22, 0x00, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x2f, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61,
	0x74, 0x6f, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_coordinator_proto_rawDescOnce sync.Once
	file_coordinator_proto_rawDescData = file_coordinator_proto_rawDesc
)

func file_coordinator_proto_rawDescGZIP() []byte {
	file_coordinator_proto_rawDescOnce.Do(func() {
		file_coordinator_proto_rawDescData = protoimpl.X.CompressGZIP(file_coordinator_proto_rawDescData)
	})
	return file_coordinator_proto_rawDescData
}

var file_coordinator_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_coordinator_proto_goTypes = []interface{}{
	(*CoordinatorInput)(nil),  // 0: coordinator.CoordinatorInput
	(*CoordinatorOutput)(nil), // 1: coordinator.CoordinatorOutput
}
var file_coordinator_proto_depIdxs = []int32{
	0, // 0: coordinator.CoordinatorService.FetchCoordinatorOutput:input_type -> coordinator.CoordinatorInput
	1, // 1: coordinator.CoordinatorService.FetchCoordinatorOutput:output_type -> coordinator.CoordinatorOutput
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_coordinator_proto_init() }
func file_coordinator_proto_init() {
	if File_coordinator_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_coordinator_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CoordinatorInput); i {
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
		file_coordinator_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CoordinatorOutput); i {
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
			RawDescriptor: file_coordinator_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_coordinator_proto_goTypes,
		DependencyIndexes: file_coordinator_proto_depIdxs,
		MessageInfos:      file_coordinator_proto_msgTypes,
	}.Build()
	File_coordinator_proto = out.File
	file_coordinator_proto_rawDesc = nil
	file_coordinator_proto_goTypes = nil
	file_coordinator_proto_depIdxs = nil
}
