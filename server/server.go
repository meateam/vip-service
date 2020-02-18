package server

import (
	"net"
	"strings"
	"time"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	ilogger "github.com/meateam/elasticsearch-logger"
	pb "github.com/meateam/vip-service/proto"
	"github.com/meateam/vip-service/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const (
	envPrefix                  = "ms"
	configPort                 = "port"
	configHealthCheckInterval  = "health_check_interval"
	configElasticAPMIgnoreURLS = "elastic_apm_ignore_urls"
)

func init() {
	viper.SetDefault(configPort, "8080")
	viper.SetDefault(configHealthCheckInterval, 3)
	viper.SetDefault(configElasticAPMIgnoreURLS, "/grpc.health.v1.Health/Check")
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
}

// VIPServer is a structure that holds the vip grpc server
// and its services configuration
type VIPServer struct {
	*grpc.Server
	logger              *logrus.Logger
	port                string
	healthCheckInterval int
	vipService          service.Service
}

// Serve accepts incoming connections on the listener `lis`, creating a new
// ServerTransport and service goroutine for each. The service goroutines
// read gRPC requests and then call the registered handlers to reply to them.
// Serve returns when `lis.Accept` fails with fatal errors. `lis` will be closed when
// this method returns.
// If `lis` is nil then Serve creates a `net.Listener` with "tcp" network listening
// on the configured `TCP_PORT`, which defaults to "8080".
// Serve will return a non-nil error unless Stop or GracefulStop is called.
func (s VIPServer) Serve(lis net.Listener) {
	listener := lis
	if lis == nil {
		l, err := net.Listen("tcp", ":"+s.port)
		if err != nil {
			s.logger.Fatalf("failed to listen: %v", err)
		}

		listener = l
	}

	s.logger.Infof("listening and serving grpc server on port %s", s.port)
	if err := s.Server.Serve(listener); err != nil {
		s.logger.Fatalf(err.Error())
	}
}

// NewServer configures and creates a grpc.Server instance with the vip service
// health check service.
// Configure using environment variables.
// `HEALTH_CHECK_INTERVAL`: Interval to update serving state of the health check server.
// `TCP_PORT`: TCP port on which the grpc server would serve on.
func NewServer(logger *logrus.Logger) *VIPServer {
	// If no logger is given, create a new default logger for the server.
	if logger == nil {
		logger = ilogger.NewLogger()
	}

	// Set up grpc server opts with logger interceptor.
	serverOpts := append(
		serverLoggerInterceptor(logger),
		grpc.MaxRecvMsgSize(10<<20),
	)

	// Create a new grpc server.
	grpcServer := grpc.NewServer(
		serverOpts...,
	)

	// Create a vip service and register it on the grpc server.
	vipService := service.NewService(logger)
	pb.RegisterVIPServer(grpcServer, vipService)

	// Create a health server and register it on the grpc server.
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	vipServer := &VIPServer{
		Server:              grpcServer,
		logger:              logger,
		port:                viper.GetString(configPort),
		healthCheckInterval: viper.GetInt(configHealthCheckInterval),
		vipService:          vipService,
	}

	// Health check validation goroutine worker.
	go vipServer.healthCheckWorker(healthServer)

	return vipServer
}

// serverLoggerInterceptor configures the logger interceptor for the vip server.
func serverLoggerInterceptor(logger *logrus.Logger) []grpc.ServerOption {
	// Create new logrus entry for logger interceptor.
	logrusEntry := logrus.NewEntry(logger)

	ignorePayload := ilogger.IgnoreServerMethodsDecider(
		append(
			strings.Split(viper.GetString(configElasticAPMIgnoreURLS), ","),
		)...,
	)

	ignoreInitialRequest := ilogger.IgnoreServerMethodsDecider(
		strings.Split(viper.GetString(configElasticAPMIgnoreURLS), ",")...,
	)

	// Shared options for the logger, with a custom gRPC code to log level function.
	loggerOpts := []grpc_logrus.Option{
		grpc_logrus.WithDecider(func(fullMethodName string, err error) bool {
			return ignorePayload(fullMethodName)
		}),
		grpc_logrus.WithLevels(grpc_logrus.DefaultClientCodeToLevel),
	}

	return ilogger.ElasticsearchLoggerServerInterceptor(
		logrusEntry,
		ignorePayload,
		ignoreInitialRequest,
		loggerOpts...,
	)
}

// healthCheckWorker is running an infinite loop that sets the serving status once
// in s.healthCheckInterval seconds.
func (s VIPServer) healthCheckWorker(healthServer *health.Server) {
	for {
		healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
		time.Sleep(time.Second * time.Duration(s.healthCheckInterval))
	}
}
