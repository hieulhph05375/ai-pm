<script lang="ts">
    import type {
        WorkloadOverview,
        ResourceWorkload,
    } from "$lib/services/resource";
    import Badge from "$lib/components/ui/Badge.svelte";

    export let data: WorkloadOverview | null = null;
    export let selectedDepartment: string = "All";

    // Format dates simply
    const formatDay = (d: string) => {
        const date = new Date(d);
        return date.toLocaleDateString("en-US", {
            weekday: "short",
            day: "numeric",
        });
    };

    // Derived states
    $: allDates = data
        ? Array.from(
              new Set(data.users.flatMap((u) => u.entries.map((e) => e.date))),
          ).sort()
        : [];

    // Departments (Role map placeholder for now)
    $: departments = [
        "All",
        ...Array.from(new Set(data?.users.map((u) => u.role).filter(Boolean))),
    ];

    // Filter users
    $: filteredUsers =
        data?.users.filter(
            (u) =>
                selectedDepartment === "All" || u.role === selectedDepartment,
        ) || [];

    // Identify bottlenecks
    $: bottlenecks = filteredUsers
        .map((u) => {
            const overloads = u.entries.filter((e) => e.load_percentage > 100);
            return {
                user: u,
                overloadDays: overloads.length,
                maxLoad: Math.max(
                    0,
                    ...u.entries.map((e) => e.load_percentage),
                ),
            };
        })
        .filter((b) => b.overloadDays > 0)
        .sort((a, b) => b.maxLoad - a.maxLoad)
        .slice(0, 5); // top 5

    function getLoadColor(percentage: number) {
        if (percentage === 0) return "bg-slate-50 border-slate-100";
        if (percentage <= 80)
            return "bg-emerald-100 border-emerald-200 text-emerald-700";
        if (percentage <= 100)
            return "bg-amber-100 border-amber-200 text-amber-700";
        return "bg-rose-500 border-rose-600 text-white font-bold";
    }
</script>

