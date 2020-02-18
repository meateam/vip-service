package db

import (
	"context"

	pb "github.com/meateam/vip-service/proto"
)

// Controller is the vip service business logic implementation.
type Controller struct {
}

// NewController returns a new controller.
func NewController() (Controller, error) {

	return Controller{}, nil
}

// HealthCheck runs store's healthcheck and returns true if healthy, otherwise returns false
// and any error if occured.
func (c Controller) HealthCheck(ctx context.Context) (bool, error) {
	return true, nil
}

// DBGetIsVIPByID returns a vip.
func (c Controller) DBGetIsVIPByID(ctx context.Context, vipID string) (pb.VIPObject, error) {
	vips := []string{
		"Shahar",
		"Yonatan",
		"Kiddon",
	}
	isVIP := stringInSlice(vipID, vips)

	//TODO
	vipRes := pb.VIPObject{VipID: vipID, IsVIP: isVIP}

	return vipRes, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
