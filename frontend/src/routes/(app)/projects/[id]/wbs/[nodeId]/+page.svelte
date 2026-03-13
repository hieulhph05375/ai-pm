<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/state";
    import {
        wbsService,
        type WBSNode,
        type WBSComment,
    } from "$lib/services/wbs";
    import { projectService, type Project } from "$lib/services/projects";
    import { userService, type User } from "$lib/services/users";
    import { authStore, authService } from "$lib/services/auth";
    import { toast } from "$lib/stores/toast";
    import Badge from "$lib/components/ui/Badge.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import ContentHeader from "$lib/components/ui/ContentHeader.svelte";
    import Avatar from "$lib/components/ui/Avatar.svelte";
    import Modal from "$lib/components/ui/Modal.svelte";
    import ConfirmDialog from "$lib/components/ui/ConfirmDialog.svelte";

    const projectId = Number(page.params.id);
    const nodeId = Number(page.params.nodeId);

    let project: Project | null = $state(null);
    let node: WBSNode | null = $state(null);
    let users: User[] = $state([]);
    let comments: WBSComment[] = $state([]);
    let loading = $state(true);
    let commentText = $state("");
    let isSavingDescription = $state(false);
    let isSubmittingComment = $state(false);
    let activeTab = $state("details"); // "details" or "comments"

    // Edit State
    let isEditing = $state(false);
    let editingNode = $state<Partial<WBSNode>>({});

    // Comment Edit/Delete State
    let editingCommentId = $state<number | null>(null);
    let editingCommentContent = $state("");
    let isSavingComment = $state(false);
    let isDeleteModalOpen = $state(false);
    let commentToDeleteId = $state<number | null>(null);
    let isEditingDescription = $state(false);
    let descriptionText = $state("");

    onMount(async () => {
        try {
            const [projectRes, nodeRes, userRes] = await Promise.all([
                projectService.get(projectId),
                wbsService.getNode(projectId, nodeId),
                userService.getUsers(),
            ]);
            project = (projectRes as any).data || projectRes;
            node = nodeRes;
            users = (userRes as any).data || [];

            if (node) {
                descriptionText = node.description || "";
                loadComments();
            }
        } catch (e: any) {
            if (!e.isAuthError) {
                console.error("[Task Detail] Load error:", e);
                toast.error("Could not load task information");
            }
        } finally {
            loading = false;
        }
    });

    function getStatusColor(
        progress: number,
    ):
        | "emerald"
        | "amber"
        | "indigo"
        | "primary"
        | "rose"
        | "pink"
        | "slate"
        | undefined {
        if (progress === 100) return "emerald";
        if (progress > 0) return "amber";
        return "indigo";
    }

    function getStatusLabel(progress: number) {
        if (progress === 100) return "Completed";
        if (progress > 0) return "In Progress";
        return "Not Started";
    }

    let assignedUser = $derived(users.find((u) => u.id === node?.assigned_to));

    async function saveDescription() {
        if (!node) return;
        isSavingDescription = true;
        try {
            await wbsService.updateNode(projectId, nodeId, {
                title: node.title,
                type: node.type,
                order_index: node.order_index,
                progress: node.progress,
                assigned_to: node.assigned_to,
                description: descriptionText,
                planned_start_date: node.planned_start_date,
                planned_end_date: node.planned_end_date,
            });
            node.description = descriptionText;
            isEditingDescription = false;
            toast.success("Description saved");
        } catch (e: any) {
            toast.error("Could not save description");
        } finally {
            isSavingDescription = false;
        }
    }

    function insertText(prefix: string, suffix: string = "") {
        const textarea = document.getElementById(
            "description-textarea",
        ) as HTMLTextAreaElement;
        if (!textarea) return;

        const start = textarea.selectionStart;
        const end = textarea.selectionEnd;
        const text = descriptionText;
        const before = text.substring(0, start);
        const selection = text.substring(start, end);
        const after = text.substring(end);

        descriptionText = before + prefix + selection + suffix + after;

        // Reset focus and selection
        setTimeout(() => {
            textarea.focus();
            textarea.setSelectionRange(
                start + prefix.length,
                end + prefix.length,
            );
        }, 0);
    }

    async function submitComment() {
        if (!commentText.trim() || isSubmittingComment) return;
        isSubmittingComment = true;
        try {
            const newComment = await wbsService.addComment(
                projectId,
                nodeId,
                commentText,
            );
            comments = [newComment, ...comments];
            commentText = "";
            toast.success("Discussion posted");
        } catch (e: any) {
            toast.error("Could not post discussion");
        } finally {
            isSubmittingComment = false;
        }
    }

    async function loadComments() {
        try {
            comments = await wbsService.getComments(projectId, nodeId);
        } catch (e) {
            console.error("Error loading comments:", e);
        }
    }

    function startInlineEdit() {
        if (!node) return;
        editingNode = {
            title: node.title,
            type: node.type,
            assigned_to: node.assigned_to,
            planned_start_date: node.planned_start_date
                ? node.planned_start_date.split("T")[0]
                : "",
            planned_end_date: node.planned_end_date
                ? node.planned_end_date.split("T")[0]
                : "",
            progress: node.progress,
        };
        descriptionText = node.description || "";
        isEditing = true;
    }

    function cancelInlineEdit() {
        isEditing = false;
        editingNode = {};
    }

    async function handleUpdate() {
        if (!node) return;
        try {
            await wbsService.updateNode(projectId, nodeId, {
                ...editingNode,
                description: descriptionText,
            });
            node = { ...node, ...editingNode, description: descriptionText };
            isEditing = false;
            toast.success("Updated successfully");
        } catch (e) {
            toast.error("Could not update information");
        }
    }

    async function quickUpdateStatus(newProgress: number) {
        if (!node || loading) return;
        const oldProgress = node.progress;

        // Optimistic update
        node.progress = newProgress;

        try {
            await wbsService.updateNode(projectId, nodeId, {
                title: node.title,
                type: node.type,
                order_index: node.order_index,
                progress: newProgress,
                assigned_to: node.assigned_to,
                description: node.description,
                planned_start_date: node.planned_start_date,
                planned_end_date: node.planned_end_date,
            });
            toast.success(`Updated: ${getStatusLabel(newProgress)}`);
        } catch (e) {
            // Revert on error
            node.progress = oldProgress;
            toast.error("Could not update status");
        }
    }

    async function quickUpdateAssignedTo(newAssignedTo: number | null) {
        if (!node || loading) return;
        const oldAssignedTo = node.assigned_to;

        // Optimistic update
        node.assigned_to = newAssignedTo;

        try {
            await wbsService.updateNode(projectId, nodeId, {
                title: node.title,
                type: node.type,
                order_index: node.order_index,
                progress: node.progress,
                assigned_to: newAssignedTo,
                description: node.description,
                planned_start_date: node.planned_start_date,
                planned_end_date: node.planned_end_date,
            });
            toast.success("Assignee updated");
        } catch (e) {
            // Revert on error
            node.assigned_to = oldAssignedTo;
            toast.error("Could not update assignee");
        }
    }

    function startEdit(comment: WBSComment) {
        editingCommentId = comment.id;
        editingCommentContent = comment.content;
    }

    function cancelEdit() {
        editingCommentId = null;
        editingCommentContent = "";
    }

    async function saveEdit(commentId: number) {
        if (!editingCommentContent.trim() || isSavingComment) return;
        isSavingComment = true;
        try {
            await wbsService.updateComment(
                projectId,
                nodeId,
                commentId,
                editingCommentContent,
            );
            comments = comments.map((c) =>
                c.id === commentId
                    ? { ...c, content: editingCommentContent }
                    : c,
            );
            cancelEdit();
            toast.success("Discussion updated");
        } catch (e) {
            toast.error("Could not update discussion");
        } finally {
            isSavingComment = false;
        }
    }

    function handleDeleteComment(commentId: number) {
        commentToDeleteId = commentId;
        isDeleteModalOpen = true;
    }

    async function confirmDelete() {
        if (commentToDeleteId === null) return;
        try {
            await wbsService.deleteComment(
                projectId,
                nodeId,
                commentToDeleteId,
            );
            comments = comments.filter((c) => c.id !== commentToDeleteId);
            toast.success("Discussion deleted");
            isDeleteModalOpen = false;
            commentToDeleteId = null;
        } catch (e) {
            toast.error("Could not delete discussion");
        }
    }
