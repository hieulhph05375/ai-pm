import { apiFetch, apiJson } from './api';
import type { Category } from './categories';

const API_URL = '/api/v1';

export interface Risk {
    id: number;
    project_id: number;
    title: string;
    description?: string | null;
    probability: number;
    impact: number;
    risk_score: number;
    status: 'Open' | 'Mitigated' | 'Closed';
    status_id?: number;
    status_cat?: Category;
    owner_id?: number | null;
    created_at: string;
    updated_at: string;
}

export type CreateRiskPayload = Omit<Risk, 'id' | 'risk_score' | 'created_at' | 'updated_at'>;

export const riskService = {
    async list(projectId: number, page = 1, limit = 10): Promise<{ items: Risk[], total: number }> {
        const res = await apiFetch(`${API_URL}/projects/${projectId}/risks?page=${page}&limit=${limit}`);
        if (!res.ok) throw new Error('Failed to load risk register');
        return res.json();
    },
    async create(projectId: number, payload: CreateRiskPayload): Promise<Risk> {
        return apiJson<Risk>(`${API_URL}/projects/${projectId}/risks`, {
            method: 'POST',
            body: JSON.stringify(payload)
        });
    },
    async update(projectId: number, riskId: number, payload: Partial<Risk>): Promise<Risk> {
        return apiJson<Risk>(`${API_URL}/projects/${projectId}/risks/${riskId}`, {
            method: 'PUT',
            body: JSON.stringify(payload)
        });
    },
    async delete(projectId: number, riskId: number): Promise<void> {
        await apiFetch(`${API_URL}/projects/${projectId}/risks/${riskId}`, { method: 'DELETE' });
    }
};
