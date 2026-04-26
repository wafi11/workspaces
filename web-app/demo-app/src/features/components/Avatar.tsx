interface AvatarProps {
  fallback: string;
  src?: string;
  size?: "1" | "2" | "3";
  radius?: "none" | "small" | "medium" | "full";
  className?: string;
}

const sizeMap = {
  "1": "w-6 h-6 text-[10px]",
  "2": "w-8 h-8 text-xs",
  "3": "w-10 h-10 text-sm",
};

const radiusMap = {
  none: "rounded-none",
  small: "rounded",
  medium: "rounded-md",
  full: "rounded-full",
};

export function Avatar({
  fallback,
  src,
  size = "2",
  radius = "full",
  className = "",
}: AvatarProps) {
  return (
    <div
      className={[
        "flex items-center justify-center font-medium shrink-0 select-none overflow-hidden",
        sizeMap[size],
        radiusMap[radius],
        className,
      ].join(" ")}
      style={{
        background: "var(--color-sidebar-surface)",
        color: "var(--color-sidebar-text-active)",
        border: "1px solid var(--color-sidebar-border)",
      }}
    >
      {src ? (
        <img src={src} alt={fallback} className="w-full h-full object-cover" />
      ) : (
        <span>{fallback}</span>
      )}
    </div>
  );
}