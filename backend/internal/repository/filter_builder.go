package repository

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/isa-ntana/to-do_case/internal/domain"
)

type scanFilter struct {
	Expression      string
	AttributeNames  map[string]string
	AttributeValues map[string]types.AttributeValue
}

func buildScanFilter(filter domain.TaskFilter) *scanFilter {
	expressions := []string{}
	names := map[string]string{}
	values := map[string]types.AttributeValue{}

	if filter.Status != nil {
		expressions = append(expressions, "#st = :status")
		names["#st"] = "status"
		values[":status"] = &types.AttributeValueMemberS{
			Value: string(*filter.Status),
		}
	}

	if filter.Priority != nil {
		expressions = append(expressions, "#pr = :priority")
		names["#pr"] = "priority"
		values[":priority"] = &types.AttributeValueMemberS{
			Value: string(*filter.Priority),
		}
	}

	if filter.DueDate != nil {
		expressions = append(expressions, "#dd = :due_date")
		names["#dd"] = "due_date"
		values[":due_date"] = &types.AttributeValueMemberS{
			Value: *filter.DueDate,
		}
	}

	if len(expressions) == 0 {
		return nil
	}

	expr := expressions[0]
	for index := 1; index < len(expressions); index++ {
		expr = fmt.Sprintf("%s AND %s", expr, expressions[index])
	}

	return &scanFilter{
		Expression:      expr,
		AttributeNames:  names,
		AttributeValues: values,
	}
}
