import { apiFetch } from './api';

const API_URL = '/api/v1';

export interface WorkloadEntry {
    user_id: number;
    full_name: string;
    email: string;
    role: string;
    date: string;
    task_count: number;
    total_hours: number;
    load_percentage: number;
}

export interface ResourceWorkload {
    user_id: number;
    full_name: string;
    email: string;
    role: string;
    entries: WorkloadEntry[];
}

export interface WorkloadOverview {
    start_date: string;
    end_date: string;
    users: ResourceWorkload[];
}

export const resourceService = {
    async getWorkload(startDate: string, endDate: string): Promise<WorkloadOverview> {
        const res = await apiFetch(`${API_URL}/resources/workload?start_date=${startDate}&end_date=${endDate}`);
        if (!res.ok) throw new Error('Failed to load workload data');
        return res.json();
    }
};
