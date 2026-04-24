import { useState } from "react";

interface SectionTerminalProps {
  terminal_url: string;
}

export function SectionTerminal({ terminal_url }: SectionTerminalProps) {
  const [isFullscreen, setIsFullscreen] = useState(false);
console.log(terminal_url)
  return (
    <div
      className={`flex flex-col bg-[var(--color-app-bg)] border border-[var(--color-app-border)] m-4 p-4 rounded-xl overflow-hidden transition-all duration-200 ${
        isFullscreen
          ? "fixed inset-4 z-50"
          : "w-full  h-full min-h-[200px]"
      }`}
    >
      {/* Title bar */}
      <div className="flex items-center justify-between px-4 py-2.5 bg-[var(--color-sidebar-bg)] border-b border-[var(--color-app-border)] shrink-0">
        {/* Traffic lights */}
        <div className="flex items-center gap-1.5">
          <div className="w-2.5 h-2.5 rounded-full bg-[#ff5f57]" />
          <div className="w-2.5 h-2.5 rounded-full bg-[#febc2e]" />
          <div className="w-2.5 h-2.5 rounded-full bg-[#28c840]" />
        </div>

        {/* Label tengah */}
        <div className="flex items-center gap-2">
          <div className="w-1.5 h-1.5 rounded-full bg-[#28c840] animate-pulse" />
          <span className="font-mono text-[11px] text-[var(--color-sidebar-text)] tracking-wide">
            terminal — workspace
          </span>
        </div>

        {/* Actions */}
        <div className="flex items-center gap-1">
          <button
            onClick={() => setIsFullscreen((p) => !p)}
            className="p-1.5 rounded-md text-[var(--color-sidebar-text-muted)] hover:text-[var(--color-primary)] hover:bg-[var(--color-sidebar-surface)] transition-colors"
            title={isFullscreen ? "Exit fullscreen" : "Fullscreen"}
          >
            {isFullscreen ? (
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M8 3v3a2 2 0 0 1-2 2H3m18 0h-3a2 2 0 0 1-2-2V3m0 18v-3a2 2 0 0 0 2-2h3M3 16h3a2 2 0 0 0 2 2v3" />
              </svg>
            ) : (
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <polyline points="15 3 21 3 21 9" />
                <polyline points="9 21 3 21 3 15" />
                <line x1="21" y1="3" x2="14" y2="10" />
                <line x1="3" y1="21" x2="10" y2="14" />
              </svg>
            )}
          </button>

          <a
            href={terminal_url}
            target="_blank"
            rel="noopener noreferrer"
            className="p-1.5 rounded-md text-[var(--color-sidebar-text-muted)] hover:text-[var(--color-primary)] hover:bg-[var(--color-sidebar-surface)] transition-colors"
            title="Open in new tab"
          >
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
              <polyline points="15 3 21 3 21 9" />
              <line x1="10" y1="14" x2="21" y2="3" />
            </svg>
          </a>
        </div>
      </div>

      {/* iframe */}
      <iframe
        src={terminal_url}
        className="w-full flex-1 border-none bg-[var(--color-app-bg)]"
        title="Workspace terminal"
        allow="clipboard-read; clipboard-write"
      />
    </div>
  );
}