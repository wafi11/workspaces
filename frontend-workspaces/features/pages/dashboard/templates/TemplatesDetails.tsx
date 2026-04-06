"use client";
import { TopBarAdmin } from "@/components/layouts/TopBarAdmin";
import { useGetTemplateDetails } from "@/features/services/templates/api";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import Image from "next/image";
import {
  Database,
  Globe,
  Tag,
  Calendar,
  ExternalLink,
  Copy,
} from "lucide-react";
import { TemplateDetailsCard } from "./Details/DetailsCard";

export function TemplateDetails({ id }: { id: string }) {
  const { data, isLoading, error } = useGetTemplateDetails(id);

  if (isLoading) {
    return (
      <>
        <TopBarAdmin title="Details Template" description="Details Template" />
        <div className="p-6 space-y-6">
          <div className="h-32 rounded-lg bg-muted animate-pulse" />
          <div className="grid grid-cols-2 gap-4">
            <div className="h-20 rounded-lg bg-muted animate-pulse" />
            <div className="h-20 rounded-lg bg-muted animate-pulse" />
          </div>
        </div>
      </>
    );
  }

  if (error || !data) {
    return (
      <>
        <TopBarAdmin title="Details Template" description="Details Template" />
        <div className="p-6">
          <p className="text-muted-foreground text-sm">
            Failed to load template.
          </p>
        </div>
      </>
    );
  }

  const template = data.data;

  return (
    <>
      <TopBarAdmin title="Details Template" description="Details Template" />

      <TemplateDetailsCard template={template} />
    </>
  );
}
