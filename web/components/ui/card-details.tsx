import { LucideIcon } from "lucide-react"

interface CardDetailsProps {
    label: string
    value: string
    sub: string
    icon: LucideIcon 
}

export function CardDetails({ icon: Icon, label, sub, value }: CardDetailsProps) {
    return (
        <div
            className="bg-muted/50 border border-border/50 rounded-lg p-4 hover:bg-muted/80 transition-colors"
        >
            <div className="flex justify-between items-start mb-3">
                <p className="text-[10px] font-bold text-muted-foreground uppercase tracking-widest">
                    {label}
                </p>
                <Icon className="w-4 h-4 text-muted-foreground/70" strokeWidth={2.5} />
            </div>
            
            <div className="space-y-1">
                <p className="text-2xl font-semibold font-mono tracking-tight leading-none">
                    {value}
                </p>
                <p className="text-[11px] text-muted-foreground/80 leading-relaxed">
                    {sub}
                </p>
            </div>
        </div>
    )
}