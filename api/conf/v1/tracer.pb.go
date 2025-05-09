// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: api/conf/v1/tracer.proto

package v1

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

// 链路追踪
type Tracer struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Batcher       string                 `protobuf:"bytes,1,opt,name=batcher,proto3" json:"batcher,omitempty"`    // stdout,otlphttp, otlpgrpc
	Endpoint      string                 `protobuf:"bytes,2,opt,name=endpoint,proto3" json:"endpoint,omitempty"`  // 端口
	Insecure      bool                   `protobuf:"varint,3,opt,name=insecure,proto3" json:"insecure,omitempty"` // 是否不安全
	Sampler       float64                `protobuf:"fixed64,4,opt,name=sampler,proto3" json:"sampler,omitempty"`  // 采样率，默认：1.0
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Tracer) Reset() {
	*x = Tracer{}
	mi := &file_api_conf_v1_tracer_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Tracer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tracer) ProtoMessage() {}

func (x *Tracer) ProtoReflect() protoreflect.Message {
	mi := &file_api_conf_v1_tracer_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tracer.ProtoReflect.Descriptor instead.
func (*Tracer) Descriptor() ([]byte, []int) {
	return file_api_conf_v1_tracer_proto_rawDescGZIP(), []int{0}
}

func (x *Tracer) GetBatcher() string {
	if x != nil {
		return x.Batcher
	}
	return ""
}

func (x *Tracer) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *Tracer) GetInsecure() bool {
	if x != nil {
		return x.Insecure
	}
	return false
}

func (x *Tracer) GetSampler() float64 {
	if x != nil {
		return x.Sampler
	}
	return 0
}

var File_api_conf_v1_tracer_proto protoreflect.FileDescriptor

const file_api_conf_v1_tracer_proto_rawDesc = "" +
	"\n" +
	"\x18api/conf/v1/tracer.proto\x12\x04conf\"t\n" +
	"\x06Tracer\x12\x18\n" +
	"\abatcher\x18\x01 \x01(\tR\abatcher\x12\x1a\n" +
	"\bendpoint\x18\x02 \x01(\tR\bendpoint\x12\x1a\n" +
	"\binsecure\x18\x03 \x01(\bR\binsecure\x12\x18\n" +
	"\asampler\x18\x04 \x01(\x01R\asamplerB\x83\x01\n" +
	"\bcom.confB\vTracerProtoP\x01Z:github.com/fzf-labs/kratos-contrib/api/conf/v1/api/conf/v1\xa2\x02\x03CXX\xaa\x02\x04Conf\xca\x02\x04Conf\xe2\x02\x10Conf\\GPBMetadata\xea\x02\x04Confb\x06proto3"

var (
	file_api_conf_v1_tracer_proto_rawDescOnce sync.Once
	file_api_conf_v1_tracer_proto_rawDescData []byte
)

func file_api_conf_v1_tracer_proto_rawDescGZIP() []byte {
	file_api_conf_v1_tracer_proto_rawDescOnce.Do(func() {
		file_api_conf_v1_tracer_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_conf_v1_tracer_proto_rawDesc), len(file_api_conf_v1_tracer_proto_rawDesc)))
	})
	return file_api_conf_v1_tracer_proto_rawDescData
}

var file_api_conf_v1_tracer_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_api_conf_v1_tracer_proto_goTypes = []any{
	(*Tracer)(nil), // 0: conf.Tracer
}
var file_api_conf_v1_tracer_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_conf_v1_tracer_proto_init() }
func file_api_conf_v1_tracer_proto_init() {
	if File_api_conf_v1_tracer_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_conf_v1_tracer_proto_rawDesc), len(file_api_conf_v1_tracer_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_conf_v1_tracer_proto_goTypes,
		DependencyIndexes: file_api_conf_v1_tracer_proto_depIdxs,
		MessageInfos:      file_api_conf_v1_tracer_proto_msgTypes,
	}.Build()
	File_api_conf_v1_tracer_proto = out.File
	file_api_conf_v1_tracer_proto_goTypes = nil
	file_api_conf_v1_tracer_proto_depIdxs = nil
}
