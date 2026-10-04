import { useContext } from "react";
import { createMemoryHistory, createRootRoute, createRoute, createRouter, RouterProvider } from "@tanstack/react-router";
import { describe, expect, it } from "vitest";
import { render } from "vitest-browser-react";
import { http, HttpResponse } from "msw";
import useSWR, { SWRConfig } from "swr";
import { Route, FrontConfigContext, genBackendPath } from "./__root";
import { UseMockServer } from "../test/backend_helper";

describe("Frontend configuration loading", () => {
    let release: () => void;
    let requests: string[];
    UseMockServer(
        http.get("*/config.json", async () => {
            await new Promise<void>((resolve) => { release = resolve; });
            return HttpResponse.json({ backendPort: 47000 });
        }),
        http.get("*/api/probe", ({ request }) => {
            requests.push(request.url);
            return HttpResponse.json({ ready: true });
        }),
    );
    function Probe() {
        const config = useContext(FrontConfigContext);
        const { data } = useSWR(genBackendPath(config, "/api/probe"), (url) => fetch(url).then(r => r.json()));
        return <div>{data ? `Connected to ${config.backendPort}` : "Connecting..."}</div>;
    }
    it("Waits for config.json before making backend requests on the configured port", async () => {
        requests = [];
        const root = createRootRoute({ component: Route.options.component });
        const probe = createRoute({ getParentRoute: () => root, path: "/", component: Probe });
        const router = createRouter({ routeTree: root.addChildren([probe]), history: createMemoryHistory({ initialEntries: ["/"] }) });
        const { getByText } = await render(<SWRConfig value={{ provider: () => new Map() }}><RouterProvider router={router} /></SWRConfig>);
        await expect.element(getByText("Loading configuration...")).toBeVisible();
        await expect.poll(() => typeof release).toBe("function");
        expect(requests).toEqual([]);
        release!();
        await expect.element(getByText("Connected to 47000")).toBeVisible();
        expect(requests).toHaveLength(1);
        expect(new URL(requests[0]).port).toBe("47000");
    });
});
