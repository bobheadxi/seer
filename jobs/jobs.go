package jobs

// Job is the interface jobs should fulfill
type Job interface {
	Name() string
	Params() map[string]interface{}
	Unique() bool
}

const jobKey = "seer.jobkey"

// Status indicates the state of a job
type Status string

const (
	// StatusQueued means job has been queued
	StatusQueued Status = "queued"
	// StatusRunning means the job is active
	StatusRunning Status = "running"
	// StatusRetrying means the job is scheduled for a retry
	StatusRetrying Status = "retrying"
	// StatusFailed means the job was unable to run successfully
	StatusFailed Status = "failed"
	// StatusDone means the job was able to run to completion
	StatusDone Status = "done"
)
