// components/workspace/WorkspaceAddons.tsx
"use client"

import { useState } from "react"
import Image from "next/image"
import { EyeIcon, EyeOff } from "lucide-react"

interface WorkspaceAddon {
  id: string
  name: string
  icon: string
  status: string
  config: Record<string, string>
}

interface WorkspaceAddonsProps {
  addons: WorkspaceAddon[]
}

const PASSWORD_KEYS = ["DB_PASSWORD", "PASSWORD", "SECRET", "TOKEN", "KEY"]

function isPasswordKey(key: string) {
  return PASSWORD_KEYS.some((k) => key.toUpperCase().includes(k))
}

function ConfigRow({ label, value }: { label: string; value: string }) {
  const [shown, setShown] = useState(false)
  const isSensitive = isPasswordKey(label)

  return (
    <div className="flex items-start gap-4 py-1">
      <span className="text-xs text-muted-foreground font-mono whitespace-nowrap w-36 shrink-0">
        {label}
      </span>
      <div className="flex items-center gap-2">
        <span className="text-xs font-mono text-foreground break-all">
          {isSensitive && !shown ? "••••••••••••" : value}
        </span>
        {isSensitive && (
          <button
            type="button"
            onClick={() => setShown((p) => !p)}
            className="text-xs text-muted-foreground underline underline-offset-2 cursor-pointer"
          >
            {shown ? <EyeIcon className="size-4"/> : <EyeOff className="size-4"/>}
          </button>
        )}
      </div>
    </div>
  )
}

function AddonCard({ addon }: { addon: WorkspaceAddon }) {
  const hasStatus = addon.status && addon.status.trim() !== ""
  const configEntries = Object.entries(addon.config)

  return (
    <div className="rounded-lg border bg-card p-4">
      {/* Header */}
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-2.5">
          {addon.icon && (
            <div className="w-6 h-6 shrink-0 relative">
              <img
                src={addon.icon}
                alt={addon.name}
                className="object-contain"
              />
            </div>
          )}
          <span className="text-sm font-medium text-foreground">
            {addon.name}
          </span>
        </div>

        {hasStatus ? (
          <span
            className={`text-xs px-2.5 py-0.5 rounded-full border font-medium ${
              addon.status === "running"
                ? "bg-green-50 text-green-700 border-green-200 dark:bg-green-950 dark:text-green-400 dark:border-green-800"
                : addon.status === "error"
                ? "bg-red-50 text-red-700 border-red-200 dark:bg-red-950 dark:text-red-400 dark:border-red-800"
                : "bg-muted text-muted-foreground border-border"
            }`}
          >
            {addon.status}
          </span>
        ) : (
          <span className="text-xs px-2.5 py-0.5 rounded-full border bg-muted text-muted-foreground">
            no status
          </span>
        )}
      </div>

     
      {/* Config */}
      {configEntries.length > 0 && (
        <div className="border-t pt-3">
          <p className="text-[11px] uppercase tracking-wide text-muted-foreground mb-2">
            Config
          </p>
          <div className="space-y-0.5">
            {configEntries.map(([key, value]) => (
              <ConfigRow key={key} label={key} value={value} />
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

export function WorkspaceAddons({ addons }: WorkspaceAddonsProps) {
  if (!addons || addons.length === 0) {
    return (
      <div className="mt-6">
        <p className="text-sm font-medium mb-3">Addons</p>
        <p className="text-sm text-muted-foreground">No addons installed.</p>
      </div>
    )
  }

  return (
    <div className="mt-6">
      <p className="text-sm font-medium mb-3">Addons</p>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
        {addons.map((addon) => (
          <AddonCard key={addon.id} addon={addon} />
        ))}
      </div>
    </div>
  )
}