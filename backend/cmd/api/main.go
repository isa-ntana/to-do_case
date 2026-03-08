package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/isa-ntana/to-do_case/internal/handler"
	"github.com/isa-ntana/to-do_case/internal/repository"
	"github.com/isa-ntana/to-do_case/internal/usecase"
	"github.com/isa-ntana/to-do_case/pkg/logger"

	_ "github.com/isa-ntana/to-do_case/docs"
)

// @title           TaskFlow API
// @version         1.0
// @description     API REST para gerenciamento de tarefas (To-Do List). Desenvolvida como desafio técnico para o Itaú.

// @contact.name    Isabela Santana
// @contact.url     https://github.com/isa-ntana/to-do_case

// @host            localhost:8080
// @BasePath        /api/v1

// @schemes         http
func main() {
	logger.Init()
	defer logger.Sync()

	dynamoClient := newDynamoClient()
	ensureTableExists(dynamoClient)

	taskRepository := repository.NewDynamoTaskRepository(dynamoClient)
	taskUseCase := usecase.NewTaskUseCase(taskRepository)
	taskHandler := handler.NewTaskHandler(taskUseCase)

	router := setupServer(taskHandler)

	port := getEnv("PORT", "8080")
	log.Printf("servidor rodando na porta %s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("falha ao iniciar servidor: %v", err)
	}
}

func newDynamoClient() *dynamodb.Client {
	endpoint := getEnv("DYNAMODB_ENDPOINT", "http://localhost:8000")
	region := getEnv("AWS_REGION", "us-east-1")

	awsConfig, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("local", "local", ""),
		),
	)
	if err != nil {
		log.Fatalf("falha ao carregar configuração AWS: %v", err)
	}

	return dynamodb.NewFromConfig(awsConfig, func(options *dynamodb.Options) {
		options.BaseEndpoint = aws.String(endpoint)
	})
}

func ensureTableExists(dynamoClient *dynamodb.Client) {
	_, err := dynamoClient.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String("tasks"),
	})
	if err == nil {
		logger.Info("tabela tasks já existe")
		return
	}

	_, err = dynamoClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String("tasks"),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		log.Fatalf("falha ao criar tabela: %v", err)
	}

	logger.Info("tabela tasks criada com sucesso")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
