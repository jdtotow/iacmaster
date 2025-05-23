// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.19.1
// source: deployment.proto

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

type GitData struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Url           string                 `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Token         string                 `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	TokenUsername string                 `protobuf:"bytes,3,opt,name=tokenUsername,proto3" json:"tokenUsername,omitempty"`
	Revision      string                 `protobuf:"bytes,4,opt,name=revision,proto3" json:"revision,omitempty"`
	ProxyUrl      string                 `protobuf:"bytes,5,opt,name=proxyUrl,proto3" json:"proxyUrl,omitempty"`
	ProxyUsername string                 `protobuf:"bytes,6,opt,name=proxyUsername,proto3" json:"proxyUsername,omitempty"`
	ProxyPassword string                 `protobuf:"bytes,7,opt,name=proxyPassword,proto3" json:"proxyPassword,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GitData) Reset() {
	*x = GitData{}
	mi := &file_deployment_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GitData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GitData) ProtoMessage() {}

func (x *GitData) ProtoReflect() protoreflect.Message {
	mi := &file_deployment_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GitData.ProtoReflect.Descriptor instead.
func (*GitData) Descriptor() ([]byte, []int) {
	return file_deployment_proto_rawDescGZIP(), []int{0}
}

func (x *GitData) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *GitData) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *GitData) GetTokenUsername() string {
	if x != nil {
		return x.TokenUsername
	}
	return ""
}

func (x *GitData) GetRevision() string {
	if x != nil {
		return x.Revision
	}
	return ""
}

func (x *GitData) GetProxyUrl() string {
	if x != nil {
		return x.ProxyUrl
	}
	return ""
}

func (x *GitData) GetProxyUsername() string {
	if x != nil {
		return x.ProxyUsername
	}
	return ""
}

func (x *GitData) GetProxyPassword() string {
	if x != nil {
		return x.ProxyPassword
	}
	return ""
}

type Deployment struct {
	state                   protoimpl.MessageState `protogen:"open.v1"`
	WorkingDir              string                 `protobuf:"bytes,1,opt,name=WorkingDir,proto3" json:"WorkingDir,omitempty"`
	HomeFolder              string                 `protobuf:"bytes,2,opt,name=HomeFolder,proto3" json:"HomeFolder,omitempty"`
	CloudDestination        string                 `protobuf:"bytes,3,opt,name=CloudDestination,proto3" json:"CloudDestination,omitempty"`
	TerraformVersion        string                 `protobuf:"bytes,4,opt,name=TerraformVersion,proto3" json:"TerraformVersion,omitempty"`
	EnvironmentParameters   map[string]string      `protobuf:"bytes,5,rep,name=EnvironmentParameters,proto3" json:"EnvironmentParameters,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Status                  string                 `protobuf:"bytes,6,opt,name=Status,proto3" json:"Status,omitempty"`
	GitData                 *GitData               `protobuf:"bytes,7,opt,name=GitData,proto3" json:"GitData,omitempty"`
	EnvironmentID           string                 `protobuf:"bytes,8,opt,name=EnvironmentID,proto3" json:"EnvironmentID,omitempty"`
	Error                   string                 `protobuf:"bytes,9,opt,name=Error,proto3" json:"Error,omitempty"`
	Activities              []string               `protobuf:"bytes,10,rep,name=Activities,proto3" json:"Activities,omitempty"`
	IaCArtifactType         string                 `protobuf:"bytes,11,opt,name=IaCArtifactType,proto3" json:"IaCArtifactType,omitempty"`
	DetectDrift             bool                   `protobuf:"varint,12,opt,name=DetectDrift,proto3" json:"DetectDrift,omitempty"`
	AutoRedeployOnGitChange bool                   `protobuf:"varint,13,opt,name=AutoRedeployOnGitChange,proto3" json:"AutoRedeployOnGitChange,omitempty"`
	Action                  string                 `protobuf:"bytes,14,opt,name=Action,proto3" json:"Action,omitempty"`
	unknownFields           protoimpl.UnknownFields
	sizeCache               protoimpl.SizeCache
}

func (x *Deployment) Reset() {
	*x = Deployment{}
	mi := &file_deployment_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Deployment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Deployment) ProtoMessage() {}

func (x *Deployment) ProtoReflect() protoreflect.Message {
	mi := &file_deployment_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Deployment.ProtoReflect.Descriptor instead.
func (*Deployment) Descriptor() ([]byte, []int) {
	return file_deployment_proto_rawDescGZIP(), []int{1}
}

func (x *Deployment) GetWorkingDir() string {
	if x != nil {
		return x.WorkingDir
	}
	return ""
}

func (x *Deployment) GetHomeFolder() string {
	if x != nil {
		return x.HomeFolder
	}
	return ""
}

func (x *Deployment) GetCloudDestination() string {
	if x != nil {
		return x.CloudDestination
	}
	return ""
}

func (x *Deployment) GetTerraformVersion() string {
	if x != nil {
		return x.TerraformVersion
	}
	return ""
}

func (x *Deployment) GetEnvironmentParameters() map[string]string {
	if x != nil {
		return x.EnvironmentParameters
	}
	return nil
}

func (x *Deployment) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Deployment) GetGitData() *GitData {
	if x != nil {
		return x.GitData
	}
	return nil
}

func (x *Deployment) GetEnvironmentID() string {
	if x != nil {
		return x.EnvironmentID
	}
	return ""
}

func (x *Deployment) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *Deployment) GetActivities() []string {
	if x != nil {
		return x.Activities
	}
	return nil
}

