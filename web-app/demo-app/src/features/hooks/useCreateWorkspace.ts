import { useForm } from "react-hook-form";

export type TypeTimeDuration = "minutes" | "hours" | "days";
export const listTypeTimeDuration: TypeTimeDuration[] = ["minutes", "hours", "days"];

export type CreateWorkspaceData = {
  name: string;
  description?: string;
  templateId?: string;
  timeDuration: number;
  typeTimeDuration: TypeTimeDuration;
};

export const useCreateWorkspace = () => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateWorkspaceData>({
    defaultValues: {
      name: "",
      description: "",
      timeDuration: 1,
      typeTimeDuration: "hours",
    },
  });

  const onSubmit = (data: CreateWorkspaceData) => {
    console.log(data);
  };

  return {
    register,
    handleSubmit: handleSubmit(onSubmit),
    errors,
  };
};