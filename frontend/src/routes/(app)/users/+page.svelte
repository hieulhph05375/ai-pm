<script lang="ts">
	import { onMount } from "svelte";
	import {
		userService,
		validatePassword,
		type User,
	} from "$lib/services/users";
	import { rbacService, type Role } from "$lib/services/rbac";
	import { toast } from "$lib/stores/toast";
	import { authStore } from "$lib/services/auth";
	import ContentHeader from "$lib/components/ui/ContentHeader.svelte";
	import DataTable from "$lib/components/ui/DataTable.svelte";
	import Button from "$lib/components/ui/Button.svelte";
	import Badge from "$lib/components/ui/Badge.svelte";
	import Avatar from "$lib/components/ui/Avatar.svelte";
	import Modal from "$lib/components/ui/Modal.svelte";
	import ConfirmDialog from "$lib/components/ui/ConfirmDialog.svelte";
	import { hasPermission } from "$lib/utils/permission";

	let users = $state<User[]>([]);
	let roles = $state<Role[]>([]);
	let total = $state(0);
	let page = $state(1);
	let limit = $state(10);
	let loading = $state(true);
	let searchQuery = $state("");
	let error = $state("");

	// Action State
	let showModal = $state(false);
	let modalMode = $state<"create" | "edit" | "reset">("create");
	let currentUser = $state<Partial<User>>({});
	let tempPassword = $state("");
	let passwordError = $state("");
	let isSubmitting = $state(false);

	// Confirm Dialog State
	let showConfirm = $state(false);
	let confirmData = $state<{
		title: string;
		message: string;
		action: () => void;
	}>({ title: "", message: "", action: () => {} });

	const columns = [
		{ key: "user", label: "Member", class: "w-64" },
		{ key: "email", label: "Email" },
		{ key: "role", label: "Role", class: "w-40" },
		{ key: "is_admin", label: "Admin", class: "w-32" },
		{ key: "status", label: "Status", class: "w-44" },
		{ key: "actions", label: "", align: "right" as const, class: "w-24" },
	];

	let showPassword = $state(false);

	onMount(async () => {
		await Promise.all([loadRoles(), loadUsers()]);
	});

	async function loadRoles() {
		try {
			roles = await rbacService.listRoles();
		} catch (e: any) {
			console.error("Failed to load roles", e);
		}
	}

	async function loadUsers() {
		loading = true;
		try {
			const res = await userService.getUsers(searchQuery, page, limit);
			users = res.data;
			total = res.total;
		} catch (e: any) {
			if (!e.isAuthError) {
				toast.error(e.message || "Error loading user list");
			}
		} finally {
			loading = false;
		}
	}

	function onPageChange(newPage: number) {
		page = newPage;
		loadUsers();
	}

	function getRoleInfo(roleId: number) {
		const role = roles.find((r) => r.id === roleId);
		if (!role) return { label: "Unknown Role", variant: "slate" };

		let variant = "slate";
		const nameUpper = role.name.toUpperCase();
		if (nameUpper.includes("ADMIN")) variant = "rose";
		else if (nameUpper.includes("MANAGER") || nameUpper === "PMO")
			variant = "purple";
		else if (nameUpper.includes("LEAD")) variant = "emerald";
		else if (nameUpper.includes("MEMBER")) variant = "primary";
		else if (nameUpper.includes("VIEWER")) variant = "slate";
		else {
			const variants = [
				"primary",
				"purple",
				"emerald",
				"sky",
				"amber",
				"rose",
				"slate",
			];
			variant = variants[role.id % variants.length];
		}

		return { label: role.name, variant };
	}

	// Error handling helpers
	let errorTimeout: any;
	function setError(msg: string) {
		error = msg;
		if (errorTimeout) clearTimeout(errorTimeout);
		errorTimeout = setTimeout(() => {
			error = "";
		}, 5000);
	}

	let passwordErrorTimeout: any;
	function setPasswordError(msg: string) {
		passwordError = msg;
		if (passwordErrorTimeout) clearTimeout(passwordErrorTimeout);
		passwordErrorTimeout = setTimeout(() => {
			passwordError = "";
		}, 5000);
	}

	// Handlers
	let searchTimeout: any;
	function handleSearch(query: string) {
		searchQuery = query;
		page = 1;
		if (searchTimeout) clearTimeout(searchTimeout);
		searchTimeout = setTimeout(loadUsers, 300);
	}

	function openCreateModal() {
		modalMode = "create";
		const defaultRole =
			roles.find((r) => r.name.toLowerCase() === "member") || roles[0];
		currentUser = {
			role_id: defaultRole?.id || 1,
			is_active: true,
			is_admin: false,
		}; // Default member
		tempPassword = "";
		showModal = true;
	}

	function openEditModal(user: User) {
		modalMode = "edit";
		currentUser = { ...user };
		showModal = true;
	}

	function openResetModal(user: User) {
		modalMode = "reset";
		currentUser = { ...user };
		tempPassword = "";
		showModal = true;
	}

	function triggerToggleStatus(user: User) {
		confirmData = {
			title: user.is_active ? "Lock account" : "Unlock account",
			message: `Are you sure you want to ${user.is_active ? "lock" : "unlock"} account ${user.email}? The user will ${user.is_active ? "not" : ""} be able to login to the system.`,
			// execute the async toggle process
			action: async () => {
				isSubmitting = true;
				try {
					await userService.toggleStatus(user.id);
					toast.success(
						`${user.is_active ? "Lock" : "Unlock"} account successfully`,
					);
					await loadUsers();
					showConfirm = false;
				} catch (e: any) {
					if (!e.isAuthError) {
						toast.error(e.message);
					}
					showConfirm = false;
				} finally {
					isSubmitting = false;
				}
			},
		};
		showConfirm = true;
	}

	async function handleModalSubmit(e: Event) {
		e.preventDefault();
		passwordError = "";

		// Validate password if creating or resetting
		if (modalMode === "create" || modalMode === "reset") {
			const err = validatePassword(tempPassword);
			if (err) {
				setPasswordError(err);
				return;
			}
		}

		isSubmitting = true;
		try {
			if (modalMode === "create") {
				await userService.createUser(currentUser, tempPassword);
			} else if (modalMode === "edit") {
				if (!currentUser.id) throw new Error("Missing User ID");
				await userService.updateUser(currentUser.id, currentUser);
			} else if (modalMode === "reset") {
				if (!currentUser.id) throw new Error("Missing User ID");
				await userService.resetPassword(currentUser.id, tempPassword);
			}
			showModal = false;
			toast.success(
				modalMode === "create"
					? "New user created successfully"
					: modalMode === "edit"
						? "Information updated successfully"
						: "Password reset successfully",
			);
			await loadUsers();
		} catch (e: any) {
			if (!e.isAuthError) {
				setError(e.message);
			}
		} finally {
			isSubmitting = false;
		}
	}
