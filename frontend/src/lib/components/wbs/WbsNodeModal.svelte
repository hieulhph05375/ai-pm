<script lang="ts">
	import Modal from "../ui/Modal.svelte";
	import Button from "../ui/Button.svelte";
	import Input from "../ui/Input.svelte";
	import type { WBSNode, WBSNodeType } from "$lib/services/wbs";
	import { categoryService, type Category } from "$lib/services/categories";
	import type { User } from "$lib/services/users";

	interface Props {
		show: boolean;
		projectId: number;
		users?: User[];
		node?: Partial<WBSNode> | null;
		parentPath?: string | null;
		parentNode?: WBSNode | null;
		onSave: (node: any) => Promise<void>;
	}

	let {
		show = $bindable(false),
		projectId,
		users = [],
		node = null,
		parentPath = null,
		parentNode = null,
		onSave,
	}: Props = $props();

	let title = $state("");
	let type = $state<WBSNodeType>("Task");
	let plannedStartDate = $state("");
	let plannedEndDate = $state("");
	let progress = $state(0);
	let estimatedEffort = $state<number>(0);
	let assignedTo = $state<number | null>(null);
	let loading = $state(false);
	let errors = $state<Record<string, string>>({});
	let wbsCategories = $state<Category[]>([]);
	let selectedTypeId = $state<number | null>(null);

	$effect(() => {
		if (show) {
			errors = {};
			if (node) {
				title = node.title || "";
				type = node.type || "Task";
				selectedTypeId = node.type_id || null;
				plannedStartDate = node.planned_start_date?.split("T")[0] || "";
				plannedEndDate = node.planned_end_date?.split("T")[0] || "";
				progress = node.progress || 0;
				estimatedEffort = node.estimated_effort ?? 0;
				assignedTo = node.assigned_to || null;
			} else {
				title = "";
				type = "Task";
				selectedTypeId = null;
				plannedStartDate =
					parentNode?.planned_start_date?.split("T")[0] || "";
				plannedEndDate =
					parentNode?.planned_end_date?.split("T")[0] || "";
				progress = 0;
				estimatedEffort = 0;
				assignedTo = null;
			}
			fetchCategories();
		}
	});

	async function fetchCategories() {
		try {
			const res = await categoryService.listCategories(
				1,
				100,
				"",
				undefined,
			);
			// Filter by type code 'WBS_NODE_TYPE' if possible, or just look for the type name
			// Since we don't have typeId in props, we'll check the type name or code in the returned data
			wbsCategories = res.data.filter(
				(c) => c.type?.code === "WBS_NODE_TYPE",
			);

			// If we still can't find them, fetch all types first
			if (wbsCategories.length === 0) {
				const typesRes = await categoryService.listTypes(1, 100);
				const wbsType = typesRes.data.find(
					(t) => t.code === "WBS_NODE_TYPE",
				);
				if (wbsType) {
					const catRes = await categoryService.listCategories(
						1,
						100,
						"",
						wbsType.id,
					);
					wbsCategories = catRes.data;
				}
			}

			// If editing and we have type_id, ensure it's set
			if (node?.type_id) {
				selectedTypeId = node.type_id;
			} else if (!selectedTypeId && wbsCategories.length > 0) {
				// Default to 'Work Package' if new task
				const defaultCat =
					wbsCategories.find((c) => c.name === "Work Package") ||
					wbsCategories[0];
				selectedTypeId = defaultCat.id;
			}
		} catch (e) {
			console.error("Error fetching WBS categories", e);
		}
	}

	function validate() {
		errors = {};
		if (!title.trim()) {
			errors.title = "Title cannot be empty";
		}
		if (!type) {
			errors.type = "Please select task type";
		}
		if (!plannedStartDate || !plannedEndDate) {
			errors.dates = "Please select both start and end dates";
		} else {
			if (new Date(plannedStartDate) > new Date(plannedEndDate)) {
				errors.dates = "End date must be after or equal to start date";
			}

			// Validate against parent node dates
			if (
				parentNode &&
				parentNode.planned_start_date &&
				parentNode.planned_end_date
			) {
				const pStart = new Date(parentNode.planned_start_date);
				const pEnd = new Date(parentNode.planned_end_date);
				const cStart = new Date(plannedStartDate);
				const cEnd = new Date(plannedEndDate);

				// Reset time to zero for date only comparison
				pStart.setHours(0, 0, 0, 0);
				pEnd.setHours(0, 0, 0, 0);
				cStart.setHours(0, 0, 0, 0);
				cEnd.setHours(0, 0, 0, 0);

				if (cStart < pStart || cStart > pEnd) {
					errors.dates = `Start date must be between ${pStart.toLocaleDateString("en-US")} - ${pEnd.toLocaleDateString("en-US")}`;
				} else if (cEnd < pStart || cEnd > pEnd) {
					errors.dates = `End date must be between ${pStart.toLocaleDateString("en-US")} - ${pEnd.toLocaleDateString("en-US")}`;
				}
			}
		}

		if (
			progress === null ||
			progress === undefined ||
			progress < 0 ||
			progress > 100
		) {
			errors.progress = "Progress must be between 0 and 100";
		}
		return Object.keys(errors).length === 0;
	}

	async function handleSubmit() {
		if (!validate()) return;
		loading = true;
		try {
			const data = {
				title,
				type,
				type_id: selectedTypeId,
				planned_start_date: plannedStartDate
					? plannedStartDate + "T00:00:00Z"
					: null,
				planned_end_date: plannedEndDate
					? plannedEndDate + "T00:00:00Z"
					: null,
				progress: Number(progress),
				estimated_effort: Number(estimatedEffort) || 0,
				assigned_to: assignedTo ? Number(assignedTo) : null,
				parent_path: parentPath,
			};
			await onSave(data);
			show = false;
		} finally {
			loading = false;
		}
	}
