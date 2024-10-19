<script lang="ts">
	import { Input } from '$lib/components/ui';
	import {
		autocompleteAddress,
		getPlaceDetails,
		type Address,
		type AddressPrediction
	} from '$lib/api';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import { slide } from 'svelte/transition';

	interface AddressAutocompleteProps extends HTMLInputAttributes {
		onaddress: (address: Address) => void;
	}
	let {
		onaddress,
		class: className,
		value = $bindable(''),
		...props
	}: AddressAutocompleteProps = $props();

	const wait = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

	const waitTime = 1200 as const;

	let predictions: AddressPrediction[] = $state([]);
	let loading = $state(false);
	let requestWait: Promise<any> = $state(wait(0));
	let focused = $state(false);
	let focusedPrediction = $state(-1);
</script>

<div class="relative flex">
	<Input
		class={className}
		{...props}
		bind:value
		onkeydown={(e) => {
			if (e.key === 'ArrowDown') {
				e.preventDefault();
				focusedPrediction = Math.min(focusedPrediction + 1, predictions.length - 1);
			} else if (e.key === 'ArrowUp') {
				e.preventDefault();
				focusedPrediction = Math.max(focusedPrediction - 1, -1);
			} else if (e.key === 'Enter' && focusedPrediction >= 0) {
				e.preventDefault();
				getPlaceDetails(predictions[focusedPrediction].placeId).then(onaddress);
				value = predictions[focusedPrediction].description;
				focusedPrediction = -1;
			}
		}}
		oninput={() => {
			if (loading || !value) return;
			loading = true;
			requestWait
				.then(() => {
					autocompleteAddress(value).then((res) => {
						focusedPrediction = -1;
						predictions = res || [];
					});
				})
				.finally(() => {
					requestWait = wait(waitTime);
					loading = false;
				});
		}}
		onfocus={() => (focused = true)}
		onblur={() => (focused = false)}
	/>
	{#if focused}
		<div
			class="border-border absolute top-12 z-10 flex w-full flex-col overflow-hidden rounded-md rounded-b-md border bg-white shadow-sm transition-[height] duration-200"
			style="height: {Math.max(predictions.length, 1) * 40 + predictions.length - 1}px"
		>
			{#each predictions as prediction, i}
				<button
					class="hover:bg-accent min-h-10 cursor-pointer truncate px-2 text-left transition-colors"
					class:bg-gray-100={i === focusedPrediction}
					onmousedown={(e) => e.preventDefault()}
					ontouchstart={(e) => e.preventDefault()}
					onclick={async (e) => {
						e.stopImmediatePropagation();

						value = prediction.description;
						await requestWait;
						requestWait = wait(waitTime);
						getPlaceDetails(prediction.placeId).then(onaddress);
						focusedPrediction = -1;

						if (e.currentTarget) {
							const inputElement = e.currentTarget.closest('input');
							if (inputElement) {
								inputElement.blur();
							}
						}
					}}
				>
					{prediction.description}
				</button>
				{#if i < predictions.length - 1}
					<div class="border-border border-t"></div>
				{/if}
			{:else}
				<p class="p-2 text-left text-muted-foreground">Keine Vorschläge</p>
			{/each}
		</div>
	{/if}
</div>
