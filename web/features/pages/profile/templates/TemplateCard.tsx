import { Badge } from "@/components/ui/badge"
import { Card } from "@/components/ui/card"
import { Templates } from "@/types"
import { Box } from "lucide-react"
import { usePathname, useRouter } from "next/navigation"

export function TemplateCard({ item }: { item: Templates }) {
  const {push} = useRouter()
  const pathname = usePathname()
  return (
    <Card onClick={() => push(`${pathname}/${item.id}`)} className="group relative p-3 hover:bg-accent/50 transition-all duration-200 border-border/40 cursor-pointer">
      <div className="flex items-center gap-3">
        <div className="shrink-0 size-10 rounded-md border border-border/60 bg-background flex items-center justify-center overflow-hidden">
          {item.icon ? (
            <img
              src={item.icon}
              alt={item.name}
              className="size-7 object-contain"
            />
          ) : (
            <Box className="size-4 text-muted-foreground/50" />
          )}
        </div>
        <div className="min-w-0 flex-1">
          <div className="flex items-center justify-between gap-2">
            <h3 className="font-medium text-sm truncate group-hover:text-primary transition-colors">
              {item.name}
            </h3>
            <Badge variant="outline" className="h-4 px-1 text-[9px] font-mono uppercase border-border/20 text-muted-foreground shrink-0">
              {item.category}
            </Badge>
          </div>
          
          <div className="flex items-center justify-between mt-0.5">
            <p className="text-[11px] text-muted-foreground truncate pr-4">
              {item.description || "Ready to deploy"}
            </p>
          </div>
        </div>
      </div>
    </Card>
  )
}