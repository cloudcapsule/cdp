package task

import (
	datapluginv1alpha "github.com/cloudcapsule/cdp/gen/proto/go/dataplugin/v1alpha"
	log "github.com/sirupsen/logrus"
	"time"
)

type PGTask struct {
	taskUUID  string
	taskRunID string
	name      string
	state     DataTaskState
	*DataTaskLog
}

//func (t *PGTask) InfoLog(msg string) {
//	l := log.WithFields(log.Fields{"runId": t.taskRunID, "task": t.name})
//	t.log = append(t.log, &DataTaskLogEntry{Message: msg, Level: Info})
//	l.Info(msg)
//
//}
//
//func (t *PGTask) WarningLog(msg string) {
//	t.log = append(t.log, &DataTaskLogEntry{Message: msg, Level: Warning})
//	log.Warning(msg)
//}
//
//func (t *PGTask) ErrorLog(msg string) {
//	t.log = append(t.log, &DataTaskLogEntry{Message: msg, Level: Error})
//	log.Error(msg)
//}

func (t *PGTask) GetName() (taskName string) {
	return "postgresql-data-task"
}

func (t *PGTask) GetRunId() (taskId string) {
	return t.taskRunID
}

func (t *PGTask) GetUUID() (taskUUID string) {
	return t.taskUUID
}

func (t *PGTask) GetState() DataTaskState {
	return t.state
}

func (t *PGTask) Start() {
	t.state = Running
	t.Info("starting task execution")
	for i := 0; i < 10; i++ {
		t.Info("task is running...")
		time.Sleep(1 * time.Second)
	}
	t.Info("task is done...")
	t.state = Done
}

func (t *PGTask) Stop() {
	log.Info("stopping task execution")
}

func (t *PGTask) Status() {
	log.Info("reporting task status")
}

func (t *PGTask) InputParams() []*datapluginv1alpha.TaskParam {
	return []*datapluginv1alpha.TaskParam{
		{
			Name:     "Hostname",
			Label:    "PG Hostname",
			Type:     "string",
			Sensitiv: false,
			Index:    0,
		},
	}
}

func NewPGTask(runId string) *PGTask {
	return &PGTask{
		taskRunID:   runId,
		name:        "pg-backup",
		DataTaskLog: NewDataTaskLog(runId, "pg"),
	}
}
