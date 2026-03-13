<script lang="ts">
    import { onMount } from "svelte";
    import {
        rbacService,
        type Role,
        type Permission,
        type RoleWithPermissions,
    } from "$lib/services/rbac";
    import { toast } from "$lib/stores/toast";
    import ContentHeader from "$lib/components/ui/ContentHeader.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import EmptyState from "$lib/components/ui/EmptyState.svelte";
    import Modal from "$lib/components/ui/Modal.svelte";
    import ConfirmDialog from "$lib/components/ui/ConfirmDialog.svelte";
    import { hasPermission } from "$lib/utils/permission";
    import { authStore } from "$lib/services/auth";

    let roles = $state<Role[]>([]);
    let allPermissions = $state<Permission[]>([]);
    let selectedRole = $state<RoleWithPermissions | null>(null);
    let isLoading = $state(true);
    let isLoadingPerms = $state(false);

    // Modal state
    let isRoleModalOpen = $state(false);
    let editingRole = $state<Partial<Role>>({});
    let isEditMode = $state(false);

    // Confirm delete
    let isConfirmOpen = $state(false);
    let deletingRoleId = $state<number | null>(null);

    // Permission assignment
    let selectedPermIds = $state<Set<number>>(new Set());

    async function fetchRoles() {
        isLoading = true;
        try {
            [roles, allPermissions] = await Promise.all([
                rbacService.listRoles(),
                rbacService.listPermissions(),
            ]);
        } catch (e: any) {
            toast.error(e.message || "Failed to load roles");
        } finally {
            isLoading = false;
        }
    }

    async function selectRole(role: Role) {
        isLoadingPerms = true;
        try {
            const roleWP = await rbacService.getRoleWithPermissions(role.id);
            selectedRole = roleWP;
            selectedPermIds = new Set(
                (roleWP?.permissions || []).map((p) => p.id),
            );
        } catch (e: any) {
            toast.error(e.message || "Failed to load role permissions");
        } finally {
            isLoadingPerms = false;
        }
    }

    function openCreateModal() {
        editingRole = {};
        isEditMode = false;
        isRoleModalOpen = true;
    }

    function openEditModal(role: Role) {
        editingRole = { ...role };
        isEditMode = true;
        isRoleModalOpen = true;
    }

    async function handleSaveRole() {
        if (!editingRole.name?.trim()) {
            toast.error("Role name is required");
            return;
        }
        try {
            if (isEditMode && editingRole.id) {
                await rbacService.updateRole(
                    editingRole.id,
                    editingRole.name,
                    editingRole.description || "",
                );
                toast.success("Role updated successfully");
            } else {
                await rbacService.createRole(
                    editingRole.name,
                    editingRole.description || "",
                );
                toast.success("Role created successfully");
            }
            isRoleModalOpen = false;
            await fetchRoles();
        } catch (e: any) {
            toast.error(e.message || "Operation failed");
        }
    }

    function confirmDelete(roleId: number) {
        deletingRoleId = roleId;
        isConfirmOpen = true;
    }

    async function handleDeleteRole() {
        if (!deletingRoleId) return;
        try {
            await rbacService.deleteRole(deletingRoleId);
            toast.success("Role deleted");
            if (selectedRole?.id === deletingRoleId) selectedRole = null;
            await fetchRoles();
        } catch (e: any) {
            toast.error(e.message || "Failed to delete role");
        } finally {
            isConfirmOpen = false;
            deletingRoleId = null;
        }
    }

    function togglePermission(permId: number) {
        if (!selectedRole) return;
        const newSet = new Set(selectedPermIds);
        if (newSet.has(permId)) {
            newSet.delete(permId);
        } else {
            newSet.add(permId);
        }
        selectedPermIds = newSet;
    }

    function selectAllPermissions() {
        if (!selectedRole || allPermissions.length === 0) return;
        selectedPermIds = new Set(allPermissions.map((p) => p.id));
    }

    function unselectAllPermissions() {
        if (!selectedRole) return;
        selectedPermIds = new Set();
    }

    async function savePermissions() {
        if (!selectedRole) return;
        try {
            await rbacService.assignPermissions(
                selectedRole.id,
                Array.from(selectedPermIds),
            );
            toast.success("Permissions updated");
        } catch (e: any) {
            toast.error(e.message || "Failed to save permissions");
        }
    }

    const groupedPermissions = $derived(() => {
        const groups: Record<string, Permission[]> = {};
        for (const p of allPermissions) {
            const group = p.name.includes(":") ? p.name.split(":")[0] : "other";
            if (!groups[group]) groups[group] = [];
            groups[group].push(p);
        }
        return groups;
    });

    onMount(fetchRoles);
