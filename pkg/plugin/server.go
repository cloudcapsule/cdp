package plugin

import (
	"context"
	"errors"
	datapluginapi "github.com/cloudcapsule/cdp/gen/proto/go/dataplugin/v1alpha"
	"github.com/cloudcapsule/cdp/pkg/task"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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
	return &datapluginapi.SubmitDataTaskResponse{TaskStatus: &datapluginapi.TaskStatus{
		RunId: id,
		Task: &datapluginapi.Task{
			Uuid:       "task-uuid",
			Name:       "pg-backup",
			TaskParams: request.Task.TaskParams,
		},
	}}, nil
}

func (s *DataPluginService) DataTaskStatus(ctx context.Context, request *datapluginapi.DataTaskStatusRequest) (*datapluginapi.DataTaskStatusResponse, error) {
	t := executor.GetActiveTask(request.RunId)
	if t == nil {
		return &datapluginapi.DataTaskStatusResponse{}, status.Error(codes.NotFound, errors.New("not found").Error())
	}
	response := &datapluginapi.DataTaskStatusResponse{
		TaskStatus: &datapluginapi.TaskStatus{
			RunId: request.RunId,
			State: string(t.GetState()),
			Task: &datapluginapi.Task{
				Uuid: t.GetUUID(),
				Name: t.GetName(),
			},
		}}
	return response, nil
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
