package usecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/isa-ntana/to-do_case/internal/domain"
	apperrors "github.com/isa-ntana/to-do_case/pkg/errors"
)

type TaskUseCase struct {
	repo domain.TaskRepository
}

func NewTaskUseCase(repo domain.TaskRepository) *TaskUseCase {
	return &TaskUseCase{repo: repo}
}

func (useCase *TaskUseCase) Create(input domain.CreateTaskInput) (*domain.Task, error) {
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

	status := input.Status
	if status == "" {
		status = domain.StatusPending
	}
	if status != domain.StatusPending && status != domain.StatusInProgress {
		return nil, &apperrors.AppError{
			Code:    400,
			Message: "status inicial deve ser pending ou in_progress",
		}
	}

	now := time.Now().UTC()
	task := &domain.Task{
		ID:          uuid.NewString(),
		Title:       input.Title,
		Description: input.Description,
		Status:      status,
		Priority:    priority,
		DueDate:     input.DueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := useCase.repo.Create(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (useCase *TaskUseCase) GetByID(id string) (*domain.Task, error) {
	return useCase.repo.FindByID(id)
}

func (useCase *TaskUseCase) GetAll(filter domain.TaskFilter) ([]*domain.Task, error) {
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

	tasks, err := useCase.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (useCase *TaskUseCase) Update(id string, input domain.UpdateTaskInput) (*domain.Task, error) {
	task, err := useCase.repo.FindByID(id)
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

	if err := useCase.repo.Update(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (useCase *TaskUseCase) Delete(id string) error {
	_, err := useCase.repo.FindByID(id)
	if err != nil {
		return err
	}

	return useCase.repo.Delete(id)
}
