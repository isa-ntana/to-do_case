# TaskFlow API

> API REST para gerenciamento de tarefas. Construída em **Go** com **Gin**, **DynamoDB Local** e **Clean Architecture**, a TaskFlow API oferece um CRUD completo de tarefas com regras de negócio, validações, logs estruturados, testes unitários e Docker.

<br>

## Tecnologias

| Tecnologia | Versão | Por quê |
|---|---|---|
| Go | 1.25 | Linguagem preferencial do desafio, alta performance e tipagem forte |
| Gin | v1.9.1 | Framework HTTP leve e idiomático para Go |
| DynamoDB Local | latest | Banco NoSQL conforme requisito do desafio |
| AWS SDK v2 | v1.24.0 | SDK oficial para integração com DynamoDB |
| Zap | v1.26.0 | Logs estruturados em JSON para ambiente corporativo |
| UUID | v1.5.0 | Geração de identificadores únicos para as tarefas |
| Docker / Docker Compose | - | Containerização e orquestração dos serviços |

<br>

## Arquitetura

```
Requisição HTTP
      ↓
  Handler          → recebe HTTP, valida formato, delega ao use case
      ↓
  Use Case         → aplica regras de negócio, orquestra o fluxo
      ↓
  Repository       → abstração do banco (interface no domínio)
      ↓
  DynamoDB         → implementação concreta do repositório
```

<br>

<details>
<summary><strong>Estrutura de Pastas</strong></summary>

<br>

```
to-do_case/
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       ├── main.go           # Ponto de entrada — inicializa e conecta as camadas
│   │       └── server.go         # Configuração do Gin, middlewares e rotas
│   │
│   ├── internal/
│   │   ├── domain/
│   │   │   ├── task.go           # Entidade Task e tipos Status/Priority
│   │   │   ├── inputs.go         # DTOs de entrada (CreateTaskInput, UpdateTaskInput)
│   │   │   ├── repository.go     # Interface TaskRepository (contrato com o banco)
│   │   │   └── validations.go    # Funções puras de validação do domínio
│   │   │
│   │   ├── usecase/
│   │   │   ├── task_usecase.go   # Casos de uso — lógica de negócio
│   │   │   └── validators.go     # Validações reutilizáveis (data, status, prioridade)
│   │   │
│   │   ├── repository/
│   │   │   ├── dynamo_task_repository.go  # Operações CRUD no DynamoDB
│   │   │   ├── mapper.go                  # Conversão entre entidade e item DynamoDB
│   │   │   └── filter_builder.go          # Montagem dinâmica de filtros para Scan
│   │   │
│   │   └── handler/
│   │       ├── task_handler.go   # Endpoints HTTP
│   │       └── response.go       # Helpers padronizados de resposta JSON
│   │
│   ├── pkg/
│   │   ├── errors/
│   │   │   └── errors.go         # Erros estruturados com código HTTP
│   │   └── logger/
│   │       └── logger.go         # Logger estruturado com Zap
│   │
│   ├── Dockerfile
│   ├── docker-compose.yml
│   ├── go.mod
│   └── go.sum
│
└── frontend/                     
```


</details>

<br>

## Como Rodar

### Opção 1 — Docker

Sobe a API e o DynamoDB Local juntos com um único comando:

```bash
cd backend
docker compose up --build
```

A API estará disponível em `http://localhost:8080`.

> O DynamoDB Local sobe automaticamente na porta `8000` e a tabela `tasks` é criada automaticamente na primeira execução.

Para parar:
```bash
docker compose down
```

---

### Opção 2 — Localmente sem Docker

**Pré-requisitos:** Go 1.25+

**1.** Sobe o DynamoDB Local isolado:
```bash
docker run -p 8000:8000 amazon/dynamodb-local
```

**2.** Em outro terminal, na pasta `backend/`:
```bash
go mod tidy
go run ./cmd/api
```

**Variáveis de ambiente:**

| Variável | Padrão | Descrição |
|---|---|---|
| `PORT` | `8080` | Porta da API |
| `DYNAMODB_ENDPOINT` | `http://localhost:8000` | Endpoint do DynamoDB |
| `AWS_REGION` | `us-east-1` | Região AWS |
| `GIN_MODE` | `debug` | Modo do Gin (`debug` ou `release`) |

<br>

## Endpoints

**Base URL:** `http://localhost:8080/api/v1`

| Método | Rota | Descrição | Status de Sucesso |
|---|---|---|---|
| `POST` | `/tasks` | Criar tarefa | 201 Created |
| `GET` | `/tasks` | Listar todas as tarefas | 200 OK |
| `GET` | `/tasks?status={status}` | Filtrar por status | 200 OK |
| `GET` | `/tasks?priority={priority}` | Filtrar por prioridade | 200 OK |
| `GET` | `/tasks?due_date={date}` | Filtrar por data de vencimento | 200 OK |
| `GET` | `/tasks/{id}` | Buscar tarefa por ID | 200 OK |
| `PUT` | `/tasks/{id}` | Atualizar tarefa | 200 OK |
| `DELETE` | `/tasks/{id}` | Deletar tarefa | 204 No Content |

