<script lang="ts">
	import { onMount } from "svelte";
	import { projectService, type Project } from "$lib/services/projects";
	import { categoryService, type Category } from "$lib/services/categories";
	import DataTable from "$lib/components/ui/DataTable.svelte";
	import Badge from "$lib/components/ui/Badge.svelte";
	import Button from "$lib/components/ui/Button.svelte";
	import Modal from "$lib/components/ui/Modal.svelte";
	import EmptyState from "$lib/components/ui/EmptyState.svelte";
	import ConfirmDialog from "$lib/components/ui/ConfirmDialog.svelte";
	import { toast } from "$lib/stores/toast";
	import { hasPermission } from "$lib/utils/permission";
	import { authStore } from "$lib/services/auth";

	let projects = $state<Project[]>([]);
	let total = $state(0);
	let page = $state(1);
	let limit = $state(10);
	let search = $state("");
	let status = $state("");
	let loading = $state(false);

	let projectStatuses = $state<Category[]>([]);
	let projectPhases = $state<Category[]>([]);
	let portfolioCategories = $state<Category[]>([]);
	let healthCategories = $state<Category[]>([]);
	let priorityLevels = $state<Category[]>([]);

	let showModal = $state(false);
	let modalMode = $state<"create" | "edit">("create");
	let currentProject = $state<Partial<Project>>({});

	let showDeleteConfirm = $state(false);
	let projectToDelete: number | null = null;
	let showExportConfirm = $state(false);

	const columns = [
		{ key: "project_info", label: "Project", class: "w-80" },
		{ key: "manager", label: "Manager", class: "w-48" },
		{ key: "status", label: "Status", class: "w-40" },
		{ key: "health", label: "Health", class: "w-32" },
		{ key: "progress", label: "Progress", class: "w-48" },
		{
			key: "budget",
			label: "Budget (USD)",
			align: "right" as const,
			class: "w-40",
		},
		{ key: "actions", label: "", align: "right" as const, class: "w-24" },
	];

	async function loadCategories() {
		try {
			const [s, p, pc, h, pr] = await Promise.all([
				categoryService.listCategories(1, 100, "", 1),
				categoryService.listCategories(1, 100, "", 2),
				categoryService.listCategories(1, 100, "", 3),
				categoryService.listCategories(1, 100, "", 4),
				categoryService.listCategories(1, 100, "", 5),
			]);
			projectStatuses = s.data;
			projectPhases = p.data;
			portfolioCategories = pc.data;
			healthCategories = h.data;
			priorityLevels = pr.data;
		} catch (err) {
			console.error("Failed to load categories", err);
		}
	}

	async function loadProjects() {
		loading = true;
		if (projectStatuses.length === 0) await loadCategories();
		try {
			const res = await projectService.list(page, limit, search, status);
			projects = res.data;
			total = res.total;
		} catch (err: any) {
			if (!err.isAuthError) {
				toast.error("Could not load project list");
			}
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		modalMode = "create";
		// Pick first active categories as defaults
		const defaultStatus = projectStatuses.find((c) => c.is_active);
		const defaultPhase = projectPhases.find((c) => c.is_active);
		const defaultHealth = healthCategories.find((c) => c.is_active);
		const defaultPriority = priorityLevels.find((c) => c.is_active);

		currentProject = {
			project_id: "",
			project_name: "",
			project_manager: "",
			planned_start_date: new Date().toISOString().split("T")[0],
			planned_end_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000)
				.toISOString()
				.split("T")[0],
			project_status_id: defaultStatus?.id,
			current_phase_id: defaultPhase?.id,
			overall_health_id: defaultHealth?.id,
			priority_level_id: defaultPriority?.id,
			strategic_score: 0,
			approved_budget: 0,
			progress: 0,
			spi: 1.0,
			cpi: 1.0,
		};
		showModal = true;
	}

	function openEditModal(project: Project) {
		modalMode = "edit";
		// Clone and format dates for input type="date"
		const formattedProject = { ...project };
		if (formattedProject.planned_start_date) {
			formattedProject.planned_start_date = new Date(
				formattedProject.planned_start_date as string,
			)
				.toISOString()
				.split("T")[0];
		}
		if (formattedProject.planned_end_date) {
			formattedProject.planned_end_date = new Date(
				formattedProject.planned_end_date as string,
			)
				.toISOString()
				.split("T")[0];
		}
		currentProject = formattedProject;
		showModal = true;
	}

	async function saveProject() {
		// Basic validation
		if (!currentProject.project_name || !currentProject.project_id) {
			toast.error("Please enter project name and ID");
			return;
		}

		try {
			// Convert dates to ISO strings safely for Go
			const projectData = { ...currentProject };

			const formatToISO = (dateVal: any) => {
				if (!dateVal || dateVal === "") return undefined;
				try {
					const d = new Date(dateVal);
					return isNaN(d.getTime()) ? undefined : d.toISOString();
				} catch (e) {
					return undefined;
				}
			};

			if (modalMode === "create") {
				projectData.planned_start_date = formatToISO(
					projectData.planned_start_date,
				);
				projectData.planned_end_date = formatToISO(
					projectData.planned_end_date,
				);

				// Sync actual dates
				projectData.actual_start_date = projectData.planned_start_date;
				projectData.actual_end_date = projectData.planned_end_date;

				await projectService.create(projectData as Project);
				toast.success("New project created");
			} else {
				projectData.planned_start_date = formatToISO(
					projectData.planned_start_date,
				);
				projectData.planned_end_date = formatToISO(
					projectData.planned_end_date,
				);

				await projectService.update(
					projectData.id!,
					projectData as Project,
				);
				toast.success("Project updated");
			}
			showModal = false;
			loadProjects();
		} catch (err: any) {
			if (!err.isAuthError) {
				console.error("Save error:", err);
				toast.error(
					"Error saving project: " +
						(err instanceof Error ? err.message : "Invalid data"),
				);
			}
		}
	}

	function confirmDelete(id: number) {
		projectToDelete = id;
		showDeleteConfirm = true;
	}

	async function handleDelete() {
		if (!projectToDelete) return;
		try {
			await projectService.delete(projectToDelete);
			toast.success("Project deleted");
			showDeleteConfirm = false;
			projectToDelete = null;
			loadProjects();
		} catch (err: any) {
			if (!err.isAuthError) {
				toast.error("Could not delete project");
			}
		}
	}

	onMount(loadProjects);

	function getHealthColor(item: Project) {
		if (item.overall_health_cat?.color) {
			return item.overall_health_cat.color;
		}
		// Fallback for legacy data
		switch (item.overall_health?.toLowerCase()) {
			case "green":
				return "#10b981"; // emerald-500
			case "yellow":
				return "#f59e0b"; // amber-500
			case "red":
				return "#f43f5e"; // rose-500
			default:
				return "#64748b"; // slate-500
		}
	}
	function triggerExport() {
		showExportConfirm = true;
	}

	async function handleExport() {
		try {
			await projectService.exportList(search, status);
			toast.success("Project report exported successfully");
		} catch (err: any) {
			if (!err.isAuthError) {
				toast.error("Could not export report");
			}
		} finally {
			showExportConfirm = false;
		}
	}
