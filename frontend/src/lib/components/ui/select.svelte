<script lang="ts">
	import type { HTMLSelectAttributes } from 'svelte/elements';
	import { cn } from '$lib/utils';

	interface Option {
		value: string;
		label: string;
	}
	interface SelectProps extends HTMLSelectAttributes {
		options: Option[];
	}

	let { options, class: className, value = $bindable(), ...props }: SelectProps = $props();
</script>

<select
	class={cn(
		'border-input ring-offset-background placeholder:text-muted-foreground focus:ring-ring flex h-10 w-full items-center justify-between whitespace-nowrap rounded-md border bg-transparent px-3 py-2 shadow-sm focus:outline-none focus:ring-1 disabled:cursor-not-allowed disabled:opacity-50 [&>span]:line-clamp-1',
		className
	)}
	bind:value
	{...props}
>
	{#each options as option}
		<option value={option.value}>{option.label}</option>
	{/each}
</select>
