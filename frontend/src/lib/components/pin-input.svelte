<script lang="ts">
	interface PinInputProps {
		length: number;
		value?: string[];
	}
	let { length, value = $bindable([]) }: PinInputProps = $props();

	let inputs: HTMLInputElement[] = $state([]);
</script>

<div class="flex gap-1">
	{#each Array.from({ length }) as _, i}
		<input
			bind:this={inputs[i]}
			type="text"
			inputmode="numeric"
			pattern="[0-9]*"
			maxlength="2"
			size="1"
			class="border-border ring-ring h-10 w-10 rounded-md border text-center shadow-sm focus:outline-none focus:ring-1"
			aria-label="Pin input"
			bind:value={value[i]}
			onkeydown={(e) => {
				if (e.key === 'Backspace' && i > 0 && !e.currentTarget.value) {
					const prevInput = inputs[i - 1];
					prevInput.focus();
				}

				if (e.key === 'ArrowLeft' && i > 0) {
					const prevInput = inputs[i - 1];
					prevInput.focus();
				} else if (e.key === 'ArrowRight' && i < length - 1) {
					const nextInput = inputs[i + 1];
					nextInput.focus();
				}
			}}
			oninput={(e) => {
				if (/^\d+$/.test(e.currentTarget.value)) {
					value[i] = e.currentTarget.value.slice(-1);
					if (i < length - 1) {
						const nextInput = inputs[i + 1];
						nextInput.focus();
					}
				} else {
					// try to find a number in the input
					const number = e.currentTarget.value.match(/\d/);
					if (number) {
						value[i] = number[0];
						if (i < length - 1) {
							const nextInput = inputs[i + 1];
							nextInput.focus();
						}
					} else {
						value[i] = '';
					}
				}
			}}
			onpaste={(e) => {
				e.preventDefault();
				if (!e.clipboardData) return;

				const pasted = e.clipboardData.getData('text').split('');
				for (let j = 0; j < length; j++) {
					if (pasted[j] && /^\d$/.test(pasted[j])) {
						value[j] = pasted[j];
					}
				}
			}}
		/>
		{#if i === length / 2 - 1}
			<div class="w-2"></div>
		{/if}
	{/each}
</div>
