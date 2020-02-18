package service

import (
	"context"

	pb "github.com/meateam/vip-service/proto"
)

// Controller is an interface for the business logic of the vip-service which uses a Store.
type Controller interface {
	DBGetIsVIPByID(ctx context.Context, vipID string) (pb.VIPObject, error)
	HealthCheck(ctx context.Context) (bool, error)
}
