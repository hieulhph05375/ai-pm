<script lang="ts">
    import { fade, fly } from "svelte/transition";
    import { quintOut } from "svelte/easing";

    let {
        show = $bindable(false),
        title = "",
        width = "max-w-xl",
        onClose,
        children,
        footer,
        tabs,
    }: {
        show: boolean;
        title: string | import("svelte").Snippet;
        width?: string;
        onClose?: () => void;
        children?: import("svelte").Snippet;
        footer?: import("svelte").Snippet;
        tabs?: import("svelte").Snippet;
    } = $props();

    function close() {
        show = false;
        if (onClose) onClose();
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === "Escape" && show) {
            close();
        }
    }
</script>

<svelte:window onkeydown={handleKeydown} />

{#if show}
    <!-- Backdrop -->
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
        class="fixed inset-0 z-50 bg-slate-900/40 backdrop-blur-sm"
        transition:fade={{ duration: 300 }}
        onclick={close}
    ></div>

    <!-- Drawer Content -->
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
        class="fixed top-0 right-0 h-full w-full {width} bg-white dark:bg-slate-900 shadow-2xl z-50 flex flex-col border-l border-white/20 dark:border-slate-700/30 overflow-hidden"
        transition:fly={{ x: 600, duration: 400, easing: quintOut }}
        onclick={(e) => e.stopPropagation()}
    >
        <!-- Header -->
        <div class="px-8 pt-10 pb-6 shrink-0 flex justify-between items-start">
            <div class="flex-1 min-w-0">
                {#if typeof title === "string"}
                    <h2
                        class="text-2xl font-outfit font-bold text-slate-900 dark:text-white leading-tight"
                    >
                        {title}
                    </h2>
                {:else}
                    {@render title()}
                {/if}
            </div>
            <button
                type="button"
                class="text-slate-400 hover:text-slate-600 transition-colors rounded-xl p-2 hover:bg-slate-100 flex items-center justify-center -mt-2 -mr-2"
                onclick={close}
            >
                <span class="material-symbols-outlined text-[20px]">close</span>
            </button>
        </div>

        {#if tabs}
            <div
                class="px-8 pb-4 shrink-0 border-b border-slate-100 dark:border-slate-800"
            >
                {@render tabs()}
            </div>
        {/if}

        <!-- Body -->
        <div
            class="px-8 flex-1 overflow-y-auto custom-scrollbar flex flex-col min-h-0 pb-10"
        >
            {#if children}
                {@render children()}
            {/if}
        </div>

        <!-- Footer -->
        {#if footer}
            <div
                class="px-8 py-6 border-t border-slate-100 dark:border-slate-800 shrink-0 flex gap-4 bg-white/50 dark:bg-slate-900/50 backdrop-blur-md"
            >
                {@render footer()}
            </div>
        {/if}
    </div>
{/if}

<style>
    .custom-scrollbar::-webkit-scrollbar {
        width: 4px;
    }
    .custom-scrollbar::-webkit-scrollbar-track {
        background: transparent;
    }
    .custom-scrollbar::-webkit-scrollbar-thumb {
        background-color: rgba(203, 213, 225, 0.4);
        border-radius: 20px;
    }
    .custom-scrollbar:hover::-webkit-scrollbar-thumb {
        background-color: rgba(148, 163, 184, 0.6);
    }
</style>
