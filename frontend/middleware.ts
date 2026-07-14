import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";
import { DEFAULT_LANG, LANG_COOKIE } from "./shared/locale/i18n";

export function middleware(req: NextRequest) {
  const res: NextResponse = NextResponse.next();

  const lang: string | undefined | null = req.cookies.get(LANG_COOKIE)?.value;
  if (!lang) {
    res.cookies.set(LANG_COOKIE, DEFAULT_LANG, {
      path: "/",
      sameSite: "lax",
      maxAge: 60 * 60 * 24 * 365,
    });
  }

  return res;
}

export const config = {
  matcher: ["/((?!_next|favicon.ico).*)"],
};
