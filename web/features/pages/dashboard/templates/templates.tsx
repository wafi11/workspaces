"use client";
import { TopBarAdmin } from "@/components/layouts/TopBarAdmin";
import { Button } from "@/components/ui/button";
import { useTemplate } from "@/features/hooks/templates/useTemplate";
import { cn } from "@/lib/utils";
import { Plus } from "lucide-react";
import { useRouter } from "next/navigation";
import { EditForm } from "./EditForm";
import { SlideOver } from "./SlideOver";
import { TemplateCard } from "./TemplatesCard";
export function TemplatesPage() {
  const {
    closeEdit,
    form,
    setOpenDialogEdit,
    handleChange,
    openEdit,
    selected,
    templateData,
  } = useTemplate();

  const { push } = useRouter();

  return (
    <>
      <TopBarAdmin title="Templates" description="Manage Workspaces Templates">
        <Button
          onClick={() => push("/dashboard/templates/create")}
          variant="outline"
          size="sm"
        >
          <Plus className="w-3.5 h-3.5 mr-2" />
          Add Templates
        </Button>
      </TopBarAdmin>

      {/* Grid */}
      <div
        className={cn(
          "grid gap-4 mt-6 transition-all duration-300",
          selected
            ? "grid-cols-1 sm:grid-cols-2"
            : "grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"
        )}
      >
        {templateData?.map((item) => (
          <TemplateCard
            key={item.id}
            item={item}
            isActive={selected?.id === item.id}
            onClick={() => openEdit(item)}
          />
        ))}
      </div>

      {/* Slide-over panel */}
      <SlideOver open={!!selected} onClose={closeEdit}>
        {selected && (
          <EditForm form={form} onChange={handleChange} onClose={closeEdit} />
        )}
      </SlideOver>
    </>
  );
}
