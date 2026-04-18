import { api } from "@/lib/api";
import { useQuery } from "@tanstack/react-query";

type Notifications = {
    id: string
    user_id:string
    retreived_id:string
    notification_type: string
    title: string
    message: string
    metadata : {}
    is_read: boolean
    created_at: string
}

export function useGetnotificationRetreived(){
    return useQuery({
        queryKey : ["notification-retrieved"],
        queryFn : async ()  => {
            const req = await api.get<Notifications[]>("/notifications/retreived")
            return req.data
        }
    })
}


export function useGetnotificationReceived(){
    return useQuery({
        queryKey : ["notification-received"],
        queryFn : async ()  => {
            const req = await api.get<Notifications[]>("/notifications/received")
            return req.data
        }
    })
}


export function useGetnotificationUnreadCount(){
    return useQuery({
        queryKey : ["notification-unreadcount"],
        queryFn : async ()  => {
            const req = await api.get("/notifications/unread-count")
            return req.data
        }
    })
}