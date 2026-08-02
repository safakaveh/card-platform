<script lang="ts">
	import CloudUploadIcon from '@lucide/svelte/icons/cloud-upload';
	import FileSpreadsheetIcon from '@lucide/svelte/icons/file-spreadsheet';
	import CircleCheckIcon from '@lucide/svelte/icons/circle-check';
	import TriangleAlertIcon from '@lucide/svelte/icons/triangle-alert';
	import XIcon from '@lucide/svelte/icons/x';
	import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';

	type ImportResult = {
		uuid: string;
		order_name: string;
		file_name: string;
		rows_imported: number;
		front_columns: string[];
		back_columns: string[];
		has_uid: boolean;
	};

	type UploadStage = 'uploading' | 'processing' | null;

	const MAX_FILE_SIZE = 100 * 1024 * 1024; // 100 MB

	let fileInput: HTMLInputElement;
	let activeRequest: XMLHttpRequest | null = null;

	let selectedFile = $state<File | null>(null);
	let orderName = $state('');
	let isDragging = $state(false);
	let isUploading = $state(false);
	let uploadStage = $state<UploadStage>(null);
	let progress = $state(0);
	let errorMessage = $state('');
	let result = $state<ImportResult | null>(null);

	const fileSize = $derived(
		selectedFile
			? new Intl.NumberFormat('fa-IR', {
					maximumFractionDigits: 1
				}).format(selectedFile.size / 1024 / 1024) + ' مگابایت'
			: ''
	);

	const formattedProgress = $derived(
		new Intl.NumberFormat('fa-IR').format(progress)
	);

	function chooseFile(file: File | undefined) {
		errorMessage = '';
		result = null;

		if (!file) return;

		if (!file.name.toLowerCase().endsWith('.csv')) {
			selectedFile = null;
			errorMessage = 'لطفاً فقط یک فایل CSV انتخاب کنید.';
			resetFileInput();
			return;
		}

		if (file.size === 0) {
			selectedFile = null;
			errorMessage = 'فایل انتخاب‌شده خالی است.';
			resetFileInput();
			return;
		}

		if (file.size > MAX_FILE_SIZE) {
			selectedFile = null;
			errorMessage = 'حجم فایل نباید بیشتر از ۱۰۰ مگابایت باشد.';
			resetFileInput();
			return;
		}

		selectedFile = file;

		if (!orderName.trim()) {
			orderName = file.name.replace(/\.csv$/i, '');
		}
	}

	function resetFileInput() {
		if (fileInput) {
			fileInput.value = '';
		}
	}

	function removeFile() {
		selectedFile = null;
		result = null;
		errorMessage = '';
		progress = 0;
		resetFileInput();
	}

	function onDrop(event: DragEvent) {
		event.preventDefault();
		isDragging = false;

		const file = event.dataTransfer?.files?.[0];
		chooseFile(file);
	}

	function getResponseError(request: XMLHttpRequest): string {
		const response = request.response;

		if (response && typeof response === 'object' && 'error' in response) {
			return String(response.error);
		}

		if (typeof response === 'string' && response.trim()) {
			return response;
		}

		return request.status
			? `بارگذاری فایل ناموفق بود. کد خطا: ${request.status}`
			: 'بارگذاری فایل ناموفق بود.';
	}

	function finishRequest() {
		activeRequest = null;
		isUploading = false;
		uploadStage = null;
	}

	function validateForm(): boolean {
		errorMessage = '';

		if (!selectedFile) {
			errorMessage = 'لطفاً ابتدا یک فایل CSV انتخاب کنید.';
			return false;
		}

		const normalizedOrderName = orderName.trim();

		if (!normalizedOrderName) {
			errorMessage = 'لطفاً نام سفارش را وارد کنید.';
			return false;
		}

		if (normalizedOrderName.length > 150) {
			errorMessage = 'نام سفارش نباید بیشتر از ۱۵۰ کاراکتر باشد.';
			return false;
		}

		return true;
	}

	function upload() {
		if (isUploading || !validateForm() || !selectedFile) return;

		result = null;
		progress = 0;
		isUploading = true;
		uploadStage = 'uploading';

		const body = new FormData();
		body.append('order_name', orderName.trim());
		body.append('file', selectedFile);

		const request = new XMLHttpRequest();
		activeRequest = request;

		request.open('POST', '/api/imports/');
		request.responseType = 'json';
		request.timeout = 30 * 60 * 1000; // 30 دقیقه

		request.upload.onprogress = (event) => {
			if (event.lengthComputable) {
				progress = Math.min(99, Math.round((event.loaded / event.total) * 100));

				if (event.loaded === event.total) {
					uploadStage = 'processing';
				}
			}
		};

		request.onload = () => {
			finishRequest();

			if (request.status >= 200 && request.status < 300) {
				progress = 100;
				result = request.response as ImportResult;
				return;
			}

			errorMessage = getResponseError(request);
		};

		request.onerror = () => {
			finishRequest();
			errorMessage = 'ارتباط با سرور برقرار نشد. لطفاً اتصال شبکه را بررسی کنید.';
		};

		request.ontimeout = () => {
			finishRequest();
			errorMessage = 'زمان بارگذاری به پایان رسید. لطفاً دوباره تلاش کنید.';
		};

		request.onabort = () => {
			finishRequest();
			errorMessage = 'بارگذاری فایل لغو شد.';
		};

		request.send(body);
	}

	function cancelUpload() {
		activeRequest?.abort();
	}
