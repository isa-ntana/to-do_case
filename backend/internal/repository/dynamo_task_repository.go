package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/isa-ntana/to-do_case/internal/domain"
	apperrors "github.com/isa-ntana/to-do_case/pkg/errors"
	"github.com/isa-ntana/to-do_case/pkg/logger"
	"go.uber.org/zap"
)

const tableName = "tasks"

type DynamoTaskRepository struct {
	client *dynamodb.Client
}

func NewDynamoTaskRepository(client *dynamodb.Client) *DynamoTaskRepository {
	return &DynamoTaskRepository{client: client}
}

func (r *DynamoTaskRepository) Create(task *domain.Task) error {
	av, err := attributevalue.MarshalMap(toItem(task))
	if err != nil {
		logger.Error("falha ao serializar tarefa", zap.Error(err))
		return apperrors.ErrInternalServer
	}

	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})
	if err != nil {
		logger.Error("falha ao inserir tarefa", zap.Error(err))
		return apperrors.ErrInternalServer
	}

	logger.Info("tarefa criada", zap.String("id", task.ID))
	return nil
}

func (r *DynamoTaskRepository) FindByID(id string) (*domain.Task, error) {
	result, err := r.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		logger.Error("falha ao buscar tarefa", zap.String("id", id), zap.Error(err))
		return nil, apperrors.ErrInternalServer
	}

	if result.Item == nil {
		return nil, apperrors.ErrTaskNotFound
	}

	var item taskItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		logger.Error("falha ao desserializar tarefa", zap.Error(err))
		return nil, apperrors.ErrInternalServer
	}

	return fromItem(item), nil
}

func (r *DynamoTaskRepository) FindAll(filter domain.TaskFilter) ([]*domain.Task, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	if f := buildScanFilter(filter); f != nil {
		input.FilterExpression = aws.String(f.Expression)
		input.ExpressionAttributeNames = f.AttributeNames
		input.ExpressionAttributeValues = f.AttributeValues
	}

	result, err := r.client.Scan(context.TODO(), input)
	if err != nil {
		logger.Error("falha ao listar tarefas", zap.Error(err))
		return nil, apperrors.ErrInternalServer
	}

	tasks := make([]*domain.Task, 0, len(result.Items))
	for _, av := range result.Items {
		var item taskItem
		if err := attributevalue.UnmarshalMap(av, &item); err != nil {
			logger.Warn("falha ao desserializar item", zap.Error(err))
			continue
		}
		tasks = append(tasks, fromItem(item))
	}

	return tasks, nil
}

func (r *DynamoTaskRepository) Update(task *domain.Task) error {
	av, err := attributevalue.MarshalMap(toItem(task))
	if err != nil {
		logger.Error("falha ao serializar tarefa para update", zap.Error(err))
		return apperrors.ErrInternalServer
	}

	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})
	if err != nil {
		logger.Error("falha ao atualizar tarefa", zap.String("id", task.ID), zap.Error(err))
		return apperrors.ErrInternalServer
	}

	logger.Info("tarefa atualizada", zap.String("id", task.ID))
	return nil
}

func (r *DynamoTaskRepository) Delete(id string) error {
	_, err := r.client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		logger.Error("falha ao deletar tarefa", zap.String("id", id), zap.Error(err))
		return apperrors.ErrInternalServer
	}

	logger.Info("tarefa deletada", zap.String("id", id))
	return nil
}
