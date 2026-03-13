<script lang="ts">
	import { toast, type Toast } from "$lib/stores/toast";
	import { flip } from "svelte/animate";
	import { fade, fly } from "svelte/transition";

	const { toasts } = $derived({ toasts: $toast });
</script>

<div
	class="fixed bottom-6 right-6 z-[100] flex flex-col gap-3 pointer-events-none"
>
	{#each toasts as t (t.id)}
		<div
			animate:flip={{ duration: 300 }}
			transition:fly={{ y: 20, duration: 400 }}
			class="pointer-events-auto flex items-center gap-3 px-4 py-3 rounded-2xl shadow-xl border min-w-[300px] max-w-md bg-white/90 backdrop-blur-xl transition-all"
			class:border-emerald-100={t.variant === "success"}
			class:border-rose-100={t.variant === "error"}
			class:border-amber-100={t.variant === "warning"}
			class:border-blue-100={t.variant === "info"}
		>
			<div
				class="size-10 rounded-xl flex items-center justify-center shrink-0"
				class:bg-emerald-50={t.variant === "success"}
				class:text-emerald-500={t.variant === "success"}
				class:bg-rose-50={t.variant === "error"}
				class:text-rose-500={t.variant === "error"}
				class:bg-amber-50={t.variant === "warning"}
				class:text-amber-500={t.variant === "warning"}
				class:bg-blue-50={t.variant === "info"}
				class:text-blue-500={t.variant === "info"}
			>
				<span class="material-symbols-outlined shrink-0">
					{#if t.variant === "success"}check_circle{:else if t.variant === "error"}error{:else if t.variant === "warning"}warning{:else}info{/if}
				</span>
			</div>

			<div class="flex-1">
				<p class="text-sm font-semibold text-slate-900 leading-tight">
					{#if t.variant === "success"}Success{:else if t.variant === "error"}System
						Error{:else if t.variant === "warning"}Warning{:else}Info{/if}
				</p>
				<p class="text-[13px] text-slate-500 mt-0.5 line-clamp-2">
					{t.message}
				</p>
			</div>

			<button
				onclick={() => toast.dismiss(t.id)}
				class="p-1 hover:bg-slate-100 rounded-lg text-slate-400 transition-colors shrink-0"
			>
				<span class="material-symbols-outlined text-sm">close</span>
			</button>
		</div>
	{/each}
</div>
