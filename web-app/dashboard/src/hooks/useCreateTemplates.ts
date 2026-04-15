import { useCreateTemplates } from "@/features/api";
import type {
    CreateTemplateRequest,
    TemplateAddon,
    TemplateFiles,
    TemplateVariable,
} from "@/types";
import { useFieldArray, useForm } from "react-hook-form";


export function useCreateTemplate() {
  const { mutate } = useCreateTemplates();

  /* ── Defaults ── */
  const defaultVariable = (): TemplateVariable => ({
    key: "",
    default_value: "",
    required: false,
    description: "",
  });

  const defaultAddon = (): TemplateAddon => ({
    name: "",
    image: "",
    description: "",
    default_config: {},
  });

  const defaultFiles = (): TemplateFiles => ({
    filename: "",
    sort_order: 0,
  });

  const defaultValues: CreateTemplateRequest = {
    name: "",
    description: "",
    image: "",
    category: "",
    icon: "",
    is_public: false,
    variables: [],
    addons: [],
    files: [],
  };

  /* ── Form ── */
  const form = useForm<CreateTemplateRequest>({ defaultValues });
  const variables = useFieldArray({ control: form.control, name: "variables" });
  const addons = useFieldArray({ control: form.control, name: "addons" });
  const files = useFieldArray({ control: form.control, name: "files" });

  

  const reset = () => {
    form.reset(defaultValues);
  };

  /* ── Submit ── */
  const handleSubmit = form.handleSubmit((data) => {
    mutate(data);
  });

  return {
    reset,
    form,
    handleSubmit,
    variables: {
      fields: variables.fields,
      append: () => variables.append(defaultVariable()),
      remove: variables.remove,
    },
    files: {
      fields: files.fields,
      append: () => files.append(defaultFiles()),
      remove: files.remove,
    },
    addons: {
      fields: addons.fields,
      append: () => addons.append(defaultAddon()),
      remove: addons.remove,
    },
  };
}
