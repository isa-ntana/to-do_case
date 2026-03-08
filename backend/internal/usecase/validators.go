package usecase

import (
	"time"

	"github.com/isa-ntana/to-do_case/internal/domain"
	apperrors "github.com/isa-ntana/to-do_case/pkg/errors"
)

func validateDueDate(dueDate string) error {
	parsed, err := time.Parse("2006-01-02", dueDate)
	if err != nil {
		return apperrors.ErrInvalidDueDateFmt
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	if parsed.Before(today) {
		return apperrors.ErrDueDateInPast
	}

	return nil
}

func validatePriority(priority domain.Priority) error {
	if priority == "" {
		return nil
	}
	if !domain.IsValidPriority(priority) {
		return apperrors.ErrInvalidPriority
	}
	return nil
}

func validateStatus(status domain.Status) error {
	if !domain.IsValidStatus(status) {
		return apperrors.ErrInvalidStatus
	}
	return nil
}
