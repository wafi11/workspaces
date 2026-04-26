import { StatCards } from '@/features/components/StatsCard'

export function SectionStatsCard() {
    return (
        <div className="flex flex-col gap-4 p-5">
            <div className="flex flex-col gap-0.5">
                <h2
                    className="text-sm font-medium"
                    style={{ color: '#e8e8e8' }}
                >
                    Good morning, Wafi
                </h2>
                <p className="text-[11px]" style={{ color: '#555' }}>
                    Here's your workspace overview
                </p>
            </div>
            <StatCards />
        </div>
    )
}
