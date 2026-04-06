export function relativeTime(iso: string): string {
  const diff = Math.floor((Date.now() - new Date(iso).getTime()) / 1000);
  if (diff < 60) return `${diff}s ago`;
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`;
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`;
  return `${Math.floor(diff / 86400)}d ago`;
}

export function formatDate(time: string) {
  if (!time) {
    return;
  }
  const date = Intl.DateTimeFormat("id-ID", {
    dateStyle: "full",
  });

  return date.format(new Date(time));
}
