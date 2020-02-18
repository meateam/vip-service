package service

import (
	pb "github.com/meateam/vip-service/proto"
)

// VIP is an interface of a vip object.
type VIP interface {
	GetID() string
	SetID(id string) error

	GetVIPID() string
	SetVIPID(vipID string) error

	MarshalProto(vip *pb.VIPObject) error
}
