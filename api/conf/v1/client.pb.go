// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: api/conf/v1/client.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
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

// 客户端
type Client struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Http          *Client_HTTP           `protobuf:"bytes,1,opt,name=http,proto3" json:"http,omitempty"` // HTTP服务
	Grpc          *Client_GRPC           `protobuf:"bytes,2,opt,name=grpc,proto3" json:"grpc,omitempty"` // GRPC服务
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Client) Reset() {
	*x = Client{}
	mi := &file_api_conf_v1_client_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Client) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Client) ProtoMessage() {}

func (x *Client) ProtoReflect() protoreflect.Message {
	mi := &file_api_conf_v1_client_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Client.ProtoReflect.Descriptor instead.
func (*Client) Descriptor() ([]byte, []int) {
	return file_api_conf_v1_client_proto_rawDescGZIP(), []int{0}
}

func (x *Client) GetHttp() *Client_HTTP {
	if x != nil {
		return x.Http
	}
	return nil
}

func (x *Client) GetGrpc() *Client_GRPC {
	if x != nil {
		return x.Grpc
	}
	return nil
}

// HTTP
type Client_HTTP struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Timeout       *durationpb.Duration   `protobuf:"bytes,1,opt,name=timeout,proto3" json:"timeout,omitempty"`       // 超时时间
	Middleware    *Middleware            `protobuf:"bytes,2,opt,name=middleware,proto3" json:"middleware,omitempty"` // 中间件
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Client_HTTP) Reset() {
	*x = Client_HTTP{}
	mi := &file_api_conf_v1_client_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Client_HTTP) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Client_HTTP) ProtoMessage() {}

func (x *Client_HTTP) ProtoReflect() protoreflect.Message {
	mi := &file_api_conf_v1_client_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Client_HTTP.ProtoReflect.Descriptor instead.
func (*Client_HTTP) Descriptor() ([]byte, []int) {
	return file_api_conf_v1_client_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Client_HTTP) GetTimeout() *durationpb.Duration {
	if x != nil {
		return x.Timeout
	}
	return nil
}

func (x *Client_HTTP) GetMiddleware() *Middleware {
	if x != nil {
		return x.Middleware
	}
	return nil
}

// gPRC
type Client_GRPC struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Timeout       *durationpb.Duration   `protobuf:"bytes,1,opt,name=timeout,proto3" json:"timeout,omitempty"`       // 超时时间
	Middleware    *Middleware            `protobuf:"bytes,2,opt,name=middleware,proto3" json:"middleware,omitempty"` // 中间件
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Client_GRPC) Reset() {
	*x = Client_GRPC{}
	mi := &file_api_conf_v1_client_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Client_GRPC) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Client_GRPC) ProtoMessage() {}

func (x *Client_GRPC) ProtoReflect() protoreflect.Message {
	mi := &file_api_conf_v1_client_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Client_GRPC.ProtoReflect.Descriptor instead.
func (*Client_GRPC) Descriptor() ([]byte, []int) {
	return file_api_conf_v1_client_proto_rawDescGZIP(), []int{0, 1}
}

func (x *Client_GRPC) GetTimeout() *durationpb.Duration {
	if x != nil {
		return x.Timeout
	}
	return nil
}

func (x *Client_GRPC) GetMiddleware() *Middleware {
	if x != nil {
		return x.Middleware
	}
	return nil
}

var File_api_conf_v1_client_proto protoreflect.FileDescriptor

const file_api_conf_v1_client_proto_rawDesc = "" +
	"\n" +
	"\x18api/conf/v1/client.proto\x12\x04conf\x1a\x1capi/conf/v1/middleware.proto\x1a\x1egoogle/protobuf/duration.proto\"\xb4\x02\n" +
	"\x06Client\x12%\n" +
	"\x04http\x18\x01 \x01(\v2\x11.conf.Client.HTTPR\x04http\x12%\n" +
	"\x04grpc\x18\x02 \x01(\v2\x11.conf.Client.GRPCR\x04grpc\x1am\n" +
	"\x04HTTP\x123\n" +
	"\atimeout\x18\x01 \x01(\v2\x19.google.protobuf.DurationR\atimeout\x120\n" +
	"\n" +
	"middleware\x18\x02 \x01(\v2\x10.conf.MiddlewareR\n" +
	"middleware\x1am\n" +
	"\x04GRPC\x123\n" +
	"\atimeout\x18\x01 \x01(\v2\x19.google.protobuf.DurationR\atimeout\x120\n" +
	"\n" +
	"middleware\x18\x02 \x01(\v2\x10.conf.MiddlewareR\n" +
	"middlewareB\x83\x01\n" +
	"\bcom.confB\vClientProtoP\x01Z:github.com/fzf-labs/kratos-contrib/api/conf/v1/api/conf/v1\xa2\x02\x03CXX\xaa\x02\x04Conf\xca\x02\x04Conf\xe2\x02\x10Conf\\GPBMetadata\xea\x02\x04Confb\x06proto3"

var (
	file_api_conf_v1_client_proto_rawDescOnce sync.Once
	file_api_conf_v1_client_proto_rawDescData []byte
)

func file_api_conf_v1_client_proto_rawDescGZIP() []byte {
	file_api_conf_v1_client_proto_rawDescOnce.Do(func() {
		file_api_conf_v1_client_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_conf_v1_client_proto_rawDesc), len(file_api_conf_v1_client_proto_rawDesc)))
	})
	return file_api_conf_v1_client_proto_rawDescData
}

var file_api_conf_v1_client_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_api_conf_v1_client_proto_goTypes = []any{
	(*Client)(nil),              // 0: conf.Client
	(*Client_HTTP)(nil),         // 1: conf.Client.HTTP
	(*Client_GRPC)(nil),         // 2: conf.Client.GRPC
	(*durationpb.Duration)(nil), // 3: google.protobuf.Duration
	(*Middleware)(nil),          // 4: conf.Middleware
}
var file_api_conf_v1_client_proto_depIdxs = []int32{
	1, // 0: conf.Client.http:type_name -> conf.Client.HTTP
	2, // 1: conf.Client.grpc:type_name -> conf.Client.GRPC
	3, // 2: conf.Client.HTTP.timeout:type_name -> google.protobuf.Duration
	4, // 3: conf.Client.HTTP.middleware:type_name -> conf.Middleware
	3, // 4: conf.Client.GRPC.timeout:type_name -> google.protobuf.Duration
	4, // 5: conf.Client.GRPC.middleware:type_name -> conf.Middleware
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_api_conf_v1_client_proto_init() }
func file_api_conf_v1_client_proto_init() {
	if File_api_conf_v1_client_proto != nil {
		return
	}
	file_api_conf_v1_middleware_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_conf_v1_client_proto_rawDesc), len(file_api_conf_v1_client_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_conf_v1_client_proto_goTypes,
		DependencyIndexes: file_api_conf_v1_client_proto_depIdxs,
		MessageInfos:      file_api_conf_v1_client_proto_msgTypes,
	}.Build()
	File_api_conf_v1_client_proto = out.File
	file_api_conf_v1_client_proto_goTypes = nil
	file_api_conf_v1_client_proto_depIdxs = nil
}
