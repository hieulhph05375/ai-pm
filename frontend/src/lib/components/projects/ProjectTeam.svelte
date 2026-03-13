<script lang="ts">
    import { onMount } from "svelte";
    import {
        projectMembersService,
        type ProjectMember,
    } from "$lib/services/projectMembers";
    import {
        projectRolesService,
        type ProjectRole,
    } from "$lib/services/projectRoles";
    import { userService, type User } from "$lib/services/users";
    import { toast } from "$lib/stores/toast";
    import Button from "$lib/components/ui/Button.svelte";
    import Badge from "$lib/components/ui/Badge.svelte";
    import Modal from "$lib/components/ui/Modal.svelte";
    import EmptyState from "$lib/components/ui/EmptyState.svelte";
    import DataTable from "$lib/components/ui/DataTable.svelte";
    import ConfirmDialog from "$lib/components/ui/ConfirmDialog.svelte";
    import { authStore } from "$lib/services/auth";
    import { hasPermission, hasProjectPermission } from "$lib/utils/permission";

    let { projectId } = $props<{ projectId: number }>();

    let members: ProjectMember[] = $state([]);
    let allUsers: User[] = $state([]);
    let loading = $state(true);
    let addingMember = $state(false);
    // Pagination state
    let page = $state(1);
    let limit = $state(10);
    let total = $state(0);

    // Dialog state
    let showRemoveConfirm = $state(false);
    let userIdToRemove = $state<number | null>(null);

    let roles: ProjectRole[] = $state([]);
    let selectedUserId = $state<number | null>(null);
    let selectedRoleId = $state<number | null>(null);

    onMount(async () => {
        await loadData();
    });

    async function loadData() {
        loading = true;
        try {
            const [mRes, uRes, rRes] = await Promise.all([
                projectMembersService.getMembers(projectId, page, limit),
                userService.getUsers("", 1, 1000),
                projectRolesService.getRoles(projectId),
            ]);
            members = mRes.data || [];
            total = mRes.total || 0;
            allUsers = uRes.data || [];
            roles = rRes || [];

            // Set default selected role
            if (roles.length > 0 && !selectedRoleId) {
                const defaultRole = roles.find((r) => r.is_default) || roles[0];
                selectedRoleId = defaultRole.id;
            }
        } catch (e: any) {
            toast.error("Could not load team members or roles");
        } finally {
            loading = false;
        }
    }

    async function handleAddMember() {
        if (!selectedUserId || !selectedRoleId) {
            toast.error("Please select a user and a role");
            return;
        }

        try {
            await projectMembersService.addMember(
                projectId,
                selectedUserId,
                selectedRoleId,
            );
            toast.success("Member added successfully");
            addingMember = false;
            selectedUserId = null;
            await loadData();
        } catch (e: any) {
            toast.error(e.message || "Could not add member");
        }
    }

    function confirmRemoveMember(userId: number) {
        userIdToRemove = userId;
        showRemoveConfirm = true;
    }

    async function handleRemoveMember() {
        if (!userIdToRemove) return;

        try {
            await projectMembersService.removeMember(projectId, userIdToRemove);
            toast.success("Member removed successfully");
            showRemoveConfirm = false;
            userIdToRemove = null;
            await loadData();
        } catch (e: any) {
            toast.error(e.message || "Could not remove member");
        }
    }

    async function handleUpdateRole(userId: number, newRoleId: number) {
        try {
            await projectMembersService.updateMemberRole(
                projectId,
                userId,
                newRoleId,
            );
            toast.success("Role updated successfully");
            await loadData();
        } catch (e: any) {
            toast.error(e.message || "Could not update role");
        }
    }

    // Filter out users who are already members
    let availableUsers = $derived(
        allUsers.filter((u) => !members.some((m) => m.user_id === u.id)),
    );
    const columns = [
        { key: "user", label: "Name & Email" },
        { key: "role", label: "Role", class: "w-48" },
        { key: "joined_at", label: "Joined At", class: "w-48" },
        { key: "actions", label: "", align: "right" as const, class: "w-24" },
    ];
</script>

