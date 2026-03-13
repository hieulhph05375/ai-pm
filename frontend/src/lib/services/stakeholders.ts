import { apiJson, apiFetch } from './api';
import type { Category } from './categories';

export interface Stakeholder {
    id?: number;
    name: string;
    role: string;
    role_id?: number;
    role_cat?: Category;
    email: string;
    phone: string;
    organization: string;
    notes: string;
    created_at?: string;
    updated_at?: string;
}

export interface ProjectStakeholder {
    project_id: number;
    stakeholder_id: number;
    project_role: string;
    role_id?: number;
    role_cat?: Category;
    created_at: string;
    stakeholder: Stakeholder;
}

export const stakeholderService = {
    async list(page: number = 1, limit: number = 10, search: string = '') {
        const query = new URLSearchParams({
            page: page.toString(),
            limit: limit.toString(),
            search
        }).toString();
        const res = await apiJson<any>(`/api/v1/stakeholders?${query}`);

        // Handle both raw array and paginated object formats
        if (Array.isArray(res)) {
            return {
                data: res,
                total: res.length
            };
        }

        return {
            data: res?.data || [],
            total: res?.total || 0
        };
    },

    async create(data: Partial<Stakeholder>) {
        return await apiJson<Stakeholder>(`/api/v1/stakeholders`, {
            method: 'POST',
            body: JSON.stringify(data)
        });
    },

    async update(id: number, data: Partial<Stakeholder>) {
        return await apiJson<Stakeholder>(`/api/v1/stakeholders/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data)
        });
    },

    async delete(id: number) {
        await apiFetch(`/api/v1/stakeholders/${id}`, {
            method: 'DELETE'
        });
    },

    // Project Mapping
    async listByProject(projectId: number) {
        const res = await apiJson<ProjectStakeholder[]>(`/api/v1/projects/${projectId}/stakeholders`);
        return Array.isArray(res) ? res : [];
    },

    async assignToProject(projectId: number, stakeholderId: number, role: string, roleId?: number) {
        await apiJson(`/api/v1/projects/${projectId}/stakeholders`, {
            method: 'POST',
            body: JSON.stringify({ stakeholder_id: stakeholderId, role, role_id: roleId })
        });
    },

    async unassignFromProject(projectId: number, stakeholderId: number) {
        await apiFetch(`/api/v1/projects/${projectId}/stakeholders/${stakeholderId}`, {
            method: 'DELETE'
        });
    }
};
