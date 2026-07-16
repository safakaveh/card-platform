import { getJsonMap } from "@/lib/server/json";
import UploadViewPage from "./upload_view";
import { getCurrentLanguage } from "@/shared/locale/actions";
import { Language } from "@/shared/locale";

export default async function IndexUploadCSV() {
  const currentLang:Language = await getCurrentLanguage()
  const dictionary = await getJsonMap(`features/upload_csv/dictionaries/page_${currentLang.lang}.json`);
  return (
    <>
      <UploadViewPage dictionary={dictionary} />
    </>
  );
}
