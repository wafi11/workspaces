"use client"

import { Badge } from "@/components/ui/badge"
import { Card } from "@/components/ui/card"
import { Workspaces } from "@/types/workspaces"
import { Calendar, Layers, Terminal } from "lucide-react"
import { DialogCreateWorkspaces } from "../dashboard/workspaces/DialogCreateWorkspaces"
import Link from "next/link"
import { useRouter } from "next/navigation"



interface SectionWorkspacesProps {
    data?: Workspaces[]
}

export function SectionWorkspaces({ data }: SectionWorkspacesProps) {
    const {push}  = useRouter()
    if (!data || data.length === 0) {
        return (
            <div className="p-8 text-center border-2 border-dashed rounded-lg">
                <p className="text-muted-foreground text-sm">No active workspaces found.</p>
                <DialogCreateWorkspaces title="Create Your Workspaces" />
            </div>
        )
    }

    return (
        <section className="space-y-4 mt-8">
            <div className="flex justify-between items-center gap-2">
                <div className="flex items-center gap-2">

                <Layers className="w-5 h-5 text-primary" />
                <h2 className="text-sm font-semibold uppercase tracking-wider text-muted-foreground">
                    Your Workspaces
                </h2>
                <Badge variant="secondary" className="ml-2 font-mono text-[10px]">
                    {data.length} Instances
                </Badge>
                </div>
            <DialogCreateWorkspaces className="mt-0 p-2 text-sm" title="Create" />

            </div>

            <div className="grid grid-cols-1 gap-3">
                {data.map((ws) => (
                    <Card key={ws.id} onClick={() => push(`/profile/workspaces/${ws.id}`)} className="overflow-hidden  bg-muted/30 border-border/50 hover:bg-muted/50 transition-colors">
                        <div className="px-4 flex flex-col md:flex-row md:items-center justify-between gap-4">
                            {/* Info Utama */}
                            <div className="flex items-start gap-3">
                                <div className="p-2 bg-background rounded-md border border-border/50">
                                    <Terminal className="w-4 h-4 text-primary" />
                                </div>
                                <div>
                                    <h3 className="text-sm leading-none mb-1.5">{ws.name}</h3>                                        
                                        <span className="flex items-center gap-1">
                                            <Calendar className="w-3 h-3" />
                                            {new Date(ws.created_at).toLocaleDateString()}
                                        </span>
                                </div>
                            </div>

                            {/* Status & Action */}
                            {/* <div className="flex items-center justify-between md:justify-end gap-4 w-full md:w-auto"> */}
                                <div className="flex items-center gap-2">
                                    <div className="relative flex h-2 w-2">
                                        <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                                        <span className="relative inline-flex rounded-full h-2 w-2 bg-green-500"></span>
                                    </div>
                                    <span className="text-xs font-medium uppercase tracking-tighter text-green-600">
                                        {ws.status}
                                    </span>
                                </div>
                                
                               
                            {/* </div> */}
                        </div>
                    </Card>
                ))}
            </div>
        </section>
    )
}