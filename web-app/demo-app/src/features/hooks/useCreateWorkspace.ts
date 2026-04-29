import { useForm } from "react-hook-form";

export type TypeTimeDuration = "minutes" | "hours" | "days";
export const listTypeTimeDuration: TypeTimeDuration[] = ["minutes", "hours", "days"];

export type CreateWorkspaceData = {
  name: string;
  description?: string;
  templateId?: string;
  timeDuration: number;
  typeTimeDuration: TypeTimeDuration;
  reqCpu : number
  reqRam: number
  limitCpu : number
  limitRam : number
};

export const useCreateWorkspace = () => {
  const {
    register,
    handleSubmit,
    setValue,
    watch,
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
    setValue,
    watch
  };
};