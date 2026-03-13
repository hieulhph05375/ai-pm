<script lang="ts">
    import {
        taskService,
        type Task,
        type TaskCreate,
        TASK_STATUSES,
        TASK_PRIORITIES,
    } from "$lib/services/tasks";
    import Drawer from "$lib/components/ui/Drawer.svelte";
    import Input from "$lib/components/ui/Input.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Avatar from "$lib/components/ui/Avatar.svelte";
    import Badge from "$lib/components/ui/Badge.svelte";
    import { categoryService, type Category } from "$lib/services/categories";
    import { toast } from "$lib/stores/toast";
    import TaskActivityLog from "./TaskActivity.svelte";

    interface Props {
        show: boolean;
        task?: Task | null; // null means Create mode
        onsave?: (task: Task) => void;
        onclose?: () => void;
    }

    let {
        show = $bindable(false),
        task = null,
        onsave,
        onclose,
    }: Props = $props();

    let isEditMode = $derived(!!task);
    let mode = $state<"view" | "edit" | "create">("view");
    let activeTab = $state<"details" | "activities" | "comments">("details");

    // Form state
    let taskTitle = $state("");
    let description = $state("");
    let assigneeId = $state<number | null>(null);
    let statusId = $state<number | undefined>(undefined);
    let priorityId = $state<number | undefined>(undefined);
    let startDate = $state("");
    let dueDate = $state("");
    let progress = $state(0);
    let saving = $state(false);

    let taskStatuses = $state<Category[]>([]);
    let taskPriorities = $state<Category[]>([]);

    let errors = $state({
        title: "",
        description: "",
        startDate: "",
        dueDate: "",
    });

    let isFormValid = $derived.by(() => {
        return (
            taskTitle.trim() !== "" &&
            description.trim() !== "" &&
            startDate !== "" &&
            dueDate !== "" &&
            new Date(startDate) <= new Date(dueDate)
        );
    });

    function validate() {
        let isValid = true;
        errors = { title: "", description: "", startDate: "", dueDate: "" };

        if (!taskTitle.trim()) {
            errors.title = "Task name is required";
            isValid = false;
        }
        if (!description.trim()) {
            errors.description = "Description is required";
            isValid = false;
        }
        if (!startDate) {
            errors.startDate = "Start date is required";
            isValid = false;
        }
        if (!dueDate) {
            errors.dueDate = "End date is required";
            isValid = false;
        }

        if (startDate && dueDate && new Date(startDate) > new Date(dueDate)) {
            errors.dueDate = "End date cannot be before start date";
            isValid = false;
        }

        return isValid;
    }

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

    $effect(() => {
        if (show) {
            if (taskStatuses.length === 0) loadCategories();
            mode = task ? "view" : "create";
            activeTab = "details";
            initForm();
        }
    });

    function initForm() {
        errors = { title: "", description: "", startDate: "", dueDate: "" };
        if (task) {
            taskTitle = task.title;
            description = task.description || "";
            statusId = task.status_id;
            priorityId = task.priority_id;
            startDate = task.start_date ? task.start_date.split("T")[0] : "";
            dueDate = task.due_date ? task.due_date.split("T")[0] : "";
            assigneeId = task.assignee_id;
            progress = task.progress || 0;
        } else {
            taskTitle = "";
            description = "";
            statusId = taskStatuses.find((c) => c.is_active)?.id;
            priorityId = taskPriorities.find((c) => c.is_active)?.id;
            startDate = "";
            dueDate = "";
            assigneeId = null;
            progress = 0;
        }
    }

    function handleClose() {
        show = false;
        onclose?.();
    }

    function handleEdit() {
        mode = "edit";
    }

    async function handleSave() {
        if (!validate()) {
            toast.error("Please correct the errors in the form");
            return;
        }

        saving = true;
        try {
            const selectedStatus = taskStatuses.find((c) => c.id === statusId);
            const selectedPriority = taskPriorities.find(
                (c) => c.id === priorityId,
            );

            const payload: Partial<TaskCreate> = {
                title: taskTitle,
                description,
                status: (selectedStatus?.name?.toUpperCase() || "TODO") as any,
                priority: (selectedPriority?.name?.toUpperCase() ||
                    "MEDIUM") as any,
                status_id: statusId,
                priority_id: priorityId,
                start_date: startDate
                    ? new Date(startDate).toISOString()
                    : null,
                due_date: dueDate ? new Date(dueDate).toISOString() : null,
                assignee_id: assigneeId,
                progress: progress,
            };

            let savedTask: Task;
            if (isEditMode) {
                savedTask = await taskService.update(task!.id, payload);
                toast.success("Task updated successfully");
            } else {
                savedTask = await taskService.create(payload as TaskCreate);
                toast.success("Task created successfully");
            }

            onsave?.(savedTask);
            handleClose();
        } catch (e: any) {
            toast.error(e.message || "An error occurred while saving the task");
        } finally {
            saving = false;
        }
    }
