<script lang="ts">
	import { authStore } from '$lib/services/auth';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	let { children } = $props();

	onMount(() => {
		const unsubscribe = authStore.subscribe((state) => {
			if (!state.isLoading && !state.token && !window.location.pathname.startsWith('/login')) {
				goto('/login');
			}
		});

		return unsubscribe;
	});
</script>

{@render children()}
