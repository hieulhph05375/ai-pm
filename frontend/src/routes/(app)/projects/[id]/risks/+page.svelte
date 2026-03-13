<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/stores";
    import {
        riskService,
        type Risk,
        type CreateRiskPayload,
    } from "$lib/services/risk";
    import { toast } from "$lib/stores/toast";
    import ContentHeader from "$lib/components/ui/ContentHeader.svelte";
    import DataTable from "$lib/components/ui/DataTable.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Badge from "$lib/components/ui/Badge.svelte";
    import { categoryService, type Category } from "$lib/services/categories";
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
    let risks = $state<Risk[]>([]);
    let loading = $state(true);
    let totalRisks = $state(0);
    let currentPage = $state(1);
    let itemsPerPage = $state(10);
    let showModal = $state(false);
    let editingRisk = $state<Risk | null>(null);
    let isSubmitting = $state(false);
    let projectMembers = $state<ProjectMember[]>([]);
    let projectRoles = $state<ProjectRole[]>([]);
    let riskStatuses = $state<Category[]>([]);

    // Form fields
    let form: CreateRiskPayload = $state({
        project_id: 0,
        title: "",
        description: "",
        probability: 3,
        impact: 3,
        status: "Open",
        owner_id: null,
    });

    const columns = [
        { key: "title", label: "Risk Description" },
        {
            key: "probability",
            label: "Prob.",
            align: "center" as const,
            class: "w-24",
        },
        {
            key: "impact",
            label: "Impact",
            align: "center" as const,
            class: "w-24",
        },
        {
            key: "risk_score",
            label: "Score",
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
            loadRisks(),
            projectMembersService
                .getMembers(projectId, 1, 1000)
                .then((res) => (projectMembers = res.data)),
            projectRolesService
                .getRoles(projectId)
                .then((res) => (projectRoles = res)),
            categoryService
                .listCategories(1, 100, "", 8)
                .then((res) => (riskStatuses = res.data)),
        ]);
    });

    async function loadRisks() {
        loading = true;
        try {
            const res = await riskService.list(
                projectId,
                currentPage,
                itemsPerPage,
            );
            risks = res.items;
            totalRisks = res.total;
        } catch (e: any) {
            if (!e.isAuthError) {
                toast.error(e.message || "Failed to load risk register");
            }
        } finally {
            loading = false;
        }
    }

    function handlePageChange(page: number) {
        currentPage = page;
        loadRisks();
    }

    function openCreate() {
        editingRisk = null;
        form = {
            project_id: projectId,
            title: "",
            description: "",
            probability: 3,
            impact: 3,
            status: "Open",
            status_id: riskStatuses.find((c) => c.is_active)?.id,
            owner_id: null,
        };
        showModal = true;
    }

    function openEdit(risk: Risk) {
        editingRisk = risk;
        form = { ...risk };
        showModal = true;
    }

    async function handleModalSubmit(e: Event) {
        e.preventDefault();
        isSubmitting = true;
        try {
            const selectedStatus = riskStatuses.find(
                (c) => c.id === form.status_id,
            );
            if (selectedStatus) {
                form.status = selectedStatus.name as any;
            }

            if (editingRisk) {
                await riskService.update(projectId, editingRisk.id, form);
                toast.success("Risk updated successfully");
            } else {
                await riskService.create(projectId, form);
                toast.success("New risk registered successfully");
            }
            showModal = false;
            await loadRisks();
        } catch (e: any) {
            if (!e.isAuthError) {
                toast.error(e.message || "Failed to save risk");
            }
        } finally {
            isSubmitting = false;
        }
    }

    async function deleteRisk(riskId: number) {
        if (!confirm("Are you sure you want to delete this risk entry?"))
            return;
        try {
            await riskService.delete(projectId, riskId);
            toast.success("Risk deleted successfully");
            await loadRisks();
        } catch (e: any) {
            if (!e.isAuthError) {
                toast.error(e.message || "Failed to delete risk");
            }
        }
    }

    function getScoreVariant(score: number) {
        if (score >= 15) return "rose";
        if (score >= 9) return "amber";
        return "emerald";
    }

    function getStatusInfo(item: Risk) {
        if (item.status_cat) {
            return {
                label: item.status_cat.name,
                color: item.status_cat.color || "slate",
            };
        }
        // Legacy fallback
        switch (item.status) {
            case "Open":
                return { label: "Open", color: "rose" };
            case "Mitigated":
                return { label: "Mitigated", color: "amber" };
            case "Closed":
                return { label: "Closed", color: "slate" };
            default:
                return { label: item.status, color: "primary" };
        }
    }

    const PROB_LABELS = [
        "",
        "Very Low (1)",
        "Low (2)",
        "Medium (3)",
        "High (4)",
        "Very High (5)",
    ];
    const IMPACT_LABELS = [
        "",
        "Insignificant (1)",
        "Minor (2)",
        "Moderate (3)",
        "Major (4)",
        "Catastrophic (5)",
    ];

    const breadcrumbs = $derived([
        { label: "Projects", href: "/projects" },
        { label: `Project #${projectId}`, href: `/projects/${projectId}` },
        { label: "Risk Register" },
    ]);
</script>

<svelte:head>
    <title>Risk Register | Project Dashboard</title>
</svelte:head>

