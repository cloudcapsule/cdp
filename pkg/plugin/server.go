package plugin

import (
	"context"
	datapluginapi "github.com/cloudcapsule/cdp/gen/proto/go/dataplugin/v1alpha"
	"github.com/cloudcapsule/cdp/pkg/task"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type DataPluginService struct {
	datapluginapi.UnimplementedDataPluginServiceServer
}

var executor *task.Executor

func (s *DataPluginService) Registration(ctx context.Context, request *datapluginapi.RegistrationRequest) (*datapluginapi.RegistrationResponse, error) {
	response := &datapluginapi.RegistrationResponse{}
	for _, dt := range executor.GetRegisteredTask() {

		t := &datapluginapi.Task{
			Uuid:       dt.GetName(),
			Name:       dt.GetName(),
			TaskParams: dt.InputParams(),
		}
		response.Tasks = append(response.Tasks, t)
	}
	return response, nil
}

func (s *DataPluginService) SubmitDataTask(ctx context.Context, request *datapluginapi.SubmitDataTaskRequest) (*datapluginapi.SubmitDataTaskResponse, error) {
	id := uuid.New().String()
	pg := task.NewPGTask(id)
	executor.ExecutionQueue <- pg
	return &datapluginapi.SubmitDataTaskResponse{}, nil
}

func (s *DataPluginService) DataTaskStatus(ctx context.Context, request *datapluginapi.DataTaskStatusRequest) (*datapluginapi.DataTaskStatusResponse, error) {
	return nil, nil
}

func (s *DataPluginService) Healthiness(ctx context.Context, request *datapluginapi.HealthinessRequest) (*datapluginapi.HealthinessResponse, error) {
	return &datapluginapi.HealthinessResponse{Message: "ok"}, nil
}

func (s *DataPluginService) Serve() {
	go func() {
		lis, err := net.Listen("tcp", viper.GetString("addr"))
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("grpc server is running on: %s", viper.GetString("addr"))

		grpcServer := grpc.NewServer()
		plugin := &DataPluginService{}
		datapluginapi.RegisterDataPluginServiceServer(grpcServer, plugin)
		reflection.Register(grpcServer)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
}

func NewDataPluginService(e *task.Executor) *DataPluginService {
	executor = e
	executor.Start()
	pluginSvc := &DataPluginService{}
	return pluginSvc
}