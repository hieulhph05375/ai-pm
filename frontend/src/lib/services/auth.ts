import { writable } from 'svelte/store';

export interface User {
    id: number;
    email: string;
    fullName: string;
    role: string;
    isAdmin: boolean;
}

export interface AuthState {
    user: User | null;
    token: string | null;
    isLoading: boolean;
}

const initialState: AuthState = {
    user: null,
    token: null,
    isLoading: false
};

export const authStore = writable<AuthState>(initialState);

export const authService = {
    async init() {
        const token = localStorage.getItem('at');
        const userJson = localStorage.getItem('user');
        if (token && userJson) {
            authStore.set({
                user: JSON.parse(userJson),
                token: token,
                isLoading: false
            });
        }
    },

    async login(email: string, pass: string) {
        authStore.update(s => ({ ...s, isLoading: true }));
        try {
            const response = await fetch('/api/v1/auth/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password: pass })
            });

            if (response.ok) {
                const data = await response.json();
                localStorage.setItem('at', data.access_token);
                localStorage.setItem('user', JSON.stringify(data.user));
                authStore.set({
                    user: data.user,
                    token: data.access_token,
                    isLoading: false
                });
                return true;
            }
        } catch (e) {
            console.error(e);
        } finally {
            authStore.update(s => ({ ...s, isLoading: false }));
        }
        return false;
    },

    logout() {
        localStorage.removeItem('at');
        localStorage.removeItem('user');
        authStore.set(initialState);
    }
};
