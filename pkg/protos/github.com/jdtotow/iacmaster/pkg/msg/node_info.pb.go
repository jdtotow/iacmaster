// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.19.1
// source: node_info.proto

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

type NodeInfo struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	NodeStatus    uint32                 `protobuf:"varint,2,opt,name=nodeStatus,proto3" json:"nodeStatus,omitempty"`
	NodeType      uint32                 `protobuf:"varint,3,opt,name=nodeType,proto3" json:"nodeType,omitempty"`
	Addr          string                 `protobuf:"bytes,4,opt,name=addr,proto3" json:"addr,omitempty"`
	Deployments   []string               `protobuf:"bytes,5,rep,name=deployments,proto3" json:"deployments,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NodeInfo) Reset() {
	*x = NodeInfo{}
	mi := &file_node_info_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NodeInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeInfo) ProtoMessage() {}

func (x *NodeInfo) ProtoReflect() protoreflect.Message {
	mi := &file_node_info_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeInfo.ProtoReflect.Descriptor instead.
func (*NodeInfo) Descriptor() ([]byte, []int) {
	return file_node_info_proto_rawDescGZIP(), []int{0}
}

func (x *NodeInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *NodeInfo) GetNodeStatus() uint32 {
	if x != nil {
		return x.NodeStatus
	}
	return 0
}

func (x *NodeInfo) GetNodeType() uint32 {
	if x != nil {
		return x.NodeType
	}
	return 0
}

func (x *NodeInfo) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *NodeInfo) GetDeployments() []string {
	if x != nil {
		return x.Deployments
	}
	return nil
}

var File_node_info_proto protoreflect.FileDescriptor

var file_node_info_proto_rawDesc = string([]byte{
	0x0a, 0x0f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x90, 0x01, 0x0a, 0x08, 0x4e, 0x6f, 0x64, 0x65, 0x49,
	0x6e, 0x66, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x6e, 0x6f, 0x64, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x6e, 0x6f, 0x64,
	0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x64, 0x64, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x70, 0x6c, 0x6f,
	0x79, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65,
	0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x64, 0x74, 0x6f, 0x74, 0x6f, 0x77, 0x2f,
	0x69, 0x61, 0x63, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x73,
	0x67, 0x3b, 0x6d, 0x73, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_node_info_proto_rawDescOnce sync.Once
	file_node_info_proto_rawDescData []byte
)

func file_node_info_proto_rawDescGZIP() []byte {
	file_node_info_proto_rawDescOnce.Do(func() {
		file_node_info_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_node_info_proto_rawDesc), len(file_node_info_proto_rawDesc)))
	})
	return file_node_info_proto_rawDescData
}

var file_node_info_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_node_info_proto_goTypes = []any{
	(*NodeInfo)(nil), // 0: msg.NodeInfo
}
var file_node_info_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_node_info_proto_init() }
func file_node_info_proto_init() {
	if File_node_info_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_node_info_proto_rawDesc), len(file_node_info_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_node_info_proto_goTypes,
		DependencyIndexes: file_node_info_proto_depIdxs,
		MessageInfos:      file_node_info_proto_msgTypes,
	}.Build()
	File_node_info_proto = out.File
	file_node_info_proto_goTypes = nil
	file_node_info_proto_depIdxs = nil
}
