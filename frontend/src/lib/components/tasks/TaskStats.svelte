<script lang="ts">
    import type { Task } from "$lib/services/tasks";
    import Card from "$lib/components/ui/Card.svelte";

    interface Props {
        tasks: Task[];
    }

    let { tasks = [] }: Props = $props();

    // Utility dates
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    const startOfWeek = new Date(today);
    startOfWeek.setDate(
        today.getDate() - today.getDay() + (today.getDay() === 0 ? -6 : 1),
    ); // Monday
    const endOfWeek = new Date(startOfWeek);
    endOfWeek.setDate(startOfWeek.getDate() + 6); // Sunday
    endOfWeek.setHours(23, 59, 59, 999);

    const startOfLastWeek = new Date(startOfWeek);
    startOfLastWeek.setDate(startOfWeek.getDate() - 7);
    const endOfLastWeek = new Date(startOfWeek);
    endOfLastWeek.setDate(startOfWeek.getDate() - 1);
    endOfLastWeek.setHours(23, 59, 59, 999);

    let stats = $derived.by(() => {
        let dueToday = 0;
        let overdue = 0;
        let pending = 0;

        let completedThisWeek = 0;
        let completedLastWeek = 0;
        let totalThisWeek = 0;

        tasks.forEach((task) => {
            // Count pending approvals (we treat TODO as pending in this context)
            if (task.status === "TODO") pending++;

            if (task.due_date) {
                const dueDate = new Date(task.due_date);
                dueDate.setHours(0, 0, 0, 0);

                if (task.status !== "DONE") {
                    if (dueDate.getTime() === today.getTime()) dueToday++;
                    else if (dueDate.getTime() < today.getTime()) overdue++;
                }
            }

            // Completion stats (this week vs last week)
            // Using updated_at as proxy for completed date if status is DONE
            if (task.updated_at) {
                const updatedDate = new Date(task.updated_at);

                // Tasks updated/created this week
                if (updatedDate >= startOfWeek && updatedDate <= endOfWeek) {
                    totalThisWeek++;
                    if (task.status === "DONE") completedThisWeek++;
                }

                // Tasks updated/created last week
                if (
                    updatedDate >= startOfLastWeek &&
                    updatedDate <= endOfLastWeek
                ) {
                    if (task.status === "DONE") completedLastWeek++;
                }
            }
        });

        // Weekly completion percentage
        const completionRate =
            totalThisWeek > 0
                ? Math.round((completedThisWeek / totalThisWeek) * 100)
                : 0;

        // Mock weekly change for UI demo (since calculating true historical % change needs snapshot data)
        const weeklyChange =
            completedLastWeek === 0
                ? "+12"
                : completionRate > completedLastWeek
                  ? `+${completionRate - completedLastWeek}`
                  : `${completionRate - completedLastWeek}`;

        return {
            dueToday,
            overdue,
            pending,
            completionRate,
            weeklyChange,
        };
    });
</script>

<div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4 mb-6">
    <!-- Tasks Due Today -->
    <Card
        class="bg-white border-slate-100 shadow-sm hover:shadow-md transition-all p-5 rounded-[20px] relative overflow-hidden"
    >
        <div class="flex items-start justify-between">
            <h3 class="text-lg font-semibold text-slate-500">Tasks Today</h3>
            <div
                class="size-10 rounded-xl bg-blue-50 text-blue-500 flex items-center justify-center"
            >
                <span class="material-symbols-outlined rounded text-[18px]"
                    >event_available</span
                >
            </div>
        </div>
        <div class="mt-4 flex items-baseline gap-2">
            <span class="text-3xl font-black font-outfit text-slate-900"
                >{stats.dueToday}</span
            >
            <span
                class="text-xs font-semibold px-2 py-0.5 rounded-md bg-slate-100 text-slate-500"
                >tasks</span
            >
        </div>
    </Card>

    <!-- Overdue Tasks -->
    <Card
        class="bg-white border-l-4 border-l-rose-500 border-slate-100 shadow-sm hover:shadow-md transition-all p-5 rounded-[20px] relative overflow-hidden"
    >
        <div class="flex items-start justify-between">
            <h3 class="text-lg font-semibold text-slate-500">Overdue</h3>
            <div
                class="size-10 rounded-xl bg-rose-50 text-rose-500 flex items-center justify-center"
            >
                <span class="material-symbols-outlined text-[18px]"
                    >warning</span
                >
            </div>
        </div>
        <div class="mt-4 flex items-baseline gap-2">
            <span class="text-3xl font-black font-outfit text-slate-900"
                >{stats.overdue}</span
            >
            {#if stats.overdue > 0}
                <span class="text-xs font-semibold text-rose-500"
                    >Immediate attention</span
                >
            {:else}
                <span class="text-xs font-semibold text-emerald-500"
                    >Excellent</span
                >
            {/if}
        </div>
    </Card>

    <!-- Pending Tasks (TODO) -->
    <Card
        class="bg-white border-slate-100 shadow-sm hover:shadow-md transition-all p-5 rounded-[20px] relative overflow-hidden"
    >
        <div class="flex items-start justify-between">
            <h3 class="text-lg font-semibold text-slate-500">To Do</h3>
            <div
                class="size-10 rounded-xl bg-amber-50 text-amber-500 flex items-center justify-center"
            >
                <span class="material-symbols-outlined text-[18px]"
                    >playlist_add_check</span
                >
            </div>
        </div>
        <div class="mt-4 flex items-baseline gap-2">
            <span class="text-3xl font-black font-outfit text-slate-900"
                >{stats.pending}</span
            >
            <span class="text-xs font-semibold text-slate-400">tasks TODO</span>
        </div>
    </Card>

    <!-- Weekly Completion -->
    <Card
        class="bg-white border-slate-100 shadow-sm hover:shadow-md transition-all p-5 rounded-[20px] relative overflow-hidden"
    >
        <div class="flex items-start justify-between">
            <h3 class="text-lg font-semibold text-slate-500">
                Weekly Completion
            </h3>
            <div
                class="size-10 rounded-xl bg-emerald-50 text-emerald-500 flex items-center justify-center"
            >
                <span class="material-symbols-outlined text-[18px]"
                    >trending_up</span
                >
            </div>
        </div>
        <div class="mt-3 flex items-baseline gap-2">
            <span class="text-3xl font-black font-outfit text-slate-900"
                >{stats.completionRate}%</span
            >
            <span class="text-xs font-semibold text-emerald-500"
                >{stats.weeklyChange}% vs last week</span
            >
        </div>
        <div
            class="mt-3 w-full h-1.5 bg-slate-100 rounded-full overflow-hidden"
        >
            <div
                class="h-full bg-primary rounded-full"
                style="width: {stats.completionRate}%"
            ></div>
        </div>
    </Card>
</div>
