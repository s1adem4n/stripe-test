<script lang="ts">
	import { createCheckoutSession, sendVerificationCode, type CheckoutParams } from '$lib/api';
	import { Button, Input, Label, Select } from '$lib/components/ui';

	const validCities = [
		['24568', 'Kaltenkirchen'],
		['24628', 'Hartenholm']
	];

	let checkoutParams: CheckoutParams = $state({
		verificationCode: '',
		name: '',
		email: '',
		phone: '',
		address: {
			line1: '',
			line2: '',
			city: validCities[0][1],
			postalCode: validCities[0][0]
		}
	});

	let emailSent = $state(false);
	let loading = $state(false);
	let invalidCode = $state(false);
</script>

<div class="mx-auto flex w-full max-w-xl flex-col gap-4 p-4">
	<h1 class="text-2xl font-bold">Checkout</h1>

	<form
		class="flex flex-col gap-2"
		onsubmit={(e) => {
			e.preventDefault();

			loading = true;
			sendVerificationCode(checkoutParams.email)
				.then(() => (emailSent = true))
				.finally(() => (loading = false));
		}}
	>
		<div class="flex flex-col gap-1">
			<Label for="email">Email</Label>
			<Input
				type="email"
				id="email"
				disabled={emailSent}
				bind:value={checkoutParams.email}
				placeholder="max@mustermann.de"
				required
			/>
		</div>
		<Button type="submit" disabled={emailSent || loading}>Verifizierungscode senden</Button>
		{#if emailSent}
			<p class="text-sm text-gray-500">
				Der Verifizierungscode wurde an die angegebene Email-Adresse gesendet.
			</p>
		{/if}
	</form>

	<form
		class="flex flex-col gap-2"
		onsubmit={(e) => {
			e.preventDefault();
			loading = true;
			createCheckoutSession(checkoutParams)
				.then((session) => {
					window.location.href = session.url;
				})
				.finally(() => {
					loading = false;
				})
				.catch(() => {
					invalidCode = true;
				});
		}}
	>
		<div class="flex flex-col gap-2">
			<div class="flex flex-col gap-1">
				<Label for="name">Name</Label>
				<Input
					type="text"
					id="name"
					bind:value={checkoutParams.name}
					placeholder="Max Mustermann"
					required
				/>
			</div>
			<div class="flex flex-col gap-1">
				<Label for="phone">Telefonnummer</Label>
				<Input
					type="tel"
					id="phone"
					bind:value={checkoutParams.phone}
					placeholder="+49 123 4567890"
					required
				/>
			</div>
			<div class="flex flex-col gap-1">
				<Label for="address-line1">Adresse</Label>
				<Input
					type="text"
					id="address-line1"
					bind:value={checkoutParams.address.line1}
					placeholder="Musterstraße 123"
					required
				/>
			</div>
			<div class="flex flex-col gap-1">
				<Label for="address-line2">Adresszusatz (optional)</Label>
				<Input
					type="text"
					id="address-line2"
					bind:value={checkoutParams.address.line2}
					placeholder="Hinterhaus"
				/>
			</div>
			<div class="flex flex-col gap-1">
				<Label for="city">Ort</Label>
				<Select
					id="city"
					bind:value={checkoutParams.address.city}
					oninput={(e) => {
						const postalCode = validCities.find(([_, city]) => city === e.currentTarget.value)?.[0];
						checkoutParams.address.postalCode = postalCode ?? '';
					}}
					options={validCities.map(([_, city]) => ({ value: city, label: city }))}
					required
				/>
			</div>
		</div>
		<p class="text-sm text-gray-500">
			Wir liefern aktuell nur in die folgenden Orte:
			{validCities.map(([postalCode, city]) => `${postalCode} ${city}`).join(', ')}
		</p>

		<div class="flex flex-col gap-1">
			<Label for="code">Verifizierungscode</Label>
			<Input
				type="text"
				id="code"
				bind:value={checkoutParams.verificationCode}
				placeholder="123456"
				oninput={() => {
					checkoutParams.verificationCode = checkoutParams.verificationCode
						.trim()
						.replaceAll(' ', '');
					const regex = /^[0-9]{6}$/;
					invalidCode = !regex.test(checkoutParams.verificationCode);
				}}
				required
			/>
		</div>
		{#if invalidCode}
			<p class="text-red-500">Ungültiger Code</p>
		{/if}

		<Button disabled={loading} type="submit">Jetzt bezahlen</Button>
	</form>
</div>
