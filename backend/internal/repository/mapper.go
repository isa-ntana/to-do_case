package repository

import (
	"time"

	"github.com/isa-ntana/to-do_case/internal/domain"
)

type taskItem struct {
	ID          string `dynamodbav:"id"`
	Title       string `dynamodbav:"title"`
	Description string `dynamodbav:"description"`
	Status      string `dynamodbav:"status"`
	Priority    string `dynamodbav:"priority"`
	DueDate     string `dynamodbav:"due_date"`
	CreatedAt   string `dynamodbav:"created_at"`
	UpdatedAt   string `dynamodbav:"updated_at"`
}

func toItem(t *domain.Task) taskItem {
	return taskItem{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		Priority:    string(t.Priority),
		DueDate:     t.DueDate,
		CreatedAt:   t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   t.UpdatedAt.Format(time.RFC3339),
	}
}

func fromItem(i taskItem) *domain.Task {
	createdAt, _ := time.Parse(time.RFC3339, i.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, i.UpdatedAt)

	return &domain.Task{
		ID:          i.ID,
		Title:       i.Title,
		Description: i.Description,
		Status:      domain.Status(i.Status),
		Priority:    domain.Priority(i.Priority),
		DueDate:     i.DueDate,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