</script>

<Drawer bind:show width="max-w-2xl" onClose={handleClose}>
    {#snippet title()}
        <div class="space-y-1">
            <Badge variant="primary" class="mb-1">
                {mode === "create"
                    ? "New Personal Task"
                    : "Task #" + (task?.id || "")}
            </Badge>
            <h2
                class="text-3xl font-outfit font-bold text-slate-900 dark:text-white leading-tight"
            >
                {mode === "create"
                    ? "Create Task"
                    : mode === "edit"
                      ? "Edit Task"
                      : "Task Details"}
            </h2>
        </div>
    {/snippet}

    {#snippet tabs()}
        {#if task}
            <div
                class="flex items-center gap-1 bg-slate-50 p-1 rounded-xl border border-slate-200 w-fit"
            >
                <button
                    class="px-5 py-2 text-[11px] font-black uppercase tracking-wider rounded-lg transition-all flex items-center gap-2 {activeTab ===
                    'details'
                        ? 'bg-white text-primary shadow-sm'
                        : 'text-slate-500 hover:text-slate-700'}"
                    onclick={() => (activeTab = "details")}
                >
                    <span class="material-symbols-outlined text-[16px]"
                        >info</span
                    >
                    Details
                </button>
                <button
                    class="px-5 py-2 text-[11px] font-black uppercase tracking-wider rounded-lg transition-all flex items-center gap-2 {activeTab ===
                    'activities'
                        ? 'bg-white text-primary shadow-sm'
                        : 'text-slate-500 hover:text-slate-700'}"
                    onclick={() => (activeTab = "activities")}
                >
                    <span class="material-symbols-outlined text-[16px]"
                        >history</span
                    >
                    Activities
                </button>
                <button
                    class="px-5 py-2 text-[11px] font-black uppercase tracking-wider rounded-lg transition-all flex items-center gap-2 {activeTab ===
                    'comments'
                        ? 'bg-white text-primary shadow-sm'
                        : 'text-slate-500 hover:text-slate-700'}"
                    onclick={() => (activeTab = "comments")}
                >
                    <span class="material-symbols-outlined text-[16px]"
                        >forum</span
                    >
                    Comments
                </button>
            </div>
        {/if}
    {/snippet}

    <div class="space-y-10">
        {#if task && activeTab !== "details"}
            <div class="min-h-[400px] flex flex-col min-h-0 pt-6">
                {#if activeTab === "activities"}
                    <div class="flex-1 min-h-0 animate-in fade-in duration-300">
                        <TaskActivityLog taskId={task.id} mode="activities" />
                    </div>
                {:else if activeTab === "comments"}
                    <div class="flex-1 min-h-0 animate-in fade-in duration-300">
                        <TaskActivityLog taskId={task.id} mode="comments" />
                    </div>
                {/if}
            </div>
        {/if}

        {#if activeTab === "details"}
            <div class="space-y-10 animate-in fade-in duration-300">
                <!-- Main Form Section -->
                <div class="space-y-6">
                    <div class="grid grid-cols-2 gap-4">
                        <div class="space-y-1.5">
                            <label
                                for="task-priority"
                                class="text-[10px] font-bold text-slate-400 uppercase tracking-widest ml-1"
                                >Priority</label
                            >
                            <div class="relative group">
                                <select
                                    id="task-priority"
                                    bind:value={priorityId}
                                    disabled={mode === "view"}
                                    class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl text-sm font-bold focus:bg-white focus:ring-4 focus:ring-primary/10 focus:border-primary outline-none transition-all appearance-none cursor-pointer disabled:opacity-70 disabled:cursor-default"
                                >
                                    {#each taskPriorities as p}
                                        <option value={p.id}>{p.name}</option>
                                    {/each}
                                </select>
                                <span
                                    class="material-symbols-outlined absolute right-4 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none text-[18px]"
                                    >unfold_more</span
                                >
                            </div>
                        </div>
                        <div class="space-y-1.5">
                            <label
                                for="task-status"
                                class="text-[10px] font-bold text-slate-400 uppercase tracking-widest ml-1"
                                >Status</label
                            >
                            <div class="relative group">
                                <select
                                    id="task-status"
                                    bind:value={statusId}
                                    disabled={mode === "view"}
                                    class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl text-sm font-bold focus:bg-white focus:ring-4 focus:ring-primary/10 focus:border-primary outline-none transition-all appearance-none cursor-pointer disabled:opacity-70 disabled:cursor-default"
                                >
                                    {#each taskStatuses as s}
                                        <option value={s.id}>{s.name}</option>
                                    {/each}
                                </select>
                                <span
                                    class="material-symbols-outlined absolute right-4 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none text-[18px]"
                                    >unfold_more</span
                                >
                            </div>
                        </div>
                    </div>

                    <!-- Progress Slider -->
                    <div
                        class="space-y-3 p-4 bg-slate-50 border border-slate-200 rounded-2xl"
                    >
                        <div class="flex items-center justify-between ml-1">
                            <label
                                for="task-progress"
                                class="text-[10px] font-bold text-slate-400 uppercase tracking-widest"
                                >Task Progress</label
                            >
                            <span
                                class="text-xs font-black text-primary px-2 py-0.5 bg-white rounded-md border border-slate-200"
                            >
                                {progress}%
                            </span>
                        </div>
                        <div class="relative flex items-center h-6">
                            <input
                                id="task-progress"
                                type="range"
                                min="0"
                                max="100"
                                step="1"
                                bind:value={progress}
                                disabled={mode === "view"}
                                class="w-full h-1.5 bg-slate-200 rounded-lg appearance-none cursor-pointer accent-primary disabled:opacity-50 disabled:cursor-default"
                            />
                            <div
                                class="absolute left-0 top-1/2 -translate-y-1/2 h-1.5 bg-primary rounded-l-lg pointer-events-none"
                                style="width: {progress}%"
                            ></div>
                        </div>
                    </div>

                    <Input
                        label="Task Name"
                        bind:value={taskTitle}
                        placeholder="What needs to be done?"
                        required
                        error={errors.title}
                        disabled={mode === "view"}
                        class="!text-lg !font-bold"
                    />

                    <div class="space-y-1.5">
                        <label
                            class="text-[10px] font-bold {errors.description
                                ? 'text-rose-500'
                                : 'text-slate-400'} uppercase tracking-widest ml-1"
                            for="desc">Description</label
                        >
                        <textarea
                            id="desc"
                            bind:value={description}
                            placeholder="Add detailed notes and requirements..."
                            disabled={mode === "view"}
                            class="w-full px-5 py-4 bg-slate-50 border {errors.description
                                ? 'border-rose-500 focus:ring-rose-500/10'
                                : 'border-slate-200 focus:ring-primary/10'} rounded-2xl text-sm focus:bg-white focus:ring-4 focus:border-primary outline-none transition-all placeholder:text-slate-400 min-h-[120px] resize-none leading-relaxed disabled:opacity-70 disabled:cursor-default shadow-inner"
                        ></textarea>
                        {#if errors.description}
                            <p
                                class="text-[11px] font-bold text-rose-500 ml-1 mt-1 flex items-center gap-1 animate-in fade-in slide-in-from-top-1 duration-200"
                            >
                                <span
                                    class="material-symbols-outlined text-[14px]"
                                    >error</span
                                >
                                {errors.description}
                            </p>
                        {/if}
                    </div>
                </div>

                <!-- Resources & Schedule Section -->
                <div class="space-y-6 pt-6 border-t border-slate-100">
                    <h3
                        class="text-[11px] font-black text-slate-400 uppercase tracking-[0.2em] flex items-center gap-2 mb-2"
                    >
                        <span class="material-symbols-outlined text-[18px]"
                            >group</span
                        >
                        Resources & Schedule
                    </h3>

                    <div class="space-y-4">
                        <div class="space-y-1.5">
                            <label
                                for="task-assignee"
                                class="text-[10px] font-bold text-slate-400 uppercase tracking-widest ml-1"
                                >Responsible Party</label
                            >
                            <div class="relative">
                                <select
                                    id="task-assignee"
                                    bind:value={assigneeId}
                                    disabled={mode === "view"}
                                    class="w-full pl-12 pr-10 py-3 bg-slate-50 border border-slate-200 rounded-xl text-sm font-semibold focus:bg-white focus:ring-4 focus:ring-primary/10 focus:border-primary outline-none transition-all appearance-none cursor-pointer disabled:opacity-70 disabled:cursor-default"
                                >
                                    <option value={null}>Unassigned</option>
                                    <option value={1}
                                        >System Admin (User 1)</option
                                    >
                                </select>
                                <div
                                    class="absolute left-3 top-1/2 -translate-y-1/2"
                                >
                                    <Avatar
                                        name={assigneeId ? "User 1" : "Guest"}
                                        size="xs"
                                    />
                                </div>
                                <span
                                    class="material-symbols-outlined absolute right-4 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none text-[18px]"
                                    >expand_more</span
                                >
                            </div>
                        </div>

                        <div class="grid grid-cols-2 gap-4">
                            <div class="space-y-1.5">
                                <label
                                    for="task-start-date"
                                    class="text-[10px] font-bold {errors.startDate
                                        ? 'text-rose-500'
                                        : 'text-slate-400'} uppercase tracking-widest ml-1"
                                    >Start Date</label
                                >
                                <div class="relative">
                                    <input
                                        type="date"
                                        id="task-start-date"
                                        bind:value={startDate}
                                        disabled={mode === "view"}
                                        class="w-full pl-4 pr-10 py-3 bg-slate-50 border {errors.startDate
                                            ? 'border-rose-500'
                                            : 'border-slate-200'} rounded-xl text-sm font-bold focus:bg-white focus:ring-4 focus:border-primary/10 transition-all outline-none disabled:opacity-70"
                                    />
                                    <span
                                        class="material-symbols-outlined absolute right-4 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none text-[18px]"
                                        >calendar_today</span
                                    >
                                </div>
                            </div>
                            <div class="space-y-1.5">
                                <label
                                    for="task-end-date"
                                    class="text-[10px] font-bold {errors.dueDate
                                        ? 'text-rose-500'
                                        : 'text-slate-400'} uppercase tracking-widest ml-1"
                                    >End Date</label
                                >
                                <div class="relative">
                                    <input
                                        type="date"
                                        id="task-end-date"
                                        bind:value={dueDate}
                                        disabled={mode === "view"}
                                        class="w-full pl-4 pr-10 py-3 bg-slate-50 border {errors.dueDate
                                            ? 'border-rose-500'
                                            : 'border-slate-200'} rounded-xl text-sm font-bold focus:bg-white focus:ring-4 focus:border-primary/10 transition-all outline-none disabled:opacity-70"
                                    />
                                    <span
                                        class="material-symbols-outlined absolute right-4 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none text-[18px]"
                                        >event_available</span
                                    >
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        {/if}
    </div>

    {#snippet footer()}
        <div class="flex-1"></div>
        <Button variant="outline" onclick={handleClose} disabled={saving}
            >Cancel</Button
        >
        {#if mode === "view"}
            <Button
                onclick={handleEdit}
                icon="edit"
                class="bg-blue-600 hover:bg-blue-700 text-white shadow-lg shadow-blue-200"
            >
                Edit Task
            </Button>
        {:else}
            <Button
                onclick={handleSave}
                disabled={saving || !isFormValid}
                icon={mode === "edit" ? "save" : "add_task"}
                class="bg-primary hover:bg-primary/90 text-white shadow-lg shadow-primary/20 disabled:opacity-50"
            >
                {mode === "edit" ? "Save Changes" : "Create Task"}
            </Button>
        {/if}
    {/snippet}
</Drawer>
