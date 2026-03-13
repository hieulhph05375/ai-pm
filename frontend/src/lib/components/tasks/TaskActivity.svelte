<script lang="ts">
    import { taskService, type TaskActivity } from "$lib/services/tasks";
    import Avatar from "$lib/components/ui/Avatar.svelte";
    import { onMount } from "svelte";

    interface Props {
        taskId: number;
        mode?: "all" | "activities" | "comments";
    }

    let { taskId, mode = "all" }: Props = $props();
    let activities = $state<TaskActivity[]>([]);
    let loading = $state(true);

    let filteredActivities = $derived.by(() => {
        if (mode === "activities") {
            return activities.filter((a) => a.action !== "COMMENT");
        }
        if (mode === "comments") {
            return activities.filter((a) => a.action === "COMMENT");
        }
        return activities;
    });

    let newComment = $state("");
    let submitting = $state(false);

    onMount(async () => {
        await loadActivities();
    });

    async function loadActivities() {
        loading = true;
        try {
            activities = await taskService.getActivities(taskId);
        } catch (e) {
            console.error(e);
        } finally {
            loading = false;
        }
    }

    function formatTime(dateStr: string) {
        const d = new Date(dateStr);
        return d.toLocaleDateString("en-US", {
            hour: "2-digit",
            minute: "2-digit",
            day: "numeric",
            month: "short",
        });
    }

    function getActionIcon(action: string) {
        switch (action) {
            case "CREATE":
                return "add_circle";
            case "UPDATE_STATUS":
                return "compare_arrows";
            case "COMMENT":
                return "chat_bubble";
            default:
                return "history_edu";
        }
    }

    function getActionColor(action: string) {
        switch (action) {
            case "CREATE":
                return "text-emerald-500 bg-emerald-50";
            case "UPDATE_STATUS":
                return "text-amber-500 bg-amber-50";
            case "COMMENT":
                return "text-primary bg-primary/10";
            default:
                return "text-slate-500 bg-slate-50";
        }
    }

    // Real comment submission
    async function submitComment() {
        if (!newComment.trim()) return;
        submitting = true;
        try {
            await taskService.addComment(taskId, newComment);
            newComment = "";
            await loadActivities(); // Reload to show the new comment
        } catch (e) {
            console.error(e);
        } finally {
            submitting = false;
        }
    }
</script>

<div
    class="flex flex-col h-full min-h-0 bg-slate-50/50 rounded-2xl border border-slate-200 overflow-hidden"
>
    {#if mode === "all"}
        <div class="px-5 py-4 border-b border-slate-200 bg-white">
            <h3
                class="text-sm font-bold text-slate-800 flex items-center gap-2"
            >
                <span
                    class="material-symbols-outlined text-[18px] text-slate-400"
                    >history</span
                >
                Activities & Comments
            </h3>
        </div>
    {/if}

    <!-- Activity Flow -->
    <div class="flex-1 overflow-y-auto p-5 space-y-5 custom-scrollbar relative">
        {#if loading}
            <div class="flex justify-center p-4">
                <div
                    class="animate-spin rounded-full h-6 w-6 border-b-2 border-primary"
                ></div>
            </div>
        {:else if filteredActivities.length === 0}
            <p class="text-xs text-center text-slate-400 mt-4">
                No {mode === "all"
                    ? "activities"
                    : mode === "activities"
                      ? "activities"
                      : "comments"} yet
            </p>
        {:else}
            <!-- Activity Timeline Line -->
            <div
                class="absolute left-9 inset-y-5 w-px bg-slate-200 -z-10"
            ></div>

            {#each filteredActivities as activity}
                <div class="flex gap-4">
                    <Avatar
                        name={String(activity.actor_id || "System")}
                        size="sm"
                        class="ring-4 ring-slate-50 flex-shrink-0"
                    />

                    <div
                        class="flex-1 bg-white rounded-2xl p-4 border border-slate-100 shadow-sm overflow-hidden"
                    >
                        <div
                            class="flex items-center justify-between gap-4 mb-2"
                        >
                            <span class="text-xs font-bold text-slate-800">
                                {activity.actor_id
                                    ? `User ${activity.actor_id}`
                                    : "System"}
                            </span>
                            <span
                                class="text-[10px] text-slate-400 font-medium whitespace-nowrap"
                                >{formatTime(activity.created_at)}</span
                            >
                        </div>

                        {#if activity.action === "COMMENT"}
                            <p
                                class="text-sm text-slate-600 leading-relaxed bg-slate-50 p-3 rounded-xl border border-slate-100"
                            >
                                {activity.new_value}
                            </p>
                        {:else}
                            <div
                                class="flex items-center gap-2 text-[13px] text-slate-600"
                            >
                                <span
                                    class="material-symbols-outlined text-[16px] {getActionColor(
                                        activity.action,
                                    ).split(' ')[0]}"
                                >
                                    {getActionIcon(activity.action)}
                                </span>
                                {#if activity.action === "CREATE"}
                                    <span>Created <strong>task</strong></span>
                                {:else if activity.action === "UPDATE_STATUS"}
                                    <span
                                        >Changed status to <strong
                                            >{activity.new_value}</strong
                                        ></span
                                    >
                                {:else}
                                    <span
                                        >Updated <strong
                                            >{activity.action}</strong
                                        ></span
                                    >
                                {/if}
                            </div>
                        {/if}
                    </div>
                </div>
            {/each}
        {/if}
    </div>

    {#if mode === "all" || mode === "comments"}
        <!-- Comment Box -->
        <div class="p-4 bg-white border-t border-slate-200 relative">
            <div class="flex gap-3">
                <Avatar name="Me" size="sm" />
                <div class="flex-1 relative">
                    <textarea
                        bind:value={newComment}
                        placeholder="Add a comment..."
                        class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 text-sm focus:bg-white focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-all resize-none h-12 min-h-[48px] placeholder:text-slate-400 custom-scrollbar pr-14"
                        onkeydown={(e) => {
                            if (e.key === "Enter" && !e.shiftKey) {
                                e.preventDefault();
                                submitComment();
                            }
                        }}
                    ></textarea>
                    <button
                        class="absolute bottom-2 right-2 p-1.5 rounded-lg text-primary hover:bg-primary/10 transition-colors disabled:opacity-50 {newComment.trim()
                            ? 'bg-primary/10'
                            : ''}"
                        title="Send (Enter)"
                        disabled={!newComment.trim() || submitting}
                        onclick={submitComment}
                    >
                        <span class="material-symbols-outlined text-[18px]"
                            >send</span
                        >
                    </button>
                </div>
            </div>
        </div>
    {/if}
</div>

<style>
    .custom-scrollbar::-webkit-scrollbar {
        width: 4px;
    }
    .custom-scrollbar::-webkit-scrollbar-track {
        background: transparent;
    }
    .custom-scrollbar::-webkit-scrollbar-thumb {
        background-color: rgba(203, 213, 225, 0.5);
        border-radius: 10px;
    }
    .custom-scrollbar:hover::-webkit-scrollbar-thumb {
        background-color: rgba(148, 163, 184, 0.8);
    }
</style>
