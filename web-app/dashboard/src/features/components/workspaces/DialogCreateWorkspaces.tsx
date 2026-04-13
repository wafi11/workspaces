import * as Dialog from "@radix-ui/react-dialog";
import { useState } from "react";
import { FormWorkspaces } from "./FormWorkspace";

export function DialogCreateWorkspaces() {
  const [open, setOpen] = useState(false);

  return (
    <Dialog.Root open={open} onOpenChange={setOpen}>
      <Dialog.Trigger asChild>
        <button
          className="px-3 py-1.5 rounded-md text-sm font-medium transition-colors cursor-pointer"
          style={{
            background: "var(--color-sidebar-surface)",
            border: "1px solid var(--color-sidebar-border)",
            color: "var(--color-sidebar-text-active)",
          }}
        >
          + New Workspace
        </button>
      </Dialog.Trigger>

      <Dialog.Portal>
        <Dialog.Overlay
          className="fixed inset-0 z-50"
          style={{ background: "rgba(0,0,0,0.7)" }}
        />
        <Dialog.Content
          className="fixed z-50 top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-full max-w-lg rounded-xl p-6 flex flex-col gap-5 focus:outline-none"
          style={{
            background: "var(--color-sidebar-bg)",
            border: "1px solid var(--color-sidebar-border)",
          }}
        >
          <div className="flex items-center justify-between">
            <Dialog.Title
              className="text-sm font-semibold"
              style={{ color: "var(--color-sidebar-text-active)" }}
            >
              Create Workspace
            </Dialog.Title>
            <Dialog.Close asChild>
              <button
                className="w-6 h-6 flex items-center justify-center rounded transition-colors cursor-pointer"
                style={{ color: "var(--color-sidebar-text-muted)" }}
              >
                ✕
              </button>
            </Dialog.Close>
          </div>

          <FormWorkspaces onSuccess={() => setOpen(false)} />
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
