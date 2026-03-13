import { apiJson, apiFetch } from './api';
import type { Category } from './categories';

export interface Project {
    id?: number;
    project_id: string;
    project_name: string;
    description?: string;
    project_manager?: string;
    sponsor?: string;
    requesting_department?: string;
    current_phase: string;
    project_status: string;
    project_status_id?: number;
    current_phase_id?: number;
    portfolio_category_id?: number;
    overall_health_id?: number;
    priority_level_id?: number;
    project_status_cat?: Category;
    current_phase_cat?: Category;
    portfolio_category_cat?: Category;
    overall_health_cat?: Category;
    priority_level_cat?: Category;
    strategic_goal?: string;
    portfolio_category?: string;
    strategic_score: number;
    priority_level: string;
    approved_budget: number;
    actual_cost: number;
    eac: number;
    capex_opex_ratio?: string;
    expected_roi: number;
    payback_period: number;
    benefit_realization_date?: string;
    planned_start_date?: string;
    actual_start_date?: string;
    planned_end_date?: string;
    actual_end_date?: string;
    progress: number;
    overall_health: string;
    spi: number;
    cpi: number;
    last_executive_summary?: string;
    estimated_effort: number;
    actual_effort: number;
    resource_risk_flag: boolean;
    missing_skills?: string;
    systemic_risk_level: string;
    open_critical_risks: number;
    compliance_impact: string;
    dependencies_summary: string;
    last_reminder_at?: string;
    created_at?: string;
    updated_at?: string;
}

export interface PMIStats {
    pv: number;
    ev: number;
    ac: number;
    spi: number;
    cpi: number;
    eac: number;
    needs_update: boolean;
}

export interface ProjectSnapshot {
    id: number;
    project_id: number;
    spi: number;
    cpi: number;
    ev: number;
    ac: number;
    pv: number;
    progress: number;
    captured_at: string;
}

export interface MilestoneSnapshot {
    id: number;
    project_id: number;
    node_id: number;
    milestone_name: string;
    planned_date: string;
    actual_date?: string;
    captured_at: string;
}

export const projectService = {
    async list(page = 1, limit = 10, search = '', status = '') {
        const query = new URLSearchParams({
            page: page.toString(),
            limit: limit.toString(),
            search,
            status
        });
        return (await apiJson<any>(`/api/v1/projects?${query}`)) || { data: [], total: 0 };
    },

    async get(id: number) {
        return await apiJson<Project>(`/api/v1/projects/${id}`);
    },

    async create(project: Project) {
        return await apiJson<Project>(`/api/v1/projects`, {
            method: 'POST',
            body: JSON.stringify(project)
        });
    },

    async update(id: number, project: Project) {
        return await apiJson<Project>(`/api/v1/projects/${id}`, {
            method: 'PUT',
            body: JSON.stringify(project)
        });
    },

    async delete(id: number) {
        await apiFetch(`/api/v1/projects/${id}`, {
            method: 'DELETE'
        });
    },

    async getPMIStats(id: number) {
        return await apiJson<PMIStats>(`/api/v1/projects/${id}/pmi-stats`);
    },

    async getTrends(id: number) {
        return await apiJson<ProjectSnapshot[]>(`/api/v1/reporting/projects/${id}/trends`);
    },

    async getMilestoneTrends(id: number) {
        return await apiJson<MilestoneSnapshot[]>(`/api/v1/reporting/projects/${id}/milestone-trends`);
    },

    async exportWBS(id: number) {
        const res = await apiFetch(`/api/v1/projects/${id}/export/wbs`);
        const blob = await res.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `wbs_export_${id}.xlsx`;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
    },

    async exportSummary(id: number) {
        const res = await apiFetch(`/api/v1/projects/${id}/export/summary`);
        const blob = await res.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `project_summary_${id}.pdf`;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
    },

    async exportList(search = '', status = '') {
        const query = new URLSearchParams({ search, status });
        const res = await apiFetch(`/api/v1/projects/export?${query}`);
        const blob = await res.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `project_list_${new Date().toISOString().split('T')[0]}.xlsx`;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
    }
};