</script>

<Modal
	bind:show
	title={node?.id ? "Edit Task" : "Add New Task"}
	maxWidth="max-w-xl"
>
	<div class="space-y-6">
		<div class="space-y-2">
			<label
				class="text-xs font-bold {errors.title
					? 'text-rose-500'
					: 'text-slate-400'} uppercase tracking-widest"
				for="title"
			>
				Title {errors.title ? `— ${errors.title}` : ""}
			</label>
			<Input
				id="title"
				bind:value={title}
				placeholder="e.g., API Infrastructure Design"
				class={errors.title ? "border-rose-300 ring-rose-100" : ""}
			/>
		</div>

		<div class="grid grid-cols-2 gap-6">
			<div class="space-y-2">
				<label
					class="text-xs font-bold {errors.type
						? 'text-rose-500'
						: 'text-slate-400'} uppercase tracking-widest"
					for="type"
				>
					Type {errors.type ? `— ${errors.type}` : ""}
				</label>
				<select
					id="type"
					bind:value={selectedTypeId}
					class="w-full bg-slate-50 border {errors.type
						? 'border-rose-200 ring-2 ring-rose-50'
						: 'border-slate-200'} rounded-xl px-4 py-2 text-sm focus:ring-2 focus:ring-primary/20 outline-none transition-all"
					onchange={(e) => {
						const cat = wbsCategories.find(
							(c) => c.id === Number(e.currentTarget.value),
						);
						if (cat) {
							// Map back to WBSNodeType for backward compatibility if needed by backend logic
							if (cat.name === "Phase") type = "Phase";
							else if (cat.name === "Milestone")
								type = "Milestone";
							else if (cat.name === "Work Package") type = "Task";
						}
					}}
				>
					{#each wbsCategories as cat}
						<option value={cat.id}>{cat.name}</option>
					{/each}
					{#if wbsCategories.length === 0}
						<option value="Phase">Phase</option>
						<option value="Milestone">Milestone</option>
						<option value="Task">Task</option>
					{/if}
				</select>
			</div>
			<div class="space-y-2">
				<label
					class="text-xs font-bold {errors.progress
						? 'text-rose-500'
						: 'text-slate-400'} uppercase tracking-widest"
					for="progress"
				>
					Progress (%) {errors.progress ? `— ${errors.progress}` : ""}
				</label>
				<input
					id="progress"
					type="number"
					bind:value={progress}
					min="0"
					max="100"
					class="w-full px-4 py-3 bg-slate-50 border {errors.progress
						? 'border-rose-300 ring-rose-100'
						: 'border-slate-200'} rounded-xl text-sm focus:ring-4 focus:ring-primary/10 focus:border-primary outline-none transition-all"
				/>
			</div>
		</div>

		<div class="space-y-2">
			<label
				class="text-xs font-bold text-slate-400 uppercase tracking-widest"
				for="pic"
			>
				PIC (Assignee)
			</label>
			<select
				id="pic"
				bind:value={assignedTo}
				class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-2 text-sm focus:ring-2 focus:ring-primary/20 outline-none transition-all"
			>
				<option value={null}>Unassigned</option>
				{#each users as user}
					<option value={user.id}>{user.full_name}</option>
				{/each}
			</select>
		</div>

		<!-- Estimated Effort -->
		<div class="space-y-2">
			<label
				class="text-xs font-bold text-slate-400 uppercase tracking-widest"
				for="effort"
			>
				Estimated Effort (Hours)
			</label>
			<div class="relative">
				<input
					id="effort"
					type="number"
					bind:value={estimatedEffort}
					min="0"
					step="0.5"
					placeholder="e.g., 8"
					class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl text-sm focus:ring-4 focus:ring-primary/10 focus:border-primary outline-none transition-all pr-12"
				/>
				<span
					class="absolute right-4 top-1/2 -translate-y-1/2 text-xs text-slate-400 font-medium"
					>hrs</span
				>
			</div>
		</div>

		<div class="space-y-2">
			<span
				class="text-xs font-bold {errors.dates
					? 'text-rose-500'
					: 'text-slate-400'} uppercase tracking-widest block"
			>
				Timeline Plan {errors.dates ? `— ${errors.dates}` : ""}
			</span>

			{#if parentNode && parentNode.planned_start_date && parentNode.planned_end_date}
				<div
					class="text-[10px] text-primary font-bold mb-1 flex items-center gap-1.5 bg-primary/5 p-2 rounded-lg border border-primary/20 shadow-sm"
				>
					<span class="material-symbols-outlined text-[14px]"
						>event_busy</span
					>
					<span class="uppercase tracking-tight"
						>Parent Constraint:</span
					>
					<span
						class="bg-white/80 px-1.5 py-0.5 rounded border border-primary/10"
						>{new Date(
							parentNode.planned_start_date,
						).toLocaleDateString("en-US")} - {new Date(
							parentNode.planned_end_date,
						).toLocaleDateString("en-US")}</span
					>
				</div>
			{/if}

			<div class="grid grid-cols-2 gap-6">
				<div class="space-y-1">
					<label
						class="text-[10px] font-medium text-slate-400"
						for="start">Start Date</label
					>
					<input
						id="start"
						type="date"
						bind:value={plannedStartDate}
						min={parentNode?.planned_start_date?.split("T")[0]}
						max={parentNode?.planned_end_date?.split("T")[0]}
						class="w-full bg-slate-50 border {errors.dates
							? 'border-rose-200 ring-4 ring-rose-50'
							: 'border-slate-200'} rounded-xl px-4 py-2.5 text-sm focus:ring-4 focus:ring-primary/10 transition-all outline-none"
					/>
				</div>
				<div class="space-y-1">
					<label
						class="text-[10px] font-medium text-slate-400"
						for="end">End Date</label
					>
					<input
						id="end"
						type="date"
						bind:value={plannedEndDate}
						min={parentNode?.planned_start_date?.split("T")[0]}
						max={parentNode?.planned_end_date?.split("T")[0]}
						class="w-full bg-slate-50 border {errors.dates
							? 'border-rose-200 ring-4 ring-rose-50'
							: 'border-slate-200'} rounded-xl px-4 py-2.5 text-sm focus:ring-4 focus:ring-primary/10 transition-all outline-none"
					/>
				</div>
			</div>
		</div>

		{#if parentPath}
			<div
				class="p-3 bg-primary/5 border border-primary/10 rounded-lg flex items-center gap-3"
			>
				<span class="material-symbols-outlined text-primary text-sm"
					>subdirectory_arrow_right</span
				>
				<span class="text-xs font-medium text-primary"
					>Add under WBS code: {parentPath}</span
				>
			</div>
		{/if}
	</div>

	{#snippet footer()}
		<Button
			variant="outline"
			onclick={() => (show = false)}
			disabled={loading}>Cancel</Button
		>
		<Button variant="primary" onclick={handleSubmit} disabled={loading}>
			{node?.id ? "Save Changes" : "Create"}
		</Button>
	{/snippet}
</Modal>
