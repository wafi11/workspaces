import { Button } from "@/components/ui/button";
import {
  EditState,
  useCreateTemplateFiles,
  useDeleteTemplateFiles,
  useGetTemplateFiles,
  useUpdateTemplateFiles,
} from "@/features/services/templates/filesApi";
import { cn } from "@/lib/utils";
import { FileCode2, Pencil, Plus, Trash2 } from "lucide-react";
import { useState } from "react";
import { FileForm } from "./CreateDetailsAddon";

interface DetailsFilesProps {
  templateId: string;
}

export function DetailsFiles({ templateId }: DetailsFilesProps) {
  const { data } = useGetTemplateFiles(templateId);

  // State untuk UI
  const [editingId, setEditingId] = useState<string | null>(null);
  const [isAdding, setIsAdding] = useState(false);

  // State untuk Data
  const [editState, setEditState] = useState<EditState>({
    filename: "",
    sort_order: 0,
  });
  const [createState, setCreateState] = useState<EditState>({
    filename: "",
    sort_order: 0,
  });

  // Hooks API
  const { mutate: updateMutate, isPending: isUpdating } =
    useUpdateTemplateFiles(editingId || "", templateId);
  const { mutate: deleteMutate } = useDeleteTemplateFiles(templateId);
  const { mutate: createMutate, isPending: isCreating } =
    useCreateTemplateFiles(templateId);

  // ── Handler Functions ──
  function startEdit(file: NonNullable<typeof data>[number]) {
    setIsAdding(false); // Tutup form add jika sedang buka
    setEditingId(file.id);
    setEditState({
      filename: file.filename,
      sort_order: file.sort_order,
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
        setCreateState({ filename: "", sort_order: 0 });
      },
    });
  }

  function handleDelete(id: string) {
    // Pastikan ID diset sebelum delete jika hook membutuhkan editingId
    setEditingId(id);
    deleteMutate(id);
  }

  return (
    <div className="space-y-2">
      {/* 1. List Data yang Sudah Ada */}
      {data?.map((file, index) => {
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
              <FileForm
                state={editState}
                setState={setEditState}
                onSave={handleSaveEdit}
                onCancel={() => setEditingId(null)}
                isLoading={isUpdating}
              />
            ) : (
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

      {/* 2. Form Tambah Baru (Inline) */}
      {isAdding && (
        <div className="rounded-lg border border-primary/50 bg-card shadow-sm shadow-primary/10">
          <FileForm
            state={createState}
            setState={setCreateState}
            onSave={handleSaveCreate}
            onCancel={() => setIsAdding(false)}
            isLoading={isCreating}
          />
        </div>
      )}

      {/* 3. Tombol Trigger Tambah Baru */}
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
          Add File
        </Button>
      )}

      {/* 4. Empty State */}
      {!data?.length && !isAdding && (
        <div className="flex flex-col items-center justify-center py-12 gap-2 border border-dashed border-border rounded-lg">
          <FileCode2 className="w-8 h-8 opacity-30 text-muted-foreground" />
          <p className="text-sm text-muted-foreground">No files defined</p>
          <p className="text-xs text-muted-foreground/60">
            Add files to configure the template startup sequence
          </p>
        </div>
      )}
    </div>
  );
}
