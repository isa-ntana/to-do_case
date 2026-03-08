package domain

import "time"

type Status string
type Priority string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusCancelled  Status = "cancelled"

	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type Task struct {
	ID          string    `json:"id"                     example:"550e8400-e29b-41d4-a716-446655440000"`
	Title       string    `json:"title"                  example:"Estudar Golang"`
	Description string    `json:"description,omitempty"  example:"Revisar conceitos de goroutines"`
	Status      Status    `json:"status"                 example:"pending"`
	Priority    Priority  `json:"priority"               example:"high"`
	DueDate     string    `json:"due_date,omitempty"     example:"2026-04-01"`
	CreatedAt   time.Time `json:"created_at"             example:"2026-03-08T22:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at"             example:"2026-03-08T22:00:00Z"`
}
