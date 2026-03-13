<script lang="ts">
    import type { Task } from "$lib/services/tasks";
    import { TASK_STATUSES, TASK_PRIORITIES } from "$lib/services/tasks";
    import Badge from "$lib/components/ui/Badge.svelte";

    interface Props {
        task: Task;
        ondragstart?: (e: DragEvent) => void;
        onclick?: () => void;
        ondelete?: () => void;
        onlogtime?: () => void;
    }

    let { task, ondragstart, onclick, ondelete, onlogtime }: Props = $props();

    function getPriorityInfo(task: Task) {
        if (task.priority_cat) {
            return {
                label: task.priority_cat.name,
                color: task.priority_cat.color || "slate",
            };
        }
        return (
            TASK_PRIORITIES.find((p) => p.value === task.priority) ?? {
                label: task.priority,
                color: "slate",
            }
        );
    }

    function formatDate(dateStr: string | null) {
        if (!dateStr) return null;
        return new Date(dateStr).toLocaleDateString("en-US", {
            day: "numeric",
            month: "short",
        });
    }

    const isOverdue = $derived(() => {
        if (!task.due_date || task.status === "DONE") return false;
        return new Date(task.due_date) < new Date();
    });
</script>

<div
    class="relative group bg-white rounded-xl p-4 shadow-sm border border-slate-100 hover:shadow-md hover:border-primary/20 transition-all duration-200 cursor-pointer select-none"
    draggable="true"
    {ondragstart}
    {onclick}
    role="button"
    tabindex="0"
    onkeydown={(e) => e.key === "Enter" && onclick?.()}
>
    <!-- Actions -->
    <div
        class="absolute top-2 right-2 flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-all z-10 bg-white shadow-sm rounded-md border border-slate-100 p-0.5"
    >
        {#if onlogtime && task.status !== "DONE"}
            <button
                type="button"
                class="p-1 text-slate-400 hover:text-primary hover:bg-primary/10 rounded transition-colors"
                onclick={(e) => {
                    e.stopPropagation();
                    onlogtime?.();
                }}
                title="Log Time"
            >
                <span class="material-symbols-outlined text-[16px]"
                    >schedule</span
                >
            </button>
        {/if}
        {#if ondelete}
            <button
                type="button"
                class="p-1 text-slate-400 hover:text-rose-500 hover:bg-rose-50 rounded transition-colors"
                onclick={(e) => {
                    e.stopPropagation();
                    ondelete?.();
                }}
                title="Delete task"
            >
                <span class="material-symbols-outlined text-[16px]">delete</span
                >
            </button>
        {/if}
    </div>

    <!-- Labels -->
    {#if task.labels && task.labels.length > 0}
        <div class="flex flex-wrap gap-1 mb-2">
            {#each task.labels.slice(0, 3) as label}
                <span
                    class="text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-full bg-primary/10 text-primary"
                    >{label}</span
                >
            {/each}
        </div>
    {/if}

    <!-- Title -->
    <p
        class="font-semibold text-slate-900 text-sm leading-snug mb-2 group-hover:text-primary transition-colors line-clamp-2"
    >
        {task.title}
    </p>

    <!-- Description -->
    {#if task.description}
        <p class="text-xs text-slate-400 leading-relaxed mb-3 line-clamp-2">
            {task.description}
        </p>
    {/if}

    <!-- Footer -->
    <div class="flex items-center justify-between gap-2 mt-1">
        {#if true}
            {@const pr = getPriorityInfo(task)}
            <Badge color={pr.color as any}>
                {pr.label}
            </Badge>
        {/if}

        {#if formatDate(task.due_date)}
            <span
                class="flex items-center gap-1 text-[11px] font-medium {isOverdue()
                    ? 'text-rose-500'
                    : 'text-slate-400'}"
            >
                <span class="material-symbols-outlined text-[13px]"
                    >{isOverdue() ? "warning" : "calendar_today"}</span
                >
                {formatDate(task.due_date)}
            </span>
        {/if}
    </div>
</div>
