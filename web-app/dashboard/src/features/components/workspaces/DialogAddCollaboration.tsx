import { useAddCollaboration } from "@/features/api/workspace-collaboration";
import { Dialog, DialogContent, DialogDescription, DialogTitle } from "@radix-ui/react-dialog";
import { useForm } from "react-hook-form";
import { Field } from "../Input";

interface Props {
  open: boolean
  onOpenChange: (v: boolean) => void
  wsId: string
}

export type AddCollaborationType = {
  email: string
  role: string
}

export function DialogAddCollaboration({ open, onOpenChange, wsId }: Props) {
  const { mutate, isPending } = useAddCollaboration({ wsId })

  const { register, handleSubmit } = useForm<AddCollaborationType>({
    defaultValues: {
      role: "viewer",
      email: ""
    }
  })

  const onSubmit = (data: AddCollaborationType) => {
    mutate({ data })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="bg-black border border-zinc-700 rounded-lg p-6 w-95
                                 fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-50
                                 shadow-[0_8px_32px_rgba(0,0,0,0.7)]">
        <DialogTitle className="text-xs font-mono text-zinc-200 uppercase tracking-widest mb-1">
          Add Collaboration
        </DialogTitle>
        <DialogDescription className="text-xs font-mono text-zinc-500 mb-5">
          Invite a collaborator to this workspace.
        </DialogDescription>

        <div className="space-y-4">
          <Field
            label="Email"
            placeholder="admin@gmail.com"
            {...register("email")}
          />
          <Field
            label="Role"
            placeholder="viewer, editor, admin..."
            {...register("role")}
          />
        </div>

        <div className="flex justify-end gap-2 mt-6">
          <button
            disabled={isPending}
            onClick={() => onOpenChange(false)}
            className="text-xs font-mono px-3 py-1.5 border border-zinc-700 
                       text-zinc-400 hover:text-zinc-200 rounded transition-colors"
          >
            Cancel
          </button>
          <button
            disabled={isPending}
            onClick={handleSubmit(onSubmit)}
            className="text-xs font-mono px-3 py-1.5 bg-zinc-200 text-zinc-900 
                       hover:bg-white rounded transition-colors disabled:opacity-50"
          >
            {isPending ? "Adding..." : "Add"}
          </button>
        </div>
      </DialogContent>
    </Dialog>
  )
}
