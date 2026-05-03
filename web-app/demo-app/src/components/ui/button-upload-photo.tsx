import { Camera, Upload, X } from "lucide-react";
import { useRef, useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

interface ButtonUploadPhotoProps {
  onClick: () => void;
}

export function ButtonUploadPhoto({ onClick }: ButtonUploadPhotoProps) {
  return (
    <button
      onClick={onClick}
      className="absolute inset-0 rounded-[10px] flex items-center justify-center bg-black/0 hover:bg-black/50 transition-all group"
    >
      <Camera className="size-4 text-white opacity-0 group-hover:opacity-100 transition-opacity" />
    </button>
  );
}

interface DialogUploadPhotoProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (base64: string) => void;
}

export function DialogUploadPhoto({ open, onClose, onSubmit }: DialogUploadPhotoProps) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [preview, setPreview] = useState<string | null>(null);
  const [base64, setBase64] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const ALLOWED = ["image/jpeg", "image/png", "image/webp"];
  const MAX_SIZE = 5 * 1024 * 1024; // 5MB

  function handleFile(file: File) {
    setError(null);

    if (!ALLOWED.includes(file.type)) {
      setError("Hanya jpg, png, webp yang diizinkan");
      return;
    }
    if (file.size > MAX_SIZE) {
      setError("Ukuran file maksimal 5MB");
      return;
    }

    const reader = new FileReader();
    reader.onload = (e) => {
      const result = e.target?.result as string;
      setPreview(result);
      // Strip prefix "data:image/jpeg;base64," — kirim base64 murni aja
      setBase64(result.split(",")[1]);
    };
    reader.readAsDataURL(file);
  }

  function handleDrop(e: React.DragEvent) {
    e.preventDefault();
    const file = e.dataTransfer.files[0];
    if (file) handleFile(file);
  }

  async function handleSubmit() {
    if (!base64) return;
    setLoading(true);
    try {
      await onSubmit(base64);
      handleClose();
    } catch (err) {
      setError("Upload gagal, coba lagi");
    } finally {
      setLoading(false);
    }
  }

  function handleClose() {
    setPreview(null);
    setBase64(null);
    setError(null);
    onClose();
  }

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Upload Foto Profil</DialogTitle>
        </DialogHeader>

        <div className="space-y-4">
          {/* Drop zone */}
          <div
            onClick={() => inputRef.current?.click()}
            onDrop={handleDrop}
            onDragOver={(e) => e.preventDefault()}
            className="relative border-2 border-dashed border-muted-foreground/25 rounded-lg h-48 flex flex-col items-center justify-center gap-2 cursor-pointer hover:border-muted-foreground/50 transition-colors overflow-hidden"
          >
            {preview ? (
              <>
                <img src={preview} alt="preview" className="absolute inset-0 w-full h-full object-cover" />
                <button
                  onClick={(e) => { e.stopPropagation(); setPreview(null); setBase64(null); }}
                  className="absolute top-2 right-2 bg-black/60 rounded-full p-1 hover:bg-black/80 transition-colors"
                >
                  <X className="size-3 text-white" />
                </button>
              </>
            ) : (
              <>
                <Upload className="size-8 text-muted-foreground" />
                <p className="text-sm text-muted-foreground">
                  Drag & drop atau <span className="text-primary underline">pilih file</span>
                </p>
                <p className="text-xs text-muted-foreground">JPG, PNG, WebP — maks 5MB</p>
              </>
            )}
          </div>

          <input
            ref={inputRef}
            type="file"
            accept="image/jpeg,image/png,image/webp"
            className="hidden"
            onChange={(e) => { const f = e.target.files?.[0]; if (f) handleFile(f); }}
          />

          {error && <p className="text-sm text-destructive">{error}</p>}

          <div className="flex gap-2 justify-end">
            <button
              onClick={handleClose}
              className="px-4 py-2 text-sm rounded-md border border-input hover:bg-accent transition-colors"
            >
              Batal
            </button>
            <button
              onClick={handleSubmit}
              disabled={!base64 || loading}
              className="px-4 py-2 text-sm rounded-md bg-primary text-primary-foreground hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {loading ? "Uploading..." : "Upload"}
            </button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}