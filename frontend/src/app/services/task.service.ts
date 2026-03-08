import { Injectable, signal, computed, inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable, tap } from 'rxjs';
import {
  Task,
  CreateTaskInput,
  UpdateTaskInput,
  TaskFilter
} from '../model/task.model';

interface ApiResponse<T> {
  data: T;
}

@Injectable({
  providedIn: 'root'
})
export class TaskService {
  private readonly httpClient = inject(HttpClient);
  private readonly apiUrl = 'http://localhost:8080/api/v1/tasks';

  private readonly tasksSignal = signal<Task[]>([]);
  private readonly loadingSignal = signal<boolean>(false);
  private readonly errorSignal = signal<string | null>(null);

  readonly tasks = this.tasksSignal.asReadonly();
  readonly loading = this.loadingSignal.asReadonly();
  readonly error = this.errorSignal.asReadonly();

  readonly pendingCount = computed(() =>
    this.tasksSignal().filter(task => task.status === 'pending').length
  );

  readonly inProgressCount = computed(() =>
    this.tasksSignal().filter(task => task.status === 'in_progress').length
  );

  readonly completedCount = computed(() =>
    this.tasksSignal().filter(task => task.status === 'completed').length
  );

  readonly totalCount = computed(() => this.tasksSignal().length);

  loadTasks(filter?: TaskFilter): void {
    this.loadingSignal.set(true);
    this.errorSignal.set(null);

    let params = new HttpParams();

    if (filter?.status) {
      params = params.set('status', filter.status);
    }
    if (filter?.priority) {
      params = params.set('priority', filter.priority);
    }
    if (filter?.due_date) {
      params = params.set('due_date', filter.due_date);
    }

    this.httpClient
      .get<ApiResponse<Task[]>>(this.apiUrl, { params })
      .pipe(tap(() => this.loadingSignal.set(false)))
      .subscribe({
        next: (response) => {
          this.tasksSignal.set(response.data ?? []);
        },
        error: (error) => {
          this.errorSignal.set('Erro ao carregar tarefas');
          this.loadingSignal.set(false);
          console.error(error);
        }
      });
  }

  getTaskById(id: string): Observable<ApiResponse<Task>> {
    return this.httpClient.get<ApiResponse<Task>>(`${this.apiUrl}/${id}`);
  }

  createTask(input: CreateTaskInput): Observable<ApiResponse<Task>> {
    return this.httpClient.post<ApiResponse<Task>>(this.apiUrl, input).pipe(
      tap(() => this.loadTasks())
    );
  }

  updateTask(id: string, input: UpdateTaskInput): Observable<ApiResponse<Task>> {
    return this.httpClient.put<ApiResponse<Task>>(`${this.apiUrl}/${id}`, input).pipe(
      tap(() => this.loadTasks())
    );
  }

  deleteTask(id: string): Observable<void> {
    return this.httpClient.delete<void>(`${this.apiUrl}/${id}`).pipe(
      tap(() => this.loadTasks())
    );
  }
}
