<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/state";
    import {
        stakeholderService,
        type ProjectStakeholder,
        type Stakeholder,
    } from "$lib/services/stakeholders";
    import { projectService, type Project } from "$lib/services/projects";
    import { toast } from "$lib/stores/toast";
    import EmptyState from "$lib/components/ui/EmptyState.svelte";
    import DataTable from "$lib/components/ui/DataTable.svelte";
    import ConfirmDialog from "$lib/components/ui/ConfirmDialog.svelte";
    import { categoryService, type Category } from "$lib/services/categories";
    import Avatar from "$lib/components/ui/Avatar.svelte";

    let projectId = $derived(Number(page.params.id));
    let project = $state<Project | null>(null);
    let projectStakeholders = $state<ProjectStakeholder[]>([]);
    let allStakeholders = $state<Stakeholder[]>([]);
    let isLoading = $state(true);
    let isAssignModalOpen = $state(false);
    let stakeholderRoles = $state<Category[]>([]);

    // Confirm Dialog State
    let isConfirmOpen = $state(false);
    let stakeholderToUnassign = $state<number | null>(null);

    let newAssignment = $state({
        stakeholder_id: 0,
        role: "",
        role_id: undefined as number | undefined,
    });

    const columns = [
        { key: "name", label: "Name & System Role" },
        { key: "role", label: "Project Role" },
        { key: "contact", label: "Contact" },
        { key: "actions", label: "Action", align: "right" as const },
    ];

    async function loadData() {
        isLoading = true;
        try {
            const [p, ps, sr, roles] = await Promise.all([
                projectService.get(projectId),
                stakeholderService.listByProject(projectId),
                stakeholderService.list(1, 1000, ""), // Get many for the dropdown
                categoryService.listCategories(1, 100, "", 9), // Role type
            ]);
            project = p;
            projectStakeholders = ps;
            allStakeholders = sr.data;
            stakeholderRoles = roles.data;
        } catch (error: any) {
            if (!error.isAuthError) {
                toast.error("Could not load data");
            }
        } finally {
            isLoading = false;
        }
    }

    async function handleAssign() {
        if (!newAssignment.stakeholder_id) return;
        // Sync role name
        const selectedRole = stakeholderRoles.find(
            (c) => c.id === newAssignment.role_id,
        );
        if (selectedRole) {
            newAssignment.role = selectedRole.name;
        }

        try {
            await stakeholderService.assignToProject(
                projectId,
                newAssignment.stakeholder_id,
                newAssignment.role,
                newAssignment.role_id,
            );
            toast.success("Stakeholder assigned successfully");
            isAssignModalOpen = false;
            loadData();
        } catch (error: any) {
            if (!error.isAuthError) {
                toast.error("An error occurred");
            }
        }
    }

    function confirmUnassign(sid: number) {
        stakeholderToUnassign = sid;
        isConfirmOpen = true;
    }

    async function handleUnassign() {
        if (!stakeholderToUnassign) return;
        try {
            await stakeholderService.unassignFromProject(
                projectId,
                stakeholderToUnassign,
            );
            toast.success("Removed");
            isConfirmOpen = false;
            stakeholderToUnassign = null;
            loadData();
        } catch (error: any) {
            if (!error.isAuthError) {
                toast.error("Could not remove");
            }
        }
    }

    onMount(loadData);
</script>

