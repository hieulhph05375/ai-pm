<script lang="ts">
	import {
		wbsService,
		type WBSNode,
		type WBSBaselineNode,
	} from "$lib/services/wbs";
	import type { User } from "$lib/services/users";
	import EmptyState from "$lib/components/ui/EmptyState.svelte";
	import { hasPermission } from "$lib/utils/permission";
	import { authStore } from "$lib/services/auth";
	import Badge from "../ui/Badge.svelte";

	interface Props {
		nodes?: WBSNode[];
		allNodes?: WBSNode[];
		collapsedPaths?: Set<string>;
		users?: User[];
		baselineNodes?: WBSBaselineNode[];
		hoveredNodeId?: number | null;
		onToggleCollapse?: (path: string) => void;
		onToggleAll?: (expand: boolean) => void;
		onEdit?: (node: WBSNode) => void;
		onAddSubtask?: (parent: WBSNode) => void;
		onDelete?: (node: WBSNode) => void;
		onStatusChange?: (node: WBSNode, progress: number) => void;
		onPicChange?: (node: WBSNode, userId: number | null) => void;
		onHover?: (id: number | null) => void;
		scrollTop?: number;
		searchText?: string;
	}

	let {
		nodes = [],
		allNodes = [],
		collapsedPaths = new Set(),
		users = [],
		baselineNodes = [],
		hoveredNodeId = null,
		scrollTop = $bindable(0),
		searchText = "",
		onToggleCollapse,
		onToggleAll,
		onEdit,
		onAddSubtask,
		onDelete,
		onStatusChange,
		onPicChange,
		onHover,
	}: Props = $props();

	function highlightText(text: string, query: string): string {
		if (!query) return text;
		const regex = new RegExp(`(${query})`, "gi");
		return text.replace(
			regex,
			'<span class="bg-amber-200 text-slate-900 rounded-sm px-0.5">$1</span>',
		);
	}

	function getDepth(path: string): number {
		return path.split(".").length - 1;
	}

	function getStatusLabel(node: WBSNode): string {
		if (node.progress === 100) return "Completed";
		if (node.progress > 0) return "In Progress";
		return "Scheduled";
	}

	let detailsEl: HTMLElement | undefined = $state();

	function handleScroll(e: Event) {
		scrollTop = (e.currentTarget as HTMLElement).scrollTop;
	}

	$effect(() => {
		if (detailsEl && detailsEl.scrollTop !== scrollTop) {
			detailsEl.scrollTop = scrollTop;
		}
	});

	function calculateDays(start: string | null, end: string | null): number {
		if (!start || !end) return 0;
		const s = new Date(start);
		const e = new Date(end);
		const diff = e.getTime() - s.getTime();
		return Math.ceil(diff / (1000 * 60 * 60 * 24)) + 1;
	}

	function calculateRemainingDays(
		end: string | null,
		progress: number,
	): number {
		if (!end || progress >= 100) return 0;
		const e = new Date(end);
		const t = new Date();
		t.setHours(0, 0, 0, 0);
		e.setHours(0, 0, 0, 0);
		if (t > e) return 0;
		const diff = e.getTime() - t.getTime();
		return Math.ceil(diff / (1000 * 60 * 60 * 24));
	}
	// Virtual Scrolling constants
	const ROW_HEIGHT = 64;
	const BUFFER = 5;

	let viewportHeight = $state(0);

	let visibleRange = $derived.by(() => {
		const startIdx = Math.max(
			0,
			Math.floor(scrollTop / ROW_HEIGHT) - BUFFER,
		);
		const endIdx = Math.min(
			nodes.length,
			Math.ceil((scrollTop + viewportHeight) / ROW_HEIGHT) + BUFFER,
		);
		return { start: startIdx, end: endIdx };
	});

	let visibleNodes = $derived(
		nodes.slice(visibleRange.start, visibleRange.end),
	);
	let totalHeight = $derived(nodes.length * ROW_HEIGHT);
	let offsetY = $derived(visibleRange.start * ROW_HEIGHT);
