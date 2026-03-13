<script lang="ts">
	import { fade, scale } from "svelte/transition";
	import Button from "./Button.svelte";

	let {
		show = $bindable(false),
		title = "Confirm",
		message = "Are you sure you want to perform this action?",
		confirmText = "Confirm",
		cancelText = "Cancel",
		variant = "danger",
		onConfirm,
		onCancel,
	}: {
		show: boolean;
		title?: string;
		message?: string;
		confirmText?: string;
		cancelText?: string;
		variant?: "primary" | "danger";
		onConfirm: () => void;
		onCancel: () => void;
	} = $props();

	function handleKeydown(e: KeyboardEvent) {
		if (!show) return;
		if (e.key === "Escape") {
			onCancel();
		} else if (e.key === "Enter") {
			onConfirm();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

{#if show}
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div
		class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm"
		transition:fade={{ duration: 250, easing: (t) => t * (2 - t) }}
		onclick={onCancel}
	>
		<!-- Dialog Content -->
		<div
			class="bg-white dark:bg-slate-900 w-full max-w-sm rounded-[1.5rem] shadow-2xl border border-white/20 dark:border-slate-700/30 overflow-hidden flex flex-col p-6 text-center"
			transition:scale={{
				duration: 250,
				start: 0.95,
				easing: (t) => t * (2 - t),
			}}
			onclick={(e) => e.stopPropagation()}
		>
			<div
				class="mx-auto flex items-center justify-center size-14 rounded-full mb-4 {variant ===
				'danger'
					? 'bg-red-100 text-red-600'
					: 'bg-primary/10 text-primary'}"
			>
				<span class="material-symbols-outlined text-3xl">
					{variant === "danger" ? "warning" : "info"}
				</span>
			</div>

			<h3
				class="text-xl font-outfit font-bold text-slate-900 dark:text-white mb-2"
			>
				{title}
			</h3>

			<p class="text-slate-500 text-sm mb-8">
				{message}
			</p>

			<div class="flex gap-3 w-full">
				<Button variant="outline" class="flex-1" onclick={onCancel}
					>{cancelText}</Button
				>
				<Button {variant} class="flex-1" onclick={onConfirm}
					>{confirmText}</Button
				>
			</div>
		</div>
	</div>
{/if}
