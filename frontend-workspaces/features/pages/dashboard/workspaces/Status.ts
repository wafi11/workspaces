export function statusAccent(status: string) {
  return (
    {
      running: "#1D9E75",
      stopped: "#888780",
      starting: "#BA7517",
      error: "#E24B4A",
    }[status] ?? "#888780"
  );
}

export function statusBadge(status: string): { bg: string; text: string } {
  return (
    {
      running: {
        bg: "bg-teal-50  dark:bg-teal-900",
        text: "text-teal-700  dark:text-teal-200",
      },
      stopped: {
        bg: "bg-gray-100 dark:bg-gray-800",
        text: "text-gray-500  dark:text-gray-300",
      },
      starting: {
        bg: "bg-amber-50 dark:bg-amber-900",
        text: "text-amber-700 dark:text-amber-200",
      },
      error: {
        bg: "bg-red-50   dark:bg-red-900",
        text: "text-red-600   dark:text-red-200",
      },
    }[status] ?? { bg: "bg-gray-100", text: "text-gray-500" }
  );
}
