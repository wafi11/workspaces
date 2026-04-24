import { TopbarAdmin } from "@/features/layout/TopbarDashboard";
import { useCreateTemplate } from "@/hooks/useCreateTemplates";
import { useResizeFormTemplates } from "@/hooks/useResizeFormTemplates";
import { useCallback } from "react";
import { FormTemplates } from "./FormTemplates";
import { FormTemplatesAddon } from "./FormTemplatesAddon";
import { FormTemplatesVariables } from "./FormTemplatesVariable";
import { FormTemplatesFiles } from "./FormFiles";
import { MainContainer } from "@/features/layout/MainContainer";
import { Button } from "@radix-ui/themes";

export function CreateTemplates() {
  const { initDrag, leftRef, rightRef, rootRef, topLeftRef, topRightRef } =
    useResizeFormTemplates();

  const { handleSubmit, form, variables, addons, files } = useCreateTemplate();

  const onDragH = useCallback((e: MouseEvent) => {
    if (!rootRef.current || !leftRef.current) return;
    const rect = rootRef.current.getBoundingClientRect();
    let pct = ((e.clientX - rect.left) / rect.width) * 100;
    pct = Math.min(Math.max(pct, 15), 85);
    leftRef.current.style.width = pct + "%";
  }, [leftRef,rootRef]);

  const onDragVLeft = useCallback((e: MouseEvent) => {
    if (!leftRef.current || !topLeftRef.current) return;
    const rect = leftRef.current.getBoundingClientRect();
    let pct = ((e.clientY - rect.top) / rect.height) * 100;
    pct = Math.min(Math.max(pct, 15), 85);
    topLeftRef.current.style.height = pct + "%";
  }, [leftRef,topLeftRef]);

  const onDragVRight = useCallback((e: MouseEvent) => {
    if (!rightRef.current || !topRightRef.current) return;
    const rect = rightRef.current.getBoundingClientRect();
    let pct = ((e.clientY - rect.top) / rect.height) * 100;
    pct = Math.min(Math.max(pct, 15), 85);
    topRightRef.current.style.height = pct + "%";
  }, [rightRef,topRightRef]);

  return (
    <>
      <MainContainer>
        <TopbarAdmin title="Create" className="mb-0 w-full" >
          <Button
            type="button"
            onClick={handleSubmit}
            className="py-2 px-4 rounded-md font-semibold flex text-xs sm:text-sm "
            style={{
              background: "var(--color-sidebar-text-active)",
              color: "var(--color-sidebar-bg)",
              cursor: "pointer",
              border: "none",
            }}
          >
            Submit
          </Button>
        </TopbarAdmin>
        <div ref={rootRef} className="flex w-full min-h-screen overflow-y-auto">
          {/* Left panel */}
          <div
            ref={leftRef}
            style={{ width: "50%" }}
            className="overflow-hidden flex flex-col"
          >
            <div
              ref={topLeftRef}
              style={{ height: "50%" }}
              className="overflow-hidden flex flex-col"
            >
              <FormTemplates form={form} />
            </div>
            <div
              onMouseDown={initDrag(onDragVLeft)} // ← onDragV bukan onDragH
              style={{
                height: 2,
                background: "#202020",
                cursor: "row-resize",
                flexShrink: 0,
              }}
            />
            <div className="flex-1 overflow-hidden flex flex-col">
              <FormTemplatesFiles files={files} form={form} />
            </div>
          </div>

          {/* Horizontal divider */}
          <div
            onMouseDown={initDrag(onDragH)}
            style={{
              width: 2,
              background: "#202020",
              cursor: "col-resize",
              flexShrink: 0,
            }}
          />

          {/* Right panels */}
          <div ref={rightRef} className="flex-1  flex flex-col overflow-hidden">
            {/* Top right */}
            <div
              ref={topRightRef}
              style={{ height: "50%" }}
              className="overflow-hidden flex flex-col"
            >
              <FormTemplatesVariables form={form} variables={variables} />
            </div>

            {/* Vertical divider */}
            <div
              onMouseDown={initDrag(onDragVRight)} // ← onDragV bukan onDragH
              style={{
                height: 2,
                background: "#202020",
                cursor: "row-resize",
                flexShrink: 0,
              }}
            />

            {/* Bottom right */}
            <div className="flex-1 overflow-hidden flex flex-col">
              <FormTemplatesAddon addons={addons} form={form} />
            </div>
          </div>
        </div>
      </MainContainer>
    </>
  );
}
