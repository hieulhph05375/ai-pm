import { apiFetch, apiJson } from './api';

export interface PaginatedResponse<T> {
    items: T[];
    total: number;
    page: number;
    limit: number;
}

const API_URL = '/api/v1';

export type TimesheetStatus = 'DRAFT' | 'SUBMITTED' | 'APPROVED' | 'REJECTED';

export interface Timesheet {
    id: number;
    user_id: number;
    project_id?: number | null;
    node_id?: number | null;
    task_id?: number | null;
    work_date: string;
    hours: number;
    description: string;
    status: TimesheetStatus;
    created_at: string;
    updated_at: string;

    user_name?: string;
    project_name?: string;
    node_title?: string;
    task_title?: string;
}

export const timesheetService = {
    async list(filter?: { page?: number; limit?: number; project_id?: number }): Promise<PaginatedResponse<Timesheet>> {
        const params = new URLSearchParams();
        if (filter?.page) params.append('page', String(filter.page));
        if (filter?.limit) params.append('limit', String(filter.limit));
        if (filter?.project_id) params.append('project_id', String(filter.project_id));

        const queryString = params.toString();
        const url = `${API_URL}/timesheets${queryString ? `?${queryString}` : ''}`;

        const res = await apiJson<any>(url);
        return {
            items: res.items || [],
            total: res.total || 0,
            page: res.page || 1,
            limit: res.limit || 20
        };
    },

    async getById(id: number): Promise<Timesheet> {
        const res = await apiJson<any>(`${API_URL}/timesheets/${id}`);
        return res;
    },

    async create(timesheet: Partial<Timesheet>): Promise<Timesheet> {
        const res = await apiJson<any>(`${API_URL}/timesheets`, {
            method: 'POST',
            body: JSON.stringify(timesheet)
        });
        return res;
    },

    async update(id: number, timesheet: Partial<Timesheet>): Promise<Timesheet> {
        const res = await apiJson<any>(`${API_URL}/timesheets/${id}`, {
            method: 'PUT',
            body: JSON.stringify(timesheet)
        });
        return res;
    },

    async delete(id: number): Promise<void> {
        await apiFetch(`${API_URL}/timesheets/${id}`, {
            method: 'DELETE'
        });
    }
};
