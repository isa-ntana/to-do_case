package domain

type CreateTaskInput struct {
	Title       string   `json:"title"  binding:"required,min=1,max=100"`
	Description string   `json:"description"`
	Priority    Priority `json:"priority"`
	DueDate     string   `json:"due_date"`
}

type UpdateTaskInput struct {
	Title       *string   `json:"title"`
	Description *string   `json:"description"`
	Status      *Status   `json:"status"`
	Priority    *Priority `json:"priority"`
	DueDate     *string   `json:"due_date"`
}

type TaskFilter struct {
	Status   *Status
	Priority *Priority
}
