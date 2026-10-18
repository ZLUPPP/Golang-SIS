package application

import (
	"errors"
	"testing"

	"log.analizer/internal/domain"
	"log.analizer/internal/storage"
)

func TestService_CreateTask(t *testing.T) {
	service := NewService(storage.NewMemoryRepository())

	task, err := service.CreateTask([]domain.LogEntry{{Level: domain.LevelInfo, Message: "test"}})
	if err != nil {
		t.Error(err)
	}

	if task.ID != 1 {
		t.Errorf("task id should be 1, was %d", task.ID)
	}

	if task.Status != domain.StatusPending {
		t.Errorf("expected task status to be pending, was %s", task.Status)
	}
}

func TestService_CreateTaskEmptyInput(t *testing.T) {
	service := NewService(storage.NewMemoryRepository())
	_, err := service.CreateTask([]domain.LogEntry{})
	if !errors.Is(err, domain.ErrEmptyInput) {
		t.Errorf("expected error to be %v, was %v", domain.ErrEmptyInput, err)
	}

}

func TestService_GetTask(t *testing.T) {
	service := NewService(storage.NewMemoryRepository())
	_, err := service.CreateTask([]domain.LogEntry{{Level: domain.LevelInfo, Message: "test"}})
	if err != nil {
		t.Error(err)
	}

	task, err := service.GetTask(1)
	if err != nil {
		t.Error(err)
	}

	if task.ID != 1 {
		t.Errorf("task id should be 1, was %d", task.ID)
	}

	_, err = service.GetTask(999)
	if !errors.Is(err, domain.ErrTaskNotFound) {
		t.Errorf("expected error to be %v, was %v", domain.ErrTaskNotFound, err)
	}
}

func TestService_ProcessTask(t *testing.T) {
	service := NewService(storage.NewMemoryRepository())

	created, err := service.CreateTask([]domain.LogEntry{
		{Level: domain.LevelInfo, Message: "started"},
		{Level: domain.LevelError, Message: "failed"},
	})
	if err != nil {
		t.Fatalf("CreateTask() error: %v", err)
	}

	processed, err := service.ProcessTask(created.ID)
	if err != nil {
		t.Fatalf("ProcessTask() error: %v", err)
	}

	if processed.Status != domain.StatusCompleted {
		t.Errorf("expected status to be completed, was %s", processed.Status)
	}

	if processed.Result == nil {
		t.Fatalf("expected result, got nil")
	}

	if processed.Result.Total != 2 {
		t.Errorf("expected total to be 2, was %d", processed.Result.Total)
	}

	if processed.Result.ErrorCount != 1 {
		t.Errorf("expected error to be 1, was %d", processed.Result.ErrorCount)
	}

	saved, err := service.GetTask(created.ID)
	if err != nil {
		t.Fatalf("GetTask() error: %v", err)
	}

	if saved.Status != domain.StatusCompleted {
		t.Errorf("expected status to be completed, was %s", saved.Status)
	}

	if saved.Result == nil {
		t.Fatal("expected result, got nil")
	}

}

func TestService_ProcessCompletedTask(t *testing.T) {
	service := NewService(storage.NewMemoryRepository())

	created, err := service.CreateTask([]domain.LogEntry{{Level: domain.LevelInfo, Message: "started"}})
	if err != nil {
		t.Fatalf("CreateTask() error: %v", err)
	}

	processed, err := service.ProcessTask(created.ID)
	if err != nil {
		t.Fatalf("ProcessTask() error: %v", err)
	}

	_, err = service.ProcessTask(processed.ID)
	if !errors.Is(err, domain.ErrTaskInvalidStatus) {
		t.Errorf("expected error to be %v, was %v", domain.ErrTaskInvalidStatus, err)
	}
}

func TestService_MultipleTasks(t *testing.T) {
	service := NewService(storage.NewMemoryRepository())

	first, err := service.CreateTask([]domain.LogEntry{{Level: domain.LevelInfo, Message: "first"}})
	if err != nil {
		t.Fatalf("CreateTask() (first) error: %v", err)
	}

	second, err := service.CreateTask([]domain.LogEntry{{Level: domain.LevelInfo, Message: "second"}})
	if err != nil {
		t.Fatalf("CreateTask() (second) error: %v", err)
	}

	if first.ID == second.ID {
		t.Errorf("first and second tasks have the same id")
	}

	_, err = service.ProcessTask(first.ID)
	if err != nil {
		t.Fatalf("ProcessTask() (first) error: %v", err)
	}

	savedFirst, _ := service.GetTask(first.ID)
	savedSecond, _ := service.GetTask(second.ID)

	if savedFirst.Status == savedSecond.Status {
		t.Errorf("first task should be completed, got %s", savedFirst.Status)
	}

	if savedFirst.Status != domain.StatusCompleted {
		t.Errorf("expected first task to be completed, was %s", savedFirst.Status)
	}

	if savedSecond.Status != domain.StatusPending {
		t.Errorf("second task should be pending, got %s", savedSecond.Status)
	}

	if savedSecond.Result != nil {
		t.Error("second task should not have a result")
	}

}