<div class="space-y-6">
    <div class="flex justify-between items-center">
        <h3 class="text-xl font-outfit font-bold text-slate-900">
            Project Team
        </h3>
        {#if hasProjectPermission($authStore.user, $authStore.token, members, roles, "project:team:create")}
            <Button
                onclick={() => (addingMember = true)}
                variant="primary"
                class="shadow-lg shadow-primary/20"
            >
                <span class="material-symbols-outlined text-[20px] mr-2">
                    person_add
                </span>
                Add Member
            </Button>
        {/if}
    </div>

    <div
        class="bg-white rounded-3xl border border-slate-200 overflow-hidden shadow-sm"
    >
        <DataTable
            {columns}
            items={members}
            {loading}
            {total}
            {page}
            {limit}
            onPageChange={(p) => {
                page = p;
                loadData();
            }}
            onLimitChange={(l) => {
                limit = l;
                page = 1;
                loadData();
            }}
        >
            {#snippet emptyState()}
                <EmptyState
                    icon="group_add"
                    title="No Team Members"
                    message="Every great project needs a team. Start by adding colleagues to this project."
                    actionLabel="Add First Member"
                    actionIcon="person_add"
                    onaction={() => (addingMember = true)}
                />
            {/snippet}

            {#snippet rowCell({ item, column })}
                {#if column.key === "user"}
                    <div class="flex items-center gap-3">
                        <div
                            class="size-10 rounded-xl bg-slate-100 flex items-center justify-center text-slate-400 font-bold group-hover:bg-primary/10 group-hover:text-primary transition-colors"
                        >
                            {item.user?.full_name?.[0] || "?"}
                        </div>
                        <div>
                            <p class="font-bold text-slate-900">
                                {item.user?.full_name || "Unknown"}
                            </p>
                            <p class="text-xs text-slate-400">
                                {item.user?.email || "N/A"}
                            </p>
                        </div>
                    </div>
                {:else if column.key === "role"}
                    {#if hasProjectPermission($authStore.user, $authStore.token, members, roles, "project:team:update")}
                        <select
                            value={item.project_role_id}
                            onchange={(e) =>
                                handleUpdateRole(
                                    item.user_id,
                                    parseInt(
                                        (e.target as HTMLSelectElement).value,
                                    ),
                                )}
                            class="bg-transparent border-none text-sm font-bold text-slate-900 focus:ring-0 cursor-pointer hover:text-primary transition-colors"
                        >
                            {#each roles as role}
                                <option value={role.id}>{role.name}</option>
                            {/each}
                        </select>
                    {:else}
                        <Badge
                            variant="primary"
                            style="background-color: {item.role
                                ?.color}20; color: {item.role?.color}"
                            class="border-none font-bold"
                        >
                            {item.role?.name || "Unknown"}
                        </Badge>
                    {/if}
                {:else if column.key === "joined_at"}
                    <span class="text-sm font-medium text-slate-500">
                        {new Date(item.created_at).toLocaleDateString()}
                    </span>
                {:else if column.key === "actions"}
                    {#if hasProjectPermission($authStore.user, $authStore.token, members, roles, "project:team:delete")}
                        <button
                            onclick={() => confirmRemoveMember(item.user_id)}
                            class="p-2 text-slate-400 hover:text-rose-500 hover:bg-rose-50 rounded-lg transition-all"
                            title="Remove member"
                        >
                            <span class="material-symbols-outlined text-[20px]"
                                >person_remove</span
                            >
                        </button>
                    {/if}
                {/if}
            {/snippet}
        </DataTable>
    </div>
</div>

<Modal bind:show={addingMember} title="Add Team Member">
    <div class="space-y-6">
        <div class="space-y-2">
            <label
                for="user-select-modal"
                class="text-[11px] font-bold text-slate-400 uppercase tracking-widest"
                >Select User</label
            >
            <select
                id="user-select-modal"
                bind:value={selectedUserId}
                class="w-full bg-slate-50 border border-slate-200 rounded-2xl px-4 py-3 text-sm focus:ring-2 focus:ring-primary/20 transition-all outline-none font-medium text-slate-700"
            >
                <option value={null}>Choose a teammate to add...</option>
                {#each availableUsers as user}
                    <option value={user.id}
                        >{user.full_name} ({user.email})</option
                    >
                {/each}
            </select>
        </div>

        <div class="space-y-2">
            <label
                for="role-select-modal"
                class="text-[11px] font-bold text-slate-400 uppercase tracking-widest"
                >Project Role</label
            >
            <select
                id="role-select-modal"
                bind:value={selectedRoleId}
                class="w-full bg-slate-50 border border-slate-200 rounded-2xl px-4 py-3 text-sm focus:ring-2 focus:ring-primary/20 transition-all outline-none font-medium text-slate-700"
            >
                {#each roles as role}
                    <option value={role.id}>{role.name}</option>
                {/each}
            </select>
        </div>

        <div
            class="p-4 bg-primary/5 rounded-2xl border border-primary/10 flex gap-3"
        >
            <span class="material-symbols-outlined text-primary text-[20px]"
                >info</span
            >
            <p class="text-xs text-primary/80 leading-relaxed font-medium">
                Assigned members will have access to all project data based on
                their role. You can update roles or remove members at any time.
            </p>
        </div>
    </div>

    {#snippet footer()}
        <Button
            variant="outline"
            onclick={() => (addingMember = false)}
            class="flex-1 justify-center">Cancel</Button
        >
        <Button
            onclick={handleAddMember}
            class="flex-1 justify-center shadow-lg shadow-primary/20"
            >Add to Team</Button
        >
    {/snippet}
</Modal>

<ConfirmDialog
    show={showRemoveConfirm}
    title="Remove Team Member"
    message="Are you sure you want to remove this member from the project? They will lose access to all project data."
    confirmText="Remove Member"
    variant="danger"
    onConfirm={handleRemoveMember}
    onCancel={() => {
        showRemoveConfirm = false;
        userIdToRemove = null;
    }}
/>
