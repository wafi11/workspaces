import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  EditState,
  useDeleteTemplateVariables,
  useGetTemplateVariables,
  useUpdateTemplateVariables,
} from "@/features/services/templates";
import { KeyRound, Pencil, Save, Trash2, X } from "lucide-react";
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

interface DetailsVariablesProps {
  templateId: string;
}

export function DetailsVariables({ templateId }: DetailsVariablesProps) {
  const { data } = useGetTemplateVariables(templateId);
  const [editingId, setEditingId] = useState<string | null>(null);
  const { mutate } = useUpdateTemplateVariables(editingId || "", templateId);
  const { mutate: deleteMutate } = useDeleteTemplateVariables(
    editingId || "",
    templateId
  );
  const [editState, setEditState] = useState<EditState>({
    key: "",
    default_value: "",
    description: "",
  });

  if (!data?.length) {
    return (
      <div className="flex flex-col items-center justify-center py-10 text-muted-foreground gap-2">
        <KeyRound className="w-8 h-8 opacity-40" />
        <p className="text-sm">No variables defined</p>
      </div>
    );
  }

  function startEdit(variable: NonNullable<typeof data>[number]) {
    setEditingId(variable.id);
    setEditState({
      key: variable.key,
      default_value: variable.default_value,
      description: variable.description,
    });
  }

  function cancelEdit() {
    setEditingId(null);
  }

  function saveEdit(id: string) {
    mutate(editState, {
      onSuccess() {
        setEditingId(null);
      },
    });
    setEditingId(null);
  }

  function handleDelete(id: string) {
    deleteMutate();
  }

  return (
    <div className="space-y-2">
      {data.map((variable) => {
        const isEditing = editingId === variable.id;

        return (
          <div
            key={variable.id}
            className={`rounded-lg border bg-card px-4 py-3 transition-colors ${
              isEditing ? "border-primary/50 bg-accent/30" : ""
            }`}
          >
            {isEditing ? (
              // — Edit mode —
              <div className="flex flex-col sm:flex-row gap-2">
                <Input
                  value={editState.key}
                  onChange={(e) =>
                    setEditState((s) => ({ ...s, key: e.target.value }))
                  }
                  placeholder="KEY"
                  className="font-mono text-sm h-8 flex-1"
                />
                <Input
                  value={editState.default_value}
                  onChange={(e) =>
                    setEditState((s) => ({
                      ...s,
                      default_value: e.target.value,
                    }))
                  }
                  placeholder="Default value"
                  className="text-sm h-8 flex-1"
                />
                <Input
                  value={editState.description}
                  onChange={(e) =>
                    setEditState((s) => ({
                      ...s,
                      description: e.target.value,
                    }))
                  }
                  placeholder="Description"
                  className="text-sm h-8 flex-1"
                />
                <div className="flex gap-1 shrink-0">
                  <Button
                    size="icon"
                    className="h-8 w-8"
                    onClick={() => saveEdit(variable.id)}
                  >
                    <Save className="w-3.5 h-3.5" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8"
                    onClick={cancelEdit}
                  >
                    <X className="w-3.5 h-3.5" />
                  </Button>
                </div>
              </div>
            ) : (
              // — View mode —
              <div className="flex items-center gap-3">
                <code className="text-sm font-mono flex-1 truncate">
                  {variable.key}
                </code>
                <span className="text-sm text-muted-foreground truncate max-w-[160px] hidden sm:block">
                  {variable.default_value || (
                    <span className="italic opacity-50">empty</span>
                  )}
                </span>
                <Badge
                  variant={variable.required ? "default" : "secondary"}
                  className="shrink-0 hidden sm:flex"
                >
                  {variable.required ? "Required" : "Optional"}
                </Badge>
                <div className="flex items-center gap-1 shrink-0">
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-7 w-7"
                    onClick={() => startEdit(variable)}
                  >
                    <Pencil className="w-3.5 h-3.5" />
                  </Button>
                  <AlertDialog>
                    <AlertDialogTrigger className="h-7 w-7 text-destructive hover:text-destructive">
                      <Trash2 className="w-3.5 h-3.5" />
                    </AlertDialogTrigger>
                    <AlertDialogContent>
                      <AlertDialogHeader>
                        <AlertDialogTitle>Delete variable?</AlertDialogTitle>
                        <AlertDialogDescription>
                          <code className="font-mono text-sm bg-muted px-1.5 py-0.5 rounded">
                            {variable.key}
                          </code>{" "}
                          will be permanently deleted and cannot be undone.
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
    </div>
  );
}
