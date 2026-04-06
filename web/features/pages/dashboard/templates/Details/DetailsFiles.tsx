import {
  EditState,
  useDeleteTemplateFiles,
  useGetTemplateFiles,
  useUpdateTemplateFiles,
} from "@/features/services/templates/filesApi";
import { FileCode2, Pencil, Trash2, X, Check } from "lucide-react";
import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

interface DetailsFilesProps {
  templateId: string;
}

export function DetailsFiles({ templateId }: DetailsFilesProps) {
  const { data } = useGetTemplateFiles(templateId);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editState, setEditState] = useState<EditState>({
    filename: "",
    sort_order: 0,
  });

  const { mutate } = useUpdateTemplateFiles(editingId || "", templateId);
  const { mutate: deleteMutate } = useDeleteTemplateFiles(
    editingId || "",
    templateId
  );

  if (!data?.length) {
    return (
      <div className="flex flex-col items-center justify-center py-12 gap-2 border border-dashed border-border rounded-lg">
        <FileCode2 className="w-8 h-8 opacity-30 text-muted-foreground" />
        <p className="text-sm text-muted-foreground">No files defined</p>
        <p className="text-xs text-muted-foreground/60">
          Add files to configure the template startup sequence
        </p>
      </div>
    );
  }

  function startEdit(file: NonNullable<typeof data>[number]) {
    setEditingId(file.id);
    setEditState({
      filename: file.filename,
      sort_order: file.sort_order,
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
    // deleteMutate dipanggil setelah editingId di-set via useEffect
    // atau langsung pakai id dari parameter — tergantung hook API kamu
    deleteMutate();
  }

  return (
    <div className="space-y-2">
      {data.map((file, index) => {
        const isEditing = editingId === file.id;

        return (
          <div
            key={file.id}
            className={cn(
              "rounded-lg border bg-card transition-all duration-200",
              isEditing
                ? "border-primary/50 shadow-sm shadow-primary/10"
                : "border-border hover:border-border/80"
            )}
          >
            {isEditing ? (
              // ── Edit mode ──
              <div className="p-3 space-y-3">
                <div className="flex gap-3">
                  <div className="flex-1 space-y-1">
                    <label className="text-[10px] uppercase tracking-widest text-muted-foreground font-medium">
                      Filename
                    </label>
                    <Input
                      value={editState.filename}
                      onChange={(e) =>
                        setEditState((s) => ({
                          ...s,
                          filename: e.target.value,
                        }))
                      }
                      placeholder="e.g. init.sh"
                      className="h-8 text-sm font-mono"
                      autoFocus
                    />
                  </div>
                  <div className="w-28 space-y-1">
                    <label className="text-[10px] uppercase tracking-widest text-muted-foreground font-medium">
                      Sort Order
                    </label>
                    <Input
                      type="number"
                      value={editState.sort_order}
                      onChange={(e) =>
                        setEditState((s) => ({
                          ...s,
                          sort_order: Number(e.target.value),
                        }))
                      }
                      className="h-8 text-sm"
                    />
                  </div>
                </div>
                <div className="flex gap-2 justify-end">
                  <Button
                    size="sm"
                    variant="ghost"
                    onClick={cancelEdit}
                    className="h-7 px-3 text-xs"
                  >
                    <X className="w-3 h-3 mr-1" />
                    Cancel
                  </Button>
                  <Button
                    size="sm"
                    onClick={saveEdit}
                    className="h-7 px-3 text-xs"
                  >
                    <Check className="w-3 h-3 mr-1" />
                    Save
                  </Button>
                </div>
              </div>
            ) : (
              // ── View mode ──
              <div className="flex items-center gap-3 px-3 py-2.5">
                <span className="text-[10px] text-muted-foreground/50 font-mono w-5 text-right shrink-0">
                  {String(index + 1).padStart(2, "0")}
                </span>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-mono text-foreground truncate">
                    {file.filename}
                  </p>
                  <p className="text-[10px] text-muted-foreground">
                    sort order: {file.sort_order}
                  </p>
                </div>
                <div className="flex gap-1 shrink-0">
                  <Button
                    size="icon"
                    variant="ghost"
                    className="h-7 w-7 text-muted-foreground hover:text-foreground"
                    onClick={() => startEdit(file)}
                  >
                    <Pencil className="w-3.5 h-3.5" />
                  </Button>
                  <Button
                    size="icon"
                    variant="ghost"
                    className="h-7 w-7 text-muted-foreground hover:text-destructive"
                    onClick={() => handleDelete(file.id)}
                  >
                    <Trash2 className="w-3.5 h-3.5" />
                  </Button>
                </div>
              </div>
            )}
          </div>
        );
      })}
    </div>
  );
}
