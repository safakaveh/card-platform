<script lang="ts">
 import { onMount } from 'svelte';
 import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
 import ClipboardListIcon from '@lucide/svelte/icons/clipboard-list';
 import TriangleAlertIcon from '@lucide/svelte/icons/triangle-alert';
 import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';
 import PlayIcon from '@lucide/svelte/icons/play';
 import Undo2Icon from '@lucide/svelte/icons/undo-2';

 type Laser = { side: string; row: number; content_type: string; content: string };
 type ReadItem = {
  card_uuid: string; order_uuid: string; order_name: string; sequence: number;
  read_at?: number; laser: Laser[]; magnet: unknown[]; mifare: unknown[];
 };
 type Report = { count: number; items: (ReadItem & { created_at: number; read: boolean })[] };

 let orderName = $state('');
 let filter = $state('');
 let report = $state<Report>({ count: 0, items: [] });
 let current = $state<ReadItem[]>([]);
 let loading = $state(false);
 let reading = $state(false);
 let error = $state('');
 const formatter = new Intl.DateTimeFormat('fa-IR', { dateStyle: 'medium', timeStyle: 'short' });
 const date = (value: number) => {
  const d = new Date(value < 10_000_000_000 ? value * 1000 : value);
  return Number.isNaN(d.getTime()) ? '—' : formatter.format(d);
 };

 async function loadReport() {
  loading = true; error = '';
  try {
   const params = new URLSearchParams({ limit: '200' });
   if (filter.trim()) params.set('order_name', filter.trim());
   const response = await fetch(`/api/data/read-report?${params}`);
   const body = await response.json();
   if (!response.ok) throw new Error(body?.error || 'دریافت گزارش ناموفق بود.');
   report = body;
  } catch (e) { error = e instanceof Error ? e.message : 'دریافت گزارش ناموفق بود.'; }
  finally { loading = false; }
 }

 async function readNext() {
  if (!orderName.trim()) { error = 'نام سفارش را وارد کنید.'; return; }
  reading = true; error = ''; current = [];
  try {
   const params = new URLSearchParams({ order_name: orderName.trim(), limit: '1' });
   const response = await fetch(`/api/data/read?${params}`);
   const body = await response.json();
   if (!response.ok) throw new Error(body?.error || 'خواندن داده ناموفق بود.');
   current = body.items || [];
   await loadReport();
   if (current.length === 0) error = 'برای این سفارش دادهٔ خوانده‌نشده‌ای باقی نمانده است.';
  } catch (e) { error = e instanceof Error ? e.message : 'خواندن داده ناموفق بود.'; }
  finally { reading = false; }
 }

 async function reset(cardUUID: string) {
  try {
   const response = await fetch(`/api/data/read/${encodeURIComponent(cardUUID)}/reset`, { method: 'POST' });
   if (!response.ok) { const body = await response.json(); throw new Error(body?.error || 'بازنشانی ناموفق بود.'); }
   await loadReport();
  } catch (e) { error = e instanceof Error ? e.message : 'بازنشانی ناموفق بود.'; }
 }
 onMount(loadReport);
</script>

