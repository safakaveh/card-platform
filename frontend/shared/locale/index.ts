//@/share/locale/index.ts
import { Languages } from "lucide-react";

export type Language = {
  lang: "fa" | "en";
  dir: "rtl" | "ltr";
};
export const LANG_COOKIE: string = "language";

export const DEFAULT_LANGUAGE: Language = { lang: "fa", dir: "rtl" };

export function normalizeLang(value?: string | null): Language {
  return value === "en" ? { lang: "en", dir: "ltr" } : { lang: "fa", dir: "rtl" };
}
