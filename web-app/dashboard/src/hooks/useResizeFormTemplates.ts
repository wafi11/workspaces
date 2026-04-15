import { useCallback, useRef } from "react"

export function useResizeFormTemplates(){
      const leftRef = useRef<HTMLDivElement>(null)
      const topRightRef = useRef<HTMLDivElement>(null)
      const topLeftRef = useRef<HTMLDivElement>(null)
      const rootRef = useRef<HTMLDivElement>(null)
      const rightRef = useRef<HTMLDivElement>(null)

      const initDrag = useCallback((
          onMove: (e: MouseEvent) => void
        ) => (e: React.MouseEvent) => {
          e.preventDefault()
          const move = (ev: MouseEvent) => onMove(ev)
          const up = () => {
            window.removeEventListener('mousemove', move)
            window.removeEventListener('mouseup', up)
          }
          window.addEventListener('mousemove', move)
          window.addEventListener('mouseup', up)
        }, [])

    
        return{
            topLeftRef,
            leftRef,
            topRightRef,
            rootRef,
            rightRef,
            initDrag
        }
        
}