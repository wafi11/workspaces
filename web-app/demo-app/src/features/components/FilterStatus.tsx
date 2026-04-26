import { faChevronDown, faCheck } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { useState, useRef, useEffect } from "react";

export type StatusFilter = "all" | "running" | "stopped" | "pending";

interface WorkspaceStatusFilterProps {
  value: StatusFilter;
  onChange: (value: StatusFilter) => void;
}

const options: { value: StatusFilter; label: string; dot?: string }[] = [
  { value: "all", label: "All status" },
  { value: "running", label: "Running", dot: "#22c55e" },
  { value: "stopped", label: "Stopped", dot: "#555" },
  { value: "pending", label: "Pending", dot: "#f59e0b" },
];

export function WorkspaceStatusFilter({ value, onChange }: WorkspaceStatusFilterProps) {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);
  const current = options.find((o) => o.value === value)!;

  // tutup kalau klik di luar
  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <div ref={ref} className="relative">
      {/* Trigger */}
      <button
        onClick={() => setOpen((o) => !o)}
        className="flex items-center gap-2 px-3 py-1.5 rounded-md cursor-pointer"
        style={{
          background: "var(--color-app-surface)",
          border: `0.5px solid ${open ? "var(--color-accent-border)" : "var(--color-app-border)"}`,
          color: "var(--color-sidebar-text)",
        }}
      >
        {current.dot && (
          <div
            className="w-1.5 h-1.5 rounded-full shrink-0"
            style={{ background: current.dot }}
          />
        )}
        <span className="text-xs">{current.label}</span>
        <FontAwesomeIcon
          icon={faChevronDown}
          style={{
            color: "var(--color-sidebar-text-muted)",
            fontSize: "10px",
            transition: "transform 0.15s",
            transform: open ? "rotate(180deg)" : "rotate(0deg)",
          }}
        />
      </button>

      {/* Dropdown */}
      {open && (
        <div
          className="absolute top-full mt-1 left-0 z-50 rounded-md overflow-hidden py-1 min-w-[140px]"
          style={{
            background: "var(--color-app-surface)",
            border: "0.5px solid var(--color-app-border)",
          }}
        >
          {options.map((o) => (
            <div
              key={o.value}
              onClick={() => { onChange(o.value); setOpen(false); }}
              className="flex items-center gap-2 px-3 py-2 cursor-pointer text-xs"
              style={{
                color: value === o.value
                  ? "var(--color-sidebar-text-active)"
                  : "var(--color-sidebar-text)",
                background: "transparent",
              }}
              onMouseEnter={(e) => (e.currentTarget.style.background = "#111")}
              onMouseLeave={(e) => (e.currentTarget.style.background = "transparent")}
            >
              {o.dot ? (
                <div className="w-1.5 h-1.5 rounded-full shrink-0" style={{ background: o.dot }} />
              ) : (
                <div className="w-1.5 h-1.5 shrink-0" />
              )}
              <span className="flex-1">{o.label}</span>
              {value === o.value && (
                <FontAwesomeIcon
                  icon={faCheck}
                  style={{ fontSize: "10px", color: "var(--color-accent-text)" }}
                />
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}