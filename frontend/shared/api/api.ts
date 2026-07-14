type Method = "GET" | "QUERY" | "POST" | "PUT" | "DELETE" | "PATCH";

// TResponse: تایپ خروجی API
// TBody: تایپ بدنه درخواستی که ارسال می‌شود (برای POST/PUT)
class ApiRequest<ResponseType, BodyType> {
  private url: string;
  private _method: Method = "GET";
  private headers: Headers = new Headers();
  private queryParams: Record<string, any> = {};
  private _body?: BodyType;

  constructor(url: string) {
    this.url = url;
    // به صورت پیش‌فرض هدر Content-Type را روی JSON می‌گذاریم
    this.headers.set("Content-Type", "application/json");
  }

  method(m: Method) {
    this._method = m;
    return this;
  }

  query(params: Record<string, any>) {
    this.queryParams = params;
    return this;
  }

  header(key: string, value: string) {
    this.headers.set(key, value);
    return this;
  }

  // متد جدید برای دریافت Body (با Type Safety)
  body(data: BodyType) {
    this._body = data;
    return this;
  }

  // متد نهایی برای اجرای درخواست و بازگرداندن پاسخ با تایپ مشخص
  async json(): Promise<ResponseType> {
    const url = new URL(this.url, typeof window !== "undefined" ? window.location.origin : undefined);

    Object.keys(this.queryParams).forEach((key) => url.searchParams.append(key, String(this.queryParams[key])));

    // ارسال درخواست از طریق HttpClient (که در مرحله قبل تعریف کردیم)
    const response = await fetch(url.toString(), {
      method: this._method,
      headers: this.headers,
      body: this._body ? JSON.stringify(this._body) : undefined,
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response.json() as Promise<ResponseType>;
  }
}

// ساختار اصلی برای شروع زنجیره درخواست‌ها
export const api = {
  get: <TRes>(url: string) => new ApiRequest<TRes, never>(url).method("GET"),

  query: <TRes, TReq = any>(url: string) => new ApiRequest<TRes, TReq>(url).method("QUERY"),

  post: <TRes, TReq = any>(url: string) => new ApiRequest<TRes, TReq>(url).method("POST"),

  put: <TRes, TReq = any>(url: string) => new ApiRequest<TRes, TReq>(url).method("PUT"),

  delete: <TRes>(url: string) => new ApiRequest<TRes, never>(url).method("DELETE"),
};
