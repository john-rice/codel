package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	gmodel "github.com/semanser/ai-coder/graph/model"
	"github.com/semanser/ai-coder/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// CreateFlow is the resolver for the createFlow field.
func (r *mutationResolver) CreateFlow(ctx context.Context) (*gmodel.Flow, error) {
	flow := models.Flow{}
	tx := r.Db.Create(&flow)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &gmodel.Flow{
		ID:    flow.ID,
		Tasks: []*gmodel.Task{},
	}, nil
}

// CreateTask is the resolver for the createTask field.
func (r *mutationResolver) CreateTask(ctx context.Context, id uint, query string) (*gmodel.Task, error) {
	type InputTaskArgs struct {
		Query string `json:"query"`
	}

	args := InputTaskArgs{Query: query}
	arg, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	flowResult := r.Db.First(&models.Flow{}, id)

	if errors.Is(flowResult.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("flow with id %d not found", id)
	}

	task := models.Task{
		Type:   models.Input,
		Status: models.Finished,
		Args:   datatypes.JSON(arg),
		FlowID: id,
	}

	tx := r.Db.Create(&task)

	if tx.Error != nil {
		return nil, fmt.Errorf("failed to create task: %w", tx.Error)
	}

	return &gmodel.Task{
		ID:     task.ID,
		Type:   gmodel.TaskType(task.Type),
		Status: gmodel.TaskStatus(task.Status),
		Args:   task.Args.String(),
	}, nil
}

// Flows is the resolver for the flows field.
func (r *queryResolver) Flows(ctx context.Context) ([]*gmodel.Flow, error) {
	flows := []models.Flow{}

	tx := r.Db.Model(&models.Flow{}).Preload("Tasks").Find(&flows)

	if tx.Error != nil {
		return nil, fmt.Errorf("failed to fetch flows: %w", tx.Error)
	}

	var gFlows []*gmodel.Flow

	for _, flow := range flows {
		var gTasks []*gmodel.Task

		for _, task := range flow.Tasks {
			gTasks = append(gTasks, &gmodel.Task{
				ID:      task.ID,
				Type:    gmodel.TaskType(task.Type),
				Status:  gmodel.TaskStatus(task.Status),
				Args:    task.Args.String(),
				Results: task.Results.String(),
			})
		}

		gFlows = append(gFlows, &gmodel.Flow{
			ID:    flow.ID,
			Tasks: gTasks,
		})
	}

	return gFlows, nil
}

// TaskAdded is the resolver for the taskAdded field.
func (r *subscriptionResolver) TaskAdded(ctx context.Context) (<-chan *gmodel.Task, error) {
	panic(fmt.Errorf("not implemented: TaskAdded - taskAdded"))
}

// TaskUpdated is the resolver for the taskUpdated field.
func (r *subscriptionResolver) TaskUpdated(ctx context.Context) (<-chan *gmodel.Task, error) {
	panic(fmt.Errorf("not implemented: TaskUpdated - taskUpdated"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }