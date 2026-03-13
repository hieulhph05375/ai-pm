<script lang="ts">
    import {
        projectService,
        type PMIStats,
        type ProjectSnapshot,
        type MilestoneSnapshot,
    } from "$lib/services/projects";
    import { onMount } from "svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Badge from "$lib/components/ui/Badge.svelte";
    import TrendChart from "$lib/components/charts/TrendChart.svelte";
    import { toast } from "$lib/stores/toast";
    import { goto } from "$app/navigation";
    import { authStore } from "$lib/services/auth";
    import { hasProjectPermission } from "$lib/utils/permission";

    // To handle permissions for the banner button
    import {
        projectMembersService,
        type ProjectMember,
    } from "$lib/services/projectMembers";
    import {
        projectRolesService,
        type ProjectRole,
    } from "$lib/services/projectRoles";

    const formatDate = (d: string | Date) =>
        new Intl.DateTimeFormat("en-US", {
            month: "short",
            day: "2-digit",
            year: "numeric",
        }).format(new Date(d));

    let { projectId } = $props();
    let stats: PMIStats | null = $state(null);
    let trends: ProjectSnapshot[] = $state([]);
    let milestoneTrends: MilestoneSnapshot[] = $state([]);
    let projectMembers: ProjectMember[] = $state([]);
    let projectRoles: ProjectRole[] = $state([]);
    let loading = $state(true);

    onMount(async () => {
        try {
            const [s, t, mt, m, r] = await Promise.all([
                projectService.getPMIStats(projectId),
                projectService.getTrends(projectId),
                projectService.getMilestoneTrends(projectId),
                projectMembersService.getMembers(projectId, 1, 1000),
                projectRolesService.getRoles(projectId),
            ]);
            stats = s;
            trends = t || [];
            milestoneTrends = mt || [];
            projectMembers = (m as any).data || [];
            projectRoles = r || [];
        } catch (e) {
            toast.error("Failed to load PMI stats or trends");
        } finally {
            loading = false;
        }
    });

    async function handleExportExcel() {
        try {
            await projectService.exportWBS(projectId);
            toast.success("Exported WBS file (Excel)");
        } catch (e) {
            toast.error("Error exporting Excel file");
        }
    }

    async function handleExportPDF() {
        try {
            await projectService.exportSummary(projectId);
            toast.success("Exported Summary Report (PDF)");
        } catch (e) {
            toast.error("Error exporting PDF file");
        }
    }

    function getMetricColor(value: number) {
        if (value >= 1)
            return "text-emerald-600 bg-emerald-50 border-emerald-100";
        if (value >= 0.85) return "text-amber-600 bg-amber-50 border-amber-100";
        return "text-rose-600 bg-rose-50 border-rose-100";
    }
</script>

