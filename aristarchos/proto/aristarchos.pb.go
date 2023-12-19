// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.1
// source: proto/aristarchos.proto

package proto

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

type AggregatorRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RootWord string `protobuf:"bytes,1,opt,name=root_word,json=rootWord,proto3" json:"root_word,omitempty"`
}

func (x *AggregatorRequest) Reset() {
	*x = AggregatorRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_aristarchos_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AggregatorRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AggregatorRequest) ProtoMessage() {}

func (x *AggregatorRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_aristarchos_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AggregatorRequest.ProtoReflect.Descriptor instead.
func (*AggregatorRequest) Descriptor() ([]byte, []int) {
	return file_proto_aristarchos_proto_rawDescGZIP(), []int{0}
}

func (x *AggregatorRequest) GetRootWord() string {
	if x != nil {
		return x.RootWord
	}
	return ""
}

type AggregatorCreationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Word        string `protobuf:"bytes,1,opt,name=word,proto3" json:"word,omitempty"`
	Rule        string `protobuf:"bytes,2,opt,name=rule,proto3" json:"rule,omitempty"`
	RootWord    string `protobuf:"bytes,3,opt,name=root_word,json=rootWord,proto3" json:"root_word,omitempty"`
	Translation string `protobuf:"bytes,4,opt,name=translation,proto3" json:"translation,omitempty"`
}

func (x *AggregatorCreationRequest) Reset() {
	*x = AggregatorCreationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_aristarchos_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AggregatorCreationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AggregatorCreationRequest) ProtoMessage() {}

func (x *AggregatorCreationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_aristarchos_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AggregatorCreationRequest.ProtoReflect.Descriptor instead.
func (*AggregatorCreationRequest) Descriptor() ([]byte, []int) {
	return file_proto_aristarchos_proto_rawDescGZIP(), []int{1}
}

func (x *AggregatorCreationRequest) GetWord() string {
	if x != nil {
		return x.Word
	}
	return ""
}

func (x *AggregatorCreationRequest) GetRule() string {
	if x != nil {
		return x.Rule
	}
	return ""
}

func (x *AggregatorCreationRequest) GetRootWord() string {
	if x != nil {
		return x.RootWord
	}
	return ""
}

func (x *AggregatorCreationRequest) GetTranslation() string {
	if x != nil {
		return x.Translation
	}
	return ""
}

type AggregatorCreationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Created bool `protobuf:"varint,1,opt,name=created,proto3" json:"created,omitempty"`
	Updated bool `protobuf:"varint,2,opt,name=updated,proto3" json:"updated,omitempty"`
}

func (x *AggregatorCreationResponse) Reset() {
	*x = AggregatorCreationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_aristarchos_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AggregatorCreationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AggregatorCreationResponse) ProtoMessage() {}

func (x *AggregatorCreationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_aristarchos_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AggregatorCreationResponse.ProtoReflect.Descriptor instead.
func (*AggregatorCreationResponse) Descriptor() ([]byte, []int) {
	return file_proto_aristarchos_proto_rawDescGZIP(), []int{2}
}

func (x *AggregatorCreationResponse) GetCreated() bool {
	if x != nil {
		return x.Created
	}
	return false
}

func (x *AggregatorCreationResponse) GetUpdated() bool {
	if x != nil {
		return x.Updated
	}
	return false
}

type SearchWordResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Word []string `protobuf:"bytes,1,rep,name=word,proto3" json:"word,omitempty"`
}

func (x *SearchWordResponse) Reset() {
	*x = SearchWordResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_aristarchos_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchWordResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchWordResponse) ProtoMessage() {}

func (x *SearchWordResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_aristarchos_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchWordResponse.ProtoReflect.Descriptor instead.
func (*SearchWordResponse) Descriptor() ([]byte, []int) {
	return file_proto_aristarchos_proto_rawDescGZIP(), []int{3}
}

func (x *SearchWordResponse) GetWord() []string {
	if x != nil {
		return x.Word
	}
	return nil
}

type ConjugationForm struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Number string `protobuf:"bytes,1,opt,name=number,proto3" json:"number,omitempty"`
	Person string `protobuf:"bytes,2,opt,name=person,proto3" json:"person,omitempty"`
	Word   string `protobuf:"bytes,3,opt,name=word,proto3" json:"word,omitempty"`
}

func (x *ConjugationForm) Reset() {
	*x = ConjugationForm{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_aristarchos_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConjugationForm) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConjugationForm) ProtoMessage() {}

