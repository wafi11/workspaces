import { api } from "@/lib/api";
import { useMutation, useQuery } from "@tanstack/react-query";

export function useFindAllPat(){
    return useQuery({
        queryKey : ["pat"],
        queryFn : async ()  => {
            const req = await api.get("/auth/pat")
            return req.data
        }       
    })
}


export function  useCreatePat(){
    return useMutation({
        mutationKey : ["pat"],
        mutationFn : async ({data} : {data : {name : string, expires_at : string | null}})  => {
            const req = await api.post("/auth/pat",data)
            return req.data
        }
    })
}


export function useDeletePat(){
    return useMutation({
        mutationKey : ["pat"],
        mutationFn : async (patId : string)  => {
            const req = await api.delete(`/auth/pat/${patId}`)
            return req.data
        }
    })
}