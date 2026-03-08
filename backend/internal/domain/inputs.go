package domain

type CreateTaskInput struct {
	Title       string   `json:"title"        binding:"required,min=3,max=100"  example:"Estudar Golang"`
	Description string   `json:"description"                                     example:"Revisar conceitos de goroutines"`
	Status      Status   `json:"status"                                          example:"pending"`
	Priority    Priority `json:"priority"                                        example:"high"`
	DueDate     string   `json:"due_date"     binding:"required"                 example:"2026-04-01"`
}

type UpdateTaskInput struct {
	Title       *string   `json:"title"        example:"Estudar Golang - Atualizado"`
	Description *string   `json:"description"  example:"Conteúdo atualizado"`
	Status      *Status   `json:"status"       example:"in_progress"`
	Priority    *Priority `json:"priority"     example:"medium"`
	DueDate     *string   `json:"due_date"     example:"2026-04-10"`
}

type TaskFilter struct {
	Status   *Status
	Priority *Priority
	DueDate  *string
}
