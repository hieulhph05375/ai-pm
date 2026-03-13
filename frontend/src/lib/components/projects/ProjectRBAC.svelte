<script lang="ts">
    import { onMount } from "svelte";
    import {
        projectRolesService,
        type ProjectRole,
        type ProjectPermission,
    } from "$lib/services/projectRoles";
    import { toast } from "$lib/stores/toast";
    import Button from "$lib/components/ui/Button.svelte";
    import Badge from "$lib/components/ui/Badge.svelte";
    import Modal from "$lib/components/ui/Modal.svelte";
    import DataTable from "$lib/components/ui/DataTable.svelte";
    import { authStore } from "$lib/services/auth";
    import { hasProjectPermission } from "$lib/utils/permission";
    import type { ProjectMember } from "$lib/services/projectMembers";

    let { projectId, projectMembers = [] } = $props<{
        projectId: number;
        projectMembers?: ProjectMember[];
    }>();

    let roles: ProjectRole[] = $state([]);
    let allPermissions: ProjectPermission[] = $state([]);
    let loading = $state(true);
    let showRoleModal = $state(false);
    let showPermissionModal = $state(false);
    let editingRole = $state<Partial<ProjectRole> | null>(null);
    let selectedRole = $state<ProjectRole | null>(null);
    let rolePermissionIds = $state<number[]>([]);

    onMount(async () => {
        await loadData();
    });

    async function loadData() {
        loading = true;
        try {
            const [rRes, pRes] = await Promise.all([
                projectRolesService.getRoles(projectId),
                projectRolesService.getAllProjectPermissions(projectId),
            ]);
            roles = rRes || [];
            allPermissions = pRes || [];
        } catch (e: any) {
            toast.error("Could not load roles or permissions");
        } finally {
            loading = false;
        }
    }

    function openCreateRole() {
        editingRole = {
            name: "",
            description: "",
            color: "#3b82f6",
            is_default: false,
        };
        showRoleModal = true;
    }

    function openEditRole(role: ProjectRole) {
        editingRole = { ...role };
        showRoleModal = true;
    }

    async function handleSaveRole() {
        if (!editingRole?.name) return;

        try {
            if (editingRole.id) {
                await projectRolesService.updateRole(
                    projectId,
                    editingRole.id,
                    editingRole,
                );
                toast.success("Role updated");
            } else {
                await projectRolesService.createRole(projectId, editingRole);
                toast.success("Role created");
            }
            showRoleModal = false;
            await loadData();
        } catch (e: any) {
            toast.error(e.message || "Could not save role");
        }
    }

    async function openManagePermissions(role: ProjectRole) {
        selectedRole = role;
        try {
            const perms = await projectRolesService.getRolePermissions(
                projectId,
                role.id,
            );
            rolePermissionIds = perms.map((p) => p.id);
            showPermissionModal = true;
        } catch (e: any) {
            toast.error("Could not load role permissions");
        }
    }

    async function handleSavePermissions() {
        if (!selectedRole) return;

        try {
            await projectRolesService.setRolePermissions(
                projectId,
                selectedRole.id,
                rolePermissionIds,
            );
            toast.success("Permissions updated");
            showPermissionModal = false;
            await loadData();
        } catch (e: any) {
            toast.error(e.message || "Could not save permissions");
        }
    }

    async function handleDeleteRole(roleId: number) {
        if (
            !confirm(
                "Are you sure you want to delete this role? Members assigned to this role will need to be reassigned.",
            )
        )
            return;

        try {
            await projectRolesService.deleteRole(projectId, roleId);
            toast.success("Role deleted");
            await loadData();
        } catch (e: any) {
            toast.error(e.message || "Could not delete role");
        }
    }

    const columns = [
        { key: "name", label: "Role Name" },
        { key: "description", label: "Description" },
        { key: "is_default", label: "Default", class: "w-24" },
        { key: "actions", label: "", align: "right" as const, class: "w-48" },
    ];

    // Group permissions by module
    let permissionsByModule = $derived(
        allPermissions.reduce(
            (acc, p) => {
                if (!acc[p.module]) acc[p.module] = [];
                acc[p.module].push(p);
                return acc;
            },
            {} as Record<string, ProjectPermission[]>,
        ),
    );
