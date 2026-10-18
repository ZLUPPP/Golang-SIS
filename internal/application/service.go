package application

import (
	"log.analizer/internal/analyzer"
	"log.analizer/internal/domain"
)

type Service struct {
	repository TaskRepository
	nextID     int
}

func NewService(repository TaskRepository) *Service {
	return &Service{
		repository: repository,
		nextID:     1,
	}
}

func (s *Service) CreateTask(input []domain.LogEntry) (domain.Task, error) {
	if len(input) == 0 {
		return domain.Task{}, domain.ErrEmptyInput
	}

	task := domain.Task{ID: s.nextID, Input: input, Status: domain.StatusPending}

	err := s.repository.Create(task)
	if err != nil {
		return domain.Task{}, err
	}

	s.nextID++
	return task, nil
}

func (s *Service) GetTask(id int) (domain.Task, error) {
	task, err := s.repository.Get(id)
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (s *Service) ProcessTask(id int) (domain.Task, error) {
	task, err := s.repository.Get(id)
	if err != nil {
		return domain.Task{}, err
	}

	if task.Status != domain.StatusPending {
		return domain.Task{}, domain.ErrTaskInvalidStatus
	}

	task.Status = domain.StatusProcessing
	err = s.repository.Update(task)
	if err != nil {
		return domain.Task{}, err
	}

	result, err := analyzer.Analyze(task.Input)
	if err != nil {
		return domain.Task{}, err
	}

	task.Result = &result
	task.Status = domain.StatusCompleted
	err = s.repository.Update(task)
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}
