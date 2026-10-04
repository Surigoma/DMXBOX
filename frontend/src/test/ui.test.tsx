import { beforeEach, describe, expect, it } from "vitest";
import { render } from "vitest-browser-react";
import { page } from "vitest/browser";
import { createRouter, createMemoryHistory, RouterProvider } from "@tanstack/react-router";
import { createTheme, ThemeProvider, CssBaseline } from "@mui/material";
import { SWRConfig } from "swr";
import { http, HttpResponse } from "msw";
import { routeTree } from "../routeTree.gen";
import { DefaultConfig, type TConfig } from "../types";
import { UseMockServer } from "./backend_helper";
import { user, UserSetup } from "./user_helper";
import { useStateStreamMock } from "./stateStream";

describe("UI navigation and configuration", () => {
    UserSetup();
    let config: TConfig;
    let browserOnly: boolean;
    let saveFailure: boolean;
    beforeEach(() => {
        browserOnly = false;
        saveFailure = false;
        config = DefaultConfig();
        config.output.osc.format = "/mute/{}";
        config.dmx.groups = { stage: { name: "舞台照明メインステージ・ウォームホワイト", devices: [{ model: "dimmer", channel: 1, max: [255] }] } };
    });
    const stream = useStateStreamMock(() => ({ fade: { stage: { level: 0.5, state: "fading", isIn: true } }, mute: { isMute: null } }));
    UseMockServer(
        http.get("*/config.json", () => HttpResponse.json({ backendPort: 8080 })),
        http.get("*/api/v1/config/all", () => HttpResponse.json(config)),
        http.post("*/api/v1/config/save", async ({ request }) => {
            if (saveFailure) return HttpResponse.error();
            config = await request.json() as TConfig;
            return HttpResponse.json({ result: "OK" });
        }),
        http.get("*/api/features", () => HttpResponse.json([])),
        http.get("*/api/v1/config/fade", () => HttpResponse.json(config.dmx.groups)),
        http.get("*/api/v1/control-state/stream", ({ request }) => stream(request)),
        http.get("*/api/v1/control-mode", () => HttpResponse.json({ browserOnly })),
        http.post("*/api/v1/control-mode", async ({ request }) => {
            browserOnly = (await request.json() as { browserOnly: boolean }).browserOnly;
            return HttpResponse.json({ browserOnly });
        }),
    );
    async function open(path = "/config", width = 1000) {
        await page.viewport(width, 800);
        const router = createRouter({ routeTree, history: createMemoryHistory({ initialEntries: [path] }) });
        const ui = await render(
            <ThemeProvider theme={createTheme({ palette: { mode: "dark" }, components: { MuiButtonBase: { defaultProps: { disableRipple: true } }, MuiDialog: { defaultProps: { transitionDuration: 0 } }, MuiAccordion: { defaultProps: { slotProps: { transition: { timeout: 0 } } } }, MuiMenu: { defaultProps: { transitionDuration: 0 } } } })}>
                <CssBaseline />
                <SWRConfig value={{ provider: () => new Map(), dedupingInterval: 0 }}>
                    <RouterProvider router={router} />
                </SWRConfig>
            </ThemeProvider>,
        );
        await expect.element(ui.getByText(path === "/config" ? "Configuration" : "Fading in · Output 50%", { exact: true })).toBeVisible();
        return ui;
    }
    async function editDevices() {
        const ui = await open();
        await user.click(ui.getByRole("button", { name: "Devices", exact: true }));
        await user.click(ui.getByTestId("OpDimmer"));
        await expect.element(ui.getByText("Unsaved changes. Press Update to save.")).toBeVisible();
        return ui;
    }
    it("Closes the mobile menu after selecting a page", async () => {
        browserOnly = true;
        const ui = await open("/", 360);
        await expect.element(ui.getByText("Browser only: ON")).toBeVisible();
        await page.screenshot();
        await user.click(ui.getByRole("button", { name: "menu" }));
        await user.click(ui.getByRole("menuitem", { name: "Config", exact: true }));
        await expect.element(ui.getByText("Configuration", { exact: true })).toBeVisible();
        await expect.element(ui.getByRole("menu")).not.toBeInTheDocument();
        await page.screenshot();
    });
    it("Fits devices and named actions within a narrow screen", async () => {
        config.dmx.groups.stage.devices = [{ model: "wclight", channel: 1, max: [255, 128, 0] }];
        const ui = await open("/config", 360);
        await user.click(ui.getByRole("button", { name: "Devices", exact: true }));
        await expect.poll(() => ui.getByTestId("DMXDevice").element().closest<HTMLElement>(".MuiCollapse-root")?.style.height).toBe("auto");
        const card = ui.getByTestId("DMXDevice").element().getBoundingClientRect();
        expect(card.right).toBeLessThanOrEqual(360);
        expect(document.documentElement.scrollWidth).toBeLessThanOrEqual(360);
        await expect.element(ui.getByRole("button", { name: "Delete group 舞台照明メインステージ・ウォームホワイト" })).toBeVisible();
        await expect.element(ui.getByRole("button", { name: "Delete device 1" })).toBeVisible();
        await ui.getByTestId("DMXDevice").screenshot();
    });
    it("Shows Browser only while the settings panel is collapsed and on the control page", async () => {
        browserOnly = true;
        const ui = await open();
        await expect.element(ui.getByRole("button", { name: "Control mode Browser only: ON" })).toHaveAttribute("aria-expanded", "false");
        await page.screenshot();
        await user.click(ui.getByRole("button", { name: "Control", exact: true }));
        await expect.element(ui.getByText("Browser only: ON")).toBeVisible();
        expect(ui.getByRole("switch", { name: "Browser only", exact: true }).elements()).toHaveLength(0);
        await page.screenshot();
    });
    it("Keeps unsaved device edits when navigation is cancelled", async () => {
        const ui = await editDevices();
        await user.click(ui.getByRole("button", { name: "Control", exact: true }));
        await expect.element(ui.getByRole("dialog", { name: "Leave configuration?" })).toBeVisible();
        await user.click(ui.getByRole("button", { name: "Stay", exact: true }));
        expect(Number(ui.getByRole("slider").first().element().getAttribute("aria-valuenow"))).toBeLessThan(255);
        await expect.element(ui.getByText("Unsaved changes. Press Update to save.")).toBeVisible();
    });
    it("Leaves only after explicitly discarding unsaved edits", async () => {
        const ui = await editDevices();
        await user.click(ui.getByRole("button", { name: "Control", exact: true }));
        await user.click(ui.getByRole("button", { name: "Discard and leave" }));
        await expect.element(ui.getByText("Fading in · Output 50%")).toBeVisible();
        expect(config.dmx.groups.stage.devices[0].max[0]).toBe(255);
    });
    it("Clears the unsaved state only after successful saving", async () => {
        const ui = await editDevices();
        await user.click(ui.getByRole("button", { name: "Update", exact: true }));
        await expect.element(ui.getByText("Unsaved changes. Press Update to save.")).not.toBeInTheDocument();
        expect(config.dmx.groups.stage.devices[0].max[0]).toBeLessThan(255);
    });
    it("Keeps the edit and leave protection when saving fails", async () => {
        const ui = await editDevices();
        saveFailure = true;
        await user.click(ui.getByRole("button", { name: "Update", exact: true }));
        await expect.element(ui.getByRole("alert").filter({ hasText: "Failed to send configuration." })).toBeVisible();
        await expect.element(ui.getByText("Unsaved changes. Press Update to save.")).toBeVisible();
        await user.click(ui.getByRole("button", { name: "Control", exact: true }));
        await expect.element(ui.getByRole("dialog", { name: "Leave configuration?" })).toBeVisible();
    });
});
