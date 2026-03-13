<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/stores";
    import { issueService, type Issue } from "$lib/services/issue";
    import { categoryService, type Category } from "$lib/services/categories";
    import { toast } from "$lib/stores/toast";
    import ContentHeader from "$lib/components/ui/ContentHeader.svelte";
    import DataTable from "$lib/components/ui/DataTable.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Badge from "$lib/components/ui/Badge.svelte";
    import Modal from "$lib/components/ui/Modal.svelte";
    import {
        projectMembersService,
        type ProjectMember,
    } from "$lib/services/projectMembers";
    import {
        projectRolesService,
        type ProjectRole,
    } from "$lib/services/projectRoles";
    import { hasProjectPermission } from "$lib/utils/permission";
    import { authStore } from "$lib/services/auth";

    let projectId = $derived(Number($page.params.id));
    let issues = $state<Issue[]>([]);
    let loading = $state(true);
    let totalIssues = $state(0);
    let currentPage = $state(1);
    let itemsPerPage = $state(10);
    let showModal = $state(false);
    let editingIssue = $state<Issue | null>(null);
    let isSubmitting = $state(false);
    let projectMembers = $state<ProjectMember[]>([]);
    let projectRoles = $state<ProjectRole[]>([]);

    let types = $state<Category[]>([]);
    let priorities = $state<Category[]>([]);
    let statuses = $state<Category[]>([]);

    let form = $state<{
        project_id: number;
        type_id: number;
        title: string;
        description: string;
        status_id: number;
        priority_id: number;
        assignee_id: number | null;
        reporter_id: number | null;
    }>({
        project_id: 0,
        type_id: 0,
        title: "",
        description: "",
        status_id: 0,
        priority_id: 0,
        assignee_id: null,
        reporter_id: null,
    });

    const columns = [
        { key: "title", label: "Issue / Bug" },
        {
            key: "type",
            label: "Type",
            align: "center" as const,
            class: "w-24",
        },
        {
            key: "priority",
            label: "Priority",
            align: "center" as const,
            class: "w-24",
        },
        {
            key: "status",
            label: "Status",
            align: "center" as const,
            class: "w-32",
        },
        { key: "actions", label: "", align: "right" as const, class: "w-24" },
    ];

    onMount(async () => {
        await Promise.all([
            loadIssues(),
            loadCategories(),
            projectMembersService
                .getMembers(projectId, 1, 1000)
                .then((res) => (projectMembers = res.data)),
            projectRolesService
                .getRoles(projectId)
                .then((res) => (projectRoles = res)),
        ]);
    });

    async function loadCategories() {
        try {
            // Find type IDs first
            const typeRes = await categoryService.listTypes(1, 100);
            const typeMap = typeRes.data.reduce(
                (acc, t) => {
                    acc[t.code] = t.id;
                    return acc;
                },
                {} as Record<string, number>,
            );

            const [tRes, pRes, sRes] = await Promise.all([
                categoryService.listCategories(
                    1,
                    50,
                    "",
                    typeMap["ISSUE_TYPE"],
                ),
                categoryService.listCategories(
                    1,
                    50,
                    "",
                    typeMap["ISSUE_PRIORITY"],
                ),
                categoryService.listCategories(
                    1,
                    50,
                    "",
                    typeMap["ISSUE_STATUS"],
                ),
            ]);

            types = tRes.data;
            priorities = pRes.data;
            statuses = sRes.data;

            // Set default form values if creating
            if (!editingIssue) {
                if (types.length)
                    form.type_id =
                        types.find((t) => t.name === "Issue")?.id ||
                        types[0].id;
                if (priorities.length)
                    form.priority_id =
                        priorities.find((p) => p.name === "Medium")?.id ||
                        priorities[0].id;
                if (statuses.length)
                    form.status_id =
                        statuses.find((s) => s.name === "Open")?.id ||
                        statuses[0].id;
            }
        } catch (e: any) {
            console.error("Failed to load categories", e);
        }
    }

    async function loadIssues() {
        loading = true;
        try {
            const res = await issueService.list(
                projectId,
                currentPage,
                itemsPerPage,
            );
            issues = res.items;
            totalIssues = res.total;
        } catch (e: any) {
            if (!e.isAuthError) {
                toast.error(e.message || "Failed to load issues list");
            }
        } finally {
            loading = false;
        }
    }

    function handlePageChange(page: number) {
        currentPage = page;
        loadIssues();
    }

    function openCreate() {
        editingIssue = null;
        form = {
            project_id: projectId,
            type_id:
                types.find((t) => t.name === "Issue")?.id || types[0]?.id || 0,
            title: "",
            description: "",
            status_id:
                statuses.find((s) => s.name === "Open")?.id ||
                statuses[0]?.id ||
                0,
            priority_id:
                priorities.find((p) => p.name === "Medium")?.id ||
                priorities[0]?.id ||
                0,
            assignee_id: null,
            reporter_id: null,
        };
        showModal = true;
    }

    function openEdit(issue: Issue) {
        editingIssue = issue;
        form = {
            project_id: issue.project_id,
            type_id: issue.type_id,
            title: issue.title,
            description: issue.description || "",
            status_id: issue.status_id,
            priority_id: issue.priority_id,
            assignee_id: issue.assignee_id || null,
            reporter_id: issue.reporter_id || null,
        };
        showModal = true;
    }

    async function handleModalSubmit(e: Event) {
        e.preventDefault();
        isSubmitting = true;
        try {
            if (editingIssue) {
                await issueService.update(
                    projectId,
                    editingIssue.id,
                    form as any,
                );
                toast.success("Issue updated successfully");
            } else {
                await issueService.create(projectId, form as any);
                toast.success("New issue created successfully");
            }
            showModal = false;
            await loadIssues();
        } catch (e: any) {
            if (!e.isAuthError) {
                toast.error(e.message || "Failed to save issue");
            }
        } finally {
            isSubmitting = false;
        }
    }

    async function deleteIssue(id: number) {
        if (!confirm("Are you sure you want to delete this tracker entry?"))
            return;
        try {
            await issueService.delete(projectId, id);
            toast.success("Issue deleted successfully");
            await loadIssues();
        } catch (e: any) {
            if (!e.isAuthError) {
                toast.error(e.message || "Failed to delete issue");
            }
        }
    }

    function getStatusVariant(name: string) {
        switch (name) {
            case "Open":
                return "rose";
            case "In Progress":
                return "amber";
            case "Resolved":
                return "emerald";
            case "Closed":
                return "slate";
            default:
                return "primary";
        }
    }

    function getPriorityVariant(name: string) {
        switch (name) {
            case "Critical":
                return "rose";
            case "High":
                return "rose";
            case "Medium":
                return "amber";
            case "Low":
                return "emerald";
            default:
                return "slate";
        }
    }

    function getTypeVariant(name: string) {
        return name === "Bug" ? "rose" : "primary";
    }

    const breadcrumbs = $derived([
        { label: "Projects", href: "/projects" },
        { label: `Project #${projectId}`, href: `/projects/${projectId}` },
        { label: "Issue & Bug Tracker" },
    ]);
