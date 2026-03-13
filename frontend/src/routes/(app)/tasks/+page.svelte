<script lang="ts">
	import { onMount } from "svelte";
	import { taskService, type Task } from "$lib/services/tasks";
	import { categoryService, type Category } from "$lib/services/categories";
	import { toast } from "$lib/stores/toast";
	import ContentHeader from "$lib/components/ui/ContentHeader.svelte";
	import Button from "$lib/components/ui/Button.svelte";
	import KanbanBoard from "$lib/components/tasks/KanbanBoard.svelte";
	import PersonalGanttChart from "$lib/components/tasks/PersonalGanttChart.svelte";
	import TaskModal from "$lib/components/tasks/TaskModal.svelte";
	import ConfirmDialog from "$lib/components/ui/ConfirmDialog.svelte";
	import DataTable from "$lib/components/ui/DataTable.svelte";
	import Badge from "$lib/components/ui/Badge.svelte";
	import TaskStats from "$lib/components/tasks/TaskStats.svelte";
	import EmptyState from "$lib/components/ui/EmptyState.svelte";
	import { TASK_STATUSES, TASK_PRIORITIES } from "$lib/services/tasks";
	import TimeEntryModal from "$lib/components/timesheets/TimeEntryModal.svelte";

	let tasks = $state<Task[]>([]);
	// Pagination State
	let currentPage = $state(1);
	let limit = $state(10);
	let totalTasks = $state(0);

	let loading = $state(true);

	// Tabs: 'board' | 'gantt' | 'list'
	let activeTab = $state<"board" | "gantt" | "list">("board");

	// Modal State
	let showModal = $state(false);
	let selectedTask = $state<Task | null>(null);

	// Confirm Delete State
	let showConfirmDelete = $state(false);
	let taskToDelete = $state<Task | null>(null);

	// Log Time State
	let showTimeEntry = $state(false);
	let logTimeTaskId = $state<number | null>(null);

	let taskStatuses = $state<Category[]>([]);
	let taskPriorities = $state<Category[]>([]);

	function handleLogTimeFromTask(task: Task) {
		logTimeTaskId = task.id;
		showTimeEntry = true;
	}

	onMount(async () => {
		await loadTasks();
	});

	async function loadCategories() {
		try {
			const [s, p] = await Promise.all([
				categoryService.listCategories(1, 100, "", 6), // Task Status
				categoryService.listCategories(1, 100, "", 7), // Task Priority
			]);
			taskStatuses = s.data;
			taskPriorities = p.data;
		} catch (err) {
			console.error("Failed to load task categories", err);
		}
	}

	async function loadTasks() {
		loading = true;
		if (taskStatuses.length === 0) await loadCategories();
		try {
			const result = await taskService.list(currentPage, limit);
			tasks = result.items;
			totalTasks = result.total;
		} catch (e: any) {
			if (!e.isAuthError) {
				console.error(e);
				toast.error("Could not load tasks");
			}
		} finally {
			loading = false;
		}
	}

	async function handlePageChange(page: number) {
		currentPage = page;
		await loadTasks();
	}

	async function handleStatusChange(task: Task, newStatusId: number) {
		// Optimistic UI update
		const originalStatusId = task.status_id;
		const originalStatusCat = task.status_cat;

		const newStatus = taskStatuses.find((s) => s.id === newStatusId);
		task.status_id = newStatusId;
		task.status_cat = newStatus;
		if (newStatus) task.status = newStatus.name as any;

		tasks = [...tasks];

		try {
			await taskService.updateStatus(task.id, task.status, task);
			toast.success(`Moved to ${newStatus?.name || "new status"}`);
		} catch (e: any) {
			// Revert on failure
			task.status_id = originalStatusId;
			task.status_cat = originalStatusCat;
			tasks = [...tasks];
			toast.error(e.message || "Error updating status");
		}
	}

	function handleTaskClick(task: Task) {
		selectedTask = task;
		showModal = true;
	}

	function handleCreateClick() {
		selectedTask = null;
		showModal = true;
	}

	function handleTaskSave(savedTask: Task) {
		const index = tasks.findIndex((t) => t.id === savedTask.id);
		if (index >= 0) {
			// Update existing
			tasks[index] = savedTask;
		} else {
			// Add new
			tasks = [savedTask, ...tasks];
		}
		tasks = [...tasks]; // trigger reactivity
	}

	function handleTaskDelete(task: Task) {
		taskToDelete = task;
		showConfirmDelete = true;
	}

	async function confirmDelete() {
		if (!taskToDelete) return;

		try {
			await taskService.delete(taskToDelete.id);
			toast.success("Task deleted successfully");
			tasks = tasks.filter((t) => t.id !== taskToDelete!.id);
			if (selectedTask?.id === taskToDelete.id) {
				showModal = false;
				selectedTask = null;
			}
		} catch (e: any) {
			toast.error(
				e.message || "An error occurred while deleting the task",
			);
		} finally {
			showConfirmDelete = false;
			taskToDelete = null;
		}
	}
	async function handleTaskUpdate(task: Task, updates: Partial<Task>) {
		// Optimistically update
		const oldStart = task.start_date;
		const oldDue = task.due_date;

		task.start_date = updates.start_date as string | null;
		task.due_date = updates.due_date as string | null;
		tasks = [...tasks];

		try {
			await taskService.update(task.id, updates);
			toast.success("Timeline updated");
		} catch (e: any) {
			task.start_date = oldStart;
			task.due_date = oldDue;
			tasks = [...tasks];
			toast.error(e.message || "Error saving timeline");
		}
	}

	// Columns for DataTable
	const listColumns = [
		{ key: "title", label: "TASK NAME" },
		{ key: "status", label: "STATUS", class: "w-40" },
		{ key: "priority", label: "PRIORITY", class: "w-40" },
		{ key: "due_date", label: "DUE DATE", class: "w-40" },
		{
			key: "actions",
			label: "",
			class: "w-16 right",
			align: "right" as const,
		},
	];

	function getStatusInfo(task: Task) {
		if (task.status_cat) {
			return {
				label: task.status_cat.name,
				color: task.status_cat.color || "slate",
			};
		}
		return (
			TASK_STATUSES.find((s) => s.value === task.status) || {
				label: task.status,
				color: "slate",
			}
		);
	}

	function getPriorityInfo(task: Task) {
		if (task.priority_cat) {
			return {
				label: task.priority_cat.name,
				color: task.priority_cat.color || "slate",
			};
		}
		return (
			TASK_PRIORITIES.find((p) => p.value === task.priority) || {
				label: task.priority,
				color: "slate",
			}
		);
	}

	function formatDate(val: string | null) {
		if (!val) return "—";
		return new Date(val).toLocaleDateString("vi-VN", {
			day: "2-digit",
			month: "2-digit",
			year: "numeric",
		});
	}
