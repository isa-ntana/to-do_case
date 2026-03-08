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

// Create godoc
// @Summary      Criar tarefa
// @Description  Cria uma nova tarefa com título, descrição, prioridade e data de vencimento
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        task  body      domain.CreateTaskInput  true  "Dados da tarefa"
// @Success      201   {object}  successResponse{data=domain.Task}
// @Failure      400   {object}  errorResponse
// @Failure      500   {object}  errorResponse
// @Router       /tasks [post]
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

// List godoc
// @Summary      Listar tarefas
// @Description  Retorna todas as tarefas, com filtros opcionais por status, prioridade e data de vencimento
// @Tags         tasks
// @Produce      json
// @Param        status    query     string  false  "Filtrar por status"      Enums(pending, in_progress, completed, cancelled)
// @Param        priority  query     string  false  "Filtrar por prioridade"  Enums(low, medium, high)
// @Param        due_date  query     string  false  "Filtrar por data de vencimento (YYYY-MM-DD)"
// @Success      200       {object}  successResponse{data=[]domain.Task}
// @Failure      500       {object}  errorResponse
// @Router       /tasks [get]
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

	if dueDateQuery := context.Query("due_date"); dueDateQuery != "" {
		filter.DueDate = &dueDateQuery
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

// GetByID godoc
// @Summary      Buscar tarefa por ID
// @Description  Retorna uma tarefa específica pelo seu ID
// @Tags         tasks
// @Produce      json
// @Param        id   path      string  true  "ID da tarefa (UUID)"
// @Success      200  {object}  successResponse{data=domain.Task}
// @Failure      404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /tasks/{id} [get]
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

// Update godoc
// @Summary      Atualizar tarefa
// @Description  Atualiza os dados de uma tarefa existente. Tarefas com status completed não podem ser editadas.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id    path      string                  true  "ID da tarefa (UUID)"
// @Param        task  body      domain.UpdateTaskInput  true  "Dados para atualização"
// @Success      200   {object}  successResponse{data=domain.Task}
// @Failure      400   {object}  errorResponse
// @Failure      404   {object}  errorResponse
// @Failure      422   {object}  errorResponse
// @Failure      500   {object}  errorResponse
// @Router       /tasks/{id} [put]
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

// Delete godoc
// @Summary      Deletar tarefa
// @Description  Remove permanentemente uma tarefa pelo seu ID
// @Tags         tasks
// @Produce      json
// @Param        id   path      string  true  "ID da tarefa (UUID)"
// @Success      204  "Tarefa deletada com sucesso"
// @Failure      404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /tasks/{id} [delete]
func (taskHandler *TaskHandler) Delete(context *gin.Context) {
	id := context.Param("id")

	if err := taskHandler.useCase.Delete(id); err != nil {
		logger.Error("erro ao deletar tarefa", zap.String("id", id), zap.Error(err))
		respondError(context, err)
		return
	}

	context.Status(http.StatusNoContent)
}
