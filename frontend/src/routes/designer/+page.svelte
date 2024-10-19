<script lang="ts">
	import { Input } from '$lib/components/ui';
	import { draggable } from '@neodrag/svelte';

	let imageFile: File | null = $state(null);
</script>

<div class="mx-auto flex w-full max-w-xl flex-col gap-4 p-4">
	<Input
		type="file"
		accept="image/*"
		oninput={(e) => {
			const file = e.currentTarget.files?.[0];
			if (file) {
				imageFile = file;
			}
		}}
	/>
	<div
		class="border-border flex w-full items-center justify-center overflow-hidden rounded-md border shadow-sm"
	>
		<img src="/designer-wood.jpg" alt="wood" />
		{#if imageFile}
			<div
				use:draggable={{
					bounds: 'parent'
				}}
				class="absolute mix-blend-multiply"
			>
				<img
					use:draggable={{
						bounds: 'parent'
					}}
					class="grayscale-100 pointer-events-none h-auto w-24"
					src={URL.createObjectURL(imageFile)}
					alt="custom"
				/>
			</div>
		{/if}
	</div>
</div>
