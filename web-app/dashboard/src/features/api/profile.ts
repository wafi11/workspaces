import { api } from "@/lib/api";
import type { ProviderUser } from "@/types/auth";
import { useQuery } from "@tanstack/react-query";

export function useProfileProviders(){
    return useQuery({
        queryKey: ["profileProviders"],
        queryFn : async ()  => {
            const req = await api.get<ProviderUser[]>("/users/providers");
            return req.data
        }
    })
}