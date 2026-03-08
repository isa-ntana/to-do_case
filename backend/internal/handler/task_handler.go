package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isa-ntana/to-do_case/internal/domain"
	"github.com/isa-ntana/to-do_case/internal/usecase"
	apperrors "github.com/isa-ntana/to-do_case/pkg/errors"
	"github.com/isa-ntana/to-do_case/pkg/logger"
	"go.uber.org/zap"
)

type TaskHandler struct {
	useCase *usecase.TaskUseCase
}

func NewTaskHandler(useCase *usecase.TaskUseCase) *TaskHandler {
	return &TaskHandler{useCase: useCase}
}

func (taskHandler *TaskHandler) RegisterRoutes(routerGroup *gin.RouterGroup) {
	tasksGroup := routerGroup.Group("/tasks")
	{
		tasksGroup.POST("", taskHandler.Create)
		tasksGroup.GET("", taskHandler.List)
		tasksGroup.GET("/:id", taskHandler.GetByID)
		tasksGroup.PUT("/:id", taskHandler.Update)
		tasksGroup.DELETE("/:id", taskHandler.Delete)
	}
}

func (taskHandler *TaskHandler) Create(context *gin.Context) {
	var input domain.CreateTaskInput
	if err := context.ShouldBindJSON(&input); err != nil {
		respondError(context, &apperrors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	task, err := taskHandler.useCase.Create(input)
	if err != nil {
		logger.Error("erro ao criar tarefa", zap.Error(err))
		respondError(context, err)
		return
	}

	respondSuccess(context, http.StatusCreated, task)
}

func (taskHandler *TaskHandler) List(context *gin.Context) {
	filter := domain.TaskFilter{}

	if statusQuery := context.Query("status"); statusQuery != "" {
		status := domain.Status(statusQuery)
		filter.Status = &status
	}

	if priorityQuery := context.Query("priority"); priorityQuery != "" {
		priority := domain.Priority(priorityQuery)
		filter.Priority = &priority
	}

	tasks, err := taskHandler.useCase.GetAll(filter)
	if err != nil {
		logger.Error("erro ao listar tarefas", zap.Error(err))
		respondError(context, err)
		return
	}

	if tasks == nil {
		tasks = []*domain.Task{}
	}

	respondSuccess(context, http.StatusOK, tasks)
}

func (taskHandler *TaskHandler) GetByID(context *gin.Context) {
	id := context.Param("id")

	task, err := taskHandler.useCase.GetByID(id)
	if err != nil {
		logger.Error("erro ao buscar tarefa", zap.String("id", id), zap.Error(err))
		respondError(context, err)
		return
	}

	respondSuccess(context, http.StatusOK, task)
}

func (taskHandler *TaskHandler) Update(context *gin.Context) {
	id := context.Param("id")

	var input domain.UpdateTaskInput
	if err := context.ShouldBindJSON(&input); err != nil {
		respondError(context, &apperrors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	task, err := taskHandler.useCase.Update(id, input)
	if err != nil {
		logger.Error("erro ao atualizar tarefa", zap.String("id", id), zap.Error(err))
		respondError(context, err)
		return
	}

	respondSuccess(context, http.StatusOK, task)
}

func (taskHandler *TaskHandler) Delete(context *gin.Context) {
	id := context.Param("id")

	if err := taskHandler.useCase.Delete(id); err != nil {
		logger.Error("erro ao deletar tarefa", zap.String("id", id), zap.Error(err))
		respondError(context, err)
		return
	}

	context.Status(http.StatusNoContent)
}
