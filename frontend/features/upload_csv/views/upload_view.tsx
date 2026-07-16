"use client";
import { Button } from "@/components/ui/button";
import { useChangeLanguage } from "@/shared/locale/uselanguage";

type Props = {
    dictionary: Map<string, string>;
};

export default function UploadViewPage({ dictionary }: Props) {
    const  {setLanguage,pending} = useChangeLanguage();
    return (
        <>
            <h5>{dictionary.get("title")}</h5>
            <h3>{dictionary.get("description")}</h3>
            <h5>{dictionary.get("test")}</h5> 
            <div>
                <Button  variant="outline" className="rounded-4xl me-3" onClick={()=>setLanguage("fa")}> فارسی </Button>
                <Button  variant="outline" onClick={()=>setLanguage("en")}> English </Button>
            </div>
        </>
    )
}
