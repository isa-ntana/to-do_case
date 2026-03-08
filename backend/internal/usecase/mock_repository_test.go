package usecase

import (
	"github.com/isa-ntana/to-do_case/internal/domain"
	"github.com/stretchr/testify/mock"
)

type mockTaskRepository struct {
	mock.Mock
}

func (mockRepo *mockTaskRepository) Create(task *domain.Task) error {
	args := mockRepo.Called(task)
	return args.Error(0)
}

func (mockRepo *mockTaskRepository) FindByID(id string) (*domain.Task, error) {
	args := mockRepo.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (mockRepo *mockTaskRepository) FindAll(filter domain.TaskFilter) ([]*domain.Task, error) {
	args := mockRepo.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Task), args.Error(1)
}

func (mockRepo *mockTaskRepository) Update(task *domain.Task) error {
	args := mockRepo.Called(task)
	return args.Error(0)
}

func (mockRepo *mockTaskRepository) Delete(id string) error {
	args := mockRepo.Called(id)
	return args.Error(0)
}
