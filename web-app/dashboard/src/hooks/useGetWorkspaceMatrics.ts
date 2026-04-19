import { API_URL } from "@/constants";
import { api } from "@/lib/api";
import type { PodMetrics } from "@/types";
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


export function useMetricSocket(onMessage: (data: any) => void) {
  const wsRef = useRef<WebSocket | null>(null);
  const onMessageRef = useRef(onMessage);
  const retryRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    onMessageRef.current = onMessage;
  }, [onMessage]);

  useEffect(() => {
    let ws: WebSocket;

    const connect = () => {
    

      ws = new WebSocket(`wss://log.wfdnstore.online/ws`);

      ws.onopen = () => console.log("[ws] connected");

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          onMessageRef.current(data);
        } catch {
          console.error("[ws] failed to parse", event.data);
        }
      };

      ws.onclose = (event) => {
          console.log("[ws] closed code:", event.code, "reason:", event.reason);

        console.log("[ws] disconnected, reconnecting in 3s...");
        retryRef.current = setTimeout(connect, 3000);
      };

      ws.onerror = (err) => console.error("[ws] error", err);

      wsRef.current = ws;
    };

    retryRef.current = setTimeout(connect, 500);

    return () => {
      if (retryRef.current) clearTimeout(retryRef.current);
      ws?.close();
    };
  }, []);

  return wsRef;
}