<br>

## Regras de Negócio

| Regra | Detalhe |
|---|---|
| **Título obrigatório** | Mínimo 3 e máximo 100 caracteres |
| **Data de vencimento obrigatória** | Formato `YYYY-MM-DD`, não pode ser no passado |
| **Status na criação** | Opcional — aceita `pending` ou `in_progress`. Quando não informada, assume: `pending` |
| **Prioridade** | aceita `low`, `medium` ou `high`. Quando não informada, assume `medium` |
| **Status válidos** | `pending`, `in_progress`, `completed`, `cancelled` |
| **Tarefa completed** | Não pode ser editada — apenas deletada (retorna `422`) |

<br>

### `POST /tasks` — Criar Tarefa

O campo `status` é opcional — quando não informado assume `pending` automaticamente.

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Estudar Golang",
    "description": "Revisar conceitos de goroutines",
    "priority": "high",
    "due_date": "2026-12-31"
  }'
```

**Com status informado:**
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Estudar Golang",
    "priority": "high",
    "due_date": "2026-12-31",
    "status": "in_progress"
  }'
```

**Resposta 201:**
```json
{
  "data": {
    "id": "cc0bfc2e-9999-4e78-bd28-03f9769bc6c8",
    "title": "Estudar Golang",
    "description": "Revisar conceitos de goroutines",
    "status": "pending",
    "priority": "high",
    "due_date": "2026-12-31",
    "created_at": "2026-03-08T17:23:36.444283878Z",
    "updated_at": "2026-03-08T17:23:36.444283878Z"
  }
}
```

---

### `GET /tasks` — Listar Tarefas

```bash
# Todas as tarefas
curl http://localhost:8080/api/v1/tasks

# Filtrar por status
curl http://localhost:8080/api/v1/tasks?status=pending

# Filtrar por prioridade
curl http://localhost:8080/api/v1/tasks?priority=high

# Filtrar por data de vencimento
curl http://localhost:8080/api/v1/tasks?due_date=2026-12-31
```

---

### `GET /tasks/{id}` — Buscar por ID

```bash
curl http://localhost:8080/api/v1/tasks/cc0bfc2e-9999-4e78-bd28-03f9769bc6c8
```

---

### `PUT /tasks/{id}` — Atualizar Tarefa

Apenas os campos enviados serão atualizados — os demais permanecem inalterados:

```bash
curl -X PUT http://localhost:8080/api/v1/tasks/cc0bfc2e-9999-4e78-bd28-03f9769bc6c8 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "in_progress"
  }'
```

---

### `DELETE /tasks/{id}` — Deletar Tarefa

```bash
curl -X DELETE http://localhost:8080/api/v1/tasks/cc0bfc2e-9999-4e78-bd28-03f9769bc6c8
```

Retorna **204 No Content** sem body.

<br>

## Modelo de Dados

```json
{
  "id":          "string (UUID)",
  "title":       "string (mín. 3, máx. 100 caracteres)",
  "description": "string (opcional)",
  "status":      "pending | in_progress | completed | cancelled",
  "priority":    "low | medium | high",
  "due_date":    "string (formato YYYY-MM-DD)",
  "created_at":  "timestamp (ISO 8601, gerado automaticamente)",
  "updated_at":  "timestamp (ISO 8601, atualizado automaticamente)"
}
```

## ❌ Respostas de Erro

Todos os erros seguem o mesmo formato:

```json
{
  "error": "mensagem descritiva do erro"
}
```

| Código | Situação |
|---|---|
| `400` | Dados inválidos — título curto, data no passado, status/prioridade inválidos |
| `404` | Tarefa não encontrada |
| `422` | Tentativa de editar tarefa com status `completed` |
| `500` | Erro interno do servidor |

## Documentação Interativa (Swagger)

Com a API rodando, acesse:
```
http://localhost:8080/swagger/index.html
```

A UI permite visualizar todos os endpoints, seus parâmetros, exemplos de request e response, e testar as chamadas diretamente no browser sem precisar de Postman ou curl.

A documentação é gerada automaticamente a partir das anotações no código via [swaggo/swag](https://github.com/swaggo/swag). Para regenerar após alterações:
```bash
swag init -g cmd/api/main.go
```
<br>

## Testes Unitários

O repositório é simulado via mock, sem necessidade de DynamoDB rodando. Para rodar basta estar no diretório backend e digitar no terminal:
```bash
go test ./internal/usecase/... -v
```

Casos cobertos:

| Método | Cenário |
|---|---|
| `Create` | Criação com sucesso, data no passado, status inicial inválido, prioridade inválida, defaults de status e prioridade |
| `GetByID` | Tarefa encontrada, tarefa não encontrada |
| `Update` | Atualização com sucesso, tarefa completed, status inválido, data no passado, tarefa não encontrada |
| `Delete` | Deleção com sucesso, tarefa não encontrada |
```
