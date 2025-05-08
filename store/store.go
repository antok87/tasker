// Package store abstracts the database layer behind store.Store methods.
package store

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-srvc/mods/sqlxmod"
	"github.com/go-srvc/srvc"
	"github.com/heppu/tasker/api"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var Now = time.Now

// Store wraps the sqlxmod module and provides an interface to interact with the database.
type Store struct {
	srvc.Module
	db NamedDB
}

func New(opts ...sqlxmod.Opt) *Store {
	s := &Store{}
	s.Module = sqlxmod.New(append(opts, setDB(s))...)
	return s
}

func (s *Store) Healthy(ctx context.Context) error {
	t := time.Time{}
	if err := s.db.GetContext(ctx, &t, "SELECT NOW()"); err != nil {
		return err
	}
	slog.Info("DB healthy", slog.Time("time_from_db", t))
	return nil
}

func (s *Store) AddTask(ctx context.Context, newTask api.NewTask) (*api.Task, error) {
	const q = `INSERT INTO tasks (name) VALUES (:name) RETURNING id, name`
	task := &api.Task{}
	err := s.db.NamedGetContext(ctx, task, q, newTask)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	return task, nil
}

func (s *Store) GetTasks(ctx context.Context) ([]api.Task, error) {
	const q = `SELECT id, name FROM tasks`
	tasks := []api.Task{}
	err := s.db.SelectContext(ctx, &tasks, q)
	if err != nil {
		return nil, fmt.Errorf("get all tasks: %w", err)
	}
	return tasks, nil
}

func (s *Store) DeleteTasks(ctx context.Context, id int64) error {
	const q = `DELETE FROM tasks WHERE id = $1`
	_, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	return nil
}
