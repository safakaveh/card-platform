//@/share/locale/client_action.ts
"use client";

import { useRouter } from "next/navigation";
import { changeLanguage, getCurrentLanguage } from "./actions";
import { Language } from ".";
import { useTransition } from "react";

export function useChangeLanguage() {
  const router = useRouter();
  const [pending, startTransition] = useTransition();

  function setLanguage(lang: string) {
    startTransition(async () => {
      await changeLanguage(lang);
      router.refresh();
    });
  }

  return {
    pending,
    setLanguage,
  };
}
