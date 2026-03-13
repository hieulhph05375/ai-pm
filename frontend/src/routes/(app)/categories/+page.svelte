<script lang="ts">
    import { onMount } from "svelte";
    import {
        categoryService,
        type CategoryType,
        type Category,
    } from "$lib/services/categories";
    import { toast } from "$lib/stores/toast";
    import ContentHeader from "$lib/components/ui/ContentHeader.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import DataTable from "$lib/components/ui/DataTable.svelte";
    import EmptyState from "$lib/components/ui/EmptyState.svelte";
    import ConfirmDialog from "$lib/components/ui/ConfirmDialog.svelte";
    import Badge from "$lib/components/ui/Badge.svelte";
    import CategoryTypeModal from "$lib/components/categories/CategoryTypeModal.svelte";
    import CategoryModal from "$lib/components/categories/CategoryModal.svelte";
    import { hasPermission } from "$lib/utils/permission";
    import { authStore } from "$lib/services/auth";

    let activeTab = $state<"types" | "categories">("types");
    let categoryTypes = $state<CategoryType[]>([]);
    let categories = $state<Category[]>([]);
    let isLoading = $state(true);

    // Pagination & Filter State
    let typePage = $state(1);
    let typeTotal = $state(0);
    let typeSearch = $state("");

    let catPage = $state(1);
    let catTotal = $state(0);
    let catSearch = $state("");
    let catTypeFilter = $state<number | undefined>(undefined);
    let limit = $state(10);

    // Modal State
    let isTypeModalOpen = $state(false);
    let editingType = $state<Partial<CategoryType>>({});

    let isCatModalOpen = $state(false);
    let editingCat = $state<Partial<Category>>({});

    // Confirm Dialog State
    let isConfirmOpen = $state(false);
    let deletionTarget = $state<{
        id: number;
        type: "type" | "category";
    } | null>(null);

    async function fetchData() {
        isLoading = true;
        try {
            if (activeTab === "types") {
                const res = await categoryService.listTypes(
                    typePage,
                    limit,
                    typeSearch,
                );
                categoryTypes = res.data;
                typeTotal = res.total;
            } else {
                const [catRes, typeRes] = await Promise.all([
                    categoryService.listCategories(
                        catPage,
                        limit,
                        catSearch,
                        catTypeFilter,
                    ),
                    // We need all types for the modal dropdown, so we fetch without pagination limit here
                    categoryService.listTypes(1, 100, ""),
                ]);
                categories = catRes.data;
                catTotal = catRes.total;
                categoryTypes = typeRes.data;
            }
        } catch (error: any) {
            toast.error("Unable to load data");
        } finally {
            isLoading = false;
        }
    }

    let searchTimeout: ReturnType<typeof setTimeout>;
    function handleSearch(tab: "types" | "categories", value: string) {
        if (tab === "types") {
            typeSearch = value;
            typePage = 1;
        } else {
            catSearch = value;
            catPage = 1;
        }
        clearTimeout(searchTimeout);
        searchTimeout = setTimeout(() => {
            fetchData();
        }, 300);
    }

    function openTypeModal(item?: CategoryType) {
        editingType = item ? { ...item } : { is_active: true };
        isTypeModalOpen = true;
    }

    function openCatModal(item?: Category) {
        editingCat = item
            ? { ...item }
            : { is_active: true, color: "#3B82F6", icon: "label" };
        isCatModalOpen = true;
    }

    function confirmDelete(id: number, type: "type" | "category") {
        deletionTarget = { id, type };
        isConfirmOpen = true;
    }

    async function handleDelete() {
        if (!deletionTarget) return;
        try {
            if (deletionTarget.type === "type") {
                await categoryService.deleteType(deletionTarget.id);
            } else {
                await categoryService.deleteCategory(deletionTarget.id);
            }
            toast.success("Deleted successfully");
            fetchData();
        } catch (error: any) {
            toast.error("Unable to delete record");
        } finally {
            isConfirmOpen = false;
            deletionTarget = null;
        }
    }

    onMount(fetchData);

    const typeColumns = [
        { key: "name", label: "Type Name" },
        { key: "code", label: "Code" },
        { key: "description", label: "Description" },
        { key: "status", label: "Status" },
        {
            key: "actions",
            label: "Actions",
            align: "right" as const,
            class: "w-24",
        },
    ];

    const catColumns = [
        { key: "name", label: "Category Name" },
        { key: "type", label: "Type" },
        { key: "parent", label: "Parent Category" },
        { key: "appearance", label: "Appearance" },
        { key: "status", label: "Status" },
        {
            key: "actions",
            label: "Actions",
            align: "right" as const,
            class: "w-24",
        },
    ];
</script>

<ContentHeader
    title="Category Configuration"
    subtitle="Standardize taxonomy across your projects. Manage labels, hierarchy, and visual markers for tasks, risks, and documentation."
