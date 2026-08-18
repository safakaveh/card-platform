<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import FileSpreadsheetIcon from '@lucide/svelte/icons/file-spreadsheet';
	import Trash2Icon from '@lucide/svelte/icons/trash-2';
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import TriangleAlertIcon from '@lucide/svelte/icons/triangle-alert';
	import CircleCheckIcon from '@lucide/svelte/icons/circle-check';
	import Clock3Icon from '@lucide/svelte/icons/clock-3';
	import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';
	import FolderOpenIcon from '@lucide/svelte/icons/folder-open';

	type ImportStatus = 'pending' | 'processing' | 'completed' | 'failed' | string;

	type ImportSummary = {
		uuid: string;
		order_name: string;
		status: ImportStatus;
		file_name: string;
		card_count: number;
		created_at: number;
	};

	let imports = $state<ImportSummary[]>([]);
	let isLoading = $state(true);
	let deletingID = $state('');
	let errorMessage = $state('');

	let loadController: AbortController | null = null;

	const importsCount = $derived(imports.length);
	const totalCards = $derived(imports.reduce((sum, item) => sum + (item.card_count || 0), 0));

	const dateFormatter = new Intl.DateTimeFormat('fa-IR', {
		dateStyle: 'medium',
		timeStyle: 'short'
	});

	function getErrorMessage(data: unknown, fallback: string): string {
		if (data && typeof data === 'object' && 'error' in data) {
			const error = (data as { error?: unknown }).error;
			if (typeof error === 'string' && error.trim()) return error;
		}

		if (typeof data === 'string' && data.trim()) {
			return data;
		}

		return fallback;
	}

	async function getResponseData(response: Response): Promise<unknown> {
		const contentType = response.headers.get('content-type') ?? '';

		if (contentType.includes('application/json')) {
			try {
				return await response.json();
			} catch {
				return null;
			}
		}

		try {
			const text = await response.text();
			return text.trim() || null;
		} catch {
			return null;
		}
	}

	function formatDate(timestamp: number): string {
		if (!timestamp || !Number.isFinite(timestamp)) return '—';

		const milliseconds = timestamp < 10_000_000_000 ? timestamp * 1000 : timestamp;
		const date = new Date(milliseconds);

		if (Number.isNaN(date.getTime())) return '—';

		return dateFormatter.format(date);
	}

	function getStatusLabel(status: ImportStatus): string {
		switch (status.toLowerCase()) {
			case 'completed':
			case 'complete':
			case 'success':
			case 'done':
				return 'تکمیل‌شده';
			case 'processing':
			case 'in_progress':
				return 'در حال پردازش';
			case 'pending':
			case 'queued':
				return 'در انتظار';
			case 'failed':
			case 'error':
				return 'ناموفق';
			default:
				return status || 'نامشخص';
		}
	}

	function getStatusClass(status: ImportStatus): string {
		switch (status.toLowerCase()) {
			case 'completed':
			case 'complete':
			case 'success':
			case 'done':
				return 'border-emerald-200 bg-emerald-50 text-emerald-700';
			case 'processing':
			case 'in_progress':
				return 'border-sky-200 bg-sky-50 text-sky-700';
			case 'pending':
			case 'queued':
				return 'border-amber-200 bg-amber-50 text-amber-700';
			case 'failed':
			case 'error':
				return 'border-red-200 bg-red-50 text-red-700';
			default:
				return 'border-slate-200 bg-slate-50 text-slate-600';
		}
	}

	async function loadImports() {
		loadController?.abort();

		const controller = new AbortController();
		loadController = controller;

		isLoading = true;
		errorMessage = '';

		try {
			const response = await fetch('/api/imports/?limit=100', {
				signal: controller.signal,
				headers: {
					Accept: 'application/json'
				}
			});

			const data = await getResponseData(response);

			if (!response.ok) {
				throw new Error(getErrorMessage(data, 'دریافت فهرست فایل‌ها ناموفق بود.'));
			}

			if (!Array.isArray(data)) {
				throw new Error('پاسخ دریافتی از سرور معتبر نیست.');
			}

			imports = data as ImportSummary[];
		} catch (error) {
			if (error instanceof DOMException && error.name === 'AbortError') return;

			errorMessage = error instanceof Error ? error.message : 'دریافت فهرست فایل‌ها ناموفق بود.';
		} finally {
			if (loadController === controller) {
				isLoading = false;
				loadController = null;
			}
		}
	}

	async function deleteImport(item: ImportSummary) {
		const confirmed = confirm(
			`آیا از حذف سفارش «${item.order_name}» و تمام کارت‌های مربوط به آن مطمئن هستید؟`
		);

		if (!confirmed) return;

		deletingID = item.uuid;
		errorMessage = '';

		try {
			const response = await fetch(`/api/imports/${encodeURIComponent(item.uuid)}`, {
				method: 'DELETE',
				headers: {
					Accept: 'application/json'
				}
			});

			const data = await getResponseData(response);

			if (!response.ok) {
				throw new Error(getErrorMessage(data, 'حذف سفارش ناموفق بود.'));
			}

			imports = imports.filter((entry) => entry.uuid !== item.uuid);
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'حذف سفارش ناموفق بود.';
		} finally {
			deletingID = '';
		}
	}

	onMount(() => {
		loadImports();
	});

	onDestroy(() => {
		loadController?.abort();
	});
