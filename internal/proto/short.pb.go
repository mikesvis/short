// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v3.21.12
// source: proto/short.proto

package protobuf

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Request and response messages
type SaveURLRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *SaveURLRequest) Reset() {
	*x = SaveURLRequest{}
	mi := &file_proto_short_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SaveURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveURLRequest) ProtoMessage() {}

func (x *SaveURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_short_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveURLRequest.ProtoReflect.Descriptor instead.
func (*SaveURLRequest) Descriptor() ([]byte, []int) {
	return file_proto_short_proto_rawDescGZIP(), []int{0}
}

func (x *SaveURLRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type SaveURLResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShortURL string `protobuf:"bytes,1,opt,name=shortURL,proto3" json:"shortURL,omitempty"`
}

func (x *SaveURLResponse) Reset() {
	*x = SaveURLResponse{}
	mi := &file_proto_short_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SaveURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveURLResponse) ProtoMessage() {}

func (x *SaveURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_short_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveURLResponse.ProtoReflect.Descriptor instead.
func (*SaveURLResponse) Descriptor() ([]byte, []int) {
	return file_proto_short_proto_rawDescGZIP(), []int{1}
}

func (x *SaveURLResponse) GetShortURL() string {
	if x != nil {
		return x.ShortURL
	}
	return ""
}

type DeleteBatchURLRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Urls []string `protobuf:"bytes,1,rep,name=urls,proto3" json:"urls,omitempty"`
}

func (x *DeleteBatchURLRequest) Reset() {
	*x = DeleteBatchURLRequest{}
	mi := &file_proto_short_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteBatchURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteBatchURLRequest) ProtoMessage() {}

func (x *DeleteBatchURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_short_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteBatchURLRequest.ProtoReflect.Descriptor instead.
func (*DeleteBatchURLRequest) Descriptor() ([]byte, []int) {
	return file_proto_short_proto_rawDescGZIP(), []int{2}
}

func (x *DeleteBatchURLRequest) GetUrls() []string {
	if x != nil {
		return x.Urls
	}
	return nil
}

type GetURLByIDRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShortURL string `protobuf:"bytes,1,opt,name=shortURL,proto3" json:"shortURL,omitempty"`
}

func (x *GetURLByIDRequest) Reset() {
	*x = GetURLByIDRequest{}
	mi := &file_proto_short_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetURLByIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetURLByIDRequest) ProtoMessage() {}

func (x *GetURLByIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_short_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetURLByIDRequest.ProtoReflect.Descriptor instead.
func (*GetURLByIDRequest) Descriptor() ([]byte, []int) {
	return file_proto_short_proto_rawDescGZIP(), []int{3}
}

func (x *GetURLByIDRequest) GetShortURL() string {
	if x != nil {
		return x.ShortURL
	}
	return ""
}

type GetURLByIDResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OriginalURL string `protobuf:"bytes,1,opt,name=originalURL,proto3" json:"originalURL,omitempty"`
}

func (x *GetURLByIDResponse) Reset() {
	*x = GetURLByIDResponse{}
	mi := &file_proto_short_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetURLByIDResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetURLByIDResponse) ProtoMessage() {}

func (x *GetURLByIDResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_short_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetURLByIDResponse.ProtoReflect.Descriptor instead.
func (*GetURLByIDResponse) Descriptor() ([]byte, []int) {
	return file_proto_short_proto_rawDescGZIP(), []int{4}
}

func (x *GetURLByIDResponse) GetOriginalURL() string {
	if x != nil {
		return x.OriginalURL
	}
	return ""
}

type GetURLByUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Urls []*URLByUserResponseElement `protobuf:"bytes,1,rep,name=urls,proto3" json:"urls,omitempty"`
}

func (x *GetURLByUserResponse) Reset() {
	*x = GetURLByUserResponse{}
	mi := &file_proto_short_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetURLByUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetURLByUserResponse) ProtoMessage() {}

func (x *GetURLByUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_short_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetURLByUserResponse.ProtoReflect.Descriptor instead.
func (*GetURLByUserResponse) Descriptor() ([]byte, []int) {
	return file_proto_short_proto_rawDescGZIP(), []int{5}
}

func (x *GetURLByUserResponse) GetUrls() []*URLByUserResponseElement {
	if x != nil {
		return x.Urls
	}
	return nil
}

type ShortenBatchURLRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Urls []*ShortenBatchURLRequestElement `protobuf:"bytes,1,rep,name=urls,proto3" json:"urls,omitempty"`
}

func (x *ShortenBatchURLRequest) Reset() {
	*x = ShortenBatchURLRequest{}
	mi := &file_proto_short_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ShortenBatchURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortenBatchURLRequest) ProtoMessage() {}

func (x *ShortenBatchURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_short_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortenBatchURLRequest.ProtoReflect.Descriptor instead.
func (*ShortenBatchURLRequest) Descriptor() ([]byte, []int) {
	return file_proto_short_proto_rawDescGZIP(), []int{6}
}

func (x *ShortenBatchURLRequest) GetUrls() []*ShortenBatchURLRequestElement {
	if x != nil {
		return x.Urls
	}
	return nil
}

type ShortenBatchURLResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Urls []*ShortenBatchURLResponseElement `protobuf:"bytes,1,rep,name=urls,proto3" json:"urls,omitempty"`
}

func (x *ShortenBatchURLResponse) Reset() {
	*x = ShortenBatchURLResponse{}
	mi := &file_proto_short_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ShortenBatchURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortenBatchURLResponse) ProtoMessage() {}

func (x *ShortenBatchURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_short_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortenBatchURLResponse.ProtoReflect.Descriptor instead.
func (*ShortenBatchURLResponse) Descriptor() ([]byte, []int) {
	return file_proto_short_proto_rawDescGZIP(), []int{7}
}

func (x *ShortenBatchURLResponse) GetUrls() []*ShortenBatchURLResponseElement {
	if x != nil {
		return x.Urls
	}
	return nil
}

type GetStatsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Urls  int32 `protobuf:"varint,1,opt,name=urls,proto3" json:"urls,omitempty"`
	Users int32 `protobuf:"varint,2,opt,name=users,proto3" json:"users,omitempty"`
}

func (x *GetStatsResponse) Reset() {
	*x = GetStatsResponse{}
	mi := &file_proto_short_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetStatsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStatsResponse) ProtoMessage() {}

func (x *GetStatsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_short_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStatsResponse.ProtoReflect.Descriptor instead.
func (*GetStatsResponse) Descriptor() ([]byte, []int) {
	return file_proto_short_proto_rawDescGZIP(), []int{8}
}

func (x *GetStatsResponse) GetUrls() int32 {
	if x != nil {
		return x.Urls
	}
	return 0
}

func (x *GetStatsResponse) GetUsers() int32 {
	if x != nil {
		return x.Users
	}
	return 0
}

// Data structure messages
type ShortenBatchURLRequestElement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CorrelationId string `protobuf:"bytes,1,opt,name=correlationId,proto3" json:"correlationId,omitempty"`
	OriginalURL   string `protobuf:"bytes,2,opt,name=originalURL,proto3" json:"originalURL,omitempty"`
}

func (x *ShortenBatchURLRequestElement) Reset() {
	*x = ShortenBatchURLRequestElement{}
	mi := &file_proto_short_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ShortenBatchURLRequestElement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortenBatchURLRequestElement) ProtoMessage() {}

func (x *ShortenBatchURLRequestElement) ProtoReflect() protoreflect.Message {
	mi := &file_proto_short_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortenBatchURLRequestElement.ProtoReflect.Descriptor instead.
func (*ShortenBatchURLRequestElement) Descriptor() ([]byte, []int) {
	return file_proto_short_proto_rawDescGZIP(), []int{9}
}

func (x *ShortenBatchURLRequestElement) GetCorrelationId() string {
	if x != nil {
		return x.CorrelationId
	}
	return ""
}

func (x *ShortenBatchURLRequestElement) GetOriginalURL() string {
	if x != nil {
		return x.OriginalURL
	}
	return ""
}

type ShortenBatchURLResponseElement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CorrelationId string `protobuf:"bytes,1,opt,name=correlationId,proto3" json:"correlationId,omitempty"`
	ShortURL      string `protobuf:"bytes,2,opt,name=shortURL,proto3" json:"shortURL,omitempty"`
}

func (x *ShortenBatchURLResponseElement) Reset() {
	*x = ShortenBatchURLResponseElement{}
	mi := &file_proto_short_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ShortenBatchURLResponseElement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortenBatchURLResponseElement) ProtoMessage() {}

func (x *ShortenBatchURLResponseElement) ProtoReflect() protoreflect.Message {
	mi := &file_proto_short_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortenBatchURLResponseElement.ProtoReflect.Descriptor instead.
func (*ShortenBatchURLResponseElement) Descriptor() ([]byte, []int) {
	return file_proto_short_proto_rawDescGZIP(), []int{10}
}

func (x *ShortenBatchURLResponseElement) GetCorrelationId() string {
	if x != nil {
		return x.CorrelationId
	}
	return ""
}

func (x *ShortenBatchURLResponseElement) GetShortURL() string {
	if x != nil {
		return x.ShortURL
	}
	return ""
}

type URLByUserResponseElement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShortURL    string `protobuf:"bytes,1,opt,name=shortURL,proto3" json:"shortURL,omitempty"`
	OriginalURL string `protobuf:"bytes,2,opt,name=originalURL,proto3" json:"originalURL,omitempty"`
}

func (x *URLByUserResponseElement) Reset() {
	*x = URLByUserResponseElement{}
	mi := &file_proto_short_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *URLByUserResponseElement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*URLByUserResponseElement) ProtoMessage() {}

func (x *URLByUserResponseElement) ProtoReflect() protoreflect.Message {
	mi := &file_proto_short_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use URLByUserResponseElement.ProtoReflect.Descriptor instead.
func (*URLByUserResponseElement) Descriptor() ([]byte, []int) {
	return file_proto_short_proto_rawDescGZIP(), []int{11}
}

func (x *URLByUserResponseElement) GetShortURL() string {
	if x != nil {
		return x.ShortURL
	}
	return ""
}

func (x *URLByUserResponseElement) GetOriginalURL() string {
	if x != nil {
		return x.OriginalURL
	}
	return ""
}

var File_proto_short_proto protoreflect.FileDescriptor

var file_proto_short_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x22, 0x0a, 0x0e, 0x53, 0x61, 0x76, 0x65, 0x55,
	0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x2d, 0x0a, 0x0f, 0x53,
	0x61, 0x76, 0x65, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x22, 0x2b, 0x0a, 0x15, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x72, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x04, 0x75, 0x72, 0x6c, 0x73, 0x22, 0x2f, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x55, 0x52,
	0x4c, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08,
	0x73, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x73, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x22, 0x36, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x55,
	0x52, 0x4c, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x20,
	0x0a, 0x0b, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x52, 0x4c, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x52, 0x4c,
	0x22, 0x4b, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x33, 0x0a, 0x04, 0x75, 0x72, 0x6c, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x2e, 0x55,
	0x52, 0x4c, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x04, 0x75, 0x72, 0x6c, 0x73, 0x22, 0x52, 0x0a,
	0x16, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x55, 0x52, 0x4c,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x38, 0x0a, 0x04, 0x75, 0x72, 0x6c, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x2e, 0x53, 0x68,
	0x6f, 0x72, 0x74, 0x65, 0x6e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x04, 0x75, 0x72, 0x6c,
	0x73, 0x22, 0x54, 0x0a, 0x17, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x42, 0x61, 0x74, 0x63,
	0x68, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a, 0x04,
	0x75, 0x72, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x73, 0x68, 0x6f,
	0x72, 0x74, 0x2e, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x55,
	0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e,
	0x74, 0x52, 0x04, 0x75, 0x72, 0x6c, 0x73, 0x22, 0x3c, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x53, 0x74,
	0x61, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75,
	0x72, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x75, 0x72, 0x6c, 0x73, 0x12,
	0x14, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05,
	0x75, 0x73, 0x65, 0x72, 0x73, 0x22, 0x67, 0x0a, 0x1d, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e,
	0x42, 0x61, 0x74, 0x63, 0x68, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x45,
	0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63,
	0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b,
	0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x52, 0x4c, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x52, 0x4c, 0x22, 0x62,
	0x0a, 0x1e, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x55, 0x52,
	0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x12, 0x24, 0x0a, 0x0d, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x55,
	0x52, 0x4c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x55,
	0x52, 0x4c, 0x22, 0x58, 0x0a, 0x18, 0x55, 0x52, 0x4c, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1a,
	0x0a, 0x08, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x12, 0x20, 0x0a, 0x0b, 0x6f, 0x72,
	0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x52, 0x4c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x52, 0x4c, 0x32, 0xa7, 0x03, 0x0a,
	0x0c, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x38, 0x0a,
	0x07, 0x53, 0x61, 0x76, 0x65, 0x55, 0x52, 0x4c, 0x12, 0x15, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74,
	0x2e, 0x53, 0x61, 0x76, 0x65, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x16, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x55, 0x52, 0x4c, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x50, 0x0a, 0x0f, 0x53, 0x68, 0x6f, 0x72, 0x74,
	0x65, 0x6e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x55, 0x52, 0x4c, 0x12, 0x1d, 0x2e, 0x73, 0x68, 0x6f,
	0x72, 0x74, 0x2e, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x55,
	0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x73, 0x68, 0x6f, 0x72,
	0x74, 0x2e, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x55, 0x52,
	0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x0a, 0x47, 0x65, 0x74,
	0x55, 0x52, 0x4c, 0x42, 0x79, 0x49, 0x44, 0x12, 0x18, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x2e,
	0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x19, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x52, 0x4c,
	0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x0c,
	0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x12, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1b, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x2e, 0x47, 0x65, 0x74,
	0x55, 0x52, 0x4c, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x46, 0x0a, 0x0e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68,
	0x55, 0x52, 0x4c, 0x12, 0x1c, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x2e, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x3b, 0x0a, 0x08, 0x47, 0x65, 0x74,
	0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x17, 0x2e,
	0x73, 0x68, 0x6f, 0x72, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_short_proto_rawDescOnce sync.Once
	file_proto_short_proto_rawDescData = file_proto_short_proto_rawDesc
)

func file_proto_short_proto_rawDescGZIP() []byte {
	file_proto_short_proto_rawDescOnce.Do(func() {
		file_proto_short_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_short_proto_rawDescData)
	})
	return file_proto_short_proto_rawDescData
}

