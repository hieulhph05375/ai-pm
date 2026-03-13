<script lang="ts">
    interface Props {
        total: number;
        page: number;
        limit: number;
        onPageChange: (page: number) => void;
        onLimitChange: (limit: number) => void;
    }

    let {
        total = 0,
        page = 1,
        limit = 10,
        onPageChange,
        onLimitChange,
    }: Props = $props();

    const startEntry = $derived(total === 0 ? 0 : (page - 1) * limit + 1);
    const endEntry = $derived(Math.min(page * limit, total));
    const totalPages = $derived(Math.ceil(total / limit));

    function handleLimitChange(e: Event) {
        const newLimit = parseInt((e.target as HTMLSelectElement).value);
        onLimitChange(newLimit);
    }
</script>

<div
    class="p-8 flex items-center justify-between border-t border-slate-100 bg-slate-50/30"
>
    <div class="flex items-center gap-6">
        <p class="text-sm text-slate-500 font-medium whitespace-nowrap">
            Showing <span class="text-slate-900"
                >{startEntry} to {endEntry}</span
            >
            of
            <span class="text-slate-900">{total.toLocaleString()}</span> results
        </p>

        <div class="flex items-center gap-2">
            <span
                class="text-xs font-bold text-slate-400 uppercase tracking-widest"
                >Show</span
            >
            <select
                class="bg-white border border-slate-200 rounded-lg px-2 py-1 text-xs font-bold text-slate-600 focus:ring-2 focus:ring-primary/20 outline-none transition-all cursor-pointer"
                value={limit}
                onchange={handleLimitChange}
            >
                <option value={10}>10</option>
                <option value={20}>20</option>
                <option value={50}>50</option>
                <option value={100}>100</option>
                <option value={200}>200</option>
            </select>
            <span
                class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-0.5"
                >per page</span
            >
        </div>
    </div>

    {#if totalPages > 1}
        <div class="flex items-center gap-1.5">
            <button
                class="size-9 flex items-center justify-center rounded-lg border border-slate-200 text-slate-400 hover:bg-white transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                disabled={page === 1}
                onclick={() => onPageChange(page - 1)}
            >
                <span class="material-symbols-outlined text-[18px]"
                    >chevron_left</span
                >
            </button>

            {#each Array(Math.min(5, totalPages)) as _, i}
                <button
                    class="size-9 flex items-center justify-center rounded-lg font-bold text-sm transition-all {page ===
                    i + 1
                        ? 'bg-primary text-white shadow-sm border-transparent'
                        : 'border border-slate-200 text-slate-600 hover:bg-white'}"
                    onclick={() => onPageChange(i + 1)}
                >
                    {i + 1}
                </button>
            {/each}

            {#if totalPages > 5}
                <span class="px-2 text-slate-400 text-xs font-bold">...</span>
                <button
                    class="size-9 flex items-center justify-center rounded-lg border border-slate-200 text-slate-600 hover:bg-white font-bold text-sm transition-colors {page ===
                    totalPages
                        ? 'bg-primary text-white'
                        : ''}"
                    onclick={() => onPageChange(totalPages)}
                >
                    {totalPages}
                </button>
            {/if}

            <button
                class="size-9 flex items-center justify-center rounded-lg border border-slate-200 text-slate-600 hover:bg-white transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                disabled={page === totalPages}
                onclick={() => onPageChange(page + 1)}
            >
                <span class="material-symbols-outlined text-[18px]"
                    >chevron_right</span
                >
            </button>
        </div>
    {/if}
</div>
