"use client";
import { TopBarAdmin } from "@/components/layouts/TopBarAdmin";
import { useTemplateDetails } from "@/features/services/templates/api";
import { TemplateDetailsCard } from "./Details/DetailsCard";
import { SubBarFilterComponent } from "./Details/SubBarFilter";
import { DetailsError, DetailsLoading } from "./Details/DetailsLoading";
import { FILTERS, useTemplate } from "@/features/hooks/templates/useTemplate";
import { DetailsVariables } from "./Details/DetailsVariables";
import { DetailsFiles } from "./Details/DetailsFiles";
import { DetailsAddons } from "./Details/DetailsAddons";

export function TemplateDetails({ id }: { id: string }) {
  const { data, isLoading, error } = useTemplateDetails(id);
  const { toggleFilter, isOpenFilter } = useTemplate();

  if (isLoading) {
    return <DetailsLoading />;
  }

  if (error || !data) {
    return <DetailsError />;
  }

  const template = data;

  return (
    <>
      <TopBarAdmin title="Details Template" description="Details Template" />
      <TemplateDetailsCard template={template} />
      <SubBarFilterComponent
        filters={FILTERS}
        selected={isOpenFilter}
        setSelected={toggleFilter}
      />
      {isOpenFilter("Variables") && <DetailsVariables templateId={id} />}
      {isOpenFilter("Files") && <DetailsFiles templateId={id} />}
      {isOpenFilter("Add-ons") && <DetailsAddons templateId={id} />}
    </>
  );
}
