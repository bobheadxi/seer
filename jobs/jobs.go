package jobs

// Job is the interface jobs should fulfill
type Job interface {
	Name() string
	Params() map[string]interface{}
	Unique() bool
}

const jobKey = "seer.jobkey"
