package jobs

// Job is the interface jobs should fulfill
type Job interface {
	Name() string
	Params() map[string]interface{}
}

const jobKey = "seer.jobkey"
