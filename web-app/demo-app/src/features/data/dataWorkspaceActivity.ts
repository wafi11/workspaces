
export const recentWorkspaces = [
  { name: "my-go-api", template: "golang-dev", status: "running", uptime: "2h 14m" },
  { name: "k8s-operator-dev", template: "golang-dev", status: "running", uptime: "45m" },
  { name: "frontend-v2", template: "node-react", status: "starting", uptime: "just now" },
  { name: "postgres-debug", template: "db-tools", status: "stopped", uptime: "3h ago" },
  { name: "redis-inspector", template: "db-tools", status: "stopped", uptime: "yesterday" },
];

export const recentActivity = [
  { type: "start", text: "frontend-v2", sub: "started", time: "just now" },
  { type: "create", text: "frontend-v2", sub: "created", time: "2 min ago" },
  { type: "stop", text: "postgres-debug", sub: "stopped", time: "3h ago" },
  { type: "template", text: "node-react", sub: "template updated", time: "yesterday" },
  { type: "start", text: "k8s-operator-dev", sub: "started", time: "yesterday" },
];

export const statusStyle: Record<string, { dot: string; badge: string; text: string }> = {
  running: {
    dot: "#22c55e",
    badge: "background: #0d2010; border: 0.5px solid #1a3a20;",
    text: "#22c55e",
  },
  starting: {
    dot: "#f59e0b",
    badge: "background: #1c1400; border: 0.5px solid #292000;",
    text: "#f59e0b",
  },
  stopped: {
    dot: "#333",
    badge: "background: #111; border: 0.5px solid #222;",
    text: "#555",
  },
};

export const activityStyle: Record<string, { bg: string; color: string; icon: string }> = {
  start: { bg: "#0d2010", color: "#22c55e", icon: "▶" },
  stop: { bg: "#1a0e0e", color: "#e24b4a", icon: "■" },
  create: { bg: "var(--color-accent-surface)", color: "var(--color-accent-text)", icon: "+" },
  template: { bg: "#1c1400", color: "#f59e0b", icon: "⊞" },
};
