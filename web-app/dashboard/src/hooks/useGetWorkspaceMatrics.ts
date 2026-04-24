import { API_URL } from "@/constants";
import { api } from "@/lib/api";
import type { PodMetrics, PodStorageMetrics } from "@/types";
import { useEffect, useRef, useState } from "react";

export function useGetMetricsWorkspaces() {
  const [metrics, setMetrics] = useState<PodMetrics[]>([]);
  const closeRef = useRef<(() => void) | null>(null);

  useEffect(() => {
    const connect = async () => {
      const close = await api.sse(`${API_URL}/metrics`, (data: PodMetrics[]) =>
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


export function useGetMetricsStorageWorkspaces() {
  const [metrics, setMetrics] = useState<PodStorageMetrics | null>(null);
  const closeRef = useRef<(() => void) | null>(null);

  useEffect(() => {
    const connect = async () => {
      const close = await api.sse(`${API_URL}/metrics/storage`, (data: PodStorageMetrics) => {
        console.log(data)
        setMetrics(data)
      }
      );
      closeRef.current = close;
    };

    connect();

    return () => {
      closeRef.current?.();
    };
  }, []);

  return { metrics };
}


