package task

import stdout "github.com/sirupsen/logrus"

type DataTaskLogEntry struct {
	Message string
	Level   DataTaskLogLevel
}

type DataTaskLogLevel int

type DataTaskLog struct {
	taskName  string
	taskRunId string
	entry     []*DataTaskLogEntry
}

var (
	Info    DataTaskLogLevel = 0
	Warning DataTaskLogLevel = 1
	Error   DataTaskLogLevel = 2
)

func (l *DataTaskLog) Log(entry *DataTaskLogEntry) {
	l.entry = append(l.entry, entry)
}

func (l *DataTaskLog) Info(msg string) {
	stdout.WithFields(stdout.Fields{"runId": l.taskRunId, "task": l.taskName}).Info(msg)
	l.Log(&DataTaskLogEntry{Message: msg, Level: Info})
}

func (l *DataTaskLog) Warning(msg string) {
	stdout.WithFields(stdout.Fields{"runId": l.taskRunId, "task": l.taskName}).Warning(msg)
	l.Log(&DataTaskLogEntry{Message: msg, Level: Warning})
}

func (l *DataTaskLog) Error(msg string) {
	stdout.WithFields(stdout.Fields{"runId": l.taskRunId, "task": l.taskName}).Error(msg)
	l.Log(&DataTaskLogEntry{Message: msg, Level: Error})

}

func NewDataTaskLog(runId, taskName string) *DataTaskLog {
	return &DataTaskLog{taskName: runId, taskRunId: taskName}
}
