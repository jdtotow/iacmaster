package worker

type JobStatus string

const StatusInit JobStatus = "init"
const StatusReady JobStatus = "ready"
const StatusRunning JobStatus = "running"
const StatusError JobStatus = "error"
const StausFailed JobStatus = "failed"
const StatusUnknown JobStatus = "unknown"

type Worker interface {
	GetType() string
	GetName() string
	StartJob() error
	StopJob() error
	GetJobStatus() JobStatus
}
