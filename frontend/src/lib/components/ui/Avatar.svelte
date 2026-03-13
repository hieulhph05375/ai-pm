<script lang="ts">
	interface Props {
		src?: string;
		alt?: string;
		name?: string;
		size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl';
		class?: string;
	}

	let {
		src,
		alt = 'User avatar',
		name = '',
		size = 'md',
		class: className = ''
	}: Props = $props();

	const sizes = {
		xs: 'size-6 text-[8px]',
		sm: 'size-8 text-[10px]',
		md: 'size-10 text-xs',
		lg: 'size-14 text-base',
		xl: 'size-20 text-xl'
	};

	let initials = $derived((name || '')
		.split(' ')
		.map(n => n[0])
		.join('')
		.toUpperCase()
		.slice(0, 2));
</script>

<div class="relative shrink-0 {sizes[size]} rounded-xl overflow-hidden ring-2 ring-white shadow-sm {className}">
	{#if src}
		<img {src} {alt} class="w-full h-full object-cover" />
	{:else}
		<div class="w-full h-full gradient-bg flex items-center justify-center text-white font-bold tracking-tighter">
			{initials || '?'}
		</div>
	{/if}
</div>
