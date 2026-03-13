import { authService } from '$lib/services/auth';

// Basic init on client side
if (typeof window !== 'undefined') {
    authService.init().catch(console.error);
}

