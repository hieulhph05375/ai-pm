<script lang="ts">
	import { fade, scale } from "svelte/transition";

	let {
		show = $bindable(false),
		title = "",
		maxWidth = "max-w-2xl",
		onClose,
		children,
		footer,
	}: {
		show: boolean;
		title: string;
		maxWidth?: string;
		onClose?: () => void;
		children?: import("svelte").Snippet;
		footer?: import("svelte").Snippet;
	} = $props();

	function close() {
		show = false;
		if (onClose) onClose();
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === "Escape" && show) {
			close();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

{#if show}
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div
		class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm"
		transition:fade={{ duration: 300, easing: (t) => t * (2 - t) }}
		onclick={close}
	>
		<!-- Modal Content -->
		<div
			class="bg-white/95 dark:bg-slate-900/95 backdrop-blur-xl w-full {maxWidth} rounded-2xl shadow-2xl border border-white/20 dark:border-slate-700/30 overflow-hidden flex flex-col max-h-[90vh]"
			transition:scale={{
				duration: 300,
				start: 0.95,
				easing: (t) => t * (2 - t),
			}}
			onclick={(e) => e.stopPropagation()}
		>
			<!-- Header -->
			<div
				class="px-8 pt-8 pb-4 flex justify-between items-center shrink-0"
			>
				<h2
					class="text-3xl font-outfit font-bold text-slate-900 dark:text-white"
				>
					{title}
				</h2>
				<button
					type="button"
					class="text-slate-400 hover:text-slate-600 transition-colors rounded-lg p-1 hover:bg-slate-100"
					onclick={close}
				>
					<span class="material-symbols-outlined">close</span>
				</button>
			</div>

			<!-- Body -->
			<div
				class="px-8 pb-8 overflow-y-auto custom-scrollbar flex-1 flex flex-col min-h-0"
			>
				{#if children}
					{@render children()}
				{/if}
			</div>

			<!-- Footer -->
			{#if footer}
				<div
					class="px-8 py-6 border-t border-slate-100 dark:border-slate-800 shrink-0 flex gap-4 mt-auto bg-white/50 dark:bg-slate-900/50 backdrop-blur-sm"
				>
					{@render footer()}
				</div>
			{/if}
		</div>
	</div>
{/if}

<style>
	.custom-scrollbar::-webkit-scrollbar {
		width: 6px;
	}
	.custom-scrollbar::-webkit-scrollbar-track {
		background: transparent;
	}
	.custom-scrollbar::-webkit-scrollbar-thumb {
		background-color: rgba(203, 213, 225, 0.5); /* slate-300 */
		border-radius: 20px;
	}
	.custom-scrollbar:hover::-webkit-scrollbar-thumb {
		background-color: rgba(148, 163, 184, 0.8); /* slate-400 */
	}
</style>