<ContentHeader
    title="Risk Register"
    subtitle="Identify, analyze, and manage project risks and mitigation strategies"
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
                Risk Register
            </h2>
        </div>
    {/snippet}

    <div class="flex items-center gap-3">
        {#if hasProjectPermission($authStore.user, $authStore.token, projectMembers, projectRoles, "project:risk:create")}
            <Button icon="add" onclick={openCreate}>Add Risk</Button>
        {/if}
    </div>
</ContentHeader>

<div class="space-y-4">
    <DataTable
        items={risks}
        {columns}
        {loading}
        total={totalRisks}
        page={currentPage}
        limit={itemsPerPage}
        onPageChange={handlePageChange}
        onLimitChange={(l) => {
            itemsPerPage = l;
            currentPage = 1;
            loadRisks();
        }}
    >
        {#snippet filterBar()}
            <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                    <!-- Filter content here if needed -->
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
            {:else if column.key === "probability"}
                <span class="text-sm font-semibold text-slate-600"
                    >{item.probability}/5</span
                >
            {:else if column.key === "impact"}
                <span class="text-sm font-semibold text-slate-600"
                    >{item.impact}/5</span
                >
            {:else if column.key === "risk_score"}
                <Badge variant={getScoreVariant(item.risk_score) as any}>
                    {item.risk_score}
                </Badge>
            {:else if column.key === "status"}
                {@const st = getStatusInfo(item)}
                <Badge color={st.color as any}>{st.label}</Badge>
            {:else if column.key === "actions"}
                <div
                    class="flex items-center justify-end gap-1 opacity-0 group-hover:opacity-100 transition-opacity"
                >
                    {#if hasProjectPermission($authStore.user, $authStore.token, projectMembers, projectRoles, "project:risk:update")}
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
                    {#if hasProjectPermission($authStore.user, $authStore.token, projectMembers, projectRoles, "project:risk:delete")}
                        <button
                            class="p-2 text-slate-400 hover:text-rose-500 hover:bg-rose-50 rounded-lg transition-all"
                            title="Delete"
                            onclick={() => deleteRisk(item.id)}
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
                        >shield_with_heart</span
                    >
                </div>
                <p class="text-slate-500 font-bold">No risks recorded</p>
                <p class="text-slate-400 text-sm mt-1 max-w-xs">
                    Start by adding potential risks to track their impact and
                    probability.
                </p>
                <Button
                    variant="outline"
                    class="mt-6"
                    icon="add"
                    onclick={openCreate}>Register First Risk</Button
                >
            </div>
        {/snippet}
    </DataTable>
</div>

<Modal
    show={showModal}
    onClose={() => (showModal = false)}
    title={editingRisk ? "Edit Risk Entry" : "Register New Risk"}
>
    <form id="risk-form" onsubmit={handleModalSubmit} class="space-y-5">
        <div class="space-y-1.5 flex flex-col">
            <label
                for="title"
                class="text-[11px] font-bold uppercase tracking-wider text-slate-500 ml-1"
            >
                Risk Title *
            </label>
            <input
                id="title"
                type="text"
                bind:value={form.title}
                required
                placeholder="Brief description of the risk"
                class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all text-sm font-medium"
            />
        </div>

        <div class="space-y-1.5 flex flex-col">
            <label
                for="description"
                class="text-[11px] font-bold uppercase tracking-wider text-slate-500 ml-1"
            >
                Mitigation Strategy / Details
            </label>
            <textarea
                id="description"
                bind:value={form.description}
                rows="3"
                placeholder="How will this be managed or avoided?"
                class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all text-sm font-medium resize-none"
            ></textarea>
        </div>

        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-1.5 flex flex-col">
                <label
                    for="probability"
                    class="text-[11px] font-bold uppercase tracking-wider text-slate-500 ml-1"
                >
                    Probability
                </label>
                <select
                    id="probability"
                    bind:value={form.probability}
                    class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all text-sm font-medium"
                >
                    {#each [1, 2, 3, 4, 5] as v}
                        <option value={v}>{PROB_LABELS[v]}</option>
                    {/each}
                </select>
            </div>
            <div class="space-y-1.5 flex flex-col">
                <label
                    for="impact"
                    class="text-[11px] font-bold uppercase tracking-wider text-slate-500 ml-1"
                >
                    Impact
                </label>
                <select
                    id="impact"
                    bind:value={form.impact}
                    class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all text-sm font-medium"
                >
                    {#each [1, 2, 3, 4, 5] as v}
                        <option value={v}>{IMPACT_LABELS[v]}</option>
                    {/each}
                </select>
            </div>
        </div>

        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-1.5 flex flex-col">
                <label
                    for="status"
                    class="text-[11px] font-bold uppercase tracking-wider text-slate-500 ml-1"
                >
                    Status
                </label>
                <select
                    id="status"
                    bind:value={form.status_id}
                    class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all text-sm font-medium"
                >
                    {#each riskStatuses as s}
                        <option value={s.id}>{s.name}</option>
                    {/each}
                </select>
            </div>

            <div class="flex flex-col justify-end">
                <div
                    class="bg-slate-50 border border-slate-100 rounded-xl px-4 py-2.5 flex items-center justify-between"
                >
                    <span
                        class="text-[10px] font-bold text-slate-400 uppercase tracking-widest leading-none"
                        >Score</span
                    >
                    <div class="flex items-center gap-2">
                        <span
                            class="text-lg font-black text-slate-900 leading-none"
                            >{form.probability * form.impact}</span
                        >
                        <Badge
                            variant={getScoreVariant(
                                form.probability * form.impact,
                            ) as any}
                            class="!px-1.5 !py-0 !text-[10px]"
                        >
                            {form.probability * form.impact >= 15
                                ? "HIGH"
                                : form.probability * form.impact >= 9
                                  ? "MED"
                                  : "LOW"}
                        </Badge>
                    </div>
                </div>
            </div>
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
            form="risk-form"
            class="flex-[2] justify-center"
            disabled={isSubmitting}
        >
            {isSubmitting
                ? "Saving..."
                : editingRisk
                  ? "Update Risk"
                  : "Register Risk"}
        </Button>
    {/snippet}
</Modal>
