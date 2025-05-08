// Package app wraps the storage layer with application logic and provides an interface for API layer.
package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-srvc/srvc"
	"github.com/heppu/tasker/api"
)

const ErrServiceNotHealthy = srvc.ErrStr("service not healthy")

type Store interface {
	Healthy(context.Context) error
	AddTask(context.Context, api.NewTask) (*api.Task, error)
	GetTasks(context.Context) ([]api.Task, error)
	DeleteTasks(context.Context, int64) error
}

type App struct {
	api.UnimplementedHandler
	s Store
}

func New(s Store) *App {
	return &App{s: s}
}

func (a *App) Healthz(ctx context.Context) (*api.Healthy, error) {
	if err := a.s.Healthy(ctx); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrServiceNotHealthy, err)
	}
	return &api.Healthy{Message: "OK"}, nil
}

func (a *App) TasksPost(ctx context.Context, req *api.NewTask) (*api.Task, error) {
	return a.s.AddTask(ctx, *req)
}

func (a *App) TasksGet(ctx context.Context) ([]api.Task, error) {
	return a.s.GetTasks(ctx)
}

func (a *App) TasksIDDelete(ctx context.Context, params api.TasksIDDeleteParams) error {
	return a.s.DeleteTasks(ctx, params.ID)
}

// NewError can be used to provide custom error responses based on the error.
func (a *App) NewError(ctx context.Context, err error) *api.ErrorRespStatusCode {
	slog.Error("internal server error", slog.String("err", err.Error()))
	return &api.ErrorRespStatusCode{
		StatusCode: 500,
		Response: api.Error{
			Error: "internal server error",
		},
	}
}
