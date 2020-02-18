package service

import (
	"context"
	"fmt"
	"time"

	pb "github.com/meateam/vip-service/proto"
	ctrlr "github.com/meateam/vip-service/service/db"
	"github.com/sirupsen/logrus"
)

const ()

// Service is the structure used for handling
type Service struct {
	logger    *logrus.Logger
	grantType string
	audience  string
}

// HealthCheck checks the health of the service, and returns a boolean accordingly.
func (s *Service) HealthCheck(mongoClientPingTimeout time.Duration) bool {
	return true
}

// NewService creates a Service and returns it.
func NewService(logger *logrus.Logger) Service {
	s := Service{logger: logger}
	return s
}

// GetIsVIPByID is the request handler for getting a vip (user, status) by file id.
func (s Service) GetIsVIPByID(ctx context.Context, req *pb.GetIsVIPByIDRequest) (*pb.GetIsVIPByIDResponse, error) {
	controller := ctrlr.Controller{}
	vipID := req.GetVipID()
	if vipID == "" {
		return nil, fmt.Errorf("vipID is required")
	}

	vip, err := controller.DBGetIsVIPByID(ctx, vipID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the vip %v", err)
	}

	//TODO
	return &pb.GetIsVIPByIDResponse{VipID: vip.VipID, IsVIP: vip.IsVIP}, nil
}