</script>

<section dir="rtl" class="mx-auto w-full max-w-6xl px-4 py-6 sm:px-6 lg:px-8">
	<div class="rounded-[2rem] bg-gradient-to-br from-sky-50 via-white to-slate-100 p-4 sm:p-6">
		<div
			class="mb-5 flex flex-col gap-4 rounded-3xl border border-white/80 bg-white/80 p-5 shadow-sm backdrop-blur sm:flex-row sm:items-end sm:justify-between"
		>
			<div>
				<p class="mb-2 text-sm font-semibold text-sky-700">مدیریت داده‌ها</p>
				<h1 class="text-2xl font-black tracking-tight text-slate-900 sm:text-3xl">
					فایل‌های بارگذاری‌شده
				</h1>
				<p class="mt-2 text-sm leading-6 text-slate-600">
					فهرست سفارش‌های واردشده، وضعیت پردازش و امکان حذف داده‌های بارگذاری‌شده.
				</p>
			</div>

			<button
				type="button"
				onclick={loadImports}
				disabled={isLoading}
				aria-label="تازه‌سازی فهرست فایل‌های بارگذاری‌شده"
				class="inline-flex items-center justify-center gap-2 rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm font-bold text-slate-700 shadow-sm transition hover:border-sky-300 hover:bg-sky-50 focus:ring-2 focus:ring-sky-300 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50"
			>
				<RefreshCwIcon class={['size-4', isLoading && 'animate-spin']} />
				تازه‌سازی
			</button>
		</div>

		<div class="mb-5 grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
			<div class="rounded-2xl border border-slate-200 bg-white/90 p-4 shadow-sm">
				<p class="text-xs font-bold text-slate-500">تعداد سفارش‌ها</p>
				<p class="mt-2 text-2xl font-black text-slate-900">
					{importsCount.toLocaleString('fa-IR')}
				</p>
			</div>

			<div class="rounded-2xl border border-slate-200 bg-white/90 p-4 shadow-sm">
				<p class="text-xs font-bold text-slate-500">مجموع کارت‌ها</p>
				<p class="mt-2 text-2xl font-black text-slate-900">
					{totalCards.toLocaleString('fa-IR')}
				</p>
			</div>

			<div
				class="rounded-2xl border border-slate-200 bg-white/90 p-4 shadow-sm sm:col-span-2 lg:col-span-1"
			>
				<p class="text-xs font-bold text-slate-500">وضعیت نمایش</p>
				<p class="mt-2 text-sm font-bold text-slate-700">
					{isLoading
						? 'در حال دریافت اطلاعات'
						: importsCount > 0
							? 'فهرست آماده است'
							: 'داده‌ای موجود نیست'}
				</p>
			</div>
		</div>

		<!-- {#if errorMessage}
			<div
				class="mb-5 flex items-start justify-between gap-4 rounded-2xl border border-red-200 bg-red-50 p-4 text-sm leading-6 text-red-800"
				role="alert"
			>
				<div class="flex items-start gap-3">
					<TriangleAlertIcon class="mt-0.5 size-5 shrink-0" />
					<span>{errorMessage}</span>
				</div>

				<button
					type="button"
					onclick={loadImports}
					class="shrink-0 rounded-lg px-2 py-1 font-bold text-red-700 transition hover:bg-red-100"
				>
					تلاش مجدد
				</button>
			</div>
		{/if} -->

		<div
			class="overflow-hidden rounded-3xl border border-slate-200 bg-white shadow-xl shadow-slate-200/60"
		>
			<div class="border-b border-slate-100 bg-slate-50/80 px-5 py-4">
				<div class="flex flex-wrap items-center justify-between gap-3">
					<div class="flex items-center gap-2">
						<span class="rounded-xl bg-sky-100 p-2 text-sky-700">
							<FolderOpenIcon class="size-5" />
						</span>
						<div>
							<p class="text-sm font-black text-slate-800">فهرست فایل‌ها</p>
							<p class="text-xs text-slate-500">آخرین ۱۰۰ فایل بارگذاری‌شده</p>
						</div>
					</div>

					{#if !isLoading && imports.length > 0}
						<p class="text-xs font-medium text-slate-500">
							نمایش {imports.length.toLocaleString('fa-IR')} مورد
						</p>
					{/if}
				</div>
			</div>

			{#if isLoading}
				<div
					class="flex flex-col items-center justify-center gap-3 p-14 text-center text-sm text-slate-500"
				>
					<LoaderCircleIcon class="size-7 animate-spin text-sky-600" />
					<span>در حال دریافت اطلاعات...</span>
				</div>
			{:else if imports.length === 0}
				<div class="flex flex-col items-center p-14 text-center">
					<span class="mb-4 rounded-2xl bg-slate-100 p-4 text-slate-500">
						<FileSpreadsheetIcon class="size-9" />
					</span>
					<strong class="text-base text-slate-800">هنوز فایلی بارگذاری نشده است</strong>
					<p class="mt-2 max-w-md text-sm leading-6 text-slate-500">
						برای شروع، یک فایل CSV شامل اطلاعات کارت‌ها بارگذاری کنید تا در این بخش نمایش داده شود.
					</p>
					<a
						href="/upload"
						class="mt-5 rounded-2xl bg-slate-900 px-5 py-3 text-sm font-bold text-white transition hover:bg-sky-700 focus:ring-2 focus:ring-sky-300 focus:outline-none"
					>
						بارگذاری اولین فایل
					</a>
				</div>
			{:else}
				<div class="overflow-x-auto">
					<table class="w-full min-w-[900px] text-right text-sm">
						<thead class="bg-slate-50 text-slate-600">
							<tr class="border-b border-slate-100">
								<th class="px-5 py-4 font-bold">نام سفارش</th>
								<th class="px-5 py-4 font-bold">نام فایل</th>
								<th class="px-5 py-4 font-bold">وضعیت</th>
								<th class="px-5 py-4 font-bold">تعداد کارت</th>
								<th class="px-5 py-4 font-bold">تاریخ ورود</th>
								<th class="px-5 py-4 font-bold">عملیات</th>
							</tr>
						</thead>

						<tbody class="divide-y divide-slate-100">
							{#each imports as item (item.uuid)}
								<tr class="transition hover:bg-sky-50/40">
									<td class="px-5 py-4 align-middle">
										<div class="flex flex-col">
											<span class="font-bold text-slate-900">
												{item.order_name || 'بدون نام'}
											</span>
											<span class="mt-1 text-xs text-slate-400">
												شناسه: {item.uuid}
											</span>
										</div>
									</td>

									<td class="max-w-56 px-5 py-4 align-middle text-slate-600">
										<span class="block truncate" title={item.file_name}>
											{item.file_name || '—'}
										</span>
									</td>

									<td class="px-5 py-4 align-middle">
										<span
											class={[
												'inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1 text-xs font-bold',
												getStatusClass(item.status)
											]}
										>
											{#if ['completed', 'complete', 'success', 'done'].includes(item.status.toLowerCase())}
												<CircleCheckIcon class="size-3.5" />
											{:else if ['processing', 'in_progress', 'pending', 'queued'].includes(item.status.toLowerCase())}
												<Clock3Icon class="size-3.5" />
											{:else}
												<TriangleAlertIcon class="size-3.5" />
											{/if}
											{getStatusLabel(item.status)}
										</span>
									</td>

									<td class="px-5 py-4 align-middle font-medium text-slate-700">
										{item.card_count.toLocaleString('fa-IR')}
									</td>

									<td class="px-5 py-4 align-middle whitespace-nowrap text-slate-600">
										{formatDate(item.created_at)}
									</td>

									<td class="px-5 py-4 align-middle">
										<button
											type="button"
											onclick={() => deleteImport(item)}
											disabled={Boolean(deletingID)}
											aria-label={`حذف سفارش ${item.order_name}`}
											class="inline-flex items-center gap-2 rounded-xl px-3 py-2 font-bold text-red-700 transition hover:bg-red-50 focus:ring-2 focus:ring-red-200 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50"
										>
											{#if deletingID === item.uuid}
												<LoaderCircleIcon class="size-4 animate-spin" />
												در حال حذف...
											{:else}
												<Trash2Icon class="size-4" />
												حذف
											{/if}
										</button>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}
		</div>
	</div>
</section>