<div class="space-y-6">
    <!-- Header with Filters & Summary -->
    <div class="flex flex-col lg:flex-row gap-6 justify-between items-start">
        <div class="space-y-2">
            <h3
                class="text-lg font-outfit font-bold text-slate-900 flex items-center gap-2"
            >
                <span
                    class="material-symbols-outlined text-[20px] text-indigo-500"
                    >grid_on</span
                >
                Resource Workload Analysis
            </h3>
            <p class="text-xs text-slate-500 font-medium">
                Daily capacity matrix highlighting organizational bottlenecks.
            </p>

            <div class="flex gap-2 mt-4 pt-2">
                {#each departments as dept}
                    <button
                        class="px-3 py-1 rounded-full text-xs font-bold border transition-colors {selectedDepartment ===
                        dept
                            ? 'bg-indigo-50 text-indigo-700 border-indigo-200'
                            : 'bg-white text-slate-600 border-slate-200 hover:bg-slate-50'}"
                        on:click={() => (selectedDepartment = dept)}
                    >
                        {dept}
                    </button>
                {/each}
            </div>
        </div>

        {#if bottlenecks.length > 0}
            <div
                class="bg-rose-50 rounded-2xl border border-rose-100 p-4 lg:w-1/3"
            >
                <div class="flex items-center gap-2 mb-3">
                    <span
                        class="material-symbols-outlined text-rose-500 text-[18px]"
                        >warning</span
                    >
                    <h4 class="text-sm font-bold text-rose-900">
                        Top Bottlenecks
                    </h4>
                </div>
                <div class="space-y-2">
                    {#each bottlenecks as b}
                        <div
                            class="flex justify-between items-center bg-white p-2 rounded-xl border border-rose-100 shadow-sm"
                        >
                            <div
                                class="truncate pr-2 text-xs font-bold text-slate-800"
                            >
                                {b.user.full_name}
                                <span
                                    class="text-[10px] text-slate-400 font-normal block truncate"
                                    >{b.user.role}</span
                                >
                            </div>
                            <Badge
                                class="bg-rose-500 text-white border-none text-[10px] shrink-0"
                            >
                                {b.maxLoad.toFixed(0)}% Load
                            </Badge>
                        </div>
                    {/each}
                </div>
            </div>
        {/if}
    </div>

    <!-- Heatmap Grid -->
    {#if data && filteredUsers.length > 0}
        <div
            class="overflow-x-auto overflow-y-auto max-h-[500px] border border-slate-200 rounded-2xl shadow-inner bg-white custom-scrollbar relative"
        >
            <table class="w-full text-left border-collapse min-w-max">
                <thead class="sticky top-0 z-10 bg-slate-50 shadow-sm">
                    <tr>
                        <th
                            class="py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-widest border-b border-r border-slate-200 bg-slate-50 sticky left-0 z-20 w-48 shrink-0"
                        >
                            Resource Name
                        </th>
                        {#each allDates as date}
                            <th
                                class="py-2 px-2 text-[10px] font-bold text-slate-500 uppercase tracking-wider border-b border-slate-100 text-center min-w-[50px]"
                            >
                                {formatDay(date)}
                            </th>
                        {/each}
                    </tr>
                </thead>
                <tbody class="text-sm divide-y divide-slate-100">
                    {#each filteredUsers as user}
                        <tr class="hover:bg-slate-50 transition-colors group">
                            <td
                                class="py-2 px-4 border-r border-slate-100 bg-white group-hover:bg-slate-50 sticky left-0 z-10 shadow-[2px_0_5px_rgba(0,0,0,0.02)]"
                            >
                                <div
                                    class="font-semibold text-slate-800 text-xs truncate w-40"
                                    title={user.full_name}
                                >
                                    {user.full_name}
                                </div>
                                <div
                                    class="text-[10px] text-slate-400 truncate w-40"
                                >
                                    {user.role || "Unspecified"}
                                </div>
                            </td>

                            {#each allDates as date}
                                {@const entry = user.entries.find(
                                    (e) => e.date === date,
                                )}
                                <td
                                    class="p-1 border-r border-slate-50 last:border-r-0 relative group/cell"
                                >
                                    {#if entry}
                                        <div
                                            class="w-full h-8 rounded-md border flex items-center justify-center text-[10px] transition-all cursor-help {getLoadColor(
                                                entry.load_percentage,
                                            )}"
                                        >
                                            {entry.load_percentage > 0
                                                ? `${entry.load_percentage.toFixed(0)}%`
                                                : "-"}
                                        </div>
                                        <!-- Tooltip -->
                                        <div
                                            class="absolute bottom-full left-1/2 -translate-x-1/2 mb-1 hidden group-hover/cell:flex flex-col items-center z-50 pointer-events-none"
                                        >
                                            <div
                                                class="bg-slate-800 text-white text-[10px] rounded-lg py-1.5 px-3 shadow-xl whitespace-nowrap"
                                            >
                                                <div
                                                    class="font-bold text-white mb-0.5"
                                                >
                                                    {entry.task_count} Tasks
                                                </div>
                                                <div class="text-slate-300">
                                                    {entry.load_percentage.toFixed(
                                                        0,
                                                    )}% ({entry.total_hours}h)
                                                </div>
                                            </div>
                                            <div
                                                class="w-2 h-2 border-t-[8px] border-t-slate-800 border-x-transparent border-x-[8px] border-b-0"
                                            ></div>
                                        </div>
                                    {:else}
                                        <div
                                            class="w-full h-8 rounded-md border border-slate-50 bg-slate-50/50 flex items-center justify-center text-[10px] text-slate-300"
                                        >
                                            -
                                        </div>
                                    {/if}
                                </td>
                            {/each}
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>

        <!-- Legend -->
        <div
            class="flex items-center gap-6 text-[11px] font-bold text-slate-500 uppercase flex-wrap"
        >
            <div class="flex items-center gap-2">
                <div
                    class="w-4 h-4 rounded bg-slate-50 border border-slate-100"
                ></div>
                 No Load (0%)
            </div>
            <div class="flex items-center gap-2">
                <div
                    class="w-4 h-4 rounded bg-emerald-100 border border-emerald-200"
                ></div>
                 Optimal (1-80%)
            </div>
            <div class="flex items-center gap-2">
                <div
                    class="w-4 h-4 rounded bg-amber-100 border border-amber-200"
                ></div>
                 Heavy (81-100%)
            </div>
            <div class="flex items-center gap-2">
                <div
                    class="w-4 h-4 rounded bg-rose-500 border border-rose-600"
                ></div>
                 Overloaded (>100%)
            </div>
        </div>
    {:else}
        <div
            class="h-40 border border-dashed border-slate-200 rounded-2xl flex items-center justify-center text-slate-400 font-medium text-sm"
        >
            No resource workload data available for this timeframe.
        </div>
    {/if}
</div>

<style>
    .custom-scrollbar::-webkit-scrollbar {
        width: 6px;
        height: 6px;
    }
    .custom-scrollbar::-webkit-scrollbar-track {
        background: transparent;
    }
    .custom-scrollbar::-webkit-scrollbar-thumb {
        background-color: #cbd5e1;
        border-radius: 20px;
    }
</style>
