<script lang="ts">
    import Modal from "$lib/components/ui/Modal.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import { projectService, type Project } from "$lib/services/projects";
    import { wbsService, type WBSNode } from "$lib/services/wbs";
    import { taskService, type Task } from "$lib/services/tasks";
    import { timesheetService, type Timesheet } from "$lib/services/timesheets";
    import { toast } from "$lib/stores/toast";

    let {
        show = $bindable(false),
        timesheet = null,
        initialTaskId = null,
        onsave,
    } = $props<{
        show: boolean;
        timesheet?: Timesheet | null;
        initialTaskId?: number | null;
        onsave: (t: Timesheet) => void;
    }>();

    let loading = $state(false);

    let formData = $state({
        project_id: null as number | null,
        node_id: null as number | null,
        task_id: null as number | null,
        work_date: new Date().toISOString().split("T")[0],
        hours: 0,
        description: "",
        is_personal_task: false,
    });

    // For dropdowns
    let projects = $state<Project[]>([]);
    let wbsNodes = $state<WBSNode[]>([]);
    let personalTasks = $state<Task[]>([]);

    // Flag to avoid infinite loops when watching `show`
    let initializing = $state(false);

    $effect(() => {
        if (show) {
            initializing = true;
            loadInitialData().finally(() => {
                if (timesheet) {
                    formData = {
                        project_id: timesheet.project_id || null,
                        node_id: timesheet.node_id || null,
                        task_id: timesheet.task_id || null,
                        work_date: timesheet.work_date.split("T")[0],
                        hours: timesheet.hours,
                        description: timesheet.description || "",
                        is_personal_task: !!timesheet.task_id,
                    };
                    if (formData.project_id) {
                        loadWbsNodes(formData.project_id);
                    }
                } else {
                    formData = {
                        project_id: null,
                        node_id: null,
                        task_id: initialTaskId ?? null,
                        work_date: new Date().toISOString().split("T")[0],
                        hours: 0,
                        description: "",
                        is_personal_task: !!initialTaskId,
                    };
                }
                initializing = false;
            });
        }
    });

    async function loadInitialData() {
        try {
            const projRes = await projectService.list();
            projects = projRes.data || [];

            const tasksRes = await taskService.list(1, 100); // Load recently assigned/created tasks
            personalTasks = tasksRes.items.filter((t) => t.status !== "DONE");
        } catch (e: any) {
            toast.error("Failed to load select dropdowns");
        }
    }

    async function loadWbsNodes(projectId: number) {
        try {
            // Only load leaves or all nodes to select from for time logging
            const res = await wbsService.listTree(projectId);
            wbsNodes = res.data;
            // Flatten tree for dropdown
            wbsNodes = flattenWBS(wbsNodes);
        } catch (e: any) {
            toast.error("Failed to load project tasks (WBS)");
        }
    }

    function flattenWBS(nodes: WBSNode[], depth: number = 0): WBSNode[] {
        let flat: WBSNode[] = [];
        for (const n of nodes) {
            // Add visual indicator of depth to title temporarily for the dropdown
            const prefix = "—".repeat(depth) + (depth > 0 ? " " : "");
            flat.push({ ...n, title: prefix + n.title });
            if (n.children && n.children.length > 0) {
                flat = flat.concat(flattenWBS(n.children, depth + 1));
            }
        }
        return flat;
    }

    function onProjectChange(event: Event) {
        const target = event.target as HTMLSelectElement;
        const pid = target.value ? parseInt(target.value) : null;
        formData.project_id = pid;
        formData.node_id = null; // reset wbs node when project changes

        if (pid) {
            loadWbsNodes(pid);
        } else {
            wbsNodes = [];
        }
    }

    async function handleSubmit(e: SubmitEvent) {
        e.preventDefault();
        if (formData.hours <= 0) {
            toast.error("Hours must be greater than 0");
            return;
        }

        if (formData.is_personal_task && !formData.task_id) {
            toast.error("Please select a personal task");
            return;
        }

        if (
            !formData.is_personal_task &&
            (!formData.project_id || !formData.node_id)
        ) {
            toast.error("Please select a Project and WBS Task");
            return;
        }

        loading = true;
        try {
            const payload: Partial<Timesheet> = {
                work_date: new Date(formData.work_date).toISOString(),
                hours: Number(formData.hours),
                description: formData.description,
            };

            if (formData.is_personal_task) {
                payload.task_id = formData.task_id;
                payload.project_id = null;
                payload.node_id = null;
            } else {
                payload.project_id = formData.project_id;
                payload.node_id = formData.node_id;
                payload.task_id = null;
            }

            let savedTimesheet: Timesheet;
            if (timesheet?.id) {
                savedTimesheet = await timesheetService.update(
                    timesheet.id,
                    payload,
                );
                toast.success("Timesheet updated");
            } else {
                savedTimesheet = await timesheetService.create(payload);
                toast.success("Time logged successfully");
            }

            onsave(savedTimesheet);
            show = false;
        } catch (err: any) {
            toast.error(err.message || "Error saving timesheet");
        } finally {
            loading = false;
        }
    }
