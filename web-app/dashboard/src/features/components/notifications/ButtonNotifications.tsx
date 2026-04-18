import { useGetnotificationUnreadCount } from "@/features/api/notifications";
import { useNavigate, useRouter } from "@tanstack/react-router";
import { Bell } from "lucide-react";

export function ButtonNotification() {
  const { data: count } = useGetnotificationUnreadCount()
  const {navigate}  = useRouter()

  return (
    <button onClick={() => navigate({to : "/notifications"})} className="relative flex items-center justify-center w-8 h-8 rounded-full
                       hover:bg-sidebar-accent/20 transition-colors">
      <Bell className="size-4 text-zinc-400" />
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