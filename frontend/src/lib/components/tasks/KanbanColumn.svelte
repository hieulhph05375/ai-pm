<script lang="ts">
    import type { Task } from "$lib/services/tasks";
    import TaskCard from "./TaskCard.svelte";

    interface Props {
        status: { value: string; label: string; color: string };
        statusId?: number;
        tasks: Task[];
        ondrop?: (e: DragEvent, statusId: number) => void;
        ondragover?: (e: DragEvent) => void;
        ontaskclick?: (task: Task) => void;
        ondragstart?: (e: DragEvent, task: Task) => void;
        ondelete?: (task: Task) => void;
        onlogtime?: (task: Task) => void;
    }

    let {
        status,
        statusId,
        tasks,
        ondrop,
        ondragover,
        ontaskclick,
        ondragstart,
        ondelete,
        onlogtime,
    }: Props = $props();

    let isDragOver = $state(false);

    const columnColors: Record<string, string> = {
        slate: "border-t-slate-300 bg-slate-50/60",
        amber: "border-t-amber-400 bg-amber-50/40",
        emerald: "border-t-emerald-400 bg-emerald-50/40",
    };

    const dotColors: Record<string, string> = {
        slate: "bg-slate-400",
        amber: "bg-amber-400",
        emerald: "bg-emerald-400",
    };
</script>

<div
    class="flex flex-col flex-1 min-w-[320px] rounded-2xl border-t-4 {columnColors[
        status.color
    ] ??
        columnColors.slate} border border-slate-200/60 transition-all duration-200 {isDragOver
        ? 'scale-[1.01] shadow-lg border-primary/30'
        : ''}"
    ondragover={(e) => {
        e.preventDefault();
        isDragOver = true;
        ondragover?.(e);
    }}
    ondragleave={() => (isDragOver = false)}
    ondrop={(e) => {
        e.preventDefault();
        isDragOver = false;
        ondrop?.(e, statusId || Number(status.value));
    }}
    role="region"
    aria-label={status.label}
>
    <!-- Column Header -->
    <div class="flex items-center justify-between px-4 pt-4 pb-3">
        <div class="flex items-center gap-2">
            <div
                class="size-2.5 rounded-full {dotColors[status.color] ??
                    dotColors.slate}"
            ></div>
            <h3
                class="font-bold text-slate-700 text-sm uppercase tracking-wider"
            >
                {status.label}
            </h3>
        </div>
        <span
            class="text-xs font-bold text-slate-400 bg-white px-2 py-0.5 rounded-full border border-slate-200 shadow-sm"
            >{tasks.length}</span
        >
    </div>

    <!-- Task Cards Stack -->
    <div class="flex flex-col gap-3 px-3 pb-4 flex-1 min-h-20">
        {#if tasks.length === 0}
            <div
                class="flex flex-col items-center justify-center py-8 text-center pointer-events-none opacity-50"
            >
                <span
                    class="material-symbols-outlined text-3xl text-slate-300 mb-1"
                    >inbox</span
                >
                <p class="text-xs text-slate-400">No tasks</p>
            </div>
        {/if}

        {#each tasks as task (task.id)}
            <TaskCard
                {task}
                ondragstart={(e) => ondragstart?.(e, task)}
                onclick={() => ontaskclick?.(task)}
                ondelete={() => ondelete?.(task)}
                onlogtime={() => onlogtime?.(task)}
            />
        {/each}

        <!-- Drop area hint when dragging over -->
        {#if isDragOver}
            <div
                class="h-16 rounded-xl border-2 border-dashed border-primary/30 bg-primary/5 flex items-center justify-center"
            >
                <span class="text-xs text-primary/60 font-medium"
                    >Drop here</span
                >
            </div>
        {/if}
    </div>
</div>
