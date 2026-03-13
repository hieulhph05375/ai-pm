import { apiJson, apiFetch } from './api';

export interface User {
    id: number;
    email: string;
    full_name: string;
    role_id: number;
    is_admin: boolean;
    is_active: boolean;
    created_at?: string;
}

export const userService = {
    async getUsers(query: string = '', page: number = 1, limit: number = 10): Promise<{ data: User[], total: number }> {
        const urlParams = new URLSearchParams();
        if (query) urlParams.append('q', query);
        urlParams.append('page', page.toString());
        urlParams.append('limit', limit.toString());

        const url = `/api/v1/users?${urlParams.toString()}`;
        const res = await apiJson<{ data: User[], total: number }>(url);
        return {
            data: res?.data || [],
            total: res?.total || 0
        };
    },

    async createUser(user: Partial<User>, password: string): Promise<User> {
        return apiJson<User>('/api/v1/users', {
            method: 'POST',
            body: JSON.stringify({ ...user, password })
        });
    },

    async updateUser(id: number, user: Partial<User>): Promise<void> {
        const res = await apiFetch(`/api/v1/users/${id}`, {
            method: 'PUT',
            body: JSON.stringify(user)
        });
        if (!res.ok) throw new Error((await res.json()).error || 'Failed to update user');
    },

    async resetPassword(id: number, password: string): Promise<void> {
        const res = await apiFetch(`/api/v1/users/${id}/reset-password`, {
            method: 'PUT',
            body: JSON.stringify({ password })
        });
        if (!res.ok) throw new Error((await res.json()).error || 'Failed to reset password');
    },

    async toggleStatus(id: number): Promise<void> {
        const res = await apiFetch(`/api/v1/users/${id}/toggle-status`, {
            method: 'PUT'
        });
        if (!res.ok) throw new Error((await res.json()).error || 'Failed to toggle status');
    }
};

// Password validation utility
export function validatePassword(password: string): string | null {
    if (password.length < 8) return 'Password must be at least 8 characters long';
    if (!/[A-Z]/.test(password)) return 'Password must contain at least 1 uppercase letter';
    return 'Password must contain at least 1 special character (!@#$%^&*...)';
    return null; // Valid
}
