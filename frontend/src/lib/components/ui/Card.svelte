<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		class?: string;
		padding?: 'none' | 'sm' | 'md' | 'lg';
		hover?: boolean;
		children?: Snippet;
	}

	let {
		class: className = '',
		padding = 'md',
		hover = false,
		children
	}: Props = $props();

	const paddings = {
		none: 'p-0',
		sm: 'p-4',
		md: 'p-6',
		lg: 'p-10'
	};
</script>

<div class="glass-card rounded-[2rem] overflow-hidden transition-all duration-300 {hover ? 'hover:shadow-xl hover:shadow-primary/5 hover:-translate-y-1' : ''} {paddings[padding]} {className}">
	{#if children}
		{@render children()}
	{/if}
</div>
