import { useEffect, useRef } from "react";

export function useWorkspaceSocket(onMessage: (data: any) => void) {
  const wsRef = useRef<WebSocket | null>(null);
  const onMessageRef = useRef(onMessage);
  const retryRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    onMessageRef.current = onMessage;
  }, [onMessage]);

  useEffect(() => {
    let ws: WebSocket;

    const connect = () => {
      
      ws = new WebSocket(`ws://localhost:8080/ws`);

      ws.onopen = () => console.log("[ws] connected");

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          onMessageRef.current(data);
        } catch {
          console.error("[ws] failed to parse", event.data);
        }
      };

      ws.onclose = () => {
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