>
    <div class="flex items-center gap-2">
        <div class="relative max-w-sm w-64 mr-2">
            <span
                class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px]"
            >
                search
            </span>
            <input
                type="text"
                placeholder="Search..."
                value={activeTab === "types" ? typeSearch : catSearch}
                oninput={(e) => handleSearch(activeTab, e.currentTarget.value)}
                class="w-full pl-10 pr-4 py-2 bg-white border border-slate-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all shadow-sm"
            />
        </div>
        {#if activeTab === "categories"}
            <div class="relative w-48 mr-2">
                <select
                    value={catTypeFilter}
                    onchange={(e) => {
                        catTypeFilter = e.currentTarget.value
                            ? parseInt(e.currentTarget.value)
                            : undefined;
                        catPage = 1;
                        fetchData();
                    }}
                    class="w-full pl-3 pr-8 py-2 bg-white border border-slate-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all shadow-sm appearance-none cursor-pointer"
                >
                    <option value="">All Types</option>
                    {#each categoryTypes as type}
                        <option value={type.id}>{type.name}</option>
                    {/each}
                </select>
                <span
                    class="material-symbols-outlined absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 text-[20px] pointer-events-none"
                >
                    expand_more
                </span>
            </div>
        {/if}
        {#if activeTab === "types"}
            {#if hasPermission($authStore.user, $authStore.token, "category:create")}
                <Button icon="add" onclick={() => openTypeModal()}
                    >Add Category Type</Button
                >
            {/if}
        {:else if hasPermission($authStore.user, $authStore.token, "category:create")}
            <Button icon="add" onclick={() => openCatModal()}
                >Add Category</Button
            >
        {/if}
    </div>
</ContentHeader>

<div
    class="bg-white dark:bg-slate-900 rounded-3xl border border-slate-200 dark:border-slate-800 flex flex-col h-[calc(100vh-215px)] overflow-hidden shadow-sm"
>
    <!-- Tabs Header -->
    <div class="flex border-b border-slate-100 dark:border-slate-800 px-6 pt-2">
        <button
            class="px-6 py-4 font-bold text-sm transition-all relative {activeTab ===
            'types'
                ? 'text-primary'
                : 'text-slate-400 hover:text-slate-600'}"
            onclick={() => {
                activeTab = "types";
                fetchData();
            }}
        >
            Category Types
            {#if activeTab === "types"}
                <div
                    class="absolute bottom-0 left-0 right-0 h-1 bg-primary rounded-t-full"
                ></div>
            {/if}
        </button>
        <button
            class="px-6 py-4 font-bold text-sm transition-all relative {activeTab ===
            'categories'
                ? 'text-primary'
                : 'text-slate-400 hover:text-slate-600'}"
            onclick={() => {
                activeTab = "categories";
                fetchData();
            }}
        >
            Categories
            {#if activeTab === "categories"}
                <div
                    class="absolute bottom-0 left-0 right-0 h-1 bg-primary rounded-t-full"
                ></div>
            {/if}
        </button>
    </div>

    <!-- Content Area -->
    <div class="flex-1 overflow-hidden p-6">
        {#if activeTab === "types"}
            <DataTable
                items={categoryTypes}
                columns={typeColumns}
                loading={isLoading}
                page={typePage}
                {limit}
                total={typeTotal}
                onPageChange={(p) => {
                    typePage = p;
                    fetchData();
                }}
                onLimitChange={(l) => {
                    limit = l;
                    typePage = 1;
                    fetchData();
                }}
            >
                {#snippet emptyState()}
                    <EmptyState
                        icon="category"
                        title="No category types found"
                        message="Start by creating a new category type for classification."
                        actionLabel="Add Category Type"
                        onaction={() => openTypeModal()}
                    />
                {/snippet}
                {#snippet rowCell({ item, column })}
                    {#if column.key === "name"}
                        <span class="font-bold text-slate-900">{item.name}</span
                        >
                    {:else if column.key === "code"}
                        <code
                            class="px-2 py-1 bg-slate-100 text-slate-600 rounded text-xs font-mono"
                            >{item.code}</code
                        >
                    {:else if column.key === "description"}
                        <span class="text-sm text-slate-500 line-clamp-1"
                            >{item.description || "—"}</span
                        >
                    {:else if column.key === "status"}
                        <Badge variant={item.is_active ? "emerald" : "slate"}
                            >{item.is_active ? "Active" : "Inactive"}</Badge
                        >
                    {:else if column.key === "actions"}
                        <div class="flex justify-end gap-1">
                            {#if hasPermission($authStore.user, $authStore.token, "category:update")}
                                <button
                                    onclick={() => openTypeModal(item)}
                                    class="p-2 text-slate-400 hover:text-primary rounded-lg transition-all"
                                >
                                    <span
                                        class="material-symbols-outlined text-[18px]"
                                        >edit</span
                                    >
                                </button>
                            {/if}
                            {#if hasPermission($authStore.user, $authStore.token, "category:delete")}
                                <button
                                    onclick={() =>
                                        confirmDelete(item.id, "type")}
                                    class="p-2 text-slate-400 hover:text-rose-600 rounded-lg transition-all"
                                >
                                    <span
                                        class="material-symbols-outlined text-[18px]"
                                        >delete</span
                                    >
                                </button>
                            {/if}
                        </div>
                    {/if}
                {/snippet}
            </DataTable>
        {:else}
            <DataTable
                items={categories}
                columns={catColumns}
                loading={isLoading}
                page={catPage}
                {limit}
                total={catTotal}
                onPageChange={(p) => {
                    catPage = p;
                    fetchData();
                }}
                onLimitChange={(l) => {
                    limit = l;
                    catPage = 1;
                    fetchData();
                }}
            >
                {#snippet emptyState()}
                    <EmptyState
                        icon="sell"
                        title="No categories found"
                        message="Create a category to start classifying data in the system."
                        actionLabel="Add Category"
                        onaction={() => openCatModal()}
                    />
                {/snippet}
                {#snippet rowCell({ item, column })}
                    {#if column.key === "name"}
                        <div
                            class="flex items-center gap-3 {item.parent_id
                                ? 'ml-8'
                                : ''}"
                        >
                            {#if item.parent_id}
                                <span
                                    class="material-symbols-outlined text-slate-400 text-[20px]"
                                >
                                    subdirectory_arrow_right
                                </span>
                            {/if}
                            <div
                                class="w-8 h-8 rounded-lg flex items-center justify-center text-white shrink-0"
                                style="background-color: {item.color ||
                                    '#64748B'}"
                            >
                                <span
                                    class="material-symbols-outlined text-[18px]"
                                >
                                    {item.icon || "sell"}
                                </span>
                            </div>
                            <span class="font-medium text-slate-900 text-sm"
                                >{item.name}</span
                            >
                        </div>
                    {:else if column.key === "type"}
                        <Badge variant="emerald"
                            >{item.type?.name || "N/A"}</Badge
                        >
                    {:else if column.key === "parent"}
                        <span class="text-sm text-slate-600">
                            {item.parent?.name || "—"}
                        </span>
                    {:else if column.key === "appearance"}
                        <div class="flex items-center gap-2">
                            <div
                                class="w-6 h-6 rounded-md flex items-center justify-center text-white"
                                style="background-color: {item.color}"
                            >
                                <span
                                    class="material-symbols-outlined text-[16px]"
                                    >{item.icon}</span
                                >
                            </div>
                            <span class="text-xs text-slate-500"
                                >{item.color}</span
                            >
                        </div>
                    {:else if column.key === "status"}
                        <Badge variant={item.is_active ? "emerald" : "slate"}
                            >{item.is_active ? "Active" : "Inactive"}</Badge
                        >
                    {:else if column.key === "actions"}
                        <div class="flex justify-end gap-1">
                            {#if hasPermission($authStore.user, $authStore.token, "category:update")}
                                <button
                                    onclick={() => openCatModal(item)}
                                    class="p-2 text-slate-400 hover:text-primary rounded-lg transition-all"
                                >
                                    <span
                                        class="material-symbols-outlined text-[18px]"
                                        >edit</span
                                    >
                                </button>
                            {/if}
                            {#if hasPermission($authStore.user, $authStore.token, "category:delete")}
                                <button
                                    onclick={() =>
                                        confirmDelete(item.id, "category")}
                                    class="p-2 text-slate-400 hover:text-rose-600 rounded-lg transition-all"
                                >
                                    <span
                                        class="material-symbols-outlined text-[18px]"
                                        >delete</span
                                    >
                                </button>
                            {/if}
                        </div>
                    {/if}
                {/snippet}
            </DataTable>
        {/if}
    </div>
</div>

<CategoryTypeModal
    show={isTypeModalOpen}
    categoryType={editingType}
    onSave={() => {
        isTypeModalOpen = false;
        fetchData();
    }}
    onCancel={() => (isTypeModalOpen = false)}
/>

<CategoryModal
    show={isCatModalOpen}
    category={editingCat}
    {categoryTypes}
    {categories}
    onSave={() => {
        isCatModalOpen = false;
        fetchData();
    }}
    onCancel={() => (isCatModalOpen = false)}
/>

<ConfirmDialog
    show={isConfirmOpen}
    title="Confirm Deletion"
    message="Are you sure you want to delete this record? This action cannot be undone."
    confirmText="Delete Record"
    variant="danger"
    onConfirm={handleDelete}
    onCancel={() => (isConfirmOpen = false)}
    cancelText="Cancel"
/>
