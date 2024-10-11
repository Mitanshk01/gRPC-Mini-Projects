// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: protofiles/player.proto

package protofiles

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

type MoveResult int32

const (
	MoveResult_SUCCESS     MoveResult = 0
	MoveResult_FAILURE     MoveResult = 1
	MoveResult_PLAYER_DEAD MoveResult = 2
	MoveResult_VICTORY     MoveResult = 3
)

// Enum value maps for MoveResult.
var (
	MoveResult_name = map[int32]string{
		0: "SUCCESS",
		1: "FAILURE",
		2: "PLAYER_DEAD",
		3: "VICTORY",
	}
	MoveResult_value = map[string]int32{
		"SUCCESS":     0,
		"FAILURE":     1,
		"PLAYER_DEAD": 2,
		"VICTORY":     3,
	}
)

func (x MoveResult) Enum() *MoveResult {
	p := new(MoveResult)
	*p = x
	return p
}

func (x MoveResult) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MoveResult) Descriptor() protoreflect.EnumDescriptor {
	return file_protofiles_player_proto_enumTypes[0].Descriptor()
}

func (MoveResult) Type() protoreflect.EnumType {
	return &file_protofiles_player_proto_enumTypes[0]
}

func (x MoveResult) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MoveResult.Descriptor instead.
func (MoveResult) EnumDescriptor() ([]byte, []int) {
	return file_protofiles_player_proto_rawDescGZIP(), []int{0}
}

type PlayerStatusResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Score        uint32    `protobuf:"varint,1,opt,name=score,proto3" json:"score,omitempty"`
	HealthPoints uint32    `protobuf:"varint,2,opt,name=health_points,json=healthPoints,proto3" json:"health_points,omitempty"`
	Position     *Position `protobuf:"bytes,3,opt,name=position,proto3" json:"position,omitempty"`
}

func (x *PlayerStatusResponse) Reset() {
	*x = PlayerStatusResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protofiles_player_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlayerStatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayerStatusResponse) ProtoMessage() {}

func (x *PlayerStatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protofiles_player_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayerStatusResponse.ProtoReflect.Descriptor instead.
func (*PlayerStatusResponse) Descriptor() ([]byte, []int) {
	return file_protofiles_player_proto_rawDescGZIP(), []int{0}
}

func (x *PlayerStatusResponse) GetScore() uint32 {
	if x != nil {
		return x.Score
	}
	return 0
}

func (x *PlayerStatusResponse) GetHealthPoints() uint32 {
	if x != nil {
		return x.HealthPoints
	}
	return 0
}

func (x *PlayerStatusResponse) GetPosition() *Position {
	if x != nil {
		return x.Position
	}
	return nil
}

type MoveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Direction string `protobuf:"bytes,1,opt,name=direction,proto3" json:"direction,omitempty"` // 'U', 'L', 'R', 'D'
}

func (x *MoveRequest) Reset() {
	*x = MoveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protofiles_player_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MoveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MoveRequest) ProtoMessage() {}

func (x *MoveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protofiles_player_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MoveRequest.ProtoReflect.Descriptor instead.
func (*MoveRequest) Descriptor() ([]byte, []int) {
	return file_protofiles_player_proto_rawDescGZIP(), []int{1}
}

func (x *MoveRequest) GetDirection() string {
	if x != nil {
		return x.Direction
	}
	return ""
}

type MoveResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result MoveResult `protobuf:"varint,1,opt,name=result,proto3,enum=player.MoveResult" json:"result,omitempty"`
}

func (x *MoveResponse) Reset() {
	*x = MoveResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protofiles_player_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MoveResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MoveResponse) ProtoMessage() {}

