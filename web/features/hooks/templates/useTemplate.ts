import { useGetTemplates } from "@/features/services/templates/api";
import { Templates } from "@/types";
import { useState } from "react";

export function useTemplate() {
  const { data } = useGetTemplates();
  const templateData = data?.data as Templates[] | undefined;

  const [selected, setSelected] = useState<Templates | null>(null);
  const [form, setForm] = useState<Partial<Templates>>({});
  const [openDialogCreate, setOpenDialogEdit] = useState<boolean>(false);
  const openEdit = (item: Templates) => {
    setSelected(item);
    setForm({ ...item });
  };

  const closeEdit = () => {
    setSelected(null);
    setForm({});
  };

  const handleChange = (key: keyof Templates, value: string | boolean) => {
    setForm((prev) => ({ ...prev, [key]: value }));
  };

  return {
    templateData,
    selected,
    form,
    openDialogCreate,
    setSelected,
    setForm,
    setOpenDialogEdit,
    openEdit,
    closeEdit,
    handleChange,
  };
}
