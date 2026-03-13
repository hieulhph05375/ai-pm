import { apiJson, apiFetch } from './api';

export interface ProjectMember {
    id: number;
    project_id: number;
    user_id: number;
    project_role_id: number;
    role?: {
        name: string;
        color: string;
    };
    user?: {
        full_name: string;
        email: string;
    };
    created_at: string;
    joined_at: string;
}

export const projectMembersService = {
    async getMembers(projectId: number, page: number = 1, limit: number = 10): Promise<{ data: ProjectMember[], total: number }> {
        const res = await apiJson<{ data: ProjectMember[], total: number }>(`/api/v1/projects/${projectId}/members?page=${page}&limit=${limit}`);
        return res || { data: [], total: 0 };
    },

    async addMember(projectId: number, userId: number, roleId: number): Promise<ProjectMember> {
        return await apiJson<ProjectMember>(`/api/v1/projects/${projectId}/members`, {
            method: 'POST',
            body: JSON.stringify({ user_id: userId, role_id: roleId })
        });
    },

    async updateMemberRole(projectId: number, userId: number, roleId: number): Promise<ProjectMember> {
        return await apiJson<ProjectMember>(`/api/v1/projects/${projectId}/members/${userId}`, {
            method: 'PUT',
            body: JSON.stringify({ role_id: roleId })
        });
    },

    async removeMember(projectId: number, userId: number): Promise<void> {
        await apiFetch(`/api/v1/projects/${projectId}/members/${userId}`, {
            method: 'DELETE'
        });
    }
};
