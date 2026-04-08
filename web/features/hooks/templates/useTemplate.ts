import { useTemplates } from "@/features/services/templates/api";
import { Templates } from "@/types";
import { useState } from "react";
export const FILTERS = ["Variables", "Files", "Add-ons"] as const;

export function useTemplate() {
  const { data } = useTemplates();
  const templateData = data as Templates[] | undefined;
  const [filter, setFilter] = useState<string>(FILTERS[0]);

  const isOpenFilter = (f: string) => filter === f;
  const toggleFilter = (f: string) =>
    setFilter((prev) => (prev === f ? "" : f));
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
    toggleFilter,
    isOpenFilter,
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
