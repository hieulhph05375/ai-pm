<script lang="ts">
	import type { Snippet } from "svelte";

	interface Props {
		type?: "button" | "submit" | "reset";
		variant?: "primary" | "glass" | "ghost" | "outline" | "danger";
		size?: "sm" | "md" | "lg";
		class?: string;
		disabled?: boolean;
		form?: string;
		onclick?: () => void;
		icon?: string;
		loading?: boolean;
		children?: Snippet;
	}

	let {
		type = "button",
		variant = "primary",
		size = "md",
		class: className = "",
		disabled = false,
		form,
		onclick,
		icon,
		loading = false,
		children,
	}: Props = $props();

	const variants = {
		primary:
			"gradient-bg text-white shadow-lg shadow-primary/20 hover:opacity-90 active:scale-[0.98]",
		glass: "glass-card text-slate-700 hover:bg-white/50 active:scale-[0.98]",
		ghost: "bg-transparent text-slate-600 hover:bg-slate-100 hover:text-primary active:scale-[0.98]",
		outline:
			"bg-white border border-slate-200 text-slate-700 hover:bg-slate-50 active:scale-[0.98]",
		danger: "bg-red-600 text-white hover:bg-red-700 active:scale-[0.98]",
	};

	const sizes = {
		sm: "px-3 py-1.5 text-xs gap-1.5 rounded-lg",
		md: "px-5 py-2.5 text-sm gap-2 rounded-xl",
		lg: "px-8 py-4 text-base gap-3 rounded-2xl font-bold",
	};

	const baseClasses =
		"inline-flex items-center justify-center font-semibold transition-all duration-200 disabled:opacity-50 disabled:pointer-events-none appearance-none";
</script>

<button
	{type}
	class="{baseClasses} {variants[variant]} {sizes[size]} {className}"
	{disabled}
	{form}
	{onclick}
>
	{#if loading}
		<div
			class="size-4 border-2 border-white/30 border-t-white rounded-full animate-spin mr-2"
		></div>
	{:else if icon}
		<span
			class="material-symbols-outlined {size === 'lg'
				? 'text-[18px]'
				: 'text-[18px]'}">{icon}</span
		>
	{/if}

	{#if children}
		{@render children()}
	{/if}
</button>
