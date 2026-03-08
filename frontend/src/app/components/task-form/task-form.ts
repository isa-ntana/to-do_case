import { Component, input, output, OnInit, inject } from '@angular/core';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Task, CreateTaskInput, UpdateTaskInput, TaskStatus, TaskPriority } from '../../model/task.model';

@Component({
  selector: 'app-task-form',
  standalone: true,
  imports: [ReactiveFormsModule],
  templateUrl: './task-form.html',
  styleUrl: './task-form.css'
})
export class TaskForm implements OnInit {
  private readonly formBuilder = inject(FormBuilder);

  readonly task = input<Task | null>(null);
  readonly loading = input<boolean>(false);

  readonly formSubmitted = output<CreateTaskInput | UpdateTaskInput>();
  readonly formCancelled = output<void>();

  taskForm!: FormGroup;

  readonly statusOptions: { value: TaskStatus; label: string }[] = [
    { value: 'pending', label: 'Pendente' },
    { value: 'in_progress', label: 'Em progresso' },
    { value: 'completed', label: 'Concluída' },
    { value: 'cancelled', label: 'Cancelada' }
  ];

  readonly priorityOptions: { value: TaskPriority; label: string }[] = [
    { value: 'low', label: 'Baixa' },
    { value: 'medium', label: 'Média' },
    { value: 'high', label: 'Alta' }
  ];

  get isEditMode(): boolean {
    return !!this.task();
  }

  get todayDate(): string {
    return new Date().toISOString().split('T')[0];
  }

  ngOnInit(): void {
    this.buildForm();
  }

  private buildForm(): void {
    const existingTask = this.task();

    this.taskForm = this.formBuilder.group({
      title: [
        existingTask?.title ?? '',
        [Validators.required, Validators.minLength(3), Validators.maxLength(100)]
      ],
      description: [existingTask?.description ?? ''],
      status: [existingTask?.status ?? 'pending'],
      priority: [existingTask?.priority ?? 'medium'],
      due_date: [existingTask?.due_date ?? '', Validators.required]
    });

    if (existingTask?.status === 'completed') {
      this.taskForm.disable();
    }
  }

  onSubmit(): void {
    if (this.taskForm.invalid) {
      this.taskForm.markAllAsTouched();
      return;
    }

    const formValue = this.taskForm.value;

    if (this.isEditMode) {
      const updateInput: UpdateTaskInput = {
        title: formValue.title,
        description: formValue.description || undefined,
        status: formValue.status,
        priority: formValue.priority,
        due_date: formValue.due_date
      };
      this.formSubmitted.emit(updateInput);
    } else {
      const createInput: CreateTaskInput = {
        title: formValue.title,
        description: formValue.description || undefined,
        status: formValue.status,
        priority: formValue.priority,
        due_date: formValue.due_date
      };
      this.formSubmitted.emit(createInput);
    }
  }

  onCancel(): void {
    this.formCancelled.emit();
  }

  hasError(fieldName: string, errorType: string): boolean {
    const field = this.taskForm.get(fieldName);
    return !!(field?.hasError(errorType) && field?.touched);
  }
}
