package domain

type TaskRepository interface {
	Create(task *Task) error
	FindByID(id string) (*Task, error)
	FindAll(filter TaskFilter) ([]*Task, error)
	Update(task *Task) error
	Delete(id string) error
}