</script>

<div class="space-y-6">
	<!-- Header Section -->
	<div
		class="flex flex-col md:flex-row md:items-center justify-between gap-4"
	>
		<div>
			<h1 class="text-2xl font-bold text-slate-900 tracking-tight">
				Project Management
			</h1>
			<p class="text-slate-500 font-medium text-sm mt-1">
				List of all projects, initiatives, and WBS of the organization.
			</p>
		</div>
		<div class="flex items-center gap-3">
			{#if hasPermission($authStore.user, $authStore.token, "project:export")}
				<Button
					variant="outline"
					onclick={triggerExport}
					class="hidden md:flex items-center gap-2"
				>
					<span class="material-symbols-outlined text-[18px]"
						>ios_share</span
					>
					Export Report
				</Button>
			{/if}
			{#if hasPermission($authStore.user, $authStore.token, "project:create")}
				<Button
					onclick={openCreateModal}
					class="flex items-center gap-2 bg-primary hover:bg-primary/90 text-white shadow-lg shadow-primary/20 transition-all active:scale-95"
				>
					<span class="material-symbols-outlined text-[18px]"
						>add_circle</span
					>
					Create New Project
				</Button>
			{/if}
		</div>
	</div>

	<!-- Controls Section -->
	<div
		class="grid grid-cols-1 md:grid-cols-12 gap-4 items-center bg-white p-2 rounded-2xl border border-slate-200 shadow-sm"
	>
		<div class="md:col-span-4 relative group">
			<span
				class="material-symbols-outlined absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-primary transition-colors"
				>search</span
			>
			<input
				type="text"
				placeholder="Search project name, ID..."
				bind:value={search}
				onkeydown={(e) => e.key === "Enter" && loadProjects()}
				class="w-full bg-slate-50 border-none rounded-xl pl-12 pr-4 py-3 focus:ring-2 focus:ring-primary/20 transition-all text-sm font-medium placeholder:text-slate-400"
			/>
		</div>
		<div class="md:col-span-3">
			<select
				bind:value={status}
				onchange={loadProjects}
				class="w-full bg-slate-50 border-none rounded-xl px-4 py-3 focus:ring-2 focus:ring-primary/20 transition-all text-sm font-medium text-slate-600"
			>
				<option value="">All statuses</option>
				{#each projectStatuses as s}
					<option value={s.name}>{s.name}</option>
				{/each}
			</select>
		</div>
	</div>

	<!-- Table Section -->
	<div
		class="bg-white rounded-3xl border border-slate-200 shadow-xl shadow-slate-200/50 overflow-hidden"
	>
		<DataTable
			{columns}
			items={projects}
			{loading}
			{total}
			{page}
			{limit}
			onPageChange={(p) => {
				page = p;
				loadProjects();
			}}
			onLimitChange={(l) => {
				limit = l;
				page = 1;
				loadProjects();
			}}
		>
			{#snippet emptyState()}
				<EmptyState
					icon="account_tree"
					actionIcon="add_circle"
					title="No projects"
					message="There are no projects in the system. Start by creating the first project."
					actionLabel="Create New Project"
					onaction={openCreateModal}
				/>
			{/snippet}
			{#snippet rowCell({ item, column })}
				{#if column.key === "project_info"}
					<div class="flex items-center gap-4">
						<div
							class="size-11 rounded-2xl bg-gradient-to-br from-primary/10 to-primary/5 flex items-center justify-center border border-primary/10 shadow-sm group-hover:scale-110 transition-transform"
						>
							<span
								class="material-symbols-outlined text-primary text-[18px]"
								>account_tree</span
							>
						</div>
						<div class="flex flex-col min-w-0">
							<a
								href="/projects/{item.id}"
								class="text-sm font-bold text-slate-900 truncate hover:text-primary cursor-pointer transition-colors leading-tight"
								>{item.project_name}</a
							>
							<span
								class="text-[11px] text-slate-400 font-bold uppercase tracking-wider mt-1 flex items-center gap-1.5"
							>
								<span
									class="size-1.5 rounded-full bg-primary/40"
								></span>
								{item.project_id}
							</span>
						</div>
					</div>
				{:else if column.key === "manager"}
					<div class="flex items-center gap-2">
						<div
							class="size-7 rounded-full bg-slate-100 flex items-center justify-center border border-slate-200"
						>
							<span
								class="material-symbols-outlined text-slate-400 text-sm"
								>person</span
							>
						</div>
						<span class="text-sm font-semibold text-slate-700"
							>{item.project_manager || "Unassigned"}</span
						>
					</div>
				{:else if column.key === "status"}
					<Badge
						color={item.current_phase_cat?.color || "indigo"}
						class="!font-bold !text-[10px] !uppercase !tracking-widest"
					>
						{item.current_phase_cat?.name || item.current_phase}
					</Badge>
				{:else if column.key === "health"}
					<div class="flex items-center gap-2">
						<div
							class="size-2 rounded-full shadow-sm"
							style="background-color: {getHealthColor(item)}"
						></div>
						<span class="text-xs font-bold text-slate-700 uppercase"
							>{item.overall_health_cat?.name ||
								item.overall_health}</span
						>
					</div>
				{:else if column.key === "progress"}
					<div class="w-full max-w-[140px] space-y-2">
						<div
							class="flex items-center justify-between text-[11px] font-bold"
						>
							<span
								class="text-slate-500 uppercase tracking-tighter"
								>Completed</span
							>
							<span class="text-primary">{item.progress}%</span>
						</div>
						<div
							class="h-1.5 w-full bg-slate-100 rounded-full overflow-hidden border border-slate-50 shadow-inner"
						>
							<div
								class="h-full bg-gradient-to-r from-primary to-primary-light rounded-full transition-all duration-1000"
								style="width: {item.progress}%"
							></div>
						</div>
					</div>
				{:else if column.key === "budget"}
					<span class="text-sm font-bold text-slate-900"
						>{new Intl.NumberFormat("en-US").format(
							item.approved_budget,
						)}</span
					>
				{:else if column.key === "actions"}
					<div class="flex items-center justify-end gap-1">
						{#if hasPermission($authStore.user, $authStore.token, "project:update")}
							<button
								onclick={() => openEditModal(item)}
								class="size-8 rounded-lg hover:bg-slate-100 flex items-center justify-center text-slate-400 hover:text-primary transition-all active:scale-90"
							>
								<span
									class="material-symbols-outlined text-[18px]"
									>edit_note</span
								>
							</button>
						{/if}
						{#if hasPermission($authStore.user, $authStore.token, "project:delete")}
							<button
								onclick={() => confirmDelete(item.id!)}
								class="size-8 rounded-lg hover:bg-rose-50 flex items-center justify-center text-slate-400 hover:text-rose-500 transition-all active:scale-90"
							>
								<span
									class="material-symbols-outlined text-[18px]"
									>delete</span
								>
							</button>
						{/if}
					</div>
				{/if}
			{/snippet}
		</DataTable>
	</div>

	<!-- Modal Form -->
	<Modal
		bind:show={showModal}
		title={modalMode === "create" ? "Create New Project" : "Edit Project"}
	>
		<div class="grid grid-cols-2 gap-4 p-4">
			<div class="col-span-2 space-y-1.5">
				<label
					for="project-name"
					class="text-[11px] font-bold uppercase text-slate-400"
					>Project Name</label
				>
				<input
					id="project-name"
					type="text"
					bind:value={currentProject.project_name}
					class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-2.5 outline-none focus:ring-2 focus:ring-primary/20 transition-all font-medium"
					placeholder="Enter project name..."
				/>
			</div>
			<div class="space-y-1.5">
				<label
					for="project-id"
					class="text-[11px] font-bold uppercase text-slate-400"
					>Project ID</label
				>
				<input
					id="project-id"
					type="text"
					bind:value={currentProject.project_id}
					class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-2.5 outline-none focus:ring-2 focus:ring-primary/20 transition-all font-medium"
					placeholder="Example: PRJ-001"
				/>
			</div>
			<div class="space-y-1.5">
				<label
					for="project-manager"
					class="text-[11px] font-bold uppercase text-slate-400"
					>Project Manager</label
				>
				<input
					id="project-manager"
					type="text"
					bind:value={currentProject.project_manager}
					class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-2.5 outline-none focus:ring-2 focus:ring-primary/20 transition-all font-medium"
					placeholder="Manager name..."
				/>
			</div>
			<div class="space-y-1.5">
				<label
					for="planned-start-date"
					class="text-[11px] font-bold uppercase text-slate-400"
					>Planned Start Date</label
				>
				<input
					id="planned-start-date"
					type="date"
					bind:value={currentProject.planned_start_date}
					class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-2.5 outline-none focus:ring-2 focus:ring-primary/20 transition-all font-medium"
				/>
			</div>
			<div class="space-y-1.5">
				<label
					for="planned-end-date"
					class="text-[11px] font-bold uppercase text-slate-400"
					>Planned End Date</label
				>
				<input
					id="planned-end-date"
					type="date"
					bind:value={currentProject.planned_end_date}
					class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-2.5 outline-none focus:ring-2 focus:ring-primary/20 transition-all font-medium"
				/>
			</div>

			<div class="space-y-1.5">
				<label
					for="project-status"
					class="text-[11px] font-bold uppercase text-slate-400"
					>Status</label
				>
				<select
					id="project-status"
					bind:value={currentProject.project_status_id}
					class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-2.5 outline-none focus:ring-2 focus:ring-primary/20 transition-all font-medium"
				>
					{#each projectStatuses as s}
						<option value={s.id}>{s.name}</option>
					{/each}
				</select>
			</div>
			<div class="space-y-1.5">
				<label
					for="current-phase"
					class="text-[11px] font-bold uppercase text-slate-400"
					>Current Phase</label
				>
				<select
					id="current-phase"
					bind:value={currentProject.current_phase_id}
					class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-2.5 outline-none focus:ring-2 focus:ring-primary/20 transition-all font-medium"
				>
					{#each projectPhases as p}
						<option value={p.id}>{p.name}</option>
					{/each}
				</select>
			</div>
			<div class="space-y-1.5">
				<label
					for="overall-health"
					class="text-[11px] font-bold uppercase text-slate-400"
					>Overall Health</label
				>
				<select
					id="overall-health"
					bind:value={currentProject.overall_health_id}
					class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-2.5 outline-none focus:ring-2 focus:ring-primary/20 transition-all font-medium"
				>
					{#each healthCategories as h}
						<option value={h.id}>{h.name}</option>
					{/each}
				</select>
			</div>
			<div class="space-y-1.5">
				<label
					for="priority-level"
					class="text-[11px] font-bold uppercase text-slate-400"
					>Priority Level</label
				>
				<select
					id="priority-level"
					bind:value={currentProject.priority_level_id}
					class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-2.5 outline-none focus:ring-2 focus:ring-primary/20 transition-all font-medium"
				>
					{#each priorityLevels as p}
						<option value={p.id}>{p.name}</option>
					{/each}
				</select>
			</div>
			<div class="space-y-1.5">
				<label
					for="portfolio-category"
					class="text-[11px] font-bold uppercase text-slate-400"
					>Portfolio Category</label
				>
				<select
					id="portfolio-category"
					bind:value={currentProject.portfolio_category_id}
					class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-2.5 outline-none focus:ring-2 focus:ring-primary/20 transition-all font-medium"
				>
					<option value={undefined}>None</option>
					{#each portfolioCategories as c}
						<option value={c.id}>{c.name}</option>
					{/each}
				</select>
			</div>

			<div class="col-span-2 flex justify-end gap-3 mt-4">
				<Button variant="outline" onclick={() => (showModal = false)}
					>Cancel</Button
				>
				<Button onclick={saveProject}>Save Project</Button>
			</div>
		</div>
	</Modal>
</div>

<ConfirmDialog
	show={showDeleteConfirm}
	title="Delete Project"
	message="Are you sure you want to delete this project? This action cannot be undone."
	confirmText="Confirm Delete"
	variant="danger"
	onConfirm={handleDelete}
	onCancel={() => {
		showDeleteConfirm = false;
		projectToDelete = null;
	}}
/>

<ConfirmDialog
	show={showExportConfirm}
	title="Export Project Report"
	message="Are you sure you want to export the project report? The file will be downloaded to your computer."
	confirmText="Export Report"
	variant="primary"
	onConfirm={handleExport}
	onCancel={() => (showExportConfirm = false)}
/>
