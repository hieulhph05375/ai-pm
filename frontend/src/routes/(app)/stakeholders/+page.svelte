<script lang="ts">
    import { onMount } from "svelte";
    import {
        stakeholderService,
        type Stakeholder,
    } from "$lib/services/stakeholders";
    import { toast } from "$lib/stores/toast";
    import EmptyState from "$lib/components/ui/EmptyState.svelte";
    import DataTable from "$lib/components/ui/DataTable.svelte";
    import ConfirmDialog from "$lib/components/ui/ConfirmDialog.svelte";
    import ContentHeader from "$lib/components/ui/ContentHeader.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Badge from "$lib/components/ui/Badge.svelte";
    import Avatar from "$lib/components/ui/Avatar.svelte";
    import { categoryService, type Category } from "$lib/services/categories";
    import { hasPermission } from "$lib/utils/permission";
    import { authStore } from "$lib/services/auth";

    let stakeholders = $state<Stakeholder[]>([]);
    let isLoading = $state(true);
    let total = $state(0);
    let page = $state(1);
    let limit = $state(10);
    let searchQuery = $state("");
    let isModalOpen = $state(false);
    let stakeholderRoles = $state<Category[]>([]);

    // Confirm Dialog State
    let isConfirmOpen = $state(false);
    let stakeholderToDelete = $state<number | null>(null);
    let editingStakeholder = $state<Partial<Stakeholder>>({
        name: "",
        role: "",
        organization: "",
        email: "",
        phone: "",
        notes: "",
    });
    let errors = $state<Record<string, string>>({});

    async function loadStakeholders() {
        isLoading = true;
        if (stakeholderRoles.length === 0) {
            try {
                const res = await categoryService.listCategories(1, 100, "", 9);
                stakeholderRoles = res.data;
            } catch (e) {
                console.error("Failed to load roles", e);
            }
        }
        try {
            const res = await stakeholderService.list(page, limit, searchQuery);
            stakeholders = res.data;
            total = res.total;
        } catch (error: any) {
            if (!error.isAuthError) {
                toast.error("Could not load stakeholder list");
            }
        } finally {
            isLoading = false;
        }
    }

    function openCreate() {
        editingStakeholder = {
            name: "",
            role: "",
            organization: "",
            email: "",
            phone: "",
            notes: "",
        };
        errors = {};
        isModalOpen = true;
    }

    function openEdit(s: Stakeholder) {
        editingStakeholder = { ...s };
        errors = {};
        isModalOpen = true;
    }

    function validate() {
        errors = {};
        if (!editingStakeholder.name?.trim()) {
            errors.name = "Please enter full name";
        }
        if (
            editingStakeholder.email &&
            !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(editingStakeholder.email)
        ) {
            errors.email = "Invalid email";
        }
        return Object.keys(errors).length === 0;
    }

    async function handleSave() {
        if (!validate()) return;

        // Sync role name for backward compatibility
        const selectedRole = stakeholderRoles.find(
            (c) => c.id === editingStakeholder.role_id,
        );
        if (selectedRole) {
            editingStakeholder.role = selectedRole.name;
        }

        try {
            if (editingStakeholder.id) {
                await stakeholderService.update(
                    editingStakeholder.id,
                    editingStakeholder,
                );
                toast.success("Updated successfully");
            } else {
                await stakeholderService.create(editingStakeholder);
                toast.success("Added successfully");
            }
            isModalOpen = false;
            loadStakeholders();
        } catch (error: any) {
            if (!error.isAuthError) {
                toast.error(error.message || "Error while saving");
            }
        }
    }

    function confirmDelete(id: number) {
        stakeholderToDelete = id;
        isConfirmOpen = true;
    }

    async function handleDelete() {
        if (!stakeholderToDelete) return;
        try {
            await stakeholderService.delete(stakeholderToDelete);
            toast.success("Deleted successfully");
            isConfirmOpen = false;
            stakeholderToDelete = null;
            loadStakeholders();
        } catch (error: any) {
            if (!error.isAuthError) {
                toast.error("Could not delete stakeholder");
            }
        }
    }

    let searchTimeout: any;
    function handleSearch() {
        clearTimeout(searchTimeout);
        searchTimeout = setTimeout(() => {
            page = 1;
            loadStakeholders();
        }, 300);
    }

    const columns = [
        { key: "name", label: "Name" },
        { key: "role", label: "Role / Organization" },
        { key: "contact", label: "Contact" },
        {
            key: "actions",
            label: "Action",
            align: "right" as const,
            class: "w-24",
        },
    ];

    onMount(loadStakeholders);
</script>

<ContentHeader
    title="Stakeholder Management"
    subtitle="Manage stakeholders across the organization"
