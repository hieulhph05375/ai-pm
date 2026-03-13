<script lang="ts">
	import { page } from "$app/state";
	import { authStore, authService } from "$lib/services/auth";
	import { hasPermission } from "$lib/utils/permission";

	const navItems = [
		{
			name: "Portfolio",
			icon: "donut_large",
			href: "/portfolio",
			requiredPerm: "portfolio:read",
		},
		{
			name: "Projects",
			icon: "account_tree",
			href: "/projects",
			requiredPerm: "project:read",
			children: [
				{
					name: "Project Management",
					href: "/projects",
					icon: "list_alt",
					requiredPerm: "project:read",
				},
				{
					name: "Stakeholders",
					href: "/stakeholders",
					icon: "groups",
					requiredPerm: "stakeholder:read",
				},
				{
					name: "Holidays",
					href: "/holidays",
					icon: "event_busy",
					requiredPerm: "holiday:read",
				},
			],
		},
		{
			name: "Tasks",
			icon: "task_alt",
			href: "/tasks",
			requiredPerm: "task:read",
		},
		{
			name: "Timesheets",
			icon: "schedule",
			href: "/timesheets",
			requiredPerm: "timesheet:read",
		},
		{
			name: "Workload",
			icon: "calendar_view_month",
			href: "/resources",
			requiredPerm: "resource:read",
		},
		{
			name: "Category Config",
			href: "/categories",
			icon: "category",
			requiredPerm: "category:read",
		},
		{
			name: "Users",
			icon: "group",
			href: "/users",
			requiredPerm: "user:read",
		},
		{
			name: "Roles & Permissions",
			href: "/roles",
			icon: "shield_person",
			requiredPerm: "role:read",
		},
	];

	let expandedItem = $state<string | null>(null);

	$effect(() => {
		for (const item of navItems) {
			if (item.children && isActive(item.href)) {
				expandedItem = item.name;
				break;
			}
		}
	});

	function isActive(href: string) {
		if (href === "/") return page.url.pathname === "/";
		return page.url.pathname.startsWith(href);
	}

	let { collapsed = $bindable(false) } = $props();
</script>

<aside
	class="{collapsed
		? 'w-20'
		: 'w-72'} glass-sidebar flex flex-col h-full z-20 transition-all duration-300 relative"
>
	<div class="p-5 {collapsed ? 'px-3' : 'px-8'}">
		<div
			class="flex items-center {collapsed
				? 'justify-center'
				: 'justify-between'} mb-10"
		>
			{#if !collapsed}
				<div class="flex items-center gap-3">
					<div
						class="size-10 rounded-xl bg-gradient-to-br from-primary to-sky-400 flex items-center justify-center text-white shadow-lg shadow-primary/20 shrink-0"
					>
						<span class="material-symbols-outlined font-bold"
							>rocket_launch</span
						>
					</div>
					<div>
						<h1
							class="font-display font-bold text-xl tracking-tight text-slate-900 leading-none"
						>
							AIVOVAN
						</h1>
						<p
							class="text-[10px] font-bold uppercase tracking-widest text-primary/60 mt-1"
						>
							Management v1.0
						</p>
					</div>
				</div>
			{:else}
				<div
					class="size-10 rounded-xl bg-gradient-to-br from-primary to-sky-400 flex items-center justify-center text-white shadow-lg shadow-primary/20"
				>
					<span class="material-symbols-outlined font-bold"
						>rocket_launch</span
					>
				</div>
			{/if}

			<button
				onclick={() => (collapsed = !collapsed)}
				class="absolute -right-3 top-12 size-6 rounded-full bg-white border border-slate-200 shadow-sm flex items-center justify-center text-slate-400 hover:text-primary transition-all z-30"
			>
				<span class="material-symbols-outlined text-[18px]">
					{collapsed ? "chevron_right" : "chevron_left"}
				</span>
			</button>
		</div>

		<nav class="flex flex-col gap-1.5">
			{#each navItems as item}
				{#if !item.requiredPerm || hasPermission($authStore.user, $authStore.token, item.requiredPerm)}
					<div class="flex flex-col gap-1">
						<a
							href={item.children
								? "javascript:void(0)"
								: item.href}
							onclick={(e) => {
								if (item.children) {
									e.preventDefault();
									expandedItem =
										expandedItem === item.name
											? null
											: item.name;
								}
							}}
							title={collapsed ? item.name : ""}
							class="flex items-center {collapsed
								? 'justify-center'
								: 'gap-3'} px-4 py-3 rounded-xl transition-all duration-200 {isActive(
								item.href,
							)
								? 'bg-primary text-white shadow-md shadow-primary/20'
								: 'text-slate-600 hover:bg-white hover:shadow-sm'}"
						>
							<span class="material-symbols-outlined text-[18px]"
								>{item.icon}</span
							>
							{#if !collapsed}
								<span
									class="font-medium text-sm flex-1 whitespace-nowrap"
									>{item.name}</span
								>
								{#if item.children}
									<span
										class="material-symbols-outlined text-[18px] transition-transform duration-200 {expandedItem ===
										item.name
											? 'rotate-180'
											: ''}"
									>
										expand_more
									</span>
								{/if}
							{/if}
						</a>

						{#if item.children && expandedItem === item.name && !collapsed}
							<div
								class="flex flex-col gap-1 ml-9 mt-1 border-l border-slate-200/50 pl-2"
							>
								{#each item.children as child}
									{#if !child.requiredPerm || hasPermission($authStore.user, $authStore.token, child.requiredPerm)}
										<a
											href={child.href}
											class="flex items-center gap-2 px-3 py-2 rounded-lg text-xs transition-all duration-200 {page
												.url.pathname === child.href
												? 'text-primary font-bold bg-primary/5'
												: 'text-slate-500 hover:text-primary hover:bg-white'}"
										>
											<span
												class="material-symbols-outlined text-[16px]"
												>{child.icon}</span
											>
											<span>{child.name}</span>
										</a>
									{/if}
								{/each}
							</div>
						{/if}
					</div>
				{/if}
			{/each}

			<button
				onclick={() => authService.logout()}
				title={collapsed ? "Logout" : ""}
				class="w-full flex items-center {collapsed
					? 'justify-center'
					: 'gap-3'} px-4 py-3 rounded-xl transition-all duration-200 text-slate-400 hover:text-rose-500 hover:bg-rose-50 mt-1"
			>
				<span class="material-symbols-outlined text-[18px]">logout</span
				>
				{#if !collapsed}
					<span class="font-medium text-sm whitespace-nowrap"
						>Logout</span
					>
				{/if}
			</button>
		</nav>
	</div>
</aside>

<style>
	:global(.glass-sidebar) {
		background: rgba(255, 255, 255, 0.7);
		backdrop-filter: blur(16px);
		-webkit-backdrop-filter: blur(16px);
		border-right: 1px solid rgba(19, 55, 236, 0.1);
	}
</style>
