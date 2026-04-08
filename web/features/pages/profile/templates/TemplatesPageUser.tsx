"use client"
import {  useTemplates } from "@/features/services/templates"
import { cn } from "@/lib/utils"
import { TopBarAdmin } from "@/components/layouts/TopBarAdmin"
import { TemplateCard } from "./TemplateCard"

export function TemplatesPageUser(){
    const {data}  = useTemplates()
    console.log(data)
    return (
        <>
        <TopBarAdmin description="Choose Your Own Templates" title="Templates"/>
        <div className={cn(
                 "grid gap-4 mt-6 transition-all duration-300 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"   
               )}
             >
               {data?.map((item) => (
                 <TemplateCard key={item.id} item={item}/>
               ))}
             </div>
        </>
    )
}