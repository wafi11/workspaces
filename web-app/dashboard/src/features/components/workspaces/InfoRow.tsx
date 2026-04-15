type InfoRowProps = {
  label: string;
  value: string;
  isLink?: boolean;
};

export function InfoRow({ label, value, isLink }: InfoRowProps) {
  return (
    <div className="bg-[#0a0a0a] px-[18px] py-3.5 flex flex-col gap-1">
      <span className="text-[10px] text-[#3f3f46] uppercase tracking-[0.1em] font-mono">
        {label}
      </span>
      {isLink ? (
        <a
          href={value}
          target="_blank"
          rel="noopener noreferrer"
          className="text-[12px] font-mono text-[#a1a1aa] underline underline-offset-[3px] decoration-[#2e2e2e] hover:text-[#e4e4e7] hover:decoration-[#3f3f46] transition-colors"
        >
          {value}
        </a>
      ) : (
        <span className="text-[12px] font-mono text-[#e4e4e7]">{value}</span>
      )}
    </div>
  );
}