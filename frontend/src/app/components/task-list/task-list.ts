import { Component, OnInit, signal, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Task, TaskFilter, TaskStatus, TaskPriority } from '../../model/task.model';
import { TaskService } from '../../services/task.service';
import { TaskCard } from '../task-card/task-card';
import { TaskDetails } from '../task-details/task-details';
import { TaskForm } from '../task-form/task-form';

@Component({
  selector: 'app-task-list',
  standalone: true,
  imports: [FormsModule, TaskCard, TaskDetails, TaskForm],
  templateUrl: './task-list.html',
  styleUrl: './task-list.css'
})
export class TaskList implements OnInit {
  private readonly taskService = inject(TaskService);

  readonly tasks = this.taskService.tasks;
  readonly loading = this.taskService.loading;
  readonly error = this.taskService.error;
  readonly totalCount = this.taskService.totalCount;
  readonly pendingCount = this.taskService.pendingCount;
  readonly inProgressCount = this.taskService.inProgressCount;
  readonly completedCount = this.taskService.completedCount;

  readonly selectedTask = signal<Task | null>(null);
  readonly showCreateForm = signal<boolean>(false);
  readonly isCreating = signal<boolean>(false);

  searchQuery = '';

  readonly filterStatus = signal<TaskStatus | ''>('');
  readonly filterPriority = signal<TaskPriority | ''>('');
  readonly filterDueDate = signal<string>('');

  readonly statusOptions: { value: TaskStatus | ''; label: string }[] = [
    { value: '', label: 'Todos os status' },
    { value: 'pending', label: 'Pendente' },
    { value: 'in_progress', label: 'Em progresso' },
    { value: 'completed', label: 'Concluída' },
    { value: 'cancelled', label: 'Cancelada' }
  ];

  readonly priorityOptions: { value: TaskPriority | ''; label: string }[] = [
    { value: '', label: 'Todas as prioridades' },
    { value: 'low', label: 'Baixa' },
    { value: 'medium', label: 'Média' },
    { value: 'high', label: 'Alta' }
  ];

  get filteredTasks(): Task[] {
    return this.tasks().filter(task =>
      task.title.toLowerCase().includes(this.searchQuery.toLowerCase()) ||
      task.description?.toLowerCase().includes(this.searchQuery.toLowerCase())
    );
  }

  ngOnInit(): void {
    this.loadTasks();
  }

  loadTasks(): void {
    const filter: TaskFilter = {};

    if (this.filterStatus()) filter.status = this.filterStatus() as TaskStatus;
    if (this.filterPriority()) filter.priority = this.filterPriority() as TaskPriority;
    if (this.filterDueDate()) filter.due_date = this.filterDueDate();

    this.taskService.loadTasks(filter);
  }

  onFilterChange(): void {
    this.loadTasks();
  }

  onClearFilters(): void {
    this.filterStatus.set('');
    this.filterPriority.set('');
    this.filterDueDate.set('');
    this.searchQuery = '';
    this.loadTasks();
  }

  get hasActiveFilters(): boolean {
    return !!(this.filterStatus() || this.filterPriority() || this.filterDueDate() || this.searchQuery);
  }

  onCardClicked(task: Task): void {
    this.selectedTask.set(task);
    this.showCreateForm.set(false);
  }

  onDrawerClosed(): void {
    this.selectedTask.set(null);
  }

  onNewTaskClick(): void {
    this.showCreateForm.set(true);
    this.selectedTask.set(null);
  }

  onCreateFormSubmit(input: any): void {
    this.isCreating.set(true);
    this.taskService.createTask(input).subscribe({
      next: () => {
        this.isCreating.set(false);
        this.showCreateForm.set(false);
      },
      error: () => {
        this.isCreating.set(false);
      }
    });
  }

  onCreateFormCancelled(): void {
    this.showCreateForm.set(false);
  }
}
