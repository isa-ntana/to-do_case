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
	uc *usecase.TaskUseCase
}

func NewTaskHandler(uc *usecase.TaskUseCase) *TaskHandler {
	return &TaskHandler{uc: uc}
}

func (h *TaskHandler) RegisterRoutes(rg *gin.RouterGroup) {
	tasks := rg.Group("/tasks")
	{
		tasks.POST("", h.Create)
		tasks.GET("", h.List)
		tasks.GET("/:id", h.GetByID)
		tasks.PUT("/:id", h.Update)
		tasks.DELETE("/:id", h.Delete)
	}
}

func (h *TaskHandler) Create(c *gin.Context) {
	var input domain.CreateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, &apperrors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	task, err := h.uc.Create(input)
	if err != nil {
		logger.Error("erro ao criar tarefa", zap.Error(err))
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusCreated, task)
}

func (h *TaskHandler) List(c *gin.Context) {
	filter := domain.TaskFilter{}

	if s := c.Query("status"); s != "" {
		status := domain.Status(s)
		filter.Status = &status
	}

	if p := c.Query("priority"); p != "" {
		priority := domain.Priority(p)
		filter.Priority = &priority
	}

	tasks, err := h.uc.GetAll(filter)
	if err != nil {
		logger.Error("erro ao listar tarefas", zap.Error(err))
		respondError(c, err)
		return
	}

	if tasks == nil {
		tasks = []*domain.Task{}
	}

	respondSuccess(c, http.StatusOK, tasks)
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	task, err := h.uc.GetByID(id)
	if err != nil {
		logger.Error("erro ao buscar tarefa", zap.String("id", id), zap.Error(err))
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, task)
}

func (h *TaskHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var input domain.UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, &apperrors.AppError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	task, err := h.uc.Update(id, input)
	if err != nil {
		logger.Error("erro ao atualizar tarefa", zap.String("id", id), zap.Error(err))
		respondError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, task)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.uc.Delete(id); err != nil {
		logger.Error("erro ao deletar tarefa", zap.String("id", id), zap.Error(err))
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