func (x *ConjugationForm) ProtoReflect() protoreflect.Message {
	mi := &file_proto_aristarchos_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConjugationForm.ProtoReflect.Descriptor instead.
func (*ConjugationForm) Descriptor() ([]byte, []int) {
	return file_proto_aristarchos_proto_rawDescGZIP(), []int{4}
}

func (x *ConjugationForm) GetNumber() string {
	if x != nil {
		return x.Number
	}
	return ""
}

func (x *ConjugationForm) GetPerson() string {
	if x != nil {
		return x.Person
	}
	return ""
}

func (x *ConjugationForm) GetWord() string {
	if x != nil {
		return x.Word
	}
	return ""
}

type Conjugation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tense  string             `protobuf:"bytes,1,opt,name=tense,proto3" json:"tense,omitempty"`
	Mood   string             `protobuf:"bytes,2,opt,name=mood,proto3" json:"mood,omitempty"`
	Aspect string             `protobuf:"bytes,3,opt,name=aspect,proto3" json:"aspect,omitempty"`
	Forms  []*ConjugationForm `protobuf:"bytes,4,rep,name=forms,proto3" json:"forms,omitempty"`
}

func (x *Conjugation) Reset() {
	*x = Conjugation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_aristarchos_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Conjugation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Conjugation) ProtoMessage() {}

func (x *Conjugation) ProtoReflect() protoreflect.Message {
	mi := &file_proto_aristarchos_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Conjugation.ProtoReflect.Descriptor instead.
func (*Conjugation) Descriptor() ([]byte, []int) {
	return file_proto_aristarchos_proto_rawDescGZIP(), []int{5}
}

func (x *Conjugation) GetTense() string {
	if x != nil {
		return x.Tense
	}
	return ""
}

func (x *Conjugation) GetMood() string {
	if x != nil {
		return x.Mood
	}
	return ""
}

func (x *Conjugation) GetAspect() string {
	if x != nil {
		return x.Aspect
	}
	return ""
}

func (x *Conjugation) GetForms() []*ConjugationForm {
	if x != nil {
		return x.Forms
	}
	return nil
}

type RootWordResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RootWord     string         `protobuf:"bytes,1,opt,name=rootWord,proto3" json:"rootWord,omitempty"`
	Translations []string       `protobuf:"bytes,2,rep,name=translations,proto3" json:"translations,omitempty"`
	Conjugations []*Conjugation `protobuf:"bytes,3,rep,name=conjugations,proto3" json:"conjugations,omitempty"`
}

func (x *RootWordResponse) Reset() {
	*x = RootWordResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_aristarchos_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RootWordResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RootWordResponse) ProtoMessage() {}

func (x *RootWordResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_aristarchos_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RootWordResponse.ProtoReflect.Descriptor instead.
func (*RootWordResponse) Descriptor() ([]byte, []int) {
	return file_proto_aristarchos_proto_rawDescGZIP(), []int{6}
}

func (x *RootWordResponse) GetRootWord() string {
	if x != nil {
		return x.RootWord
	}
	return ""
}

func (x *RootWordResponse) GetTranslations() []string {
	if x != nil {
		return x.Translations
	}
	return nil
}

func (x *RootWordResponse) GetConjugations() []*Conjugation {
	if x != nil {
		return x.Conjugations
	}
	return nil
}

type HealthResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Health bool `protobuf:"varint,1,opt,name=health,proto3" json:"health,omitempty"`
}

func (x *HealthResponse) Reset() {
	*x = HealthResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_aristarchos_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthResponse) ProtoMessage() {}

func (x *HealthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_aristarchos_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthResponse.ProtoReflect.Descriptor instead.
func (*HealthResponse) Descriptor() ([]byte, []int) {
	return file_proto_aristarchos_proto_rawDescGZIP(), []int{7}
}

func (x *HealthResponse) GetHealth() bool {
	if x != nil {
		return x.Health
	}
	return false
}

type HealthRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *HealthRequest) Reset() {
	*x = HealthRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_aristarchos_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthRequest) ProtoMessage() {}

