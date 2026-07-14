export interface TokenProvider {
  getAccessToken(): Promise<string | null>;
}
