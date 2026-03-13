<script lang="ts">
    import { onMount } from "svelte";
    import {
        portfolioService,
        type PortfolioOverview,
    } from "$lib/services/portfolio";
    import {
        resourceService,
        type WorkloadOverview,
    } from "$lib/services/resource";
    import ResourceHeatmap from "$lib/components/charts/ResourceHeatmap.svelte";

    let overview: PortfolioOverview | null = null;
    let workloadData: WorkloadOverview | null = null;
    let loading = true;
    let error = "";

    onMount(async () => {
        try {
            const end = new Date();
            end.setDate(end.getDate() + 28);
            const startDate = new Date().toISOString().split("T")[0];
            const endDate = end.toISOString().split("T")[0];

            const [o, w] = await Promise.all([
                portfolioService.getOverview(),
                resourceService.getWorkload(startDate, endDate),
            ]);
            overview = o;
            workloadData = w;
        } catch (e) {
            error = "Failed to load portfolio data. Please try again.";
        } finally {
            loading = false;
        }
    });

    function healthColor(health: string) {
        if (health === "Green")
            return "text-emerald-600 bg-emerald-50 border-emerald-200";
        if (health === "Yellow")
            return "text-amber-600 bg-amber-50 border-amber-200";
        return "text-rose-600 bg-rose-50 border-rose-200";
    }

    function healthLabel(health: string) {
        if (health === "Green") return "Good";
        if (health === "Yellow") return "Needs Attention";
        return "At Risk";
    }

    function formatCurrency(val: number) {
        return new Intl.NumberFormat("vi-VN", {
            style: "currency",
            currency: "VND",
            maximumFractionDigits: 0,
        }).format(val);
    }

    // SVG Donut chart
    function donutSegments(g: number, y: number, r: number, total: number) {
        if (total === 0) return [];
        const segments = [
            { value: g, color: "#10b981", label: "Good" },
            { value: y, color: "#f59e0b", label: "Needs Attention" },
            { value: r, color: "#ef4444", label: "Risk" },
        ];
        const cx = 60,
            cy = 60,
            radius = 50,
            inner = 30;
        let startAngle = -Math.PI / 2;
        return segments
            .filter((s) => s.value > 0)
            .map((s) => {
                const angle = (s.value / total) * 2 * Math.PI;
                const endAngle = startAngle + angle;
                const x1 = cx + radius * Math.cos(startAngle);
                const y1 = cy + radius * Math.sin(startAngle);
                const x2 = cx + radius * Math.cos(endAngle);
                const y2 = cy + radius * Math.sin(endAngle);
                const ix1 = cx + inner * Math.cos(startAngle);
                const iy1 = cy + inner * Math.sin(startAngle);
                const ix2 = cx + inner * Math.cos(endAngle);
                const iy2 = cy + inner * Math.sin(endAngle);
                const large = angle > Math.PI ? 1 : 0;
                const path = `M ${ix1} ${iy1} L ${x1} ${y1} A ${radius} ${radius} 0 ${large} 1 ${x2} ${y2} L ${ix2} ${iy2} A ${inner} ${inner} 0 ${large} 0 ${ix1} ${iy1} Z`;
                const result = {
                    path,
                    color: s.color,
                    label: s.label,
                    value: s.value,
                };
                startAngle = endAngle;
                return result;
            });
    }
</script>

<svelte:head>
    <title>Portfolio Dashboard | Enterprise PM</title>
    <meta
        name="description"
        content="Portfolio overview dashboard for PMO leaders"
    />
</svelte:head>

<div
    class="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50/30 to-indigo-50/20 p-6"
