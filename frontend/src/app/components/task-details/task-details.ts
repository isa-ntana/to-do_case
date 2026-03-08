import { Component, input, output, signal, inject } from '@angular/core';
import { Task, UpdateTaskInput } from '../../model/task.model';
import { TaskService } from '../../services/task.service';
import { TaskForm } from '../task-form/task-form';

@Component({
  selector: 'app-task-details',
  standalone: true,
  imports: [TaskForm],
  templateUrl: './task-details.html',
  styleUrl: './task-details.css'
})
export class TaskDetails {
  private readonly taskService = inject(TaskService);

  readonly task = input.required<Task>();
  readonly closed = output<void>();
  readonly taskDeleted = output<void>();
  readonly taskUpdated = output<void>();

  readonly isEditMode = signal<boolean>(false);
  readonly isLoading = signal<boolean>(false);
  readonly showDeleteConfirm = signal<boolean>(false);

  getStatusLabel(status: string): string {
    const labels: Record<string, string> = {
      pending: 'Pendente',
      in_progress: 'Em progresso',
      completed: 'Concluída',
      cancelled: 'Cancelada'
    };
    return labels[status] ?? status;
  }

  getPriorityLabel(priority: string): string {
    const labels: Record<string, string> = {
      low: 'Baixa',
      medium: 'Média',
      high: 'Alta'
    };
    return labels[priority] ?? priority;
  }

  formatDate(dateString: string): string {
    if (!dateString) return '';
    const [year, month, day] = dateString.split('-');
    return `${day}/${month}/${year}`;
  }

  formatDateTime(dateString: string): string {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString('pt-BR', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }

  onClose(): void {
    this.isEditMode.set(false);
    this.showDeleteConfirm.set(false);
    this.closed.emit();
  }

  onEditClick(): void {
    this.isEditMode.set(true);
  }

  onMarkCompleted(): void {
    this.isLoading.set(true);
    this.taskService.updateTask(this.task().id, { status: 'completed' }).subscribe({
      next: () => {
        this.isLoading.set(false);
        this.taskUpdated.emit();
        this.onClose();
      },
      error: () => {
        this.isLoading.set(false);
      }
    });
  }

  onDeleteClick(): void {
    this.showDeleteConfirm.set(true);
  }

  onDeleteConfirm(): void {
    this.isLoading.set(true);
    this.taskService.deleteTask(this.task().id).subscribe({
      next: () => {
        this.isLoading.set(false);
        this.taskDeleted.emit();
        this.onClose();
      },
      error: () => {
        this.isLoading.set(false);
      }
    });
  }

  onDeleteCancel(): void {
    this.showDeleteConfirm.set(false);
  }

  onFormSubmit(input: UpdateTaskInput): void {
    this.isLoading.set(true);
    this.taskService.updateTask(this.task().id, input).subscribe({
      next: () => {
        this.isLoading.set(false);
        this.isEditMode.set(false);
        this.taskUpdated.emit();
        this.onClose();
      },
      error: () => {
        this.isLoading.set(false);
      }
    });
  }

  onFormCancelled(): void {
    this.isEditMode.set(false);
  }
}
