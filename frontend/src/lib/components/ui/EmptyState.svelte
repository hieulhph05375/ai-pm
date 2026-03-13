<script lang="ts">
    interface Props {
        icon?: string;
        title?: string;
        message?: string;
        actionLabel?: string;
        actionIcon?: string;
        onaction?: () => void;
    }

    let {
        icon = "assignment",
        title = "No data available",
        message = "Start by adding new data.",
        actionLabel = "Add New",
        actionIcon = "add",
        onaction,
    }: Props = $props();
</script>

<div
    class="flex-1 flex flex-col items-center justify-center p-8 bg-slate-50/50 rounded-2xl w-full py-20 animate-in fade-in duration-500"
>
    <div class="flex flex-col items-center text-center max-w-md">
        <div class="relative mb-10 group">
            <!-- Decorative backgrounds -->
            <div
                class="size-32 bg-white rounded-3xl flex items-center justify-center rotate-6 translate-x-4 opacity-40 border border-slate-200 shadow-xl transition-transform group-hover:rotate-12 duration-500"
            ></div>
            <div
                class="size-32 bg-white rounded-3xl absolute inset-0 flex items-center justify-center -rotate-3 border border-slate-200 shadow-2xl transition-transform group-hover:-rotate-6 duration-500"
            >
                <div
                    class="size-20 rounded-2xl bg-primary/5 flex items-center justify-center"
                >
                    <span
                        class="material-symbols-outlined text-5xl text-primary/40 group-hover:scale-110 transition-transform duration-500"
                        >{icon}</span
                    >
                </div>
            </div>

            <!-- Action indicator -->
            <div
                class="absolute -bottom-2 -right-2 size-10 rounded-full bg-primary shadow-lg shadow-primary/30 flex items-center justify-center border-4 border-white animate-bounce-subtle"
            >
                <span
                    class="material-symbols-outlined text-white text-[18px] font-bold"
                    >{actionIcon}</span
                >
            </div>
        </div>

        <h3
            class="font-display font-bold text-2xl text-slate-900 mb-3 tracking-tight"
        >
            {title}
        </h3>
        <p class="text-slate-500 text-sm leading-relaxed mb-10 px-4">
            {message}
        </p>

        {#if onaction}
            <button
                onclick={onaction}
                class="bg-slate-900 text-white px-8 py-3.5 rounded-xl font-bold text-sm flex items-center gap-2 shadow-xl shadow-slate-900/10 hover:bg-primary transition-all transform hover:-translate-y-0.5 active:translate-y-0 active:scale-95"
            >
                <span class="material-symbols-outlined">{actionIcon}</span>
                {actionLabel}
            </button>
        {/if}
    </div>
</div>

<style>
    @keyframes bounce-subtle {
        0%,
        100% {
            transform: translateY(0);
        }
        50% {
            transform: translateY(-3px);
        }
    }
    .animate-bounce-subtle {
        animation: bounce-subtle 3s infinite ease-in-out;
    }
</style>
