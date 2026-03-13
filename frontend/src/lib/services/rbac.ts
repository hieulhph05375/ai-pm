import { apiJson, apiFetch } from './api';

export interface Role {
    id: number;
    name: string;
    description: string;
    created_at?: string;
}

export interface Permission {
    id: number;
    name: string;
    description: string;
}

export interface RoleWithPermissions extends Role {
    permissions: Permission[];
}

export const rbacService = {
    async listRoles(): Promise<Role[]> {
        const res = await apiJson<{ data: Role[] }>('/api/v1/roles');
        return res?.data || [];
    },

    async getRoleWithPermissions(id: number): Promise<RoleWithPermissions | null> {
        const res = await apiJson<{ data: RoleWithPermissions }>(`/api/v1/roles/${id}`);
        return res?.data || null;
    },

    async createRole(name: string, description: string): Promise<Role> {
        const res = await apiJson<{ data: Role }>('/api/v1/roles', {
            method: 'POST',
            body: JSON.stringify({ name, description })
        });
        return res.data;
    },

    async updateRole(id: number, name: string, description: string): Promise<void> {
        const res = await apiFetch(`/api/v1/roles/${id}`, {
            method: 'PUT',
            body: JSON.stringify({ name, description })
        });
        if (!res.ok) throw new Error((await res.json()).error || 'Failed to update role');
    },

    async deleteRole(id: number): Promise<void> {
        const res = await apiFetch(`/api/v1/roles/${id}`, { method: 'DELETE' });
        if (!res.ok) throw new Error((await res.json()).error || 'Failed to delete role');
    },

    async assignPermissions(roleId: number, permissionIds: number[]): Promise<void> {
        const res = await apiFetch(`/api/v1/roles/${roleId}/permissions`, {
            method: 'PUT',
            body: JSON.stringify({ permission_ids: permissionIds })
        });
        if (!res.ok) throw new Error((await res.json()).error || 'Failed to assign permissions');
    },

    async listPermissions(): Promise<Permission[]> {
        const res = await apiJson<{ data: Permission[] }>('/api/v1/permissions');
        return res?.data || [];
    }
};
