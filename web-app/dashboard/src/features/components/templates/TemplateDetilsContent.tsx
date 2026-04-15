import { MainContainer } from "@/features/layout/MainContainer";
import { SectionAddOn } from "./SectionAddon";
import { useTemplateDetails } from "./TemplateDetailsContext";
import { TemplateVariable } from "./TemplateVariable";
import { TemplateFiles } from "./TemplateFiles";

export function TemplateDetailsContent() {
  const {
    data,
   
  } = useTemplateDetails()

  if (!data) return null

  return (
    <MainContainer>
      {/* ── Header ────────────────────────────────────────── */}
      <div className="flex items-center gap-3 px-5 py-4 border-b border-[#111]">
        {data.icon && (
          <img
            src={data.icon}
            width={32}
            height={32}
            className="rounded-sm border border-[#1c1c1c] object-contain shrink-0"
          />
        )}
        <div className="flex flex-col gap-0.5 flex-1 min-w-0">
          <div className="flex items-center gap-2">
            <span className="text-[13px] text-[#e4e4e7] font-medium tracking-[0.01em]">
              {data.name}
            </span>
            {data.is_public && (
              <span className="text-[9px] px-1.5 py-0.5 rounded-[3px] bg-[#111] border border-[#1c1c1c] text-sidebar-text tracking-widest uppercase">
                public
              </span>
            )}
          </div>
          <span className="text-[11px] text-sidebar-text">{data.description}</span>
        </div>
        <span className="text-[9px] px-2 py-0.5 rounded-[3px] bg-[#111] border border-[#1c1c1c] text-sidebar-text tracking-widest uppercase shrink-0">
          {data.category}
        </span>
      </div>

      <div className="flex flex-col gap-0">

        {/* ── Variables ──────────────────────────────────── */}
       <TemplateVariable />

        {/* ── Addons ────────────────────────────────────── */}
        <SectionAddOn />

        {/* ── Files ─────────────────────────────────────── */}
        <TemplateFiles />
      </div>
    </MainContainer>
  )
}