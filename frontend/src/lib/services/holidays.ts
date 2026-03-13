import { apiJson, apiFetch } from './api';
import type { Category } from './categories';

export interface Holiday {
    id?: number;
    name: string;
    date: string; // ISO Date YYYY-MM-DD
    type: string; // 'state' or 'company'
    type_id?: number;
    type_cat?: Category;
    is_recurring: boolean;
    created_at?: string;
    updated_at?: string;
}

export const holidayService = {
    async list(start?: string, end?: string, page: number = 1, limit: number = 10) {
        let url = `/api/v1/holidays`;
        const params = new URLSearchParams();
        if (start) params.append('start', start);
        if (end) params.append('end', end);
        params.append('page', page.toString());
        params.append('limit', limit.toString());
        url += `?${params.toString()}`;

        return (await apiJson<{ data: Holiday[], total: number }>(url)) || { data: [], total: 0 };
    },

    async create(data: Partial<Holiday>) {
        return await apiJson<Holiday>(`/api/v1/holidays`, {
            method: 'POST',
            body: JSON.stringify(data)
        });
    },

    async update(id: number, data: Partial<Holiday>) {
        return await apiJson<Holiday>(`/api/v1/holidays/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data)
        });
    },

    async delete(id: number) {
        await apiFetch(`/api/v1/holidays/${id}`, {
            method: 'DELETE'
        });
    }
};
