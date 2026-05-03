import { api } from "@/lib/api";
import type { User } from "@/types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

export function useGetProfile(){
    return useQuery({
        queryKey : ["profile"],
        queryFn : async ()  => {
            const req = await api.get<User>("/user/profile")
            console.log("data profile : ",req.data)
            return req.data
        }
    })
}

export function useUpdatePhotoProfile(){
    const queryClient = useQueryClient()
    return useMutation({
        mutationKey : ["update-photo-profile"],
        mutationFn : async (data : {avatar_base64 : string})  => {
    
            const req = await api.patch("/user",data)
            return req.data
        },
        onSuccess : ()  => {
            queryClient.invalidateQueries({queryKey : ["profile"]})

        }
    })
}

export function useUpdateProfileName(){
    const queryClient = useQueryClient()
    return useMutation({
        mutationKey : ["update-photo-name"],
        mutationFn : async (data : {name : string})  => {
            const req = await api.patch("/user",data)
            return req.data
        },
        onSuccess : ()  => {
            queryClient.invalidateQueries({queryKey : ["profile"]})
        }
    })
}