<script lang="ts">
    import type { Task } from "$lib/services/tasks";
    import Badge from "$lib/components/ui/Badge.svelte";
    import GanttChart from "$lib/components/ui/GanttChart.svelte";
    import Pagination from "$lib/components/ui/Pagination.svelte";
    import type { GanttItem } from "$lib/types/gantt";

    interface Props {
        tasks: Task[];
        page?: number;
        total?: number;
        limit?: number;
        onPageChange?: (page: number) => void;
        onLimitChange?: (limit: number) => void;
        ontaskupdate?: (task: Task, updates: Partial<Task>) => void;
        ontaskclick?: (task: Task) => void;
        ondelete?: (task: Task) => void;
    }

    let {
        tasks,
        page = 1,
        total = 0,
        limit = 10,
        onPageChange,
        onLimitChange,
        ontaskupdate,
        ontaskclick,
        ondelete,
    }: Props = $props();

    let viewMode = $state<"Day" | "Week" | "Month" | "Quarter">("Day");
    let scrollTop = $state(0);
    let hoveredItemId = $state<string | number | null>(null);

    // Map tasks to GanttItems
    let ganttItems = $derived.by((): GanttItem[] => {
        return tasks
            .map((t) => ({
                id: t.id,
                title: t.title,
                startDate: t.start_date || t.due_date,
                dueDate: t.due_date || t.start_date,
                progress: t.progress || 0,
                type: "task" as const,
            }))
            .sort((a, b) => {
                const dateA = new Date(a.startDate || 0).getTime();
                const dateB = new Date(b.startDate || 0).getTime();
                return dateA - dateB;
            });
    });

    let sidebarEl: HTMLElement | undefined = $state();

    $effect(() => {
        if (sidebarEl && sidebarEl.scrollTop !== scrollTop) {
            sidebarEl.scrollTop = scrollTop;
        }
    });

    function handleSidebarScroll(e: Event) {
        const target = e.target as HTMLElement;
        if (target.scrollTop !== scrollTop) {
            scrollTop = target.scrollTop;
        }
    }

    let noDateTasksCount = $derived(
        tasks.filter((t) => !t.start_date && !t.due_date).length,
    );

    function handleProgressChange(item: GanttItem, progress: number) {
        const task = tasks.find((t) => t.id === item.id);
        if (task) {
            ontaskupdate?.(task, { progress });
        }
    }

    function handleDatesChange(item: GanttItem, newStart: Date, newEnd: Date) {
        const task = tasks.find((t) => t.id === item.id);
        if (task) {
            ontaskupdate?.(task, {
                start_date: newStart.toISOString(),
                due_date: newEnd.toISOString(),
            });
        }
    }
</script>

<div
    class="flex flex-col h-full bg-white rounded-2xl border border-slate-200 overflow-hidden relative"
>
    <!-- Top info bar -->
    <div
        class="flex items-center justify-between px-4 py-2 bg-slate-50 border-b border-slate-200 shrink-0"
    >
        <div class="flex items-center gap-4">
            <span class="text-xs font-bold text-slate-500 uppercase"
                >Timeline (BETA)</span
            >
            {#if noDateTasksCount > 0}
                <Badge variant="amber" class="text-[10px]"
                    >{noDateTasksCount} tasks without dates</Badge
                >
            {/if}
        </div>
        <div class="flex items-center gap-2">
            <select
                bind:value={viewMode}
                class="bg-white border border-slate-200 rounded-lg px-2 py-1 text-xs font-semibold focus:outline-none focus:ring-2 focus:ring-primary/20"
            >
                <option value="Day">Day</option>
                <option value="Week">Week</option>
                <option value="Month">Month</option>
                <option value="Quarter">Quarter</option>
            </select>
        </div>
    </div>

    <div
        class="flex-1 overflow-x-auto overflow-y-hidden custom-scrollbar relative flex min-h-0"
    >
        <!-- Sidebar: Task List -->
        <div
            class="w-64 shrink-0 bg-white border-r border-slate-200 z-20 flex flex-col sticky left-0 h-full"
        >
            <!-- Header -->
            <div
                class="h-12 border-b border-slate-200 bg-slate-50 flex items-center px-4 shrink-0"
            >
                <span
                    class="text-xs font-bold text-slate-500 uppercase tracking-widest"
                    >Task Name</span
                >
            </div>

            <!-- List -->
            <div
                class="flex-1 overflow-y-auto no-scrollbar pb-32"
                bind:this={sidebarEl}
                onscroll={handleSidebarScroll}
            >
                {#each ganttItems as item}
                    {@const task = tasks.find((t) => t.id === item.id)}
                    <!-- svelte-ignore a11y_click_events_have_key_events -->
                    <!-- svelte-ignore a11y_no_static_element_interactions -->
                    <div
                        class="h-[64px] border-b border-slate-100 flex items-center justify-between px-4 hover:bg-primary/5 cursor-pointer transition-colors text-sm font-medium text-slate-700 group relative"
                        style="background-color: {hoveredItemId === item.id
                            ? 'rgba(19, 55, 236, 0.08)'
                            : 'transparent'}"
                        onclick={() => task && ontaskclick?.(task)}
                        onmouseenter={() => (hoveredItemId = item.id)}
                        onmouseleave={() => (hoveredItemId = null)}
                    >
                        <span class="truncate pr-4">{item.title}</span>
                        {#if ondelete && task}
                            <button
                                type="button"
                                class="absolute right-2 p-1 text-slate-300 hover:text-rose-500 hover:bg-rose-50 rounded-md opacity-0 group-hover:opacity-100 transition-all shrink-0 bg-white"
                                onclick={(e) => {
                                    e.stopPropagation();
                                    ondelete?.(task);
                                }}
                                title="Delete Task"
                            >
                                <span
                                    class="material-symbols-outlined text-[16px]"
                                    >delete</span
                                >
                            </button>
                        {/if}
                    </div>
                {/each}
            </div>
        </div>

        <!-- Gantt Chart Component (Ported from WBS) -->
        <GanttChart
            items={ganttItems}
            bind:scrollTop
            bind:hoveredItemId
            {viewMode}
            onProgressChange={handleProgressChange}
            onDatesChange={handleDatesChange}
        />
    </div>

    <!-- Full-width Pagination Footer -->
    {#if total > 0 || onLimitChange}
        <div class="shrink-0 border-t border-slate-200">
            <Pagination
                {total}
                {page}
                {limit}
                onPageChange={(p) => onPageChange?.(p)}
                onLimitChange={(l) => onLimitChange?.(l)}
            />
        </div>
    {/if}
</div>

<style>
    /* Same custom scrollbar standard */
    .custom-scrollbar::-webkit-scrollbar {
        width: 6px;
        height: 6px;
    }
    .custom-scrollbar::-webkit-scrollbar-track {
        background: transparent;
    }
    .custom-scrollbar::-webkit-scrollbar-thumb {
        background-color: rgba(203, 213, 225, 0.5);
        border-radius: 10px;
    }
    .custom-scrollbar:hover::-webkit-scrollbar-thumb {
        background-color: rgba(148, 163, 184, 0.8);
    }
</style>