<section dir="rtl" class="mx-auto w-full max-w-6xl px-4 py-6 sm:px-6 lg:px-8">
 <div class="rounded-[2rem] bg-gradient-to-br from-indigo-50 via-white to-sky-50 p-4 sm:p-6">
  <div class="mb-5 rounded-3xl border border-white/80 bg-white/85 p-5 shadow-sm">
   <p class="mb-2 text-sm font-semibold text-indigo-700">گزارش و خواندن ترتیبی</p>
   <h1 class="text-2xl font-black text-slate-900 sm:text-3xl">داده‌های خوانده‌شده</h1>
   <p class="mt-2 text-sm leading-6 text-slate-600">هر بار فقط اولین ردیف خوانده‌نشدهٔ سفارش برگردانده و همان لحظه خوانده‌شده می‌شود.</p>
   <div class="mt-5 flex flex-col gap-3 sm:flex-row">
    <input bind:value={orderName} placeholder="نام سفارش" class="flex-1 rounded-xl border border-slate-200 bg-white px-4 py-3 text-sm outline-none focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100" />
    <button type="button" onclick={readNext} disabled={reading} class="inline-flex items-center justify-center gap-2 rounded-xl bg-indigo-600 px-5 py-3 text-sm font-bold text-white hover:bg-indigo-700 disabled:opacity-50"><PlayIcon class="size-4" />{reading ? 'در حال خواندن...' : 'خواندن ردیف بعدی'}</button>
   </div>
  </div>

  {#if error}<div class="mb-5 flex items-center gap-3 rounded-2xl border border-red-200 bg-red-50 p-4 text-sm text-red-800"><TriangleAlertIcon class="size-5" />{error}</div>{/if}

  {#if current.length}
   <div class="mb-5 rounded-3xl border border-emerald-200 bg-emerald-50 p-5">
    <h2 class="font-black text-emerald-900">ردیف فعلی</h2>
    {#each current as item}
     <p class="mt-2 text-sm text-emerald-800">کارت <span class="font-mono">{item.card_uuid}</span> با {item.laser.length} مقدار لیزر خوانده شد.</p>
     <div class="mt-3 grid gap-2 sm:grid-cols-2">{#each item.laser as value}<div class="rounded-lg bg-white/80 p-3 text-sm"><span class="font-bold">{value.side} / {value.content_type}:</span> {value.content}</div>{/each}</div>
    {/each}
   </div>
  {/if}

  <div class="mb-5 flex items-center justify-between rounded-2xl border border-indigo-100 bg-white p-5 shadow-sm">
   <div><p class="text-xs font-bold text-slate-500">تعداد ردیف‌های گزارش</p><p class="mt-2 text-3xl font-black text-indigo-700">{report.count.toLocaleString('fa-IR')}</p></div>
   <div class="flex gap-2"><input bind:value={filter} placeholder="فیلتر نام سفارش" class="rounded-xl border border-slate-200 px-3 py-2 text-sm" /><button type="button" onclick={loadReport} disabled={loading} class="rounded-xl border border-slate-200 p-3 hover:bg-indigo-50"><RefreshCwIcon class={['size-4', loading && 'animate-spin']} /></button></div>
  </div>

  <div class="overflow-hidden rounded-3xl border border-slate-200 bg-white shadow-xl">
   {#if loading}<div class="flex items-center justify-center gap-3 p-14 text-sm text-slate-500"><LoaderCircleIcon class="size-6 animate-spin text-indigo-600" />در حال دریافت گزارش...</div>
   {:else if report.items.length === 0}<div class="flex flex-col items-center p-14 text-center"><ClipboardListIcon class="mb-4 size-10 text-slate-400" /><strong class="text-slate-800">موردی پیدا نشد</strong></div>
   {:else}<div class="overflow-x-auto"><table class="w-full min-w-[900px] text-right text-sm"><thead class="bg-slate-50 text-slate-600"><tr class="border-b border-slate-100"><th class="px-5 py-4 font-bold">سفارش</th><th class="px-5 py-4 font-bold">شناسه کارت</th><th class="px-5 py-4 font-bold">وضعیت</th><th class="px-5 py-4 font-bold">تاریخ خواندن</th><th class="px-5 py-4"></th></tr></thead><tbody class="divide-y divide-slate-100">{#each report.items as item (item.card_uuid)}<tr class="hover:bg-indigo-50/40"><td class="px-5 py-4 font-bold">{item.order_name}</td><td class="px-5 py-4 font-mono text-xs text-slate-500">{item.card_uuid}</td><td class="px-5 py-4">{item.read ? 'خوانده شده' : 'خوانده نشده'}</td><td class="px-5 py-4 text-slate-600">{item.read_at ? date(item.read_at) : '—'}</td><td class="px-5 py-4 text-left">{#if item.read}<button type="button" onclick={() => reset(item.card_uuid)} class="inline-flex items-center gap-1 rounded-lg border border-amber-200 px-3 py-2 text-xs font-bold text-amber-700 hover:bg-amber-50"><Undo2Icon class="size-4" />خوانده‌نشده</button>{/if}</td></tr>{/each}</tbody></table></div>{/if}
  </div>
 </div>
</section>
