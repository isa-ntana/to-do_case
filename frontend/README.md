#  TaskFlow — Frontend

> Interface web da aplicação TaskFlow, desenvolvida em Angular 21 com design minimalista inspirado na identidade visual do Itaú.

<br>

##  Tecnologias

| Tecnologia | Versão | Por quê |
|---|---|---|
| Angular | 21.2.1 | Framework moderno com Signals e standalone components |
| TypeScript | 5.9.3 | Tipagem estática que garante consistência com o backend |
| Reactive Forms | - | Controle total sobre validação e estado do formulário |
| Angular Signals | - | Gerenciamento de estado reativo sem complexidade extra |
| Inter (Google Fonts) | - | Tipografia limpa e moderna, padrão em produtos de produtividade |
| CSS puro com variáveis | - | Sistema de design sem dependências externas |

<br>

## Arquitetura

O projeto segue a estrutura de componentes do Angular com separação clara de responsabilidades:

```
Interface do usuário
      ↓
  Components       → renderizam a UI e capturam interações
      ↓
  Services         → toda comunicação HTTP com o backend
      ↓
  Models           → contratos de dados compartilhados
```

<br>

<details>
<summary><strong>📁 Estrutura de Pastas</strong></summary>

<br>

```
frontend/
└── src/
    ├── main.ts                         # Ponto de entrada da aplicação
    ├── index.html                      # HTML base
    ├── styles.css                      # Design tokens e estilos globais
    └── app/
        ├── app.ts                      # Componente raiz
        ├── app.config.ts               # Configuração global (HttpClient, providers)
        │
        ├── models/
        │   └── task.model.ts           # Interfaces e tipos (Task, Status, Priority)
        │
        ├── services/
        │   └── task.service.ts         # Chamadas HTTP e estado global com Signals
        │
        └── components/
            ├── task-list/              # Tela principal — lista, filtros e contadores
            ├── task-card/              # Card horizontal de cada tarefa
            ├── task-details/           # Drawer lateral com detalhes e ações
            └── task-form/              # Formulário de criação e edição
```

**Por que standalone components?**
Angular 17+ adotou standalone como padrão — elimina NgModules desnecessários, cada componente declara suas próprias dependências e é mais fácil de entender isoladamente.

**Por que Signals em vez de RxJS BehaviorSubject?**
Signals são a nova primitiva reativa do Angular — mais simples, mais performáticos e com menos boilerplate. O `computed()` recalcula automaticamente valores derivados como contadores de status.

</details>

<br>

## Como Rodar

**Pré-requisitos:** Node.js 18+, Angular CLI

**1.** Instala as dependências:
```bash
cd frontend
npm install
```

**2.** Sobe o servidor de desenvolvimento:
```bash
ng serve
```

A aplicação estará disponível em `http://localhost:4200`.

> O backend precisa estar rodando em `http://localhost:8080` para a integração funcionar. Veja as instruções no [README do backend](../backend/README.md).

<br>

##  Funcionalidades

**Tela principal:**
- Título **TaskFlow** com contadores de tarefas por status
- Busca por título ou descrição em tempo real
- Filtros por status, prioridade e data de vencimento
- Cards horizontais com indicador visual de prioridade
- Destaque visual para tarefas atrasadas

**Criar tarefa:**
- Modal centralizado com formulário completo
- Validações em tempo real — título, data no passado
- Status inicial opcional (padrão: pendente)

**Detalhes da tarefa:**
- Drawer lateral animado ao clicar no card
- Exibe todos os campos incluindo datas de criação e atualização
- Ações: editar, marcar como concluída, deletar
- Confirmação antes de deletar
- Tarefas concluídas não podem ser editadas

<br>

##  Design System

Todas as decisões visuais são centralizadas em variáveis CSS no `styles.css`:

| Token | Valor | Uso |
|---|---|---|
| `--color-primary` | `#FF6200` | Laranja Itaú — botões, destaques, foco |
| `--color-background` | `#F8F8F8` | Fundo geral da aplicação |
| `--color-surface` | `#FFFFFF` | Cards e painéis |
| `--color-status-pending` | `#F59E0B` | Amarelo — tarefas pendentes |
| `--color-status-in-progress` | `#3B82F6` | Azul — em progresso |
| `--color-status-completed` | `#10B981` | Verde — concluídas |
| `--color-priority-high` | `#EF4444` | Vermelho — prioridade alta |

<br>

##  Componentes

### `task-list`
Componente principal. Orquestra todos os outros, gerencia filtros e estados de loading/erro/vazio.

### `task-card`
Card horizontal com barra lateral colorida por prioridade, badge de status e detecção automática de tarefas atrasadas.

### `task-details`
Drawer que desliza da direita com animação. Alterna entre modo visualização e modo edição sem fechar o painel.

### `task-form`
Formulário reativo usado tanto para criação quanto para edição. Validações: título obrigatório (3-100 chars), data obrigatória e não pode ser no passado.

<br>

## Integração com o Backend

O `TaskService` centraliza toda comunicação HTTP:

| Método | Endpoint | Descrição |
|---|---|---|
| `loadTasks(filter?)` | `GET /tasks` | Carrega lista com filtros opcionais |
| `getTaskById(id)` | `GET /tasks/:id` | Busca tarefa específica |
| `createTask(input)` | `POST /tasks` | Cria nova tarefa |
| `updateTask(id, input)` | `PUT /tasks/:id` | Atualiza tarefa existente |
| `deleteTask(id)` | `DELETE /tasks/:id` | Remove tarefa |

A URL base da API está em `task.service.ts`:
```typescript
private readonly apiUrl = 'http://localhost:8080/api/v1/tasks';
```