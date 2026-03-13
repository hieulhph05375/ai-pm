import { apiJson } from './api';

export const settingService = {
    async getAll() {
        return (await apiJson<Record<string, any>>(`/api/v1/settings`)) || {};
    },

    async update(key: string, value: any) {
        return await apiJson<{ key: string, value: any }>(`/api/v1/settings/${key}`, {
            method: 'PUT',
            body: JSON.stringify({ value })
        });
    }
};
