<script lang="ts">
	interface Props {
		type?: string;
		placeholder?: string;
		value?: string;
		label?: string;
		id?: string;
		icon?: string;
		rightIcon?: string;
		class?: string;
		wrapperClass?: string;
		disabled?: boolean;
		required?: boolean;
		error?: string;
		oninput?: (e: Event) => void;
	}

	let {
		type = "text",
		placeholder = "",
		value = $bindable(),
		label = "",
		id,
		icon,
		rightIcon,
		class: className = "",
		wrapperClass = "",
		disabled = false,
		required = false,
		error = "",
		oninput,
	}: Props = $props();

	// Generate random ID if none provided
	const inputId = $derived(
		id || `input-${Math.random().toString(36).substring(2, 9)}`,
	);
</script>

<div class="space-y-1.5 {wrapperClass}">
	{#if label}
		<label
			for={inputId}
			class="text-xs font-bold text-slate-500 uppercase tracking-widest ml-1"
			>{label}</label
		>
	{/if}
	<div class="relative group">
		{#if icon}
			<span
				class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[18px] group-focus-within:text-primary transition-colors"
			>
				{icon}
			</span>
		{/if}

		<input
			id={inputId}
			{type}
			{placeholder}
			bind:value
			{disabled}
			{required}
			{oninput}
			class="w-full {icon ? 'pl-10' : 'px-4'} {rightIcon
				? 'pr-10'
				: 'px-4'} py-3 bg-slate-50 dark:bg-slate-900/50 border {error
				? 'border-rose-500 focus:ring-rose-500/10'
				: 'border-slate-200 dark:border-slate-700 focus:ring-primary/10'} rounded-xl text-sm focus:ring-4 focus:border-primary outline-none transition-all text-slate-900 dark:text-white placeholder:text-slate-400/70 disabled:opacity-50 {className}"
		/>

		{#if rightIcon}
			<button
				type="button"
				class="absolute right-3 top-1/2 -translate-y-1/2 {error
					? 'text-rose-400'
					: 'text-slate-400'} hover:text-slate-600 transition-colors"
			>
				<span class="material-symbols-outlined text-[18px]"
					>{rightIcon}</span
				>
			</button>
		{/if}
	</div>

	{#if error}
		<p
			class="text-[11px] font-bold text-rose-500 ml-1 mt-1 flex items-center gap-1 animate-in fade-in slide-in-from-top-1 duration-200"
		>
			<span class="material-symbols-outlined text-[14px]">error</span>
			{error}
		</p>
	{/if}
</div>
