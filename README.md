# TaskFlow

> Aplicação fullstack de gerenciamento de tarefas. O TaskFlow é uma To-Do List com CRUD completo, regras de negócio, filtros, validações e interface moderna — construído com **Go + DynamoDB** no backend e **Angular** no frontend.

<br>

## Estrutura do Repositório

```
to-do_case/
├── backend/        # API REST em Go com Gin e DynamoDB Local
└── frontend/       # Interface web em Angular 21
```

Cada parte tem seu próprio README com instruções detalhadas:

- 📖 [README do Backend](./backend/README.md)
- 📖 [README do Frontend](./frontend/README.md)

<br>

##  Tecnologias

### Backend
| Tecnologia | Descrição |
|---|---|
| Go 1.25 | Linguagem principal |
| Gin | Framework HTTP para Go |
| DynamoDB Local | Banco de dados NoSQL |
| AWS SDK v2 | Integração com DynamoDB |
| Zap | Logs estruturados em JSON |
| Docker / Docker Compose | Containerização dos serviços |

### Frontend
| Tecnologia | Descrição |
|---|---|
| Angular 21 | Framework com Signals e standalone components |
| TypeScript | Tipagem estática |
| Reactive Forms | Formulários com validação |
| CSS puro | Sistema de design com variáveis |

<br>

## Como Rodar o Projeto Completo

### Opção 1 — Docker Compose + Angular Dev Server

**1.** Sobe o backend e o DynamoDB com Docker:
```bash
cd backend
docker compose up --build
```

**2.** Em outro terminal, sobe o frontend:
```bash
cd frontend
npm install
ng serve
```

**3.** Acessa a aplicação em `http://localhost:4200`

---

### Opção 2 — Tudo local sem Docker

**Pré-requisitos:** Go 1.25+, Node.js 18+, Angular CLI, Docker (só para o DynamoDB)

**1.** Sobe o DynamoDB Local:
```bash
docker run -p 8000:8000 amazon/dynamodb-local
```

**2.** Sobe o backend:
```bash
cd backend
go mod tidy
go run ./cmd/api
```

**3.** Sobe o frontend:
```bash
cd frontend
npm install
ng serve
```

**4.** Acessa a aplicação em `http://localhost:4200`

<br>

##  Arquitetura

O projeto foi construído seguindo **Clean Architecture** no backend e separação clara de responsabilidades no frontend.

```
┌─────────────────────────────────────┐
│           Angular (porta 4200)      │
│  task-list → task-card              │
│  task-form → task-details           │
│         TaskService                 │
└──────────────┬──────────────────────┘
               │ HTTP REST
┌──────────────▼──────────────────────┐
│        Go + Gin (porta 8080)        │
│  Handler → UseCase → Repository     │
└──────────────┬──────────────────────┘
               │ AWS SDK v2
┌──────────────▼──────────────────────┐
│      DynamoDB Local (porta 8000)    │
│           Tabela: tasks             │
└─────────────────────────────────────┘
```

<br>

##  Endpoints da API

Base URL: `http://localhost:8080/api/v1`

| Método | Rota | Descrição |
|---|---|---|
| `POST` | `/tasks` | Criar tarefa |
| `GET` | `/tasks` | Listar tarefas |
| `GET` | `/tasks?status=` | Filtrar por status |
| `GET` | `/tasks?priority=` | Filtrar por prioridade |
| `GET` | `/tasks?due_date=` | Filtrar por data de vencimento |
| `GET` | `/tasks/:id` | Buscar por ID |
| `PUT` | `/tasks/:id` | Atualizar tarefa |
| `DELETE` | `/tasks/:id` | Deletar tarefa |

<br>

## Regras de Negócio

| Regra | Detalhe |
|---|---|
| **Título obrigatório** | Mínimo 3 e máximo 100 caracteres |
| **Data de vencimento obrigatória** | Formato `YYYY-MM-DD`, não pode ser no passado |
| **Status na criação** | Opcional — aceita `pending` ou `in_progress`. Padrão: `pending` |
| **Prioridade padrão** | Quando não informada, assume `medium` |
| **Status válidos** | `pending`, `in_progress`, `completed`, `cancelled` |
| **Prioridades válidas** | `low`, `medium`, `high` |
| **Tarefa completed** | Não pode ser editada — apenas deletada |
