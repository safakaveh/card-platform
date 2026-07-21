<script lang="ts">
	import { NavigationMenu as NavigationMenuPrimitive } from "bits-ui";
	import { cn } from "$lib/utils.js";

	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: NavigationMenuPrimitive.ListProps = $props();
</script>

<NavigationMenuPrimitive.List
	bind:ref
	data-slot="navigation-menu-list"
	class={cn(
		// استفاده از justify-start باعث می‌شود در RTL راست‌چین و در LTR چپ‌چین شود.
		// اگر اصرار دارید منو همیشه وسط باشد، همان justify-center را نگه دارید.
		"gap-0 group flex flex-1 list-none items-center justify-start",
		className
	)}
	{...restProps}
/>