</script>

<ContentHeader
	title="User Management"
	subtitle="List of system accounts and access permissions"
>
	<div class="flex items-center gap-3">
		<div
			class="flex bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl px-4 py-2 w-64 items-center gap-2 focus-within:ring-2 focus-within:ring-primary/20 focus-within:border-primary transition-all"
		>
			<span class="material-symbols-outlined text-slate-400 text-sm"
				>search</span
			>
			<input
				class="bg-transparent border-none focus:ring-0 text-sm w-full outline-none placeholder:text-slate-400 text-slate-900 dark:text-white"
				placeholder="Search members..."
				type="text"
				oninput={(e) => handleSearch(e.currentTarget.value)}
			/>
		</div>
		{#if hasPermission($authStore.user, $authStore.token, "user:create")}
			<Button icon="add" onclick={openCreateModal}>Add Member</Button>
		{/if}
	</div>
</ContentHeader>
<div class="space-y-4">
	<DataTable
		items={users}
		{columns}
		{loading}
		{total}
		{page}
		{limit}
		{onPageChange}
		onLimitChange={(l) => {
			limit = l;
			page = 1;
			loadUsers();
		}}
	>
		{#snippet rowCell({ item, column })}
			{#if column.key === "user"}
				<div class="flex items-center gap-3">
					<Avatar name={item.full_name} size="md" />
					<div>
						<p
							class="font-bold text-slate-900 text-sm leading-tight"
						>
							{item.full_name}
						</p>
						<p
							class="text-[11px] text-slate-400 font-medium mt-0.5"
						>
							{item.email}
						</p>
					</div>
				</div>
			{:else if column.key === "email"}
				<span class="text-sm font-medium text-slate-600"
					>{item.email}</span
				>
			{:else if column.key === "role"}
				{@const role = getRoleInfo(item.role_id)}
				<Badge variant={role.variant as any}>{role.label}</Badge>
			{:else if column.key === "is_admin"}
				{#if item.is_admin}
					<Badge
						variant="rose"
						class="!bg-rose-500 !text-white !border-none !text-[10px] !px-2 !py-0.5"
						>ADMIN</Badge
					>
				{:else}
					<span class="text-[11px] font-medium text-slate-300 ml-2"
						>—</span
					>
				{/if}
			{:else if column.key === "status"}
				<div class="flex items-center gap-2">
					<div
						class="size-2 rounded-full {item.is_active
							? 'bg-emerald-500'
							: 'bg-rose-500'} shadow-sm"
					></div>
					<span class="text-sm font-semibold text-slate-700"
						>{item.is_active ? "Active" : "Locked"}</span
					>
				</div>
			{:else if column.key === "actions"}
				<div
					class="flex items-center justify-end gap-1 opacity-0 group-hover:opacity-100 transition-opacity"
				>
					{#if hasPermission($authStore.user, $authStore.token, "user:update")}
						<button
							class="p-2 text-slate-400 hover:text-sky-500 hover:bg-sky-50 rounded-lg transition-all"
							title="Edit"
							onclick={() => openEditModal(item)}
						>
							<span class="material-symbols-outlined text-[18px]"
								>edit</span
							>
						</button>
						<button
							class="p-2 text-slate-400 hover:text-amber-500 hover:bg-amber-50 rounded-lg transition-all"
							title="Reset Password"
							onclick={() => openResetModal(item)}
						>
							<span class="material-symbols-outlined text-[18px]"
								>key</span
							>
						</button>
						<button
							class="p-2 text-slate-400 hover:text-rose-500 hover:bg-rose-50 rounded-lg transition-all"
							title={item.is_active ? "Lock Account" : "Unlock"}
							onclick={() => triggerToggleStatus(item)}
						>
							<span class="material-symbols-outlined text-[18px]"
								>{item.is_active ? "lock" : "lock_open"}</span
							>
						</button>
					{/if}
				</div>
			{/if}
		{/snippet}
	</DataTable>
</div>

<!-- Modal Component -->
<Modal
	show={showModal}
	onClose={() => (showModal = false)}
	title={modalMode === "create"
		? "Add New Member"
		: modalMode === "edit"
			? "Edit Information"
			: "Reset Password"}
>
	<form id="user-form" onsubmit={handleModalSubmit} class="space-y-5">
		{#if error}
			<div
				class="p-3 bg-rose-50 border border-rose-100 text-rose-600 rounded-xl text-xs font-medium flex items-start gap-2 animate-in fade-in slide-in-from-top-2 duration-300"
			>
				<span class="material-symbols-outlined text-sm mt-0.5"
					>error</span
				>
				<span class="flex-1">{error}</span>
			</div>
		{/if}
		{#if modalMode === "create" || modalMode === "edit"}
			<div class="space-y-1.5 flex flex-col">
				<label
					for="fullName"
					class="text-[11px] font-bold uppercase tracking-wider text-slate-500 ml-1"
					>Full Name</label
				>
				<input
					id="fullName"
					type="text"
					bind:value={currentUser.full_name}
					required
					class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all text-sm font-medium"
				/>
			</div>

			<div class="space-y-1.5 flex flex-col">
				<label
					for="email"
					class="text-[11px] font-bold uppercase tracking-wider text-slate-500 ml-1"
					>Email</label
				>
				<input
					id="email"
					type="email"
					bind:value={currentUser.email}
					required
					class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all text-sm font-medium"
				/>
			</div>

			<div class="space-y-1.5 flex flex-col">
				<label
					for="role"
					class="text-[11px] font-bold uppercase tracking-wider text-slate-500 ml-1"
					>Role</label
				>
				<select
					id="role"
					bind:value={currentUser.role_id}
					class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all text-sm font-medium"
				>
					{#each roles as role}
						<option value={role.id}>{role.name}</option>
					{/each}
				</select>
			</div>

			{#if $authStore.user?.isAdmin}
				<div
					class="flex items-center gap-3 p-4 bg-slate-50 rounded-2xl border border-slate-100 hover:border-primary/20 transition-all cursor-pointer"
					onclick={() =>
						(currentUser.is_admin = !currentUser.is_admin)}
				>
					<div class="flex-1">
						<p
							class="text-sm font-bold text-slate-900 leading-tight"
						>
							System Administrator
						</p>
						<p
							class="text-[11px] text-slate-400 font-medium mt-0.5"
						>
							Allow user to bypass all permission barriers
						</p>
					</div>
					<div
						class="size-6 bg-white border-2 {currentUser.is_admin
							? 'border-primary bg-primary/10'
							: 'border-slate-300'} rounded-lg flex items-center justify-center transition-all"
					>
						{#if currentUser.is_admin}
							<span
								class="material-symbols-outlined text-primary text-lg font-bold"
								>check</span
							>
						{/if}
					</div>
				</div>
			{/if}
		{/if}

		{#if modalMode === "create" || modalMode === "reset"}
			<div
				class="space-y-1.5 flex flex-col pt-2 border-t border-slate-100 mt-2"
			>
				<label
					for="password"
					class="text-[11px] font-bold uppercase tracking-wider text-slate-500 ml-1"
					>Password {modalMode === "reset" ? "new" : ""}</label
				>
				<div class="relative">
					<input
						id="password"
						type={showPassword ? "text" : "password"}
						bind:value={tempPassword}
						required
						placeholder="Enter password"
						class="w-full bg-slate-50 border {passwordError
							? 'border-rose-400'
							: 'border-slate-200'} rounded-xl px-4 py-3 pr-12 focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all text-sm font-medium"
						oninput={() => (passwordError = "")}
					/>
					<button
						type="button"
						class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition-colors"
						onclick={() => (showPassword = !showPassword)}
					>
						<span class="material-symbols-outlined text-[18px]"
							>{showPassword
								? "visibility_off"
								: "visibility"}</span
						>
					</button>
				</div>
				{#if passwordError}
					<p
						class="text-xs text-rose-500 font-medium mt-1 ml-1 flex items-center gap-1"
					>
						<span class="material-symbols-outlined text-sm"
							>error</span
						>
						{passwordError}
					</p>
				{/if}
			</div>
		{/if}
	</form>

	{#snippet footer()}
		<Button
			variant="outline"
			class="flex-1"
			onclick={() => (showModal = false)}>Cancel</Button
		>
		<Button
			type="submit"
			form="user-form"
			class="flex-[2] justify-center"
			disabled={isSubmitting}
		>
			{isSubmitting
				? "Processing..."
				: modalMode === "create"
					? "Create Account"
					: "Update"}
		</Button>
	{/snippet}
</Modal>

<ConfirmDialog
	show={showConfirm}
	title={confirmData.title}
	message={confirmData.message}
	onConfirm={confirmData.action}
	onCancel={() => (showConfirm = false)}
/>
