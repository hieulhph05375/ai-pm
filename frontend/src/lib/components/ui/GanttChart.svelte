<script lang="ts">
	import { tick } from "svelte";
	import type { GanttItem, GanttDependency } from "$lib/types/gantt";
	import type { Holiday } from "$lib/services/holidays";

	interface Props {
		items: GanttItem[];
		dependencies?: GanttDependency[];
		holidays?: Holiday[];
		baselineNodes?: any[];
		restDays?: number[];
		onScroll?: (e: Event) => void;
		viewMode?: "Day" | "Week" | "Month" | "Quarter";
		hoveredItemId?: string | number | null;
		onHover?: (id: string | number | null) => void;
		onProgressChange?: (item: GanttItem, progress: number) => void;
		onDatesChange?: (item: GanttItem, newStart: Date, newEnd: Date) => void;
		onCreateDependency?: (
			fromId: string | number,
			toId: string | number,
		) => void;
		onDeleteDependency?: (dep: GanttDependency) => void;
		scrollTop?: number;
	}

	let {
		items = [],
		dependencies = [],
		holidays = [],
		baselineNodes = [],
		restDays = [0, 6],
		viewMode = "Week",
		hoveredItemId = $bindable(null),
		scrollTop = $bindable(0),
		onHover,
		onProgressChange,
		onDatesChange,
		onCreateDependency,
		onDeleteDependency,
	}: Props = $props();

	let headerEl: HTMLElement | undefined = $state();
	let ganttBodyEl: HTMLElement | undefined = $state();
	let ganttBodyClientHeight = $state(0);
	let svgEl: SVGElement | undefined = $state();

	// Drag-to-progress state
	let draggingItemId = $state<string | number | null>(null);
	let dragProgress = $state<Record<string | number, number>>({});

	// Drag-bar (move dates) state
	let dragBarItemId = $state<string | number | null>(null);
	let dragBarOffset = $state<Record<string | number, number>>({});

	// Resize state
	let resizingItemId = $state<string | number | null>(null);
	let resizingType = $state<"left" | "right" | null>(null);
	let resizingOffset = $state<Record<string | number, number>>({});

	// Connect-mode state
	let connectingFromId = $state<string | number | null>(null);
	let connectMousePos = $state({ x: 0, y: 0 });

	// Row positions for SVG arrows (nodeId → y center px from top)
	const ROW_HEIGHT = 64;

	let ganttScrollLeft = $state(0);

	function handleBodyScroll(e: Event) {
		const target = e.target as HTMLElement;
		if (headerEl) headerEl.scrollLeft = target.scrollLeft;
		scrollTop = target.scrollTop;
		ganttScrollLeft = target.scrollLeft;
	}

	$effect(() => {
		if (ganttBodyEl && ganttBodyEl.scrollTop !== scrollTop) {
			ganttBodyEl.scrollTop = scrollTop;
		}
	});

	function normalizeToLocalMidnight(
		d: Date | string | undefined | null,
	): Date {
		if (!d) return new Date();
		const date =
			typeof d === "string" ? new Date(d.split("T")[0]) : new Date(d);
		return new Date(date.getFullYear(), date.getMonth(), date.getDate());
	}

	const today = normalizeToLocalMidnight(new Date());

	// todayLeftPx calculation removed because we are highlighting columns instead
	// let todayLeftPx = $derived(getPixelOffset(today) + 40 - ganttScrollLeft);

	function getTimelineDates(mode: string, itemsList: GanttItem[]) {
		let minDate = new Date(today);
		let maxDate = new Date(today);

		if (itemsList && itemsList.length > 0) {
			let minTime = Infinity;
			let maxTime = -Infinity;
			itemsList.forEach((n) => {
				if (n.startDate) {
					const s = normalizeToLocalMidnight(n.startDate).getTime();
					if (s < minTime) minTime = s;
				}
				if (n.dueDate) {
					const e = normalizeToLocalMidnight(n.dueDate).getTime();
					if (e > maxTime) maxTime = e;
				}
			});
			if (minTime !== Infinity) minDate = new Date(minTime);
			if (maxTime !== -Infinity) maxDate = new Date(maxTime);
		}

		const dates: Date[] = [];
		const start = normalizeToLocalMidnight(minDate);

		if (mode === "Day" || mode === "Week") {
			// Pad 3 days before minDate
			start.setDate(start.getDate() - 3);
			if (mode === "Week") {
				// Snap to previous Monday
				const day = start.getDay();
				const diff = start.getDate() - day + (day === 0 ? -6 : 1);
				start.setDate(diff);
			}
		} else if (mode === "Month") {
			// Pad 1 month before minDate
			start.setDate(1);
			start.setMonth(start.getMonth() - 1);
		} else {
			// Snap to quarter start
			start.setDate(1);
			start.setMonth(Math.floor(start.getMonth() / 3) * 3);
		}

		// Calculate total iterations needed to cover maxDate + padding
		let curr = new Date(start);
		const paddedMax = new Date(maxDate);
		if (mode === "Day")
			paddedMax.setDate(paddedMax.getDate() + 2); // Pad 2 days after
		else if (mode === "Week")
			paddedMax.setDate(paddedMax.getDate() + 7); // Pad 1 week after
		else if (mode === "Month") paddedMax.setMonth(paddedMax.getMonth() + 1); // Pad 1 month after
		// For Quarter, we don't add extra quarter padding beyond the snap

		while (
			curr <= paddedMax ||
			dates.length < (mode === "Quarter" ? 4 : 6)
		) {
			// Ensure at least 4 columns for Quarter, 20 for Week, 10 for others
			dates.push(new Date(curr));
			if (mode === "Day") curr.setDate(curr.getDate() + 1);
			else if (mode === "Week") curr.setDate(curr.getDate() + 7);
			else if (mode === "Month") curr.setMonth(curr.getMonth() + 1);
			else curr.setMonth(curr.getMonth() + 3);
		}

		return dates;
	}

	function getWeekNumber(d: Date) {
		const date = new Date(
			Date.UTC(d.getFullYear(), d.getMonth(), d.getDate()),
		);
		const dayNum = date.getUTCDay() || 7;
		date.setUTCDate(date.getUTCDate() + 4 - dayNum);
		const yearStart = new Date(Date.UTC(date.getUTCFullYear(), 0, 1));
		return Math.ceil(
			((date.getTime() - yearStart.getTime()) / 86400000 + 1) / 7,
		);
	}

	let timelineDates = $derived(getTimelineDates(viewMode, items));
	let ganttContainerWidth = $state(0);

	let UNIT_WIDTH = $derived.by(() => {
		const baseWidth = 40;
		if (ganttContainerWidth > 0 && timelineDates.length > 0) {
			const availableWidth = ganttContainerWidth - 40; // Subtract the first 40px column
			const minWidth = availableWidth / timelineDates.length;
			return Math.max(baseWidth, minWidth);
		}
		return baseWidth;
	});

	function getPixelOffset(date: Date) {
		const firstDate = timelineDates[0];
		if (!firstDate) return 0;

		if (viewMode === "Day") {
			const diffMs = date.getTime() - firstDate.getTime();
			return (diffMs / (1000 * 60 * 60 * 24)) * UNIT_WIDTH;
		}

		if (viewMode === "Week") {
			const diffMs = date.getTime() - firstDate.getTime();
			return (diffMs / (1000 * 60 * 60 * 24 * 7)) * UNIT_WIDTH;
		}

		if (viewMode === "Month") {
			const yearDiff = date.getFullYear() - firstDate.getFullYear();
			const monthDiff = date.getMonth() - firstDate.getMonth();
			const totalMonths = yearDiff * 12 + monthDiff;

			// Add fraction of the month
			const daysInMonth = new Date(
				date.getFullYear(),
				date.getMonth() + 1,
				0,
			).getDate();
			const dayFraction = (date.getDate() - 1) / daysInMonth;

			return (totalMonths + dayFraction) * UNIT_WIDTH;
		}

		if (viewMode === "Quarter") {
			const yearDiff = date.getFullYear() - firstDate.getFullYear();
			const monthDiff = date.getMonth() - firstDate.getMonth();
			const totalQuarters = (yearDiff * 12 + monthDiff) / 3;

			// Quarter fraction is harder, but this is close enough for Gantt
			return totalQuarters * UNIT_WIDTH;
		}

		return 0;
	}

	function scrollToToday() {
		if (!ganttBodyEl) return;
		const _mode = viewMode; // reactive Dependency trigger
		tick().then(() => {
			if (!ganttBodyEl) return;
			const offset = getPixelOffset(today) + 40; // match bar padding offset
			const containerWidth = ganttBodyEl.clientWidth;
			ganttBodyEl.scrollTo({
				left: Math.max(0, offset - containerWidth / 2 + UNIT_WIDTH / 2),
				behavior: "smooth",
			});
		});
	}

	$effect(() => {
		scrollToToday();
	});

	function getHolidayForDate(date: Date) {
		const list = Array.isArray(holidays) ? holidays : [];
		if (list.length === 0) return null;

		const yyyy = date.getFullYear();
		const mm = String(date.getMonth() + 1).padStart(2, "0");
		const dd = String(date.getDate()).padStart(2, "0");
		const dStr = `${yyyy}-${mm}-${dd}`;
		// Recurring: match Month-Day
		const monthDay = `${mm}-${dd}`;
		return (
			list.find((h) => {
				// Backend might return "2026-03-05T00:00:00Z", so we match the prefix
				const hDate = h.date.split("T")[0];
				if (h.is_recurring) return hDate.substring(5) === monthDay;
				return hDate === dStr;
			}) || null
		);
	}

	function isWeekend(date: Date) {
		const day = date.getDay();
		return restDays.includes(day);
	}

	function getHolidaysInRange(start: Date, end: Date) {
		const result: { date: Date; holiday: Holiday }[] = [];
		const curr = normalizeToLocalMidnight(start);
		const finish = normalizeToLocalMidnight(end);

		while (curr <= finish) {
			const hol = getHolidayForDate(curr);
			if (hol) {
				result.push({ date: new Date(curr), holiday: hol });
			}
			curr.setDate(curr.getDate() + 1);
		}
		return result;
	}

	function getBarPosition(item: GanttItem) {
		if (!item.startDate || !item.dueDate) {
			return {
				left: "0px",
				width: "0px",
				visible: false,
				leftNum: 0,
				widthNum: 0,
			};
		}
		const start = normalizeToLocalMidnight(item.startDate);
		const end = normalizeToLocalMidnight(item.dueDate);

		let pixelOffset = getPixelOffset(start);
		let pixelWidth = getPixelOffset(end) - pixelOffset;

		if (resizingItemId === item.id) {
			const offset = resizingOffset[item.id] || 0;
			if (resizingType === "left") {
				pixelOffset += offset;
				pixelWidth -= offset;
			} else {
				pixelWidth += offset;
			}
		}

		pixelOffset += 40;
		const w = Math.max(pixelWidth, 20); // allow smaller width when resizing
		return {
			left: `${pixelOffset}px`,
			width: `${w}px`,
			leftNum: pixelOffset,
			widthNum: w,
			visible: true,
		};
	}

	function getBaselinePosition(baselineItem: any) {
		if (
			!baselineItem.planned_start_date ||
			!baselineItem.planned_end_date
		) {
			return {
				left: "0px",
				width: "0px",
				visible: false,
			};
		}
		const start = normalizeToLocalMidnight(baselineItem.planned_start_date);
		const end = normalizeToLocalMidnight(baselineItem.planned_end_date);

		let pixelOffset = getPixelOffset(start);
		let pixelWidth = getPixelOffset(end) - pixelOffset;

		pixelOffset += 40;
		const w = Math.max(pixelWidth, 20);
		return {
			left: `${pixelOffset}px`,
			width: `${w}px`,
			visible: true,
		};
	}

	// Total content height for sizing background columns
	let contentHeight = $derived(items.length * ROW_HEIGHT);
	// Full column height: at least contentHeight+128, or fill the visible container if taller
	let columnHeight = $derived(
		Math.max(contentHeight + 128, ganttBodyClientHeight),
	);
	// Total content width for SVG
	let svgWidth = $derived(timelineDates.length * UNIT_WIDTH + 200);
	let svgHeight = $derived(contentHeight + 20);

	// Compute bar positions map for arrow drawing
	let barPositions = $derived.by(() => {
		const map: Record<
			string | number,
			{ x1: number; x2: number; y: number }
		> = {};
		items.forEach((item, idx) => {
			const pos = getBarPosition(item);
			if (pos.visible) {
				map[item.id] = {
					x1: pos.leftNum,
					x2: pos.leftNum + pos.widthNum,
					y: idx * ROW_HEIGHT + ROW_HEIGHT / 2,
				};
			}
		});
		return map;
	});

	function formatDate(date: Date) {
		if (viewMode === "Day" || viewMode === "Week") {
			return date.toLocaleDateString("en-US", {
				month: "short",
				day: "numeric",
			});
		}
		return date.toLocaleDateString("en-US", {
			month: "short",
			year: "2-digit",
		});
	}

	function isCurrentPeriod(date: Date) {
		if (viewMode === "Day") {
			return date.toDateString() === today.toDateString();
		}
		const dStart = date.getTime();
		let dEnd = dStart;
		if (viewMode === "Week") {
			dEnd += 7 * 24 * 60 * 60 * 1000;
		} else if (viewMode === "Month") {
			const next = new Date(date);
			next.setMonth(next.getMonth() + 1);
			dEnd = next.getTime();
		} else if (viewMode === "Quarter") {
			const next = new Date(date);
			next.setMonth(next.getMonth() + 3);
			dEnd = next.getTime();
		}
		const tTime = today.getTime();
		return tTime >= dStart && tTime < dEnd;
	}

	function getDisplayProgress(item: GanttItem): number {
		if (draggingItemId === item.id && dragProgress[item.id] !== undefined) {
			return dragProgress[item.id];
		}
		return item.progress;
	}

	function getTimeWindowMs(mode: string) {
		if (mode === "Day") return 1000 * 60 * 60 * 24;
		if (mode === "Week") return 1000 * 60 * 60 * 24 * 7;
		if (mode === "Month") return 1000 * 60 * 60 * 24 * 30;
		return 1000 * 60 * 60 * 24 * 90;
	}

	let todayPx = $derived.by(() => {
		const dayMs = 1000 * 60 * 60 * 24;
		const offset = getPixelOffset(today);
		let width = UNIT_WIDTH;
		if (viewMode === "Week") width = UNIT_WIDTH / 7;
		else if (viewMode === "Month")
			width = UNIT_WIDTH / 30; // Approx
		else if (viewMode === "Quarter") width = UNIT_WIDTH / 90; // Approx

		// For more accuracy in Month/Quarter, we could calculate next day's offset
		// but simple division is usually enough for a visual indicator
		return { left: offset + 40, width };
	});

	// --- Drag to move gantt bar ---
	function startDragBar(e: MouseEvent, item: GanttItem) {
		if (connectingFromId !== null) return;
		// Don't drag bar if clicking on handles or interactive elements inside
		const target = e.target as Element;
		if (target.closest(".drag-handle, .cursor-crosshair")) return;

		e.preventDefault();
		dragBarItemId = item.id;
		dragBarOffset[item.id] = 0;
		const startX = e.clientX;
		const initialStart = new Date(item.startDate!);
		const initialEnd = new Date(item.dueDate!);

		function onMove(me: MouseEvent) {
			dragBarOffset = {
				...dragBarOffset,
				[item.id]: me.clientX - startX,
			};
		}

		function onUp() {
			const dx = dragBarOffset[item.id] ?? 0;
			dragBarItemId = null;
			document.removeEventListener("mousemove", onMove);
			document.removeEventListener("mouseup", onUp);
			document.body.style.cursor = "";

			if (dx !== 0 && Math.abs(dx) > 2) {
				const msDelta = (dx / UNIT_WIDTH) * getTimeWindowMs(viewMode);
				const newStart = new Date(initialStart.getTime() + msDelta);
				const newEnd = new Date(initialEnd.getTime() + msDelta);
				onDatesChange?.(item, newStart, newEnd);
			}
			dragBarOffset = { ...dragBarOffset, [item.id]: 0 };
		}

		document.body.style.cursor = "grab";
		document.addEventListener("mousemove", onMove);
		document.addEventListener("mouseup", onUp);
	}

	// --- Resize bar ---
	function startResize(
		e: MouseEvent,
		item: GanttItem,
		side: "left" | "right",
	) {
		e.stopPropagation();
		e.preventDefault();
		resizingItemId = item.id;
		resizingType = side;
		resizingOffset[item.id] = 0;
		const startX = e.clientX;
		const initialStart = new Date(item.startDate!);
		const initialEnd = new Date(item.dueDate!);

		function onMove(me: MouseEvent) {
			resizingOffset = {
				...resizingOffset,
				[item.id]: me.clientX - startX,
			};
		}

		function onUp() {
			const dx = resizingOffset[item.id] ?? 0;
			resizingItemId = null;
			resizingType = null;
			document.removeEventListener("mousemove", onMove);
			document.removeEventListener("mouseup", onUp);
			document.body.style.cursor = "";

			if (dx !== 0 && Math.abs(dx) > 2) {
				const msDelta = (dx / UNIT_WIDTH) * getTimeWindowMs(viewMode);
				let newStart = initialStart;
				let newEnd = initialEnd;

				if (side === "left") {
					newStart = new Date(initialStart.getTime() + msDelta);
				} else {
					newEnd = new Date(initialEnd.getTime() + msDelta);
				}

				if (newStart <= newEnd) {
					onDatesChange?.(item, newStart, newEnd);
				}
			}
			resizingOffset = { ...resizingOffset, [item.id]: 0 };
		}

		document.body.style.cursor = "ew-resize";
		document.addEventListener("mousemove", onMove);
		document.addEventListener("mouseup", onUp);
	}

	// --- Drag to set progress ---
	function startDrag(e: MouseEvent, item: GanttItem, barWidth: number) {
		if (connectingFromId !== null) return; // don't drag in connect mode
		e.preventDefault();
		e.stopPropagation();
		draggingItemId = item.id;
		dragProgress[item.id] = item.progress;
		const target = e.currentTarget as HTMLElement;
		const barElement = target.closest(".gantt-bar") || target.parentElement;
		if (!barElement) return;
		const rect = barElement.getBoundingClientRect();

		function onMove(me: MouseEvent) {
			const offset = me.clientX - rect.left;
			const newProgress = Math.round(
				Math.min(100, Math.max(0, (offset / barWidth) * 100)),
			);
			dragProgress = { ...dragProgress, [item.id]: newProgress };
		}
		function onUp() {
			const finalProgress = dragProgress[item.id] ?? item.progress;
			draggingItemId = null;
			document.removeEventListener("mousemove", onMove);
			document.removeEventListener("mouseup", onUp);
			document.body.style.cursor = "";
			if (finalProgress !== item.progress) {
				onProgressChange?.(item, finalProgress);
			}
		}
		document.body.style.cursor = "ew-resize";
		document.addEventListener("mousemove", onMove);
		document.addEventListener("mouseup", onUp);
	}

	// --- Connect mode (Drag & Drop) ---
	function startConnect(e: MouseEvent, itemId: string | number) {
		e.stopPropagation();
		e.preventDefault();
		connectingFromId = itemId;

		if (svgEl) {
			const rect = svgEl.getBoundingClientRect();
			connectMousePos = {
				x: e.clientX - rect.left,
				y: e.clientY - rect.top,
			};
		}

		function onMove(me: MouseEvent) {
			if (svgEl) {
				const rect = svgEl.getBoundingClientRect();
				connectMousePos = {
					x: me.clientX - rect.left,
					y: me.clientY - rect.top,
				};
			}
		}

		function onUp(me: MouseEvent) {
			document.removeEventListener("mousemove", onMove);
			document.removeEventListener("mouseup", onUp);

			// Find drop element
			const previewEl = document.getElementById("preview-line");
			if (previewEl) previewEl.style.pointerEvents = "none"; // so we don't hover it

			const elem = document.elementFromPoint(me.clientX, me.clientY);
			const targetBar = elem?.closest(".gantt-bar");
			if (targetBar) {
				const targetId = Number(targetBar.getAttribute("data-node-id"));
				if (targetId && targetId !== connectingFromId) {
					const exists = dependencies.some(
						(d) =>
							d.fromId === connectingFromId &&
							d.toId === targetId,
					);
					if (!exists) {
						onCreateDependency?.(connectingFromId!, targetId);
					}
				}
			}

			connectingFromId = null;
		}

		document.addEventListener("mousemove", onMove);
		document.addEventListener("mouseup", onUp);
	}

	function arrowPath(
		pred: { x1: number; x2: number; y: number },
		succ: { x1: number; x2: number; y: number },
	): string {
		const x1 = pred.x2 + 9; // center of right dot
		const y1 = pred.y;
		const x2 = succ.x1 - 14; // edge of left dot for arrowhead to fit
		const y2 = succ.y;
		const cx = (x1 + x2) / 2;
		return `M ${x1} ${y1} C ${cx} ${y1}, ${cx} ${y2}, ${x2} ${y2}`;
	}
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div class="flex-1 flex flex-col overflow-hidden bg-white/50 relative">
	<!-- Timeline Header -->
	<div
		class="h-12 border-b border-slate-200/50 flex flex-none items-center overflow-x-hidden no-scrollbar gantt-grid bg-white/50"
		bind:this={headerEl}
		bind:clientWidth={ganttContainerWidth}
	>
		<div
			class="flex-none flex items-center h-full relative"
			style="width: {svgWidth}px"
		>
			<div
				class="w-10 flex-none h-full border-r border-slate-200/50 bg-slate-50/30"
			></div>
			{#each timelineDates as date}
				{@const isDayView = viewMode === "Day"}
				{@const hol = isDayView ? getHolidayForDate(date) : null}
				{@const weekend = isDayView ? isWeekend(date) : false}
				<div
					class="flex-none text-center text-[10px] font-bold border-r border-slate-100 uppercase h-full flex flex-col items-center justify-center leading-tight gap-0.5
					{hol ? (hol.type === 'state' ? 'bg-rose-500/10' : 'bg-sky-500/10') : ''}
					{weekend ? 'bg-slate-500/5' : ''} text-slate-400"
					style="width: {UNIT_WIDTH}px"
				>
					{#if viewMode === "Week"}
						<span>W{getWeekNumber(date)}</span>
					{:else if viewMode === "Quarter"}
						<span
							>Q{Math.floor(date.getMonth() / 3) + 1}
							{date.getFullYear()}</span
						>
					{:else if viewMode === "Day"}
						<span class="text-[10px] font-black"
							>{date.getDate()}</span
						>
						<span class="text-[9px] text-slate-400"
							>{date.getMonth() + 1}</span
						>
						<span class="text-[7px] opacity-80"
							>{date.toLocaleDateString("en-US", {
								weekday: "short",
							})}</span
						>
					{:else}
						<span>{formatDate(date)}</span>
					{/if}
				</div>
			{/each}

			<!-- Today Highlight in Header -->
			<div
				class="absolute top-0 h-full bg-primary/[0.04] border-x border-primary/20 z-20 pointer-events-none"
				style="left: {todayPx.left}px; width: {todayPx.width}px;"
			></div>
		</div>
	</div>

	<!-- Main Timeline Grid -->
	<div
		class="flex-1 overflow-auto gantt-grid relative"
		onscroll={handleBodyScroll}
		bind:this={ganttBodyEl}
		bind:clientHeight={ganttBodyClientHeight}
	>
		<div class="relative w-fit min-w-full">
			<!-- Today Highlight Strip (Full Height) -->
			<div
				class="absolute top-0 bg-primary/[0.03] border-x border-primary/10 z-30 pointer-events-none"
				style="left: {todayPx.left}px; width: {todayPx.width}px; height: {columnHeight}px"
			></div>

			<!-- Holiday/Weekend Background Columns -->
			<div
				class="absolute top-0 left-0 flex pointer-events-none z-0"
				style="width: {svgWidth}px; height: {columnHeight}px"
			>
				<div
					class="w-10 flex-none h-full border-r border-slate-200/50 bg-slate-50/30"
				></div>
				{#each timelineDates as date}
					{@const isDayView = viewMode === "Day"}
					{@const hol = isDayView ? getHolidayForDate(date) : null}
					{@const weekend = isDayView ? isWeekend(date) : false}
					<div
						class="flex-none h-full border-r border-slate-100/50 relative
					{hol ? (hol.type === 'state' ? 'bg-rose-500/15' : 'bg-sky-500/15') : ''}
					{weekend ? 'bg-slate-500/10' : ''}"
						style="width: {UNIT_WIDTH}px"
					>
						{#if hol}
							<div
								class="sticky top-1 text-[8px] font-bold text-center {hol.type ===
								'state'
									? 'text-rose-400'
									: 'text-sky-400'} uppercase leading-none px-0.5"
							>
								{hol.name}
							</div>
						{/if}
					</div>
				{/each}
			</div>

			<!-- SVG layer for dependency arrows -->
			<svg
				bind:this={svgEl}
				class="absolute top-0 left-0 pointer-events-none z-10"
				width={svgWidth}
				height={svgHeight}
				style="min-width: {svgWidth}px"
			>
				<defs>
					<marker
						id="arrowhead"
						markerWidth="8"
						markerHeight="6"
						refX="6"
						refY="3"
						orient="auto"
					>
						<polygon
							points="0 0, 8 3, 0 6"
							fill="#6366f1"
							opacity="0.8"
						/>
					</marker>
					<marker
						id="arrowhead-del"
						markerWidth="8"
						markerHeight="6"
						refX="6"
						refY="3"
						orient="auto"
					>
						<polygon
							points="0 0, 8 3, 0 6"
							fill="#ef4444"
							opacity="0.8"
						/>
					</marker>
				</defs>
				{#each dependencies as dep}
					{@const pred = barPositions[dep.fromId]}
					{@const succ = barPositions[dep.toId]}
					{#if pred && succ}
						<path
							d={arrowPath(pred, succ)}
							fill="none"
							stroke="#6366f1"
							stroke-width="1.5"
							stroke-dasharray="4 2"
							opacity="0.7"
							marker-end="url(#arrowhead)"
							class="dep-line pointer-events-auto cursor-pointer hover:stroke-red-500 hover:opacity-100 transition-all"
							onclick={(e) => {
								e.stopPropagation();
								onDeleteDependency?.(dep);
							}}
						>
							<title>Click to remove connection</title>
						</path>
					{/if}
				{/each}

				{#if connectingFromId !== null && barPositions[connectingFromId]}
					{@const pred = barPositions[connectingFromId]}
					<path
						id="preview-line"
						d={arrowPath(pred, {
							x1: connectMousePos.x + 14,
							x2: connectMousePos.x + 14,
							y: connectMousePos.y,
						})}
						fill="none"
						stroke="#6366f1"
						stroke-width="1.5"
						stroke-dasharray="4 2"
						opacity="1"
						marker-end="url(#arrowhead)"
						class="pointer-events-none"
					/>
				{/if}
			</svg>

			<div
				class="flex flex-col relative z-20 pb-32"
				style="width: {svgWidth}px"
			>
				{#each items as item, i (item.id)}
					{@const pos = getBarPosition(item)}
					{@const displayProgress = getDisplayProgress(item)}
					{@const baselineNode = baselineNodes.find(
						(b) => b.node_id === item.id,
					)}
					{@const baselinePos = baselineNode
						? getBaselinePosition(baselineNode)
						: { visible: false }}

					<div
						role="row"
						tabindex="0"
						class="h-[64px] min-h-[64px] max-h-[64px] flex-none flex items-center relative border-b border-slate-100 group overflow-visible transition-all duration-200 gantt-grid"
						style="width: {svgWidth}px; background-color: {hoveredItemId ===
						item.id
							? 'rgba(19, 55, 236, 0.08)'
							: 'transparent'}"
						onmouseenter={() => onHover?.(item.id)}
						onmouseleave={() => onHover?.(null)}
					>
						{#if (draggingItemId === item.id || resizingItemId === item.id || dragBarItemId === item.id) && pos.visible}
							<!-- Ghost Bar (Original Reference) -->
							<div
								class="absolute top-1/2 -translate-y-1/2 h-[22px] rounded-lg border-2 border-dashed border-primary/20 bg-primary/5 pointer-events-none z-10"
								style="left: {pos.left}; width: {pos.width}"
							></div>
						{/if}

						<!-- Render Baseline Bar Behind Actual Bar -->
						{#if baselinePos && "left" in baselinePos && baselinePos.visible}
							<div
								class="absolute top-[40%] translate-y-1 h-[8px] rounded-sm bg-slate-300 pointer-events-none z-10 opacity-70"
								style="left: {baselinePos.left}; width: {baselinePos.width}"
								title="Baseline: {new Date(
									baselineNode?.planned_start_date || '',
								).toLocaleDateString('en-US')} - {new Date(
									baselineNode?.planned_end_date || '',
								).toLocaleDateString('en-US')}"
							></div>
						{/if}

						{#if pos.visible}
							{@const itemHolidays = getHolidaysInRange(
								new Date(item.startDate!),
								new Date(item.dueDate!),
							)}
							<div
								role="button"
								tabindex="0"
								aria-label="Gantt Bar for {item.title}"
								class="gantt-bar absolute top-1/2 -translate-y-1/2 h-[22px] rounded-lg shadow-lg flex items-center px-4 select-none transition-all duration-200
							{item.progress === 100
									? 'shadow-indigo-500/20'
									: item.progress > 0
										? 'shadow-amber-500/20'
										: 'bg-slate-100 border border-slate-200'}
							{itemHolidays.length > 0 ? '!border-rose-300 ring-2 ring-rose-500/10' : ''}
							{draggingItemId === item.id ||
								resizingItemId === item.id ||
								dragBarItemId === item.id
									? 'scale-[1.02] shadow-xl z-50 ring-2 ring-primary/20'
									: ''}
							{draggingItemId === item.id
									? 'cursor-ew-resize'
									: dragBarItemId === item.id
										? 'cursor-grabbing'
										: 'cursor-default'}
							hover:cursor-grab outline-none focus:ring-2 focus:ring-primary/40"
								style="left: calc({pos.left} + {dragBarOffset[
									item.id
								] || 0}px); width: {pos.width}"
								data-node-id={item.id}
								onmousedown={(e) => startDragBar(e, item)}
								onkeydown={(e) => {
									if (e.key === "Enter" || e.key === " ") {
										e.preventDefault();
										// No specific action for Enter yet, but makes it focusable
									}
								}}
							>
								{#if itemHolidays.length > 0}
									<div
										class="absolute -top-[16px] left-1/2 -translate-x-1/2 flex items-center gap-1 bg-white shadow-xl border border-rose-100 rounded-full px-2 py-0.5 z-40"
									>
										<span
											style="font-size: 12px"
											class="material-symbols-outlined text-[10px] text-rose-500 font-bold"
											>warning</span
										>
										<span
											class="text-[7px] font-black text-rose-600 uppercase whitespace-nowrap"
											>Overlaps {itemHolidays.length} holidays</span
										>
									</div>
								{/if}
								{#if displayProgress > 0}
									<div
										class="absolute inset-y-0 left-0 bg-gradient-to-r {displayProgress ===
										100
											? 'from-indigo-500 to-sky-400'
											: 'from-amber-400 to-orange-500'} rounded-l-lg {displayProgress ===
										100
											? 'rounded-r-lg'
											: ''} transition-all max-w-full overflow-hidden"
										style="width: {displayProgress}%"
									></div>
									<span
										class="text-[10px] {displayProgress < 15
											? 'text-slate-600 pl-2'
											: 'text-white'} font-black tracking-tight relative z-10 whitespace-nowrap pointer-events-none"
									>
										{displayProgress}%
									</span>
								{:else}
									<span
										class="text-[10px] text-slate-400 font-bold tracking-tight whitespace-nowrap pointer-events-none"
									>
										{item.type === "milestone"
											? "MILESTONE"
											: "NOT STARTED"}
									</span>
								{/if}

								<!-- Resize Handles -->
								{#if item.type !== "milestone" && connectingFromId === null}
									<!-- Left Resize Handle -->
									<div
										role="button"
										tabindex="0"
										aria-label="Resize left"
										class="absolute left-0 top-0 bottom-0 w-2 cursor-ew-resize hover:bg-primary/20 z-30 outline-none"
										onmousedown={(e) =>
											startResize(e, item, "left")}
									></div>
									<!-- Right Resize Handle -->
									<div
										role="button"
										tabindex="0"
										aria-label="Resize right"
										class="absolute right-0 top-0 bottom-0 w-2 cursor-ew-resize hover:bg-primary/20 z-30 outline-none"
										onmousedown={(e) =>
											startResize(e, item, "right")}
									></div>
								{/if}

								<!-- Drag handle (progress edge) -->
								{#if connectingFromId === null}
									<div
										role="button"
										tabindex="0"
										aria-label="Update progress"
										class="drag-handle absolute top-0 bottom-0 flex items-center justify-center cursor-ew-resize z-40 group/drag {draggingItemId ===
										item.id
											? 'opacity-100 scale-110'
											: 'opacity-0 group-hover:opacity-100 hover:scale-110'} transition-all duration-200 outline-none"
										style="left: max(-8px, calc({displayProgress}% - 8px)); width: 16px;"
										onmousedown={(e) =>
											startDrag(e, item, pos.widthNum)}
										title="Drag along chart to update progress"
									>
										<div
											class="w-2.5 h-full rounded bg-white shadow-sm border border-slate-200 flex flex-col justify-center items-center gap-[2px] cursor-ew-resize"
										>
											<div
												class="w-[1px] h-2 bg-slate-300"
											></div>
											<div
												class="w-[1px] h-2 bg-slate-300"
											></div>
										</div>
									</div>
								{/if}

								<!-- Connect endpoints (left & right dots) -->
								{#if onCreateDependency}
									<!-- Left Connect Dot -->
									<div
										role="button"
										tabindex="0"
										aria-label="Link from start"
										class="absolute -left-[14px] top-1/2 -translate-y-1/2 size-[10px] rounded-full border-[2px] border-primary/40 bg-white hover:border-primary hover:bg-primary/10 hover:scale-125 transition-all z-30 shadow-sm cursor-crosshair pointer-events-auto outline-none focus:ring-2 focus:ring-primary"
										onmousedown={(e) => {
											e.preventDefault();
											e.stopPropagation();
											startConnect(e, item.id);
										}}
										title="Drag to link"
									></div>

									<!-- Right Connect Dot -->
									<div
										role="button"
										tabindex="0"
										aria-label="Link from end"
										class="absolute -right-[14px] top-1/2 -translate-y-1/2 size-[10px] rounded-full border-[2px] border-primary/40 bg-white hover:border-primary hover:bg-primary/10 hover:scale-125 transition-all z-30 shadow-sm cursor-crosshair pointer-events-auto outline-none focus:ring-2 focus:ring-primary"
										onmousedown={(e) => {
											e.preventDefault();
											e.stopPropagation();
											startConnect(e, item.id);
										}}
										title="Drag to link"
									></div>
								{/if}
							</div>
						{:else}
							<div
								class="absolute left-20 top-1/2 -translate-y-1/2 h-7 w-40 rounded-lg border-2 border-dashed border-slate-100 flex items-center px-4 opacity-50"
							>
								<span
									class="text-[10px] text-slate-300 font-bold tracking-tight"
									>NO DATES SET</span
								>
							</div>
						{/if}

						<!-- Row Hover Highlight -->
						<div
							class="absolute inset-0 bg-primary/0 group-hover:bg-primary/5 pointer-events-none transition-colors"
						></div>
					</div>
				{/each}
			</div>
		</div>
	</div>
</div>

<style>
	.gantt-grid {
		background-size: 50px 100%;
	}
	.no-scrollbar::-webkit-scrollbar {
		display: none;
	}
	.no-scrollbar {
		-ms-overflow-style: none;
		scrollbar-width: none;
	}
	.dep-line:hover {
		stroke: #ef4444;
	}
</style>
