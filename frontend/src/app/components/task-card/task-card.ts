import { Component, input, output } from '@angular/core';
import { Task } from '../../model/task.model';

@Component({
  selector: 'app-task-card',
  standalone: true,
  templateUrl: './task-card.html',
  styleUrl: './task-card.css'
})
export class TaskCard {
  readonly task = input.required<Task>();

  readonly cardClicked = output<Task>();

  onCardClick(): void {
    this.cardClicked.emit(this.task());
  }

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

  isOverdue(): boolean {
    if (!this.task().due_date) return false;
    if (this.task().status === 'completed' || this.task().status === 'cancelled') return false;
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    const dueDate = new Date(this.task().due_date + 'T00:00:00');
    return dueDate < today;
  }
}
