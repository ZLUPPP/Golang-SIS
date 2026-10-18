package storage

import (
	"errors"
	"testing"

	"log.analizer/internal/domain"
)

func TestMemoryRepository_CreateAndGet(t *testing.T) {
	repository := NewMemoryRepository()

	task := domain.Task{
		ID:     1,
		Input:  []domain.LogEntry{{Level: domain.LevelInfo, Message: "first log"}},
		Status: domain.StatusPending,
	}

	err := repository.Create(task)
	if err != nil {
		t.Fatalf("Create() returned unexpected error: %v", err)
	}

	got, err := repository.Get(task.ID)
	if err != nil {
		t.Fatalf("Get() returned unexpected error: %v", err)
	}

	if got.ID != task.ID {
		t.Errorf("expected ID %d, got %d", task.ID, got.ID)
	}

	if got.Status != task.Status {
		t.Errorf("expected status %q, got %q", task.Status, got.Status)
	}
}

func TestMemoryRepository_Update(t *testing.T) {
	repository := NewMemoryRepository()

	task := domain.Task{
		ID:     1,
		Input:  []domain.LogEntry{{Level: domain.LevelInfo, Message: "first log"}},
		Status: domain.StatusProcessing,
	}

	err := repository.Create(task)
	if err != nil {
		t.Fatalf("error creating task: %v", err)
	}

	task.Status = domain.StatusCompleted

	err = repository.Update(task)
	if err != nil {
		t.Fatalf("error updating task: %v", err)
	}

	updatedTask, err := repository.Get(task.ID)
	if err != nil {
		t.Fatalf("error getting task: %v", err)
	}

	if updatedTask.Status != task.Status {
		t.Errorf("task status should be %s but is %s", task.Status, updatedTask.Status)
	}
}

func TestMemoryRepository_UpdateNotFound(t *testing.T) {
	repository := NewMemoryRepository()

	err := repository.Update(domain.Task{
		ID:     9,
		Input:  []domain.LogEntry{{Level: domain.LevelInfo, Message: "first log"}},
		Status: domain.StatusProcessing,
	})
	if !errors.Is(err, domain.ErrTaskNotFound) {
		t.Errorf("error should be %v but is %v", domain.ErrTaskNotFound, err)
	}
}