</script>

<div class="w-[500px] flex flex-col sidebar-glass z-40">
	<div
		class="h-12 border-b border-slate-200/50 flex items-center px-6 justify-between flex-shrink-0"
	>
		<span class="text-xs font-bold uppercase tracking-widest text-slate-400"
			>Task Details</span
		>

		<button
			class="p-1.5 rounded-lg hover:bg-slate-100 text-slate-400 hover:text-primary transition-all flex items-center justify-center"
			onclick={() => onToggleAll?.(collapsedPaths.size > 0)}
			title={collapsedPaths.size > 0 ? "Expand all" : "Collapse all"}
		>
			<span class="material-symbols-outlined text-[18px]">
				{collapsedPaths.size > 0 ? "unfold_more" : "unfold_less"}
			</span>
		</button>
	</div>
	<div
		class="flex-1 overflow-y-auto no-scrollbar relative"
		bind:this={detailsEl}
		bind:clientHeight={viewportHeight}
		onscroll={handleScroll}
	>
		<!-- Placeholder to enforce total height -->
		<div
			style="height: {totalHeight}px; width: 100%; pointer-events: none; position: absolute; top: 0; left: 0;"
		></div>

		<!-- Visible Items Container -->
		<div style="transform: translateY({offsetY}px); width: 100%;">
			{#each visibleNodes as node, idx (node.id)}
				{@const globalIdx = visibleRange.start + idx}
				{@const depth = getDepth(node.path)}
				{@const hasChildren = node.has_children}
				{@const isCollapsed = collapsedPaths.has(node.path)}
				{@const totalDays = calculateDays(
					node.planned_start_date ?? null,
					node.planned_end_date ?? null,
				)}
				{@const remaining = calculateRemainingDays(
					node.planned_end_date ?? null,
					node.progress,
				)}
				{@const baselineNode = baselineNodes.find(
					(bn) => bn.node_id === node.id,
				)}
				{@const slippage =
					baselineNode && node.planned_end_date
						? wbsService.calculateVarianceDays(
								node.planned_end_date,
								baselineNode.planned_end_date,
							)
						: null}
				<div
					role="row"
					tabindex="0"
					class="h-[64px] min-h-[64px] max-h-[64px] p-3 px-6 border-b border-slate-100 transition-all duration-200 group relative flex flex-col justify-center flex-none"
					class:bg-primary-5={hoveredNodeId === node.id}
					class:bg-slate-50={globalIdx % 2 !== 0 &&
						hoveredNodeId !== node.id}
					style="padding-left: {32 +
						depth * 20}px; background-color: {hoveredNodeId ===
					node.id
						? 'rgba(19, 55, 236, 0.08)'
						: ''}"
					onmouseenter={() => onHover?.(node.id)}
					onmouseleave={() => onHover?.(null)}
				>
					{#if hasChildren}
						<button
							class="absolute top-1/2 -translate-y-1/2 size-5 flex items-center justify-center rounded hover:bg-slate-200/50 transition-all text-slate-400 hover:text-primary z-20"
							style="left: {12 + depth * 20}px"
							onclick={(e) => {
								e.stopPropagation();
								onToggleCollapse?.(node.path);
							}}
						>
							<span
								class="material-symbols-outlined text-[14px] transition-transform duration-200 {isCollapsed
									? '-rotate-90'
									: ''}"
								style="font-size:14px;"
							>
								expand_more
							</span>
						</button>
					{/if}
					<div class="flex items-start justify-between mb-1">
						<div class="flex flex-col gap-0.5 min-w-0">
							<div class="flex items-center gap-2">
								<span
									class="text-[10px] font-bold text-primary/60 font-mono"
									>{node.path}</span
								>
								{#if node.type_cat}
									<Badge
										color={node.type_cat.color || ""}
										class="!px-1.5 !py-0.5"
									>
										{#if node.type_cat.icon}
											<span
												class="material-symbols-outlined text-[10px] mr-1"
												>{node.type_cat.icon}</span
											>
										{/if}
										{node.type_cat.name}
									</Badge>
								{:else}
									<Badge
										variant={node.type === "Phase"
											? "slate"
											: node.type === "Milestone"
												? "amber"
												: "indigo"}
										class="!px-1.5 !py-0.5"
									>
										{node.type}
									</Badge>
								{/if}
								<a
									href="/projects/{node.project_id}/wbs/{node.id}"
									class="font-semibold text-sm text-slate-900 hover:text-primary active:text-primary-dark transition-colors truncate decoration-primary/30 hover:underline underline-offset-4"
									onclick={(e) => e.stopPropagation()}
								>
									{@html highlightText(
										node.title,
										searchText,
									)}
								</a>
							</div>
						</div>
						<div class="flex items-center gap-2 mt-1 flex-shrink-0">
							{#if totalDays > 0}
								<div
									class="flex items-center gap-0.5 text-[9px] font-bold text-slate-400 bg-slate-50 px-1.5 py-0.5 rounded border border-slate-100"
									title="Total duration"
								>
									<span
										style="font-size: 12px;"
										class="material-symbols-outlined text-[10px]"
										>calendar_today</span
									>
									{totalDays}d
								</div>
							{/if}
							{#if slippage !== null}
								<div
									class="flex items-center gap-0.5 text-[9px] font-black {slippage >
									0
										? 'text-rose-600 bg-rose-50 border-rose-100'
										: slippage < 0
											? 'text-emerald-600 bg-emerald-50 border-emerald-100'
											: 'text-slate-400 bg-slate-50 border-slate-100'} px-1.5 py-0.5 rounded border"
									title="Variance from baseline"
								>
									<span
										style="font-size: 11px;"
										class="material-symbols-outlined text-[10px]"
										>history</span
									>
									{slippage > 0 ? "+" : ""}{slippage}d
								</div>
							{/if}
							<div
								class="flex items-center ml-2 border-l border-slate-200 pl-2 gap-2"
							>
								<div
									class="flex items-center gap-0.5 text-[9px] font-bold text-indigo-600 bg-indigo-50 px-1.5 py-0.5 rounded border border-indigo-100"
									title="Estimated Effort (Hours)"
								>
									<span
										style="font-size: 11px;"
										class="material-symbols-outlined text-[10px]"
										>functions</span
									>
									{node.estimated_effort ?? 0}h
								</div>
								<div
									class="flex items-center gap-0.5 text-[9px] font-bold {(node.actual_effort ||
										0) > (node.estimated_effort || 0) &&
									(node.estimated_effort || 0) > 0
										? 'text-rose-600 bg-rose-50 border-rose-100'
										: 'text-emerald-600 bg-emerald-50 border-emerald-100'} px-1.5 py-0.5 rounded border"
									title="Actual Effort Logged (Hours)"
								>
									<span
										style="font-size: 11px;"
										class="material-symbols-outlined text-[10px]"
										>schedule</span
									>
									{node.actual_effort ?? 0}h
								</div>
							</div>
							<span
								class="size-1.5 rounded-full shadow-[0_0_6px] ml-2
					{node.progress === 100
									? 'bg-emerald-500 shadow-emerald-500/50'
									: node.progress > 0
										? 'bg-amber-500 shadow-amber-500/50'
										: 'bg-slate-300 shadow-slate-300/50'}"
							></span>
						</div>
					</div>

					<div class="flex items-center gap-2">
						<div
							class="size-[14px] rounded-full bg-slate-100 border border-slate-200 flex items-center justify-center overflow-hidden"
						>
							<span
								class="material-symbols-outlined text-[8px] text-slate-400"
								>person</span
							>
						</div>

						<select
							class="text-[10px] text-slate-500 bg-transparent border-none outline-none cursor-pointer hover:text-primary transition-colors pr-4 appearance-none"
							value={node.assigned_to != null
								? String(node.assigned_to)
								: ""}
							onclick={(e) => e.stopPropagation()}
							onchange={(e) => {
								const val = e.currentTarget.value
									? Number(e.currentTarget.value)
									: null;
								onPicChange?.(node, val);
							}}
							disabled={!hasPermission(
								$authStore.user,
								$authStore.token,
								"project:update",
							)}
						>
							<option value="">Unassigned</option>
							{#each users as user}
								<option value={String(user.id)}
									>{user.full_name}</option
								>
							{/each}
						</select>

						<!-- Quick Status Change Badge (Dropdown) -->
						<select
							class="ml-2 px-1.5 py-0.5 text-[8px] font-black rounded uppercase tracking-tighter border-none cursor-pointer outline-none transition-all
                        {node.progress === 100
								? 'bg-emerald-100 text-emerald-700 hover:bg-emerald-200'
								: node.progress > 0
									? 'bg-amber-100 text-amber-700 hover:bg-amber-200'
									: 'bg-slate-100 text-slate-500 hover:bg-slate-200'}"
							value={node.progress === 100
								? "100"
								: node.progress > 0
									? "50"
									: "0"}
							onclick={(e) => e.stopPropagation()}
							onchange={(e) => {
								const val = Number(e.currentTarget.value);
								onStatusChange?.(node, val);
							}}
							disabled={!hasPermission(
								$authStore.user,
								$authStore.token,
								"project:update",
							)}
						>
							<option value="0">On Hold</option>
							<option value="50">Running</option>
							<option value="100">Completed</option>
						</select>
					</div>

					<!-- Hover Actions -->
					<div
						class="absolute right-4 top-1/2 -translate-y-1/2 flex items-center gap-3 opacity-0 group-hover:opacity-100 transition-opacity bg-white/80 backdrop-blur-sm pl-2 pr-1 py-1 rounded-lg"
					>
						{#if hasPermission($authStore.user, $authStore.token, "project:update")}
							<button
								class="size-6 text-slate-400 hover:text-primary hover:scale-110 flex items-center justify-center transition-all"
								onclick={(e) => {
									e.stopPropagation();
									onEdit?.(node);
								}}
								title="Edit Task"
							>
								<span
									class="material-symbols-outlined text-[14px]"
									style="font-size:14px;">edit</span
								>
							</button>
							{#if depth < 4}
								<button
									class="size-6 text-slate-400 hover:text-primary hover:scale-110 flex items-center justify-center transition-all"
									onclick={(e) => {
										e.stopPropagation();
										onAddSubtask?.(node);
									}}
									title="Add Subtask"
								>
									<span
										class="material-symbols-outlined text-[14px]"
										style="font-size:14px;">add_circle</span
									>
								</button>
							{/if}
							<button
								class="size-[6px] text-slate-400 hover:text-rose-500 hover:scale-110 flex items-center justify-center transition-all"
								onclick={(e) => {
									e.stopPropagation();
									onDelete?.(node);
								}}
								title="Delete Task"
							>
								<span
									class="material-symbols-outlined text-[6px]"
									style="font-size:14px;">delete</span
								>
							</button>
						{/if}
					</div>
				</div>
			{/each}
		</div>

		{#if nodes.length === 0}
			<div class="px-6 py-12">
				<EmptyState
					icon="folder_open"
					title="No tasks found"
					message="No tasks match your current filters."
					actionIcon="add_circle"
					actionLabel="Add new task"
					onaction={() => onAddSubtask?.(allNodes[0])}
				/>
			</div>
		{/if}
	</div>
</div>

<style>
	.no-scrollbar::-webkit-scrollbar {
		display: none;
	}
	.no-scrollbar {
		-ms-overflow-style: none;
		scrollbar-width: none;
	}
</style>
