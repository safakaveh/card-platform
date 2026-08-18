<script lang="ts">
 import * as NavigationMenu from "$lib/components/ui/navigation-menu/index.js";
 import { navigationMenuTriggerStyle } from "$lib/components/ui/navigation-menu/navigation-menu-trigger.svelte";
 import ClipboardPenIcon from "@lucide/svelte/icons/clipboard-pen";
 import HouseIcon from '@lucide/svelte/icons/house';
 import LogoutIcon from '@lucide/svelte/icons/log-out';
 import UploadIcon from '@lucide/svelte/icons/cloud-upload';
 
 import { IsMobile } from "$lib/hooks/is-mobile.svelte.js";
 
 const isMobile = new IsMobile();
 let isShuttingDown = $state(false);
 let shutdownError = $state('');
 let shutdownRequested = $state(false);

 async function requestShutdown() {
  if (isShuttingDown) return;
  if (!window.confirm('آیا از خروج از برنامه مطمئن هستید؟')) return;

  isShuttingDown = true;
  shutdownError = '';
  try {
   const response = await fetch('/system/shutdown', {
    method: 'POST',
    headers: { Accept: 'application/json' }
   });
   if (!response.ok) throw new Error('خروج از برنامه ناموفق بود.');
   shutdownRequested = true;
  } catch (error) {
   shutdownError = error instanceof Error ? error.message : 'خروج از برنامه ناموفق بود.';
  } finally {
   isShuttingDown = false;
  }
 }
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
 
      <button
       type="button"
       onclick={requestShutdown}
       disabled={isShuttingDown}
       class="flex w-full items-center gap-2 rounded-md px-2 py-2 text-right text-sm transition hover:bg-red-50 hover:text-red-700 disabled:cursor-not-allowed disabled:opacity-50"
      >
       <LogoutIcon />
       {isShuttingDown ? 'در حال خروج...' : 'خروج از برنامه'}
      </button>
     </li>
    </ul>
   </NavigationMenu.Content>
  </NavigationMenu.Item>

  <NavigationMenu.Item class="block">
   <NavigationMenu.Trigger>گزارش ها</NavigationMenu.Trigger>
 
   <NavigationMenu.Content>
    <ul class="grid w-[250px] gap-4 p-2">
     <li>
       <NavigationMenu.Link href="/reports" class="flex-row items-center gap-2">
       <UploadIcon />
       بر اساس فایل ها
      </NavigationMenu.Link>
 
       <NavigationMenu.Link href="/reports" class="flex-row items-center gap-2">
       <ClipboardPenIcon />
       بر اساس فعالیت ها
      </NavigationMenu.Link>
 
     </li>
    </ul>
   </NavigationMenu.Content>
  </NavigationMenu.Item>

</NavigationMenu.List>
</NavigationMenu.Root>

{#if shutdownRequested || shutdownError}
 <div class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/30 p-4 backdrop-blur-sm" dir="rtl">
  <div class="w-full max-w-md rounded-3xl border border-white/80 bg-white p-6 text-right shadow-2xl">
   {#if shutdownRequested}
    <h2 class="text-xl font-black text-slate-900">برنامه با موفقیت بسته شد</h2>
    <p class="mt-3 text-sm leading-7 text-slate-600">
     سرور برنامه خاموش شد. برای استفاده‌ی دوباره، برنامه را مجدداً اجرا کنید.
    </p>
    <button
     type="button"
     class="mt-5 w-full rounded-xl bg-slate-900 px-4 py-3 text-sm font-bold text-white hover:bg-sky-700"
     onclick={() => (shutdownRequested = false)}
    >متوجه شدم</button>
   {:else}
    <h2 class="text-xl font-black text-red-700">خروج انجام نشد</h2>
    <p class="mt-3 text-sm leading-7 text-slate-600">{shutdownError}</p>
    <button
     type="button"
     class="mt-5 w-full rounded-xl bg-slate-900 px-4 py-3 text-sm font-bold text-white hover:bg-sky-700"
     onclick={() => (shutdownError = '')}
    >بستن</button>
   {/if}
  </div>
 </div>
{/if}
