import { useState } from "react"
import { Maximize2, Minimize2, Terminal, TerminalIcon } from "lucide-react"
import { section } from "framer-motion/client"

interface SectionTerminalProps {
    url?: string
}

export function SectionTerminal({ url }: SectionTerminalProps) {
    const [isExpanded, setIsExpanded] = useState(false)

    return (
        <section className="flex flex-col gap-2">
            <div className="text-xl font-semibold items-center flex gap-2">
                <TerminalIcon className="size-4"/>
                 <h2 className="text-sm font-semibold uppercase tracking-wider text-muted-foreground">
                   Terminal Access
                </h2>
            </div>
        <div className={`flex flex-col border border-[#2a2a2a] rounded-md overflow-hidden transition-all duration-300 ${
            isExpanded ? "fixed inset-4 z-50 shadow-2xl" : "w-full h-full"
        }`}>
            {/* Header bar */}
            <div className="flex items-center justify-between px-3 py-1.5 bg-[#1a1a1a] border-b border-[#2a2a2a]">
                <div className="flex items-center gap-2">
                    <Terminal className="w-3.5 h-3.5 text-[#555]" />
                    <span className="text-[11px] text-[#555] font-mono">terminal</span>
                </div>
                <button
                    onClick={() => setIsExpanded(!isExpanded)}
                    className="p-1 rounded hover:bg-[#2a2a2a] text-[#555] hover:text-[#999] transition-colors"
                >
                    {isExpanded ? (
                        <Minimize2 className="w-3.5 h-3.5" />
                    ) : (
                        <Maximize2 className="w-3.5 h-3.5" />
                    )}
                </button>
            </div>

            {/* Terminal content */}
            <div className="flex-1 bg-[#0d0d0d] overflow-hidden">
                {url ? (
                    <iframe
                        src={`https://${url}`}
                        className="w-full h-full border-0"
                        title="workspace-terminal"
                    />
                ) : (
                    <div className="w-full h-full flex items-center justify-center flex-col gap-2 text-[#2a2a2a]">
                        <span className="w-2 h-4 bg-blue-500 animate-pulse" />
                        <span className="text-[11px] font-mono">terminal url not available</span>
                    </div>
                )}
            </div>
        </div>
            </section>
    )
}