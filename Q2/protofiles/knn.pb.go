// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: protofiles/knn.proto

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

type DataPoint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Coordinates []float32 `protobuf:"fixed32,1,rep,packed,name=coordinates,proto3" json:"coordinates,omitempty"`
}

func (x *DataPoint) Reset() {
	*x = DataPoint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protofiles_knn_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataPoint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataPoint) ProtoMessage() {}

func (x *DataPoint) ProtoReflect() protoreflect.Message {
	mi := &file_protofiles_knn_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataPoint.ProtoReflect.Descriptor instead.
func (*DataPoint) Descriptor() ([]byte, []int) {
	return file_protofiles_knn_proto_rawDescGZIP(), []int{0}
}

func (x *DataPoint) GetCoordinates() []float32 {
	if x != nil {
		return x.Coordinates
	}
	return nil
}

type Neighbor struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Point    *DataPoint `protobuf:"bytes,1,opt,name=point,proto3" json:"point,omitempty"`
	Distance float32    `protobuf:"fixed32,2,opt,name=distance,proto3" json:"distance,omitempty"`
}

func (x *Neighbor) Reset() {
	*x = Neighbor{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protofiles_knn_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Neighbor) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Neighbor) ProtoMessage() {}

func (x *Neighbor) ProtoReflect() protoreflect.Message {
	mi := &file_protofiles_knn_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Neighbor.ProtoReflect.Descriptor instead.
func (*Neighbor) Descriptor() ([]byte, []int) {
	return file_protofiles_knn_proto_rawDescGZIP(), []int{1}
}

func (x *Neighbor) GetPoint() *DataPoint {
	if x != nil {
		return x.Point
	}
	return nil
}

func (x *Neighbor) GetDistance() float32 {
	if x != nil {
		return x.Distance
	}
	return 0
}

type KNNRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueryPoint *DataPoint `protobuf:"bytes,1,opt,name=query_point,json=queryPoint,proto3" json:"query_point,omitempty"`
	K          int32      `protobuf:"varint,2,opt,name=k,proto3" json:"k,omitempty"`
}

func (x *KNNRequest) Reset() {
	*x = KNNRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protofiles_knn_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KNNRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KNNRequest) ProtoMessage() {}

func (x *KNNRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protofiles_knn_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KNNRequest.ProtoReflect.Descriptor instead.
func (*KNNRequest) Descriptor() ([]byte, []int) {
	return file_protofiles_knn_proto_rawDescGZIP(), []int{2}
}

func (x *KNNRequest) GetQueryPoint() *DataPoint {
	if x != nil {
		return x.QueryPoint
	}
	return nil
}

func (x *KNNRequest) GetK() int32 {
	if x != nil {
		return x.K
	}
	return 0
}

type KNNResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Neighbors []*Neighbor `protobuf:"bytes,1,rep,name=neighbors,proto3" json:"neighbors,omitempty"`
}

func (x *KNNResponse) Reset() {
	*x = KNNResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protofiles_knn_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KNNResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KNNResponse) ProtoMessage() {}

