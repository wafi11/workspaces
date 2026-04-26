import { faSearch } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

interface SearchBarProps {
  placeholder?: string;
  value: string;
  onChange: (value: string) => void;
}

export function SearchBar({ placeholder, value, onChange }: SearchBarProps) {
  return (
    <div
      className="flex items-center gap-2 px-4 py-1.5 w-full rounded-md"
      style={{
        background: "var(--color-app-surface)",
        border: "0.5px solid var(--color-app-border)",
      }}
    >
      <FontAwesomeIcon
        icon={faSearch}
        style={{ color: "var(--color-sidebar-text-muted)", fontSize: "11px" }}
      />
      <input
        type="text"
        className="w-full bg-transparent text-xs focus:outline-none placeholder:text-[var(--color-sidebar-text-muted)]"
        placeholder={placeholder ?? "Search..."}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        style={{ color: "var(--color-sidebar-text-active)" }}
      />
    </div>
  );
}