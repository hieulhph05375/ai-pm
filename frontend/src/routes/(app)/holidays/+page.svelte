<script lang="ts">
    import { onMount } from "svelte";
    import { holidayService, type Holiday } from "$lib/services/holidays";
    import { settingService } from "$lib/services/settings";
    import { toast } from "$lib/stores/toast";
    import EmptyState from "$lib/components/ui/EmptyState.svelte";
    import ConfirmDialog from "$lib/components/ui/ConfirmDialog.svelte";
    import DataTable from "$lib/components/ui/DataTable.svelte";
    import ContentHeader from "$lib/components/ui/ContentHeader.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Badge from "$lib/components/ui/Badge.svelte";
    import { categoryService, type Category } from "$lib/services/categories";
    import { hasPermission } from "$lib/utils/permission";
    import { authStore } from "$lib/services/auth";

    let holidays = $state<Holiday[]>([]);
    let total = $state(0);
    let page = $state(1);
    let limit = $state(10);
    let isLoading = $state(true);
    let isModalOpen = $state(false);
    let restDays = $state<number[]>([0, 6]);
    let isUpdatingSettings = $state(false);
    let editingHoliday = $state<Partial<Holiday>>({
        name: "",
        date: new Date().toISOString().split("T")[0],
        type: "state",
        is_recurring: false,
    });
    let holidayCategories = $state<Category[]>([]);

    // Confirm Dialog State
    let isConfirmOpen = $state(false);
    let holidayToDelete = $state<number | null>(null);

    const columns = [
        { key: "date", label: "Date", class: "w-40" },
        { key: "name", label: "Holiday Name" },
        { key: "type", label: "Category", class: "w-40" },
        { key: "actions", label: "", align: "right" as const, class: "w-24" },
    ];

    async function loadHolidays() {
        isLoading = true;
        try {
            const res = await holidayService.list(
                undefined,
                undefined,
                page,
                limit,
            );
            holidays = res.data;
            total = res.total;
        } catch (error: any) {
            if (!error.isAuthError) {
                toast.error("Could not load holiday list");
            }
        } finally {
            isLoading = false;
        }
    }

    async function loadSettings() {
        try {
            const settings = await settingService.getAll();
            if (settings.rest_days) {
                restDays = settings.rest_days;
            }
        } catch (error) {
            console.error("Failed to load settings", error);
        }
    }

    async function toggleRestDay(day: number) {
        isUpdatingSettings = true;
        try {
            let newRestDays = [...restDays];
            if (newRestDays.includes(day)) {
                newRestDays = newRestDays.filter((d) => d !== day);
            } else {
                newRestDays.push(day);
            }
            await settingService.update("rest_days", newRestDays);
            restDays = newRestDays;
            toast.success("Settings updated successfully");
        } catch (error: any) {
            if (!error.isAuthError) {
                toast.error("Could not update settings");
            }
        } finally {
            isUpdatingSettings = false;
        }
    }

    function onPageChange(newPage: number) {
        page = newPage;
        loadHolidays();
    }

    function openCreate() {
        editingHoliday = {
            name: "",
            date: new Date().toISOString().split("T")[0],
            type: "state",
            type_id: holidayCategories.find((c) => c.is_active)?.id,
            is_recurring: false,
        };
        isModalOpen = true;
    }

    function openEdit(h: Holiday) {
        editingHoliday = { ...h };
        isModalOpen = true;
    }

    async function handleSave() {
        // Sync type name
        const selectedCat = holidayCategories.find(
            (c) => c.id === editingHoliday.type_id,
        );
        if (selectedCat) {
            editingHoliday.type = selectedCat.name as any;
        }

        try {
            if (editingHoliday.id) {
                await holidayService.update(editingHoliday.id, editingHoliday);
                toast.success("Updated successfully");
            } else {
                await holidayService.create(editingHoliday);
                toast.success("Added successfully");
            }
            isModalOpen = false;
            loadHolidays();
        } catch (error: any) {
            if (!error.isAuthError) {
                toast.error("An error occurred while saving");
            }
        }
    }

    function confirmDelete(id: number) {
        holidayToDelete = id;
        isConfirmOpen = true;
    }

    async function handleDelete() {
        if (!holidayToDelete) return;
        try {
            await holidayService.delete(holidayToDelete);
            toast.success("Deleted successfully");
            isConfirmOpen = false;
            holidayToDelete = null;
            loadHolidays();
        } catch (error: any) {
            if (!error.isAuthError) {
                toast.error("Could not delete");
            }
        }
    }

    onMount(async () => {
        try {
            const res = await categoryService.listCategories(1, 100, "", 10);
            holidayCategories = res.data;
        } catch (e) {
            console.error("Failed to load holiday categories", e);
        }
        loadHolidays();
        loadSettings();
    });
