<script lang="ts" generics="T">
	import type { Snippet } from "svelte";
	import Card from "./Card.svelte";
	import Pagination from "./Pagination.svelte";

	interface Props {
		items: T[];
		columns: {
			key: keyof T | string;
			label: string;
			class?: string;
			align?: "left" | "right" | "center";
		}[];
		filterBar?: Snippet;
		headerCell?: Snippet<[{ column: Props["columns"][0] }]>;
		rowCell?: Snippet<[{ item: T; column: Props["columns"][0] }]>;
		loading?: boolean;

		// Pagination
		total?: number;
		page?: number;
		limit?: number;
		onPageChange?: (page: number) => void;
		onLimitChange?: (limit: number) => void;
		emptyState?: Snippet;
	}

	let {
		items = [],
		columns,
		filterBar,
		headerCell,
		rowCell,
		emptyState,
		loading = false,
		total = 0,
		page = 1,
		limit = 10,
		onPageChange,
		onLimitChange,
	}: Props = $props();

	const startEntry = $derived((page - 1) * limit + 1);
	const endEntry = $derived(Math.min(page * limit, total));
	const totalPages = $derived(Math.ceil(total / limit));
</script>

<Card padding="none" class="flex flex-col h-full overflow-hidden">
	<!-- Advanced Filter Bar -->
	{#if filterBar}
		<div class="p-6 border-b border-slate-100 bg-white/50 shrink-0">
			{@render filterBar()}
		</div>
	{/if}

	<!-- Table Container - Flex-1 and overflow-auto for vertical scrolling -->
	<div class="flex-1 overflow-auto relative min-h-0">
		{#if loading}
			<div
				class="absolute inset-0 bg-white/50 backdrop-blur-[1px] z-10 flex items-center justify-center"
			>
				<div
					class="size-8 border-4 border-primary/20 border-t-primary rounded-full animate-spin"
				></div>
			</div>
		{/if}

		<table class="w-full text-left border-collapse min-w-full">
			<thead class="sticky top-0 z-20 bg-slate-50 shadow-sm">
				<tr>
					{#each columns as column}
						<th
							class="px-6 py-5 text-xs font-bold uppercase tracking-wider text-slate-400 {column.align ===
							'right'
								? 'text-right'
								: ''} {column.class || ''}"
						>
							{#if headerCell}
								{@render headerCell({ column })}
							{:else}
								{column.label}
							{/if}
						</th>
					{/each}
				</tr>
			</thead>
			<tbody class="divide-y divide-slate-100 bg-white">
				{#each items as item}
					<tr class="hover:bg-slate-50/50 transition-colors group">
						{#each columns as column}
							<td
								class="px-6 py-5 text-sm text-slate-700 {column.align ===
								'right'
									? 'text-right'
									: ''} {column.class || ''}"
							>
								{#if rowCell}
									{@render rowCell({ item, column })}
								{:else}
									{String(item[column.key as keyof T] || "")}
								{/if}
							</td>
						{/each}
					</tr>
				{:else}
					<tr>
						<td colspan={columns.length} class="p-0 border-b-0">
							{#if emptyState}
								{@render emptyState()}
							{:else}
								<div
									class="px-6 py-20 text-center text-slate-400"
								>
									<span
										class="material-symbols-outlined text-4xl mb-2 opacity-20"
										>inventory_2</span
									>
									<p class="text-sm font-medium">
										No data to display
									</p>
								</div>
							{/if}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>

	<!-- Pagination Footer - Shrink-0 to prevent compression -->
	{#if total > 0}
		<div class="shrink-0">
			<Pagination
				{total}
				{page}
				{limit}
				onPageChange={(p) => onPageChange?.(p)}
				onLimitChange={(l) => onLimitChange?.(l)}
			/>
		</div>
	{/if}
</Card>
