package data

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/evergreen-ci/evergreen"
	"github.com/evergreen-ci/evergreen/model/task"
	"github.com/evergreen-ci/gimlet"
	"github.com/pkg/errors"
)

// GenerateTasks parses JSON files for `generate.tasks` and creates the new builds and tasks.
func GenerateTasks(ctx context.Context, settings *evergreen.Settings, taskID string, jsonFiles []json.RawMessage) error {
	t, err := task.FindOneId(ctx, taskID)
	if err != nil {
		return errors.Wrapf(err, "finding task '%s'", taskID)
	}
	if t == nil {
		return errors.Errorf("task '%s' not found", taskID)
	}

	// Don't continue if the generator has already run
	// Return status code 400 to prevent retries
	if t.GeneratedTasks {
		return gimlet.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    evergreen.TasksAlreadyGeneratedError,
		}
	}

	var files task.GeneratedJSONFiles
	for _, f := range jsonFiles {
		files = append(files, string(f))
	}
	if _, err := task.GeneratedJSONInsertWithS3Fallback(ctx, settings, t, files, evergreen.ProjectStorageMethodDB); err != nil {
		return errors.Wrapf(err, "inserting generated JSON files for task '%s'", t.Id)
	}

	return nil
}

// GeneratePoll checks to see if a `generate.tasks` job has finished.
func GeneratePoll(ctx context.Context, taskID string) (bool, string, error) {
	t, err := task.FindOneId(ctx, taskID)
	if err != nil {
		return false, "", errors.Wrapf(err, "finding task '%s'", taskID)
	}
	if t == nil {
		return false, "", errors.Errorf("task '%s' not found", taskID)
	}

	return t.GeneratedTasks, t.GenerateTasksError, nil
}