>
    <!-- Header -->
    <div class="mb-8">
        <h1 class="text-3xl font-black text-slate-900 tracking-tight">
            Portfolio Dashboard
        </h1>
        <p class="text-slate-500 mt-1 text-sm font-medium">
            Project Portfolio Overview — PMO View
        </p>
    </div>

    {#if loading}
        <!-- Loading skeleton -->
        <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4 mb-8">
            {#each Array(4) as _}
                <div
                    class="h-32 rounded-2xl bg-white/60 animate-pulse border border-slate-100"
                ></div>
            {/each}
        </div>
    {:else if error}
        <!-- Error state -->
        <div
            class="flex flex-col items-center justify-center py-24 text-center"
        >
            <span class="material-symbols-outlined text-6xl text-slate-300 mb-4"
                >error_outline</span
            >
            <h2 class="text-xl font-bold text-slate-700 mb-2">
                Failed to load data
            </h2>
            <p class="text-slate-400 text-sm">{error}</p>
        </div>
    {:else if overview}
        <!-- Stats Cards Row -->
        <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4 mb-8">
            <div
                class="bg-white/80 backdrop-blur-sm rounded-2xl border border-slate-200/60 shadow-sm p-5 hover:shadow-md transition-shadow"
            >
                <div class="flex items-center justify-between mb-3">
                    <span
                        class="text-xs font-bold text-slate-500 uppercase tracking-wider"
                        >Total Projects</span
                    >
                    <span
                        class="material-symbols-outlined text-indigo-500 text-xl"
                        >folder_open</span
                    >
                </div>
                <div class="text-4xl font-black text-slate-900">
                    {overview.total_projects}
                </div>
                <div class="flex gap-2 mt-3 flex-wrap">
                    <span
                        class="text-[11px] bg-sky-50 text-sky-700 border border-sky-200 px-2 py-0.5 rounded-full font-semibold"
                    >
                        {overview.active_projects} Active
                    </span>
                    <span
                        class="text-[11px] bg-emerald-50 text-emerald-700 border border-emerald-200 px-2 py-0.5 rounded-full font-semibold"
                    >
                        {overview.completed_projects} Completed
                    </span>
                    <span
                        class="text-[11px] bg-slate-50 text-slate-600 border border-slate-200 px-2 py-0.5 rounded-full font-semibold"
                    >
                        {overview.on_hold_projects} On Hold
                    </span>
                </div>
            </div>

            <div
                class="bg-white/80 backdrop-blur-sm rounded-2xl border border-slate-200/60 shadow-sm p-5 hover:shadow-md transition-shadow"
            >
                <div class="flex items-center justify-between mb-3">
                    <span
                        class="text-xs font-bold text-slate-500 uppercase tracking-wider"
                        >Total Budget</span
                    >
                    <span
                        class="material-symbols-outlined text-emerald-500 text-xl"
                        >payments</span
                    >
                </div>
                <div class="text-2xl font-black text-slate-900 leading-tight">
                    {formatCurrency(overview.total_budget)}
                </div>
                <p class="text-xs text-slate-400 mt-2">
                    Total approved budget across portfolio
                </p>
            </div>

            <div
                class="bg-white/80 backdrop-blur-sm rounded-2xl border border-slate-200/60 shadow-sm p-5 hover:shadow-md transition-shadow"
            >
                <div class="flex items-center justify-between mb-3">
                    <span
                        class="text-xs font-bold text-slate-500 uppercase tracking-wider"
                        >Portfolio Health</span
                    >
                    <span
                        class="material-symbols-outlined text-amber-500 text-xl"
                        >health_and_safety</span
                    >
                </div>
                <div class="flex items-center gap-3">
                    <!-- Donut SVG -->
                    <svg width="120" height="120" viewBox="0 0 120 120">
                        {#each donutSegments(overview.green_projects, overview.yellow_projects, overview.red_projects, overview.total_projects) as seg}
                            <path d={seg.path} fill={seg.color} opacity="0.9" />
                        {/each}
                        <text
                            x="60"
                            y="56"
                            text-anchor="middle"
                            font-size="16"
                            font-weight="800"
                            fill="#1e293b">{overview.total_projects}</text
                        >
                        <text
                            x="60"
                            y="70"
                            text-anchor="middle"
                            font-size="7"
                            fill="#94a3b8">projects</text
                        >
                    </svg>
                    <div class="flex flex-col gap-1.5 text-xs">
                        <span class="flex items-center gap-1.5"
                            ><span
                                class="size-2 rounded-full bg-emerald-500 inline-block"
                            ></span><span class="text-slate-600"
                                >{overview.green_projects} Good</span
                            ></span
                        >
                        <span class="flex items-center gap-1.5"
                            ><span
                                class="size-2 rounded-full bg-amber-400 inline-block"
                            ></span><span class="text-slate-600"
                                >{overview.yellow_projects} Attention</span
                            ></span
                        >
                        <span class="flex items-center gap-1.5"
                            ><span
                                class="size-2 rounded-full bg-rose-500 inline-block"
                            ></span><span class="text-slate-600"
                                >{overview.red_projects} Risk</span
                            ></span
                        >
                    </div>
                </div>
            </div>

            <div
                class="bg-gradient-to-br from-rose-500 to-orange-500 rounded-2xl border border-rose-300/40 shadow-sm p-5 hover:shadow-md transition-shadow text-white"
            >
                <div class="flex items-center justify-between mb-3">
                    <span
                        class="text-xs font-bold text-rose-100 uppercase tracking-wider"
                        >High Risk</span
                    >
                    <span class="material-symbols-outlined text-white text-xl"
                        >warning</span
                    >
                </div>
                <div class="text-4xl font-black">{overview.red_projects}</div>
                <p class="text-xs text-rose-100 mt-2">
                    Projects requiring immediate action
                </p>
            </div>
        </div>

        <!-- High Risk Projects Table -->
        {#if overview.high_risk_projects.length > 0}
            <div
                class="bg-white/80 backdrop-blur-sm rounded-2xl border border-slate-200/60 shadow-sm overflow-hidden"
            >
                <div
                    class="px-6 py-4 border-b border-slate-100 flex items-center gap-3"
                >
                    <span
                        class="material-symbols-outlined text-rose-500 text-lg"
                        >priority_high</span
                    >
                    <h2 class="text-sm font-bold text-slate-800">
                        High Risk Projects — Action Required
                    </h2>
                    <span
                        class="ml-auto bg-rose-100 text-rose-700 text-xs font-bold px-2.5 py-1 rounded-full border border-rose-200"
                    >
                        {overview.high_risk_projects.length} projects
                    </span>
                </div>
                <div class="overflow-x-auto">
                    <table class="w-full text-sm">
                        <thead>
                            <tr
                                class="bg-slate-50/80 border-b border-slate-100"
                            >
                                <th
                                    class="text-left px-6 py-3 text-xs font-bold text-slate-500 uppercase tracking-wider"
                                    >Project</th
                                >
                                <th
                                    class="text-left px-6 py-3 text-xs font-bold text-slate-500 uppercase tracking-wider"
                                    >Status</th
                                >
                                <th
                                    class="text-left px-6 py-3 text-xs font-bold text-slate-500 uppercase tracking-wider"
                                    >Health</th
                                >
                                <th
                                    class="text-right px-6 py-3 text-xs font-bold text-slate-500 uppercase tracking-wider"
                                    >Progress</th
                                >
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-slate-100">
                            {#each overview.high_risk_projects as project (project.id)}
                                <tr
                                    class="hover:bg-rose-50/40 transition-colors"
                                >
                                    <td class="px-6 py-4">
                                        <a
                                            href="/projects/{project.id}"
                                            class="font-semibold text-slate-800 hover:text-indigo-600 transition-colors"
                                        >
                                            {project.project_name}
                                        </a>
                                        <div
                                            class="text-xs text-slate-400 font-mono mt-0.5"
                                        >
                                            {project.project_id}
                                        </div>
                                    </td>
                                    <td class="px-6 py-4">
                                        <span
                                            class="text-xs font-semibold text-slate-600 bg-slate-100 px-2 py-1 rounded-full"
                                        >
                                            {project.project_status}
                                        </span>
                                    </td>
                                    <td class="px-6 py-4">
                                        <span
                                            class="text-xs font-bold px-2.5 py-1 rounded-full border {healthColor(
                                                project.overall_health,
                                            )}"
                                        >
                                            {healthLabel(
                                                project.overall_health,
                                            )}
                                        </span>
                                    </td>
                                    <td class="px-6 py-4 text-right">
                                        <div
                                            class="flex items-center justify-end gap-2"
                                        >
                                            <div
                                                class="w-20 bg-slate-100 rounded-full h-1.5 overflow-hidden"
                                            >
                                                <div
                                                    class="h-1.5 rounded-full bg-gradient-to-r from-indigo-500 to-sky-400 transition-all"
                                                    style="width: {project.progress}%"
                                                ></div>
                                            </div>
                                            <span
                                                class="text-xs font-bold text-slate-700 w-9 text-right"
                                                >{project.progress}%</span
                                            >
                                        </div>
                                    </td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                </div>
            </div>
        {:else}
            <div
                class="bg-white/80 backdrop-blur-sm rounded-2xl border border-slate-200/60 shadow-sm p-12 flex flex-col items-center text-center"
            >
                <span
                    class="material-symbols-outlined text-6xl text-emerald-300 mb-4"
                    >verified</span
                >
                <h2 class="text-xl font-bold text-slate-700 mb-2">
                    No High Risk Projects
                </h2>
                <p class="text-slate-400 text-sm">
                    All projects are in good standing or require routine
                    monitoring.
                </p>
            </div>
        {/if}

        <!-- Resource Workload Analysis (Heatmap) -->
        {#if workloadData}
            <div
                class="mt-8 bg-white/80 backdrop-blur-sm rounded-3xl border border-slate-200/60 shadow-sm p-8"
            >
                <ResourceHeatmap data={workloadData} />
            </div>
        {/if}
    {/if}
</div>