</script>

<TaskModal bind:show={showModal} task={selectedTask} onsave={handleTaskSave} />

<TimeEntryModal
	bind:show={showTimeEntry}
	initialTaskId={logTimeTaskId}
	onsave={() => toast.success("Time logged!")}
/>

<ContentHeader
	title="Task Management"
	subtitle="Track and manage personal tasks (Kanban & Timeline)"
>
	<div
		class="flex items-center bg-slate-100 p-1 rounded-lg border border-slate-200 mr-2"
	>
		<button
			class="px-4 py-1.5 text-sm font-medium rounded-md transition-colors {activeTab ===
			'board'
				? 'bg-white text-primary shadow-sm'
				: 'text-slate-500 hover:text-slate-700'}"
			onclick={() => (activeTab = "board")}
		>
			<span
				class="material-symbols-outlined text-[18px] align-middle mr-1 border-b border-transparent"
				>view_kanban</span
			>
			Kanban Board
		</button>
		<button
			class="px-4 py-1.5 text-sm font-medium rounded-md transition-colors {activeTab ===
			'gantt'
				? 'bg-white text-primary shadow-sm'
				: 'text-slate-500 hover:text-slate-700'}"
			onclick={() => (activeTab = "gantt")}
		>
			<span
				class="material-symbols-outlined text-[18px] align-middle mr-1 border-b border-transparent"
				>timeline</span
			>
			Timeline (Gantt)
		</button>
		<button
			class="px-4 py-1.5 text-sm font-medium rounded-md transition-colors {activeTab ===
			'list'
				? 'bg-white text-primary shadow-sm'
				: 'text-slate-500 hover:text-slate-700'}"
			onclick={() => (activeTab = "list")}
		>
			<span
				class="material-symbols-outlined text-[18px] align-middle mr-1 border-b border-transparent"
				>format_list_bulleted</span
			>
			List View
		</button>
	</div>

	<Button icon="add" onclick={handleCreateClick}>Create Task</Button>
