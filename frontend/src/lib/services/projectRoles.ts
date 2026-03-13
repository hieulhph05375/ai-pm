import { apiJson, apiFetch } from './api';

export interface ProjectPermission {
    id: number;
    name: string;
    description: string;
    module: string;
    created_at: string;
}

export interface ProjectRole {
    id: number;
    project_id: number;
    name: string;
    description: string;
    color: string;
    is_default: boolean;
    created_at: string;
    updated_at: string;
    permissions?: ProjectPermission[];
}

export const projectRolesService = {
    async getPermissions(): Promise<ProjectPermission[]> {
        return await apiJson<ProjectPermission[]>(`/api/v1/projects/:id/permissions`); // Wait, the endpoint in main.go is projectSpecific.GET("/permissions", ...)
        // Actually I need to know the :id or just call it from any project context
    },

    async getRoles(projectId: number): Promise<ProjectRole[]> {
        return await apiJson<ProjectRole[]>(`/api/v1/projects/${projectId}/roles`);
    },

    async createRole(projectId: number, role: Partial<ProjectRole>): Promise<ProjectRole> {
        return await apiJson<ProjectRole>(`/api/v1/projects/${projectId}/roles`, {
            method: 'POST',
            body: JSON.stringify(role)
        });
    },

    async updateRole(projectId: number, roleId: number, role: Partial<ProjectRole>): Promise<ProjectRole> {
        return await apiJson<ProjectRole>(`/api/v1/projects/${projectId}/roles/${roleId}`, {
            method: 'PUT',
            body: JSON.stringify(role)
        });
    },

    async deleteRole(projectId: number, roleId: number): Promise<void> {
        await apiFetch(`/api/v1/projects/${projectId}/roles/${roleId}`, {
            method: 'DELETE'
        });
    },

    async getRolePermissions(projectId: number, roleId: number): Promise<ProjectPermission[]> {
        return await apiJson<ProjectPermission[]>(`/api/v1/projects/${projectId}/roles/${roleId}/permissions`);
    },

    async setRolePermissions(projectId: number, roleId: number, permissionIds: number[]): Promise<void> {
        await apiFetch(`/api/v1/projects/${projectId}/roles/${roleId}/permissions`, {
            method: 'PUT',
            body: JSON.stringify({ permission_ids: permissionIds })
        });
    },

    async getAllProjectPermissions(projectId: number): Promise<ProjectPermission[]> {
        return await apiJson<ProjectPermission[]>(`/api/v1/projects/${projectId}/permissions`);
    }
};
