<script lang="ts">
    import { onMount } from "svelte";
    import ContentHeader from "$lib/components/ui/ContentHeader.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import DataTable from "$lib/components/ui/DataTable.svelte";
    import Badge from "$lib/components/ui/Badge.svelte";
    import ConfirmDialog from "$lib/components/ui/ConfirmDialog.svelte";
    import TimeEntryModal from "$lib/components/timesheets/TimeEntryModal.svelte";
    import { timesheetService, type Timesheet } from "$lib/services/timesheets";
    import { toast } from "$lib/stores/toast";
    import EmptyState from "$lib/components/ui/EmptyState.svelte";

    let timesheets = $state<Timesheet[]>([]);
    let loading = $state(true);

    let currentPage = $state(1);
    let limit = $state(10);
    let totalItems = $state(0);

    // Modal State
    let showModal = $state(false);
    let selectedTimesheet = $state<Timesheet | null>(null);

    // Delete State
    let showConfirmDelete = $state(false);
    let timesheetToDelete = $state<Timesheet | null>(null);

    onMount(async () => {
        await loadTimesheets();
    });

    async function loadTimesheets() {
        loading = true;
        try {
            const res = await timesheetService.list({
                page: currentPage,
                limit,
            });
            timesheets = res.items || [];
            totalItems = res.total;
        } catch (e: any) {
            toast.error(e.message || "Failed to load timesheets");
        } finally {
            loading = false;
        }
    }

    async function handlePageChange(page: number) {
        currentPage = page;
        await loadTimesheets();
    }

    function handleCreateClick() {
        selectedTimesheet = null;
        showModal = true;
    }

    function handleEditClick(t: Timesheet) {
        selectedTimesheet = t;
        showModal = true;
    }

    function handleDeleteClick(t: Timesheet) {
        timesheetToDelete = t;
        showConfirmDelete = true;
    }

    async function confirmDelete() {
        if (!timesheetToDelete) return;

        try {
            await timesheetService.delete(timesheetToDelete.id);
            toast.success("Time entry deleted");
            await loadTimesheets();
        } catch (e: any) {
            toast.error(e.message || "Failed to delete time entry");
        } finally {
            showConfirmDelete = false;
            timesheetToDelete = null;
        }
    }

    function handleTimeSaved() {
        loadTimesheets();
    }

    function formatDate(val: string) {
        return new Date(val).toLocaleDateString("vi-VN", {
            weekday: "short",
            day: "2-digit",
            month: "2-digit",
            year: "numeric",
        });
    }

    // Columns
    const columns = [
        { key: "work_date", label: "DATE", class: "w-32" },
        { key: "context", label: "PROJECT / TASK" },
        { key: "hours", label: "HOURS", class: "w-24 text-center" },
        {
            key: "description",
            label: "DESCRIPTION",
            class: "hidden md:table-cell",
        },
        { key: "status", label: "STATUS", class: "w-28" },
        {
            key: "actions",
            label: "",
            class: "w-20 right",
            align: "right" as const,
        },
    ];

    function getStatusColor(status: string) {
        switch (status) {
            case "APPROVED":
                return "emerald";
            case "SUBMITTED":
                return "amber";
            case "REJECTED":
                return "rose";
            default:
                return "slate";
        }
    }
</script>

<TimeEntryModal
    bind:show={showModal}
    timesheet={selectedTimesheet}
    onsave={handleTimeSaved}
/>

<ConfirmDialog
    bind:show={showConfirmDelete}
    title="Delete Time Entry"
    message="Are you sure you want to delete this time entry? This action cannot be undone."
    confirmText="Delete"
    onConfirm={confirmDelete}
    onCancel={() => {
        showConfirmDelete = false;
        timesheetToDelete = null;
    }}
/>

<ContentHeader
    title="My Timesheets"
    subtitle="Log and track your working hours across all projects and tasks."
>
    <Button icon="more_time" onclick={handleCreateClick}>Log Time</Button>
</ContentHeader>

<div class="h-[calc(100vh-170px)] pb-6 flex flex-col">
    <div
        class="flex-1 bg-white rounded-2xl border border-slate-200/60 shadow-sm overflow-hidden flex flex-col"
    >
        <DataTable
            items={timesheets}
            {columns}
            {loading}
            total={totalItems}
            page={currentPage}
            {limit}
            onPageChange={handlePageChange}
            onLimitChange={(l) => {
                limit = l;
                currentPage = 1;
                loadTimesheets();
            }}
        >
            {#snippet emptyState()}
                <EmptyState
                    icon="more_time"
                    title="No time entries yet"
                    message="Start logging your hours. Track time against projects and tasks to stay on top of your workload."
                    actionLabel="Log Time"
                    actionIcon="add"
                    onaction={handleCreateClick}
                />
            {/snippet}
            {#snippet rowCell({ item, column })}
                {#if column.key === "work_date"}
                    <span class="font-medium text-slate-800"
                        >{formatDate(item.work_date)}</span
                    >
                {:else if column.key === "context"}
                    <div class="flex flex-col">
                        {#if item.task_id}
                            <div class="flex items-center gap-1.5 align-middle">
                                <span
                                    class="material-symbols-outlined text-[16px] text-indigo-500"
                                    >task_alt</span
                                >
                                <span class="font-medium text-slate-900"
                                    >{item.task_title}</span
                                >
                            </div>
                            <span class="text-xs text-slate-400 ml-5"
                                >Personal Task</span
                            >
                        {:else}
                            <div class="flex items-center gap-1.5 align-middle">
                                <span
                                    class="material-symbols-outlined text-[16px] text-sky-500"
                                    >account_tree</span
                                >
                                <span class="font-medium text-slate-900"
                                    >{item.node_title}</span
                                >
                            </div>
                            <span class="text-xs text-slate-400 ml-5"
                                >{item.project_name}</span
                            >
                        {/if}
                    </div>
                {:else if column.key === "hours"}
                    <div
                        class="font-bold text-slate-700 bg-slate-50 px-2 py-1 rounded-md inline-block"
                    >
                        {item.hours}h
                    </div>
                {:else if column.key === "description"}
                    <span
                        class="text-slate-500 text-sm italic truncate block max-w-sm"
                    >
                        {item.description || "—"}
                    </span>
                {:else if column.key === "status"}
                    <Badge variant={getStatusColor(item.status) as any}
                        >{item.status}</Badge
                    >
                {:else if column.key === "actions"}
                    <div class="flex justify-end gap-1">
                        <button
                            type="button"
                            class="p-1.5 text-slate-400 hover:text-primary hover:bg-primary/5 rounded-md transition-all"
                            onclick={() => handleEditClick(item)}
                            title="Edit entry"
                        >
                            <span class="material-symbols-outlined text-[18px]"
                                >edit</span
                            >
                        </button>
                        <button
                            type="button"
                            class="p-1.5 text-slate-400 hover:text-rose-500 hover:bg-rose-50 rounded-md transition-all"
                            onclick={() => handleDeleteClick(item)}
                            title="Delete entry"
                        >
                            <span class="material-symbols-outlined text-[18px]"
                                >delete</span
                            >
                        </button>
                    </div>
                {/if}
            {/snippet}
        </DataTable>
    </div>
</div>
