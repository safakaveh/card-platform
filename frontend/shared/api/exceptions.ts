export class ApiException extends Error {
  constructor(
    message: string,

    readonly status: number,
  ) {
    super(message);
  }
}
