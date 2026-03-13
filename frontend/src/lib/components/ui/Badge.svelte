<script lang="ts">
	import type { Snippet } from "svelte";

	interface Props {
		variant?:
			| "primary"
			| "indigo"
			| "emerald"
			| "amber"
			| "rose"
			| "pink"
			| "slate";
		class?: string;
		color?: string;
		style?: string;
		children?: Snippet;
	}

	let {
		variant = "primary",
		class: className = "",
		color = "",
		style = "",
		children,
	}: Props = $props();

	function getCustomStyle(c: string) {
		if (!c || !c.startsWith("#")) return "";
		// If it's a hex color, we apply it to background (with alpha) and text
		return `background-color: ${c}15; color: ${c}; border-color: ${c}30; ${style}`;
	}

	const variants = {
		primary: "bg-primary/5 text-primary border-primary/10",
		indigo: "bg-indigo-50 text-indigo-600 border-indigo-100",
		emerald: "bg-emerald-50 text-emerald-600 border-emerald-100",
		amber: "bg-amber-50 text-amber-600 border-amber-100",
		rose: "bg-rose-50 text-rose-600 border-rose-100",
		pink: "bg-pink-50 text-pink-600 border-pink-100",
		slate: "bg-slate-50 text-slate-600 border-slate-200",
	};

	const baseClasses =
		"inline-flex w-fit items-center px-2.5 py-1 rounded-lg text-[10px] font-bold uppercase tracking-wider border transition-colors whitespace-nowrap";
</script>

<span
	class="{baseClasses} {color.startsWith('#')
		? ''
		: variants[variant]} {className}"
	style={getCustomStyle(color)}
>
	{#if children}
		{@render children()}
	{/if}
</span>
