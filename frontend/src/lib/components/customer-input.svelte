<script lang="ts">
	import { Input, Label, Select } from './ui';
	import { AddressAutocomplete } from '.';
	import type { Address } from '$lib/api';
	import Button from './ui/button.svelte';

	interface CustomerInputProps {
		name: string;
		phone: string;
		address: Address;
		allowedCities: { postalCode: string; city: string }[];
		onsubmit: (name: string, phone: string, address: Address) => void;
	}
	let {
		name = $bindable(),
		phone = $bindable(),
		address = $bindable(),
		allowedCities,
		onsubmit
	}: CustomerInputProps = $props();

	let addressStage: 'autocomplete' | 'manual' = $state('autocomplete');
	let error = $state('');

	const validateCity = (postalCode: string, city: string) => {
		if (!allowedCities.some((c) => c.postalCode === postalCode && c.city === city)) {
			error = 'Wir liefern aktuell leider nicht nach ' + city;
			return false;
		} else {
			error = '';
			return true;
		}
	};
	const formatCity = (postalCode: string, city: string) => `${postalCode} ${city}`;
</script>

<form
	class="flex flex-col gap-3"
	onsubmit={(e) => {
		e.preventDefault();
		onsubmit(name, phone, address);
	}}
>
	<div class="flex flex-col gap-1">
		<Label for="name">Name</Label>
		<Input type="text" id="name" bind:value={name} placeholder="Max Mustermann" required />
	</div>
	<div class="flex flex-col gap-1">
		<Label for="phone">Telefonnummer</Label>
		<Input type="tel" id="phone" bind:value={phone} placeholder="+49 123 4567890" required />
	</div>

	{#if addressStage === 'autocomplete'}
		<div class="flex flex-col gap-1">
			<Label for="autocomplete">Adresse</Label>
			<AddressAutocomplete
				placeholder="Straße, Hausnummer, Ort"
				onaddress={(res) => {
					if (validateCity(res.postalCode, res.city)) {
						address = res;
						addressStage = 'manual';
					}
				}}
			/>
			<button
				class="text-left text-sm text-blue-500 hover:underline"
				onclick={() => (addressStage = 'manual')}
			>
				Adresse manuell eingeben
			</button>
		</div>
	{:else if addressStage === 'manual'}
		<div class="flex flex-col gap-1">
			<Label for="address">Adresse</Label>
			<Input
				type="text"
				id="address"
				bind:value={address.line1}
				placeholder="Straße, Hausnummer"
				required
			/>
		</div>
		<div class="flex flex-col gap-1">
			<Label for="address-line2">Adresszusatz (optional)</Label>
			<Input type="text" id="address-line2" bind:value={address.line2} placeholder="Hinterhaus" />
		</div>
		<div class="flex flex-col gap-1">
			<Label for="city">Ort</Label>
			<Select
				id="city"
				options={allowedCities.map((c) => ({
					value: formatCity(c.postalCode, c.city),
					label: formatCity(c.postalCode, c.city)
				}))}
				oninput={(e) => {
					const [postalCode, city] = e.currentTarget.value.split(' ', 2);
					address.postalCode = postalCode;
					address.city = city;
				}}
				value={formatCity(address.postalCode, address.city)}
				required
			></Select>
		</div>
	{/if}
	<Button disabled={error.length > 0} type="submit">Bezahlen</Button>

	{#if error}
		<p class="text-sm text-red-500">{error}</p>
	{/if}
</form>
