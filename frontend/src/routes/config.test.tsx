import { beforeEach, describe, expect, it } from "vitest";
import { render } from "vitest-browser-react";
import { page } from "vitest/browser";
import { SWRConfig } from "swr";
import { http, HttpResponse } from "msw";
import { createRouter, createMemoryHistory, RouterProvider } from "@tanstack/react-router";
import { routeTree } from "../routeTree.gen";
import { Suspense } from "react";
import { UseMockServer } from "../test/backend_helper";
import { user, UserSetup } from "../test/user_helper";
import { DefaultConfig } from "../types";

describe("Configuration control mode", () => {
    UserSetup();
    let failure: "http" | "network" | "json" | "lost-response" | undefined;
    let marker: string | null;
    let browserOnly: boolean;
    let loadFailure: boolean;
    let configSaves: number;
    beforeEach(() => {
        failure = undefined;
        marker = null;
        browserOnly = false;
        loadFailure = false;
        configSaves = 0;
    });
    UseMockServer(
        http.get("*/config.json", () => HttpResponse.json({ backendPort: 8080 })),
        http.get("*/api/v1/config/all", () => {
            const config = DefaultConfig();
            config.output.osc.format = "/mute/{}";
            return HttpResponse.json(config);
        }),
        http.post("*/api/v1/config/save", () => {
            configSaves++;
            return HttpResponse.json({ result: "OK" });
        }),
        http.get("*/api/v1/control-mode", () => loadFailure ? HttpResponse.json({}, { status: 503 }) : HttpResponse.json({ browserOnly })),
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
    async function createPage(expand = true) {
        const router = createRouter({ routeTree, history: createMemoryHistory({ initialEntries: ["/config"] }) });
        const result = await render(
            <SWRConfig value={{ provider: () => new Map(), dedupingInterval: 0 }}>
                <Suspense fallback="Loading..."><RouterProvider router={router} /></Suspense>
            </SWRConfig>,
        );
        await expect.element(result.getByRole("button", { name: /^Control mode/ })).toBeVisible();
        if (expand) await user.click(result.getByRole("button", { name: /^Control mode/ }));
        return result;
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
    it("Keeps Browser only collapsed and separate from the saved configuration", async () => {
        const { getByRole, getByText } = await createPage(false);
        await expect.element(getByRole("button", { name: /^Control mode/ })).toHaveAttribute("aria-expanded", "false");
        await user.click(getByRole("button", { name: /^Control mode/ }));
        const toggle = getByRole("switch", { name: "Browser only" });
        await expect.element(toggle).toBeVisible();
        await expect.element(getByText(/Changes take effect immediately/)).toBeVisible();
        expect(toggle.element().closest("form")).toBeNull();
        await user.click(toggle);
        await expect.element(toggle).toBeChecked();
        expect(configSaves).toBe(0);
        await user.click(getByRole("button", { name: "Update", exact: true }));
        await expect.element(getByRole("alert")).toHaveTextContent("Success");
        expect(configSaves).toBe(1);
        await expect.element(toggle).toBeChecked();
    });
    it("Allows configuration use while control mode retrieval fails", async () => {
        loadFailure = true;
        const { getByRole } = await createPage();
        await expect.element(getByRole("alert")).toHaveTextContent("Failed to load control mode.");
        await expect.element(getByRole("switch", { name: "Browser only" })).toBeDisabled();
        await expect.element(getByRole("button", { name: "Update", exact: true })).toBeEnabled();
        loadFailure = false;
        await user.click(getByRole("button", { name: "Retry" }));
        await expect.element(getByRole("switch", { name: "Browser only" })).toBeEnabled();
        await expect.element(getByRole("alert")).not.toBeInTheDocument();
    });
    it("Shows the active mode and explanation on a narrow screen", async () => {
        await page.viewport(360, 640);
        try {
            browserOnly = true;
            const { getByRole, getByTestId, getByText } = await createPage();
            await expect.element(getByRole("switch", { name: "Browser only" })).toBeChecked();
            await expect.element(getByText(/Changes take effect immediately/)).toBeVisible();
            const panel = getByTestId("ControlModeSettings").element().getBoundingClientRect();
            expect(panel.width).toBeLessThanOrEqual(360);
            await getByTestId("ControlModeSettings").screenshot();
        } finally {
            await page.viewport(800, 600);
        }
    });
});