<div class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
    {#if stats?.needs_update}
        <!-- Reminder Banner -->
        <div
            class="bg-amber-50 border border-amber-200 rounded-3xl p-6 flex flex-col md:flex-row items-center justify-between gap-6 animate-in slide-in-from-top-4 duration-500 shadow-xl shadow-amber-500/10"
        >
            <div class="flex items-center gap-6">
                <div
                    class="size-14 rounded-2xl bg-gradient-to-br from-amber-400 to-amber-500 flex items-center justify-center shadow-lg shadow-amber-500/30"
                >
                    <span class="material-symbols-outlined text-white text-3xl"
                        >notifications_active</span
                    >
                </div>
                <div class="space-y-1">
                    <h3 class="text-lg font-outfit font-bold text-amber-900">
                        Progress Update Required
                    </h3>
                    <p class="text-sm text-amber-800/80 font-medium">
                        This project hasn't been updated recently. Please review
                        and update the WBS progress and actual costs.
                    </p>
                </div>
            </div>
            {#if hasProjectPermission($authStore.user, $authStore.token, projectMembers, projectRoles, "project:update")}
                <Button
                    onclick={() => goto(`/projects/${projectId}/wbs`)}
                    class="bg-amber-900 hover:bg-black text-amber-50 shadow-lg shadow-amber-900/20 whitespace-nowrap"
                >
                    Update Now
                </Button>
            {/if}
        </div>
    {/if}
    {#if loading}
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 animate-pulse">
            {#each Array(3) as _}
                <div
                    class="h-40 bg-slate-100 rounded-3xl border border-slate-200"
                ></div>
            {/each}
        </div>
    {:else if stats}
        <!-- Export Actions -->
        <div
            class="flex flex-wrap items-center justify-between gap-4 bg-slate-50 p-6 rounded-3xl border border-slate-200/60 shadow-inner"
        >
            <div class="space-y-1">
                <h3
                    class="text-sm font-bold text-slate-900 uppercase tracking-widest"
                >
                    Export Reports
                </h3>
                <p class="text-xs text-slate-500 font-medium">
                    Download reports in Excel or PDF format
                </p>
            </div>
            <div class="flex gap-3">
                <Button
                    variant="outline"
                    onclick={handleExportExcel}
                    class="bg-white hover:bg-emerald-50 hover:text-emerald-600 hover:border-emerald-200 border-slate-200 transition-all active:scale-95 gap-2"
                >
                    <span class="material-symbols-outlined text-[18px]"
                        >table_view</span
                    >
                    Export WBS (Excel)
                </Button>
                <Button
                    variant="outline"
                    onclick={handleExportPDF}
                    class="bg-white hover:bg-blue-50 hover:text-blue-600 hover:border-blue-200 border-slate-200 transition-all active:scale-95 gap-2"
                >
                    <span class="material-symbols-outlined text-[18px]"
                        >picture_as_pdf</span
                    >
                    PMI Report (PDF)
                </Button>
            </div>
        </div>

        <!-- Main Stats -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <!-- SPI -->
            <div
                class="bg-white p-8 rounded-3xl border border-slate-200 shadow-sm relative overflow-hidden group"
            >
                <div
                    class="absolute right-0 top-0 w-24 h-24 bg-primary/5 rounded-full blur-2xl -translate-y-1/2 translate-x-1/2"
                ></div>
                <div class="flex justify-between items-start mb-6">
                    <div class="space-y-1">
                        <span
                            class="text-[10px] font-bold text-slate-400 uppercase tracking-widest"
                            >Schedule Performance Index</span
                        >
                        <h4
                            class="text-lg font-outfit font-bold text-slate-900"
                        >
                            SPI Index
                        </h4>
                    </div>
                    <div
                        class="size-10 rounded-xl bg-slate-50 border border-slate-100 flex items-center justify-center text-slate-400"
                    >
                        <span class="material-symbols-outlined text-[20px]"
                            >schedule</span
                        >
                    </div>
                </div>
                <div class="flex items-baseline gap-3">
                    <span
                        class="text-5xl font-outfit font-bold text-slate-900 tracking-tight"
                        >{stats.spi}</span
                    >
                    <Badge class={getMetricColor(stats.spi)}>
                        {stats.spi >= 1
                            ? "On schedule"
                            : stats.spi >= 0.85
                              ? "Warning"
                              : "Behind schedule"}
                    </Badge>
                </div>
                <p
                    class="text-xs text-slate-500 mt-4 font-medium leading-relaxed"
                >
                    SPI measures time efficiency. A value >= 1.0 means the
                    project is on or ahead of schedule.
                </p>
            </div>

            <!-- CPI -->
            <div
                class="bg-white p-8 rounded-3xl border border-slate-200 shadow-sm relative overflow-hidden group"
            >
                <div
                    class="absolute right-0 top-0 w-24 h-24 bg-primary/5 rounded-full blur-2xl -translate-y-1/2 translate-x-1/2"
                ></div>
                <div class="flex justify-between items-start mb-6">
                    <div class="space-y-1">
                        <span
                            class="text-[10px] font-bold text-slate-400 uppercase tracking-widest"
                            >Cost Performance Index</span
                        >
                        <h4
                            class="text-lg font-outfit font-bold text-slate-900"
                        >
                            CPI Index
                        </h4>
                    </div>
                    <div
                        class="size-10 rounded-xl bg-slate-50 border border-slate-100 flex items-center justify-center text-slate-400"
                    >
                        <span class="material-symbols-outlined text-[20px]"
                            >payments</span
                        >
                    </div>
                </div>
                <div class="flex items-baseline gap-3">
                    <span
                        class="text-5xl font-outfit font-bold text-slate-900 tracking-tight"
                        >{stats.cpi}</span
                    >
                    <Badge class={getMetricColor(stats.cpi)}>
                        {stats.cpi >= 1
                            ? "Under budget"
                            : stats.cpi >= 0.85
                              ? "Warning"
                              : "Over budget"}
                    </Badge>
                </div>
                <p
                    class="text-xs text-slate-500 mt-4 font-medium leading-relaxed"
                >
                    CPI measures cost efficiency. A value >= 1.0 means the
                    project is under or on budget.
                </p>
            </div>

            <!-- EAC -->
            <div
                class="bg-white p-8 rounded-3xl border border-slate-200 shadow-sm relative overflow-hidden group"
            >
                <div
                    class="absolute right-0 top-0 w-24 h-24 bg-primary/5 rounded-full blur-2xl -translate-y-1/2 translate-x-1/2"
                ></div>
                <div class="flex justify-between items-start mb-6">
                    <div class="space-y-1">
                        <span
                            class="text-[10px] font-bold text-slate-400 uppercase tracking-widest"
                            >Estimate At Completion</span
                        >
                        <h4
                            class="text-lg font-outfit font-bold text-slate-900"
                        >
                            Estimate At Completion
                        </h4>
                    </div>
                    <div
                        class="size-10 rounded-xl bg-slate-50 border border-slate-100 flex items-center justify-center text-slate-400"
                    >
                        <span class="material-symbols-outlined text-[20px]"
                            >functions</span
                        >
                    </div>
                </div>
                <span
                    class="text-4xl font-outfit font-bold text-slate-900 tracking-tight"
                >
                    ${new Intl.NumberFormat("en-US").format(stats.eac)}
                </span>
                <p
                    class="text-xs text-slate-500 mt-4 font-medium leading-relaxed"
                >
                    Total expected cost at completion based on current cost
                    performance (CPI).
                </p>
            </div>
        </div>

        <!-- Earned Value Components -->
        <div class="bg-white rounded-3xl border border-slate-200 p-8 shadow-sm">
            <h3
                class="text-lg font-outfit font-bold text-slate-900 mb-8 flex items-center gap-3"
            >
                <div
                    class="size-8 rounded-lg bg-primary/10 flex items-center justify-center text-primary"
                >
                    <span class="material-symbols-outlined text-[18px]"
                        >finance</span
                    >
                </div>
                Earned Value Analysis
            </h3>

            <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
                <div class="space-y-4">
                    <div class="flex justify-between items-end">
                        <span
                            class="text-sm font-bold text-slate-500 uppercase tracking-widest"
                            >Planned Value (PV)</span
                        >
                        <span class="text-xl font-bold text-slate-900"
                            >${new Intl.NumberFormat("en-US").format(
                                stats.pv,
                            )}</span
                        >
                    </div>
                    <div
                        class="h-3 w-full bg-slate-100 rounded-full overflow-hidden"
                    >
                        <div
                            class="h-full bg-slate-300 rounded-full"
                            style="width: 100%"
                        ></div>
                    </div>
                    <p class="text-[11px] text-slate-400 font-medium">
                        The authorized budget assigned to scheduled work.
                    </p>
                </div>

                <div class="space-y-4">
                    <div class="flex justify-between items-end">
                        <span
                            class="text-sm font-bold text-primary uppercase tracking-widest"
                            >Earned Value (EV)</span
                        >
                        <span class="text-xl font-bold text-primary"
                            >${new Intl.NumberFormat("en-US").format(
                                stats.ev,
                            )}</span
                        >
                    </div>
                    <div
                        class="h-3 w-full bg-slate-100 rounded-full overflow-hidden"
                    >
                        <div
                            class="h-full bg-primary rounded-full transition-all duration-1000"
                            style="width: {(stats.ev / stats.pv) * 100}%"
                        ></div>
                    </div>
                    <p class="text-[11px] text-slate-400 font-medium">
                        The measure of work performed expressed in terms of the
                        budget authorized for that work.
                    </p>
                </div>

                <div class="space-y-4">
                    <div class="flex justify-between items-end">
                        <span
                            class="text-sm font-bold text-slate-500 uppercase tracking-widest"
                            >Actual Cost (AC)</span
                        >
                        <span
                            class="text-xl font-bold {stats.ac > stats.ev
                                ? 'text-rose-500'
                                : 'text-emerald-500'}"
                            >${new Intl.NumberFormat("en-US").format(
                                stats.ac,
                            )}</span
                        >
                    </div>
                    <div
                        class="h-3 w-full bg-slate-100 rounded-full overflow-hidden"
                    >
                        <div
                            class="h-full {stats.ac > stats.ev
                                ? 'bg-rose-500'
                                : 'bg-emerald-500'} rounded-full transition-all duration-1000"
                            style="width: {(stats.ac / stats.pv) * 100}%"
                        ></div>
                    </div>
                    <p class="text-[11px] text-slate-400 font-medium">
                        The realized cost incurred for the work performed on an
                        activity during a specific time period.
                    </p>
                </div>
            </div>
        </div>

        <!-- Performance Trends Section -->
        {#if trends.length > 0}
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mt-8">
                <!-- SPI Trend -->
                <div
                    class="bg-white rounded-3xl border border-slate-200 p-8 shadow-sm"
                >
                    <div class="flex items-center justify-between mb-6">
                        <div class="space-y-1">
                            <h3
                                class="text-lg font-outfit font-bold text-slate-900 flex items-center gap-2"
                            >
                                SPI Trend
                            </h3>
                            <p class="text-xs text-slate-500 font-medium">
                                Historical Schedule Performance
                            </p>
                        </div>
                    </div>

                    <TrendChart
                        data={trends.map((t) => ({
                            date: t.captured_at,
                            value: t.spi,
                        }))}
                        color="secondary"
                        showBaseline={true}
                        baselineValue={1.0}
                    />
                </div>

                <!-- CPI Trend -->
                <div
                    class="bg-white rounded-3xl border border-slate-200 p-8 shadow-sm"
                >
                    <div class="flex items-center justify-between mb-6">
                        <div class="space-y-1">
                            <h3
                                class="text-lg font-outfit font-bold text-slate-900 flex items-center gap-2"
                            >
                                CPI Trend
                            </h3>
                            <p class="text-xs text-slate-500 font-medium">
                                Historical Cost Performance
                            </p>
                        </div>
                    </div>

                    <TrendChart
                        data={trends.map((t) => ({
                            date: t.captured_at,
                            value: t.cpi,
                        }))}
                        color="primary"
                        showBaseline={true}
                        baselineValue={1.0}
                    />
                </div>
            </div>
        {/if}

        <!-- Milestone Trend Analysis -->
        {#if milestoneTrends.length > 0}
            <div
                class="bg-white rounded-3xl border border-slate-200 p-8 shadow-sm mt-6"
            >
                <div class="flex items-center justify-between mb-6">
                    <div class="space-y-1">
                        <h3
                            class="text-lg font-outfit font-bold text-slate-900 flex items-center gap-2"
                        >
                            <span
                                class="material-symbols-outlined text-[20px] text-amber-500"
                                >flag</span
                            >
                            Milestone Trends
                        </h3>
                        <p class="text-xs text-slate-500 font-medium">
                            Tracking completion dates of major milestones
                        </p>
                    </div>
                </div>

                <div
                    class="overflow-hidden rounded-2xl border border-slate-100"
                >
                    <table class="w-full text-left border-collapse">
                        <thead>
                            <tr class="bg-slate-50/80">
                                <th
                                    class="py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-widest border-b border-slate-100"
                                    >Milestone</th
                                >
                                <th
                                    class="py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-widest border-b border-slate-100"
                                    >Planned Date</th
                                >
                                <th
                                    class="py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-widest border-b border-slate-100"
                                    >Most Recent Projection</th
                                >
                                <th
                                    class="py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-widest border-b border-slate-100"
                                    >Status</th
                                >
                            </tr>
                        </thead>
                        <tbody class="text-sm">
                            {#each Object.values(milestoneTrends.reduce((acc, curr) => {
                                        // Group by node_id, keep the latest
                                        if (!acc[curr.node_id] || new Date(curr.captured_at) > new Date(acc[curr.node_id].captured_at)) {
                                            acc[curr.node_id] = curr;
                                        }
                                        return acc;
                                    }, {} as Record<number, MilestoneSnapshot>)) as milestone}
                                <tr
                                    class="border-b border-slate-50 hover:bg-slate-50/50 transition-colors"
                                >
                                    <td
                                        class="py-3 px-4 font-medium text-slate-900"
                                        >{milestone.milestone_name}</td
                                    >
                                    <td
                                        class="py-3 px-4 text-slate-600 font-mono text-xs"
                                        >{formatDate(
                                            milestone.planned_date,
                                        )}</td
                                    >
                                    <td
                                        class="py-3 px-4 text-slate-600 font-mono text-xs"
                                    >
                                        {formatDate(
                                            milestone.actual_date ||
                                                milestone.planned_date,
                                        )}
                                    </td>
                                    <td class="py-3 px-4">
                                        {#if milestone.actual_date && new Date(milestone.actual_date) > new Date(milestone.planned_date)}
                                            <Badge
                                                class="bg-rose-50 text-rose-600 border-rose-100"
                                                >Delayed</Badge
                                            >
                                        {:else if milestone.actual_date}
                                            <Badge
                                                class="bg-emerald-50 text-emerald-600 border-emerald-100"
                                                >Completed On Time</Badge
                                            >
                                        {:else}
                                            <Badge
                                                class="bg-slate-50 text-slate-600 border-slate-200"
                                                >Pending</Badge
                                            >
                                        {/if}
                                    </td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                </div>
            </div>
        {/if}
    {/if}
</div>
