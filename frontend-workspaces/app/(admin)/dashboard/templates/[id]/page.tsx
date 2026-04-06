import { TemplateDetails } from "@/features/pages/dashboard/templates/TemplatesDetails";

export default async function Page({ params }: { params: { id: string } }) {
  const { id } = await params;
  return <TemplateDetails id={id} />;
}