func (x *Deployment) GetIaCArtifactType() string {
	if x != nil {
		return x.IaCArtifactType
	}
	return ""
}

func (x *Deployment) GetDetectDrift() bool {
	if x != nil {
		return x.DetectDrift
	}
	return false
}

func (x *Deployment) GetAutoRedeployOnGitChange() bool {
	if x != nil {
		return x.AutoRedeployOnGitChange
	}
	return false
}

func (x *Deployment) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

var File_deployment_proto protoreflect.FileDescriptor

var file_deployment_proto_rawDesc = string([]byte{
	0x0a, 0x10, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x03, 0x6d, 0x73, 0x67, 0x22, 0xdb, 0x01, 0x0a, 0x07, 0x47, 0x69, 0x74, 0x44,
	0x61, 0x74, 0x61, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x24, 0x0a, 0x0d, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a,
	0x08, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x55, 0x72, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x55, 0x72, 0x6c, 0x12, 0x24, 0x0a, 0x0d, 0x70, 0x72, 0x6f,
	0x78, 0x79, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0d, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x24, 0x0a, 0x0d, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x50, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x8a, 0x05, 0x0a, 0x0a, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x57, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x44,
	0x69, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x57, 0x6f, 0x72, 0x6b, 0x69, 0x6e,
	0x67, 0x44, 0x69, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x48, 0x6f, 0x6d, 0x65, 0x46, 0x6f, 0x6c, 0x64,
	0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x48, 0x6f, 0x6d, 0x65, 0x46, 0x6f,
	0x6c, 0x64, 0x65, 0x72, 0x12, 0x2a, 0x0a, 0x10, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x44, 0x65, 0x73,
	0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10,
	0x43, 0x6c, 0x6f, 0x75, 0x64, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x2a, 0x0a, 0x10, 0x54, 0x65, 0x72, 0x72, 0x61, 0x66, 0x6f, 0x72, 0x6d, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x54, 0x65, 0x72, 0x72,
	0x61, 0x66, 0x6f, 0x72, 0x6d, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x60, 0x0a, 0x15,
	0x45, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x50, 0x61, 0x72, 0x61, 0x6d,
	0x65, 0x74, 0x65, 0x72, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x6d, 0x73,
	0x67, 0x2e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x45, 0x6e, 0x76,
	0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65,
	0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x15, 0x45, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e,
	0x6d, 0x65, 0x6e, 0x74, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x16,
	0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x26, 0x0a, 0x07, 0x47, 0x69, 0x74, 0x44, 0x61, 0x74,
	0x61, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x6d, 0x73, 0x67, 0x2e, 0x47, 0x69,
	0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x07, 0x47, 0x69, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x24,
	0x0a, 0x0d, 0x45, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x45, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65,
	0x6e, 0x74, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x41, 0x63,
	0x74, 0x69, 0x76, 0x69, 0x74, 0x69, 0x65, 0x73, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a,
	0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x69, 0x65, 0x73, 0x12, 0x28, 0x0a, 0x0f, 0x49, 0x61,
	0x43, 0x41, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0f, 0x49, 0x61, 0x43, 0x41, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x44, 0x65, 0x74, 0x65, 0x63, 0x74, 0x44, 0x72,
	0x69, 0x66, 0x74, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x44, 0x65, 0x74, 0x65, 0x63,
	0x74, 0x44, 0x72, 0x69, 0x66, 0x74, 0x12, 0x38, 0x0a, 0x17, 0x41, 0x75, 0x74, 0x6f, 0x52, 0x65,
	0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x4f, 0x6e, 0x47, 0x69, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x08, 0x52, 0x17, 0x41, 0x75, 0x74, 0x6f, 0x52, 0x65, 0x64,
	0x65, 0x70, 0x6c, 0x6f, 0x79, 0x4f, 0x6e, 0x47, 0x69, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x48, 0x0a, 0x1a, 0x45, 0x6e, 0x76, 0x69,
	0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x6a, 0x64, 0x74, 0x6f, 0x74, 0x6f, 0x77, 0x2f, 0x69, 0x61, 0x63, 0x6d, 0x61, 0x73, 0x74,
	0x65, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x73, 0x67, 0x3b, 0x6d, 0x73, 0x67, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_deployment_proto_rawDescOnce sync.Once
	file_deployment_proto_rawDescData []byte
)

func file_deployment_proto_rawDescGZIP() []byte {
	file_deployment_proto_rawDescOnce.Do(func() {
		file_deployment_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_deployment_proto_rawDesc), len(file_deployment_proto_rawDesc)))
	})
	return file_deployment_proto_rawDescData
}

var file_deployment_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_deployment_proto_goTypes = []any{
	(*GitData)(nil),    // 0: msg.GitData
	(*Deployment)(nil), // 1: msg.Deployment
	nil,                // 2: msg.Deployment.EnvironmentParametersEntry
}
var file_deployment_proto_depIdxs = []int32{
	2, // 0: msg.Deployment.EnvironmentParameters:type_name -> msg.Deployment.EnvironmentParametersEntry
	0, // 1: msg.Deployment.GitData:type_name -> msg.GitData
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_deployment_proto_init() }
func file_deployment_proto_init() {
	if File_deployment_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_deployment_proto_rawDesc), len(file_deployment_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_deployment_proto_goTypes,
		DependencyIndexes: file_deployment_proto_depIdxs,
		MessageInfos:      file_deployment_proto_msgTypes,
	}.Build()
	File_deployment_proto = out.File
	file_deployment_proto_goTypes = nil
	file_deployment_proto_depIdxs = nil
}
