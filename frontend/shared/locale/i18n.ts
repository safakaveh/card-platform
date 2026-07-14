import { AppRouterInstance } from "next/dist/shared/lib/app-router-context.shared-runtime";
export type Lang = "fa" | "en";
export type Dir = "rtl" | "ltr";

export const DEFAULT_LANG: Lang = "fa";
export const LANG_COOKIE = "lang";

export function dirOf(lang: Lang): Dir {
  return lang === "fa" ? "rtl" : "ltr";
}

export function normalizeLang(v: string | undefined | null): Lang {
  return v === "en" ? "en" : "fa";
}

function setCookie(name: string, value: string, days = 365) {
  const maxAge = days * 24 * 60 * 60;
  document.cookie = `${name}=${encodeURIComponent(value)}; Path=/; Max-Age=${maxAge}; SameSite=Lax`;
}

export const changeLanguage = (lang: Lang, router: AppRouterInstance) => {
  setCookie(LANG_COOKIE, lang);
  router.refresh();
};
