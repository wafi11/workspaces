import { TopBarAdmin } from "@/components/layouts/TopBarAdmin";

export function DetailsLoading() {
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

export function DetailsError() {
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
