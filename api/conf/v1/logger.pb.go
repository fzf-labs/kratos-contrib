// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: api/conf/v1/logger.proto

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

// 日志
type Logger struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Type          string                 `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"` // 类型 std zap zerolog
	Zap           *Logger_Zap            `protobuf:"bytes,2,opt,name=zap,proto3" json:"zap,omitempty"`   // Zap
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Logger) Reset() {
	*x = Logger{}
	mi := &file_api_conf_v1_logger_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Logger) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Logger) ProtoMessage() {}

func (x *Logger) ProtoReflect() protoreflect.Message {
	mi := &file_api_conf_v1_logger_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Logger.ProtoReflect.Descriptor instead.
func (*Logger) Descriptor() ([]byte, []int) {
	return file_api_conf_v1_logger_proto_rawDescGZIP(), []int{0}
}

func (x *Logger) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Logger) GetZap() *Logger_Zap {
	if x != nil {
		return x.Zap
	}
	return nil
}

// Zap
type Logger_Zap struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Filename      string                 `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`      // 文件名
	Level         string                 `protobuf:"bytes,2,opt,name=level,proto3" json:"level,omitempty"`            // 日志级别
	MaxSize       int32                  `protobuf:"varint,3,opt,name=maxSize,proto3" json:"maxSize,omitempty"`       // 最大大小
	MaxAge        int32                  `protobuf:"varint,4,opt,name=maxAge,proto3" json:"maxAge,omitempty"`         // 最大年龄
	MaxBackups    int32                  `protobuf:"varint,5,opt,name=maxBackups,proto3" json:"maxBackups,omitempty"` // 最大备份
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Logger_Zap) Reset() {
	*x = Logger_Zap{}
	mi := &file_api_conf_v1_logger_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Logger_Zap) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Logger_Zap) ProtoMessage() {}

func (x *Logger_Zap) ProtoReflect() protoreflect.Message {
	mi := &file_api_conf_v1_logger_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Logger_Zap.ProtoReflect.Descriptor instead.
func (*Logger_Zap) Descriptor() ([]byte, []int) {
	return file_api_conf_v1_logger_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Logger_Zap) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *Logger_Zap) GetLevel() string {
	if x != nil {
		return x.Level
	}
	return ""
}

func (x *Logger_Zap) GetMaxSize() int32 {
	if x != nil {
		return x.MaxSize
	}
	return 0
}

func (x *Logger_Zap) GetMaxAge() int32 {
	if x != nil {
		return x.MaxAge
	}
	return 0
}

func (x *Logger_Zap) GetMaxBackups() int32 {
	if x != nil {
		return x.MaxBackups
	}
	return 0
}

var File_api_conf_v1_logger_proto protoreflect.FileDescriptor

const file_api_conf_v1_logger_proto_rawDesc = "" +
	"\n" +
	"\x18api/conf/v1/logger.proto\x12\x04conf\"\xcc\x01\n" +
	"\x06Logger\x12\x12\n" +
	"\x04type\x18\x01 \x01(\tR\x04type\x12\"\n" +
	"\x03zap\x18\x02 \x01(\v2\x10.conf.Logger.ZapR\x03zap\x1a\x89\x01\n" +
	"\x03Zap\x12\x1a\n" +
	"\bfilename\x18\x01 \x01(\tR\bfilename\x12\x14\n" +
	"\x05level\x18\x02 \x01(\tR\x05level\x12\x18\n" +
	"\amaxSize\x18\x03 \x01(\x05R\amaxSize\x12\x16\n" +
	"\x06maxAge\x18\x04 \x01(\x05R\x06maxAge\x12\x1e\n" +
	"\n" +
	"maxBackups\x18\x05 \x01(\x05R\n" +
	"maxBackupsB\x83\x01\n" +
	"\bcom.confB\vLoggerProtoP\x01Z:github.com/fzf-labs/kratos-contrib/api/conf/v1/api/conf/v1\xa2\x02\x03CXX\xaa\x02\x04Conf\xca\x02\x04Conf\xe2\x02\x10Conf\\GPBMetadata\xea\x02\x04Confb\x06proto3"

var (
	file_api_conf_v1_logger_proto_rawDescOnce sync.Once
	file_api_conf_v1_logger_proto_rawDescData []byte
)

func file_api_conf_v1_logger_proto_rawDescGZIP() []byte {
	file_api_conf_v1_logger_proto_rawDescOnce.Do(func() {
		file_api_conf_v1_logger_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_conf_v1_logger_proto_rawDesc), len(file_api_conf_v1_logger_proto_rawDesc)))
	})
	return file_api_conf_v1_logger_proto_rawDescData
}

var file_api_conf_v1_logger_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_conf_v1_logger_proto_goTypes = []any{
	(*Logger)(nil),     // 0: conf.Logger
	(*Logger_Zap)(nil), // 1: conf.Logger.Zap
}
var file_api_conf_v1_logger_proto_depIdxs = []int32{
	1, // 0: conf.Logger.zap:type_name -> conf.Logger.Zap
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_conf_v1_logger_proto_init() }
func file_api_conf_v1_logger_proto_init() {
	if File_api_conf_v1_logger_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_conf_v1_logger_proto_rawDesc), len(file_api_conf_v1_logger_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_conf_v1_logger_proto_goTypes,
		DependencyIndexes: file_api_conf_v1_logger_proto_depIdxs,
		MessageInfos:      file_api_conf_v1_logger_proto_msgTypes,
	}.Build()
	File_api_conf_v1_logger_proto = out.File
	file_api_conf_v1_logger_proto_goTypes = nil
	file_api_conf_v1_logger_proto_depIdxs = nil
}
