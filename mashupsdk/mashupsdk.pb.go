// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.15.8
// source: mashupsdk/mashupsdk.proto

package mashupsdk

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

type MashupEmpty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AuthToken string `protobuf:"bytes,1,opt,name=authToken,proto3" json:"authToken,omitempty"`
}

func (x *MashupEmpty) Reset() {
	*x = MashupEmpty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mashupsdk_mashupsdk_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MashupEmpty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MashupEmpty) ProtoMessage() {}

func (x *MashupEmpty) ProtoReflect() protoreflect.Message {
	mi := &file_mashupsdk_mashupsdk_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MashupEmpty.ProtoReflect.Descriptor instead.
func (*MashupEmpty) Descriptor() ([]byte, []int) {
	return file_mashupsdk_mashupsdk_proto_rawDescGZIP(), []int{0}
}

func (x *MashupEmpty) GetAuthToken() string {
	if x != nil {
		return x.AuthToken
	}
	return ""
}

// The response message with mashup credentials
type MashupConnectionConfigs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AuthToken   string `protobuf:"bytes,1,opt,name=authToken,proto3" json:"authToken,omitempty"`
	CallerToken string `protobuf:"bytes,2,opt,name=callerToken,proto3" json:"callerToken,omitempty"`
	Port        int64  `protobuf:"varint,3,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *MashupConnectionConfigs) Reset() {
	*x = MashupConnectionConfigs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mashupsdk_mashupsdk_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MashupConnectionConfigs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MashupConnectionConfigs) ProtoMessage() {}

func (x *MashupConnectionConfigs) ProtoReflect() protoreflect.Message {
	mi := &file_mashupsdk_mashupsdk_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MashupConnectionConfigs.ProtoReflect.Descriptor instead.
func (*MashupConnectionConfigs) Descriptor() ([]byte, []int) {
	return file_mashupsdk_mashupsdk_proto_rawDescGZIP(), []int{1}
}

func (x *MashupConnectionConfigs) GetAuthToken() string {
	if x != nil {
		return x.AuthToken
	}
	return ""
}

func (x *MashupConnectionConfigs) GetCallerToken() string {
	if x != nil {
		return x.CallerToken
	}
	return ""
}

func (x *MashupConnectionConfigs) GetPort() int64 {
	if x != nil {
		return x.Port
	}
	return 0
}

// The query message with display position information.
type MashupDisplayHint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Focused bool  `protobuf:"varint,1,opt,name=focused,proto3" json:"focused,omitempty"`
	Xpos    int64 `protobuf:"varint,2,opt,name=xpos,proto3" json:"xpos,omitempty"`
	Ypos    int64 `protobuf:"varint,3,opt,name=ypos,proto3" json:"ypos,omitempty"`
	Width   int64 `protobuf:"varint,4,opt,name=width,proto3" json:"width,omitempty"`
	Height  int64 `protobuf:"varint,5,opt,name=height,proto3" json:"height,omitempty"`
}

func (x *MashupDisplayHint) Reset() {
	*x = MashupDisplayHint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mashupsdk_mashupsdk_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MashupDisplayHint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MashupDisplayHint) ProtoMessage() {}

func (x *MashupDisplayHint) ProtoReflect() protoreflect.Message {
	mi := &file_mashupsdk_mashupsdk_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MashupDisplayHint.ProtoReflect.Descriptor instead.
func (*MashupDisplayHint) Descriptor() ([]byte, []int) {
	return file_mashupsdk_mashupsdk_proto_rawDescGZIP(), []int{2}
}

func (x *MashupDisplayHint) GetFocused() bool {
	if x != nil {
		return x.Focused
	}
	return false
}

func (x *MashupDisplayHint) GetXpos() int64 {
	if x != nil {
		return x.Xpos
	}
	return 0
}

func (x *MashupDisplayHint) GetYpos() int64 {
	if x != nil {
		return x.Ypos
	}
	return 0
}

func (x *MashupDisplayHint) GetWidth() int64 {
	if x != nil {
		return x.Width
	}
	return 0
}

func (x *MashupDisplayHint) GetHeight() int64 {
	if x != nil {
		return x.Height
	}
	return 0
}

type MashupDisplayBundle struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AuthToken         string             `protobuf:"bytes,1,opt,name=authToken,proto3" json:"authToken,omitempty"`
	MashupDisplayHint *MashupDisplayHint `protobuf:"bytes,2,opt,name=mashupDisplayHint,proto3" json:"mashupDisplayHint,omitempty"`
}

func (x *MashupDisplayBundle) Reset() {
	*x = MashupDisplayBundle{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mashupsdk_mashupsdk_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MashupDisplayBundle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MashupDisplayBundle) ProtoMessage() {}

func (x *MashupDisplayBundle) ProtoReflect() protoreflect.Message {
	mi := &file_mashupsdk_mashupsdk_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MashupDisplayBundle.ProtoReflect.Descriptor instead.
func (*MashupDisplayBundle) Descriptor() ([]byte, []int) {
	return file_mashupsdk_mashupsdk_proto_rawDescGZIP(), []int{3}
}

func (x *MashupDisplayBundle) GetAuthToken() string {
	if x != nil {
		return x.AuthToken
	}
	return ""
}

func (x *MashupDisplayBundle) GetMashupDisplayHint() *MashupDisplayHint {
	if x != nil {
		return x.MashupDisplayHint
	}
	return nil
}

// The response message containing the any messages to log
type ShutdownReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ShutdownReply) Reset() {
	*x = ShutdownReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mashupsdk_mashupsdk_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShutdownReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShutdownReply) ProtoMessage() {}

func (x *ShutdownReply) ProtoReflect() protoreflect.Message {
	mi := &file_mashupsdk_mashupsdk_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShutdownReply.ProtoReflect.Descriptor instead.
func (*ShutdownReply) Descriptor() ([]byte, []int) {
	return file_mashupsdk_mashupsdk_proto_rawDescGZIP(), []int{4}
}

func (x *ShutdownReply) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type MashupDetailedElement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Basisid       int64               `protobuf:"varint,1,opt,name=basisid,proto3" json:"basisid,omitempty"`
	Id            int64               `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	State         *MashupElementState `protobuf:"bytes,3,opt,name=state,proto3" json:"state,omitempty"`
	Name          string              `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	Alias         string              `protobuf:"bytes,5,opt,name=alias,proto3" json:"alias,omitempty"`
	Description   string              `protobuf:"bytes,6,opt,name=description,proto3" json:"description,omitempty"`
	Renderer      string              `protobuf:"bytes,7,opt,name=renderer,proto3" json:"renderer,omitempty"`
	Colabrenderer string              `protobuf:"bytes,8,opt,name=colabrenderer,proto3" json:"colabrenderer,omitempty"`
	Genre         string              `protobuf:"bytes,9,opt,name=genre,proto3" json:"genre,omitempty"`
	Subgenre      string              `protobuf:"bytes,10,opt,name=subgenre,proto3" json:"subgenre,omitempty"`
	Parentids     []int64             `protobuf:"varint,11,rep,packed,name=parentids,proto3" json:"parentids,omitempty"`
	Childids      []int64             `protobuf:"varint,12,rep,packed,name=childids,proto3" json:"childids,omitempty"`
}

func (x *MashupDetailedElement) Reset() {
	*x = MashupDetailedElement{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mashupsdk_mashupsdk_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MashupDetailedElement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MashupDetailedElement) ProtoMessage() {}

func (x *MashupDetailedElement) ProtoReflect() protoreflect.Message {
	mi := &file_mashupsdk_mashupsdk_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MashupDetailedElement.ProtoReflect.Descriptor instead.
func (*MashupDetailedElement) Descriptor() ([]byte, []int) {
	return file_mashupsdk_mashupsdk_proto_rawDescGZIP(), []int{5}
}

func (x *MashupDetailedElement) GetBasisid() int64 {
	if x != nil {
		return x.Basisid
	}
	return 0
}

func (x *MashupDetailedElement) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *MashupDetailedElement) GetState() *MashupElementState {
	if x != nil {
		return x.State
	}
	return nil
}

func (x *MashupDetailedElement) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *MashupDetailedElement) GetAlias() string {
	if x != nil {
		return x.Alias
	}
	return ""
}

func (x *MashupDetailedElement) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *MashupDetailedElement) GetRenderer() string {
	if x != nil {
		return x.Renderer
	}
	return ""
}

func (x *MashupDetailedElement) GetColabrenderer() string {
	if x != nil {
		return x.Colabrenderer
	}
	return ""
}

func (x *MashupDetailedElement) GetGenre() string {
	if x != nil {
		return x.Genre
	}
	return ""
}

func (x *MashupDetailedElement) GetSubgenre() string {
	if x != nil {
		return x.Subgenre
	}
	return ""
}

func (x *MashupDetailedElement) GetParentids() []int64 {
	if x != nil {
		return x.Parentids
	}
	return nil
}

func (x *MashupDetailedElement) GetChildids() []int64 {
	if x != nil {
		return x.Childids
	}
	return nil
}

type MashupDetailedElementBundle struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AuthToken        string                   `protobuf:"bytes,1,opt,name=authToken,proto3" json:"authToken,omitempty"`
	DetailedElements []*MashupDetailedElement `protobuf:"bytes,2,rep,name=detailedElements,proto3" json:"detailedElements,omitempty"`
}

func (x *MashupDetailedElementBundle) Reset() {
	*x = MashupDetailedElementBundle{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mashupsdk_mashupsdk_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MashupDetailedElementBundle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MashupDetailedElementBundle) ProtoMessage() {}

func (x *MashupDetailedElementBundle) ProtoReflect() protoreflect.Message {
	mi := &file_mashupsdk_mashupsdk_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MashupDetailedElementBundle.ProtoReflect.Descriptor instead.
func (*MashupDetailedElementBundle) Descriptor() ([]byte, []int) {
	return file_mashupsdk_mashupsdk_proto_rawDescGZIP(), []int{6}
}

func (x *MashupDetailedElementBundle) GetAuthToken() string {
	if x != nil {
		return x.AuthToken
	}
	return ""
}

func (x *MashupDetailedElementBundle) GetDetailedElements() []*MashupDetailedElement {
	if x != nil {
		return x.DetailedElements
	}
	return nil
}

type MashupElementState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	State int64 `protobuf:"varint,2,opt,name=state,proto3" json:"state,omitempty"`
}

func (x *MashupElementState) Reset() {
	*x = MashupElementState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mashupsdk_mashupsdk_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MashupElementState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MashupElementState) ProtoMessage() {}

func (x *MashupElementState) ProtoReflect() protoreflect.Message {
	mi := &file_mashupsdk_mashupsdk_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MashupElementState.ProtoReflect.Descriptor instead.
func (*MashupElementState) Descriptor() ([]byte, []int) {
	return file_mashupsdk_mashupsdk_proto_rawDescGZIP(), []int{7}
}

func (x *MashupElementState) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *MashupElementState) GetState() int64 {
	if x != nil {
		return x.State
	}
	return 0
}

type MashupElementStateBundle struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AuthToken     string                `protobuf:"bytes,1,opt,name=authToken,proto3" json:"authToken,omitempty"`
	ElementStates []*MashupElementState `protobuf:"bytes,2,rep,name=elementStates,proto3" json:"elementStates,omitempty"`
}

func (x *MashupElementStateBundle) Reset() {
	*x = MashupElementStateBundle{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mashupsdk_mashupsdk_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MashupElementStateBundle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MashupElementStateBundle) ProtoMessage() {}

func (x *MashupElementStateBundle) ProtoReflect() protoreflect.Message {
	mi := &file_mashupsdk_mashupsdk_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MashupElementStateBundle.ProtoReflect.Descriptor instead.
func (*MashupElementStateBundle) Descriptor() ([]byte, []int) {
	return file_mashupsdk_mashupsdk_proto_rawDescGZIP(), []int{8}
}

func (x *MashupElementStateBundle) GetAuthToken() string {
	if x != nil {
		return x.AuthToken
	}
	return ""
}

func (x *MashupElementStateBundle) GetElementStates() []*MashupElementState {
	if x != nil {
		return x.ElementStates
	}
	return nil
}

var File_mashupsdk_mashupsdk_proto protoreflect.FileDescriptor

var file_mashupsdk_mashupsdk_proto_rawDesc = []byte{
	0x0a, 0x19, 0x6d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x73, 0x64, 0x6b, 0x2f, 0x6d, 0x61, 0x73, 0x68,
	0x75, 0x70, 0x73, 0x64, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x6d, 0x61, 0x73,
	0x68, 0x75, 0x70, 0x73, 0x64, 0x6b, 0x22, 0x2b, 0x0a, 0x0b, 0x4d, 0x61, 0x73, 0x68, 0x75, 0x70,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x75, 0x74, 0x68, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x22, 0x6d, 0x0a, 0x17, 0x4d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x43, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x12, 0x1c,
	0x0a, 0x09, 0x61, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x61, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x20, 0x0a, 0x0b,
	0x63, 0x61, 0x6c, 0x6c, 0x65, 0x72, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x63, 0x61, 0x6c, 0x6c, 0x65, 0x72, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x12,
	0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x70, 0x6f,
	0x72, 0x74, 0x22, 0x83, 0x01, 0x0a, 0x11, 0x4d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x44, 0x69, 0x73,
	0x70, 0x6c, 0x61, 0x79, 0x48, 0x69, 0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x66, 0x6f, 0x63, 0x75,
	0x73, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x66, 0x6f, 0x63, 0x75, 0x73,
	0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x78, 0x70, 0x6f, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x04, 0x78, 0x70, 0x6f, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x79, 0x70, 0x6f, 0x73, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x79, 0x70, 0x6f, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x77, 0x69,
	0x64, 0x74, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x77, 0x69, 0x64, 0x74, 0x68,
	0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x22, 0x7f, 0x0a, 0x13, 0x4d, 0x61, 0x73, 0x68,
	0x75, 0x70, 0x44, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x42, 0x75, 0x6e, 0x64, 0x6c, 0x65, 0x12,
	0x1c, 0x0a, 0x09, 0x61, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x61, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x4a, 0x0a,
	0x11, 0x6d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x44, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x48, 0x69,
	0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x6d, 0x61, 0x73, 0x68, 0x75,
	0x70, 0x73, 0x64, 0x6b, 0x2e, 0x4d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x44, 0x69, 0x73, 0x70, 0x6c,
	0x61, 0x79, 0x48, 0x69, 0x6e, 0x74, 0x52, 0x11, 0x6d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x44, 0x69,
	0x73, 0x70, 0x6c, 0x61, 0x79, 0x48, 0x69, 0x6e, 0x74, 0x22, 0x29, 0x0a, 0x0d, 0x53, 0x68, 0x75,
	0x74, 0x64, 0x6f, 0x77, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0xf0, 0x02, 0x0a, 0x15, 0x4d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x44,
	0x65, 0x74, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x18,
	0x0a, 0x07, 0x62, 0x61, 0x73, 0x69, 0x73, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x07, 0x62, 0x61, 0x73, 0x69, 0x73, 0x69, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x33, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x6d, 0x61, 0x73, 0x68, 0x75, 0x70,
	0x73, 0x64, 0x6b, 0x2e, 0x4d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e,
	0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x65, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x65, 0x72, 0x12, 0x24, 0x0a, 0x0d, 0x63, 0x6f, 0x6c, 0x61, 0x62, 0x72, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6f,
	0x6c, 0x61, 0x62, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x67,
	0x65, 0x6e, 0x72, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x67, 0x65, 0x6e, 0x72,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x75, 0x62, 0x67, 0x65, 0x6e, 0x72, 0x65, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x75, 0x62, 0x67, 0x65, 0x6e, 0x72, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x69, 0x64, 0x73, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x03,
	0x52, 0x09, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x69, 0x64, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x63,
	0x68, 0x69, 0x6c, 0x64, 0x69, 0x64, 0x73, 0x18, 0x0c, 0x20, 0x03, 0x28, 0x03, 0x52, 0x08, 0x63,
	0x68, 0x69, 0x6c, 0x64, 0x69, 0x64, 0x73, 0x22, 0x89, 0x01, 0x0a, 0x1b, 0x4d, 0x61, 0x73, 0x68,
	0x75, 0x70, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e,
	0x74, 0x42, 0x75, 0x6e, 0x64, 0x6c, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x75, 0x74, 0x68, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x75, 0x74, 0x68,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x4c, 0x0a, 0x10, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x65,
	0x64, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x20, 0x2e, 0x6d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x73, 0x64, 0x6b, 0x2e, 0x4d, 0x61, 0x73, 0x68,
	0x75, 0x70, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e,
	0x74, 0x52, 0x10, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x45, 0x6c, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x22, 0x3a, 0x0a, 0x12, 0x4d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x45, 0x6c, 0x65,
	0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61,
	0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22,
	0x7d, 0x0a, 0x18, 0x4d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x42, 0x75, 0x6e, 0x64, 0x6c, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x61,
	0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x61, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x43, 0x0a, 0x0d, 0x65, 0x6c, 0x65,
	0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x1d, 0x2e, 0x6d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x73, 0x64, 0x6b, 0x2e, 0x4d, 0x61, 0x73,
	0x68, 0x75, 0x70, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52,
	0x0d, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x73, 0x32, 0xe8,
	0x04, 0x0a, 0x0c, 0x4d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12,
	0x51, 0x0a, 0x05, 0x53, 0x68, 0x61, 0x6b, 0x65, 0x12, 0x22, 0x2e, 0x6d, 0x61, 0x73, 0x68, 0x75,
	0x70, 0x73, 0x64, 0x6b, 0x2e, 0x4d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x1a, 0x22, 0x2e, 0x6d,
	0x61, 0x73, 0x68, 0x75, 0x70, 0x73, 0x64, 0x6b, 0x2e, 0x4d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x43,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73,
	0x22, 0x00, 0x12, 0x4a, 0x0a, 0x08, 0x4f, 0x6e, 0x52, 0x65, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x1e,
	0x2e, 0x6d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x73, 0x64, 0x6b, 0x2e, 0x4d, 0x61, 0x73, 0x68, 0x75,
	0x70, 0x44, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x42, 0x75, 0x6e, 0x64, 0x6c, 0x65, 0x1a, 0x1c,
	0x2e, 0x6d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x73, 0x64, 0x6b, 0x2e, 0x4d, 0x61, 0x73, 0x68, 0x75,
	0x70, 0x44, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x48, 0x69, 0x6e, 0x74, 0x22, 0x00, 0x12, 0x55,
	0x0a, 0x11, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x45, 0x6c, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x12, 0x16, 0x2e, 0x6d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x73, 0x64, 0x6b, 0x2e,
	0x4d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x26, 0x2e, 0x6d, 0x61,
	0x73, 0x68, 0x75, 0x70, 0x73, 0x64, 0x6b, 0x2e, 0x4d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x44, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x42, 0x75, 0x6e,
	0x64, 0x6c, 0x65, 0x22, 0x00, 0x12, 0x68, 0x0a, 0x14, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x4d,
	0x61, 0x73, 0x68, 0x75, 0x70, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x26, 0x2e,
	0x6d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x73, 0x64, 0x6b, 0x2e, 0x4d, 0x61, 0x73, 0x68, 0x75, 0x70,
	0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x42,
	0x75, 0x6e, 0x64, 0x6c, 0x65, 0x1a, 0x26, 0x2e, 0x6d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x73, 0x64,
	0x6b, 0x2e, 0x4d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x65, 0x64,
	0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x42, 0x75, 0x6e, 0x64, 0x6c, 0x65, 0x22, 0x00, 0x12,
	0x67, 0x0a, 0x19, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x4d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x45,
	0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x23, 0x2e, 0x6d,
	0x61, 0x73, 0x68, 0x75, 0x70, 0x73, 0x64, 0x6b, 0x2e, 0x4d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x45,
	0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x42, 0x75, 0x6e, 0x64, 0x6c,
	0x65, 0x1a, 0x23, 0x2e, 0x6d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x73, 0x64, 0x6b, 0x2e, 0x4d, 0x61,
	0x73, 0x68, 0x75, 0x70, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x42, 0x75, 0x6e, 0x64, 0x6c, 0x65, 0x22, 0x00, 0x12, 0x51, 0x0a, 0x1d, 0x52, 0x65, 0x73, 0x65,
	0x74, 0x47, 0x33, 0x6e, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x45, 0x6c, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x73, 0x12, 0x16, 0x2e, 0x6d, 0x61, 0x73, 0x68,
	0x75, 0x70, 0x73, 0x64, 0x6b, 0x2e, 0x4d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x1a, 0x16, 0x2e, 0x6d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x73, 0x64, 0x6b, 0x2e, 0x4d, 0x61,
	0x73, 0x68, 0x75, 0x70, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x3c, 0x0a, 0x08, 0x53,
	0x68, 0x75, 0x74, 0x64, 0x6f, 0x77, 0x6e, 0x12, 0x16, 0x2e, 0x6d, 0x61, 0x73, 0x68, 0x75, 0x70,
	0x73, 0x64, 0x6b, 0x2e, 0x4d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a,
	0x16, 0x2e, 0x6d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x73, 0x64, 0x6b, 0x2e, 0x4d, 0x61, 0x73, 0x68,
	0x75, 0x70, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x41, 0x0a, 0x0e, 0x6e, 0x75, 0x74,
	0x65, 0x2e, 0x6d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x73, 0x64, 0x6b, 0x42, 0x09, 0x4d, 0x61, 0x73,
	0x68, 0x75, 0x70, 0x53, 0x64, 0x6b, 0x50, 0x01, 0x5a, 0x22, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x72, 0x6a, 0x72, 0x69, 0x65, 0x6b, 0x65, 0x2f, 0x6e, 0x75,
	0x74, 0x65, 0x2f, 0x6d, 0x61, 0x73, 0x68, 0x75, 0x70, 0x73, 0x64, 0x6b, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_mashupsdk_mashupsdk_proto_rawDescOnce sync.Once
	file_mashupsdk_mashupsdk_proto_rawDescData = file_mashupsdk_mashupsdk_proto_rawDesc
)

func file_mashupsdk_mashupsdk_proto_rawDescGZIP() []byte {
	file_mashupsdk_mashupsdk_proto_rawDescOnce.Do(func() {
		file_mashupsdk_mashupsdk_proto_rawDescData = protoimpl.X.CompressGZIP(file_mashupsdk_mashupsdk_proto_rawDescData)
	})
	return file_mashupsdk_mashupsdk_proto_rawDescData
}

var file_mashupsdk_mashupsdk_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_mashupsdk_mashupsdk_proto_goTypes = []interface{}{
	(*MashupEmpty)(nil),                 // 0: mashupsdk.MashupEmpty
	(*MashupConnectionConfigs)(nil),     // 1: mashupsdk.MashupConnectionConfigs
	(*MashupDisplayHint)(nil),           // 2: mashupsdk.MashupDisplayHint
	(*MashupDisplayBundle)(nil),         // 3: mashupsdk.MashupDisplayBundle
	(*ShutdownReply)(nil),               // 4: mashupsdk.ShutdownReply
	(*MashupDetailedElement)(nil),       // 5: mashupsdk.MashupDetailedElement
	(*MashupDetailedElementBundle)(nil), // 6: mashupsdk.MashupDetailedElementBundle
	(*MashupElementState)(nil),          // 7: mashupsdk.MashupElementState
	(*MashupElementStateBundle)(nil),    // 8: mashupsdk.MashupElementStateBundle
}
var file_mashupsdk_mashupsdk_proto_depIdxs = []int32{
	2,  // 0: mashupsdk.MashupDisplayBundle.mashupDisplayHint:type_name -> mashupsdk.MashupDisplayHint
	7,  // 1: mashupsdk.MashupDetailedElement.state:type_name -> mashupsdk.MashupElementState
	5,  // 2: mashupsdk.MashupDetailedElementBundle.detailedElements:type_name -> mashupsdk.MashupDetailedElement
	7,  // 3: mashupsdk.MashupElementStateBundle.elementStates:type_name -> mashupsdk.MashupElementState
	1,  // 4: mashupsdk.MashupServer.Shake:input_type -> mashupsdk.MashupConnectionConfigs
	3,  // 5: mashupsdk.MashupServer.OnResize:input_type -> mashupsdk.MashupDisplayBundle
	0,  // 6: mashupsdk.MashupServer.GetMashupElements:input_type -> mashupsdk.MashupEmpty
	6,  // 7: mashupsdk.MashupServer.UpsertMashupElements:input_type -> mashupsdk.MashupDetailedElementBundle
	8,  // 8: mashupsdk.MashupServer.UpsertMashupElementsState:input_type -> mashupsdk.MashupElementStateBundle
	0,  // 9: mashupsdk.MashupServer.ResetG3nDetailedElementStates:input_type -> mashupsdk.MashupEmpty
	0,  // 10: mashupsdk.MashupServer.Shutdown:input_type -> mashupsdk.MashupEmpty
	1,  // 11: mashupsdk.MashupServer.Shake:output_type -> mashupsdk.MashupConnectionConfigs
	2,  // 12: mashupsdk.MashupServer.OnResize:output_type -> mashupsdk.MashupDisplayHint
	6,  // 13: mashupsdk.MashupServer.GetMashupElements:output_type -> mashupsdk.MashupDetailedElementBundle
	6,  // 14: mashupsdk.MashupServer.UpsertMashupElements:output_type -> mashupsdk.MashupDetailedElementBundle
	8,  // 15: mashupsdk.MashupServer.UpsertMashupElementsState:output_type -> mashupsdk.MashupElementStateBundle
	0,  // 16: mashupsdk.MashupServer.ResetG3nDetailedElementStates:output_type -> mashupsdk.MashupEmpty
	0,  // 17: mashupsdk.MashupServer.Shutdown:output_type -> mashupsdk.MashupEmpty
	11, // [11:18] is the sub-list for method output_type
	4,  // [4:11] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_mashupsdk_mashupsdk_proto_init() }
func file_mashupsdk_mashupsdk_proto_init() {
	if File_mashupsdk_mashupsdk_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mashupsdk_mashupsdk_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MashupEmpty); i {
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
		file_mashupsdk_mashupsdk_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MashupConnectionConfigs); i {
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
		file_mashupsdk_mashupsdk_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MashupDisplayHint); i {
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
		file_mashupsdk_mashupsdk_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MashupDisplayBundle); i {
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
		file_mashupsdk_mashupsdk_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShutdownReply); i {
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
		file_mashupsdk_mashupsdk_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MashupDetailedElement); i {
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
		file_mashupsdk_mashupsdk_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MashupDetailedElementBundle); i {
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
		file_mashupsdk_mashupsdk_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MashupElementState); i {
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
		file_mashupsdk_mashupsdk_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MashupElementStateBundle); i {
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
			RawDescriptor: file_mashupsdk_mashupsdk_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_mashupsdk_mashupsdk_proto_goTypes,
		DependencyIndexes: file_mashupsdk_mashupsdk_proto_depIdxs,
		MessageInfos:      file_mashupsdk_mashupsdk_proto_msgTypes,
	}.Build()
	File_mashupsdk_mashupsdk_proto = out.File
	file_mashupsdk_mashupsdk_proto_rawDesc = nil
	file_mashupsdk_mashupsdk_proto_goTypes = nil
	file_mashupsdk_mashupsdk_proto_depIdxs = nil
}