var file_proto_short_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_proto_short_proto_goTypes = []any{
	(*SaveURLRequest)(nil),                 // 0: short.SaveURLRequest
	(*SaveURLResponse)(nil),                // 1: short.SaveURLResponse
	(*DeleteBatchURLRequest)(nil),          // 2: short.DeleteBatchURLRequest
	(*GetURLByIDRequest)(nil),              // 3: short.GetURLByIDRequest
	(*GetURLByIDResponse)(nil),             // 4: short.GetURLByIDResponse
	(*GetURLByUserResponse)(nil),           // 5: short.GetURLByUserResponse
	(*ShortenBatchURLRequest)(nil),         // 6: short.ShortenBatchURLRequest
	(*ShortenBatchURLResponse)(nil),        // 7: short.ShortenBatchURLResponse
	(*GetStatsResponse)(nil),               // 8: short.GetStatsResponse
	(*ShortenBatchURLRequestElement)(nil),  // 9: short.ShortenBatchURLRequestElement
	(*ShortenBatchURLResponseElement)(nil), // 10: short.ShortenBatchURLResponseElement
	(*URLByUserResponseElement)(nil),       // 11: short.URLByUserResponseElement
	(*emptypb.Empty)(nil),                  // 12: google.protobuf.Empty
}
var file_proto_short_proto_depIdxs = []int32{
	11, // 0: short.GetURLByUserResponse.urls:type_name -> short.URLByUserResponseElement
	9,  // 1: short.ShortenBatchURLRequest.urls:type_name -> short.ShortenBatchURLRequestElement
	10, // 2: short.ShortenBatchURLResponse.urls:type_name -> short.ShortenBatchURLResponseElement
	0,  // 3: short.ShortService.SaveURL:input_type -> short.SaveURLRequest
	6,  // 4: short.ShortService.ShortenBatchURL:input_type -> short.ShortenBatchURLRequest
	3,  // 5: short.ShortService.GetURLByID:input_type -> short.GetURLByIDRequest
	12, // 6: short.ShortService.GetURLByUser:input_type -> google.protobuf.Empty
	2,  // 7: short.ShortService.DeleteBatchURL:input_type -> short.DeleteBatchURLRequest
	12, // 8: short.ShortService.GetStats:input_type -> google.protobuf.Empty
	1,  // 9: short.ShortService.SaveURL:output_type -> short.SaveURLResponse
	7,  // 10: short.ShortService.ShortenBatchURL:output_type -> short.ShortenBatchURLResponse
	4,  // 11: short.ShortService.GetURLByID:output_type -> short.GetURLByIDResponse
	5,  // 12: short.ShortService.GetURLByUser:output_type -> short.GetURLByUserResponse
	12, // 13: short.ShortService.DeleteBatchURL:output_type -> google.protobuf.Empty
	8,  // 14: short.ShortService.GetStats:output_type -> short.GetStatsResponse
	9,  // [9:15] is the sub-list for method output_type
	3,  // [3:9] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_proto_short_proto_init() }
func file_proto_short_proto_init() {
	if File_proto_short_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_short_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_short_proto_goTypes,
		DependencyIndexes: file_proto_short_proto_depIdxs,
		MessageInfos:      file_proto_short_proto_msgTypes,
	}.Build()
	File_proto_short_proto = out.File
	file_proto_short_proto_rawDesc = nil
	file_proto_short_proto_goTypes = nil
	file_proto_short_proto_depIdxs = nil
}
