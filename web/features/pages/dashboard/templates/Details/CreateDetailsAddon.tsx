import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { EditState } from "@/features/services/templates/addOnApi";
import { Check, X } from "lucide-react";

export function AddOnForm({
  state,
  setState,
  onSave,
  onCancel,
}: {
  state: EditState;
  setState: (s: any) => void;
  onSave: () => void;
  onCancel: () => void;
}) {
  return (
    <div className="p-3 space-y-3">
      <div className="grid grid-cols-2 gap-3">
        <div className="space-y-1">
          <label className="text-[10px] uppercase tracking-widest text-muted-foreground font-medium">
            Name
          </label>
          <Input
            value={state.name}
            onChange={(e) =>
              setState((s: any) => ({ ...s, name: e.target.value }))
            }
            placeholder="e.g. Redis"
            className="h-8 text-sm"
            autoFocus
          />
        </div>
        <div className="space-y-1">
          <label className="text-[10px] uppercase tracking-widest text-muted-foreground font-medium">
            Image
          </label>
          <Input
            value={state.image}
            onChange={(e) =>
              setState((s: any) => ({ ...s, image: e.target.value }))
            }
            placeholder="e.g. redis:7-alpine"
            className="h-8 text-sm font-mono"
          />
        </div>
      </div>
      <div className="space-y-1">
        <label className="text-[10px] uppercase tracking-widest text-muted-foreground font-medium">
          Description
        </label>
        <Textarea
          value={state.description}
          onChange={(e) =>
            setState((s: any) => ({ ...s, description: e.target.value }))
          }
          placeholder="Short description of this add-on..."
          className="text-sm resize-none min-h-[60px]"
        />
      </div>
      <div className="flex gap-2 justify-end">
        <Button
          size="sm"
          variant="ghost"
          onClick={onCancel}
          className="h-7 px-3 text-xs"
        >
          <X className="w-3 h-3 mr-1" />
          Cancel
        </Button>
        <Button size="sm" onClick={onSave} className="h-7 px-3 text-xs">
          <Check className="w-3 h-3 mr-1" />
          Save
        </Button>
      </div>
    </div>
  );
}
