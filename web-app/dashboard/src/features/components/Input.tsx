
// Field sekarang forward ref biar kompatibel sama register()
import { forwardRef } from "react"

export const Field = forwardRef<
  HTMLInputElement,
  { label: string; placeholder: string } & React.InputHTMLAttributes<HTMLInputElement>
>(({ label, placeholder, ...props }, ref) => {
  return (
    <div>
      <label className="block text-[10px] font-mono text-zinc-500 uppercase tracking-widest mb-1.5">
        {label}
      </label>
      <input
        ref={ref}
        placeholder={placeholder}
        className="w-full bg-zinc-950 border border-zinc-700 focus:border-zinc-500
                   text-zinc-200 text-xs font-mono px-3 py-2 rounded outline-none transition-colors
                   placeholder:text-zinc-600"
        {...props}
      />
    </div>
  )
})
Field.displayName = "Field"