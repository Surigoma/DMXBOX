import { afterEach } from "vitest";
import { HttpResponse } from "msw";
import type { TControlState } from "../types";

// Exercise real fetch readers, including arbitrary byte boundaries and disconnects.
export function useStateStreamMock(snapshot: () => TControlState) {
    const cleanups = new Set<() => void>();
    afterEach(() => { for (const close of cleanups) close(); cleanups.clear(); });
    return (request: Request) => {
        let timer: ReturnType<typeof setInterval>;
        let finish: () => void;
        const body = new ReadableStream<Uint8Array>({
            start(controller) {
                let previous = "";
                const stop = () => { clearInterval(timer); cleanups.delete(finish); request.signal.removeEventListener("abort", finish); };
                finish = () => { stop(); controller.close(); };
                const send = () => {
                    try {
                        const line = JSON.stringify(snapshot()) + "\n";
                        if (line !== previous) {
                            const bytes = new TextEncoder().encode(line);
                            controller.enqueue(bytes.slice(0, 7));
                            controller.enqueue(bytes.slice(7));
                            previous = line;
                        }
                    } catch (error) { stop(); controller.error(error); }
                };
                timer = setInterval(send, 50);
                cleanups.add(finish);
                request.signal.addEventListener("abort", finish, { once: true });
                send();
            },
            cancel() { clearInterval(timer); cleanups.delete(finish); request.signal.removeEventListener("abort", finish); },
        });
        return new HttpResponse(body, { headers: { "Content-Type": "application/x-ndjson" } });
    };
}
