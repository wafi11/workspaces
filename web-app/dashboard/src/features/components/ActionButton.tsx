type Variant = "default" | "warn" | "danger";

export function ActionBtn({
  label,
  variant,
  disabled,
  type,
  className,
  onClick,
}: {
  className?: string;
  disabled?: boolean;
  type?: "submit" | "button";
  label: string;
  variant: Variant;
  onClick?: () => void;
}) {
  const styles: Record<
    Variant,
    { color: string; border: string; hoverBg: string }
  > = {
    default: {
      color: "var(--color-sidebar-text-active)",
      border: "var(--color-sidebar-border)",
      hoverBg: "var(--color-sidebar-surface)",
    },
    warn: { color: "#f59e0b", border: "#292000", hoverBg: "#1c1400" },
    danger: { color: "#f87171", border: "#2a1212", hoverBg: "#1f0d0d" },
  };
  const s = styles[variant];

  return (
    <button
      type={type ?? "button"}
      disabled={disabled}
      onClick={onClick}
      className={`flex-1 py-1.5 rounded-md text-xs font-medium transition-colors ${className ?? ""}`}
      style={{
        color: disabled ? "var(--color-sidebar-text-muted)" : s.color,
        border: `1px solid ${s.border}`,
        background: "transparent",
        cursor: disabled ? "not-allowed" : "pointer",
        opacity: disabled ? 0.5 : 1,
      }}
      onMouseEnter={(e) => {
        if (!disabled) e.currentTarget.style.background = s.hoverBg;
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.background = "transparent";
      }}
    >
      {label}
    </button>
  );
}