func (x *KNNResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protofiles_knn_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KNNResponse.ProtoReflect.Descriptor instead.
func (*KNNResponse) Descriptor() ([]byte, []int) {
	return file_protofiles_knn_proto_rawDescGZIP(), []int{3}
}

func (x *KNNResponse) GetNeighbors() []*Neighbor {
	if x != nil {
		return x.Neighbors
	}
	return nil
}

var File_protofiles_knn_proto protoreflect.FileDescriptor

var file_protofiles_knn_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2f, 0x6b, 0x6e, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x6b, 0x6e, 0x6e, 0x22, 0x2d, 0x0a, 0x09, 0x44,
	0x61, 0x74, 0x61, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x6f, 0x6f, 0x72,
	0x64, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x02, 0x52, 0x0b, 0x63,
	0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x73, 0x22, 0x4c, 0x0a, 0x08, 0x4e, 0x65,
	0x69, 0x67, 0x68, 0x62, 0x6f, 0x72, 0x12, 0x24, 0x0a, 0x05, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x6b, 0x6e, 0x6e, 0x2e, 0x44, 0x61, 0x74, 0x61,
	0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x05, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x08,
	0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x08,
	0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x22, 0x4b, 0x0a, 0x0a, 0x4b, 0x4e, 0x4e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2f, 0x0a, 0x0b, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x6b, 0x6e,
	0x6e, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x0a, 0x71, 0x75, 0x65,
	0x72, 0x79, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x0c, 0x0a, 0x01, 0x6b, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x01, 0x6b, 0x22, 0x3a, 0x0a, 0x0b, 0x4b, 0x4e, 0x4e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x09, 0x6e, 0x65, 0x69, 0x67, 0x68, 0x62, 0x6f, 0x72,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x6b, 0x6e, 0x6e, 0x2e, 0x4e, 0x65,
	0x69, 0x67, 0x68, 0x62, 0x6f, 0x72, 0x52, 0x09, 0x6e, 0x65, 0x69, 0x67, 0x68, 0x62, 0x6f, 0x72,
	0x73, 0x32, 0x47, 0x0a, 0x0a, 0x4b, 0x4e, 0x4e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x39, 0x0a, 0x15, 0x46, 0x69, 0x6e, 0x64, 0x4b, 0x4e, 0x65, 0x61, 0x72, 0x65, 0x73, 0x74, 0x4e,
	0x65, 0x69, 0x67, 0x68, 0x62, 0x6f, 0x72, 0x73, 0x12, 0x0f, 0x2e, 0x6b, 0x6e, 0x6e, 0x2e, 0x4b,
	0x4e, 0x4e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x6b, 0x6e, 0x6e, 0x2e,
	0x4e, 0x65, 0x69, 0x67, 0x68, 0x62, 0x6f, 0x72, 0x30, 0x01, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4d, 0x69, 0x74, 0x61, 0x6e, 0x73, 0x68,
	0x6b, 0x30, 0x31, 0x2f, 0x44, 0x53, 0x5f, 0x48, 0x57, 0x34, 0x2f, 0x51, 0x32, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protofiles_knn_proto_rawDescOnce sync.Once
	file_protofiles_knn_proto_rawDescData = file_protofiles_knn_proto_rawDesc
)

func file_protofiles_knn_proto_rawDescGZIP() []byte {
	file_protofiles_knn_proto_rawDescOnce.Do(func() {
		file_protofiles_knn_proto_rawDescData = protoimpl.X.CompressGZIP(file_protofiles_knn_proto_rawDescData)
	})
	return file_protofiles_knn_proto_rawDescData
}

var file_protofiles_knn_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_protofiles_knn_proto_goTypes = []any{
	(*DataPoint)(nil),   // 0: knn.DataPoint
	(*Neighbor)(nil),    // 1: knn.Neighbor
	(*KNNRequest)(nil),  // 2: knn.KNNRequest
	(*KNNResponse)(nil), // 3: knn.KNNResponse
}
var file_protofiles_knn_proto_depIdxs = []int32{
	0, // 0: knn.Neighbor.point:type_name -> knn.DataPoint
	0, // 1: knn.KNNRequest.query_point:type_name -> knn.DataPoint
	1, // 2: knn.KNNResponse.neighbors:type_name -> knn.Neighbor
	2, // 3: knn.KNNService.FindKNearestNeighbors:input_type -> knn.KNNRequest
	1, // 4: knn.KNNService.FindKNearestNeighbors:output_type -> knn.Neighbor
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_protofiles_knn_proto_init() }
func file_protofiles_knn_proto_init() {
	if File_protofiles_knn_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protofiles_knn_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*DataPoint); i {
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
		file_protofiles_knn_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Neighbor); i {
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
		file_protofiles_knn_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*KNNRequest); i {
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
		file_protofiles_knn_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*KNNResponse); i {
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
			RawDescriptor: file_protofiles_knn_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protofiles_knn_proto_goTypes,
		DependencyIndexes: file_protofiles_knn_proto_depIdxs,
		MessageInfos:      file_protofiles_knn_proto_msgTypes,
	}.Build()
	File_protofiles_knn_proto = out.File
	file_protofiles_knn_proto_rawDesc = nil
	file_protofiles_knn_proto_goTypes = nil
	file_protofiles_knn_proto_depIdxs = nil
}
