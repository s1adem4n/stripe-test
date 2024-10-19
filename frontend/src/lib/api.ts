export const BASE_URL = import.meta.env.DEV
	? window.location.origin.replace('5173', '8090')
	: window.location.origin;

export const API_URL = `${BASE_URL}/api`;

export interface CheckoutParams {
	verificationCode: string;
	email: string;
	name: string;
	phone: string;
	address: Address;
}

export interface Address {
	line1: string;
	line2: string;
	city: string;
	postalCode: string;
}
export interface AddressPrediction {
	placeId: string;
	description: string;
}

export async function createCheckoutSession(params: CheckoutParams) {
	const response = await fetch(`${API_URL}/checkout`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(params)
	});

	return response.json() as Promise<{ url: string }>;
}

export async function sendVerificationCode(email: string) {
	const response = await fetch(`${API_URL}/verify?email=${email}`);
	if (!response.ok) {
		throw new Error('Failed to send verification code');
	}

	return;
}

export async function autocompleteAddress(input: string) {
	const response = await fetch(`${API_URL}/autocomplete-address?input=${input}`);
	if (!response.ok) {
		throw new Error('Failed to autocomplete address');
	}

	return (response.json() || []) as Promise<AddressPrediction[]>;
}
export async function getPlaceDetails(placeId: string) {
	const response = await fetch(`${API_URL}/place-details?placeId=${placeId}`);
	if (!response.ok) {
		throw new Error('Failed to get place details');
	}

	return response.json() as Promise<Address>;
}
