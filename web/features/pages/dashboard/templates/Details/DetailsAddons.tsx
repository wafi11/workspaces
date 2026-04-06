import { Button } from "@/components/ui/button";
import {
  EditState,
  useCreateTemplateAddOn,
  useDeleteTemplateAddOn,
  useGetTemplateAddOns,
  useUpdateTemplateAddOn,
} from "@/features/services/templates/addOnApi";
import { cn } from "@/lib/utils";
import { Pencil, Plus, Puzzle, Trash2 } from "lucide-react";
import { useState } from "react";
import { AddOnForm } from "./CreateDetailsAddon";

interface DetailsAddonsProps {
  templateId: string;
}

export function DetailsAddons({ templateId }: DetailsAddonsProps) {
  const { data } = useGetTemplateAddOns(templateId);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [isAdding, setIsAdding] = useState(false);
  const [editState, setEditState] = useState<EditState>({
    name: "",
    image: "",
    description: "",
  });
  const [createState, setCreateState] = useState<EditState>({
    name: "",
    image: "",
    description: "",
  });

  const { mutate } = useUpdateTemplateAddOn(editingId || "", templateId);
  const { mutate: deleteMutate } = useDeleteTemplateAddOn(
    editingId || "",
    templateId
  );
  const { mutate: createMutate } = useCreateTemplateAddOn(templateId);

  function startEdit(addon: NonNullable<typeof data>[number]) {
    setIsAdding(false);
    setEditingId(addon.id);
    setEditState({
      name: addon.name,
      image: addon.image,
      description: addon.description,
    });
  }

  function cancelEdit() {
    setEditingId(null);
  }

  function saveEdit() {
    mutate(editState, {
      onSuccess() {
        setEditingId(null);
      },
    });
  }

  function handleDelete(id: string) {
    setEditingId(id);
    deleteMutate();
  }

  function startAdd() {
    setEditingId(null);
    setCreateState({ name: "", image: "", description: "" });
    setIsAdding(true);
  }

  function cancelAdd() {
    setIsAdding(false);
  }

  function saveAdd() {
    createMutate(createState, {
      onSuccess() {
        setIsAdding(false);
      },
    });
  }

  return (
    <div className="space-y-2">
      {/* Empty state */}
      {!data?.length && !isAdding && (
        <div className="flex flex-col items-center justify-center py-12 gap-2 border border-dashed border-border rounded-lg">
          <Puzzle className="w-8 h-8 opacity-30 text-muted-foreground" />
          <p className="text-sm text-muted-foreground">No add-ons defined</p>
          <p className="text-xs text-muted-foreground/60">
            Add integrations or tools to extend this template
          </p>
        </div>
      )}

      {/* Existing items */}
      {data?.map((addon, index) => {
        const isEditing = editingId === addon.id;
        return (
          <div
            key={addon.id}
            className={cn(
              "rounded-lg border bg-card transition-all duration-200",
              isEditing
                ? "border-primary/50 shadow-sm shadow-primary/10"
                : "border-border hover:border-border/80"
            )}
          >
            {isEditing ? (
              <AddOnForm
                state={editState}
                setState={setEditState}
                onSave={saveEdit}
                onCancel={cancelEdit}
              />
            ) : (
              <div className="flex items-start gap-3 px-3 py-2.5">
                <span className="text-[10px] text-muted-foreground/50 font-mono w-5 text-right shrink-0 pt-0.5">
                  {String(index + 1).padStart(2, "0")}
                </span>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <p className="text-sm font-medium text-foreground truncate">
                      {addon.name}
                    </p>
                    <span className="text-[10px] font-mono text-muted-foreground/60 bg-muted px-1.5 py-0.5 rounded truncate max-w-[180px]">
                      {addon.image}
                    </span>
                  </div>
                  {addon.description && (
                    <p className="text-xs text-muted-foreground mt-0.5 truncate">
                      {addon.description}
                    </p>
                  )}
                </div>
                <div className="flex gap-1 shrink-0">
                  <Button
                    size="icon"
                    variant="ghost"
                    className="h-7 w-7 text-muted-foreground hover:text-foreground"
                    onClick={() => startEdit(addon)}
                  >
                    <Pencil className="w-3.5 h-3.5" />
                  </Button>
                  <Button
                    size="icon"
                    variant="ghost"
                    className="h-7 w-7 text-muted-foreground hover:text-destructive"
                    onClick={() => handleDelete(addon.id)}
                  >
                    <Trash2 className="w-3.5 h-3.5" />
                  </Button>
                </div>
              </div>
            )}
          </div>
        );
      })}

      {/* Add new inline form */}
      {isAdding && (
        <div className="rounded-lg border border-primary/50 bg-card shadow-sm shadow-primary/10">
          <AddOnForm
            state={createState}
            setState={setCreateState}
            onSave={saveAdd}
            onCancel={cancelAdd}
          />
        </div>
      )}

      {/* Add button */}
      {!isAdding && (
        <Button
          size="sm"
          variant="outline"
          className="w-full h-8 text-xs border-dashed"
          onClick={startAdd}
        >
          <Plus className="w-3.5 h-3.5 mr-1.5" />
          Add Add-on
        </Button>
      )}
    </div>
  );
}
