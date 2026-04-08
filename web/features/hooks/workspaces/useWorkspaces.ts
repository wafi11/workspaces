import { useCreateWorkspaces } from "@/features/services/workspaces/api";
import { WorkspaceRequest } from "@/types/workspaces";
import { useState } from "react";

export function useWorkspaces(){
       const {mutate}  = useCreateWorkspaces()
    const [form, setForm] = useState<WorkspaceRequest>({
        name: "",
        description: "",
        env_vars: {
            "password": ""
        }
    });

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { name, value } = e.target;
        setForm(prev => ({ ...prev, [name]: value }));
    };

    // Handler khusus untuk envVars
    const handleEnvChange = (key: string, value: string) => {
        setForm(prev => ({
            ...prev,
            env_vars: { ...prev.env_vars, [key]: value }
        }));
    };

    const handleSubmit = () => {
      mutate({
        ...form,
        name : form.name.trim()
      })
    };


    return {
        mutate,
        form,
        setForm,
        handleChange,
        handleEnvChange,
        handleSubmit
    }

}