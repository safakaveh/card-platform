//@/share/locale/action.ts
"use server";

import { cookies } from "next/headers";
import { revalidatePath } from "next/cache";
import { DEFAULT_LANGUAGE, LANG_COOKIE, Language, normalizeLang } from ".";

async function changeLanguageCookies(lang: Language) {
  const cookieStore = await cookies();
  cookieStore.set(LANG_COOKIE, JSON.stringify(lang), {
    path: "/",
    sameSite: "lax",
    maxAge: 60 * 60 * 24 * 365,
  });
  revalidatePath("/", "layout");
}

export async function getCurrentLanguage(): Promise<Language> {
  const cookieStore = await cookies();
  const strCurrentLang: string | undefined = cookieStore.get(LANG_COOKIE)?.value;

  if (!strCurrentLang) {
    await changeLanguageCookies(DEFAULT_LANGUAGE);
    return DEFAULT_LANGUAGE;
  }
  return JSON.parse(strCurrentLang);
}

export async function changeLanguage(lang: string) {
  await changeLanguageCookies(normalizeLang(lang));
}
