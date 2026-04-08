import { Button } from "@/components/ui/button";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useWorkspaces } from "@/features/hooks/workspaces/useWorkspaces";
import { cn } from "@/lib/utils";

interface DialogCreateWorkspacesProps {
    className? : string
    title : string
}

export function DialogCreateWorkspaces({className,title} : DialogCreateWorkspacesProps) {
 const {
    form,
    handleChange,
    handleEnvChange,
    handleSubmit,
 }  = useWorkspaces()
    return (
        <Dialog>
            <DialogTrigger className={cn("mt-4 bg-card p-4 font-semibold",className)}>
                {title}
            </DialogTrigger>
            
            <DialogContent className="sm:max-w-[425px]">
                <DialogHeader>
                    <DialogTitle>Create Workspace</DialogTitle>
                    <DialogDescription>
                        Isi detail workspace baru Anda di bawah ini.
                    </DialogDescription>
                </DialogHeader>

                <div className="grid gap-4 py-4">
                    {/* Input Nama */}
                    <div className="grid gap-2">
                        <Label htmlFor="name">Workspace Name</Label>
                        <Input
                            id="name"
                            name="name"
                            value={form.name}
                            onChange={handleChange}
                            placeholder="Project Alpha"
                        />
                    </div>

                    {/* Input Deskripsi */}
                    <div className="grid gap-2">
                        <Label htmlFor="description">Description</Label>
                        <Textarea
                            id="description"
                            name="description"
                            value={form.description}
                            onChange={handleChange}
                            placeholder="Deskripsi singkat workspace..."
                        />
                    </div>

                    {/* Input Env Vars (Password) */}
                    <div className="grid gap-2">
                        <Label htmlFor="password">Environment Password</Label>
                        <Input
                            id="password"
                            type="password"
                            value={form.env_vars.password}
                            onChange={(e) => handleEnvChange("password", e.target.value)}
                            placeholder="••••••••"
                        />
                    </div>
                </div>

                <DialogFooter>
                    <Button type="submit" onClick={handleSubmit}>
                        Save Workspace
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}