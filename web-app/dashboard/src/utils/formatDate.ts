export function formatDate(iso: string | null) {
  if (!iso) return "—";
  return new Intl.DateTimeFormat("id-ID", {
    day: "2-digit",
    month: "short",
    hour: "2-digit",
    minute: "2-digit",
  }).format(new Date(iso));
}
