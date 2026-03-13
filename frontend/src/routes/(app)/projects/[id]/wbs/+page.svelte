<script lang="ts">
	import { onMount, tick } from "svelte";
	import { projectService, type Project } from "$lib/services/projects";
	import {
		wbsService,
		type WBSNode,
		type WBSDependency,
		type WBSBaseline,
		type WBSBaselineNode,
	} from "$lib/services/wbs";
	import { userService, type User } from "$lib/services/users";
	import { holidayService, type Holiday } from "$lib/services/holidays";
	import { settingService } from "$lib/services/settings";
	import { toast } from "$lib/stores/toast";
	import {
		projectMembersService,
		type ProjectMember,
	} from "$lib/services/projectMembers";
	import {
		projectRolesService,
		type ProjectRole,
	} from "$lib/services/projectRoles";
	import WbsTaskDetails from "$lib/components/wbs/WbsTaskDetails.svelte";
	import WbsGanttChart from "$lib/components/wbs/WbsGanttChart.svelte";
	import WbsNodeModal from "$lib/components/wbs/WbsNodeModal.svelte";
	import ConfirmDialog from "$lib/components/ui/ConfirmDialog.svelte";
	import Pagination from "$lib/components/ui/Pagination.svelte";
	import { hasPermission, hasProjectPermission } from "$lib/utils/permission";
	import { authStore } from "$lib/services/auth";

	let { data } = $props();
	let project: Project | null = $state(null);
	// Nodes State
	let nodes: WBSNode[] = $state([]);
	let dependencies: WBSDependency[] = $state([]);
	let users: User[] = $state([]);
	let holidays: Holiday[] = $state([]);
	let baselines: WBSBaseline[] = $state([]);
	let pinnedBaselineIds = $state<number[]>([]);
	let baselineMapping = $state<Record<number, WBSBaselineNode[]>>({});
	let projectMembers: ProjectMember[] = $state([]);
	let projectRoles: ProjectRole[] = $state([]);
	let restDays: number[] = $state([0, 6]);

	let loading = $state(true);
	let viewMode: "Day" | "Week" | "Month" | "Quarter" = $state("Day");

	let sharedScrollTop = $state(0);
	let searchText = $state("");
	let hoveredNodeId = $state<number | null>(null);

	// Status Filter State
	let selectedStatus = $state(""); // '', 'todo', 'doing', 'done'
	let selectedPicId = $state<number | null>(null);

	// Variance Display State
	let showVariance = $state(true);

	// Pagination State
	let page = $state(1);
	let limit = $state(50);
	let totalItems = $state(0);

	let collapsedPaths = $state<Set<string>>(new Set());
	let loadedPaths = $state<Set<string>>(new Set());

	// Filtered Nodes: Handling local collapse and potentially future search results
	let filteredNodes = $derived.by(() => {
		const isFiltering =
			searchText !== "" ||
			selectedPicId !== null ||
			selectedStatus !== "";

		if (isFiltering) {
			return nodes;
		}

		return nodes.filter((node) => {
			const isHidden = Array.from(collapsedPaths).some((p) =>
				node.path.startsWith(p + "."),
			);
			return !isHidden;
		});
	});

	// Watch filters and fetch from server
	let searchTimeout: any;
	let currentController: AbortController | null = null;

	$effect(() => {
		const query = searchText;
		const pic = selectedPicId;
		const status = selectedStatus;
		const p = page;
		const l = limit;

		sharedScrollTop = 0;
		clearTimeout(searchTimeout);
		searchTimeout = setTimeout(() => {
			if (currentController) currentController.abort();
			currentController = new AbortController();

			loading = true;
			// For initial load or filter, we clear everything and start fresh
			fetchWBS(
				{ search: query, assignedTo: pic, status, page: p, limit: l },
				currentController.signal,
			).finally(() => {
				loading = false;
			});
		}, 300);

		return () => {
			clearTimeout(searchTimeout);
			if (currentController) currentController.abort();
		};
	});

	async function handleToggleCollapse(path: string) {
		if (collapsedPaths.has(path)) {
			// Expanding
			collapsedPaths.delete(path);
			collapsedPaths = new Set(collapsedPaths);

			// Check if children are loaded
			const node = nodes.find((n) => n.path === path);
			if (node && !loadedPaths.has(path)) {
				loading = true;
				try {
					const res = await wbsService.listTree(Number(data.id), {
						parentPath: path,
						limit: 500, // Fetch all children for this parent
						fields: [
							"id",
							"title",
							"path",
							"type",
							"planned_start_date",
							"planned_end_date",
							"progress",
							"assigned_to",
							"project_id",
							"has_children",
							"type_id",
						],
					});

					// Insert children after parent in the flat list
					const index = nodes.findIndex((n) => n.id === node.id);
					if (index !== -1) {
						nodes.splice(index + 1, 0, ...res.data);
						nodes = [...nodes];
						console.log("Added nodes, new length:", nodes.length);
						loadedPaths.add(path);
						loadedPaths = new Set(loadedPaths);
					}
				} finally {
					loading = false;
				}
			}
		} else {
			// Collapsing
			collapsedPaths.add(path);
			collapsedPaths = new Set(collapsedPaths);
		}
	}

	function handleToggleAll(expand: boolean) {
		if (expand) {
			collapsedPaths = new Set();
		} else {
			const newCollapsed = new Set<string>();
			nodes.forEach((node) => {
				const hasChildren = nodes.some((n) =>
					n.path.startsWith(node.path + "."),
				);
				if (hasChildren) {
					newCollapsed.add(node.path);
				}
			});
			collapsedPaths = newCollapsed;
		}
	}

	// Modal State
	let showModal = $state(false);
	let selectedNode = $state<WBSNode | null>(null);
	let parentPathForNewNode = $state<string | null>(null);
	let selectedParentNode = $state<WBSNode | null>(null);

	// Confirm State
	let showDeleteConfirm = $state(false);
	let nodeToDelete = $state<WBSNode | null>(null);

	// PIC Search/Dropdown state (restored)
	let picSearchText = $state("");
	let showPicDropdown = $state(false);
	let showBaselineDropdown = $state(false);

	let filteredPics = $derived(
		users.filter((user) =>
			user.full_name?.toLowerCase().includes(picSearchText.toLowerCase()),
		),
	);

	let selectedPicName = $derived(
		selectedPicId
			? users.find((u) => u.id === selectedPicId)?.full_name
			: "All Members",
	);

	async function fetchWBS(
		filter?: {
			search?: string;
			assignedTo?: number | null;
			status?: string;
			page?: number;
			limit?: number;
		},
		signal?: AbortSignal,
	) {
		try {
			// Optimization: Request only needed fields for the tree view
			const fields = [
				"id",
				"title",
				"path",
				"type",
				"planned_start_date",
				"planned_end_date",
				"progress",
				"assigned_to",
				"project_id",
				"has_children",
				"type_id",
			];

			const [nodesRes, depsRes] = await Promise.all([
				wbsService.listTree(Number(data.id), { ...filter, fields }),
				wbsService.listDependencies(Number(data.id)),
			]);

			if (signal?.aborted) return;

			nodes = nodesRes.data;
			totalItems = nodesRes.total;
			dependencies = depsRes;

			// On initial load (no filter), mark roots as loaded if they have no search
			if (!filter?.search && !filter?.assignedTo && !filter?.status) {
				loadedPaths = new Set();
				// Note: We don't know yet which nodes have children on backend
				// But we'll fetch them on expand.
			}
		} catch (e: any) {
			if (e.name === "AbortError") return;
			if (!e.isAuthError) {
				console.error("[WBS UI] Fetch error:", e);
				toast.error("Could not load WBS data");
			}
		}
	}

	onMount(async () => {
		try {
			// Don't call fetchWBS here, the $effect will trigger it on init
			const [
				projectRes,
				userRes,
				holidayRes,
				baselineRes,
				membersRes,
				rolesRes,
			] = await Promise.all([
				projectService.get(Number(data.id)),
				userService.getUsers(),
				holidayService.list(),
				wbsService.getBaselines(Number(data.id)),
				projectMembersService.getMembers(Number(data.id), 1, 1000),
				projectRolesService.getRoles(Number(data.id)),
			]);
			project = (projectRes as any).data || projectRes;
			users = (userRes as any).data || [];
			holidays = (holidayRes as any).data || [];
			baselines = (baselineRes as any) || [];
			projectMembers = (membersRes as any).data || [];
			projectRoles = rolesRes || [];
			// Load rest_days setting
			try {
				const settings = await settingService.getAll();
				if (settings?.rest_days) restDays = settings.rest_days;
			} catch {
				/* Use default if settings fail */
			}
		} catch (e: any) {
			if (!e.isAuthError) {
				console.error("[WBS UI] onMount error:", e);
				toast.error("Could not load project info");
			}
		} finally {
			loading = false;
		}
	});

	let isCreatingBaseline = $state(false);
	async function handleCreateBaseline() {
		if (isCreatingBaseline) return;
		isCreatingBaseline = true;
		try {
			const name = `Baseline ${new Date().toLocaleDateString("vi-VN")} ${new Date().toLocaleTimeString("vi-VN")}`;
			const desc = "Auto-generated baseline from snapshot.";
			const newB = await wbsService.createBaseline(
				Number(data.id),
				name,
				desc,
			);
			if (newB) {
				toast.success("Baseline saved successfully");
				baselines = await wbsService.getBaselines(Number(data.id));
			}
		} catch (e: any) {
			if (!e.isAuthError) {
				toast.error(e.message || "Error saving baseline");
			}
		} finally {
			isCreatingBaseline = false;
		}
	}

	const baselineColors = [
		"#94a3b8", // Slate
		"#6366f1", // Indigo
		"#10b981", // Emerald
		"#f59e0b", // Amber
		"#ef4444", // Red
		"#8b5cf6", // Violet
	];

	async function handleBaselineToggle(baselineId: number) {
		if (pinnedBaselineIds.includes(baselineId)) {
			pinnedBaselineIds = pinnedBaselineIds.filter(
				(id) => id !== baselineId,
			);
			const newMapping = { ...baselineMapping };
			delete newMapping[baselineId];
			baselineMapping = newMapping;
		} else {
			if (pinnedBaselineIds.length >= 5) {
				toast.warning("Maximum 5 baselines for comparison");
				return;
			}
			try {
				const bNodes = await wbsService.getBaselineNodes(
					Number(data.id),
					baselineId,
				);
				baselineMapping = {
					...baselineMapping,
					[baselineId]: bNodes,
				};
				pinnedBaselineIds = [...pinnedBaselineIds, baselineId];
			} catch (err) {
				toast.error("Could not load baseline details");
			}
		}
	}

	let baselineGanttData = $derived.by(() => {
		return pinnedBaselineIds.map((id, index) => {
			const b = baselines.find((bl) => bl.id === id);
			return {
				id,
				name: b?.name || "Baseline",
				color: baselineColors[index % baselineColors.length],
				nodes: baselineMapping[id] || [],
			};
		});
	});

	function handleEdit(node: WBSNode) {
		selectedNode = node;
		parentPathForNewNode = null;
		selectedParentNode = null;
		showModal = true;
	}

	function handleAddSubtask(parent: WBSNode) {
		selectedNode = null;
		parentPathForNewNode = parent.path;
		selectedParentNode = parent;
		showModal = true;
	}

	function handleNewTask() {
		selectedNode = null;
		parentPathForNewNode = null;
		selectedParentNode = null;
		showModal = true;
	}

	async function handleSaveNode(formData: any) {
		try {
			if (selectedNode?.id) {
				await wbsService.updateNode(
					Number(data.id),
					selectedNode.id,
					formData,
				);
				toast.success("Updated successfully");
			} else {
				await wbsService.createNode(Number(data.id), formData);
				toast.success("Added successfully");
			}
			await fetchWBS();
		} catch (e: any) {
			if (!e.isAuthError) {
				toast.error(e.message || "WBS operation error");
			}
		}
	}

	function handleDelete(node: WBSNode) {
		nodeToDelete = node;
		showDeleteConfirm = true;
	}

	async function confirmDelete() {
		if (!nodeToDelete) return;
		try {
			await wbsService.deleteNode(Number(data.id), nodeToDelete.id);
			toast.success("Deleted successfully");
			await fetchWBS();
		} catch (e: any) {
			if (!e.isAuthError) {
				toast.error(e.message || "Error deleting node");
			}
		} finally {
			showDeleteConfirm = false;
			nodeToDelete = null;
		}
	}

	async function handleStatusChange(node: WBSNode, progress: number) {
		try {
			await wbsService.updateNode(Number(data.id), node.id, {
				title: node.title,
				type: node.type,
				order_index: node.order_index,
				progress: progress,
				assigned_to: node.assigned_to ?? null,
				planned_start_date: node.planned_start_date ?? null,
				planned_end_date: node.planned_end_date ?? null,
			});
			toast.success("Status updated successfully");
			await fetchWBS();
		} catch (e: any) {
			if (!e.isAuthError) {
				toast.error(e.message || "Error updating status");
			}
		}
	}

	async function handlePicChange(node: WBSNode, userId: number | null) {
		try {
			await wbsService.updateNode(Number(data.id), node.id, {
				title: node.title,
				type: node.type,
				order_index: node.order_index,
				progress: node.progress,
				assigned_to: userId,
				planned_start_date: node.planned_start_date ?? null,
				planned_end_date: node.planned_end_date ?? null,
			});
			toast.success("Assignee updated successfully");
			await fetchWBS();
		} catch (e: any) {
			if (!e.isAuthError) {
				toast.error(e.message || "Error assigning person");
			}
		}
	}

	async function handleDatesChange(
		node: WBSNode,
		newStart: Date,
		newEnd: Date,
	) {
		// Validate against parent node dates
		if (node.path.includes(".")) {
			const parentPath = node.path.substring(
				0,
				node.path.lastIndexOf("."),
			);
			const parentNode = nodes.find((n) => n.path === parentPath);
			if (
				parentNode &&
				parentNode.planned_start_date &&
				parentNode.planned_end_date
			) {
				const pStart = new Date(parentNode.planned_start_date);
				const pEnd = new Date(parentNode.planned_end_date);

				// Reset time for date-only comparison
				const ps = new Date(
					pStart.getFullYear(),
					pStart.getMonth(),
					pStart.getDate(),
				);
				const pe = new Date(
					pEnd.getFullYear(),
					pEnd.getMonth(),
					pEnd.getDate(),
				);
				const ns = new Date(
					newStart.getFullYear(),
					newStart.getMonth(),
					newStart.getDate(),
				);
				const ne = new Date(
					newEnd.getFullYear(),
					newEnd.getMonth(),
					newEnd.getDate(),
				);

				// Allow equal dates (child can start/end on same day as parent)
				if (ns < ps || ne > pe) {
					toast.error(
						`Sub-task time cannot exceed parent task range: ${ps.toLocaleDateString("en-US")} → ${pe.toLocaleDateString("en-US")}`,
					);
					return;
				}
			}
		}

		try {
			await wbsService.updateNode(Number(data.id), node.id, {
				title: node.title,
				type: node.type,
				order_index: node.order_index,
				progress: node.progress,
				assigned_to: node.assigned_to,
				planned_start_date: newStart.toISOString(),
				planned_end_date: newEnd.toISOString(),
			});
			toast.success("Timeline updated successfully");
			const originalNode = nodes.find((n) => n.id === node.id);
			if (originalNode) {
				originalNode.planned_start_date = newStart.toISOString();
				originalNode.planned_end_date = newEnd.toISOString();
				nodes = [...nodes];
				console.log("Added nodes, new length:", nodes.length);
			}
		} catch (e: any) {
			if (!e.isAuthError) {
				toast.error(e.message || "Error updating timeline");
			}
		}
	}

	async function handleCreateDependency(fromId: number, toId: number) {
		try {
			await wbsService.createDependency(
				Number(data.id),
				fromId,
				toId,
				"FS",
			);
			toast.success("Dependency created successfully");
			dependencies = await wbsService.listDependencies(Number(data.id));
		} catch (e: any) {
			if (!e.isAuthError) {
				toast.error(e.message || "Error creating dependency");
			}
		}
	}

	async function handleDeleteDependency(dep: WBSDependency) {
		try {
			await wbsService.deleteDependency(Number(data.id), dep.id);
			toast.success("Dependency deleted successfully");
			dependencies = await wbsService.listDependencies(Number(data.id));
		} catch (e: any) {
			if (!e.isAuthError) {
				toast.error(e.message || "Error deleting dependency");
			}
		}
	}

	function handlePicSelect(user: User | null) {
		selectedPicId = user ? user.id : null;
		picSearchText = "";
		showPicDropdown = false;
	}

	function handleOutsideClick(e: MouseEvent) {
		const target = e.target as HTMLElement;
		if (!target.closest(".pic-dropdown-container")) {
			showPicDropdown = false;
		}
		if (!target.closest(".baseline-dropdown-container")) {
			showBaselineDropdown = false;
		}
	}

	$effect(() => {
		if (showPicDropdown || showBaselineDropdown) {
			window.addEventListener("click", handleOutsideClick);
			return () =>
				window.removeEventListener("click", handleOutsideClick);
		}
	});