func (x *HealthRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_aristarchos_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthRequest.ProtoReflect.Descriptor instead.
func (*HealthRequest) Descriptor() ([]byte, []int) {
	return file_proto_aristarchos_proto_rawDescGZIP(), []int{8}
}

var File_proto_aristarchos_proto protoreflect.FileDescriptor

var file_proto_aristarchos_proto_rawDesc = []byte{
	0x0a, 0x17, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x72, 0x63,
	0x68, 0x6f, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x6f, 0x6c, 0x79, 0x6d, 0x70,
	0x69, 0x61, 0x5f, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x72, 0x63, 0x68, 0x6f, 0x73, 0x22, 0x30,
	0x0a, 0x11, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x6f, 0x6f, 0x74, 0x5f, 0x77, 0x6f, 0x72, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x6f, 0x6f, 0x74, 0x57, 0x6f, 0x72, 0x64,
	0x22, 0x82, 0x01, 0x0a, 0x19, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x77, 0x6f,
	0x72, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x75, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x72, 0x75, 0x6c, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x6f, 0x6f, 0x74, 0x5f, 0x77,
	0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x6f, 0x6f, 0x74, 0x57,
	0x6f, 0x72, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6c,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x50, 0x0a, 0x1a, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61,
	0x74, 0x6f, 0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x18, 0x0a,
	0x07, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x22, 0x28, 0x0a, 0x12, 0x53, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x57, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x77, 0x6f, 0x72,
	0x64, 0x22, 0x55, 0x0a, 0x0f, 0x43, 0x6f, 0x6e, 0x6a, 0x75, 0x67, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x46, 0x6f, 0x72, 0x6d, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06,
	0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x65,
	0x72, 0x73, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x8b, 0x01, 0x0a, 0x0b, 0x43, 0x6f, 0x6e,
	0x6a, 0x75, 0x67, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x65, 0x6e, 0x73,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x65, 0x6e, 0x73, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x6d, 0x6f, 0x6f, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6d, 0x6f,
	0x6f, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x73, 0x70, 0x65, 0x63, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x61, 0x73, 0x70, 0x65, 0x63, 0x74, 0x12, 0x3a, 0x0a, 0x05, 0x66, 0x6f,
	0x72, 0x6d, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x6f, 0x6c, 0x79, 0x6d,
	0x70, 0x69, 0x61, 0x5f, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x72, 0x63, 0x68, 0x6f, 0x73, 0x2e,
	0x43, 0x6f, 0x6e, 0x6a, 0x75, 0x67, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x6f, 0x72, 0x6d, 0x52,
	0x05, 0x66, 0x6f, 0x72, 0x6d, 0x73, 0x22, 0x98, 0x01, 0x0a, 0x10, 0x52, 0x6f, 0x6f, 0x74, 0x57,
	0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72,
	0x6f, 0x6f, 0x74, 0x57, 0x6f, 0x72, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72,
	0x6f, 0x6f, 0x74, 0x57, 0x6f, 0x72, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x44, 0x0a, 0x0c, 0x63,
	0x6f, 0x6e, 0x6a, 0x75, 0x67, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x20, 0x2e, 0x6f, 0x6c, 0x79, 0x6d, 0x70, 0x69, 0x61, 0x5f, 0x61, 0x72, 0x69, 0x73,
	0x74, 0x61, 0x72, 0x63, 0x68, 0x6f, 0x73, 0x2e, 0x43, 0x6f, 0x6e, 0x6a, 0x75, 0x67, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x63, 0x6f, 0x6e, 0x6a, 0x75, 0x67, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x22, 0x28, 0x0a, 0x0e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x06, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x22, 0x0f, 0x0a, 0x0d, 0x48,
	0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x32, 0xa3, 0x03, 0x0a,
	0x0b, 0x41, 0x72, 0x69, 0x73, 0x74, 0x61, 0x72, 0x63, 0x68, 0x6f, 0x73, 0x12, 0x73, 0x0a, 0x0e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4e, 0x65, 0x77, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x2e,
	0x2e, 0x6f, 0x6c, 0x79, 0x6d, 0x70, 0x69, 0x61, 0x5f, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x72,
	0x63, 0x68, 0x6f, 0x73, 0x2e, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2f,
	0x2e, 0x6f, 0x6c, 0x79, 0x6d, 0x70, 0x69, 0x61, 0x5f, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x72,
	0x63, 0x68, 0x6f, 0x73, 0x2e, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x60, 0x0a, 0x0d, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x26, 0x2e, 0x6f, 0x6c, 0x79, 0x6d, 0x70, 0x69, 0x61, 0x5f, 0x61, 0x72, 0x69,
	0x73, 0x74, 0x61, 0x72, 0x63, 0x68, 0x6f, 0x73, 0x2e, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61,
	0x74, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x6f, 0x6c, 0x79,
	0x6d, 0x70, 0x69, 0x61, 0x5f, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x72, 0x63, 0x68, 0x6f, 0x73,
	0x2e, 0x52, 0x6f, 0x6f, 0x74, 0x57, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x68, 0x0a, 0x13, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x53,
	0x65, 0x61, 0x72, 0x63, 0x68, 0x57, 0x6f, 0x72, 0x64, 0x73, 0x12, 0x26, 0x2e, 0x6f, 0x6c, 0x79,
	0x6d, 0x70, 0x69, 0x61, 0x5f, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x72, 0x63, 0x68, 0x6f, 0x73,
	0x2e, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x27, 0x2e, 0x6f, 0x6c, 0x79, 0x6d, 0x70, 0x69, 0x61, 0x5f, 0x61, 0x72, 0x69,
	0x73, 0x74, 0x61, 0x72, 0x63, 0x68, 0x6f, 0x73, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x57,
	0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x53, 0x0a,
	0x06, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x12, 0x22, 0x2e, 0x6f, 0x6c, 0x79, 0x6d, 0x70, 0x69,
	0x61, 0x5f, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x72, 0x63, 0x68, 0x6f, 0x73, 0x2e, 0x48, 0x65,
	0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x6f, 0x6c,
	0x79, 0x6d, 0x70, 0x69, 0x61, 0x5f, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x72, 0x63, 0x68, 0x6f,
	0x73, 0x2e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x6f, 0x64, 0x79, 0x73, 0x73, 0x65, 0x69, 0x61, 0x2d, 0x67, 0x72, 0x65, 0x65, 0x6b, 0x2f,
	0x6f, 0x6c, 0x79, 0x6d, 0x70, 0x69, 0x61, 0x2f, 0x61, 0x72, 0x69, 0x73, 0x74, 0x61, 0x72, 0x63,
	0x68, 0x6f, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_proto_aristarchos_proto_rawDescOnce sync.Once
	file_proto_aristarchos_proto_rawDescData = file_proto_aristarchos_proto_rawDesc
)

func file_proto_aristarchos_proto_rawDescGZIP() []byte {
	file_proto_aristarchos_proto_rawDescOnce.Do(func() {
		file_proto_aristarchos_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_aristarchos_proto_rawDescData)
	})
	return file_proto_aristarchos_proto_rawDescData
}

