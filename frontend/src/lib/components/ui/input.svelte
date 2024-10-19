<script lang="ts">
	import type { HTMLInputAttributes } from 'svelte/elements';
	import { cn } from '$lib/utils';

	interface InputProps extends HTMLInputAttributes {}

	let { children, class: className, value = $bindable(), ...props }: InputProps = $props();
</script>

<input
	class={cn(
		'border-input file:text-foreground placeholder:text-muted-foreground focus-visible:ring-ring flex h-10 w-full rounded-md border bg-transparent px-3 py-1 shadow-sm transition-colors file:border-0 file:bg-transparent file:font-medium focus-visible:outline-none focus-visible:ring-1 disabled:cursor-not-allowed disabled:opacity-50',
		className
	)}
	bind:value
	{...props}
/>
