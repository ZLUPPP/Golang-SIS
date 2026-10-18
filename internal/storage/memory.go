package storage

import "log.analizer/internal/domain"

type MemoryRepository struct {
	tasks map[int]domain.Task
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{tasks: make(map[int]domain.Task)}
}

func (r *MemoryRepository) Create(task domain.Task) error {
	_, ok := r.tasks[task.ID]
	if ok {
		return domain.ErrTaskAlreadyExists
	}

	r.tasks[task.ID] = task
	return nil
}

func (r *MemoryRepository) Get(id int) (domain.Task, error) {
	val, ok := r.tasks[id]
	if !ok {
		return domain.Task{}, domain.ErrTaskNotFound
	}
	return val, nil
}

func (r *MemoryRepository) Update(task domain.Task) error {
	_, ok := r.tasks[task.ID]
	if !ok {
		return domain.ErrTaskNotFound
	}

	r.tasks[task.ID] = task
	return nil
}