</script>

<div class="space-y-6">
    <div class="flex justify-between items-center">
        <div>
            <h3 class="text-xl font-outfit font-bold text-slate-900">
                Project Roles & Permissions
            </h3>
            <p class="text-sm text-slate-500 mt-1">
                Define custom roles for this project and assign granular
                permissions.
            </p>
        </div>
        {#if hasProjectPermission($authStore.user, $authStore.token, projectMembers, roles, "project:roles:create")}
            <Button
                onclick={openCreateRole}
                variant="primary"
                class="shadow-lg shadow-primary/20"
            >
                <span class="material-symbols-outlined text-[20px] mr-2"
                    >add_moderator</span
                >
                Create Role
            </Button>
        {/if}
    </div>

    <div
        class="bg-white rounded-3xl border border-slate-200 overflow-hidden shadow-sm"
    >
        <DataTable {columns} items={roles} {loading}>
            {#snippet rowCell({ item, column })}
                {#if column.key === "name"}
                    <div class="flex items-center gap-3">
                        <div
                            class="size-3 rounded-full"
                            style="background-color: {item.color}"
                        ></div>
                        <span class="font-bold text-slate-900">{item.name}</span
                        >
                    </div>
                {:else if column.key === "description"}
                    <span class="text-sm text-slate-500"
                        >{item.description || "No description"}</span
                    >
                {:else if column.key === "is_default"}
                    {#if item.is_default}
                        <Badge variant="emerald">Yes</Badge>
                    {:else}
                        <span class="text-slate-300">---</span>
                    {/if}
                {:else if column.key === "actions"}
                    <div class="flex justify-end gap-2">
                        {#if hasProjectPermission($authStore.user, $authStore.token, projectMembers, roles, "project:roles:update")}
                            <button
                                onclick={() => openManagePermissions(item)}
                                class="p-2 text-slate-400 hover:text-primary hover:bg-primary/5 rounded-lg transition-all"
                                title="Manage Permissions"
                            >
                                <span
                                    class="material-symbols-outlined text-[20px]"
                                    >rule</span
                                >
                            </button>
                            <button
                                onclick={() => openEditRole(item)}
                                class="p-2 text-slate-400 hover:text-slate-600 hover:bg-slate-50 rounded-lg transition-all"
                                title="Edit Role"
                            >
                                <span
                                    class="material-symbols-outlined text-[20px]"
                                    >edit</span
                                >
                            </button>
                        {/if}
                        {#if !item.is_default && hasProjectPermission($authStore.user, $authStore.token, projectMembers, roles, "project:roles:delete")}
                            <button
                                onclick={() => handleDeleteRole(item.id)}
                                class="p-2 text-slate-400 hover:text-rose-500 hover:bg-rose-50 rounded-lg transition-all"
                                title="Delete Role"
                            >
                                <span
                                    class="material-symbols-outlined text-[20px]"
                                    >delete</span
                                >
                            </button>
                        {/if}
                    </div>
                {/if}
            {/snippet}
        </DataTable>
    </div>
</div>

<!-- Role Edit Modal -->
<Modal
    bind:show={showRoleModal}
    title={editingRole?.id ? "Edit Role" : "Create Role"}
>
    {#if editingRole}
        <div class="space-y-6">
            <div class="space-y-2">
                <label
                    for="role-name"
                    class="text-[11px] font-bold text-slate-400 uppercase tracking-widest"
                    >Role Name</label
                >
                <input
                    id="role-name"
                    type="text"
                    bind:value={editingRole.name}
                    placeholder="e.g. Lead Developer, Stakeholder"
                    class="w-full bg-slate-50 border border-slate-200 rounded-2xl px-4 py-3 text-sm focus:ring-2 focus:ring-primary/20 transition-all outline-none font-medium"
                />
            </div>
            <div class="space-y-2">
                <label
                    for="role-desc"
                    class="text-[11px] font-bold text-slate-400 uppercase tracking-widest"
                    >Description</label
                >
                <textarea
                    id="role-desc"
                    bind:value={editingRole.description}
                    placeholder="What can this role do?"
                    rows="2"
                    class="w-full bg-slate-50 border border-slate-200 rounded-2xl px-4 py-3 text-sm focus:ring-2 focus:ring-primary/20 transition-all outline-none font-medium"
                ></textarea>
            </div>
            <div class="grid grid-cols-2 gap-4">
                <div class="space-y-2">
                    <label
                        for="role-color"
                        class="text-[11px] font-bold text-slate-400 uppercase tracking-widest"
                        >Color</label
                    >
                    <div class="flex gap-2">
                        <input
                            id="role-color"
                            type="color"
                            bind:value={editingRole.color}
                            class="size-10 rounded-lg overflow-hidden border-none p-0 bg-transparent cursor-pointer"
                        />
                        <input
                            type="text"
                            bind:value={editingRole.color}
                            class="flex-1 bg-slate-50 border border-slate-200 rounded-lg px-3 py-2 text-xs font-mono uppercase"
                        />
                    </div>
                </div>
                <div class="flex items-center gap-3 pt-6">
                    <input
                        type="checkbox"
                        id="is_default"
                        bind:checked={editingRole.is_default}
                        class="size-4 rounded border-slate-300 text-primary focus:ring-primary/20"
                    />
                    <label
                        for="is_default"
                        class="text-sm font-bold text-slate-700"
                        >Set as Default</label
                    >
                </div>
            </div>
        </div>
    {/if}
    {#snippet footer()}
        <Button variant="outline" onclick={() => (showRoleModal = false)}
            >Cancel</Button
        >
        <Button
            onclick={handleSaveRole}
            variant="primary"
            class="shadow-lg shadow-primary/20">Save Role</Button
        >
    {/snippet}
</Modal>

<!-- Permissions Modal -->
<Modal
    bind:show={showPermissionModal}
    title="Manage Role Permissions"
    maxWidth="max-w-4xl"
>
    {#if selectedRole}
        <div class="space-y-6">
            <div class="flex items-center gap-3 mb-2">
                <div
                    class="size-4 rounded-full"
                    style="background-color: {selectedRole.color}"
                ></div>
                <h4 class="text-lg font-outfit font-bold text-slate-900">
                    {selectedRole.name}
                </h4>
            </div>

            <div class="max-h-[60vh] overflow-y-auto pr-2 space-y-8">
                {#each Object.entries(permissionsByModule) as [module, perms]}
                    <div class="space-y-4">
                        <h5
                            class="text-[11px] font-bold text-slate-400 uppercase tracking-widest flex items-center gap-2"
                        >
                            <span class="size-1.5 rounded-full bg-primary/40"
                            ></span>
                            {module}
                        </h5>
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                            {#each perms as p}
                                <label
                                    class="flex items-start gap-3 p-3 rounded-2xl border {rolePermissionIds.includes(
                                        p.id,
                                    )
                                        ? 'bg-primary/5 border-primary/20'
                                        : 'bg-slate-50/50 border-slate-200'} cursor-pointer hover:border-primary/40 transition-all"
                                >
                                    <input
                                        type="checkbox"
                                        value={p.id}
                                        multiple
                                        bind:group={rolePermissionIds}
                                        class="mt-1 size-4 rounded border-slate-300 text-primary focus:ring-primary/20"
                                    />
                                    <div>
                                        <p
                                            class="text-sm font-bold text-slate-900"
                                        >
                                            {p.description}
                                        </p>
                                        <p
                                            class="text-[10px] font-mono text-slate-400 mt-1 uppercase"
                                        >
                                            {p.name}
                                        </p>
                                    </div>
                                </label>
                            {/each}
                        </div>
                    </div>
                {/each}
            </div>
        </div>
    {/if}
    {#snippet footer()}
        <Button variant="outline" onclick={() => (showPermissionModal = false)}
            >Cancel</Button
        >
        <Button
            onclick={handleSavePermissions}
            variant="primary"
            class="shadow-lg shadow-primary/20">Save Permissions</Button
        >
    {/snippet}
</Modal>
