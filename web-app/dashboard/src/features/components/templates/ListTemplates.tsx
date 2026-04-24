import { useTemplates } from "@/features/api";
import type { User } from "@/types";
import { TemplateCard } from "./TemplateCard";

interface ListTemplatesProps {
  profile : User
}

export function ListTemplates({profile} : ListTemplatesProps) {
  const { data: templates, isLoading } = useTemplates();

  if (isLoading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 p-6">
        {[...Array(6)].map((_, i) => (
          <div key={i} className="h-48 rounded-xl bg-sidebar-surface animate-pulse border border-sidebar-border" />
        ))}
      </div>
    );
  }

  return (
    <div className="p-6 space-y-3 overflow-y-auto">
      {/* Grid Layout */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {templates?.map((template) => (
          <TemplateCard profile={profile} template={template}/>
        ))}
      </div>
    </div>
  );
}