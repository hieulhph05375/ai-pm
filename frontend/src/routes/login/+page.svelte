<script lang="ts">
	import { authService } from "$lib/services/auth";
	import { goto } from "$app/navigation";
	import { page } from "$app/stores";

	let email = $state("");
	let password = $state("");
	let error = $state("");
	let isLoading = $state(false);

	async function handleLogin(e: Event) {
		e.preventDefault();
		error = "";
		isLoading = true;
		const success = await authService.login(email, password);
		if (success) {
			// Get redirect path from query param or default to '/'
			const redirectTo = $page.url.searchParams.get("redirectTo") || "/";
			goto(redirectTo);
		} else {
			error = "Login failed. Please check your email and password.";
		}
		isLoading = false;
	}
</script>

<div
	class="flex h-screen w-full font-display bg-background-light dark:bg-background-dark text-slate-900 dark:text-slate-100 antialiased overflow-hidden"
>
	<!-- Left Side: Abstract Branding Illustration -->
	<div
		class="hidden lg:flex lg:w-1/2 relative overflow-hidden gradient-bg items-center justify-center p-20"
	>
		<div
			class="absolute inset-0 opacity-20"
			style="background-image: radial-gradient(circle at 2px 2px, white 1px, transparent 0); background-size: 40px 40px;"
		></div>
		<div class="relative z-10 text-white space-y-8 max-w-lg">
			<div class="flex items-center gap-3">
				<div
					class="size-10 bg-white/20 backdrop-blur-md rounded-lg flex items-center justify-center border border-white/30"
				>
					<span
						class="material-symbols-outlined text-white text-[18px]"
						>rocket_launch</span
					>
				</div>
				<span class="text-2xl font-outfit font-extrabold tracking-tight"
					>AIVOVAN_2026</span
				>
			</div>
			<div class="space-y-4">
				<h1 class="text-5xl font-outfit font-extrabold leading-tight">
					Orchestrate your SaaS ecosystem.
				</h1>
				<p class="text-white/80 text-lg leading-relaxed">
					The next generation of high-end project management. Secure,
					scalable, and beautifully designed for peak performance.
				</p>
			</div>
			<div class="flex gap-4 pt-4">
				<div class="h-1 w-12 bg-white rounded-full"></div>
				<div class="h-1 w-4 bg-white/30 rounded-full"></div>
				<div class="h-1 w-4 bg-white/30 rounded-full"></div>
			</div>
		</div>
		<!-- Decorative Shapes -->
		<div
			class="absolute -bottom-20 -left-20 size-80 bg-white/10 rounded-full blur-3xl"
		></div>
		<div
			class="absolute -top-20 -right-20 size-96 bg-sky-400/20 rounded-full blur-3xl"
		></div>
	</div>

	<!-- Right Side: Login Form -->
	<div
		class="w-full lg:w-1/2 flex items-center justify-center p-6 sm:p-12 relative bg-slate-50 dark:bg-background-dark"
	>
		<!-- Background Image Pattern for Depth -->
		<div
			class="absolute inset-0 opacity-[0.03] pointer-events-none"
			style="background-image: url('https://lh3.googleusercontent.com/aida-public/AB6AXuDT6fAxJHXhrg3sNRcuiu2tLT--IOZKflLdDgaK7TQFhkwgXrDSjVKYzpr8dxnqrJqRboeX8VUFd1oCX254CYKx0enuZtT3zKNAZ4iNXwQ58-ervzYZVrzCjWCaeTpBuj9KF3ogCWRo2Qf2gTqiuPeeaQjTJZvOh7VC37VXxr8OCiuZgiBz7SoXDBSUeRrSsV53-N2ieyGxvoKyV64N3iyVEPdmL0EVbUqlV5l9qlO8feve_EfB4QZk87ygH-IEzn9fk9_IjHahyX49');"
		></div>

		<div
			class="glass-card w-full max-w-[480px] p-10 rounded-xl shadow-2xl space-y-8 relative z-10"
		>
			<div class="space-y-2 text-center lg:text-left">
				<h2
					class="text-3xl font-outfit font-extrabold text-slate-900 dark:text-white"
				>
					Welcome Back
				</h2>
				<p class="text-slate-500 dark:text-slate-400 text-sm">
					Please enter your credentials to access your workspace.
				</p>
			</div>

			{#if error}
				<div
					class="p-4 text-sm text-red-700 bg-red-100 rounded-lg animate-in fade-in slide-in-from-top-4 duration-300"
					role="alert"
				>
					{error}
				</div>
			{/if}

			<form onsubmit={handleLogin} class="space-y-6">
				<div class="space-y-4">
					<div class="space-y-1.5">
						<label
							for="email"
							class="text-sm font-semibold text-slate-700 dark:text-slate-300 ml-1"
							>Email Address</label
						>
						<div class="relative">
							<span
								class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[18px]"
								>mail</span
							>
							<input
								bind:value={email}
								id="email"
								type="email"
								required
								class="w-full pl-10 pr-4 py-3.5 bg-white/50 dark:bg-slate-900/50 border border-slate-200 dark:border-slate-700 rounded-lg focus:ring-2 focus:ring-primary/50 focus:border-primary outline-none transition-all text-slate-900 dark:text-white placeholder:text-slate-400"
								placeholder="name@company.com"
							/>
						</div>
					</div>

					<div class="space-y-1.5">
						<div class="flex justify-between items-center px-1">
							<label
								for="password"
								class="text-sm font-semibold text-slate-700 dark:text-slate-300"
								>Password</label
							>
							<a
								class="text-xs font-bold text-primary hover:underline"
								href="/">Forgot Password?</a
							>
						</div>
						<div class="relative">
							<span
								class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[18px]"
								>lock</span
							>
							<input
								bind:value={password}
								id="password"
								type="password"
								required
								class="w-full pl-10 pr-12 py-3.5 bg-white/50 dark:bg-slate-900/50 border border-slate-200 dark:border-slate-700 rounded-lg focus:ring-2 focus:ring-primary/50 focus:border-primary outline-none transition-all text-slate-900 dark:text-white placeholder:text-slate-400"
								placeholder="••••••••"
							/>
							<button
								type="button"
								class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 transition-colors"
							>
								<span class="material-symbols-outlined text-[18px]"
									>visibility</span
								>
							</button>
						</div>
					</div>
				</div>

				<button
					type="submit"
					disabled={isLoading}
					class="w-full gradient-bg text-white font-bold py-4 rounded-lg shadow-lg shadow-primary/20 hover:shadow-primary/40 hover:opacity-95 transition-all active:scale-[0.98] flex items-center justify-center gap-2"
				>
					{#if isLoading}
						<div
							class="h-5 w-5 border-2 border-white/30 border-t-white rounded-full animate-spin"
						></div>
					{:else}
						<span>Sign In</span>
						<span class="material-symbols-outlined text-lg"
							>arrow_forward</span
						>
					{/if}
				</button>
			</form>
		</div>

		<div class="absolute bottom-8 text-xs text-slate-400 flex gap-6">
			<a class="hover:text-primary transition-colors" href="/#"
				>Privacy Policy</a
			>
			<a class="hover:text-primary transition-colors" href="/#"
				>Terms of Service</a
			>
			<a class="hover:text-primary transition-colors" href="/#">Support</a
			>
		</div>
	</div>
</div>
