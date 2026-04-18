// dropdown.tsx
"use client"

import { Ellipsis, Pencil, Trash2, Users } from "lucide-react"
import { useEffect, useRef, useState } from "react"
import { DialogAddCollaboration } from "./DialogAddCollaboration"

export const dropdown_choices = ["Add Collaboration", "Update", "Delete"] as const
type DropdownChoice = (typeof dropdown_choices)[number]

export function DropdownChoice({wsId} : {wsId : string}) {
  const [open, setOpen] = useState(false)
  const [dialogOpen, setDialogOpen] = useState(false)
  const ref = useRef<HTMLDivElement>(null)

  // close on outside click
  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setOpen(false)
      }
    }
    document.addEventListener("mousedown", handler)
    return () => document.removeEventListener("mousedown", handler)
  }, [])

  const handleSelect = (choice: DropdownChoice) => {
    setOpen(false)
    if (choice === "Add Collaboration") {
      setDialogOpen(true)
    }
    // handle Update / Delete di sini
  }

  return (
    <>
      <div ref={ref} className="relative">
        {/* Trigger */}
        <button
          onClick={() => setOpen((v) => !v)}
          className="flex items-center gap-2 bg-zinc-900 border border-zinc-700 hover:border-zinc-500 
                     text-zinc-300 text-xs px-1.5 py-1.5 rounded-full font-mono transition-colors"
        >
          <Ellipsis
            size={12}
            className={`text-zinc-500 transition-transform ${open ? "rotate-180" : ""}`}
          />
        </button>

        {/* Menu */}
        {open && (
          <div
            className="absolute top-[calc(100%+6px)] left-0 z-50 min-w-45 
                       bg-zinc-900 border border-zinc-700 rounded-md overflow-hidden
                       shadow-[0_8px_24px_rgba(0,0,0,0.6)]
                       animate-in fade-in-0 slide-in-from-top-1 duration-150"
          >
            <MenuItem icon={<Users size={13} />} onClick={() => handleSelect("Add Collaboration")}>
              Add Collaboration
            </MenuItem>
            <MenuItem icon={<Pencil size={13} />} onClick={() => handleSelect("Update")}>
              Update
            </MenuItem>
            <MenuItem
              icon={<Trash2 size={13} />}
              onClick={() => handleSelect("Delete")}
              danger
            >
              Delete
            </MenuItem>
          </div>
        )}
      </div>

      {/* Dialog trigger programmatically */}
      <DialogAddCollaboration wsId={wsId} open={dialogOpen} onOpenChange={setDialogOpen} />
    </>
  )
}

function MenuItem({
  icon,
  children,
  onClick,
  danger,
}: {
  icon: React.ReactNode
  children: React.ReactNode
  onClick: () => void
  danger?: boolean
}) {
  return (
    <button
      onClick={onClick}
      className={`flex items-center gap-2.5 w-full px-3.5 py-2.5 text-xs font-mono
                  border-b border-zinc-800 last:border-0 transition-colors
                  ${danger
                    ? "text-red-400 hover:bg-red-950/40 hover:text-red-300"
                    : "text-zinc-400 hover:bg-zinc-800 hover:text-zinc-200"
                  }`}
    >
      {icon}
      {children}
    </button>
  )
}