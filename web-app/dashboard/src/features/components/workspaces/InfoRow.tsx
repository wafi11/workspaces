type InfoRowProps = {
  label: string;
  value: string;
  isLink?: boolean;
};

export function InfoRow({ label, value, isLink }: InfoRowProps) {
  return (
    <div className="bg-[#0a0a0a] px-4.5 py-3.5 flex flex-col gap-1">
      <span className="text-[10px] text-sidebar-text-muted uppercase tracking-widest font-mono">
        {label}
      </span>
      {isLink ? (
        <a
          href={value}
          target="_blank"
          rel="noopener noreferrer"
          className="text-[12px] font-mono text-primary-hover underline underline-offset-[3px] decoration-[#2e2e2e] hover:text-[#e4e4e7] hover:decoration-sidebar-text-muted transition-colors"
        >
          {value}
        </a>
      ) : (
        <span className="text-[12px] font-mono text-[#e4e4e7]">{value}</span>
      )}
    </div>
  );
}