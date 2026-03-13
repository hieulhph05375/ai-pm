import { get } from 'svelte/store';
import { authStore, authService } from './auth';
import { goto } from '$app/navigation';

/**
 * Centralized API fetch that auto-handles 401 → redirect to login.
 */
export async function apiFetch(url: string, options: RequestInit = {}): Promise<Response> {
    const token = get(authStore).token;

    const headers: Record<string, string> = {
        ...(options.headers as Record<string, string> || {}),
    };

    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    if (options.body && !headers['Content-Type']) {
        headers['Content-Type'] = 'application/json';
    }

    const res = await fetch(url, { ...options, headers });

    if (res.status === 401) {
        // Session expired — clear auth
        authService.logout();

        // Hard redirect to login to ensure clean state
        if (typeof window !== 'undefined') {
            const currentPath = encodeURIComponent(window.location.pathname + window.location.search);
            const redirectUrl = `/login?redirectTo=${currentPath}`;

            // Try SvelteKit navigation first, fallback to hard reload if it fails
            goto(redirectUrl).catch(() => {
                window.location.href = redirectUrl;
            });
        }

        // Throw a specific error that can be ignored by global toast handlers
        const error = new Error('Session expired');
        (error as any).isAuthError = true;
        throw error;
    }

    return res;
}

/**
 * Helper to get JSON body, throwing readable error from API response.
 */
export async function apiJson<T>(url: string, options: RequestInit = {}): Promise<T> {
    const res = await apiFetch(url, options);
    if (!res.ok) {
        let errMsg = `Error ${res.status}`;
        try {
            const body = await res.json();
            errMsg = body.error || errMsg;
        } catch { }
        throw new Error(errMsg);
    }
    return res.json() as Promise<T>;
}
