import { useState } from "react";

export function CopyableValue({
  value, highlight, mono, copyable,
}: {
  value: string; highlight?: boolean; mono?: boolean; copyable?: boolean;
}) {
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    if (!copyable) return;
    navigator.clipboard.writeText(value).catch(() => {});
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  };

  return (
    <span
      onClick={handleCopy}
      className={`flex items-center gap-1 text-xs max-w-[260px] truncate text-right
        ${mono ? "font-mono" : ""}
        ${highlight ? "text-[#a1a1aa]" : "text-[#71717a]"}
        ${copyable ? "cursor-pointer hover:text-[#a1a1aa] transition-colors" : ""}
        ${copied ? "!text-green-400" : ""}
      `}
    >
      {copied ? "copied!" : value}
      {copyable && !copied && (
        <svg className="w-3 h-3 opacity-30 flex-shrink-0" viewBox="0 0 12 12" fill="none" stroke="currentColor" strokeWidth="1.2">
          <rect x="4" y="4" width="7" height="7" rx="1.2" />
          <path d="M2 8V2h6" />
        </svg>
      )}
    </span>
  );
}