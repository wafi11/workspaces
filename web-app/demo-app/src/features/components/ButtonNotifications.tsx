import { faBell } from "@fortawesome/free-solid-svg-icons/faBell";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { useNavigate } from "react-router-dom";
import { useGetnotificationUnreadCount } from "../api/notifications/unread-count";

export function ButtonNotification() {
  const count = useGetnotificationUnreadCount()
  const navigate  = useNavigate()
  
//   useWorkspaceSocket((data: { type: string; message: string }) => {
//     if (data.type === "notification.unread") {
//       queryClient.setQueryData(
//         ["notification-unreadcount"],
//         (old: number) => (old ?? 0) + 1
//       );
//     }
//   });

  return (
    <button onClick={() => navigate("/notifications")} className="relative flex items-center justify-center w-8 h-8 rounded-full
                       hover:bg-sidebar-accent/20 transition-colors">
      <FontAwesomeIcon icon={faBell} className="size-4" />
      {count !== undefined && count > 0 && (
        <span className="absolute -top-1 -right-1 min-w-4 h-4 px-1
                         bg-red-500 rounded-full
                         text-[10px] font-mono text-white
                         flex items-center justify-center leading-none">
          {count > 99 ? "99+" : count}
        </span>
      )}
    </button>
  )
}