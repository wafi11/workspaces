import { useEffect, useState } from "react"

export function useCountdown(expiresAt: string) {
    const [timeLeft, setTimeLeft] = useState("")

    useEffect(() => {
        const calc = () => {
            const diff = new Date(expiresAt).getTime() - Date.now()

            if (diff <= 0) {
                setTimeLeft("expired")
                return
            }

            const h = Math.floor(diff / 1000 / 60 / 60)
            const m = Math.floor((diff / 1000 / 60) % 60)
            const s = Math.floor((diff / 1000) % 60)

            setTimeLeft(
                [
                    h > 0 ? `${h}h` : null,
                    m > 0 ? `${m}m` : null,
                    `${s}s`,
                ]
                    .filter(Boolean)
                    .join(" ") + " remaining"
            )
        }

        calc()
        const interval = setInterval(calc, 1000)
        return () => clearInterval(interval)
    }, [expiresAt])

    return timeLeft
}