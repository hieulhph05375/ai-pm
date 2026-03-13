<script lang="ts">
	import { page } from "$app/state";
	import { authStore, authService } from "$lib/services/auth";
	import Avatar from "../ui/Avatar.svelte";
	import Input from "../ui/Input.svelte";
	import NotificationBell from "./NotificationBell.svelte";

	let searchQuery = $state("");

	// Improved breadcrumb logic
	const breadcrumbs = $derived.by(() => {
		const paths = page.url.pathname.split("/").filter((p) => p !== "");
		const bcs = paths.map((p, i) => {
			const label = p.charAt(0).toUpperCase() + p.slice(1);
			const href = "/" + paths.slice(0, i + 1).join("/");
			return { label, href };
		});

		// Always start with Dashboard
		if (page.url.pathname === "/")
			return [{ label: "Dashboard", href: "/" }];
		return [{ label: "Dashboard", href: "/" }, ...bcs];
	});
</script>

<header
	class="h-20 flex items-center justify-between px-10 glass-header border-b border-slate-100 z-10 shrink-0"
>
	<div class="flex items-center gap-6">
		<div class="flex items-center gap-2 text-sm font-medium">
			{#each breadcrumbs as bc, i}
				{#if i > 0}
					<span
						class="material-symbols-outlined text-slate-300 text-sm"
						>chevron_right</span
					>
				{/if}
				<a
					href={bc.href}
					class={i === breadcrumbs.length - 1
						? "text-slate-900"
						: "text-slate-400 hover:text-primary transition-colors"}
				>
					{bc.label}
				</a>
			{/each}
		</div>

		<div class="h-6 w-px bg-slate-200"></div>

		<div class="w-80">
			<Input
				placeholder="Search tasks, team members..."
				icon="search"
				bind:value={searchQuery}
				wrapperClass="!space-y-0"
				class="!py-2 !bg-slate-100/50 !border-none"
			/>
		</div>
	</div>

	<div class="flex items-center gap-4">
		<NotificationBell />
		<button
			class="size-10 flex items-center justify-center rounded-xl bg-white border border-slate-200 text-slate-600 hover:text-primary transition-colors"
		>
			<span class="material-symbols-outlined">chat_bubble</span>
		</button>

		<div
			class="flex items-center gap-3 pl-4 ml-2 border-l border-slate-200"
		>
			{#if $authStore.user}
				<div class="text-right">
					<p class="text-sm font-bold text-slate-900 leading-none">
						{$authStore.user.fullName}
					</p>
					<p
						class="text-[10px] font-medium text-slate-400 mt-1 uppercase tracking-tight"
					>
						Admin
					</p>
				</div>
				<Avatar name={$authStore.user.fullName} size="sm" />
				<button
					class="p-2 text-slate-400 hover:text-rose-500 hover:bg-rose-50 rounded-lg transition-all ml-1"
					title="Logout"
					onclick={() => authService.logout()}
				>
					<span class="material-symbols-outlined text-[18px]"
						>logout</span
					>
				</button>
			{:else}
				<div class="text-right">
					<p class="text-sm font-bold text-slate-900 leading-none">
						Guest
					</p>
				</div>
				<Avatar name="Guest" size="sm" />
			{/if}
		</div>
	</div>
</header>

<style>
	:global(.glass-header) {
		background: rgba(255, 255, 255, 0.8);
		backdrop-filter: blur(8px);
	}
</style>