</script>

<ContentHeader
    title="Holiday Calendar System"
    subtitle="Manage State and Company Holidays (2-tier)"
>
    {#if hasPermission($authStore.user, $authStore.token, "holiday:create")}
        <Button icon="add" onclick={openCreate}>Add Holiday</Button>
    {/if}
</ContentHeader>

<div class="space-y-8">
    <!-- Stats -->

    <!-- Stats -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div
            class="bg-white rounded-2xl border shadow-sm p-6 bg-rose-50 border-rose-100"
        >
            <div class="flex items-center gap-4">
                <div
                    class="size-12 rounded-2xl bg-rose-500/10 flex items-center justify-center text-rose-600"
                >
                    <span class="material-symbols-outlined font-bold"
                        >account_balance</span
                    >
                </div>
                <div>
                    <h3 class="text-rose-900 font-bold text-xl">
                        {holidays.filter((h) => h.type === "state").length}
                    </h3>
                    <p
                        class="text-rose-600 text-xs font-bold uppercase tracking-wider"
                    >
                        State Holidays
                    </p>
                </div>
            </div>
        </div>
        <div
            class="bg-white rounded-2xl border shadow-sm p-6 bg-sky-50 border-sky-100"
        >
            <div class="flex items-center gap-4">
                <div
                    class="size-12 rounded-2xl bg-sky-500/10 flex items-center justify-center text-sky-600"
                >
                    <span class="material-symbols-outlined font-bold"
                        >apartment</span
                    >
                </div>
                <div>
                    <h3 class="text-sky-900 font-bold text-xl">
                        {holidays.filter((h) => h.type === "company").length}
                    </h3>
                    <p
                        class="text-sky-600 text-xs font-bold uppercase tracking-wider"
                    >
                        Company Holidays
                    </p>
                </div>
            </div>
        </div>
        <div
            class="bg-white rounded-2xl border border-slate-200 shadow-sm p-6 bg-slate-50 border-slate-100"
        >
            <div class="flex items-center gap-4">
                <div
                    class="size-12 rounded-2xl bg-slate-500/10 flex items-center justify-center text-slate-600"
                >
                    <span class="material-symbols-outlined font-bold"
                        >list_alt</span
                    >
                </div>
                <div>
                    <h3 class="text-slate-900 font-bold text-xl">
                        {holidays.length}
                    </h3>
                    <p
                        class="text-slate-600 text-xs font-bold uppercase tracking-wider"
                    >
                        Total Holidays
                    </p>
                </div>
            </div>
        </div>
    </div>

    <!-- Permanent Non-Working Days info -->
    <div
        class="bg-indigo-50/50 border border-indigo-100 rounded-2xl p-6 mb-8 flex items-center justify-between"
    >
        <div class="flex items-center gap-5">
            <div
                class="size-14 rounded-2xl bg-white shadow-sm flex items-center justify-center text-indigo-600 border border-indigo-100"
            >
                <span class="material-symbols-outlined text-[28px]"
                    >weekend</span
                >
            </div>
            <div>
                <h4 class="font-bold text-slate-900 text-lg">
                    Weekend Configuration
                </h4>
                <p class="text-slate-500 text-sm">
                    Select the weekend days that are automatically recorded as
                    holidays for projects.
                </p>
            </div>
        </div>
        <div class="hidden sm:flex items-center gap-3">
            {#if hasPermission($authStore.user, $authStore.token, "setting:update")}
                <button
                    onclick={() => toggleRestDay(6)}
                    disabled={isUpdatingSettings}
                    class="px-4 py-2 rounded-xl text-sm font-bold uppercase tracking-wider transition-all border-2 flex items-center gap-2 {restDays.includes(
                        6,
                    )
                        ? 'bg-indigo-100 text-indigo-700 border-indigo-200 hover:bg-indigo-200'
                        : 'bg-white text-slate-400 border-slate-200 hover:border-slate-300'} disabled:opacity-50"
                >
                    {#if restDays.includes(6)}
                        <span class="material-symbols-outlined text-[16px]"
                            >check_circle</span
                        >
                    {:else}
                        <span class="material-symbols-outlined text-[16px]"
                            >radio_button_unchecked</span
                        >
                    {/if}
                    Saturday
                </button>
                <button
                    onclick={() => toggleRestDay(0)}
                    disabled={isUpdatingSettings}
                    class="px-4 py-2 rounded-xl text-sm font-bold uppercase tracking-wider transition-all border-2 flex items-center gap-2 {restDays.includes(
                        0,
                    )
                        ? 'bg-indigo-100 text-indigo-700 border-indigo-200 hover:bg-indigo-200'
                        : 'bg-white text-slate-400 border-slate-200 hover:border-slate-300'} disabled:opacity-50"
                >
                    {#if restDays.includes(0)}
                        <span class="material-symbols-outlined text-[16px]"
                            >check_circle</span
                        >
                    {:else}
                        <span class="material-symbols-outlined text-[16px]"
                            >radio_button_unchecked</span
                        >
                    {/if}
                    Sunday
                </button>
            {/if}
        </div>
    </div>

    <!-- List -->
    <DataTable
        items={holidays}
        {columns}
        loading={isLoading}
        {total}
        {page}
        {limit}
        {onPageChange}
        onLimitChange={(l) => {
            limit = l;
            page = 1;
            loadHolidays();
        }}
    >
        {#snippet emptyState()}
            <EmptyState
                icon="calendar_month"
                actionIcon="add"
                title="No holidays yet"
                message="The holiday list is currently empty. Please add state or company holidays."
                actionLabel="Add Holiday"
                onaction={openCreate}
            />
        {/snippet}
        {#snippet rowCell({ item: h, column })}
            {#if column.key === "date"}
                <div class="font-bold text-slate-900">
                    {new Date(h.date).toLocaleDateString("vi-VN")}
                </div>
                {#if h.is_recurring}
                    <span
                        class="inline-flex items-center gap-1 text-[10px] bg-emerald-50 text-emerald-600 px-2 py-0.5 rounded-full font-bold uppercase mt-1"
                    >
                        <span class="material-symbols-outlined text-[10px]"
                            >sync</span
                        >
                        Recurring every year
                    </span>
                {/if}
            {:else if column.key === "name"}
                <div class="text-sm font-medium text-slate-700">
                    {h.name}
                </div>
            {:else if column.key === "type"}
                {@const cat =
                    holidayCategories.find((c) => c.id === h.type_id) ||
                    h.type_cat}
                <Badge
                    color={cat?.color ||
                        (h.type === "state" ? "#f43f5e" : "#6366f1")}
                >
                    <span class="inline-flex items-center gap-1">
                        <span class="material-symbols-outlined text-[14px]">
                            {h.type === "state"
                                ? "account_balance"
                                : "apartment"}
                        </span>
                        <span>
                            {cat?.name ||
                                (h.type === "state" ? "State" : "Company")}
                        </span>
                    </span>
                </Badge>
            {:else if column.key === "actions"}
                <div
                    class="flex justify-end gap-1 opacity-0 group-hover:opacity-100 transition-opacity"
                >
                    {#if hasPermission($authStore.user, $authStore.token, "holiday:update")}
                        <button
                            onclick={() => openEdit(h)}
                            class="p-2 text-slate-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-all"
                            title="Edit"
                        >
                            <span class="material-symbols-outlined text-[18px]"
                                >edit</span
                            >
                        </button>
                    {/if}
                    {#if hasPermission($authStore.user, $authStore.token, "holiday:delete")}
                        <button
                            onclick={() => confirmDelete(h.id!)}
                            class="p-2 text-slate-400 hover:text-rose-600 hover:bg-rose-50 rounded-lg transition-all"
                            title="Delete"
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
            class="absolute inset-0 bg-slate-900/40 backdrop-blur-sm w-full h-full border-none cursor-default"
            onclick={() => (isModalOpen = false)}
            aria-label="Close modal"
        ></button>
        <div
            class="bg-white rounded-2xl w-full max-w-md relative shadow-2xl animate-in fade-in zoom-in duration-200"
        >
            <div
                class="px-6 py-4 border-b border-slate-100 flex justify-between items-center"
            >
                <h3 class="font-display font-bold text-lg text-slate-900">
                    {editingHoliday.id ? "Update Holiday" : "Add New Holiday"}
                </h3>
                <button
                    onclick={() => (isModalOpen = false)}
                    class="text-slate-400 hover:text-slate-600"
                >
                    <span class="material-symbols-outlined">close</span>
                </button>
            </div>

            <div class="p-6 space-y-5">
                <div>
                    <label
                        for="holiday-name"
                        class="block text-xs font-bold text-slate-500 uppercase mb-1.5 ml-1"
                        >Holiday Name</label
                    >
                    <input
                        id="holiday-name"
                        bind:value={editingHoliday.name}
                        type="text"
                        class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500/20 outline-none transition-all"
                        placeholder="e.g., Lunar New Year"
                    />
                </div>
                <div>
                    <label
                        for="holiday-date"
                        class="block text-xs font-bold text-slate-500 uppercase mb-1.5 ml-1"
                        >Date</label
                    >
                    <input
                        id="holiday-date"
                        bind:value={editingHoliday.date}
                        type="date"
                        class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500/20 outline-none transition-all"
                    />
                </div>
                <div>
                    <span
                        class="block text-xs font-bold text-slate-500 uppercase mb-1.5 ml-1"
                        >Category</span
                    >
                    <div class="grid grid-cols-2 gap-3">
                        {#each holidayCategories as cat}
                            <button
                                type="button"
                                onclick={() =>
                                    (editingHoliday.type_id = cat.id)}
                                class="px-4 py-2.5 rounded-xl border-2 transition-all flex items-center justify-center gap-2 {editingHoliday.type_id ===
                                cat.id
                                    ? 'bg-primary/5 border-primary text-primary'
                                    : 'bg-white border-slate-200 text-slate-400 hover:border-slate-300'}"
                                style={editingHoliday.type_id === cat.id
                                    ? `background-color: ${cat.color}15; border-color: ${cat.color}; color: ${cat.color};`
                                    : ""}
                            >
                                <span
                                    class="material-symbols-outlined text-[18px]"
                                    >{cat.name.includes("State")
                                        ? "account_balance"
                                        : "apartment"}</span
                                >
                                <span class="text-sm font-bold">{cat.name}</span
                                >
                            </button>
                        {/each}
                    </div>
                </div>
                <div class="flex items-center gap-3 bg-slate-50 p-3 rounded-xl">
                    <input
                        type="checkbox"
                        bind:checked={editingHoliday.is_recurring}
                        id="recurring"
                        class="size-4 rounded accent-indigo-600"
                    />
                    <label
                        for="recurring"
                        class="text-sm font-medium text-slate-700 cursor-pointer"
                        >Recurring every year</label
                    >
                </div>
            </div>

            <div
                class="p-6 bg-slate-50/50 border-t border-slate-100 flex justify-end gap-3 rounded-b-2xl"
            >
                <button
                    onclick={() => (isModalOpen = false)}
                    class="px-5 py-2.5 rounded-xl font-medium text-slate-600 hover:bg-slate-200 transition-all"
                    >Cancel</button
                >
                <button
                    onclick={handleSave}
                    disabled={!editingHoliday.name}
                    class="px-8 py-2.5 rounded-xl font-medium bg-indigo-600 text-white shadow-lg shadow-indigo-200 hover:bg-indigo-700 transition-all disabled:opacity-50"
                    >Save</button
                >
            </div>
        </div>
    </div>
{/if}

<ConfirmDialog
    show={isConfirmOpen}
    title="Delete Holiday"
    message="Are you sure you want to delete this holiday? This action cannot be undone."
    confirmText="Confirm Delete"
    variant="danger"
    onConfirm={handleDelete}
    onCancel={() => {
        isConfirmOpen = false;
        holidayToDelete = null;
    }}
/>
