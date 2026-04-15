import { SSE_URL } from "@/constants";
import { useState, useEffect, useRef } from "react";
import { api } from "@/lib/api";
interface ContainerMetrics {
  name: string;
  cpu_cores: number;
  memory_mb: number;
}

interface PodMetrics {
  pod_name: string;
  app_name: string;
  containers: ContainerMetrics[];
  total_cpu_cores: number;
  total_memory_mb: number;
}

export function useGetMetricsWorkspaces() {
  const [metrics, setMetrics] = useState<PodMetrics[]>([]);
  const closeRef = useRef<(() => void) | null>(null);

  useEffect(() => {
    const connect = async () => {
      const close = await api.sse(`${SSE_URL}/metrics`, (data: PodMetrics[]) =>
        setMetrics(data)
      );
      closeRef.current = close;
    };

    connect();

    return () => {
      closeRef.current?.();
    };
  }, []);

  const byApp = metrics.reduce(
    (acc, pod) => {
      acc[pod.app_name] = pod;
      return acc;
    },
    {} as Record<string, PodMetrics>
  );

  return { metrics, byApp };
}
