<script lang="ts">
	import { Input } from '$lib/components/ui';
	import { onMount } from 'svelte';

	let imageFile: File | null = $state(null);
	let woodImage: HTMLImageElement | null = $state(null);
	let canvasImage: HTMLCanvasElement | null = $state(null);
	let canvas: HTMLCanvasElement = $state({} as HTMLCanvasElement);
	let ctx: CanvasRenderingContext2D = $state({} as CanvasRenderingContext2D);
	let transform = $state({ x: 0, y: 0, width: 0, height: 0 });

	let grabbing = $state(false);

	const grayscale = (image: HTMLImageElement, opacity = 1) => {
		const tempCanvas = document.createElement('canvas');
		const tempCtx = tempCanvas.getContext('2d');
		if (!tempCtx) return null;

		if (image.width > image.height) {
			tempCanvas.width = Math.min(canvas.width, image.width);
			tempCanvas.height = (image.height / image.width) * tempCanvas.width;
		} else {
			tempCanvas.height = Math.min(canvas.height, image.height);
			tempCanvas.width = (image.width / image.height) * tempCanvas.height;
		}
		tempCtx.filter = 'grayscale(100%)';
		tempCtx.globalAlpha = opacity;
		tempCtx.drawImage(image, 0, 0, tempCanvas.width, tempCanvas.height);

		return tempCanvas;
	};

	const upscaleCanvas = () => {
		const dpr = window.devicePixelRatio || 1;
		const rect = canvas.getBoundingClientRect();
		canvas.width = rect.width * dpr;
		canvas.height = rect.height * dpr;
		ctx.scale(dpr, dpr);
	};

	const draw = async () => {
		ctx.clearRect(0, 0, canvas.width, canvas.height);
		ctx.globalCompositeOperation = 'multiply';

		if (woodImage) {
			ctx.drawImage(woodImage, 0, 0, canvas.width, canvas.height);
		}

		if (canvasImage) {
			ctx.drawImage(canvasImage, transform.x, transform.y, transform.width, transform.height);
		}

		requestAnimationFrame(draw);
	};

	onMount(() => {
		if (!canvas) return;
		const newCtx = canvas.getContext('2d');
		if (!newCtx) return;
		ctx = newCtx;

		upscaleCanvas();
		requestAnimationFrame(draw);
		const wood = new Image();
		wood.src = '/designer-wood.jpg';
		wood.onload = () => {
			woodImage = wood;
		};
	});

	$effect(() => {
		if (!imageFile) return;

		const reader = new FileReader();
		reader.onload = (e) => {
			const image = new Image();
			image.src = e.target?.result as string;
			image.onload = () => {
				const aspectRatio = image.width / image.height;
				const height = canvas.height / 2;
				transform.height = height;
				const width = height * aspectRatio;
				transform.width = width;

				transform.x = (canvas.width - width) / 2;
				transform.y = (canvas.height - height) / 2;

				canvasImage = grayscale(image, 0.9);
			};
		};
		reader.readAsDataURL(imageFile);
	});
</script>

<div class="mx-auto flex w-full max-w-xl flex-col gap-4 p-4">
	<Input
		type="file"
		accept="image/*"
		oninput={(e) => {
			const file = e.currentTarget.files?.[0];
			if (!file) return;
			imageFile = file;
		}}
	/>
	<canvas
		bind:this={canvas}
		class:cursor-grab={canvasImage}
		class:cursor-grabbing={grabbing}
		onmousedown={() => {
			if (!canvasImage) return;
			grabbing = true;
		}}
		onmouseup={() => (grabbing = false)}
		onmouseleave={() => (grabbing = false)}
		onmousemove={(e) => {
			if (!grabbing) return;
			transform.x += e.movementX;
			transform.y += e.movementY;
		}}
		class="border-border rounded-md border shadow-sm"
		style="width: 100%; height: 400px;"
	></canvas>
</div>
