import { faPlay, faScroll, faStop, faTrash } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

export type ViewMode = "list" | "grid";

const actionStyle = {
  base: {
    background: "var(--color-app-surface)",
    border: "0.5px solid var(--color-app-border)",
  },
};

export function WorkspaceActions({ status }: { status: string }) {
  return (
    <div className="flex items-center gap-1">
      {status === "running" ? (
        <button
          className="w-7 h-7 rounded flex items-center justify-center"
          style={actionStyle.base}
          title="Stop"
        >
          <FontAwesomeIcon icon={faStop} style={{ fontSize: "10px", color: "#e24b4a" }} />
        </button>
      ) : (
        <button
          className="w-7 h-7 rounded flex items-center justify-center"
          style={actionStyle.base}
          title="Start"
        >
          <FontAwesomeIcon icon={faPlay} style={{ fontSize: "10px", color: "#22c55e" }} />
        </button>
      )}
      <button
        className="w-7 h-7 rounded flex items-center justify-center"
        style={actionStyle.base}
        title="Logs"
      >
        <FontAwesomeIcon icon={faScroll} style={{ fontSize: "10px", color: "var(--color-sidebar-text)" }} />
      </button>
      <button
        className="w-7 h-7 rounded flex items-center justify-center"
        style={actionStyle.base}
        title="Delete"
      >
        <FontAwesomeIcon icon={faTrash} style={{ fontSize: "10px", color: "#e24b4a" }} />
      </button>
    </div>
  );
}
