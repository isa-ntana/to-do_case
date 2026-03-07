package domain

func IsValidStatus(s Status) bool {
	switch s {
	case StatusPending, StatusInProgress, StatusCompleted, StatusCancelled:
		return true
	}
	return false
}

func IsValidPriority(p Priority) bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	}
	return false
}
