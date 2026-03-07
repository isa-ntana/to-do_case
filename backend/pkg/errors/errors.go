package errors

import "fmt"

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("code %d: %s", e.Code, e.Message)
}

var (
	ErrTaskNotFound      = &AppError{Code: 404, Message: "tarefa não encontrada"}
	ErrTaskCompleted     = &AppError{Code: 422, Message: "tarefas concluídas não podem ser editadas"}
	ErrInvalidStatus     = &AppError{Code: 400, Message: "status inválido: use pending, in_progress, completed ou cancelled"}
	ErrInvalidPriority   = &AppError{Code: 400, Message: "prioridade inválida: use low, medium ou high"}
	ErrDueDateInPast     = &AppError{Code: 400, Message: "a data de vencimento não pode ser no passado"}
	ErrInvalidDueDateFmt = &AppError{Code: 400, Message: "formato de data inválido: use YYYY-MM-DD"}
	ErrTitleRequired     = &AppError{Code: 400, Message: "o título é obrigatório (mínimo 3, máximo 100 caracteres)"}
	ErrInternalServer    = &AppError{Code: 500, Message: "erro interno do servidor"}
)

func IsNotFound(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == 404
	}
	return false
}
