import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  EditState,
  useDeleteTemplateVariables,
  useGetTemplateVariables,
  useUpdateTemplateVariables,
  useCreateTemplateVariables, // Pastikan hook ini diimport
} from "@/features/services/templates";
import { KeyRound, Pencil, Plus, Save, Trash2, X } from "lucide-react";
import { useState } from "react";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { cn } from "@/lib/utils";

interface DetailsVariablesProps {
  templateId: string;
}

// ── Komponen Form Reusable ──
function VariableForm({
  state,
  setState,
  onSave,
  onCancel,
  isPending,
}: {
  state: EditState;
  setState: React.Dispatch<React.SetStateAction<EditState>>;
  onSave: () => void;
  onCancel: () => void;
  isPending?: boolean;
}) {
  return (
    <div className="flex flex-col gap-3 p-1">
      <div className="flex flex-col sm:flex-row gap-2">
        <Input
          value={state.key}
          onChange={(e) => setState((s) => ({ ...s, key: e.target.value }))}
          placeholder="VARIABLE_KEY"
          className="font-mono text-sm h-8 flex-1"
          autoFocus
        />
        <Input
          value={state.default_value}
          onChange={(e) =>
            setState((s) => ({ ...s, default_value: e.target.value }))
          }
          placeholder="Default value"
          className="text-sm h-8 flex-1"
        />
        <Input
          value={state.description}
          onChange={(e) =>
            setState((s) => ({ ...s, description: e.target.value }))
          }
          placeholder="Description"
          className="text-sm h-8 flex-1"
        />
      </div>
      <div className="flex gap-2 justify-end">
        <Button
          variant="ghost"
          size="sm"
          className="h-7 px-3 text-xs"
          onClick={onCancel}
          disabled={isPending}
        >
          <X className="w-3 h-3 mr-1" />
          Cancel
        </Button>
        <Button
          size="sm"
          className="h-7 px-3 text-xs"
          onClick={onSave}
          disabled={isPending}
        >
          <Save className="w-3 h-3 mr-1" />
          {isPending ? "Saving..." : "Save"}
        </Button>
      </div>
    </div>
  );
}

export function DetailsVariables({ templateId }: DetailsVariablesProps) {
  const { data } = useGetTemplateVariables(templateId);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [isAdding, setIsAdding] = useState(false);

  const [editState, setEditState] = useState<EditState>({
    key: "",
    default_value: "",
    description: "",
  });
  const [createState, setCreateState] = useState<EditState>({
    key: "",
    default_value: "",
    description: "",
  });

  // Hooks API
  const { mutate: updateMutate, isPending: isUpdating } =
    useUpdateTemplateVariables(editingId || "", templateId);
  const { mutate: deleteMutate } = useDeleteTemplateVariables(templateId);
  const { mutate: createMutate, isPending: isCreating } =
    useCreateTemplateVariables(templateId);

  // ── Handlers ──
  function startEdit(variable: NonNullable<typeof data>[number]) {
    setIsAdding(false);
    setEditingId(variable.id);
    setEditState({
      key: variable.key,
      default_value: variable.default_value,
      description: variable.description,
    });
  }

  function handleSaveEdit() {
    updateMutate(editState, {
      onSuccess: () => setEditingId(null),
    });
  }

  function handleSaveCreate() {
    createMutate(createState, {
      onSuccess: () => {
        setIsAdding(false);
        setCreateState({ key: "", default_value: "", description: "" });
      },
    });
  }

  function handleDelete(id: string) {
    // Set editingId sebentar agar hook deleteMutate mendapatkan ID yang benar
    setEditingId(id);
    deleteMutate(id, {
      onSettled: () => setEditingId(null),
    });
  }

  return (
    <div className="space-y-2">
      {/* 1. List Variables */}
      {data?.map((variable) => {
        const isEditing = editingId === variable.id;

        return (
          <div
            key={variable.id}
            className={cn(
              "rounded-lg border bg-card px-4 py-3 transition-all",
              isEditing
                ? "border-primary/50 bg-accent/30 shadow-sm shadow-primary/5"
                : "border-border hover:border-border/80"
            )}
          >
            {isEditing ? (
              <VariableForm
                state={editState}
                setState={setEditState}
                onSave={handleSaveEdit}
                onCancel={() => setEditingId(null)}
                isPending={isUpdating}
              />
            ) : (
              <div className="flex items-center gap-3">
                <code className="text-sm font-mono flex-1 truncate bg-muted/50 px-1.5 py-0.5 rounded">
                  {variable.key}
                </code>
                <span className="text-sm text-muted-foreground truncate max-w-[160px] hidden sm:block">
                  {variable.default_value || (
                    <span className="italic opacity-50 text-[12px]">
                      no default
                    </span>
                  )}
                </span>
                <Badge
                  variant={variable.required ? "default" : "secondary"}
                  className="shrink-0 hidden sm:flex text-[10px] h-5"
                >
                  {variable.required ? "Required" : "Optional"}
                </Badge>
                <div className="flex items-center gap-1 shrink-0">
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-7 w-7 text-muted-foreground hover:text-foreground"
                    onClick={() => startEdit(variable)}
                  >
                    <Pencil className="w-3.5 h-3.5" />
                  </Button>

                  <AlertDialog>
                    <AlertDialogTrigger className="h-7 w-7 text-muted-foreground hover:text-destructive">
                      <Trash2 className="w-3.5 h-3.5" />
                    </AlertDialogTrigger>
                    <AlertDialogContent>
                      <AlertDialogHeader>
                        <AlertDialogTitle>Delete variable?</AlertDialogTitle>
                        <AlertDialogDescription>
                          Variable{" "}
                          <code className="font-mono text-sm bg-muted px-1.5 py-0.5 rounded">
                            {variable.key}
                          </code>{" "}
                          akan dihapus permanen.
                        </AlertDialogDescription>
                      </AlertDialogHeader>
                      <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction
                          className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                          onClick={() => handleDelete(variable.id)}
                        >
                          Delete
                        </AlertDialogAction>
                      </AlertDialogFooter>
                    </AlertDialogContent>
                  </AlertDialog>
                </div>
              </div>
            )}
          </div>
        );
      })}

      {/* 2. Create Form Inline */}
      {isAdding && (
        <div className="rounded-lg border border-primary/50 bg-card px-4 py-3 shadow-sm shadow-primary/10">
          <VariableForm
            state={createState}
            setState={setCreateState}
            onSave={handleSaveCreate}
            onCancel={() => setIsAdding(false)}
            isPending={isCreating}
          />
        </div>
      )}

      {/* 3. Add Trigger Button */}
      {!isAdding && (
        <Button
          size="sm"
          variant="outline"
          className="w-full h-8 text-xs border-dashed"
          onClick={() => {
            setEditingId(null);
            setIsAdding(true);
          }}
        >
          <Plus className="w-3.5 h-3.5 mr-1.5" />
          Add Variable
        </Button>
      )}

      {/* 4. Empty State */}
      {!data?.length && !isAdding && (
        <div className="flex flex-col items-center justify-center py-10 text-muted-foreground gap-2 border border-dashed rounded-lg">
          <KeyRound className="w-8 h-8 opacity-30" />
          <p className="text-sm">No variables defined</p>
          <p className="text-xs opacity-60">
            Define environment variables for your services
          </p>
        </div>
      )}
    </div>
  );
}
