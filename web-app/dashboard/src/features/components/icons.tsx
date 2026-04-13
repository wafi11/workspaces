export const Ico = ({ d, d2 }: { d: string; d2?: string }) => (
  <svg
    width="16"
    height="16"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth={1.75}
    strokeLinecap="round"
    strokeLinejoin="round"
    className="shrink-0"
  >
    <path d={d} />
    {d2 && <path d={d2} />}
  </svg>
);
