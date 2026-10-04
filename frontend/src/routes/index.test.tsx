import { beforeEach, describe, expect, it } from "vitest";
import { render } from "vitest-browser-react";
import { SWRConfig } from "swr";
import { http, HttpResponse } from "msw";
import { Route } from "./index";
import { Suspense } from "react";
import { UseMockServer } from "../test/backend_helper";
import { user, UserSetup } from "../test/user_helper";

describe("Control mode", () => {
    UserSetup();
    let failure: "http" | "network" | "json" | "lost-response" | undefined;
    let marker: string | null;
    let browserOnly: boolean;
    beforeEach(() => {
        failure = undefined;
        marker = null;
        browserOnly = false;
    });
    UseMockServer(
        http.get("*/config.json", () => HttpResponse.json({ backendPort: 8080 })),
        http.get("*/api/v1/config/fade", () => HttpResponse.json({})),
        http.get("*/api/features", () => HttpResponse.json([])),
        http.get("*/api/v1/control-mode", () => HttpResponse.json({ browserOnly })),
        http.post("*/api/v1/control-mode", async ({ request }) => {
            marker = request.headers.get("X-DMXBOX-Control");
            if (failure === "http") return HttpResponse.json({}, { status: 403 });
            if (failure === "network") return HttpResponse.error();
            if (failure === "json") return HttpResponse.json({});
            browserOnly = (await request.json() as { browserOnly: boolean }).browserOnly;
            if (failure === "lost-response") return HttpResponse.error();
            return HttpResponse.json({ browserOnly });
        }),
    );
    function createPage() {
        const ControlPage = Route.options.component!;
        return render(
            <SWRConfig value={{ provider: () => new Map(), dedupingInterval: 0 }}>
                <Suspense fallback="Loading...">
                    <ControlPage />
                </Suspense>
            </SWRConfig>,
        );
    }
    it("Sends the manual control marker and can enable and disable mode", async () => {
        const { getByRole } = await createPage();
        const toggle = getByRole("switch", { name: "Browser only" });
        await expect.element(toggle).toBeEnabled();
        await user.click(toggle);
        await expect.element(toggle).toBeChecked();
        expect(marker).toBe("web-ui");
        await user.click(toggle);
        await expect.element(toggle).not.toBeChecked();
        expect(browserOnly).toBe(false);
    });
    for (const kind of ["http", "network", "json"] as const) {
        it(`Shows ${kind} failure and allows retry`, async () => {
            const { getByRole } = await createPage();
            const toggle = getByRole("switch", { name: "Browser only" });
            await expect.element(toggle).toBeEnabled();
            failure = kind;
            await user.click(toggle);
            await expect.element(getByRole("alert")).toHaveTextContent("Failed to change control mode.");
            await expect.element(toggle).toBeEnabled();
            await expect.element(toggle).not.toBeChecked();
            failure = undefined;
            await user.click(toggle);
            await expect.element(toggle).toBeChecked();
            await expect.element(getByRole("alert")).not.toBeInTheDocument();
        });
    }
    it("Reads back the applied mode when the POST response is lost", async () => {
        const { getByRole } = await createPage();
        const toggle = getByRole("switch", { name: "Browser only" });
        await expect.element(toggle).toBeEnabled();
        failure = "lost-response";
        await user.click(toggle);
        await expect.element(getByRole("alert")).toBeVisible();
        await expect.element(toggle).toBeEnabled();
        await expect.element(toggle).toBeChecked();
    });
});
