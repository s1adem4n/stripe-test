<script lang="ts">
	import type { CheckoutParams } from '$lib/api';
	import { CustomerInput } from '$lib/components';
	import PinInput from '$lib/components/pin-input.svelte';

	let checkoutParams: CheckoutParams = $state({
		verificationCode: '',
		name: '',
		email: '',
		phone: '',
		address: {
			line1: '',
			line2: '',
			city: '',
			postalCode: ''
		}
	});
</script>

<div class="mx-auto flex w-full max-w-xl flex-col gap-4 p-4">
	<PinInput length={6} />
	<CustomerInput
		bind:name={checkoutParams.name}
		bind:phone={checkoutParams.phone}
		bind:address={checkoutParams.address}
		allowedCities={[
			{ postalCode: '24568', city: 'Kaltenkirchen' },
			{ postalCode: '24628', city: 'Hartenholm' }
		]}
		onsubmit={(name, phone, address) => {
			console.log(name, phone, address);
		}}
	/>
</div>
