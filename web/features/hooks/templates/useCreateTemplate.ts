import { useCreateTemplates } from "@/features/services/templates/api";
import {
  CreateTemplateRequest,
  TemplateAddon,
  TemplateFiles,
  TemplateVariable,
} from "@/types";
import { useState } from "react";
import { useForm, useFieldArray } from "react-hook-form";

export const STEPS = [
  "Basic Info",
  "Variables",
  "Addons",
  "Files",
  "Review",
] as const;
export type Step = (typeof STEPS)[number];

export function useCreateTemplate() {
  const [step, setStep] = useState<number>(0);

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

  /* ── Navigation ── */
  const isFirst = step === 0;
  const isLast = step === STEPS.length - 1;

  const next = async () => {
    // validate only fields relevant to current step before advancing
    const fieldsPerStep: (keyof CreateTemplateRequest)[][] = [
      ["name", "category", "description", "image"], // step 0
      ["variables"],
      ["addons"], // step 2
      ["files"], // step 2
      [], // step 3 review
    ];

    const valid = await form.trigger(fieldsPerStep[step] as any);
    if (valid && !isLast) setStep((s) => s + 1);
  };

  const prev = () => {
    if (!isFirst) setStep((s) => s - 1);
  };

  const goTo = (i: number) => setStep(i);

  const reset = () => {
    form.reset(defaultValues);
    setStep(0);
  };

  /* ── Submit ── */
  const handleSubmit = form.handleSubmit((data) => {
    mutate(data);
    console.log("submit", data);
  });

  return {
    step,
    stepLabel: STEPS[step],
    steps: STEPS,
    isFirst,
    isLast,
    next,
    prev,
    goTo,
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
