package usecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/isa-ntana/to-do_case/internal/domain"
	apperrors "github.com/isa-ntana/to-do_case/pkg/errors"
)

type TaskUseCase struct {
	repo domain.TaskRepositorys
}

func NewTaskUseCase(repo domain.TaskRepository) *TaskUseCase {
	return &TaskUseCase{repo: repo}
}

func (uc *TaskUseCase) Create(input domain.CreateTaskInput) (*domain.Task, error) {
	if err := validateDueDate(input.DueDate); err != nil {
		return nil, err
	}

	priority := input.Priority
	if priority == "" {
		priority = domain.PriorityMedium
	}

	if err := validatePriority(priority); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	task := &domain.Task{
		ID:          uuid.NewString(),
		Title:       input.Title,
		Description: input.Description,
		Status:      domain.StatusPending,
		Priority:    priority,
		DueDate:     input.DueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.repo.Create(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (uc *TaskUseCase) GetByID(id string) (*domain.Task, error) {
	return uc.repo.FindByID(id)
}

func (uc *TaskUseCase) GetAll(filter domain.TaskFilter) ([]*domain.Task, error) {
	if filter.Status != nil {
		if err := validateStatus(*filter.Status); err != nil {
			return nil, err
		}
	}

	if filter.Priority != nil {
		if err := validatePriority(*filter.Priority); err != nil {
			return nil, err
		}
	}

	tasks, err := uc.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (uc *TaskUseCase) Update(id string, input domain.UpdateTaskInput) (*domain.Task, error) {
	task, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if task.Status == domain.StatusCompleted {
		return nil, apperrors.ErrTaskCompleted
	}

	if input.Title != nil {
		task.Title = *input.Title
	}

	if input.Description != nil {
		task.Description = *input.Description
	}

	if input.Status != nil {
		if err := validateStatus(*input.Status); err != nil {
			return nil, err
		}
		task.Status = *input.Status
	}

	if input.Priority != nil {
		if err := validatePriority(*input.Priority); err != nil {
			return nil, err
		}
		task.Priority = *input.Priority
	}

	if input.DueDate != nil {
		if err := validateDueDate(*input.DueDate); err != nil {
			return nil, err
		}
		task.DueDate = *input.DueDate
	}

	task.UpdatedAt = time.Now().UTC()

	if err := uc.repo.Update(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (uc *TaskUseCase) Delete(id string) error {
	_, err := uc.repo.FindByID(id)
	if err != nil {
		return err
	}

	return uc.repo.Delete(id)
}
