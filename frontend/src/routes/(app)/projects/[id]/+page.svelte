<script lang="ts">
    import { onMount } from "svelte";
    import { projectService, type Project } from "$lib/services/projects";
    import { toast } from "$lib/stores/toast";
    import Badge from "$lib/components/ui/Badge.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import ProjectPmi from "$lib/components/projects/ProjectPMI.svelte";
    import ProjectTeam from "$lib/components/projects/ProjectTeam.svelte";
    import ProjectRbac from "$lib/components/projects/ProjectRBAC.svelte";
    import { goto } from "$app/navigation";
    import {
        projectMembersService,
        type ProjectMember,
    } from "$lib/services/projectMembers";
    import {
        projectRolesService,
        type ProjectRole,
    } from "$lib/services/projectRoles";
    import { hasPermission, hasProjectPermission } from "$lib/utils/permission";
    import { authStore } from "$lib/services/auth";

    let { data } = $props();
    let project: Project | null = $state(null);
    let pmiLoading = $state(false);
    let loading = $state(true);
    let activeTab = $state("dashboard"); // dashboard, team, pmi, rbac
    let projectMembers = $state<ProjectMember[]>([]);
    let projectRoles = $state<ProjectRole[]>([]);

    onMount(async () => {
        try {
            const id = Number(data.id);
            const [p, m, r] = await Promise.all([
                projectService.get(id),
                projectMembersService.getMembers(id, 1, 1000),
                projectRolesService.getRoles(id),
            ]);
            project = p;
            projectMembers = (m as any).data || [];
            projectRoles = r || [];
        } catch (e: any) {
            if (!e.isAuthError) {
                toast.error("Could not load project information");
            }
        } finally {
            loading = false;
        }
    });

    function getHealthColor(health: string) {
        switch (health?.toLowerCase()) {
            case "green":
                return "bg-emerald-500";
            case "yellow":
                return "bg-amber-500";
            case "red":
                return "bg-rose-500";
            default:
                return "bg-slate-500";
        }
    }
</script>

