package task

import (
	datapluginv1alpha "github.com/cloudcapsule/cdp/gen/proto/go/dataplugin/v1alpha"
	log "github.com/sirupsen/logrus"
)

type PGTask struct {
	name string
}

func (t *PGTask) GetName() (taskName string) {
	return "postgresql-data-task"
}

func (t *PGTask) Start() {
	log.Info("starting task execution")
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
