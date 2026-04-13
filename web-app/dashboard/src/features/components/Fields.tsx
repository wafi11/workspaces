export function Field({
  label,
  error,
  children,
}: {
  label: string;
  error?: string;
  children: React.ReactNode;
}) {
  return (
    <div className="flex flex-col gap-1">
      <label
        className="text-[11px]"
        style={{ color: "var(--color-sidebar-text-muted)" }}
      >
        {label}
      </label>
      {children}
      {error && (
        <span className="text-[11px]" style={{ color: "#f87171" }}>
          {error}
        </span>
      )}
    </div>
  );
}
