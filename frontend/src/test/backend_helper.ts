import { beforeAll, afterAll, afterEach } from "vitest";
import { setupServer } from "msw/native";
import { HttpHandler } from "msw";

export function UseMockServer(...handers: HttpHandler[]) {
    const mockServer = setupServer(...handers);
    beforeAll(async () => {
        await mockServer.listen({
            onUnhandledRequest: "error",
        });
    });
    afterEach(() => {
        mockServer.resetHandlers();
    });
    afterAll(() => {
        mockServer.dispose();
    });
}
