package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/evergreen-ci/evergreen/model/task"
	"github.com/evergreen-ci/evergreen/rest/model"
	"github.com/evergreen-ci/utility"
)

// Type is the resolver for the type field.
func (r *podResolver) Type(ctx context.Context, obj *model.APIPod) (string, error) {
	return string(obj.Type), nil
}

// Status is the resolver for the status field.
func (r *podResolver) Status(ctx context.Context, obj *model.APIPod) (string, error) {
	return string(obj.Status), nil
}

// Task is the resolver for the task field.
func (r *podResolver) Task(ctx context.Context, obj *model.APIPod) (*model.APITask, error) {
	task, err := task.FindByIdExecution(utility.FromStringPtr(obj.TaskRuntimeInfo.RunningTaskID), obj.TaskRuntimeInfo.RunningTaskExecution)
	if err != nil {
		return nil, InternalServerError.Send(ctx, fmt.Sprintf("Error finding task %s: %s", utility.FromStringPtr(obj.TaskRuntimeInfo.RunningTaskID), err.Error()))
	}
	if task == nil {
		return nil, nil
	}
	apiTask := &model.APITask{}
	err = apiTask.BuildFromService(task, &model.APITaskArgs{
		LogURL: r.sc.GetURL(),
	})
	if err != nil {
		return nil, InternalServerError.Send(ctx, fmt.Sprintf("Error building API task from service: %s", err.Error()))
	}
	return apiTask, nil
}

// Os is the resolver for the os field.
func (r *taskContainerCreationOptsResolver) Os(ctx context.Context, obj *model.APIPodTaskContainerCreationOptions) (string, error) {
	return string(obj.OS), nil
}

// Arch is the resolver for the arch field.
func (r *taskContainerCreationOptsResolver) Arch(ctx context.Context, obj *model.APIPodTaskContainerCreationOptions) (string, error) {
	return string(obj.Arch), nil
}

// Pod returns PodResolver implementation.
func (r *Resolver) Pod() PodResolver { return &podResolver{r} }

// TaskContainerCreationOpts returns TaskContainerCreationOptsResolver implementation.
func (r *Resolver) TaskContainerCreationOpts() TaskContainerCreationOptsResolver {
	return &taskContainerCreationOptsResolver{r}
}

type podResolver struct{ *Resolver }
type taskContainerCreationOptsResolver struct{ *Resolver }