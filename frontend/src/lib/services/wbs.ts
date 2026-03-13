import { apiFetch, apiJson } from './api';
import { authService } from './auth';

const API_URL = '/api/v1';

export type WBSNodeType = 'Phase' | 'Milestone' | 'Task' | 'Sub-task';
export type DependencyType = 'FS' | 'SF' | 'SS' | 'FF';

export interface WBSBaseline {
    id: number;
    project_id: number;
    name: string;
    description: string;
    created_by: number;
    created_at: string;
}

export interface WBSBaselineNode {
    baseline_id: number;
    node_id: number;
    path: string;
    planned_start_date?: string | null;
    planned_end_date?: string | null;
    progress: number;
    planned_value: number;
    actual_cost: number;
}

export interface WBSNode {
    id: number;
    project_id: number;
    title: string;
    type: WBSNodeType;
    type_id?: number | null;
    type_cat?: {
        id: number;
        name: string;
        color?: string | null;
        icon?: string | null;
    } | null;
    path: string;
    order_index: number;
    planned_start_date?: string | null;
    planned_end_date?: string | null;
    actual_start_date?: string | null;
    actual_end_date?: string | null;
    progress: number;
    planned_value: number;
    actual_cost: number;
    estimated_effort?: number | null;
    actual_effort?: number | null;
    assigned_to?: number | null;
    description?: string | null;
    created_at?: string;
    updated_at?: string;
    children?: WBSNode[];
    has_children?: boolean;
}

export interface WBSComment {
    id: number;
    project_id: number;
    node_id: number;
    user_id: number;
    user_name?: string;
    content: string;
    created_at: string;
    updated_at: string;
}

export interface WBSDependency {
    id: number;
    project_id: number;
    predecessor_id: number;
    successor_id: number;
    type: DependencyType;
    created_at?: string;
}

// Removed getHeaders as apiFetch handles it now

export const wbsService = {
    async listTree(projectId: number, filter?: {
        search?: string;
        status?: string;
        assignedTo?: number | null;
        parentPath?: string;
        fields?: string[];
        page?: number;
        limit?: number;
    }): Promise<{ data: WBSNode[], total: number }> {
        const params = new URLSearchParams();
        if (filter?.search) params.append('search', filter.search);
        if (filter?.status) params.append('status', filter.status);
        if (filter?.parentPath) params.append('parent_path', filter.parentPath);
        if (filter?.fields && filter.fields.length > 0) {
            params.append('fields', filter.fields.join(','));
        }
        if (filter?.assignedTo !== undefined && filter?.assignedTo !== null) {
            params.append('assigned_to', String(filter.assignedTo));
        }
        if (filter?.page !== undefined) params.append('page', String(filter.page));
        if (filter?.limit !== undefined) params.append('limit', String(filter.limit));

        const queryString = params.toString();
        const url = `${API_URL}/projects/${projectId}/wbs${queryString ? `?${queryString}` : ''}`;

        console.log(`[WBS Service] Fetching: ${url}`);

        const res = await apiJson<any>(url);
        return { data: res.data || [], total: res.total || 0 };
    },

    async getNode(projectId: number, nodeId: number): Promise<WBSNode> {
        const res = await apiJson<any>(`${API_URL}/projects/${projectId}/wbs/${nodeId}`);
        return res.data;
    },

    async createNode(projectId: number, node: Partial<WBSNode> & { parent_path?: string }): Promise<WBSNode> {
        const res = await apiJson<any>(`${API_URL}/projects/${projectId}/wbs`, {
            method: 'POST',
            body: JSON.stringify(node)
        });
        return res.data;
    },

    async updateNode(projectId: number, nodeId: number, node: Partial<WBSNode>): Promise<void> {
        await apiJson(`${API_URL}/projects/${projectId}/wbs/${nodeId}`, {
            method: 'PUT',
            body: JSON.stringify(node)
        });
    },

    async deleteNode(projectId: number, nodeId: number): Promise<void> {
        await apiFetch(`${API_URL}/projects/${projectId}/wbs/${nodeId}`, {
            method: 'DELETE'
        });
    },

    async listDependencies(projectId: number): Promise<WBSDependency[]> {
        const res = await apiFetch(`${API_URL}/projects/${projectId}/wbs/dependencies`);
        const json = await res.json();
        return json.data || [];
    },

    async createDependency(projectId: number, predecessorId: number, successorId: number, type: string = 'FS'): Promise<WBSDependency> {
        const res = await apiFetch(`${API_URL}/projects/${projectId}/wbs/dependencies`, {
            method: 'POST',
            body: JSON.stringify({ predecessor_id: predecessorId, successor_id: successorId, type })
        });
        const json = await res.json();
        return json.data;
    },

    async deleteDependency(projectId: number, depId: number): Promise<void> {
        await apiFetch(`${API_URL}/projects/${projectId}/wbs/dependencies/${depId}`, {
            method: 'DELETE'
        });
    },

    // Comment methods
    async getComments(projectId: number, nodeId: number): Promise<WBSComment[]> {
        const res = await apiJson<any>(`${API_URL}/projects/${projectId}/wbs/${nodeId}/comments`);
        return res.data || [];
    },

    async addComment(projectId: number, nodeId: number, content: string): Promise<WBSComment> {
        const res = await apiJson<any>(`${API_URL}/projects/${projectId}/wbs/${nodeId}/comments`, {
            method: 'POST',
            body: JSON.stringify({ content })
        });
        return res.data;
    },

    async updateComment(projectId: number, nodeId: number, commentId: number, content: string): Promise<void> {
        await apiJson(`${API_URL}/projects/${projectId}/wbs/${nodeId}/comments/${commentId}`, {
            method: 'PUT',
            body: JSON.stringify({ content })
        });
    },

    async deleteComment(projectId: number, nodeId: number, commentId: number): Promise<void> {
        await apiFetch(`${API_URL}/projects/${projectId}/wbs/${nodeId}/comments/${commentId}`, {
            method: 'DELETE'
        });
    },

    // Baseline methods
    async getBaselines(projectId: number): Promise<WBSBaseline[]> {
        const res = await apiJson<any>(`${API_URL}/projects/${projectId}/wbs-baselines`);
        return res.data || [];
    },

    async createBaseline(projectId: number, name: string, description: string): Promise<WBSBaseline> {
        const res = await apiJson<any>(`${API_URL}/projects/${projectId}/wbs-baselines`, {
            method: 'POST',
            body: JSON.stringify({ name, description })
        });
        return res.data;
    },

    async getBaselineNodes(projectId: number, baselineId: number): Promise<WBSBaselineNode[]> {
        const res = await apiJson<any>(`${API_URL}/projects/${projectId}/wbs-baselines/${baselineId}/nodes`);
        return res.data || [];
    },

    // Helper: Calculate variance in days between two dates
    calculateVarianceDays(date1?: string | null, date2?: string | null): number | null {
        if (!date1 || !date2) return null;
        const d1 = new Date(date1.split('T')[0]);
        const d2 = new Date(date2.split('T')[0]);
        const diffMs = d1.getTime() - d2.getTime();
        return Math.floor(diffMs / (1000 * 60 * 60 * 24));
    }
}
