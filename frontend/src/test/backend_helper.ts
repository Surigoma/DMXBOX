import { beforeAll, afterAll, afterEach } from "vitest";
import { setupWorker } from "msw/browser";
import { HttpHandler } from "msw";

// Tests run in Vitest Browser Mode, so mock the network with the Service Worker
// (`msw/native` was removed in msw v3). The worker script is served by the
// `msw/vite` plugin in "worker-only" mode under Vitest (see vite.config.ts).
export function UseMockServer(...handers: HttpHandler[]) {
    const mockWorker = setupWorker(...handers);
    beforeAll(async () => {
        await mockWorker.start({
            onUnhandledFrame: "error",
            quiet: true,
        });
    });
    afterEach(() => {
        mockWorker.resetHandlers();
    });
    afterAll(async () => {
        await mockWorker.stop();
    });
}
