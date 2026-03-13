import { apiFetch, apiJson } from './api';

export interface Notification {
    id: number;
    user_id: number;
    type: string;
    title: string;
    body?: string;
    ref_id?: number;
    ref_type?: string;
    is_read: boolean;
    created_at: string;
}

export interface NotificationListResponse {
    items: Notification[];
    total: number;
    unread_count: number;
    page: number;
    limit: number;
}

export const notificationService = {
    async list(page = 1, limit = 20): Promise<NotificationListResponse> {
        const res = await apiFetch(`/api/v1/notifications?page=${page}&limit=${limit}`);
        return res.json();
    },

    async getUnreadCount(): Promise<number> {
        const res = await apiFetch('/api/v1/notifications/unread-count');
        const data = await res.json();
        return data.unread_count ?? 0;
    },

    async markRead(id: number): Promise<void> {
        await apiFetch(`/api/v1/notifications/${id}/read`, { method: 'PUT' });
    },

    async markAllRead(): Promise<void> {
        await apiFetch('/api/v1/notifications/read-all', { method: 'PUT' });
    }
};
