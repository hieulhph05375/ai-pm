import { apiJson, apiFetch } from './api';
import type { Category } from './categories';

// Matches backend entity/task.go (ADR-011)
export interface Task {
    id: number;
    title: string;
    description: string;
    status: 'TODO' | 'IN_PROGRESS' | 'DONE';
    priority: 'LOW' | 'MEDIUM' | 'HIGH' | 'URGENT';
    status_id?: number;
    priority_id?: number;
    status_cat?: Category;
    priority_cat?: Category;
    assignee_id: number | null;
    created_by: number | null;
    start_date: string | null;
    due_date: string | null;
    progress: number; // 0-100
    labels: string[];
    created_at: string;
    updated_at: string;
}

export interface TaskCreate {
    title: string;
    description?: string;
    status?: Task['status'];
    priority?: Task['priority'];
    status_id?: number;
    priority_id?: number;
    assignee_id?: number | null;
    start_date?: string | null;
    due_date?: string | null;
    progress?: number;
    labels?: string[];
}

export interface TaskActivity {
    id: number;
    task_id: number;
    actor_id: number | null;
    action: string;
    old_value: string;
    new_value: string;
    created_at: string;
}

export const TASK_STATUSES: { value: Task['status']; label: string; color: string }[] = [
    { value: 'TODO', label: 'Todo', color: 'slate' },
    { value: 'IN_PROGRESS', label: 'In Progress', color: 'amber' },
    { value: 'DONE', label: 'Done', color: 'emerald' },
];

export const TASK_PRIORITIES: { value: Task['priority']; label: string; color: string }[] = [
    { value: 'LOW', label: 'Low', color: 'slate' },
    { value: 'MEDIUM', label: 'Medium', color: 'primary' },
    { value: 'HIGH', label: 'High', color: 'amber' },
    { value: 'URGENT', label: 'Urgent', color: 'rose' },
];

class TaskService {
    private baseUrl = '/api/v1/tasks';

    async list(page = 1, limit = 10): Promise<{ items: Task[], total: number }> {
        const query = new URLSearchParams({
            page: page.toString(),
            limit: limit.toString()
        }).toString();
        const response = await apiJson<{ items: Task[], total: number }>(`${this.baseUrl}?${query}`);
        return response || { items: [], total: 0 };
    }

    async getById(id: number): Promise<Task> {
        return await apiJson<Task>(`${this.baseUrl}/${id}`);
    }

    async create(task: TaskCreate): Promise<Task> {
        return await apiJson<Task>(this.baseUrl, {
            method: 'POST',
            body: JSON.stringify(task)
        });
    }

    async update(id: number, task: Partial<TaskCreate>): Promise<Task> {
        return await apiJson<Task>(`${this.baseUrl}/${id}`, {
            method: 'PUT',
            body: JSON.stringify(task)
        });
    }

    async updateStatus(id: number, status: Task['status'], currentTask: Task): Promise<Task> {
        return this.update(id, {
            title: currentTask.title,
            description: currentTask.description,
            priority: currentTask.priority,
            status,
            assignee_id: currentTask.assignee_id,
            start_date: currentTask.start_date,
            due_date: currentTask.due_date,
            progress: currentTask.progress,
            labels: currentTask.labels,
        });
    }

    async delete(id: number): Promise<void> {
        await apiFetch(`${this.baseUrl}/${id}`, { method: 'DELETE' });
    }

    async getActivities(id: number): Promise<TaskActivity[]> {
        return (await apiJson<TaskActivity[]>(`${this.baseUrl}/${id}/activities`)) || [];
    }

    async addComment(id: number, content: string): Promise<void> {
        await apiJson(`${this.baseUrl}/${id}/comments`, {
            method: 'POST',
            body: JSON.stringify({ content })
        });
    }
}

export const taskService = new TaskService();
