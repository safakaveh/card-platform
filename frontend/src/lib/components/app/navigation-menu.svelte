<script lang="ts">
 import * as NavigationMenu from "$lib/components/ui/navigation-menu/index.js";
 import { navigationMenuTriggerStyle } from "$lib/components/ui/navigation-menu/navigation-menu-trigger.svelte";
 import ClipboardPenIcon from "@lucide/svelte/icons/clipboard-pen";
 import HouseIcon from '@lucide/svelte/icons/house';
 import LogoutIcon from '@lucide/svelte/icons/log-out';
 import UploadIcon from '@lucide/svelte/icons/cloud-upload';
 
 import { IsMobile } from "$lib/hooks/is-mobile.svelte.js";
 
 const isMobile = new IsMobile();
</script>
 
<NavigationMenu.Root viewport={isMobile.current}>
 <NavigationMenu.List class="flex-wrap">
  
  <NavigationMenu.Item>
   <NavigationMenu.Link>
    {#snippet child()}
     <a href="/" class={navigationMenuTriggerStyle()}><HouseIcon class="me-1"/>خانه</a>
    {/snippet}
   </NavigationMenu.Link>
  </NavigationMenu.Item>

  <NavigationMenu.Item class="block">
   <NavigationMenu.Trigger>فایل</NavigationMenu.Trigger>
 
   <NavigationMenu.Content>
    <ul class="grid w-[250px] gap-4 p-2">
     <li>
      <NavigationMenu.Link href="/upload" class="flex-row items-center gap-2">
       <UploadIcon />
       بارگذاری فایل CSV
      </NavigationMenu.Link>
 
      <NavigationMenu.Link href="/editfile" class="flex-row items-center gap-2">
       <ClipboardPenIcon />
       ویرایش فایل های بارگذاری شده
      </NavigationMenu.Link>
 
      <NavigationMenu.Link href="##" class="flex-row items-center gap-2">
       <LogoutIcon />
       خروج از برنامه
      </NavigationMenu.Link>
     </li>
    </ul>
   </NavigationMenu.Content>
  </NavigationMenu.Item>

  <NavigationMenu.Item class="block">
   <NavigationMenu.Trigger>گزارش ها</NavigationMenu.Trigger>
 
   <NavigationMenu.Content>
    <ul class="grid w-[250px] gap-4 p-2">
     <li>
      <NavigationMenu.Link href="##" class="flex-row items-center gap-2">
       <UploadIcon />
       بر اساس فایل ها
      </NavigationMenu.Link>
 
      <NavigationMenu.Link href="##" class="flex-row items-center gap-2">
       <ClipboardPenIcon />
       بر اساس فعالیت ها
      </NavigationMenu.Link>
 
     </li>
    </ul>
   </NavigationMenu.Content>
  </NavigationMenu.Item>

</NavigationMenu.List>
</NavigationMenu.Root>
