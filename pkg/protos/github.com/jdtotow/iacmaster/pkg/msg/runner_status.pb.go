// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.19.1
// source: runner_status.proto

package msg

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RunnerStatus struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Status        string                 `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	Address       string                 `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RunnerStatus) Reset() {
	*x = RunnerStatus{}
	mi := &file_runner_status_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RunnerStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RunnerStatus) ProtoMessage() {}

func (x *RunnerStatus) ProtoReflect() protoreflect.Message {
	mi := &file_runner_status_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RunnerStatus.ProtoReflect.Descriptor instead.
func (*RunnerStatus) Descriptor() ([]byte, []int) {
	return file_runner_status_proto_rawDescGZIP(), []int{0}
}

func (x *RunnerStatus) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RunnerStatus) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *RunnerStatus) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

var File_runner_status_proto protoreflect.FileDescriptor

var file_runner_status_proto_rawDesc = string([]byte{
	0x0a, 0x13, 0x72, 0x75, 0x6e, 0x6e, 0x65, 0x72, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x54, 0x0a, 0x0c, 0x52, 0x75,
	0x6e, 0x6e, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a,
	0x64, 0x74, 0x6f, 0x74, 0x6f, 0x77, 0x2f, 0x69, 0x61, 0x63, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72,
	0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x73, 0x67, 0x3b, 0x6d, 0x73, 0x67, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_runner_status_proto_rawDescOnce sync.Once
	file_runner_status_proto_rawDescData []byte
)

func file_runner_status_proto_rawDescGZIP() []byte {
	file_runner_status_proto_rawDescOnce.Do(func() {
		file_runner_status_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_runner_status_proto_rawDesc), len(file_runner_status_proto_rawDesc)))
	})
	return file_runner_status_proto_rawDescData
}

var file_runner_status_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_runner_status_proto_goTypes = []any{
	(*RunnerStatus)(nil), // 0: msg.RunnerStatus
}
var file_runner_status_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_runner_status_proto_init() }
func file_runner_status_proto_init() {
	if File_runner_status_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_runner_status_proto_rawDesc), len(file_runner_status_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_runner_status_proto_goTypes,
		DependencyIndexes: file_runner_status_proto_depIdxs,
		MessageInfos:      file_runner_status_proto_msgTypes,
	}.Build()
	File_runner_status_proto = out.File
	file_runner_status_proto_goTypes = nil
	file_runner_status_proto_depIdxs = nil
}
