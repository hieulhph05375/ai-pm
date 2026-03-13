import { apiFetch, apiJson } from './api';

const API_URL = '/api/v1';

export interface Category {
    id: number;
    name: string;
    color: string;
    icon: string;
}

export interface Issue {
    id: number;
    project_id: number;
    type_id: number;
    priority_id: number;
    status_id: number;
    type?: Category;
    priority?: Category;
    status?: Category;
    title: string;
    description?: string | null;
    assignee_id?: number | null;
    reporter_id?: number | null;
    created_at: string;
    updated_at: string;
}

export type CreateIssuePayload = Omit<Issue, 'id' | 'created_at' | 'updated_at'>;

export const issueService = {
    async list(projectId: number, page = 1, limit = 10): Promise<{ items: Issue[], total: number }> {
        const res = await apiFetch(`${API_URL}/projects/${projectId}/issues?page=${page}&limit=${limit}`);
        if (!res.ok) throw new Error('Failed to load issues list');
        return res.json();
    },
    async create(projectId: number, payload: CreateIssuePayload): Promise<Issue> {
        return apiJson<Issue>(`${API_URL}/projects/${projectId}/issues`, {
            method: 'POST',
            body: JSON.stringify(payload)
        });
    },
    async update(projectId: number, issueId: number, payload: Partial<Issue>): Promise<Issue> {
        return apiJson<Issue>(`${API_URL}/projects/${projectId}/issues/${issueId}`, {
            method: 'PUT',
            body: JSON.stringify(payload)
        });
    },
    async delete(projectId: number, issueId: number): Promise<void> {
        await apiFetch(`${API_URL}/projects/${projectId}/issues/${issueId}`, { method: 'DELETE' });
    }
};
