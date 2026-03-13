<script lang="ts">
	import Sidebar from '$lib/components/layout/Sidebar.svelte';
	import TopHeader from '$lib/components/layout/TopHeader.svelte';
	import RouteGuard from '$lib/components/RouteGuard.svelte';
	import { page } from '$app/state';

	let { children } = $props();

	// Check if this is a full-screen tool like WBS
	const isFullScreen = $derived(page.url.pathname.includes('/wbs'));
	let isSidebarCollapsed = $state(false);
</script>

<RouteGuard>
	<div class="flex h-screen overflow-hidden bg-mesh font-sans text-slate-900">
		<Sidebar bind:collapsed={isSidebarCollapsed} />
		
		<main class="flex-1 flex flex-col overflow-hidden">
			{#if !isFullScreen}
				<TopHeader />
			{/if}
			
			<div class="flex-1 overflow-y-auto {isFullScreen ? 'p-0' : 'p-10'} bg-mesh">
				<div class={isFullScreen ? 'h-full w-full' : 'max-w-[1440px] mx-auto'}>
					{@render children()}
				</div>
			</div>
		</main>
	</div>
</RouteGuard>

<style>
	:global(.bg-mesh) {
		background-color: #f6f6f8;
		background-image: 
			radial-gradient(at 0% 0%, rgba(19, 55, 236, 0.05) 0px, transparent 50%),
			radial-gradient(at 100% 0%, rgba(125, 211, 252, 0.1) 0px, transparent 50%);
	}
</style>
