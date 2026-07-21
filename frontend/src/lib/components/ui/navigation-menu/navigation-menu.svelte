<script lang="ts">
	import { NavigationMenu as NavigationMenuPrimitive } from "bits-ui";
	import { cn } from "$lib/utils.js";
	import NavigationMenuViewport from "./navigation-menu-viewport.svelte";

	let {
		ref = $bindable(null),
		class: className,
		viewport = true,
		children,
		dir = "rtl", // اضافه کردن dir به عنوان مقدار پیش‌فرض فارسی یا دریافت پویا
		...restProps
	}: NavigationMenuPrimitive.RootProps & {
		viewport?: boolean;
	} = $props();
</script>

<NavigationMenuPrimitive.Root
	bind:ref
	{dir}
	data-slot="navigation-menu"
	data-viewport={viewport}
	class={cn(
		// تغییر justify-center به justify-start برای تراز درست در راست/چپ
		// اگر اصرار داری منو همیشه و در هر حالتی دقیقاً وسط صفحه باشد، justify-center را نگه دار.
		"group/navigation-menu relative flex max-w-max flex-1 items-center justify-start",
		className
	)}
	{...restProps}
>
	{@render children?.()}
	{#if viewport}
		<NavigationMenuViewport />
	{/if}
</NavigationMenuPrimitive.Root>
