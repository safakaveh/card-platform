export class ApiResult<T> {
  constructor(
    public readonly data: T,

    public readonly status: number,
  ) {}
}
