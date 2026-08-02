<script lang="ts">
	let { children } = $props();
</script>

<svelte:head>
	<title>ویرایش فایل</title>
</svelte:head>

<div class="w-full">
	{@render children()}
</div>
