<script lang="ts">
	import { authService } from '$lib/services/auth';
	import { browser } from '$app/environment';
	import favicon from '$lib/assets/favicon.svg';
	import RouteGuard from '$lib/components/RouteGuard.svelte';
	import Toast from '$lib/components/ui/Toast.svelte';
	import '../app.css';

	let { children } = $props();

	if (browser) {
		authService.init();
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<RouteGuard>
	{@render children()}
</RouteGuard>

<Toast />
