package task

import (
	datapluginv1alpha "github.com/cloudcapsule/cdp/gen/proto/go/dataplugin/v1alpha"
)

type DataTask interface {
	GetUUID() string
	GetName() string
	GetRunId() string
	Start()
	Stop()
	Status()
	GetState() DataTaskState
	InputParams() []*datapluginv1alpha.TaskParam
	Log(entry *DataTaskLogEntry)
}

type DataTaskState string

var (
	Running DataTaskState = "running"
	Done    DataTaskState = "done"
	Failed  DataTaskState = "failed"
)
