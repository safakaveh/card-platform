<script lang="ts">
	import { NavigationMenu as NavigationMenuPrimitive } from "bits-ui";
	import { cn } from "$lib/utils.js";

	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: NavigationMenuPrimitive.ContentProps = $props();
</script>

<NavigationMenuPrimitive.Content
	bind:ref
	data-slot="navigation-menu-content"
	class={cn(
		"data-[motion^=from-]:animate-in data-[motion^=to-]:animate-out data-[motion^=from-]:fade-in data-[motion^=to-]:fade-out",
		
		// اصلاح انیمیشن‌ها برای پشتیبانی از هر دو جهت (LTR/RTL)
		// استفاده از کلاس‌های منطقی Tailwind (اگر نسخه 3.3+ هستید) یا هندل کردن دستی:
		"ltr:data-[motion=from-end]:slide-in-from-right-52 ltr:data-[motion=from-start]:slide-in-from-left-52",
		"rtl:data-[motion=from-end]:slide-in-from-left-52 rtl:data-[motion=from-start]:slide-in-from-right-52",
		"ltr:data-[motion=to-end]:slide-out-to-right-52 ltr:data-[motion=to-start]:slide-out-to-left-52",
		"rtl:data-[motion=to-end]:slide-out-to-left-52 rtl:data-[motion=to-start]:slide-out-to-right-52",

		"group-data-[viewport=false]/navigation-menu:bg-popover group-data-[viewport=false]/navigation-menu:text-popover-foreground group-data-[viewport=false]/navigation-menu:data-open:animate-in group-data-[viewport=false]/navigation-menu:data-closed:animate-out group-data-[viewport=false]/navigation-menu:data-closed:zoom-out-95 group-data-[viewport=false]/navigation-menu:data-open:zoom-in-95 group-data-[viewport=false]/navigation-menu:data-open:fade-in-0 group-data-[viewport=false]/navigation-menu:data-closed:fade-out-0 group-data-[viewport=false]/navigation-menu:ring-foreground/10",
		
		// اصلاح Padding: به جای pr-2.5 از ps یا pe استفاده کنید
		"p-2 ps-2.5 ease-[cubic-bezier(0.22,1,0.36,1)]", 
		
		"group-data-[viewport=false]/navigation-menu:rounded-md group-data-[viewport=false]/navigation-menu:shadow group-data-[viewport=false]/navigation-menu:ring-1 group-data-[viewport=false]/navigation-menu:duration-300 top-0",
		
		// اصلاح موقعیت: inset-inline-start-0 به جای left-0
		"start-0 w-full group-data-[viewport=false]/navigation-menu:top-full group-data-[viewport=false]/navigation-menu:mt-1.5 group-data-[viewport=false]/navigation-menu:overflow-hidden **:data-[slot=navigation-menu-link]:focus:ring-0 **:data-[slot=navigation-menu-link]:focus:outline-none md:absolute md:w-auto",
		
		className
	)}
	{...restProps}
/>
