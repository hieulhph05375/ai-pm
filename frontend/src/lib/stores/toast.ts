import { writable } from 'svelte/store';

export type ToastVariant = 'success' | 'error' | 'warning' | 'info';

export interface Toast {
    id: string;
    message: string;
    variant: ToastVariant;
}

function createToastStore() {
    const { subscribe, update } = writable<Toast[]>([]);

    function show(message: string, variant: ToastVariant = 'info', duration = 5000) {
        const id = Math.random().toString(36).substring(2, 9);
        const toast: Toast = { id, message, variant };

        update(toasts => [...toasts, toast]);

        setTimeout(() => {
            dismiss(id);
        }, duration);
    }

    function dismiss(id: string) {
        update(toasts => toasts.filter(t => t.id !== id));
    }

    return {
        subscribe,
        success: (msg: string) => show(msg, 'success'),
        error: (msg: string) => show(msg, 'error'),
        warning: (msg: string) => show(msg, 'warning'),
        info: (msg: string) => show(msg, 'info'),
        dismiss,
    };
}

export const toast = createToastStore();
