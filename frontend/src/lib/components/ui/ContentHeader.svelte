<script lang="ts">
	import type { Snippet } from "svelte";

	interface Props {
		title?: string;
		subtitle?: string | Snippet;
		children?: Snippet;
		titleSnippet?: Snippet;
	}

	let { title, subtitle = "", children, titleSnippet }: Props = $props();
</script>

<div class="flex items-center justify-between mb-8">
	<div class="flex-1 min-w-0">
		{#if titleSnippet}
			{@render titleSnippet()}
		{:else if title}
			<h2
				class="font-display text-2xl font-bold text-slate-900 tracking-tight mb-1 truncate"
			>
				{title}
			</h2>
		{/if}

		{#if typeof subtitle === "function"}
			{@render subtitle()}
		{:else if subtitle}
			<p class="text-slate-500 font-medium truncate">{subtitle}</p>
		{/if}
	</div>

	<div class="flex gap-3 ml-4 shrink-0">
		{#if children}
			{@render children()}
		{/if}
	</div>
</div>
