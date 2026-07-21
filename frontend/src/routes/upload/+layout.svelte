<script lang="ts">
	let { children } = $props();
</script>

<svelte:head>
	<title>بارگذاری فایل</title>
</svelte:head>

<div class="flex min-h-screen items-center justify-center">
	{@render children()}
</div>