</script>

<div class="min-h-screen flex flex-col bg-slate-50/30">
    <!-- Breadcrumbs/Header Section -->
    <header
        class="bg-white/80 backdrop-blur-md border-b border-slate-200/60 sticky top-0 z-30 px-6 py-4"
    >
        <div class="max-w-7xl mx-auto flex items-center justify-between">
            <nav
                class="flex items-center gap-2 text-[10px] font-bold text-slate-400 uppercase tracking-[0.1em]"
            >
                <a href="/projects" class="hover:text-primary transition-colors"
                    >Projects</a
                >
                <span class="material-symbols-outlined text-[16px]"
                    >chevron_right</span
                >
                {#if project}
                    <a
                        href="/projects/{projectId}"
                        class="hover:text-primary transition-colors max-w-[150px] truncate"
                        >{project.project_name}</a
                    >
                {/if}
                <span class="material-symbols-outlined text-[16px]"
                    >chevron_right</span
                >
                <a
                    href="/projects/{projectId}/wbs"
                    class="hover:text-primary transition-colors"
                    >WBS & Timeline</a
                >
                {#if node}
                    <span class="material-symbols-outlined text-[16px]"
                        >chevron_right</span
                    >
                    <span
                        class="text-slate-900 truncate max-w-[250px] font-black"
                        >{node.title}</span
                    >
                {/if}
            </nav>

            <div class="flex items-center gap-4">
                <button
                    class="size-9 flex items-center justify-center rounded-xl bg-white border border-slate-200 text-slate-400 hover:text-primary transition-all hover:border-primary/20"
                >
                    <span class="material-symbols-outlined text-[18px]"
                        >notifications</span
                    >
                </button>
                <button
                    class="size-9 flex items-center justify-center rounded-xl bg-white border border-slate-200 text-slate-400 hover:text-primary transition-all hover:border-primary/20"
                >
                    <span class="material-symbols-outlined text-[18px]"
                        >chat_bubble</span
                    >
                </button>

                <div
                    class="flex items-center gap-3 pl-4 ml-1 border-l border-slate-200/60"
                >
                    {#if $authStore.user}
                        <div class="text-right hidden sm:block">
                            <p
                                class="text-[11px] font-bold text-slate-900 leading-none"
                            >
                                {$authStore.user.fullName}
                            </p>
                            <p
                                class="text-[9px] font-bold text-slate-400 mt-1 uppercase tracking-tight"
                            >
                                Admin
                            </p>
                        </div>
                        <Avatar name={$authStore.user.fullName} size="sm" />
                        <button
                            class="p-1.5 text-slate-400 hover:text-rose-500 hover:bg-rose-50 rounded-lg transition-all ml-1"
                            title="Logout"
                            onclick={() => authService.logout()}
                        >
                            <span class="material-symbols-outlined text-[18px]"
                                >logout</span
                            >
                        </button>
                    {:else}
                        <div class="text-right">
                            <p
                                class="text-xs font-bold text-slate-900 leading-none"
                            >
                                Guest
                            </p>
                        </div>
                        <Avatar name="Guest" size="sm" />
                    {/if}
                </div>
            </div>
        </div>
    </header>

    <div class="flex-1 overflow-y-auto p-10 bg-mesh">
        {#if loading}
            <div class="animate-pulse space-y-6">
                <div class="h-32 bg-slate-100 rounded-3xl"></div>
                <div class="h-10 w-64 bg-slate-100 rounded-xl"></div>
                <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
                    <div class="h-64 bg-slate-100 rounded-3xl"></div>
                    <div class="h-64 bg-slate-100 rounded-3xl col-span-2"></div>
                </div>
            </div>
        {:else if node}
            <!-- Tabs Navigation -->
            <div
                class="flex items-center gap-1 bg-white p-1 rounded-2xl border border-slate-200 w-fit shadow-sm mb-8"
            >
                <button
                    onclick={() => (activeTab = "details")}
                    class="px-6 py-2.5 rounded-xl text-sm font-bold transition-all flex items-center gap-2 {activeTab ===
                    'details'
                        ? 'bg-primary text-white shadow-lg shadow-primary/20'
                        : 'text-slate-500 hover:bg-slate-50'}"
                >
                    <span class="material-symbols-outlined text-[18px]"
                        >info</span
                    >
                    Task Details
                </button>
                <button
                    onclick={() => (activeTab = "comments")}
                    class="px-6 py-2.5 rounded-xl text-sm font-bold transition-all flex items-center gap-2 {activeTab ===
                    'comments'
                        ? 'bg-primary text-white shadow-lg shadow-primary/20'
                        : 'text-slate-500 hover:bg-slate-50'}"
                >
                    <span class="material-symbols-outlined text-[18px]"
                        >forum</span
                    >
                    Discussion
                    {#if comments.length > 0}
                        <span
                            class="size-5 rounded-full bg-white/20 flex items-center justify-center text-[10px]"
                            >{comments.length}</span
                        >
                    {/if}
                </button>
            </div>

            <ContentHeader
                title={node.title || ""}
                subtitle="{node.path} • {node.type}"
            >
                {#snippet titleSnippet()}
                    {#if isEditing}
                        <div class="flex-1 max-w-2xl">
                            <input
                                type="text"
                                bind:value={editingNode.title}
                                class="w-full bg-white border border-slate-200 focus:ring-2 focus:ring-primary/20 rounded-xl px-4 py-2 text-xl font-bold outline-none transition-all"
                                placeholder="Enter task title..."
                            />
                        </div>
                    {:else}
                        <h2
                            class="font-display text-2xl font-bold text-slate-900 tracking-tight truncate"
                        >
                            {node?.title || ""}
                        </h2>
                    {/if}
                {/snippet}
                <div class="flex items-center gap-3">
                    {#if isEditing}
                        <Button variant="outline" onclick={cancelInlineEdit}
                            >Cancel</Button
                        >
                        <Button icon="save" onclick={handleUpdate}>Save</Button>
                    {:else}
                        <Button icon="edit" onclick={startInlineEdit}
                            >Edit</Button
                        >
                    {/if}
                </div>
            </ContentHeader>

            <!-- New Timeline Info Bar below Header -->
            <div
                class="flex items-center gap-6 mb-8 px-5 py-3 bg-slate-50/50 border border-slate-200/60 rounded-2xl w-fit shadow-sm transition-all hover:bg-slate-50"
            >
                <div class="flex items-center gap-2.5">
                    <div
                        class="size-7 rounded-lg bg-white border border-slate-200 flex items-center justify-center shadow-sm"
                    >
                        <span
                            class="material-symbols-outlined text-[16px] text-slate-500"
                            >calendar_today</span
                        >
                    </div>
                    <span
                        class="text-[10px] font-bold text-slate-400 uppercase tracking-widest"
                        >Planned Period</span
                    >
                </div>
                <div class="h-4 w-px bg-slate-200/80"></div>
                <div class="flex items-center gap-6 text-sm">
                    <div class="flex items-center gap-2.5">
                        <span
                            class="text-[10px] text-slate-400 font-bold uppercase tracking-tight"
                            >Start</span
                        >
                        {#if isEditing}
                            <input
                                type="date"
                                bind:value={editingNode.planned_start_date}
                                class="bg-white border border-slate-200 rounded-lg px-2 py-1 text-slate-900 outline-none focus:ring-2 focus:ring-primary/20 text-xs font-bold"
                            />
                        {:else}
                            <span
                                class="font-bold text-slate-700 bg-white px-2 py-0.5 rounded-md border border-slate-100 shadow-sm"
                            >
                                {node.planned_start_date
                                    ? new Date(
                                          node.planned_start_date,
                                      ).toLocaleDateString("vi-VN")
                                    : "---"}
                            </span>
                        {/if}
                    </div>
                    <span class="text-slate-300 font-light">→</span>
                    <div class="flex items-center gap-2.5">
                        <span
                            class="text-[10px] text-slate-400 font-bold uppercase tracking-tight"
                            >End</span
                        >
                        {#if isEditing}
                            <input
                                type="date"
                                bind:value={editingNode.planned_end_date}
                                class="bg-white border border-slate-200 rounded-lg px-2 py-1 text-slate-900 outline-none focus:ring-2 focus:ring-primary/20 text-xs font-bold"
                            />
                        {:else}
                            <span
                                class="font-bold text-slate-700 bg-white px-2 py-0.5 rounded-md border border-slate-100 shadow-sm"
                            >
                                {node.planned_end_date
                                    ? new Date(
                                          node.planned_end_date,
                                      ).toLocaleDateString("vi-VN")
                                    : "---"}
                            </span>
                        {/if}
                    </div>
                </div>
            </div>

            {#if activeTab === "details"}
                <div
                    class="grid grid-cols-1 lg:grid-cols-3 gap-6 animate-in fade-in slide-in-from-bottom-4 duration-300"
                >
                    <!-- Main Info -->
                    <div class="lg:col-span-2 space-y-6">
                        <div
                            class="bg-white rounded-3xl border border-slate-200 shadow-sm relative overflow-hidden flex flex-col"
                        >
                            <div
                                class="p-6 border-b border-slate-100 flex items-center justify-between bg-slate-50/10"
                            >
                                <h3
                                    class="text-xs font-bold text-slate-400 uppercase tracking-widest flex items-center gap-2"
                                >
                                    <span
                                        class="material-symbols-outlined text-[18px]"
                                        >description</span
                                    >
                                    Detailed Description
                                </h3>
                            </div>

                            <div class="p-6">
                                {#if isEditing}
                                    <div class="space-y-4">
                                        <!-- Minimal Toolbar -->
                                        <div
                                            class="flex items-center gap-1 p-1 bg-slate-50 rounded-xl border border-slate-200/50 w-fit"
                                        >
                                            <button
                                                onclick={() =>
                                                    insertText("**", "**")}
                                                class="size-8 flex items-center justify-center rounded-lg hover:bg-white hover:shadow-sm text-slate-600 transition-all"
                                                title="Bold"
                                            >
                                                <span
                                                    class="material-symbols-outlined text-[18px]"
                                                    >format_bold</span
                                                >
                                            </button>
                                            <button
                                                onclick={() =>
                                                    insertText("_", "_")}
                                                class="size-8 flex items-center justify-center rounded-lg hover:bg-white hover:shadow-sm text-slate-600 transition-all"
                                                title="Italic"
                                            >
                                                <span
                                                    class="material-symbols-outlined text-[18px]"
                                                    >format_italic</span
                                                >
                                            </button>
                                            <div
                                                class="w-px h-4 bg-slate-200 mx-1"
                                            ></div>
                                            <button
                                                onclick={() => insertText("- ")}
                                                class="size-8 flex items-center justify-center rounded-lg hover:bg-white hover:shadow-sm text-slate-600 transition-all"
                                                title="List"
                                            >
                                                <span
                                                    class="material-symbols-outlined text-[18px]"
                                                    >format_list_bulleted</span
                                                >
                                            </button>
                                        </div>

                                        <div class="relative group">
                                            <textarea
                                                id="description-textarea"
                                                class="w-full bg-white border border-slate-200 focus:border-primary/30 rounded-2xl p-5 text-sm text-slate-600 outline-none transition-all placeholder:text-slate-400 min-h-[300px] resize-none leading-relaxed shadow-inner"
                                                placeholder="Add detailed description for this task..."
                                                maxlength="2000"
                                                bind:value={descriptionText}
                                            ></textarea>
                                        </div>

                                        <div
                                            class="flex items-center justify-between"
                                        >
                                            <span
                                                class="text-[10px] font-bold text-slate-400"
                                            >
                                                {descriptionText.length} / 2000 characters
                                            </span>
                                        </div>
                                    </div>
                                {:else}
                                    <div
                                        class="text-sm text-slate-600 leading-relaxed whitespace-pre-wrap min-h-[200px]"
                                    >
                                        {#if node?.description}
                                            {node.description}
                                        {:else}
                                            <p class="italic text-slate-400">
                                                No detailed description yet.
                                            </p>
                                        {/if}
                                    </div>
                                {/if}
                            </div>
                        </div>
                    </div>

                    <!-- Sidebar Info -->
                    <div class="space-y-6">
                        <!-- Status & PIC -->
                        <div
                            class="bg-white rounded-3xl p-6 border border-slate-200 shadow-sm space-y-6"
                        >
                            <div>
                                <h4
                                    class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mb-3"
                                >
                                    Status & Progress
                                </h4>
                                {#if isEditing}
                                    <div class="space-y-3">
                                        <div
                                            class="flex items-center justify-between"
                                        >
                                            <span
                                                class="text-xs font-bold text-slate-500 uppercase"
                                                >Progress (%)</span
                                            >
                                            <input
                                                type="number"
                                                min="0"
                                                max="100"
                                                bind:value={
                                                    editingNode.progress
                                                }
                                                class="w-20 px-3 py-1.5 bg-slate-50 border border-slate-200 rounded-lg text-sm font-bold outline-none focus:ring-2 focus:ring-primary/20 transition-all text-center"
                                            />
                                        </div>
                                        <input
                                            type="range"
                                            min="0"
                                            max="100"
                                            bind:value={editingNode.progress}
                                            class="w-full h-1.5 bg-slate-100 rounded-full appearance-none cursor-pointer accent-primary"
                                        />
                                        <div
                                            class="flex justify-between text-[10px] font-bold text-slate-400"
                                        >
                                            <span>0%</span>
                                            <span>50%</span>
                                            <span>100%</span>
                                        </div>
                                    </div>
                                {:else}
                                    <div class="relative group mb-6">
                                        <select
                                            class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-2xl focus:ring-2 focus:ring-primary/20 outline-none transition-all text-sm font-bold appearance-none cursor-pointer text-slate-700"
                                            value={node.progress}
                                            onchange={(e) =>
                                                quickUpdateStatus(
                                                    Number(
                                                        e.currentTarget.value,
                                                    ),
                                                )}
                                        >
                                            <option value={0}
                                                >Not Started</option
                                            >
                                            <option value={50}
                                                >In Progress</option
                                            >
                                            <option value={100}
                                                >Completed</option
                                            >
                                        </select>
                                        <span
                                            class="material-symbols-outlined absolute right-4 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none"
                                            >unfold_more</span
                                        >
                                    </div>

                                    <div
                                        class="flex items-center justify-between mb-3"
                                    >
                                        <div class="flex items-center gap-2">
                                            <Badge
                                                variant={getStatusColor(
                                                    node.progress,
                                                )}
                                            >
                                                {getStatusLabel(node.progress)}
                                            </Badge>
                                        </div>
                                        <span
                                            class="text-sm font-black text-slate-900"
                                            >{node.progress}%</span
                                        >
                                    </div>
                                    <div
                                        class="h-2 bg-slate-100 rounded-full overflow-hidden"
                                    >
                                        <div
                                            class="h-full bg-primary transition-all duration-500"
                                            style="width: {node.progress}%"
                                        ></div>
                                    </div>
                                {/if}
                            </div>

                            <div class="pt-6 border-t border-slate-100">
                                <h4
                                    class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mb-3"
                                >
                                    Assignee
                                </h4>
                                {#if isEditing}
                                    <div class="relative">
                                        <select
                                            bind:value={editingNode.assigned_to}
                                            class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-2xl focus:ring-2 focus:ring-primary/20 outline-none transition-all text-sm font-medium appearance-none"
                                        >
                                            <option value={null}
                                                >Unassigned</option
                                            >
                                            {#each users as user}
                                                <option value={user.id}
                                                    >{user.full_name}</option
                                                >
                                            {/each}
                                        </select>
                                        <span
                                            class="material-symbols-outlined absolute right-4 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none"
                                            >expand_more</span
                                        >
                                    </div>
                                {:else}
                                    <div class="relative group">
                                        <select
                                            value={node.assigned_to}
                                            onchange={(e) =>
                                                quickUpdateAssignedTo(
                                                    e.currentTarget.value
                                                        ? Number(
                                                              e.currentTarget
                                                                  .value,
                                                          )
                                                        : null,
                                                )}
                                            class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-2xl focus:ring-2 focus:ring-primary/20 outline-none transition-all text-sm font-bold appearance-none cursor-pointer text-slate-700"
                                        >
                                            <option value={null}
                                                >Unassigned</option
                                            >
                                            {#each users as user}
                                                <option value={user.id}
                                                    >{user.full_name}</option
                                                >
                                            {/each}
                                        </select>
                                        <span
                                            class="material-symbols-outlined absolute right-4 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none"
                                            >unfold_more</span
                                        >
                                    </div>

                                    {#if assignedUser}
                                        <div
                                            class="mt-3 flex items-center gap-3 bg-primary/5 p-3 rounded-2xl border border-primary/10"
                                        >
                                            <Avatar
                                                size="sm"
                                                name={assignedUser.full_name}
                                            />
                                            <div>
                                                <p
                                                    class="text-sm font-bold text-slate-900 leading-tight"
                                                >
                                                    {assignedUser.full_name}
                                                </p>
                                                <p
                                                    class="text-[10px] text-primary font-bold uppercase tracking-wider"
                                                >
                                                    Member
                                                </p>
                                            </div>
                                        </div>
                                    {/if}
                                {/if}
                            </div>
                        </div>
                    </div>
                </div>
            {:else}
                <!-- Comments Tab -->
                <div
                    class="bg-white rounded-3xl border border-slate-200 shadow-sm overflow-hidden flex flex-col animate-in fade-in slide-in-from-bottom-4 duration-300 min-h-[600px]"
                >
                    <div
                        class="p-6 border-b border-slate-100 flex items-center justify-between bg-slate-50/30"
                    >
                        <h3
                            class="text-sm font-bold text-slate-900 flex items-center gap-2"
                        >
                            <span class="material-symbols-outlined text-primary"
                                >forum</span
                            >
                            Team Discussion ({comments.length})
                        </h3>
                    </div>

                    <div
                        class="flex-1 p-6 space-y-6 overflow-y-auto custom-scrollbar"
                    >
                        {#each comments as comment}
                            <div
                                class="flex gap-4 group animate-in slide-in-from-left-2 duration-300"
                            >
                                <Avatar
                                    size="sm"
                                    name={comment.user_name || "User"}
                                />
                                <div class="flex-1 space-y-1">
                                    <div class="flex items-center gap-2">
                                        <span
                                            class="text-sm font-bold text-slate-900"
                                            >{comment.user_name ||
                                                "Member"}</span
                                        >
                                        <span
                                            class="text-[10px] font-bold text-slate-400 uppercase tracking-tighter"
                                        >
                                            {new Date(
                                                comment.created_at,
                                            ).toLocaleString("vi-VN")}
                                        </span>
                                        {#if $authStore.user?.id === comment.user_id && editingCommentId !== comment.id}
                                            <div
                                                class="flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity"
                                            >
                                                <button
                                                    onclick={() =>
                                                        startEdit(comment)}
                                                    class="flex items-center gap-1 text-[10px] font-bold text-slate-400 hover:text-primary transition-colors"
                                                >
                                                    <span
                                                        style="font-size: 14px;"
                                                        class="material-symbols-outlined text-[14px]"
                                                        >edit</span
                                                    >
                                                    Edit
                                                </button>
                                                <button
                                                    onclick={() =>
                                                        handleDeleteComment(
                                                            comment.id,
                                                        )}
                                                    class="flex items-center gap-1 text-[10px] font-bold text-slate-400 hover:text-rose-500 transition-colors"
                                                >
                                                    <span
                                                        style="font-size: 14px;"
                                                        class="material-symbols-outlined text-[14px]"
                                                        >delete</span
                                                    >
                                                    Delete
                                                </button>
                                            </div>
                                        {/if}
                                    </div>
                                    {#if editingCommentId === comment.id}
                                        <div class="space-y-2 mt-2">
                                            <textarea
                                                class="w-full bg-white border border-slate-200 focus:border-primary/30 rounded-xl p-3 text-sm text-slate-600 outline-none transition-all resize-none shadow-sm"
                                                bind:value={
                                                    editingCommentContent
                                                }
                                                rows="3"
                                                onkeydown={(e) => {
                                                    if (e.key === "Enter") {
                                                        e.preventDefault();
                                                        saveEdit(comment.id);
                                                    }
                                                }}
                                            ></textarea>
                                            <div class="flex justify-end gap-2">
                                                <button
                                                    onclick={cancelEdit}
                                                    class="px-3 py-1.5 rounded-lg text-[10px] font-bold text-slate-500 hover:bg-slate-100 transition-all"
                                                >
                                                    Cancel
                                                </button>
                                                <button
                                                    onclick={() =>
                                                        saveEdit(comment.id)}
                                                    disabled={!editingCommentContent.trim() ||
                                                        isSavingComment}
                                                    class="px-3 py-1.5 rounded-lg text-[10px] font-bold bg-primary text-white hover:bg-primary/90 transition-all shadow-sm disabled:opacity-50"
                                                >
                                                    {isSavingComment
                                                        ? "Saving..."
                                                        : "Save changes"}
                                                </button>
                                            </div>
                                        </div>
                                    {:else}
                                        <div
                                            class="bg-slate-50 rounded-2xl p-4 text-sm text-slate-600 border border-slate-100/50 relative shadow-sm"
                                        >
                                            {comment.content}
                                        </div>
                                    {/if}
                                </div>
                            </div>
                        {:else}
                            <div
                                class="h-full flex flex-col items-center justify-center gap-4 text-slate-400 opacity-50"
                            >
                                <div
                                    class="size-20 rounded-full bg-slate-100 flex items-center justify-center"
                                >
                                    <span
                                        class="material-symbols-outlined text-4xl"
                                        >chat_bubble</span
                                    >
                                </div>
                                <p class="text-sm font-medium italic">
                                    No discussions for this task yet.
                                </p>
                            </div>
                        {/each}
                    </div>

                    <div class="p-6 bg-slate-50/50 border-t border-slate-100">
                        <div class="flex gap-4">
                            <Avatar size="sm" name="Me" />
                            <div class="relative flex-1">
                                <textarea
                                    class="w-full bg-white border border-slate-200 focus:border-primary/30 rounded-2xl p-4 pr-14 text-sm text-slate-600 outline-none transition-all placeholder:text-slate-400 min-h-[100px] resize-none shadow-sm"
                                    placeholder="Write a reply or ask a question..."
                                    bind:value={commentText}
                                    onkeydown={(e) => {
                                        if (e.key === "Enter") {
                                            e.preventDefault();
                                            submitComment();
                                        }
                                    }}
                                ></textarea>
                                <button
                                    class="absolute right-3 bottom-3 size-10 rounded-xl bg-primary text-white flex items-center justify-center hover:scale-110 active:scale-95 transition-all shadow-lg shadow-primary/20 disabled:opacity-50"
                                    onclick={submitComment}
                                    disabled={!commentText.trim() ||
                                        isSubmittingComment}
                                >
                                    <span class="material-symbols-outlined"
                                        >send</span
                                    >
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            {/if}
        {/if}
    </div>
</div>

<ConfirmDialog
    bind:show={isDeleteModalOpen}
    title="Confirm Delete Discussion"
    message="Are you sure you want to delete this discussion? This action cannot be undone."
    confirmText="Confirm Delete"
    variant="danger"
    onConfirm={confirmDelete}
    onCancel={() => {
        isDeleteModalOpen = false;
        commentToDeleteId = null;
    }}
/>

<style>
    .custom-scrollbar::-webkit-scrollbar {
        width: 6px;
    }
    .custom-scrollbar::-webkit-scrollbar-track {
        background: transparent;
    }
    .custom-scrollbar::-webkit-scrollbar-thumb {
        background: #e2e8f0;
        border-radius: 10px;
    }
    .custom-scrollbar::-webkit-scrollbar-thumb:hover {
        background: #cbd5e1;
    }
</style>