</script>

<section
	dir="rtl"
	class="mx-auto min-h-full w-full max-w-4xl rounded-3xl bg-gradient-to-br from-sky-50 via-white to-indigo-50 p-4 sm:p-6"
>
	<div class="mb-6 text-center">
		<p class="mb-2 text-sm font-semibold text-sky-700">ورود اطلاعات کارت‌ها</p>

		<h1 class="text-2xl font-black text-slate-900 sm:text-3xl">
			بارگذاری فایل CSV
		</h1>

		<p class="mx-auto mt-3 max-w-2xl text-sm leading-7 text-slate-600">
			فایل به‌صورت جریانی پردازش می‌شود؛ بنابراین برای فایل‌های حجیم نیز حافظه‌ی زیادی مصرف
			نخواهد شد.
		</p>
	</div>

	<div
		class="rounded-3xl border border-white/80 bg-white/95 p-5 shadow-xl shadow-slate-200/70 backdrop-blur sm:p-8"
	>
		<label class="mb-2 block text-sm font-bold text-slate-700" for="order-name">
			نام سفارش
		</label>

		<input
			id="order-name"
			bind:value={orderName}
			disabled={isUploading}
			maxlength="150"
			placeholder="مثلاً کارت پرسنلی مرداد"
			class="mb-5 w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none transition placeholder:text-slate-400 focus:border-sky-500 focus:ring-2 focus:ring-sky-200 disabled:cursor-not-allowed disabled:opacity-60"
		/>

		<p id="file-help" class="sr-only">
			فقط فایل‌های CSV با حجم حداکثر ۱۰۰ مگابایت قابل انتخاب هستند.
		</p>

		<button
			type="button"
			aria-label="انتخاب یا رها کردن فایل CSV"
			aria-describedby="file-help"
			disabled={isUploading}
			onclick={() => fileInput.click()}
			ondragenter={(event) => {
				event.preventDefault();
				isDragging = true;
			}}
			ondragover={(event) => {
				event.preventDefault();
				isDragging = true;
			}}
			ondragleave={(event) => {
				event.preventDefault();
				isDragging = false;
			}}
			ondrop={onDrop}
			class={[
				'flex min-h-64 w-full flex-col items-center justify-center rounded-2xl border-2 border-dashed px-6 text-center transition',
				isDragging
					? 'border-sky-500 bg-sky-50'
					: 'border-slate-300 bg-slate-50/70 hover:border-sky-400 hover:bg-sky-50/60',
				isUploading ? 'cursor-not-allowed opacity-60' : 'cursor-pointer'
			]}
		>
			<span class="mb-4 rounded-2xl bg-sky-100 p-4 text-sky-700">
				{#if selectedFile}
					<FileSpreadsheetIcon class="size-10" />
				{:else}
					<CloudUploadIcon class="size-10" />
				{/if}
			</span>

			{#if selectedFile}
				<strong class="max-w-full truncate text-base text-slate-900">
					{selectedFile.name}
				</strong>

				<span class="mt-2 text-sm text-slate-500">
					{fileSize}
				</span>

				<span class="mt-3 text-xs font-semibold text-sky-700">
					برای انتخاب فایل دیگر کلیک کنید
				</span>
			{:else}
				<strong class="text-base text-slate-900">
					فایل را اینجا رها کنید یا کلیک کنید
				</strong>

				<span class="mt-2 text-sm text-slate-500">
					فقط فرمت CSV، حداکثر ۱۰۰ مگابایت
				</span>
			{/if}
		</button>

		<input
			class="hidden"
			type="file"
			accept=".csv,text/csv"
			bind:this={fileInput}
			disabled={isUploading}
			onchange={(event) => {
				const input = event.currentTarget as HTMLInputElement;
				chooseFile(input.files?.[0]);
			}}
		/>

		{#if selectedFile && !isUploading}
			<button
				type="button"
				onclick={removeFile}
				class="mt-3 flex w-full items-center justify-center gap-2 rounded-xl border border-red-200 px-4 py-2.5 text-sm font-bold text-red-700 transition hover:bg-red-50"
			>
				<XIcon class="size-4" />
				حذف فایل انتخاب‌شده
			</button>
		{/if}

		<div class="mt-5 grid gap-3 rounded-2xl bg-slate-50 p-4 text-sm text-slate-600 sm:grid-cols-2">
			<p>
				<code class="font-bold text-sky-700">frn_*</code>
				اطلاعات روی قسمت جلوی کارت
			</p>

			<p>
				<code class="font-bold text-indigo-700">bck_*</code>
				اطلاعات روی قسمت پشت کارت
			</p>

			<p>
				<code class="font-bold text-sky-700">frn_uid</code>
				لیزر UID در جلوی کارت
			</p>

			<p>
				<code class="font-bold text-indigo-700">bck_uid</code>
				لیزر UID در پشت کارت
			</p>
		</div>

		{#if isUploading}
			<div class="mt-6" aria-live="polite">
				<div class="mb-2 flex items-center justify-between text-sm font-semibold text-slate-700">
					<span class="flex items-center gap-2">
						<LoaderCircleIcon class="size-4 animate-spin text-sky-600" />

						{uploadStage === 'processing'
							? 'در حال پردازش و ذخیره‌سازی...'
							: 'در حال ارسال فایل...'}
					</span>

					<span>{formattedProgress}٪</span>
				</div>

				<div
					class="h-3 overflow-hidden rounded-full bg-slate-200"
					role="progressbar"
					aria-valuemin="0"
					aria-valuemax="100"
					aria-valuenow={progress}
				>
					<div
						class="h-full rounded-full bg-gradient-to-l from-sky-500 to-blue-600 transition-all duration-300"
						style={`width: ${progress}%`}
					></div>
				</div>

				<button
					type="button"
					onclick={cancelUpload}
					class="mt-3 w-full rounded-xl border border-red-200 px-4 py-2.5 text-sm font-bold text-red-700 transition hover:bg-red-50"
				>
					لغو بارگذاری
				</button>
			</div>
		{/if}

		{#if errorMessage}
			<div
				class="mt-5 flex items-start gap-3 rounded-xl border border-red-200 bg-red-50 p-4 text-sm leading-6 text-red-800"
				role="alert"
			>
				<TriangleAlertIcon class="mt-0.5 size-5 shrink-0" />
				<span>{errorMessage}</span>
			</div>
		{/if}

		{#if result}
			<div class="mt-5 rounded-xl border border-emerald-200 bg-emerald-50 p-5 text-emerald-900">
				<div class="flex items-center gap-2 font-black">
					<CircleCheckIcon class="size-6" />
					فایل با موفقیت ذخیره شد
				</div>

				<div class="mt-4 grid gap-3 text-sm sm:grid-cols-2">
					<p>
						تعداد کارت‌ها:
						<strong>{result.rows_imported.toLocaleString('fa-IR')}</strong>
					</p>

					<p>
						ستون‌های جلو:
						<strong>{result.front_columns.length.toLocaleString('fa-IR')}</strong>
					</p>

					<p>
						ستون‌های پشت:
						<strong>{result.back_columns.length.toLocaleString('fa-IR')}</strong>
					</p>

					<p>
						ستون UID:
						<strong>{result.has_uid ? 'دارد' : 'ندارد'}</strong>
					</p>
				</div>
			</div>
		{/if}

		<button
			type="button"
			onclick={upload}
			disabled={!selectedFile || !orderName.trim() || isUploading}
			class="mt-6 w-full rounded-xl bg-slate-900 px-5 py-3.5 text-sm font-bold text-white shadow-lg transition hover:bg-sky-700 focus:outline-none focus:ring-2 focus:ring-sky-400 disabled:cursor-not-allowed disabled:opacity-40"
		>
			{isUploading ? 'لطفاً منتظر بمانید...' : 'شروع بارگذاری و ذخیره'}
		</button>
	</div>
</section>
