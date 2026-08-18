# Card Platform

سامانه‌ی Go برای بارگذاری CSV، شخصی‌سازی و آماده‌سازی چاپ کارت‌های بانکی. رابط وب به‌صورت embedded در باینری قرار دارد و داده‌ها در SQLite ذخیره می‌شوند.

## اجرا

نیازمندی: Go 1.26 یا جدیدتر.

```bash
go run ./cmd/app
```

پس از اجرا، آدرس رابط وب در لاگ نمایش داده می‌شود. تنظیمات پایگاه‌داده از متغیر `DB_DATASOURCE_NAME` خوانده می‌شود و مقدار پیش‌فرض `data.db` است.

## قالب CSV

ردیف اول باید header باشد. ستون‌های `frn_` و `bck_` به‌ترتیب داده‌های لیزر جلو و پشت هستند؛ `frn_img_`/`bck_img_` مسیر تصویر را می‌پذیرند و محتوای فایل در SQLite ذخیره می‌شود. ستون‌های `trk1_` تا `trk3_` در ترک‌های مگنت ذخیره می‌شوند. مقدار `1` یا `true` در `frn_uid`/`bck_uid` یک رکورد `UUID_PENDING` در `mifare_data` ایجاد می‌کند تا در مرحله‌ی خواندن کارت با UUID واقعی جایگزین شود.

هر کارت و هر مرحله‌ی پردازش در `card_status_history` ثبت می‌شود؛ وضعیت اولیه‌ی واردسازی `loaded` است.

## API

- `POST /api/imports/` با فرم multipart شامل `file` و اختیاری `order_name`
- `GET /api/imports/` فهرست واردسازی‌ها
- `GET /api/imports/{id}` جزئیات واردسازی
- `DELETE /api/imports/{id}` حذف سفارش و داده‌های وابسته

منوی «فایل» در رابط وب شامل بارگذاری CSV، ویرایش header و خروج از برنامه است. فایل‌های وب در `internal/web/build` قرار گرفته و با `go:embed` بسته می‌شوند.

## توسعه

```bash
go test ./...
GOOS=windows GOARCH=amd64 go build -o bin/myapp.exe ./cmd/app
```
