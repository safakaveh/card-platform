<script lang="ts">
	import { onMount } from 'svelte';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';

	type ImportSummary = { card_count?: number; status?: string };
	let ordersCount = $state<number | null>(null);
	let cardsCount = $state<number | null>(null);
	let pendingCount = $state<number | null>(null);

	async function loadDashboard() {
		try {
			const [importsResponse, pendingResponse] = await Promise.all([
				fetch('/api/imports/?limit=100', { headers: { Accept: 'application/json' } }),
				fetch('/api/data/pending?limit=1000', { headers: { Accept: 'application/json' } })
			]);
			if (importsResponse.ok) {
				const imports = (await importsResponse.json()) as ImportSummary[];
				ordersCount = imports.length;
				cardsCount = imports.reduce((total, item) => total + (item.card_count || 0), 0);
			}
			if (pendingResponse.ok) {
				const pending = (await pendingResponse.json()) as { count?: number };
				pendingCount = pending.count || 0;
			}
		} catch {
			// The guide remains useful when the API is temporarily unavailable.
		}
	}

	onMount(loadDashboard);
</script>

<svelte:head>
	<title>راهنمای استفاده از دستگاه XYZ</title>
</svelte:head>

<div dir="rtl" class="flex min-h-full items-center justify-center px-4 py-8">
	<Card.Root class="mx-auto w-full max-w-3xl overflow-hidden shadow-md">
		<div class="grid gap-3 border-b bg-slate-50/70 p-4 sm:grid-cols-3">
			{#each [
				{ label: 'تعداد سفارش‌ها', value: ordersCount },
				{ label: 'مجموع کارت‌ها', value: cardsCount },
				{ label: 'UIDهای در انتظار', value: pendingCount }
			] as metric}
				<div class="rounded-xl border bg-white p-3 text-center shadow-sm">
					<p class="text-xs font-bold text-muted-foreground">{metric.label}</p>
					<p class="mt-1 text-2xl font-black text-primary">
						{metric.value === null ? '—' : metric.value.toLocaleString('fa-IR')}
					</p>
				</div>
			{/each}
		</div>
		<Card.Header class="space-y-3 bg-muted/30 text-right">
			<div class="flex items-center gap-3">
				<div
					class="flex h-11 w-11 items-center justify-center rounded-xl bg-primary text-lg font-bold text-primary-foreground"
				>
					XYZ
				</div>

				<div>
					<Card.Title class="text-xl font-bold sm:text-2xl">
						به برنامه دستگاه XYZ خوش آمدید
					</Card.Title>

					<Card.Description class="mt-1 text-sm sm:text-base">
						پیاده‌سازی و طراحی توسط شرکت تراشه هوشمند
					</Card.Description>
				</div>
			</div>
		</Card.Header>

		<Card.Content class="space-y-6 pt-6 text-right">
			<div class="rounded-lg border border-blue-200 bg-blue-50 p-4 text-blue-950 dark:border-blue-900 dark:bg-blue-950/30 dark:text-blue-100">
				<p class="font-medium">
					لطفاً پیش از استفاده از برنامه، نکات زیر را با دقت مطالعه فرمایید.
				</p>
			</div>

			<ol class="space-y-4">
				<li class="flex gap-3">
					<span
						class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-primary text-sm font-bold text-primary-foreground"
					>
						۱
					</span>

					<p class="leading-7">
						فایل بارگذاری‌شده باید یک فایل
						<Badge variant="secondary">CSV</Badge>
						استاندارد و دارای سطر عنوان ستون‌ها
						<code class="rounded bg-muted px-1.5 py-0.5 text-sm">header</code>
						باشد.
					</p>
				</li>

				<li class="flex gap-3">
					<span
						class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-primary text-sm font-bold text-primary-foreground"
					>
						۲
					</span>

					<p class="leading-7">
						اگر فایل شما با فرمت Excel است، کافی است آن را با فرمت
						<Badge variant="secondary">.csv</Badge>
						ذخیره کنید.
					</p>
				</li>

				<li class="flex gap-3">
					<span
						class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-primary text-sm font-bold text-primary-foreground"
					>
						۳
					</span>

					<p class="leading-7">
						ستون‌هایی که عنوان آن‌ها با
						<code class="rounded bg-emerald-100 px-1.5 py-0.5 text-sm text-emerald-800 dark:bg-emerald-950 dark:text-emerald-200">
							frn_
						</code>
						شروع می‌شود، روی قسمت جلوی کارت لیزر خواهند شد.
					</p>
				</li>

				<li class="flex gap-3">
					<span
						class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-primary text-sm font-bold text-primary-foreground"
					>
						۴
					</span>

					<p class="leading-7">
						ستون‌هایی که عنوان آن‌ها با
						<code class="rounded bg-amber-100 px-1.5 py-0.5 text-sm text-amber-800 dark:bg-amber-950 dark:text-amber-200">
							bck_
						</code>
						شروع می‌شود، روی قسمت پشت کارت لیزر خواهند شد.
					</p>
				</li>

				<li class="flex gap-3">
					<span
						class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-primary text-sm font-bold text-primary-foreground"
					>
						۵
					</span>

					<p class="leading-7">
						برای لیزر کردن UID کارت‌ها، کافی است یک ستون خالی با یکی از عنوان‌های زیر
						در فایل CSV ایجاد کنید:
					</p>
				</li>
			</ol>

			<div class="grid gap-3 sm:grid-cols-2">
				<div class="rounded-lg border bg-muted/40 p-3 text-center">
					<code class="text-sm font-semibold">frn_uid</code>
					<p class="mt-1 text-xs text-muted-foreground">
						لیزر UID در قسمت جلوی کارت
					</p>
				</div>

				<div class="rounded-lg border bg-muted/40 p-3 text-center">
					<code class="text-sm font-semibold">bck_uid</code>
					<p class="mt-1 text-xs text-muted-foreground">
						لیزر UID در قسمت پشت کارت
					</p>
				</div>
			</div>

			<div class="rounded-lg border border-red-200 bg-red-50 p-4 text-sm leading-7 text-red-900 dark:border-red-900 dark:bg-red-950/30 dark:text-red-100">
				<strong>نکته:</strong>
				مطمئن شوید نام ستون‌ها دقیقاً مطابق الگوی اعلام‌شده باشد؛ حروف اضافه یا فاصله
				ممکن است باعث شناسایی نشدن ستون شود.
			</div>
		</Card.Content>

		<Separator />

		<Card.Footer class="justify-center px-6 py-5 text-center">
			<p class="text-sm leading-7 text-muted-foreground">
				از توجه شما سپاسگزاریم و امیدواریم از استفاده از دستگاه ما رضایت کامل داشته باشید.
			</p>
		</Card.Footer>
	</Card.Root>
</div>