</ContentHeader>

<div class="relative z-10 -mt-2 mb-4" style="zoom: 0.7;">
	<TaskStats {tasks} />
</div>

<div class="pb-6 h-[calc(100vh-215px)] overflow-hidden flex flex-col">
	{#if loading}
		<div class="flex items-center justify-center h-full">
			<div
				class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"
			></div>
		</div>
	{:else if activeTab === "board"}
		{#if tasks.length === 0}
			<div
				class="flex-1 flex flex-col justify-center py-12 bg-white rounded-2xl border border-slate-200"
			>
				<EmptyState
					icon="task_alt"
					title="No tasks found"
					message="You don't have any tasks scheduled. Create a new task to get started."
					actionLabel="Create Task"
					onaction={handleCreateClick}
				/>
			</div>
		{:else}
			<KanbanBoard
				{tasks}
				statusCategories={taskStatuses}
				onstatuschange={handleStatusChange}
				ontaskclick={handleTaskClick}
				ondelete={handleTaskDelete}
				onlogtime={handleLogTimeFromTask}
			/>
		{/if}
	{:else if activeTab === "gantt"}
		{#if tasks.length === 0}
			<div
				class="flex-1 flex flex-col justify-center py-12 bg-white rounded-2xl border border-slate-200"
			>
				<EmptyState
					icon="task_alt"
					title="No tasks found"
					message="You don't have any tasks scheduled. Create a new task to get started."
					actionLabel="Create Task"
					onaction={handleCreateClick}
				/>
			</div>
		{:else}
			<PersonalGanttChart
				{tasks}
				page={currentPage}
				total={totalTasks}
				{limit}
				onPageChange={handlePageChange}
				onLimitChange={(l) => {
					limit = l;
					currentPage = 1;
					loadTasks();
				}}
				ontaskclick={handleTaskClick}
				ontaskupdate={handleTaskUpdate}
				ondelete={handleTaskDelete}
			/>
		{/if}
	{:else if activeTab === "list"}
		<DataTable
			items={tasks}
			columns={listColumns}
			loading={false}
			total={totalTasks}
			page={currentPage}
			{limit}
			onPageChange={handlePageChange}
			onLimitChange={(l) => {
				limit = l;
				currentPage = 1;
				loadTasks();
			}}
		>
			{#snippet emptyState()}
				<EmptyState
					icon="task_alt"
					title="No tasks found"
					message="You don't have any tasks scheduled. Create a new task to get started."
					actionLabel="Create Task"
					onaction={handleCreateClick}
				/>
			{/snippet}
			{#snippet rowCell({ item, column })}
				{#if column.key === "title"}
					<div
						class="font-semibold text-slate-900 cursor-pointer hover:text-primary transition-colors"
						onclick={() => handleTaskClick(item)}
						role="button"
						tabindex="0"
						onkeydown={(e) =>
							e.key === "Enter" && handleTaskClick(item)}
					>
						{item.title}
					</div>
				{:else if column.key === "status"}
					{@const st = getStatusInfo(item)}
					<Badge color={st.color as any}>{st.label}</Badge>
				{:else if column.key === "priority"}
					{@const pr = getPriorityInfo(item)}
					<Badge color={pr.color as any}>{pr.label}</Badge>
				{:else if column.key === "due_date"}
					<span class="text-slate-500 font-medium text-sm"
						>{formatDate(item.due_date)}</span
					>
				{:else if column.key === "actions"}
					<div class="flex justify-end pr-2">
						<button
							type="button"
							class="p-1.5 text-slate-300 hover:text-rose-500 hover:bg-rose-50 rounded-md transition-all"
							onclick={() => handleTaskDelete(item)}
							title="Delete task"
						>
							<span class="material-symbols-outlined text-[18px]"
								>delete</span
							>
						</button>
					</div>
				{/if}
			{/snippet}
		</DataTable>
	{/if}
</div>

<ConfirmDialog
	bind:show={showConfirmDelete}
	title="Confirm Deletion"
	message={`Are you sure you want to delete task "${taskToDelete?.title || ""}"? This action cannot be undone.`}
	confirmText="Delete Task"
	onConfirm={confirmDelete}
	onCancel={() => {
		showConfirmDelete = false;
		taskToDelete = null;
	}}
/>
