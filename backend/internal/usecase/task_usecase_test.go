package usecase

import (
	"testing"
	"time"

	"github.com/isa-ntana/to-do_case/internal/domain"
	apperrors "github.com/isa-ntana/to-do_case/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ─── Helpers ───────────────────────────────────────────────────────────────

func futureDateString() string {
	return time.Now().AddDate(0, 0, 7).Format("2006-01-02")
}

func makeTask(status domain.Status) *domain.Task {
	return &domain.Task{
		ID:        "test-id",
		Title:     "Tarefa de teste",
		Status:    status,
		Priority:  domain.PriorityMedium,
		DueDate:   futureDateString(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// ─── Testes: Create ────────────────────────────────────────────────────────

func TestCreate_Success(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	taskUseCase := NewTaskUseCase(mockRepo)

	input := domain.CreateTaskInput{
		Title:   "Estudar Golang",
		DueDate: futureDateString(),
	}

	mockRepo.On("Create", mock.AnythingOfType("*domain.Task")).Return(nil)

	task, err := taskUseCase.Create(input)

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "Estudar Golang", task.Title)
	assert.Equal(t, domain.StatusPending, task.Status)
	assert.Equal(t, domain.PriorityMedium, task.Priority)
	mockRepo.AssertExpectations(t)
}

func TestCreate_PastDueDate_ReturnsError(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	taskUseCase := NewTaskUseCase(mockRepo)

	input := domain.CreateTaskInput{
		Title:   "Tarefa atrasada",
		DueDate: "2020-01-01",
	}

	task, err := taskUseCase.Create(input)

	assert.Nil(t, task)
	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreate_InvalidInitialStatus_ReturnsError(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	taskUseCase := NewTaskUseCase(mockRepo)

	input := domain.CreateTaskInput{
		Title:   "Tarefa inválida",
		DueDate: futureDateString(),
		Status:  domain.StatusCompleted,
	}

	task, err := taskUseCase.Create(input)

	assert.Nil(t, task)
	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreate_DefaultPriorityAndStatus(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	taskUseCase := NewTaskUseCase(mockRepo)

	input := domain.CreateTaskInput{
		Title:   "Tarefa sem prioridade",
		DueDate: futureDateString(),
	}

	mockRepo.On("Create", mock.AnythingOfType("*domain.Task")).Return(nil)

	task, err := taskUseCase.Create(input)

	assert.NoError(t, err)
	assert.Equal(t, domain.PriorityMedium, task.Priority)
	assert.Equal(t, domain.StatusPending, task.Status)
}

func TestCreate_InvalidPriority_ReturnsError(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	taskUseCase := NewTaskUseCase(mockRepo)

	input := domain.CreateTaskInput{
		Title:    "Tarefa com prioridade inválida",
		DueDate:  futureDateString(),
		Priority: domain.Priority("urgente"),
	}

	task, err := taskUseCase.Create(input)

	assert.Nil(t, task)
	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Create")
}

// ─── Testes: GetByID ───────────────────────────────────────────────────────

func TestGetByID_Success(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	taskUseCase := NewTaskUseCase(mockRepo)

	existingTask := makeTask(domain.StatusPending)
	mockRepo.On("FindByID", "test-id").Return(existingTask, nil)

	task, err := taskUseCase.GetByID("test-id")

	assert.NoError(t, err)
	assert.Equal(t, "test-id", task.ID)
	mockRepo.AssertExpectations(t)
}

func TestGetByID_NotFound_ReturnsError(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	taskUseCase := NewTaskUseCase(mockRepo)

	mockRepo.On("FindByID", "id-inexistente").Return(nil, apperrors.ErrTaskNotFound)

	task, err := taskUseCase.GetByID("id-inexistente")

	assert.Nil(t, task)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// ─── Testes: Update ────────────────────────────────────────────────────────

func TestUpdate_Success(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	taskUseCase := NewTaskUseCase(mockRepo)

	existingTask := makeTask(domain.StatusPending)
	newTitle := "Título atualizado"
	input := domain.UpdateTaskInput{Title: &newTitle}

	mockRepo.On("FindByID", "test-id").Return(existingTask, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Task")).Return(nil)

	task, err := taskUseCase.Update("test-id", input)

	assert.NoError(t, err)
	assert.Equal(t, "Título atualizado", task.Title)
	mockRepo.AssertExpectations(t)
}

func TestUpdate_CompletedTask_ReturnsError(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	taskUseCase := NewTaskUseCase(mockRepo)

	completedTask := makeTask(domain.StatusCompleted)
	newTitle := "Tentando editar"
	input := domain.UpdateTaskInput{Title: &newTitle}

	mockRepo.On("FindByID", "test-id").Return(completedTask, nil)

	task, err := taskUseCase.Update("test-id", input)

	assert.Nil(t, task)
	assert.ErrorIs(t, err, apperrors.ErrTaskCompleted)
	mockRepo.AssertNotCalled(t, "Update")
}

func TestUpdate_InvalidStatus_ReturnsError(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	taskUseCase := NewTaskUseCase(mockRepo)

	existingTask := makeTask(domain.StatusPending)
	invalidStatus := domain.Status("invalido")
	input := domain.UpdateTaskInput{Status: &invalidStatus}

	mockRepo.On("FindByID", "test-id").Return(existingTask, nil)

	task, err := taskUseCase.Update("test-id", input)

	assert.Nil(t, task)
	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Update")
}

func TestUpdate_PastDueDate_ReturnsError(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	taskUseCase := NewTaskUseCase(mockRepo)

	existingTask := makeTask(domain.StatusPending)
	pastDate := "2020-01-01"
	input := domain.UpdateTaskInput{DueDate: &pastDate}

	mockRepo.On("FindByID", "test-id").Return(existingTask, nil)

	task, err := taskUseCase.Update("test-id", input)

	assert.Nil(t, task)
	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Update")
}

func TestUpdate_TaskNotFound_ReturnsError(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	taskUseCase := NewTaskUseCase(mockRepo)

	newTitle := "Título qualquer"
	input := domain.UpdateTaskInput{Title: &newTitle}

	mockRepo.On("FindByID", "id-inexistente").Return(nil, apperrors.ErrTaskNotFound)

	task, err := taskUseCase.Update("id-inexistente", input)

	assert.Nil(t, task)
	assert.ErrorIs(t, err, apperrors.ErrTaskNotFound)
	mockRepo.AssertNotCalled(t, "Update")
}

// ─── Testes: Delete ────────────────────────────────────────────────────────

func TestDelete_Success(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	taskUseCase := NewTaskUseCase(mockRepo)

	existingTask := makeTask(domain.StatusPending)
	mockRepo.On("FindByID", "test-id").Return(existingTask, nil)
	mockRepo.On("Delete", "test-id").Return(nil)

	err := taskUseCase.Delete("test-id")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDelete_NotFound_ReturnsError(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	taskUseCase := NewTaskUseCase(mockRepo)

	mockRepo.On("FindByID", "id-inexistente").Return(nil, apperrors.ErrTaskNotFound)

	err := taskUseCase.Delete("id-inexistente")

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Delete")
}
