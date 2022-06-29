package plugin

import (
	"context"
	"errors"
	dpv1alpha "github.com/cloudcapsule/cdp/gen/proto/go/dataplugin/v1alpha"
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
	dpv1alpha.UnimplementedDataPluginServiceServer
}

var executor *task.Executor

func (s *DataPluginService) Registration(ctx context.Context, request *dpv1alpha.RegistrationRequest) (*dpv1alpha.RegistrationResponse, error) {
	response := &dpv1alpha.RegistrationResponse{}
	response.PluginId = viper.GetString("plugin-id")
	for _, dt := range executor.GetRegisteredTask() {
		t := &dpv1alpha.DataTask{
			Uuid:       dt.GetName(),
			Name:       dt.GetName(),
			TaskParams: dt.InputParams(),
		}
		response.DataTasks = append(response.DataTasks, t)
	}
	return response, nil
}

func (s *DataPluginService) SubmitDataTask(ctx context.Context, request *dpv1alpha.SubmitDataTaskRequest) (*dpv1alpha.SubmitDataTaskResponse, error) {
	id := uuid.New().String()
	pg := task.NewPGTask(id)
	executor.ExecutionQueue <- pg
	return &dpv1alpha.SubmitDataTaskResponse{TaskStatus: &dpv1alpha.TaskStatus{
		RunId: id,
		Task: &dpv1alpha.DataTask{
			Uuid:       "task-uuid",
			Name:       "pg-backup",
			TaskParams: request.Task.TaskParams,
		},
	}}, nil
}

func (s *DataPluginService) DataTaskStatus(ctx context.Context, request *dpv1alpha.DataTaskStatusRequest) (*dpv1alpha.DataTaskStatusResponse, error) {
	t := executor.GetActiveTask(request.RunId)
	if t == nil {
		return &dpv1alpha.DataTaskStatusResponse{}, status.Error(codes.NotFound, errors.New("not found").Error())
	}
	response := &dpv1alpha.DataTaskStatusResponse{
		TaskStatus: &dpv1alpha.TaskStatus{
			RunId: request.RunId,
			State: string(t.GetState()),
			Task: &dpv1alpha.DataTask{
				Uuid: t.GetUUID(),
				Name: t.GetName(),
			},
		}}
	return response, nil
}

func (s *DataPluginService) Healthiness(ctx context.Context, request *dpv1alpha.HealthinessRequest) (*dpv1alpha.HealthinessResponse, error) {
	return &dpv1alpha.HealthinessResponse{Message: "ok"}, nil
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
		dpv1alpha.RegisterDataPluginServiceServer(grpcServer, plugin)
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
