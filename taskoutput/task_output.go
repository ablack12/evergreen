package taskoutput

import (
	"github.com/evergreen-ci/evergreen"
)

// TaskOutput is the versioned entry point for coordinating persistent storage
// of a task run's output data.
type TaskOutput struct {
	TaskLogs    TaskLogOutput    `bson:"task_logs,omitempty" json:"task_logs"`
	TestLogs    TestLogOutput    `bson:"test_logs,omitempty" json:"test_logs"`
	TestResults TestResultOutput `bson:"test_results,omitempty" json:"test_results"`
}

// TaskOptions represents the task-level information required for accessing
// task logs belonging to a task run.
type TaskOptions struct {
	// ProjectID is the project ID of the task run.
	ProjectID string `bson:"-" json:"-"`
	// TaskID is the task ID of the task run.
	TaskID string `bson:"-" json:"-"`
	// Execution is the execution number of the task run.
	Execution int `bson:"-" json:"-"`
	// ResultsService is the service test results are stored in.
	ResultsService string `bson:"-" json:"-"`
}

// InitializeTaskOutput initializes the task output for a new task run.
func InitializeTaskOutput(env evergreen.Environment) *TaskOutput {
	settings := env.Settings()

	return &TaskOutput{
		TaskLogs: TaskLogOutput{
			Version:      1,
			BucketConfig: settings.Buckets.LogBucket,
		},
		TestLogs: TestLogOutput{
			Version:      1,
			BucketConfig: settings.Buckets.LogBucket,
		},
		TestResults: TestResultOutput{
			Version:      0,
			BucketConfig: settings.Buckets.TestResultsBucket,
		},
	}
}
