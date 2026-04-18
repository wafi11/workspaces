export function StatusBadge({ status }: { status: string }) {
  const map: Record<string, { bg: string; color: string; dot: string }> = {
    running: { bg: "#0d1f0f", color: "#4ade80", dot: "#22c55e" },
    stopped: { bg: "#1a1a1a", color: "#71717a", dot: "#3f3f46" },
    paused: { bg: "#1c1400", color: "#fbbf24", dot: "#f59e0b" },
    pending: { bg: "#0d1829", color: "#60a5fa", dot: "#3b82f6" },
    error: { bg: "#1f0d0d", color: "#f87171", dot: "#ef4444" },
  };
  const s = map[status] ?? map.stopped;
  return (
    <span
      className="flex items-center w-fit gap-1.5 text-[11px] font-medium px-2 py-0.5 rounded-full"
      style={{
        background: s.bg,
        color: s.color,
        border: `1px solid ${s.dot}22`,
      }}
    >
      <span
        className="w-1.5 h-1.5 rounded-full shrink-0"
        style={{ background: s.dot }}
      />
      {status}
    </span>
  );
}
