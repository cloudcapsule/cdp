package task

import datapluginv1alpha "github.com/cloudcapsule/cdp/gen/proto/go/dataplugin/v1alpha"

type DataTask interface {
	GetName() (taskName string)
	Start()
	Stop()
	Status()
	InputParams() []*datapluginv1alpha.TaskParam
}
