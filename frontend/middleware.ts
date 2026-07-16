import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";
import { DEFAULT_LANGUAGE, LANG_COOKIE } from "./shared/locale";

export function middleware(req: NextRequest) {
  const res: NextResponse = NextResponse.next();

  const lang: string | undefined | null = req.cookies.get(LANG_COOKIE)?.value;
  if (!lang) {
    res.cookies.set(LANG_COOKIE, JSON.stringify(DEFAULT_LANGUAGE), {
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
