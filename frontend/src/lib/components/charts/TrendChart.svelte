<script lang="ts">
    const formatDate = (d: string | Date, opts: Intl.DateTimeFormatOptions) =>
        new Intl.DateTimeFormat("en-US", opts).format(new Date(d));

    // Props
    export let data: { date: string; value: number; tooltip?: string }[] = [];
    export let title: string = "";
    export let color: "primary" | "secondary" | "emerald" | "rose" | "amber" =
        "primary";
    export let showBaseline: boolean = true;
    export let baselineValue: number = 1.0;

    // SVG Config
    const width = 400;
    const height = 150;
    const padding = { top: 20, right: 20, bottom: 30, left: 30 };
    const innerWidth = width - padding.left - padding.right;
    const innerHeight = height - padding.top - padding.bottom;

    // Derived state
    $: values = data.map((d) => d.value);
    $: minVal = showBaseline
        ? Math.min(...values, baselineValue - 0.2)
        : Math.min(...values);
    $: maxVal = showBaseline
        ? Math.max(...values, baselineValue + 0.2)
        : Math.max(...values);
    $: range = maxVal - minVal || 1;

    // Scale functions
    $: xScale = (index: number) =>
        padding.left + (index / Math.max(1, data.length - 1)) * innerWidth;
    $: yScale = (val: number) =>
        height - padding.bottom - ((val - minVal) / range) * innerHeight;

    // Path generator
    $: pathData =
        data.length > 0
            ? `M ${data.map((d, i) => `${xScale(i)},${yScale(d.value)}`).join(" L ")}`
            : "";

    const colorMap = {
        primary: "var(--color-primary, #0ea5e9)",
        secondary: "var(--color-secondary, #64748b)",
        emerald: "#10b981",
        rose: "#f43f5e",
        amber: "#f59e0b",
    };

    // Interaction State
    let hoveredPoint: number | null = null;
</script>

<div class="trend-chart-container relative group">
    {#if title}
        <h4 class="text-sm font-bold text-slate-700 mb-2">{title}</h4>
    {/if}

    {#if data.length === 0}
        <div
            class="h-[150px] flex items-center justify-center text-sm text-slate-400 bg-slate-50 rounded-xl w-[400px]"
        >
            No trend data available
        </div>
    {:else}
        <svg
            {width}
            {height}
            class="overflow-visible svg-container w-full h-auto max-w-full"
            viewBox="0 0 {width} {height}"
        >
            <!-- Axes context -->
            <line
                x1={padding.left}
                y1={height - padding.bottom}
                x2={width - padding.right}
                y2={height - padding.bottom}
                stroke="#e2e8f0"
                stroke-width="1"
            />

            <!-- Dates on X axis (start and end only if many, or all if few) -->
            {#if data.length > 0}
                <text
                    x={padding.left}
                    y={height - padding.bottom + 15}
                    fill="#94a3b8"
                    font-size="10"
                    dominant-baseline="hanging"
                    >{formatDate(data[0].date, {
                        month: "short",
                        day: "2-digit",
                    })}</text
                >
                {#if data.length > 1}
                    <text
                        x={width - padding.right}
                        y={height - padding.bottom + 15}
                        fill="#94a3b8"
                        font-size="10"
                        text-anchor="end"
                        dominant-baseline="hanging"
                        >{formatDate(data[data.length - 1].date, {
                            month: "short",
                            day: "2-digit",
                        })}</text
                    >
                {/if}
            {/if}

            <!-- Target Baseline -->
            {#if showBaseline}
                <line
                    x1={padding.left}
                    y1={yScale(baselineValue)}
                    x2={width - padding.right}
                    y2={yScale(baselineValue)}
                    stroke="#94a3b8"
                    stroke-width="1"
                    stroke-dasharray="4"
                />
                <text
                    x={padding.left - 5}
                    y={yScale(baselineValue)}
                    fill="#94a3b8"
                    font-size="10"
                    text-anchor="end"
                    dominant-baseline="middle">{baselineValue.toFixed(1)}</text
                >
            {/if}

            <!-- Trend Line -->
            <path
                d={pathData}
                fill="none"
                stroke={colorMap[color]}
                stroke-width="2.5"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="transition-all duration-500 ease-out"
            />

            <!-- Value Points -->
            {#each data as d, i}
                <!-- svelte-ignore a11y-mouse-events-have-key-events // purely visual hover -->
                <!-- svelte-ignore a11y-no-static-element-interactions -->
                <circle
                    cx={xScale(i)}
                    cy={yScale(d.value)}
                    r={hoveredPoint === i ? 5 : 3}
                    fill="white"
                    stroke={colorMap[color]}
                    stroke-width="2"
                    class="transition-all duration-200 cursor-pointer origin-center"
                    on:mouseover={() => (hoveredPoint = i)}
                    on:mouseleave={() => (hoveredPoint = null)}
                />
            {/each}

            <!-- Dynamic Tooltip drawn via SVG on hover -->
            {#if hoveredPoint !== null}
                {@const point = data[hoveredPoint]}
                {@const cx = xScale(hoveredPoint)}
                {@const cy = yScale(point.value)}
                {@const ttWidth = 100}
                {@const ttHeight = 40}
                {@const ttX = cx > width / 2 ? cx - ttWidth - 10 : cx + 10}
                {@const ttY = cy - ttHeight / 2}

                <g class="pointer-events-none fade-in">
                    <rect
                        x={ttX}
                        y={ttY}
                        width={ttWidth}
                        height={ttHeight}
                        fill="white"
                        rx="6"
                        stroke="#e2e8f0"
                        stroke-width="1"
                        filter="drop-shadow(0 4px 6px rgb(0 0 0 / 0.05))"
                    />
                    <text
                        x={ttX + 10}
                        y={ttY + 16}
                        font-size="11"
                        font-weight="bold"
                        fill="#334155"
                        >{formatDate(point.date, {
                            month: "short",
                            day: "2-digit",
                            year: "numeric",
                        })}</text
                    >
                    <text
                        x={ttX + 10}
                        y={ttY + 30}
                        font-size="12"
                        font-weight="bold"
                        fill={colorMap[color]}
                        >{point.tooltip || point.value.toFixed(2)}</text
                    >
                </g>
            {/if}
        </svg>
    {/if}
</div>

<style>
    .svg-container {
        font-family: inherit;
    }
    .fade-in {
        animation: fadeIn 0.15s ease-out forwards;
    }
    @keyframes fadeIn {
        from {
            opacity: 0;
            transform: translateY(2px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }
</style>
