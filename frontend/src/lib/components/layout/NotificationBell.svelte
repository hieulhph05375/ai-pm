<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import {
        notificationService,
        type Notification,
    } from "$lib/services/notifications";

    let unreadCount = $state(0);
    let notifications = $state<Notification[]>([]);
    let open = $state(false);
    let loading = $state(false);
    let dropdownRef: HTMLDivElement | undefined = $state();
    let pollInterval: ReturnType<typeof setInterval>;

    async function loadUnreadCount() {
        try {
            unreadCount = await notificationService.getUnreadCount();
        } catch (e) {
            // Silently fail - non-critical
        }
    }

    async function loadNotifications() {
        loading = true;
        try {
            const res = await notificationService.list(1, 20);
            notifications = res.items ?? [];
            unreadCount = res.unread_count ?? 0;
        } catch (e) {
            notifications = [];
        } finally {
            loading = false;
        }
    }

    async function toggleOpen() {
        open = !open;
        if (open) {
            await loadNotifications();
        }
    }

    async function markRead(id: number) {
        await notificationService.markRead(id);
        notifications = notifications.map((n) =>
            n.id === id ? { ...n, is_read: true } : n,
        );
        unreadCount = Math.max(0, unreadCount - 1);
    }

    async function markAllRead() {
        await notificationService.markAllRead();
        notifications = notifications.map((n) => ({ ...n, is_read: true }));
        unreadCount = 0;
    }

    function getTypeIcon(type: string): string {
        const map: Record<string, string> = {
            TIMESHEET_REMINDER: "schedule",
            EFFORT_OVERRUN: "warning",
            DEADLINE_SOON: "event",
            ISSUE_STALE: "bug_report",
        };
        return map[type] ?? "notifications";
    }

    function getTypeColor(type: string): string {
        const map: Record<string, string> = {
            TIMESHEET_REMINDER: "text-blue-500 bg-blue-50",
            EFFORT_OVERRUN: "text-rose-500 bg-rose-50",
            DEADLINE_SOON: "text-amber-500 bg-amber-50",
            ISSUE_STALE: "text-purple-500 bg-purple-50",
        };
        return map[type] ?? "text-slate-500 bg-slate-50";
    }

    function formatRelative(dateStr: string): string {
        const diff = Date.now() - new Date(dateStr).getTime();
        const mins = Math.floor(diff / 60000);
        if (mins < 1) return "Just now";
        if (mins < 60) return `${mins}m ago`;
        const hrs = Math.floor(mins / 60);
        if (hrs < 24) return `${hrs}h ago`;
        return `${Math.floor(hrs / 24)}d ago`;
    }

    function handleClickOutside(e: MouseEvent) {
        if (dropdownRef && !dropdownRef.contains(e.target as Node)) {
            open = false;
        }
    }

    onMount(() => {
        loadUnreadCount();
        // Poll every 30s
        pollInterval = setInterval(loadUnreadCount, 30_000);
        document.addEventListener("click", handleClickOutside, true);
    });

    onDestroy(() => {
        clearInterval(pollInterval);
        document.removeEventListener("click", handleClickOutside, true);
    });
</script>

<div class="relative" bind:this={dropdownRef}>
    <!-- Bell Button -->
    <button
        class="relative size-10 flex items-center justify-center rounded-xl bg-white border border-slate-200 text-slate-600 hover:text-primary hover:border-primary/30 transition-all"
        onclick={toggleOpen}
        aria-label="Notifications"
    >
        <span class="material-symbols-outlined">notifications</span>
        {#if unreadCount > 0}
            <span
                class="absolute -top-1 -right-1 min-w-[18px] h-[18px] flex items-center justify-center bg-rose-500 text-white text-[9px] font-black rounded-full px-1 shadow-sm animate-bounce"
            >
                {unreadCount > 99 ? "99+" : unreadCount}
            </span>
        {/if}
    </button>

    <!-- Dropdown -->
    {#if open}
        <div
            class="absolute right-0 top-12 w-80 bg-white rounded-2xl shadow-2xl border border-slate-100 z-50 overflow-hidden"
        >
            <!-- Header -->
            <div
                class="px-4 py-3 border-b border-slate-100 flex items-center justify-between"
            >
                <div class="flex items-center gap-2">
                    <span class="text-sm font-bold text-slate-900"
                        >Notifications</span
                    >
                    {#if unreadCount > 0}
                        <span
                            class="text-[10px] font-black text-primary bg-primary/10 px-1.5 py-0.5 rounded-full"
                            >{unreadCount} new</span
                        >
                    {/if}
                </div>
                {#if unreadCount > 0}
                    <button
                        class="text-[11px] font-semibold text-primary hover:text-primary/70 transition-colors"
                        onclick={markAllRead}
                    >
                        Mark all read
                    </button>
                {/if}
            </div>

            <!-- Notification list -->
            <div class="max-h-80 overflow-y-auto divide-y divide-slate-50">
                {#if loading}
                    <div class="py-8 text-center text-slate-400 text-sm">
                        Loading...
                    </div>
                {:else if notifications.length === 0}
                    <div class="py-8 text-center">
                        <span
                            class="material-symbols-outlined text-slate-300 text-3xl"
                            >notifications_none</span
                        >
                        <p class="text-sm text-slate-400 mt-1">
                            No notifications
                        </p>
                    </div>
                {:else}
                    {#each notifications as notif (notif.id)}
                        <button
                            onclick={() => !notif.is_read && markRead(notif.id)}
                            class="w-full text-left px-4 py-3 flex gap-3 items-start transition-all {notif.is_read
                                ? 'opacity-60'
                                : 'hover:bg-slate-50 bg-primary/[0.02]'}"
                        >
                            <span
                                class="mt-0.5 size-7 flex-none flex items-center justify-center rounded-full text-[14px] {getTypeColor(
                                    notif.type,
                                )}"
                            >
                                <span
                                    class="material-symbols-outlined"
                                    style="font-size:14px"
                                    >{getTypeIcon(notif.type)}</span
                                >
                            </span>
                            <div class="flex-1 min-w-0">
                                <p
                                    class="text-sm font-semibold text-slate-900 truncate"
                                >
                                    {notif.title}
                                </p>
                                {#if notif.body}
                                    <p
                                        class="text-[11px] text-slate-500 mt-0.5 line-clamp-2"
                                    >
                                        {notif.body}
                                    </p>
                                {/if}
                                <p class="text-[10px] text-slate-400 mt-1">
                                    {formatRelative(notif.created_at)}
                                </p>
                            </div>
                            {#if !notif.is_read}
                                <span
                                    class="mt-2 size-2 flex-none rounded-full bg-primary"
                                ></span>
                            {/if}
                        </button>
                    {/each}
                {/if}
            </div>
        </div>
    {/if}
</div>