<div class="max-w-[1440px] mx-auto space-y-6">
    <!-- Breadcrumb -->
    <a
        href="/projects"
        class="inline-flex items-center gap-2 text-sm font-bold text-slate-500 hover:text-primary transition-colors hover:translate-x-1 duration-300"
    >
        <span class="material-symbols-outlined text-[18px]">arrow_back</span>
        Back to list
    </a>

    {#if loading}
        <!-- Loading Skeleton -->
        <div class="animate-pulse space-y-6">
            <div class="h-40 bg-slate-200 rounded-3xl"></div>
            <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
                {#each Array(4) as _}
                    <div class="h-32 bg-slate-200 rounded-2xl"></div>
                {/each}
            </div>
        </div>
    {:else if project}
        <!-- Header -->
        <div
            class="bg-white rounded-3xl p-8 border border-slate-200 shadow-xl shadow-slate-200/50 flex flex-col lg:flex-row gap-6 justify-between items-start lg:items-center relative overflow-hidden group"
        >
            <div
                class="absolute right-0 top-0 w-64 h-64 bg-primary/5 rounded-full blur-3xl -translate-y-1/2 translate-x-1/2 group-hover:bg-primary/10 transition-colors duration-1000"
            ></div>

            <div
                class="relative z-10 flex flex-col md:flex-row gap-6 items-start md:items-center"
            >
                <div
                    class="size-20 rounded-2xl bg-gradient-to-br from-primary to-primary-light flex items-center justify-center shadow-lg shadow-primary/30"
                >
                    <span class="material-symbols-outlined text-white text-4xl"
                        >rocket_launch</span
                    >
                </div>
                <div>
                    <h1
                        class="text-3xl font-outfit font-bold text-slate-900 tracking-tight"
                    >
                        {project.project_name}
                    </h1>
                    <div class="flex flex-wrap items-center gap-3 mt-3">
                        <Badge
                            variant="primary"
                            class="!bg-blue-50 !text-blue-600 !border-blue-100 uppercase tracking-widest"
                            >{project.current_phase}</Badge
                        >
                        <span
                            class="text-sm font-bold text-slate-400 uppercase tracking-widest flex items-center gap-1.5"
                        >
                            <span class="size-1.5 rounded-full bg-slate-300"
                            ></span>
                            {project.project_id}
                        </span>
                    </div>
                </div>
            </div>

            <div class="relative z-10 flex flex-wrap gap-3 w-full lg:w-auto">
                {#if hasProjectPermission($authStore.user, $authStore.token, projectMembers, projectRoles, "project:read")}
                    <Button
                        onclick={() => goto(`/projects/${project?.id}/wbs`)}
                        class="bg-primary hover:bg-primary/90 flex-1 lg:flex-none justify-center gap-2 shadow-lg shadow-primary/20 transition-all active:scale-95"
                    >
                        <span class="material-symbols-outlined text-[18px]"
                            >account_tree</span
                        >
                        Open WBS
                    </Button>
                {/if}
                {#if hasProjectPermission($authStore.user, $authStore.token, projectMembers, projectRoles, "project:stakeholder:read")}
                    <Button
                        variant="outline"
                        onclick={() =>
                            goto(`/projects/${project?.id}/stakeholders`)}
                        class="hover:bg-slate-50 flex-1 lg:flex-none justify-center gap-2 transition-all active:scale-95 border-slate-200"
                    >
                        <span class="material-symbols-outlined text-[18px]"
                            >groups</span
                        >
                        Stakeholders
                    </Button>
                {/if}
                <Button
                    variant="outline"
                    onclick={() => goto(`/projects/${project?.id}/risks`)}
                    class="hover:bg-rose-50 hover:text-rose-600 hover:border-rose-200 flex-1 lg:flex-none justify-center gap-2 transition-all active:scale-95 border-slate-200"
                >
                    <span class="material-symbols-outlined text-[18px]"
                        >shield_alert</span
                    >
                    Risks
                </Button>
                <Button
                    variant="outline"
                    onclick={() => goto(`/projects/${project?.id}/issues`)}
                    class="hover:bg-amber-50 hover:text-amber-600 hover:border-amber-200 flex-1 lg:flex-none justify-center gap-2 transition-all active:scale-95 border-slate-200"
                >
                    <span class="material-symbols-outlined text-[18px]"
                        >bug_report</span
                    >
                    Issues
                </Button>
            </div>
        </div>

        <!-- Tabs -->
        <div
            class="flex items-center gap-2 p-1.5 bg-white border border-slate-200 rounded-2xl w-fit shadow-sm"
        >
            <button
                onclick={() => (activeTab = "dashboard")}
                class="px-6 py-2 rounded-xl text-sm font-bold transition-all {activeTab ===
                'dashboard'
                    ? 'bg-primary text-white shadow-lg shadow-primary/20'
                    : 'text-slate-500 hover:bg-slate-50'}"
            >
                Dashboard
            </button>
            <button
                onclick={() => (activeTab = "pmi")}
                class="px-6 py-2 rounded-xl text-sm font-bold transition-all {activeTab ===
                'pmi'
                    ? 'bg-primary text-white shadow-lg shadow-primary/20'
                    : 'text-slate-500 hover:bg-slate-50'}"
            >
                Reports & PMI
            </button>
            <button
                onclick={() => (activeTab = "team")}
                class="px-6 py-2 rounded-xl text-sm font-bold transition-all {activeTab ===
                'team'
                    ? 'bg-primary text-white shadow-lg shadow-primary/20'
                    : 'text-slate-500 hover:bg-slate-50'}"
            >
                Team Members
            </button>
            <button
                onclick={() => (activeTab = "rbac")}
                class="px-6 py-2 rounded-xl text-sm font-bold transition-all {activeTab ===
                'rbac'
                    ? 'bg-primary text-white shadow-lg shadow-primary/20'
                    : 'text-slate-500 hover:bg-slate-50'}"
            >
                Roles & Permissions
            </button>
        </div>

        {#if activeTab === "dashboard"}
            <!-- Status overview content ... -->

            <!-- Metric Cards -->
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                <!-- Progress -->
                <div
                    class="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm hover:shadow-xl hover:shadow-primary/5 hover:-translate-y-1 transition-all duration-300"
                >
                    <p
                        class="text-[11px] font-bold uppercase tracking-widest text-slate-400 mb-4"
                    >
                        Overall Progress
                    </p>
                    <div class="flex items-end gap-2 mb-3">
                        <span
                            class="text-4xl font-outfit font-bold text-slate-900"
                            >{project.progress}%</span
                        >
                    </div>
                    <div
                        class="h-2 w-full bg-slate-100 rounded-full overflow-hidden shadow-inner"
                    >
                        <div
                            class="h-full bg-gradient-to-r from-primary to-primary-light rounded-full transition-all duration-1000"
                            style="width: {project.progress}%"
                        ></div>
                    </div>
                </div>

                <!-- Health -->
                <div
                    class="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm hover:shadow-xl hover:shadow-primary/5 hover:-translate-y-1 transition-all duration-300"
                >
                    <p
                        class="text-[11px] font-bold uppercase tracking-widest text-slate-400 mb-4"
                    >
                        Project Health
                    </p>
                    <div class="flex items-center gap-4 mt-2">
                        <div
                            class="size-14 rounded-full {getHealthColor(
                                project.overall_health,
                            ).replace(
                                'bg-',
                                'bg-',
                            )}/10 flex items-center justify-center border border-{getHealthColor(
                                project.overall_health,
                            ).split('-')[1]}-200"
                        >
                            <div
                                class="size-6 rounded-full {getHealthColor(
                                    project.overall_health,
                                )} shadow-lg shadow-{getHealthColor(
                                    project.overall_health,
                                ).split('-')[1]}-500/50"
                            ></div>
                        </div>
                        <span
                            class="text-2xl font-outfit font-bold text-slate-900 uppercase tracking-tight"
                            >{project.overall_health}</span
                        >
                    </div>
                </div>

                <!-- Budget -->
                <div
                    class="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm hover:shadow-xl hover:shadow-primary/5 hover:-translate-y-1 transition-all duration-300"
                >
                    <p
                        class="text-[11px] font-bold uppercase tracking-widest text-slate-400 mb-4"
                    >
                        Approved Budget
                    </p>
                    <p
                        class="text-3xl font-outfit font-bold text-slate-900 tracking-tight"
                    >
                        ${new Intl.NumberFormat("en-US").format(
                            project.approved_budget,
                        )}
                    </p>
                    <div
                        class="flex items-center gap-2 mt-3 text-sm font-bold {project.actual_cost >
                        project.approved_budget
                            ? 'text-rose-500'
                            : 'text-emerald-500'}"
                    >
                        <span class="material-symbols-outlined text-[16px]"
                            >{project.actual_cost > project.approved_budget
                                ? "trending_up"
                                : "trending_down"}</span
                        >
                        <span
                            >Actual: ${new Intl.NumberFormat("en-US").format(
                                project.actual_cost,
                            )}</span
                        >
                    </div>
                </div>

                <!-- Dates -->
                <div
                    class="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm hover:shadow-xl hover:shadow-primary/5 hover:-translate-y-1 transition-all duration-300"
                >
                    <p
                        class="text-[11px] font-bold uppercase tracking-widest text-slate-400 mb-4"
                    >
                        Timeline
                    </p>
                    <div class="space-y-3">
                        <div class="flex justify-between items-center text-sm">
                            <span
                                class="text-slate-400 font-bold uppercase tracking-wider text-[11px]"
                                >Start:</span
                            >
                            <span
                                class="text-slate-900 font-bold bg-slate-50 px-3 py-1 rounded-lg border border-slate-100"
                                >{project.planned_start_date
                                    ? new Date(
                                          project.planned_start_date,
                                      ).toLocaleDateString("vi-VN")
                                    : "---"}</span
                            >
                        </div>
                        <div class="flex justify-between items-center text-sm">
                            <span
                                class="text-slate-400 font-bold uppercase tracking-wider text-[11px]"
                                >End:</span
                            >
                            <span
                                class="text-slate-900 font-bold bg-slate-50 px-3 py-1 rounded-lg border border-slate-100"
                                >{project.planned_end_date
                                    ? new Date(
                                          project.planned_end_date,
                                      ).toLocaleDateString("vi-VN")
                                    : "---"}</span
                            >
                        </div>
                    </div>
                </div>
            </div>

            <!-- Info Grid -->
            <div class="grid grid-cols-1 xl:grid-cols-3 gap-6">
                <!-- Left Column: Details -->
                <div
                    class="xl:col-span-2 bg-white rounded-3xl p-8 border border-slate-200 shadow-sm"
                >
                    <h3
                        class="text-lg font-outfit font-bold text-slate-900 mb-6 flex items-center gap-3"
                    >
                        <div
                            class="size-8 rounded-lg bg-primary/10 flex items-center justify-center text-primary"
                        >
                            <span class="material-symbols-outlined text-[18px]"
                                >description</span
                            >
                        </div>
                        Project Description
                    </h3>
                    <p
                        class="text-slate-600 font-medium leading-relaxed min-h-[80px]"
                    >
                        {project.description ||
                            "No detailed description for this project yet."}
                    </p>

                    <h3
                        class="text-lg font-outfit font-bold text-slate-900 mt-10 mb-6 flex items-center gap-3"
                    >
                        <div
                            class="size-8 rounded-lg bg-amber-500/10 flex items-center justify-center text-amber-500"
                        >
                            <span class="material-symbols-outlined text-[18px]"
                                >flag</span
                            >
                        </div>
                        Strategic Goals
                    </h3>
                    <p class="text-slate-600 font-medium leading-relaxed">
                        {project.strategic_goal ||
                            "Strategic goals not set yet."}
                    </p>
                </div>

                <!-- Right Column: Sidebar info -->
                <div class="space-y-6">
                    <!-- Team -->
                    <div
                        class="bg-white rounded-3xl p-8 border border-slate-200 shadow-sm"
                    >
                        <h3
                            class="text-xs font-bold text-slate-400 mb-6 uppercase tracking-widest flex items-center gap-2 border-b border-slate-100 pb-4"
                        >
                            <span class="material-symbols-outlined text-[16px]"
                                >groups</span
                            >
                            Key Personnel
                        </h3>
                        <div class="space-y-6">
                            <div class="flex items-center gap-4">
                                <div
                                    class="size-12 rounded-2xl bg-slate-50 border border-slate-100 flex items-center justify-center text-slate-400"
                                >
                                    <span
                                        class="material-symbols-outlined text-[18px]"
                                        >person</span
                                    >
                                </div>
                                <div>
                                    <p
                                        class="text-[10px] font-bold text-slate-400 uppercase tracking-widest"
                                    >
                                        Project Manager
                                    </p>
                                    <p
                                        class="text-base font-bold text-slate-900 mt-0.5"
                                    >
                                        {project.project_manager ||
                                            "Unassigned"}
                                    </p>
                                </div>
                            </div>
                            <div class="flex items-center gap-4">
                                <div
                                    class="size-12 rounded-2xl bg-amber-50 border border-amber-100 flex items-center justify-center text-amber-500/50"
                                >
                                    <span
                                        class="material-symbols-outlined text-[18px]"
                                        >workspace_premium</span
                                    >
                                </div>
                                <div>
                                    <p
                                        class="text-[10px] font-bold text-slate-400 uppercase tracking-widest"
                                    >
                                        Sponsor
                                    </p>
                                    <p
                                        class="text-base font-bold text-slate-900 mt-0.5"
                                    >
                                        {project.sponsor || "Unassigned"}
                                    </p>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- KPIs -->
                    <div
                        class="bg-white rounded-3xl p-8 border border-slate-200 shadow-sm"
                    >
                        <h3
                            class="text-xs font-bold text-slate-400 mb-6 uppercase tracking-widest flex items-center gap-2 border-b border-slate-100 pb-4"
                        >
                            <span class="material-symbols-outlined text-[16px]"
                                >monitoring</span
                            >
                            Performance Indicators (KPI)
                        </h3>
                        <div class="space-y-4">
                            <div
                                class="flex justify-between items-center p-4 rounded-2xl bg-slate-50 border border-slate-100"
                            >
                                <div>
                                    <span
                                        class="text-[11px] font-bold text-slate-500 uppercase tracking-widest block"
                                        >SPI (Schedule)</span
                                    >
                                    <span
                                        class="text-xs font-medium text-slate-400 mt-1 block"
                                        >Ratio of actual progress to plan</span
                                    >
                                </div>
                                <span
                                    class="text-2xl font-outfit font-bold {project.spi >=
                                    1
                                        ? 'text-emerald-500'
                                        : 'text-rose-500'}">{project.spi}</span
                                >
                            </div>
                            <div
                                class="flex justify-between items-center p-4 rounded-2xl bg-slate-50 border border-slate-100"
                            >
                                <div>
                                    <span
                                        class="text-[11px] font-bold text-slate-500 uppercase tracking-widest block"
                                        >CPI (Cost)</span
                                    >
                                    <span
                                        class="text-xs font-medium text-slate-400 mt-1 block"
                                        >Ratio of value to actual cost</span
                                    >
                                </div>
                                <span
                                    class="text-2xl font-outfit font-bold {project.cpi >=
                                    1
                                        ? 'text-emerald-500'
                                        : 'text-rose-500'}">{project.cpi}</span
                                >
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        {:else if activeTab === "team"}
            <ProjectTeam projectId={project.id!} />
        {:else if activeTab === "pmi"}
            <ProjectPmi projectId={project.id!} />
        {:else if activeTab === "rbac"}
            <ProjectRbac projectId={project.id!} {projectMembers} />
        {/if}
    {/if}
</div>