<div class="max-w-7xl mx-auto space-y-6">
    <!-- Breadcrumb -->
    <a
        href="/projects/{projectId}"
        class="inline-flex items-center gap-2 text-sm font-bold text-slate-500 hover:text-primary transition-all"
    >
        <span class="material-symbols-outlined text-[18px]">arrow_back</span>
        Back to project details
    </a>

    <div
        class="flex justify-between items-center bg-white p-8 rounded-3xl border border-slate-200 shadow-xl shadow-slate-200/50"
    >
        <div>
            <h1
                class="text-3xl font-display font-bold text-slate-900 tracking-tight"
            >
                Stakeholders - {project?.project_name || "..."}
            </h1>
            <p class="text-slate-500 mt-1">
                Manage stakeholders directly involved in this project
            </p>
        </div>
        <button
            onclick={() => (isAssignModalOpen = true)}
            class="bg-primary text-white px-5 py-2.5 rounded-xl font-medium shadow-lg shadow-primary/20 hover:scale-[1.02] active:scale-[0.98] transition-all flex items-center gap-2"
        >
            <span class="material-symbols-outlined">person_add</span>
            Assign Stakeholder
        </button>
    </div>

    <DataTable
        items={projectStakeholders}
        {columns}
        loading={isLoading}
        total={projectStakeholders.length}
    >
        {#snippet emptyState()}
            <EmptyState
                icon="person_add"
                actionIcon="person_add"
                title="No stakeholders yet"
                message="This project has no stakeholders assigned yet."
                actionLabel="Assign Stakeholder"
                onaction={() => (isAssignModalOpen = true)}
            />
        {/snippet}

        {#snippet rowCell({
            item,
            column,
        }: {
            item: ProjectStakeholder;
            column: any;
        })}
            {#if column.key === "name"}
                <div class="flex items-center gap-3">
                    <Avatar name={item.stakeholder.name} size="sm" />
                    <div>
                        <div
                            class="font-bold text-slate-900 group-hover:text-primary transition-colors text-sm"
                        >
                            {item.stakeholder.name}
                        </div>
                        <div
                            class="text-[11px] text-slate-400 font-medium uppercase tracking-tighter mt-0.5"
                        >
                            {item.stakeholder.role_cat?.name ||
                                item.stakeholder.role} | {item.stakeholder
                                .organization}
                        </div>
                    </div>
                </div>
            {:else if column.key === "role"}
                <span
                    class="bg-primary/5 text-primary text-[11px] font-bold px-3 py-1.5 rounded-full border border-primary/10"
                    style={item.role_cat?.color
                        ? `background-color: ${item.role_cat.color}15; color: ${item.role_cat.color}; border-color: ${item.role_cat.color}30;`
                        : ""}
                >
                    {item.role_cat?.name || item.project_role || "Member"}
                </span>
            {:else if column.key === "contact"}
                <div class="space-y-1">
                    <div
                        class="text-sm text-slate-600 font-medium flex items-center gap-1.5"
                    >
                        <span
                            class="material-symbols-outlined text-[16px] text-slate-400"
                            >mail</span
                        >
                        {item.stakeholder.email || "-"}
                    </div>
                    <div
                        class="text-xs text-slate-400 flex items-center gap-1.5"
                    >
                        <span
                            class="material-symbols-outlined text-[16px] text-slate-400"
                            >phone</span
                        >
                        {item.stakeholder.phone || "-"}
                    </div>
                </div>
            {:else if column.key === "actions"}
                <div
                    class="flex justify-end gap-1 opacity-0 group-hover:opacity-100 transition-opacity"
                >
                    <button
                        onclick={() => confirmUnassign(item.stakeholder_id)}
                        class="p-2 text-slate-400 hover:text-rose-600 hover:bg-rose-50 rounded-lg transition-all"
                        title="Remove from project"
                    >
                        <span class="material-symbols-outlined text-[18px]"
                            >person_remove</span
                        >
                    </button>
                </div>
            {/if}
        {/snippet}
    </DataTable>
</div>

<!-- Assign Modal -->
{#if isAssignModalOpen}
    <div class="fixed inset-0 z-[100] flex items-center justify-center p-4">
        <button
            class="absolute inset-0 bg-slate-900/40 backdrop-blur-sm border-none w-full h-full cursor-default"
            onclick={() => (isAssignModalOpen = false)}
            aria-label="Close modal"
        ></button>
        <div
            class="bg-white rounded-3xl w-full max-w-md relative shadow-2xl animate-in fade-in zoom-in duration-200 overflow-hidden"
        >
            <div class="p-8">
                <h3 class="font-display font-bold text-xl text-slate-900 mb-2">
                    Assign Stakeholder
                </h3>
                <p class="text-slate-500 text-sm mb-6">
                    Select from the system-wide stakeholder list to assign to
                    this project.
                </p>

                <div class="space-y-4">
                    <div>
                        <label
                            for="stakeholder-select"
                            class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2 ml-1"
                            >Select Stakeholder</label
                        >
                        <select
                            id="stakeholder-select"
                            bind:value={newAssignment.stakeholder_id}
                            class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-2xl focus:ring-2 focus:ring-primary/20 outline-none transition-all appearance-none cursor-pointer"
                        >
                            <option value={0}>-- Select from list --</option>
                            {#each allStakeholders as s}
                                <option value={s.id}
                                    >{s.name} ({s.organization})</option
                                >
                            {/each}
                        </select>
                    </div>

                    <select
                        id="project-role-select"
                        bind:value={newAssignment.role_id}
                        class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-2xl focus:ring-2 focus:ring-primary/20 outline-none transition-all"
                    >
                        <option value={undefined}>-- Select Role --</option>
                        {#each stakeholderRoles as r}
                            <option value={r.id}>{r.name}</option>
                        {/each}
                    </select>
                </div>

                <div class="flex gap-3 mt-8">
                    <button
                        onclick={() => (isAssignModalOpen = false)}
                        class="flex-1 py-3 rounded-2xl font-bold text-slate-500 hover:bg-slate-100 transition-all"
                    >
                        Cancel
                    </button>
                    <button
                        onclick={handleAssign}
                        disabled={!newAssignment.stakeholder_id}
                        class="flex-1 py-3 rounded-2xl font-bold bg-primary text-white shadow-lg shadow-primary/20 hover:scale-[1.02] active:scale-[0.98] transition-all disabled:opacity-50"
                    >
                        Assign to Project
                    </button>
                </div>
            </div>
        </div>
    </div>
{/if}

<ConfirmDialog
    show={isConfirmOpen}
    title="Remove Stakeholder"
    message="Are you sure you want to remove this stakeholder from the project?"
    confirmText="Confirm Remove"
    variant="danger"
    onConfirm={handleUnassign}
    onCancel={() => {
        isConfirmOpen = false;
        stakeholderToUnassign = null;
    }}
/>