</script>

<svelte:head>
    <title>Issue Tracker | Project Dashboard</title>
</svelte:head>

<ContentHeader
    title="Issue & Bug Tracker"
    subtitle="Track and resolve project issues, bugs, and blockers efficiently"
>
    {#snippet titleSnippet()}
        <div class="flex flex-col gap-1 mb-1">
            <div
                class="flex items-center gap-2 text-[11px] font-bold text-slate-400 uppercase tracking-widest"
            >
                {#each breadcrumbs as crumb, i}
                    {#if i > 0}
                        <span class="material-symbols-outlined text-[12px]"
                            >chevron_right</span
                        >
                    {/if}
                    {#if crumb.href}
                        <a
                            href={crumb.href}
                            class="hover:text-primary transition-colors"
                            >{crumb.label}</a
                        >
                    {:else}
                        <span class="text-slate-500">{crumb.label}</span>
                    {/if}
                {/each}
            </div>
            <h2
                class="font-display text-2xl font-bold text-slate-900 tracking-tight"
            >
                Issue & Bug Tracker
            </h2>
        </div>
    {/snippet}

    <div class="flex items-center gap-3">
        {#if hasProjectPermission($authStore.user, $authStore.token, projectMembers, projectRoles, "project:issue:create")}
            <Button icon="add" onclick={openCreate}>Add Issue</Button>
        {/if}
    </div>
</ContentHeader>

<div class="space-y-4">
    <DataTable
        items={issues}
        {columns}
        {loading}
        total={totalIssues}
        page={currentPage}
        limit={itemsPerPage}
        onPageChange={handlePageChange}
        onLimitChange={(l) => {
            itemsPerPage = l;
            currentPage = 1;
            loadIssues();
        }}
    >
        {#snippet filterBar()}
            <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                    <!-- Filter bar content here if needed -->
                </div>
            </div>
        {/snippet}
        {#snippet rowCell({ item, column })}
            {#if column.key === "title"}
                <div>
                    <p class="font-bold text-slate-900 text-sm leading-tight">
                        {item.title}
                    </p>
                    {#if item.description}
                        <p
                            class="text-[11px] text-slate-400 font-medium mt-0.5 line-clamp-1"
                        >
                            {item.description}
                        </p>
                    {/if}
                </div>
            {:else if column.key === "type"}
                <Badge variant={getTypeVariant(item.type?.name || "") as any}
                    >{item.type?.name || "Unknown"}</Badge
                >
            {:else if column.key === "priority"}
                <Badge
                    variant={getPriorityVariant(
                        item.priority?.name || "",
                    ) as any}>{item.priority?.name || "Unknown"}</Badge
                >
            {:else if column.key === "status"}
                <Badge
                    variant={getStatusVariant(item.status?.name || "") as any}
                    >{item.status?.name || "Unknown"}</Badge
                >
            {:else if column.key === "actions"}
                <div
                    class="flex items-center justify-end gap-1 opacity-0 group-hover:opacity-100 transition-opacity"
                >
                    {#if hasProjectPermission($authStore.user, $authStore.token, projectMembers, projectRoles, "project:issue:update")}
                        <button
                            class="p-2 text-slate-400 hover:text-primary hover:bg-primary/5 rounded-lg transition-all"
                            title="Edit"
                            onclick={() => openEdit(item)}
                        >
                            <span class="material-symbols-outlined text-[18px]"
                                >edit</span
                            >
                        </button>
                    {/if}
                    {#if hasProjectPermission($authStore.user, $authStore.token, projectMembers, projectRoles, "project:issue:delete")}
                        <button
                            class="p-2 text-slate-400 hover:text-rose-500 hover:bg-rose-50 rounded-lg transition-all"
                            title="Delete"
                            onclick={() => deleteIssue(item.id)}
                        >
                            <span class="material-symbols-outlined text-[18px]"
                                >delete</span
                            >
                        </button>
                    {/if}
                </div>
            {/if}
        {/snippet}

        {#snippet emptyState()}
            <div
                class="flex flex-col items-center justify-center py-20 text-center"
            >
                <div
                    class="size-16 bg-slate-50 rounded-full flex items-center justify-center mb-4"
                >
                    <span
                        class="material-symbols-outlined text-3xl text-slate-200"
                        >bug_report</span
                    >
                </div>
                <p class="text-slate-500 font-bold">No issues found</p>
                <p class="text-slate-400 text-sm mt-1 max-w-xs">
                    The tracker is clean! Use the button above to log a new bug
                    or project issue.
                </p>
                <Button
                    variant="outline"
                    class="mt-6"
                    icon="add"
                    onclick={openCreate}>Create First Issue</Button
                >
            </div>
        {/snippet}
    </DataTable>
</div>

<Modal
    show={showModal}
    onClose={() => (showModal = false)}
    title={editingIssue ? "Edit Issue Entry" : "Create New Issue"}
>
    <form id="issue-form" onsubmit={handleModalSubmit} class="space-y-5">
        <div class="space-y-1.5 flex flex-col">
            <label
                for="title"
                class="text-[11px] font-bold uppercase tracking-wider text-slate-500 ml-1"
            >
                Title *
            </label>
            <input
                id="title"
                type="text"
                bind:value={form.title}
                required
                placeholder="Short summary of the issue"
                class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all text-sm font-medium"
            />
        </div>

        <div class="space-y-1.5 flex flex-col">
            <label
                for="description"
                class="text-[11px] font-bold uppercase tracking-wider text-slate-500 ml-1"
            >
                Description
            </label>
            <textarea
                id="description"
                bind:value={form.description}
                rows="3"
                placeholder="Detailed explanation, steps to reproduce, or notes..."
                class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all text-sm font-medium resize-none"
            ></textarea>
        </div>

        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-1.5 flex flex-col">
                <label
                    for="type_id"
                    class="text-[11px] font-bold uppercase tracking-wider text-slate-500 ml-1"
                >
                    Type
                </label>
                <select
                    id="type_id"
                    bind:value={form.type_id}
                    class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all text-sm font-medium"
                >
                    {#each types as t}
                        <option value={t.id}>{t.name}</option>
                    {/each}
                </select>
            </div>
            <div class="space-y-1.5 flex flex-col">
                <label
                    for="priority_id"
                    class="text-[11px] font-bold uppercase tracking-wider text-slate-500 ml-1"
                >
                    Priority
                </label>
                <select
                    id="priority_id"
                    bind:value={form.priority_id}
                    class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all text-sm font-medium"
                >
                    {#each priorities as p}
                        <option value={p.id}>{p.name}</option>
                    {/each}
                </select>
            </div>
        </div>

        <div class="space-y-1.5 flex flex-col">
            <label
                for="status_id"
                class="text-[11px] font-bold uppercase tracking-wider text-slate-500 ml-1"
            >
                Status
            </label>
            <select
                id="status_id"
                bind:value={form.status_id}
                class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all text-sm font-medium"
            >
                {#each statuses as s}
                    <option value={s.id}>{s.name}</option>
                {/each}
            </select>
        </div>
    </form>

    {#snippet footer()}
        <Button
            variant="outline"
            class="flex-1"
            onclick={() => (showModal = false)}>Cancel</Button
        >
        <Button
            type="submit"
            form="issue-form"
            class="flex-[2] justify-center"
            disabled={isSubmitting}
        >
            {isSubmitting
                ? "Saving..."
                : editingIssue
                  ? "Update Issue"
                  : "Create Issue"}
        </Button>
    {/snippet}
</Modal>
