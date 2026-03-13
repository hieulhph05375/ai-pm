<script lang="ts">
    import { onMount } from "svelte";
    import {
        resourceService,
        type WorkloadOverview,
    } from "$lib/services/resource";

    let overview: WorkloadOverview | null = null;
    let loading = true;
    let error = "";

    // Date range - default current month
    const now = new Date();
    let startDate = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, "0")}-01`;
    let endDate = new Date(now.getFullYear(), now.getMonth() + 1, 0)
        .toISOString()
        .split("T")[0];

    // Build date columns from start to end
    $: dateCols = buildDateRange(startDate, endDate);

    function buildDateRange(start: string, end: string): string[] {
        const dates: string[] = [];
        const cur = new Date(start);
        const last = new Date(end);
        while (cur <= last) {
            dates.push(cur.toISOString().split("T")[0]);
            cur.setDate(cur.getDate() + 1);
        }
        return dates;
    }

    async function loadWorkload() {
        loading = true;
        error = "";
        try {
            overview = await resourceService.getWorkload(startDate, endDate);
        } catch (e) {
            error = "Failed to load workload data. Please try again.";
        } finally {
            loading = false;
        }
    }

    onMount(loadWorkload);

    function getHours(
        user: NonNullable<WorkloadOverview["users"]>[number],
        date: string,
    ): number {
        return user.entries.find((e) => e.date === date)?.total_hours ?? 0;
    }

    function heatColor(hours: number): string {
        if (hours === 0) return "bg-slate-50 text-transparent";
        if (hours <= 4) return "bg-emerald-100 text-emerald-700";
        if (hours <= 8) return "bg-amber-200 text-amber-800";
        if (hours <= 12) return "bg-orange-300 text-orange-900";
        return "bg-rose-400 text-rose-900";
    }

    function isWeekend(date: string): boolean {
        const d = new Date(date);
        return d.getDay() === 0 || d.getDay() === 6;
    }

    function formatDate(date: string): string {
        const d = new Date(date);
        return `${d.getDate()}`;
    }

    function formatMonth(date: string): string {
        const d = new Date(date);
        return d.toLocaleDateString("vi-VN", {
            month: "short",
            day: "numeric",
        });
    }
</script>

<svelte:head>
    <title>Resource Allocation | Enterprise PM</title>
    <meta name="description" content="Resource allocation workload heatmap" />
</svelte:head>

<div
    class="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50/30 to-indigo-50/20 p-6"
>
    <!-- Header -->
    <div class="mb-6 flex flex-wrap items-center justify-between gap-4">
        <div>
            <h1 class="text-3xl font-black text-slate-900 tracking-tight">
                Resource Allocation
            </h1>
            <p class="text-slate-500 mt-1 text-sm font-medium">
                Workload Heatmap — Daily Workload Tracking
            </p>
        </div>
        <div
            class="flex items-center gap-3 bg-white/80 backdrop-blur-sm border border-slate-200/60 rounded-2xl px-4 py-3 shadow-sm"
        >
            <div class="flex flex-col gap-1">
                <label
                    for="start-date"
                    class="text-[10px] font-bold text-slate-400 uppercase tracking-wider"
                    >From Date</label
                >
                <input
                    id="start-date"
                    type="date"
                    bind:value={startDate}
                    class="text-sm font-semibold text-slate-700 border-none outline-none bg-transparent cursor-pointer"
                />
            </div>
            <span class="material-symbols-outlined text-slate-300 text-xl"
                >arrow_forward</span
            >
            <div class="flex flex-col gap-1">
                <label
                    for="end-date"
                    class="text-[10px] font-bold text-slate-400 uppercase tracking-wider"
                    >To Date</label
                >
                <input
                    id="end-date"
                    type="date"
                    bind:value={endDate}
                    class="text-sm font-semibold text-slate-700 border-none outline-none bg-transparent cursor-pointer"
                />
            </div>
            <button
                onclick={loadWorkload}
                class="ml-2 px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-bold rounded-xl transition-colors"
            >
                Update
            </button>
        </div>
    </div>

    <!-- Legend -->
    <div class="flex gap-3 mb-5 flex-wrap">
        <span class="text-xs font-semibold text-slate-500">Level:</span>
        {#each [{ label: "Empty", cls: "bg-slate-100" }, { label: "1-4h", cls: "bg-emerald-100" }, { label: "4-8h", cls: "bg-amber-200" }, { label: "8-12h", cls: "bg-orange-300" }, { label: ">12h", cls: "bg-rose-400" }] as item}
            <span
                class="flex items-center gap-1.5 text-xs text-slate-600 font-medium"
            >
                <span class="size-3 rounded {item.cls} inline-block"></span>
                {item.label}
            </span>
        {/each}
    </div>

    {#if loading}
        <div
            class="bg-white/80 backdrop-blur-sm rounded-2xl border border-slate-200/60 shadow-sm p-12 flex items-center justify-center"
        >
            <div
                class="animate-spin material-symbols-outlined text-4xl text-indigo-400"
            >
                refresh
            </div>
        </div>
    {:else if error}
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
    {:else if !overview || overview.users.length === 0}
        <div
            class="bg-white/80 backdrop-blur-sm rounded-2xl border border-slate-200/60 shadow-sm p-16 flex flex-col items-center text-center"
        >
            <span class="material-symbols-outlined text-6xl text-slate-200 mb-4"
                >grid_view</span
            >
            <h2 class="text-xl font-bold text-slate-700 mb-2">
                No allocation data
            </h2>
            <p class="text-slate-400 text-sm">
                No tasks assigned during this period.
            </p>
        </div>
    {:else}
        <div
            class="bg-white/80 backdrop-blur-sm rounded-2xl border border-slate-200/60 shadow-sm overflow-hidden"
        >
            <div class="overflow-x-auto">
                <table class="border-collapse" style="min-width: max-content;">
                    <thead>
                        <tr class="bg-slate-50/80 border-b border-slate-100">
                            <th
                                class="sticky left-0 z-10 bg-slate-50/95 backdrop-blur-sm px-5 py-3 text-left text-xs font-bold text-slate-500 uppercase tracking-wider min-w-[160px] border-r border-slate-100"
                            >
                                Member
                            </th>
                            {#each dateCols as date}
                                <th
                                    class="py-3 px-1 text-center text-[10px] font-bold min-w-[32px] {isWeekend(
                                        date,
                                    )
                                        ? 'bg-rose-50/60 text-rose-400'
                                        : 'text-slate-400'}"
                                    title={formatMonth(date)}
                                >
                                    {formatDate(date)}
                                </th>
                            {/each}
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-slate-50">
                        {#each overview.users as user (user.user_id)}
                            <tr
                                class="hover:bg-slate-50/50 transition-colors group"
                            >
                                <td
                                    class="sticky left-0 z-10 bg-white/95 group-hover:bg-slate-50/95 backdrop-blur-sm px-5 py-3 border-r border-slate-100"
                                >
                                    <div class="flex items-center gap-2">
                                        <div
                                            class="size-7 rounded-full bg-gradient-to-br from-indigo-500 to-sky-400 flex items-center justify-center text-white text-[10px] font-black flex-shrink-0"
                                        >
                                            {user.full_name
                                                .charAt(0)
                                                .toUpperCase()}
                                        </div>
                                        <div>
                                            <p
                                                class="text-xs font-bold text-slate-800 whitespace-nowrap"
                                            >
                                                {user.full_name}
                                            </p>
                                            <p
                                                class="text-[10px] text-slate-400"
                                            >
                                                {user.email}
                                            </p>
                                        </div>
                                    </div>
                                </td>
                                {#each dateCols as date}
                                    {@const hours = getHours(user, date)}
                                    <td
                                        class="px-1 py-3 text-center"
                                        title="{formatMonth(date)}: {hours}h"
                                    >
                                        {#if hours > 0}
                                            <span
                                                class="inline-flex items-center justify-center size-7 rounded-lg text-[9px] font-black {heatColor(
                                                    hours,
                                                )}"
                                            >
                                                {hours}
                                            </span>
                                        {:else}
                                            <span
                                                class="inline-flex items-center justify-center size-7 rounded-lg {isWeekend(
                                                    date,
                                                )
                                                    ? 'bg-rose-50'
                                                    : 'bg-slate-50'}"
                                            >
                                            </span>
                                        {/if}
                                    </td>
                                {/each}
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        </div>
    {/if}
</div>
