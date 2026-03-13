<script lang="ts">
    import type { Task } from "$lib/services/tasks";
    import { TASK_STATUSES } from "$lib/services/tasks";
    import type { Category } from "$lib/services/categories";
    import KanbanColumn from "./KanbanColumn.svelte";

    interface Props {
        tasks: Task[];
        statusCategories?: Category[];
        onstatuschange?: (task: Task, newStatusId: number) => void;
        ontaskclick?: (task: Task) => void;
        ondelete?: (task: Task) => void;
        onlogtime?: (task: Task) => void;
    }

    let {
        tasks,
        statusCategories = [],
        onstatuschange,
        ontaskclick,
        ondelete,
        onlogtime,
    }: Props = $props();

    let draggedTask: Task | null = null;

    function handleDragStart(e: DragEvent, task: Task) {
        draggedTask = task;
        e.dataTransfer!.effectAllowed = "move";
        // Add dragging style via data attribute
        (e.currentTarget as HTMLElement)?.setAttribute("data-dragging", "true");
    }

    function handleDrop(e: DragEvent, targetStatusId: number) {
        if (!draggedTask || draggedTask.status_id === targetStatusId) {
            draggedTask = null;
            return;
        }
        onstatuschange?.(draggedTask, targetStatusId);
        draggedTask = null;
    }

    function getTasksByStatus(statusId: number) {
        return tasks.filter((t) => t.status_id === statusId);
    }
</script>

<div class="flex w-full gap-5 overflow-x-auto pb-4" style="min-height: 400px;">
    {#each statusCategories as status (status.id)}
        <KanbanColumn
            status={{
                value: status.id.toString(),
                label: status.name,
                color: status.color || "slate",
            }}
            statusId={status.id}
            tasks={getTasksByStatus(status.id)}
            ondrop={handleDrop}
            {ontaskclick}
            {ondelete}
            {onlogtime}
            ondragstart={handleDragStart}
        />
    {/each}
</div>
