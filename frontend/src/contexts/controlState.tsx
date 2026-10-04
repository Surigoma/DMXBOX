import { createContext, useContext, useEffect, useState, type ReactNode } from "react";
import { FrontConfigContext, genBackendPath } from "../routes/__root";
import { ControlState, type TControlState } from "../types";

const StateContext = createContext<TControlState | undefined>(undefined);

export function ControlStateProvider({ children }: { children: ReactNode }) {
    const config = useContext(FrontConfigContext);
    const url = genBackendPath(config, "/api/v1/control-state/stream");
    const [state, setState] = useState<{ url: string; data: TControlState }>();
    useEffect(() => {
        let stopped = false;
        let retry: ReturnType<typeof setTimeout>;
        let connection: AbortController;
        let heartbeat: ReturnType<typeof setTimeout>;
        async function connect() {
            connection = new AbortController();
            const keepAlive = () => {
                clearTimeout(heartbeat);
                heartbeat = setTimeout(() => connection.abort(), 35000);
            };
            keepAlive();
            let reader: ReadableStreamDefaultReader<Uint8Array> | undefined;
            try {
                const response = await fetch(url, { signal: connection.signal, headers: { Accept: "application/x-ndjson" } });
                if (!response.ok || !response.body || !response.headers.get("Content-Type")?.startsWith("application/x-ndjson")) {
                    throw new Error("Control state stream unavailable");
                }
                reader = response.body.getReader();
                const decoder = new TextDecoder();
                let pending = "";
                while (!stopped) {
                    const { value, done } = await reader.read();
                    if (done) throw new Error("Control state stream disconnected");
                    keepAlive();
                    pending += decoder.decode(value, { stream: true });
                    let end: number;
                    while ((end = pending.indexOf("\n")) !== -1) {
                        const line = pending.slice(0, end).trim();
                        pending = pending.slice(end + 1);
                        if (line && !stopped) setState({ url, data: ControlState.parse(JSON.parse(line)) });
                    }
                    if (pending.length > 1024 * 1024) throw new Error("Control state frame too large");
                }
            } catch {
                if (!stopped) setState(undefined);
            } finally {
                clearTimeout(heartbeat);
                connection.abort();
                await reader?.cancel().catch(() => undefined);
                if (!stopped) retry = setTimeout(() => { void connect(); }, 3000);
            }
        }
        void connect();
        return () => {
            stopped = true;
            clearTimeout(retry);
            clearTimeout(heartbeat);
            connection.abort();
        };
    }, [url]);
    return <StateContext.Provider value={state?.url === url ? state.data : undefined}>{children}</StateContext.Provider>;
}

export function useControlState() { return useContext(StateContext); }