var file_proto_aristarchos_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_proto_aristarchos_proto_goTypes = []interface{}{
	(*AggregatorRequest)(nil),          // 0: olympia_aristarchos.AggregatorRequest
	(*AggregatorCreationRequest)(nil),  // 1: olympia_aristarchos.AggregatorCreationRequest
	(*AggregatorCreationResponse)(nil), // 2: olympia_aristarchos.AggregatorCreationResponse
	(*SearchWordResponse)(nil),         // 3: olympia_aristarchos.SearchWordResponse
	(*ConjugationForm)(nil),            // 4: olympia_aristarchos.ConjugationForm
	(*Conjugation)(nil),                // 5: olympia_aristarchos.Conjugation
	(*RootWordResponse)(nil),           // 6: olympia_aristarchos.RootWordResponse
	(*HealthResponse)(nil),             // 7: olympia_aristarchos.HealthResponse
	(*HealthRequest)(nil),              // 8: olympia_aristarchos.HealthRequest
}
var file_proto_aristarchos_proto_depIdxs = []int32{
	4, // 0: olympia_aristarchos.Conjugation.forms:type_name -> olympia_aristarchos.ConjugationForm
	5, // 1: olympia_aristarchos.RootWordResponse.conjugations:type_name -> olympia_aristarchos.Conjugation
	1, // 2: olympia_aristarchos.Aristarchos.CreateNewEntry:input_type -> olympia_aristarchos.AggregatorCreationRequest
	0, // 3: olympia_aristarchos.Aristarchos.RetrieveEntry:input_type -> olympia_aristarchos.AggregatorRequest
	0, // 4: olympia_aristarchos.Aristarchos.RetrieveSearchWords:input_type -> olympia_aristarchos.AggregatorRequest
	8, // 5: olympia_aristarchos.Aristarchos.Health:input_type -> olympia_aristarchos.HealthRequest
	2, // 6: olympia_aristarchos.Aristarchos.CreateNewEntry:output_type -> olympia_aristarchos.AggregatorCreationResponse
	6, // 7: olympia_aristarchos.Aristarchos.RetrieveEntry:output_type -> olympia_aristarchos.RootWordResponse
	3, // 8: olympia_aristarchos.Aristarchos.RetrieveSearchWords:output_type -> olympia_aristarchos.SearchWordResponse
	7, // 9: olympia_aristarchos.Aristarchos.Health:output_type -> olympia_aristarchos.HealthResponse
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_aristarchos_proto_init() }
func file_proto_aristarchos_proto_init() {
	if File_proto_aristarchos_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_aristarchos_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AggregatorRequest); i {
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
		file_proto_aristarchos_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AggregatorCreationRequest); i {
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
		file_proto_aristarchos_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AggregatorCreationResponse); i {
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
		file_proto_aristarchos_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchWordResponse); i {
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
		file_proto_aristarchos_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConjugationForm); i {
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
		file_proto_aristarchos_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Conjugation); i {
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
		file_proto_aristarchos_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RootWordResponse); i {
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
		file_proto_aristarchos_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HealthResponse); i {
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
		file_proto_aristarchos_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HealthRequest); i {
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
			RawDescriptor: file_proto_aristarchos_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_aristarchos_proto_goTypes,
		DependencyIndexes: file_proto_aristarchos_proto_depIdxs,
		MessageInfos:      file_proto_aristarchos_proto_msgTypes,
	}.Build()
	File_proto_aristarchos_proto = out.File
	file_proto_aristarchos_proto_rawDesc = nil
	file_proto_aristarchos_proto_goTypes = nil
	file_proto_aristarchos_proto_depIdxs = nil
}