>
    <div class="flex items-center gap-3">
        <div
            class="flex bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl px-4 py-2 w-64 items-center gap-2 focus-within:ring-2 focus-within:ring-primary/20 focus-within:border-primary transition-all"
        >
            <span class="material-symbols-outlined text-slate-400 text-sm"
                >search</span
            >
            <input
                class="bg-transparent border-none focus:ring-0 text-sm w-full outline-none placeholder:text-slate-400 text-slate-900 dark:text-white"
                placeholder="Search stakeholders..."
                type="text"
                bind:value={searchQuery}
                oninput={handleSearch}
            />
        </div>
        {#if hasPermission($authStore.user, $authStore.token, "stakeholder:create")}
            <Button icon="add" onclick={openCreate}>Add Stakeholder</Button>
        {/if}
    </div>
</ContentHeader>

<div class="space-y-4">
    <!-- List -->
    <DataTable
        items={stakeholders}
        {columns}
        loading={isLoading}
        {total}
        {page}
        {limit}
        onPageChange={(p) => {
            page = p;
            loadStakeholders();
        }}
        onLimitChange={(l) => {
            limit = l;
            page = 1;
            loadStakeholders();
        }}
    >
        {#snippet emptyState()}
            <EmptyState
                icon="groups"
                actionIcon="person_add"
                title="No stakeholders yet"
                message="Start by adding a new stakeholder to the system."
                actionLabel="Add New Stakeholder"
                onaction={openCreate}
            />
        {/snippet}
        {#snippet rowCell({ item, column })}
            {#if column.key === "name"}
                <div class="flex items-center gap-3">
                    <Avatar name={item.name} size="md" />
                    <div>
                        <p
                            class="font-bold text-slate-900 text-sm leading-tight"
                        >
                            {item.name}
                        </p>
                        {#if item.notes}
                            <p
                                class="text-[11px] text-slate-400 font-medium mt-0.5 truncate max-w-[200px]"
                            >
                                {item.notes}
                            </p>
                        {/if}
                    </div>
                </div>
            {:else if column.key === "role"}
                <div>
                    <div class="text-sm font-medium text-slate-700">
                        {item.role_cat?.name || item.role || "N/A"}
                    </div>
                    <div class="text-xs text-slate-400 mt-0.5">
                        {item.organization || "Internal"}
                    </div>
                </div>
            {:else if column.key === "contact"}
                <div class="space-y-1">
                    <div
                        class="text-sm text-slate-600 flex items-center gap-1.5 font-medium"
                    >
                        <span
                            class="material-symbols-outlined text-[16px] text-slate-400"
                            >mail</span
                        >
                        {item.email || "-"}
                    </div>
                    <div
                        class="text-xs text-slate-400 flex items-center gap-1.5"
                    >
                        <span
                            class="material-symbols-outlined text-[16px] text-slate-400"
                            >phone</span
                        >
                        {item.phone || "-"}
                    </div>
                </div>
            {:else if column.key === "actions"}
                <div class="flex justify-end gap-1">
                    {#if hasPermission($authStore.user, $authStore.token, "stakeholder:update")}
                        <button
                            onclick={() => openEdit(item)}
                            class="p-2 text-slate-400 hover:text-primary hover:bg-primary/5 rounded-lg transition-all"
                        >
                            <span class="material-symbols-outlined text-[18px]"
                                >edit</span
                            >
                        </button>
                    {/if}
                    {#if hasPermission($authStore.user, $authStore.token, "stakeholder:delete")}
                        <button
                            onclick={() => confirmDelete(item.id!)}
                            class="p-2 text-slate-400 hover:text-rose-600 hover:bg-rose-50 rounded-lg transition-all"
                        >
                            <span class="material-symbols-outlined text-[18px]"
                                >delete</span
                            >
                        </button>
                    {/if}
                </div>
            {/if}
        {/snippet}
    </DataTable>
</div>

<!-- Modal -->
{#if isModalOpen}
    <div class="fixed inset-0 z-[100] flex items-center justify-center p-4">
        <button
            class="absolute inset-0 bg-slate-900/40 backdrop-blur-sm border-none w-full h-full cursor-default"
            onclick={() => (isModalOpen = false)}
            aria-label="Close modal"
        ></button>
        <div
            class="bg-white rounded-2xl w-full max-w-lg relative shadow-2xl animate-in fade-in zoom-in duration-200 overflow-hidden"
        >
            <div
                class="px-6 py-4 border-b border-slate-100 flex justify-between items-center bg-slate-50/50"
            >
                <h3 class="font-display font-bold text-lg text-slate-900">
                    {editingStakeholder.id
                        ? "Update Stakeholder"
                        : "Add New Stakeholder"}
                </h3>
                <button
                    onclick={() => (isModalOpen = false)}
                    class="text-slate-400 hover:text-slate-600"
                >
                    <span class="material-symbols-outlined">close</span>
                </button>
            </div>

            <div class="p-6 space-y-4">
                <div class="grid grid-cols-2 gap-4">
                    <div class="col-span-2">
                        <label
                            for="stakeholder-name"
                            class="block text-xs font-bold {errors.name
                                ? 'text-rose-500'
                                : 'text-slate-500'} uppercase mb-1.5 ml-1"
                        >
                            Full Name <span class="text-rose-500">*</span>
                            {#if errors.name}<span
                                    class="text-rose-500 normal-case font-normal italic ml-2"
                                    >{errors.name}</span
                                >{/if}
                        </label>
                        <input
                            id="stakeholder-name"
                            bind:value={editingStakeholder.name}
                            type="text"
                            class="w-full px-4 py-2.5 bg-slate-50 border {errors.name
                                ? 'border-rose-300 focus:ring-rose-500/20'
                                : 'border-slate-200 focus:ring-primary/20'} rounded-xl focus:ring-2 outline-none transition-all"
                            placeholder="e.g., John Doe"
                        />
                    </div>
                    <div>
                        <label
                            for="stakeholder-role"
                            class="block text-xs font-bold text-slate-500 uppercase mb-1.5 ml-1"
                            >Role</label
                        >
                        <select
                            id="stakeholder-role"
                            bind:value={editingStakeholder.role_id}
                            class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-primary/20 outline-none transition-all"
                        >
                            {#each stakeholderRoles as r}
                                <option value={r.id}>{r.name}</option>
                            {/each}
                        </select>
                    </div>
                    <div>
                        <label
                            for="stakeholder-org"
                            class="block text-xs font-bold text-slate-500 uppercase mb-1.5 ml-1"
                            >Organization</label
                        >
                        <input
                            id="stakeholder-org"
                            bind:value={editingStakeholder.organization}
                            type="text"
                            class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-primary/20 outline-none transition-all"
                            placeholder="e.g., IBM / FPT"
                        />
                    </div>
                    <div>
                        <label
                            for="stakeholder-email"
                            class="block text-xs font-bold {errors.email
                                ? 'text-rose-500'
                                : 'text-slate-500'} uppercase mb-1.5 ml-1"
                        >
                            Email
                            {#if errors.email}<span
                                    class="text-rose-500 normal-case font-normal italic ml-2"
                                    >{errors.email}</span
                                >{/if}
                        </label>
                        <input
                            id="stakeholder-email"
                            bind:value={editingStakeholder.email}
                            type="email"
                            class="w-full px-4 py-2.5 bg-slate-50 border {errors.email
                                ? 'border-rose-300 focus:ring-rose-500/20'
                                : 'border-slate-200 focus:ring-primary/20'} rounded-xl focus:ring-2 outline-none transition-all"
                            placeholder="email@example.com"
                        />
                    </div>
                    <div>
                        <label
                            for="stakeholder-phone"
                            class="block text-xs font-bold text-slate-500 uppercase mb-1.5 ml-1"
                            >Phone</label
                        >
                        <input
                            id="stakeholder-phone"
                            bind:value={editingStakeholder.phone}
                            type="text"
                            class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-primary/20 outline-none transition-all"
                            placeholder="09xxx..."
                        />
                    </div>
                    <div class="col-span-2">
                        <label
                            for="stakeholder-notes"
                            class="block text-xs font-bold text-slate-500 uppercase mb-1.5 ml-1"
                            >Notes</label
                        >
                        <textarea
                            id="stakeholder-notes"
                            bind:value={editingStakeholder.notes}
                            rows="2"
                            class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-primary/20 outline-none transition-all"
                            placeholder="Describe interest, influence..."
                        ></textarea>
                    </div>
                </div>
            </div>

            <div
                class="p-6 bg-slate-50/50 border-t border-slate-100 flex justify-end gap-3"
            >
                <button
                    onclick={() => (isModalOpen = false)}
                    class="px-5 py-2.5 rounded-xl font-medium text-slate-600 hover:bg-slate-200 transition-all"
                >
                    Cancel
                </button>
                <button
                    onclick={handleSave}
                    disabled={!editingStakeholder.name}
                    class="px-8 py-2.5 rounded-xl font-medium bg-primary text-white shadow-lg shadow-primary/20 hover:scale-[1.02] active:scale-[0.98] transition-all disabled:opacity-50"
                >
                    Save
                </button>
            </div>
        </div>
    </div>
{/if}

<ConfirmDialog
    show={isConfirmOpen}
    title="Delete Stakeholder"
    message="Are you sure you want to delete this stakeholder? This action cannot be undone."
    confirmText="Confirm Delete"
    variant="danger"
    onConfirm={handleDelete}
    onCancel={() => {
        isConfirmOpen = false;
        stakeholderToDelete = null;
    }}
/>
