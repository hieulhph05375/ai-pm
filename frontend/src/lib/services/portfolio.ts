import { apiFetch } from './api';

const API_URL = '/api/v1';

export interface PortfolioProject {
    id: number;
    project_id: string;
    project_name: string;
    project_status: string;
    overall_health: string;
    progress: number;
    approved_budget: number;
    actual_cost: number;
    planned_end_date?: string | null;
}

export interface PortfolioOverview {
    total_projects: number;
    active_projects: number;
    completed_projects: number;
    on_hold_projects: number;
    total_budget: number;
    green_projects: number;
    yellow_projects: number;
    red_projects: number;
    high_risk_projects: PortfolioProject[];
}

export const portfolioService = {
    async getOverview(): Promise<PortfolioOverview> {
        const res = await apiFetch(`${API_URL}/portfolio/overview`);
        if (!res.ok) throw new Error('Failed to load portfolio data');
        return res.json();
    }
};
