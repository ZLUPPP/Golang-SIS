package application

import "log.analizer/internal/domain"

type TaskRepository interface {
	Create(task domain.Task) error
	Get(id int) (domain.Task, error)
	Update(task domain.Task) error
}