</script>

<Modal bind:show title={timesheet ? "Edit Time Entry" : "Log Time"}>
    <form onsubmit={handleSubmit} class="space-y-4">
        <!-- Entry Type Toggle -->
        <div class="flex items-center space-x-6 pb-2 border-b border-slate-100">
            <label class="flex items-center space-x-2 cursor-pointer">
                <input
                    type="radio"
                    name="entry_type"
                    checked={!formData.is_personal_task}
                    onchange={() => {
                        formData.is_personal_task = false;
                        formData.task_id = null;
                    }}
                    class="w-4 h-4 text-primary border-slate-300 focus:ring-primary"
                />
                <span class="text-sm font-medium text-slate-700"
                    >Project WBS Task</span
                >
            </label>
            <label class="flex items-center space-x-2 cursor-pointer">
                <input
                    type="radio"
                    name="entry_type"
                    checked={formData.is_personal_task}
                    onchange={() => {
                        formData.is_personal_task = true;
                        formData.project_id = null;
                        formData.node_id = null;
                    }}
                    class="w-4 h-4 text-primary border-slate-300 focus:ring-primary"
                />
                <span class="text-sm font-medium text-slate-700"
                    >Personal Task</span
                >
            </label>
        </div>

        {#if formData.is_personal_task}
            <div>
                <label
                    for="personal-task"
                    class="block text-sm font-medium text-slate-700 mb-1"
                    >Select Personal Task *</label
                >
                <select
                    id="personal-task"
                    bind:value={formData.task_id}
                    class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/50 text-slate-900 bg-white"
                    required
                >
                    <option value={null}>-- Select a Task --</option>
                    {#each personalTasks as targetTask}
                        <option value={targetTask.id}>{targetTask.title}</option
                        >
                    {/each}
                </select>
            </div>
        {:else}
            <div>
                <label
                    for="project-select"
                    class="block text-sm font-medium text-slate-700 mb-1"
                    >Project *</label
                >
                <select
                    id="project-select"
                    value={formData.project_id}
                    onchange={onProjectChange}
                    class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/50 text-slate-900 bg-white"
                    required
                >
                    <option value={null}>-- Select Project --</option>
                    {#each projects as proj}
                        <option value={proj.id}>{proj.project_name}</option>
                    {/each}
                </select>
            </div>
            <div>
                <label
                    for="wbs-select"
                    class="block text-sm font-medium text-slate-700 mb-1"
                    >WBS Task *</label
                >
                <select
                    id="wbs-select"
                    bind:value={formData.node_id}
                    class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/50 text-slate-900 bg-white"
                    required
                    disabled={!formData.project_id || wbsNodes.length === 0}
                >
                    <option value={null}>-- Select WBS Task --</option>
                    {#each wbsNodes as node}
                        <option
                            value={node.id}
                            disabled={node.type !== "Task" &&
                                node.type !== "Sub-task"}
                        >
                            {node.title}
                            {node.type !== "Task" && node.type !== "Sub-task"
                                ? `(${node.type} - Cannot log directly)`
                                : ""}
                        </option>
                    {/each}
                </select>
            </div>
        {/if}

        <div class="grid grid-cols-2 gap-4">
            <div>
                <label
                    for="work-date"
                    class="block text-sm font-medium text-slate-700 mb-1"
                    >Date *</label
                >
                <input
                    id="work-date"
                    type="date"
                    bind:value={formData.work_date}
                    required
                    class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/50 text-slate-900"
                />
            </div>
            <div>
                <label
                    for="hours"
                    class="block text-sm font-medium text-slate-700 mb-1"
                    >Hours *</label
                >
                <input
                    id="hours"
                    type="number"
                    bind:value={formData.hours}
                    required
                    min="0.25"
                    max="24"
                    step="0.25"
                    class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/50 text-slate-900"
                />
            </div>
        </div>

        <div>
            <label
                for="description"
                class="block text-sm font-medium text-slate-700 mb-1"
                >Description / Notes</label
            >
            <textarea
                id="description"
                bind:value={formData.description}
                rows="3"
                placeholder="What did you work on?"
                class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/50 text-slate-900"
            ></textarea>
        </div>

        <div class="flex justify-end gap-3 pt-4 border-t border-slate-100">
            <Button
                variant="outline"
                type="button"
                onclick={() => (show = false)}
            >
                Cancel
            </Button>
            <Button type="submit" {loading}>
                {timesheet ? "Update Entry" : "Log Time"}
            </Button>
        </div>
    </form>
</Modal>
