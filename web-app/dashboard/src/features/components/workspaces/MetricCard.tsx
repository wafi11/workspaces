import { Cpu, MemoryStick } from "lucide-react";
import { useEffect, useState } from "react";
import { Line, LineChart, ResponsiveContainer, Tooltip, YAxis } from "recharts";
export interface MetricPoint {
  time: string;
  cpu: number;
  memory: number;
}

export function MetricCard({
  label,
  cpu,
  memory,
}: {
  label: string;
  cpu: number;
  memory: number;
}) {
  const [history, setHistory] = useState<MetricPoint[]>([]);

  useEffect(() => {
    setHistory((prev) => {
      const next = [
        ...prev,
        {
          time: new Date().toLocaleTimeString(),
          cpu: parseFloat(cpu.toFixed(3)),
          memory,
        },
      ];
      return next.slice(-20);
    });
  }, [cpu, memory]);

  return (
    <div className="flex flex-col w-full gap-2 py-3 rounded bg-[#0a0a0a] px-4">
      <span className="text-[10px] font-mono text-[#555] uppercase tracking-widest">
        {label}
      </span>

      {/* Current values */}
      <div className="flex items-center gap-4">
        <div className="flex items-center gap-1.5">
          <Cpu className="w-3 h-3 text-blue-500" />
          <span className="text-[12px] font-mono text-[#aaa]">
            {cpu.toFixed(3)} <span className="text-[#444]">cores</span>
          </span>
        </div>
        <div className="flex items-center gap-1.5">
          <MemoryStick className="w-3 h-3 text-purple-500" />
          <span className="text-[12px] font-mono text-[#aaa]">
            {memory} <span className="text-[#444]">MB</span>
          </span>
        </div>
      </div>

      {/* CPU Chart */}
      <div className="mt-1">
        <span className="text-[9px] font-mono text-[#333] uppercase">cpu</span>
        <ResponsiveContainer width="100%" height={50}>
          <LineChart data={history}>
            <Line
              type="monotone"
              dataKey="cpu"
              stroke="#3b82f6"
              strokeWidth={1}
              dot={false}
              isAnimationActive={false}
            />
            <YAxis hide domain={["auto", "auto"]} />
            <Tooltip
              contentStyle={{
                background: "#0a0a0a",
                border: "1px solid #1a1a1a",
                fontSize: 10,
              }}
              labelStyle={{ display: "none" }}
              formatter={(val) => [`${val} cores`, "cpu"]}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>

      {/* Memory Chart */}
      <div>
        <span className="text-[9px] font-mono text-[#333] uppercase">
          memory
        </span>
        <ResponsiveContainer width="100%" height={50}>
          <LineChart data={history}>
            <Line
              type="monotone"
              dataKey="memory"
              stroke="#a855f7"
              strokeWidth={1}
              dot={false}
              isAnimationActive={false}
            />
            <YAxis hide domain={["auto", "auto"]} />
            <Tooltip
              contentStyle={{
                background: "#0a0a0a",
                border: "1px solid #1a1a1a",
                fontSize: 10,
              }}
              labelStyle={{ display: "none" }}
              formatter={(val) => [`${val} MB`, "memory"]}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}
