import { TemplateCreate } from "@/features/pages/profile/templates/TemplateCreate";

export default async function Page({params} : {params : {id : string}}){
    const {id}  = await params
    return (
       <TemplateCreate id={id}/>
    )
}