func (x *MoveResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protofiles_player_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MoveResponse.ProtoReflect.Descriptor instead.
func (*MoveResponse) Descriptor() ([]byte, []int) {
	return file_protofiles_player_proto_rawDescGZIP(), []int{2}
}

func (x *MoveResponse) GetResult() MoveResult {
	if x != nil {
		return x.Result
	}
	return MoveResult_SUCCESS
}

var File_protofiles_player_proto protoreflect.FileDescriptor

var file_protofiles_player_proto_rawDesc = []byte{
	0x0a, 0x17, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2f, 0x70, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x6c, 0x61, 0x79, 0x65,
	0x72, 0x1a, 0x17, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2f, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7f, 0x0a, 0x14, 0x50, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x68, 0x65, 0x61, 0x6c,
	0x74, 0x68, 0x5f, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x0c, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x12, 0x2c, 0x0a,
	0x08, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x08, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x2b, 0x0a, 0x0b, 0x4d,
	0x6f, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x64, 0x69,
	0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x64,
	0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x3a, 0x0a, 0x0c, 0x4d, 0x6f, 0x76, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x70, 0x6c, 0x61, 0x79, 0x65,
	0x72, 0x2e, 0x4d, 0x6f, 0x76, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x06, 0x72, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x2a, 0x44, 0x0a, 0x0a, 0x4d, 0x6f, 0x76, 0x65, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x55, 0x43, 0x43, 0x45, 0x53, 0x53, 0x10, 0x00, 0x12,
	0x0b, 0x0a, 0x07, 0x46, 0x41, 0x49, 0x4c, 0x55, 0x52, 0x45, 0x10, 0x01, 0x12, 0x0f, 0x0a, 0x0b,
	0x50, 0x4c, 0x41, 0x59, 0x45, 0x52, 0x5f, 0x44, 0x45, 0x41, 0x44, 0x10, 0x02, 0x12, 0x0b, 0x0a,
	0x07, 0x56, 0x49, 0x43, 0x54, 0x4f, 0x52, 0x59, 0x10, 0x03, 0x32, 0x91, 0x01, 0x0a, 0x0d, 0x50,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x45, 0x0a, 0x0f,
	0x47, 0x65, 0x74, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x14, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x1c, 0x2e, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2e, 0x50,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a, 0x0c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4d,
	0x6f, 0x76, 0x65, 0x12, 0x13, 0x2e, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2e, 0x4d, 0x6f, 0x76,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x70, 0x6c, 0x61, 0x79, 0x65,
	0x72, 0x2e, 0x4d, 0x6f, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x2c,
	0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4d, 0x69, 0x74,
	0x61, 0x6e, 0x73, 0x68, 0x6b, 0x30, 0x31, 0x2f, 0x44, 0x53, 0x5f, 0x48, 0x57, 0x34, 0x2f, 0x51,
	0x31, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protofiles_player_proto_rawDescOnce sync.Once
	file_protofiles_player_proto_rawDescData = file_protofiles_player_proto_rawDesc
)

func file_protofiles_player_proto_rawDescGZIP() []byte {
	file_protofiles_player_proto_rawDescOnce.Do(func() {
		file_protofiles_player_proto_rawDescData = protoimpl.X.CompressGZIP(file_protofiles_player_proto_rawDescData)
	})
	return file_protofiles_player_proto_rawDescData
}

var file_protofiles_player_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_protofiles_player_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_protofiles_player_proto_goTypes = []any{
	(MoveResult)(0),              // 0: player.MoveResult
	(*PlayerStatusResponse)(nil), // 1: player.PlayerStatusResponse
	(*MoveRequest)(nil),          // 2: player.MoveRequest
	(*MoveResponse)(nil),         // 3: player.MoveResponse
	(*Position)(nil),             // 4: common.Position
	(*EmptyMessage)(nil),         // 5: common.EmptyMessage
}
var file_protofiles_player_proto_depIdxs = []int32{
	4, // 0: player.PlayerStatusResponse.position:type_name -> common.Position
	0, // 1: player.MoveResponse.result:type_name -> player.MoveResult
	5, // 2: player.PlayerService.GetPlayerStatus:input_type -> common.EmptyMessage
	2, // 3: player.PlayerService.RegisterMove:input_type -> player.MoveRequest
	1, // 4: player.PlayerService.GetPlayerStatus:output_type -> player.PlayerStatusResponse
	3, // 5: player.PlayerService.RegisterMove:output_type -> player.MoveResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_protofiles_player_proto_init() }
func file_protofiles_player_proto_init() {
	if File_protofiles_player_proto != nil {
		return
	}
	file_protofiles_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_protofiles_player_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*PlayerStatusResponse); i {
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
		file_protofiles_player_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*MoveRequest); i {
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
		file_protofiles_player_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*MoveResponse); i {
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
			RawDescriptor: file_protofiles_player_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protofiles_player_proto_goTypes,
		DependencyIndexes: file_protofiles_player_proto_depIdxs,
		EnumInfos:         file_protofiles_player_proto_enumTypes,
		MessageInfos:      file_protofiles_player_proto_msgTypes,
	}.Build()
	File_protofiles_player_proto = out.File
	file_protofiles_player_proto_rawDesc = nil
	file_protofiles_player_proto_goTypes = nil
	file_protofiles_player_proto_depIdxs = nil
}