</script>

<div
	class="h-full w-full flex flex-col bg-background-light overflow-hidden relative"
>
	<!-- Header Section -->
	<header
		class="h-20 glass-panel flex items-center justify-between px-8 z-40 shrink-0"
	>
		<div class="flex items-center gap-6">
			<a
				href="/projects/{data.id}"
				class="size-10 rounded-xl bg-slate-100 flex items-center justify-center text-slate-500 hover:text-primary transition-all hover:bg-white border border-slate-200 shadow-sm"
			>
				<span class="material-symbols-outlined">arrow_back</span>
			</a>
			<div class="flex flex-col">
				<h2 class="font-header text-2xl leading-none mb-1">
					Project WBS & Timeline
				</h2>
				<span
					class="text-xs font-semibold text-slate-400 uppercase tracking-wider"
					>{project?.project_name || "Loading..."}</span
				>
			</div>
		</div>

		<div class="flex items-center gap-4">
			<div class="relative baseline-dropdown-container">
				<button
					class="flex items-center gap-2 px-4 py-2 rounded-xl border border-slate-200 bg-white text-xs font-bold {pinnedBaselineIds.length >
					0
						? 'text-primary border-primary/30'
						: 'text-slate-600'} hover:bg-slate-50 transition-all shadow-sm min-w-48"
					onclick={(e) => {
						e.stopPropagation();
						showBaselineDropdown = !showBaselineDropdown;
					}}
				>
					<span class="material-symbols-outlined text-lg"
						>history</span
					>
					<span class="truncate flex-1 text-left">
						{pinnedBaselineIds.length > 0
							? `${pinnedBaselineIds.length} Baselines Pinned`
							: "Compare Baselines"}
					</span>
					<span
						class="material-symbols-outlined text-sm transition-transform {showBaselineDropdown
							? 'rotate-180'
							: ''}">expand_more</span
					>
				</button>

				{#if showBaselineDropdown}
					<div
						class="absolute top-full right-0 mt-2 w-64 bg-white border border-slate-200 rounded-xl shadow-xl z-[100] p-2 animate-in fade-in slide-in-from-top-2 duration-200"
					>
						<div class="max-h-60 overflow-y-auto custom-scrollbar">
							{#each baselines as b, i}
								{@const isPinned = pinnedBaselineIds.includes(
									b.id,
								)}
								{@const color = isPinned
									? baselineColors[
											pinnedBaselineIds.indexOf(b.id) %
												baselineColors.length
										]
									: null}
								<button
									class="w-full text-left px-3 py-2 rounded-lg text-xs hover:bg-slate-50 transition-colors flex items-center gap-3 {isPinned
										? 'bg-primary/5 text-primary font-bold'
										: 'text-slate-600 font-medium'}"
									onclick={() => handleBaselineToggle(b.id)}
								>
									<div
										class="size-3 rounded-full border border-slate-300"
										style="background-color: {color ||
											'transparent'}"
									></div>
									<span class="truncate flex-1">{b.name}</span
									>
									<span
										class="material-symbols-outlined text-sm"
									>
										{isPinned
											? "check_box"
											: "check_box_outline_blank"}
									</span>
								</button>
							{/each}
							{#if baselines.length === 0}
								<div
									class="p-4 text-center text-xs text-slate-400"
								>
									No baselines found
								</div>
							{/if}
						</div>
					</div>
				{/if}
			</div>

			{#if hasProjectPermission($authStore.user, $authStore.token, projectMembers, projectRoles, "project:update")}
				<button
					class="bg-white border border-slate-200 text-slate-600 px-4 py-2 rounded-xl font-bold text-sm flex items-center gap-2 hover:bg-slate-50 transition-all shadow-sm disabled:opacity-50"
					onclick={handleCreateBaseline}
					disabled={isCreatingBaseline || nodes.length === 0}
				>
					{#if isCreatingBaseline}
						<span
							class="material-symbols-outlined text-lg animate-spin"
							>refresh</span
						>
					{:else}
						<span class="material-symbols-outlined text-lg"
							>save</span
						>
					{/if}
					Save Baseline
				</button>
			{/if}

			<div class="w-px h-6 bg-slate-200 mx-2"></div>

			{#if hasProjectPermission($authStore.user, $authStore.token, projectMembers, projectRoles, "project:wbs:create")}
				<button
					class="bg-primary text-white px-6 py-2.5 rounded-xl font-bold text-sm flex items-center gap-2 hover:opacity-90 transition-all shadow-lg shadow-primary/20"
					onclick={handleNewTask}
				>
					<span class="material-symbols-outlined text-lg">add</span>
					New Task
				</button>
			{/if}
		</div>
	</header>

	<!-- Filter & View Bar -->
	<div
		class="h-14 bg-white/40 border-b border-slate-200/50 px-8 flex items-center gap-6 z-30 backdrop-blur-sm shrink-0"
	>
		<div class="flex items-center gap-2">
			<span
				class="text-[10px] font-bold text-slate-400 uppercase tracking-widest"
				>Filter & View</span
			>
		</div>
		<div class="flex items-center gap-4">
			<div
				class="relative flex items-center bg-white border border-slate-200 rounded-lg shadow-sm hover:border-primary/50 transition-all"
			>
				<span
					class="material-symbols-outlined absolute left-2 text-slate-400 text-[16px] pointer-events-none"
				>
					calendar_month
				</span>
				<select
					class="pl-8 pr-8 py-1.5 text-xs font-semibold text-slate-700 bg-transparent border-none appearance-none outline-none cursor-pointer w-32"
					bind:value={viewMode}
				>
					<option value="Day">View: Day</option>
					<option value="Week">View: Week</option>
					<option value="Month">View: Month</option>
					<option value="Quarter">View: Quarter</option>
				</select>
				<span
					class="material-symbols-outlined absolute right-2 text-slate-400 text-sm pointer-events-none"
				>
					expand_more
				</span>
			</div>
			<div class="h-6 w-px bg-slate-200/60 mx-1"></div>

			<div
				class="relative flex items-center bg-white border border-slate-200 rounded-lg shadow-sm hover:border-primary/50 transition-all"
			>
				<span
					class="material-symbols-outlined absolute left-2 text-slate-400 text-[16px] pointer-events-none"
					>search</span
				>
				<input
					type="text"
					class="pl-8 pr-4 py-1.5 text-xs font-semibold text-slate-700 bg-transparent border-none outline-none w-48 placeholder:text-slate-400"
					placeholder="Search..."
					bind:value={searchText}
				/>
			</div>

			<div class="h-6 w-px bg-slate-200/60 mx-1"></div>

			<div class="relative pic-dropdown-container">
				<button
					class="flex items-center gap-2 px-3 py-1.5 rounded-lg border border-slate-200 bg-white text-xs font-medium {selectedPicId
						? 'text-primary border-primary/30'
						: 'text-slate-600'} hover:border-primary/50 transition-all shadow-sm w-44"
					onclick={(e) => {
						e.stopPropagation();
						showPicDropdown = !showPicDropdown;
					}}
				>
					<span class="material-symbols-outlined text-sm">person</span
					>
					<span class="truncate flex-1 text-left"
						>{selectedPicName}</span
					>
					<span
						class="material-symbols-outlined text-sm transition-transform {showPicDropdown
							? 'rotate-180'
							: ''}">expand_more</span
					>
				</button>

				{#if showPicDropdown}
					<div
						class="absolute top-full left-0 mt-2 w-64 bg-white border border-slate-200 rounded-xl shadow-xl z-[100] p-2 animate-in fade-in slide-in-from-top-2 duration-200"
					>
						<div class="relative mb-2">
							<span
								class="material-symbols-outlined absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-400 text-sm"
								>search</span
							>
							<input
								type="text"
								class="w-full pl-8 pr-3 py-1.5 bg-slate-50 border border-slate-200 rounded-lg text-xs outline-none focus:ring-2 focus:ring-primary/10 transition-all"
								placeholder="Search assignee..."
								bind:value={picSearchText}
								onclick={(e) => e.stopPropagation()}
							/>
						</div>
						<div class="max-h-60 overflow-y-auto custom-scrollbar">
							<button
								class="w-full text-left px-3 py-2 rounded-lg text-xs hover:bg-slate-50 transition-colors flex items-center justify-between {selectedPicId ===
								null
									? 'bg-primary/5 text-primary font-bold'
									: 'text-slate-600 font-medium'}"
								onclick={() => handlePicSelect(null)}
							>
								All members
								{#if selectedPicId === null}
									<span
										class="material-symbols-outlined text-sm"
										>check</span
									>
								{/if}
							</button>
							<div class="h-px bg-slate-100 my-1 mx-1"></div>
							{#each filteredPics as user}
								<button
									class="w-full text-left px-3 py-2 rounded-lg text-xs hover:bg-slate-50 transition-colors flex items-center gap-2 {selectedPicId ===
									user.id
										? 'bg-primary/5 text-primary font-bold'
										: 'text-slate-600 font-medium'}"
									onclick={() => handlePicSelect(user)}
								>
									<div
										class="size-6 rounded-full bg-slate-200 flex items-center justify-center text-[10px] text-slate-500 font-bold overflow-hidden"
									>
										{user.full_name?.charAt(0)}
									</div>
									<span class="truncate flex-1"
										>{user.full_name}</span
									>
									{#if selectedPicId === user.id}
										<span
											class="material-symbols-outlined text-sm"
											>check</span
										>
									{/if}
								</button>
							{/each}
							{#if filteredPics.length === 0}
								<div
									class="p-4 text-center text-[10px] text-slate-400"
								>
									No results found
								</div>
							{/if}
						</div>
					</div>
				{/if}
			</div>

			<div class="h-6 w-px bg-slate-200/60 mx-1"></div>

			<div
				class="relative flex items-center bg-white border border-slate-200 rounded-lg shadow-sm hover:border-primary/50 transition-all"
			>
				<span
					class="material-symbols-outlined absolute left-2 text-slate-400 text-[16px] pointer-events-none"
				>
					assignment_turned_in
				</span>
				<select
					class="pl-8 pr-8 py-1.5 text-xs font-semibold {selectedStatus !==
					''
						? 'text-primary'
						: 'text-slate-700'} bg-transparent border-none appearance-none outline-none cursor-pointer w-36"
					bind:value={selectedStatus}
				>
					<option value="">Status: All</option>
					<option value="todo">Scheduled (0%)</option>
					<option value="doing">In Progress (1-99%)</option>
					<option value="done">Completed (100%)</option>
				</select>
				<span
					class="material-symbols-outlined absolute right-2 text-slate-400 text-sm pointer-events-none"
				>
					expand_more
				</span>
			</div>
			<div class="h-6 w-px bg-slate-200/60 mx-1"></div>

			<button
				class="flex items-center gap-2 px-3 py-1.5 rounded-lg border border-slate-200 {showVariance
					? 'bg-primary/5 text-primary border-primary/20'
					: 'bg-white text-slate-400'} text-xs font-bold hover:border-primary/50 transition-all shadow-sm"
				onclick={() => (showVariance = !showVariance)}
			>
				<span class="material-symbols-outlined text-sm"
					>{showVariance ? "visibility" : "visibility_off"}</span
				>
				Variance
			</button>

			<div class="h-6 w-px bg-slate-200/60 mx-1"></div>

			<button
				class="text-xs font-medium text-primary hover:underline transition-all disabled:opacity-30 disabled:no-underline"
				onclick={() => {
					selectedPicId = null;
					searchText = "";
					selectedStatus = "";
					picSearchText = "";
				}}
				disabled={!selectedPicId &&
					searchText === "" &&
					selectedStatus === ""}
			>
				Clear Filters
			</button>
			<div class="h-4 w-px bg-slate-200"></div>
			<button
				class="p-1.5 rounded-lg hover:bg-slate-100 text-slate-400 transition-colors"
			>
				<span class="material-symbols-outlined text-lg">fullscreen</span
				>
			</button>
		</div>
	</div>

	<!-- Main Content Area (Split Pane) -->
	<div class="flex-1 flex overflow-hidden">
		{#if loading}
			<div class="flex-1 flex items-center justify-center bg-white/50">
				<div class="flex flex-col items-center gap-4">
					<div
						class="size-12 border-4 border-primary/20 border-t-primary rounded-full animate-spin"
					></div>
					<p class="text-sm font-medium text-slate-500">
						Loading WBS structure...
					</p>
				</div>
			</div>
		{:else}
			<div class="contents">
				<WbsTaskDetails
					nodes={filteredNodes}
					allNodes={nodes}
					{users}
					baselineNodes={showVariance
						? baselineGanttData[0]?.nodes || []
						: []}
					{hoveredNodeId}
					{collapsedPaths}
					onToggleCollapse={handleToggleCollapse}
					onToggleAll={handleToggleAll}
					bind:scrollTop={sharedScrollTop}
					onEdit={handleEdit}
					onAddSubtask={handleAddSubtask}
					onDelete={handleDelete}
					onStatusChange={handleStatusChange}
					onPicChange={handlePicChange}
					onHover={(id) => (hoveredNodeId = id)}
					{searchText}
				/>
			</div>

			<div class="contents">
				<WbsGanttChart
					nodes={filteredNodes}
					{dependencies}
					{holidays}
					baselineSets={showVariance ? baselineGanttData : []}
					{restDays}
					bind:scrollTop={sharedScrollTop}
					{viewMode}
					{hoveredNodeId}
					onHover={(id) => (hoveredNodeId = id)}
					onProgressChange={handleStatusChange}
					onDatesChange={handleDatesChange}
					onCreateDependency={handleCreateDependency}
					onDeleteDependency={handleDeleteDependency}
				/>
			</div>
		{/if}
	</div>

	<WbsNodeModal
		bind:show={showModal}
		projectId={Number(data.id)}
		{users}
		node={selectedNode}
		parentPath={parentPathForNewNode}
		parentNode={selectedParentNode}
		onSave={handleSaveNode}
	/>

	<!-- Confirm Delete -->
	<ConfirmDialog
		bind:show={showDeleteConfirm}
		title="Confirm Delete"
		message="Are you sure you want to delete {nodeToDelete?.title}? All sub-tasks will be deleted as well."
		confirmText="Delete"
		variant="danger"
		onConfirm={confirmDelete}
		onCancel={() => (showDeleteConfirm = false)}
	/>

	<!-- Footer Pagination -->
	<div class="shrink-0 z-30 w-full bg-slate-50/50">
		<Pagination
			total={totalItems}
			{page}
			{limit}
			onPageChange={(p) => (page = p)}
			onLimitChange={(l) => {
				limit = l;
				page = 1;
			}}
		/>
	</div>
</div>

<!-- Desktop Only Overlay -->
<div
	class="lg:hidden fixed inset-0 bg-background-dark/95 z-[100] flex flex-col items-center justify-center p-8 text-center"
>
	<span class="material-symbols-outlined text-6xl text-primary mb-4"
		>desktop_windows</span
	>
	<h2 class="text-2xl font-bold text-white mb-2">Desktop Only View</h2>
	<p class="text-slate-400 text-sm">
		WBS Management is optimized for 1440px+ displays.
	</p>
</div>

<style>
	.contents {
		display: contents;
	}
	:global(.contents > *) {
		height: 100%;
	}
</style>