</script>

<svelte:head>
    <title>Roles & Permissions | Enterprise PM</title>
</svelte:head>

<div class="page-wrapper">
    <ContentHeader
        title="Roles & Permissions"
        subtitle="Manage system roles and assign fine-grained permissions"
    >
        {#snippet children()}
            {#if hasPermission($authStore.user, $authStore.token, "role:create")}
                <Button variant="primary" onclick={openCreateModal}>
                    <span class="material-symbols-outlined">add</span>
                    New Role
                </Button>
            {/if}
        {/snippet}
    </ContentHeader>

    <div class="rbac-layout">
        <!-- Roles List Panel -->
        <div class="roles-panel">
            <div class="panel-header">
                <h3>Roles</h3>
                <span class="role-count">{roles.length} roles</span>
            </div>

            {#if isLoading}
                <div class="loading-container">
                    <div class="loading-spinner"></div>
                    <p>Loading...</p>
                </div>
            {:else if roles.length === 0}
                <EmptyState
                    icon="shield_person"
                    title="No roles yet"
                    message="Create roles to organize user access levels."
                    actionLabel="Create Role"
                    onaction={openCreateModal}
                />
            {:else}
                <div class="roles-list">
                    {#each roles as role}
                        <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
                        <div
                            class="role-item"
                            class:active={selectedRole?.id === role.id}
                            onclick={() => selectRole(role)}
                            role="button"
                            tabindex="0"
                            onkeydown={(e) =>
                                e.key === "Enter" && selectRole(role)}
                        >
                            <div class="role-item-icon">
                                <span class="material-symbols-outlined"
                                    >shield_person</span
                                >
                            </div>
                            <div class="role-item-content">
                                <div class="role-name">{role.name}</div>
                                {#if role.description}
                                    <div class="role-desc">
                                        {role.description}
                                    </div>
                                {/if}
                            </div>
                            <div class="role-item-actions">
                                {#if hasPermission($authStore.user, $authStore.token, "role:update")}
                                    <button
                                        class="icon-btn"
                                        onclick={(e) => {
                                            e.stopPropagation();
                                            openEditModal(role);
                                        }}
                                        title="Edit"
                                    >
                                        <span class="material-symbols-outlined"
                                            >edit</span
                                        >
                                    </button>
                                {/if}
                                {#if hasPermission($authStore.user, $authStore.token, "role:delete")}
                                    <button
                                        class="icon-btn danger"
                                        onclick={(e) => {
                                            e.stopPropagation();
                                            confirmDelete(role.id);
                                        }}
                                        title="Delete"
                                    >
                                        <span class="material-symbols-outlined"
                                            >delete</span
                                        >
                                    </button>
                                {/if}
                            </div>
                        </div>
                    {/each}
                </div>
            {/if}
        </div>

        <!-- Permission Assignment Panel -->
        <div class="permissions-panel">
            {#if !selectedRole}
                <div class="no-selection">
                    <span class="material-symbols-outlined no-selection-icon"
                        >touch_app</span
                    >
                    <h3>Select a role</h3>
                    <p>
                        Choose a role on the left to view and edit its
                        permissions
                    </p>
                </div>
            {:else if isLoadingPerms}
                <div class="loading-container">
                    <div class="loading-spinner"></div>
                    <p>Loading permissions...</p>
                </div>
            {:else}
                <div class="perm-panel-header">
                    <div>
                        <h3>
                            Permissions: <span class="role-highlight"
                                >{selectedRole.name}</span
                            >
                        </h3>
                        <p class="perm-count">
                            {selectedPermIds.size} / {allPermissions.length} permissions
                            selected
                        </p>
                    </div>
                    <div class="perm-actions">
                        {#if hasPermission($authStore.user, $authStore.token, "role:update")}
                            <Button
                                variant="ghost"
                                onclick={selectAllPermissions}
                                >Select All</Button
                            >
                            <Button
                                variant="ghost"
                                onclick={unselectAllPermissions}
                                >Clear All</Button
                            >
                            <Button variant="primary" onclick={savePermissions}>
                                <span class="material-symbols-outlined"
                                    >save</span
                                >
                                Save Changes
                            </Button>
                        {/if}
                    </div>
                </div>

                {#if allPermissions.length === 0}
                    <EmptyState
                        icon="key"
                        title="No permissions defined"
                        message="System permissions are seeded into the database during initialization."
                    />
                {:else}
                    <div class="permissions-groups">
                        {#each Object.entries(groupedPermissions()) as [group, perms]}
                            <div class="perm-group">
                                <div class="perm-group-header">
                                    <span class="material-symbols-outlined"
                                        >folder</span
                                    >
                                    <span class="perm-group-name">{group}</span>
                                    <span class="perm-group-count"
                                        >{perms.length}</span
                                    >
                                </div>
                                <div class="perm-list">
                                    {#each perms as perm}
                                        <label class="perm-item">
                                            <input
                                                type="checkbox"
                                                checked={selectedPermIds.has(
                                                    perm.id,
                                                )}
                                                onchange={() =>
                                                    togglePermission(perm.id)}
                                                disabled={!hasPermission(
                                                    $authStore.user,
                                                    $authStore.token,
                                                    "role:update",
                                                )}
                                            />
                                            <div class="perm-info">
                                                <span class="perm-name"
                                                    >{perm.name}</span
                                                >
                                                {#if perm.description}
                                                    <span class="perm-desc"
                                                        >{perm.description}</span
                                                    >
                                                {/if}
                                            </div>
                                        </label>
                                    {/each}
                                </div>
                            </div>
                        {/each}
                    </div>
                {/if}
            {/if}
        </div>
    </div>
</div>

<!-- Role Create/Edit Modal -->
<Modal
    bind:show={isRoleModalOpen}
    title={isEditMode ? "Edit Role" : "New Role"}
>
    <div class="modal-form">
        <div class="form-group">
            <label for="role-name"
                >Role Name <span class="required">*</span></label
            >
            <input
                id="role-name"
                type="text"
                bind:value={editingRole.name}
                placeholder="e.g. Project Manager"
                class="form-input"
            />
        </div>
        <div class="form-group">
            <label for="role-desc">Description</label>
            <textarea
                id="role-desc"
                bind:value={editingRole.description}
                placeholder="Describe the purpose of this role..."
                class="form-input"
                rows="3"
            ></textarea>
        </div>
    </div>
    {#snippet footer()}
        <Button
            variant="ghost"
            onclick={() => {
                isRoleModalOpen = false;
            }}>Cancel</Button
        >
        <Button variant="primary" onclick={handleSaveRole}>
            {isEditMode ? "Save Changes" : "Create Role"}
        </Button>
    {/snippet}
</Modal>

<ConfirmDialog
    bind:show={isConfirmOpen}
    title="Delete Role"
    message="Are you sure you want to delete this role? This action cannot be undone."
    confirmText="Delete"
    variant="danger"
    onConfirm={handleDeleteRole}
    onCancel={() => {
        isConfirmOpen = false;
    }}
/>

<style>
    .page-wrapper {
        display: flex;
        flex-direction: column;
        gap: 1.5rem;
        padding: 1.5rem;
        min-height: 0;
        height: 100%;
    }

    .rbac-layout {
        display: grid;
        grid-template-columns: 360px 1fr;
        gap: 1.5rem;
        flex: 1;
        min-height: 500px;
    }

    .roles-panel,
    .permissions-panel {
        background: var(--bg-card, #fff);
        border: 1px solid var(--border-color, #e2e8f0);
        border-radius: 1rem;
        overflow: hidden;
        display: flex;
        flex-direction: column;
    }

    .panel-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 1rem 1.25rem;
        border-bottom: 1px solid var(--border-color, #e2e8f0);
    }

    .panel-header h3 {
        font-size: 0.9375rem;
        font-weight: 600;
        color: #0f172a;
        margin: 0;
    }

    .role-count {
        font-size: 0.75rem;
        color: #64748b;
        background: #f8fafc;
        padding: 0.2rem 0.6rem;
        border-radius: 999px;
        border: 1px solid #e2e8f0;
    }

    .roles-list {
        flex: 1;
        overflow-y: auto;
        padding: 0.5rem;
    }

    .role-item {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        padding: 0.75rem 1rem;
        border-radius: 0.75rem;
        border: 1.5px solid transparent;
        cursor: pointer;
        width: 100%;
        transition: all 0.15s ease;
    }

    .role-item:hover {
        background: #f8fafc;
    }

    .role-item.active {
        background: #eff6ff;
        border-color: #93c5fd;
    }

    .role-item-icon {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 36px;
        height: 36px;
        border-radius: 0.5rem;
        background: #dbeafe;
        color: #2563eb;
        flex-shrink: 0;
    }

    .role-item-content {
        flex: 1;
        min-width: 0;
    }

    .role-name {
        font-size: 0.875rem;
        font-weight: 600;
        color: #0f172a;
    }

    .role-desc {
        font-size: 0.75rem;
        color: #64748b;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .role-item-actions {
        display: flex;
        gap: 0.25rem;
        opacity: 0;
        transition: opacity 0.15s;
    }

    .role-item:hover .role-item-actions {
        opacity: 1;
    }

    .icon-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 28px;
        height: 28px;
        border: none;
        border-radius: 0.375rem;
        background: transparent;
        cursor: pointer;
        color: #94a3b8;
        transition: all 0.15s;
    }

    .icon-btn:hover {
        background: #f1f5f9;
        color: #0f172a;
    }

    .icon-btn.danger:hover {
        background: #fef2f2;
        color: #dc2626;
    }

    .icon-btn .material-symbols-outlined {
        font-size: 16px;
    }

    .no-selection {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 100%;
        gap: 1rem;
        color: #94a3b8;
        text-align: center;
        padding: 2rem;
    }

    .no-selection-icon {
        font-size: 4rem;
        opacity: 0.3;
    }

    .no-selection h3 {
        margin: 0;
        font-size: 1.125rem;
        color: #475569;
    }

    .no-selection p {
        margin: 0;
        font-size: 0.875rem;
    }

    .perm-panel-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 1rem 1.25rem;
        border-bottom: 1px solid #e2e8f0;
    }

    .perm-panel-header h3 {
        font-size: 0.9375rem;
        font-weight: 600;
        color: #0f172a;
        margin: 0;
    }

    .perm-actions {
        display: flex;
        gap: 0.5rem;
        align-items: center;
    }

    .role-highlight {
        color: #2563eb;
    }

    .perm-count {
        font-size: 0.75rem;
        color: #64748b;
        margin: 0.25rem 0 0;
    }

    .permissions-groups {
        flex: 1;
        overflow-y: auto;
        padding: 1rem 1.25rem;
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .perm-group {
        border: 1px solid #e2e8f0;
        border-radius: 0.75rem;
        overflow: hidden;
    }

    .perm-group-header {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.625rem 0.875rem;
        background: #f8fafc;
        border-bottom: 1px solid #e2e8f0;
    }

    .perm-group-header .material-symbols-outlined {
        font-size: 16px;
        color: #94a3b8;
    }

    .perm-group-name {
        font-size: 0.8125rem;
        font-weight: 600;
        color: #475569;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        flex: 1;
    }

    .perm-group-count {
        font-size: 0.75rem;
        color: #94a3b8;
        background: #fff;
        padding: 0.125rem 0.5rem;
        border-radius: 999px;
        border: 1px solid #e2e8f0;
    }

    .perm-list {
        padding: 0.5rem;
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
        gap: 0.5rem;
    }

    .perm-item {
        display: flex;
        align-items: flex-start;
        gap: 0.75rem;
        padding: 0.625rem 0.75rem;
        border-radius: 0.375rem;
        cursor: pointer;
        transition: background 0.15s;
    }

    .perm-item:hover {
        background: #f8fafc;
    }

    .perm-item input[type="checkbox"] {
        margin-top: 2px;
        width: 16px;
        height: 16px;
        accent-color: #2563eb;
        flex-shrink: 0;
        cursor: pointer;
    }

    .perm-info {
        display: flex;
        flex-direction: column;
        gap: 0.125rem;
    }

    .perm-name {
        font-size: 0.8125rem;
        font-weight: 500;
        color: #0f172a;
        font-family: "JetBrains Mono", monospace;
    }

    .perm-desc {
        font-size: 0.75rem;
        color: #64748b;
    }

    .modal-form {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .form-group {
        display: flex;
        flex-direction: column;
        gap: 0.375rem;
    }

    .form-group label {
        font-size: 0.875rem;
        font-weight: 500;
        color: #475569;
    }

    .required {
        color: #ef4444;
    }

    .form-input {
        width: 100%;
        padding: 0.625rem 0.875rem;
        border: 1.5px solid #e2e8f0;
        border-radius: 0.5rem;
        background: #fff;
        color: #0f172a;
        font-size: 0.875rem;
        transition: border-color 0.15s;
        font-family: inherit;
        resize: vertical;
        box-sizing: border-box;
    }

    .form-input:focus {
        outline: none;
        border-color: #3b82f6;
    }

    .loading-container {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 3rem;
        gap: 1rem;
        color: #94a3b8;
    }

    .loading-spinner {
        width: 32px;
        height: 32px;
        border: 3px solid #e2e8f0;
        border-top-color: #3b82f6;
        border-radius: 50%;
        animation: spin 0.7s linear infinite;
    }

    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }
</style>
