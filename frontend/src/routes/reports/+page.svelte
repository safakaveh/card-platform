<script lang="ts">
 import { onMount } from 'svelte';
 import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
 import ClipboardListIcon from '@lucide/svelte/icons/clipboard-list';
 import TriangleAlertIcon from '@lucide/svelte/icons/triangle-alert';
 import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';

 type PendingCard = {
  card_uuid: string;
  order_uuid: string;
  order_name: string;
  block_no: number;
  created_at: number;
 };
 type PendingResponse = { count: number; items: PendingCard[] };

 let data = $state<PendingResponse>({ count: 0, items: [] });
 let loading = $state(true);
 let error = $state('');

 const formatter = new Intl.DateTimeFormat('fa-IR', { dateStyle: 'medium', timeStyle: 'short' });
 function date(value: number) {
  const d = new Date(value < 10_000_000_000 ? value * 1000 : value);
  return Number.isNaN(d.getTime()) ? '—' : formatter.format(d);
 }

 async function load() {
  loading = true;
  error = '';
  try {
   const response = await fetch('/api/data/pending?limit=100', { headers: { Accept: 'application/json' } });
   const body = await response.json();
   if (!response.ok) throw new Error(body?.error || 'دریافت گزارش ناموفق بود.');
   data = body as PendingResponse;
  } catch (reason) {
   error = reason instanceof Error ? reason.message : 'دریافت گزارش ناموفق بود.';
  } finally {
   loading = false;
  }
 }
 onMount(load);
</script>

<section dir="rtl" class="mx-auto w-full max-w-6xl px-4 py-6 sm:px-6 lg:px-8">
 <div class="rounded-[2rem] bg-gradient-to-br from-indigo-50 via-white to-sky-50 p-4 sm:p-6">
  <div class="mb-5 flex flex-col gap-4 rounded-3xl border border-white/80 bg-white/85 p-5 shadow-sm sm:flex-row sm:items-end sm:justify-between">
   <div>
    <p class="mb-2 text-sm font-semibold text-indigo-700">گزارش‌ها</p>
    <h1 class="text-2xl font-black text-slate-900 sm:text-3xl">کارت‌های در انتظار تخصیص UID</h1>
    <p class="mt-2 text-sm leading-6 text-slate-600">کارت‌هایی که در فایل CSV برای آن‌ها UID درخواست شده و هنوز مقدار واقعی دریافت نکرده‌اند.</p>
   </div>
   <button type="button" onclick={load} disabled={loading} class="inline-flex items-center justify-center gap-2 rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm font-bold text-slate-700 hover:border-indigo-300 hover:bg-indigo-50 disabled:opacity-50">
    <RefreshCwIcon class={['size-4', loading && 'animate-spin']} /> تازه‌سازی
   </button>
  </div>

  {#if error}
   <div class="mb-5 flex items-center gap-3 rounded-2xl border border-red-200 bg-red-50 p-4 text-sm text-red-800"><TriangleAlertIcon class="size-5" />{error}</div>
  {/if}

  <div class="mb-5 rounded-2xl border border-indigo-100 bg-white p-5 shadow-sm">
   <p class="text-xs font-bold text-slate-500">تعداد کارت‌های در انتظار</p>
   <p class="mt-2 text-3xl font-black text-indigo-700">{data.count.toLocaleString('fa-IR')}</p>
  </div>

  <div class="overflow-hidden rounded-3xl border border-slate-200 bg-white shadow-xl">
   {#if loading}
    <div class="flex items-center justify-center gap-3 p-14 text-sm text-slate-500"><LoaderCircleIcon class="size-6 animate-spin text-indigo-600" />در حال دریافت گزارش...</div>
   {:else if data.items.length === 0}
    <div class="flex flex-col items-center p-14 text-center"><ClipboardListIcon class="mb-4 size-10 text-slate-400" /><strong class="text-slate-800">موردی در انتظار نیست</strong><p class="mt-2 text-sm text-slate-500">همه‌ی UIDهای ثبت‌شده آماده هستند.</p></div>
   {:else}
    <div class="overflow-x-auto">
     <table class="w-full min-w-[760px] text-right text-sm">
      <thead class="bg-slate-50 text-slate-600"><tr class="border-b border-slate-100"><th class="px-5 py-4 font-bold">نام سفارش</th><th class="px-5 py-4 font-bold">شناسه کارت</th><th class="px-5 py-4 font-bold">بلوک</th><th class="px-5 py-4 font-bold">تاریخ</th></tr></thead>
      <tbody class="divide-y divide-slate-100">{#each data.items as item (item.card_uuid)}<tr class="hover:bg-indigo-50/40"><td class="px-5 py-4 font-bold text-slate-900">{item.order_name || 'بدون نام'}</td><td class="px-5 py-4 font-mono text-xs text-slate-500">{item.card_uuid}</td><td class="px-5 py-4 text-slate-700">{item.block_no.toLocaleString('fa-IR')}</td><td class="px-5 py-4 text-slate-600">{date(item.created_at)}</td></tr>{/each}</tbody>
     </table>
    </div>
   {/if}
  </div>
 </div>
</section>
