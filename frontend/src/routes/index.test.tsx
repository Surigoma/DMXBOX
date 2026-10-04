import { beforeEach, describe, expect, it } from "vitest";
import { render } from "vitest-browser-react";
import { SWRConfig } from "swr";
import { http, HttpResponse } from "msw";
import { Route } from "./index";
import { Suspense } from "react";
import { UseMockServer } from "../test/backend_helper";
import type { TDMXGroupMap, TFadeState } from "../types";
import { useStateStreamMock } from "../test/stateStream";

describe("Control page", () => {
    let groups: TDMXGroupMap;
    let fadeStates: Record<string, TFadeState>;
    let fadeStateFailure: boolean;
    let features: string[];
    let isMute: boolean | null;
    beforeEach(() => {
        groups = {};
        fadeStates = {};
        fadeStateFailure = false;
        features = [];
        isMute = null;
    });
    const stream = useStateStreamMock(() => {
        if (fadeStateFailure) throw new Error("Disconnected");
        return { fade: fadeStates, mute: { isMute } };
    });
    let streamRequests: number;
    beforeEach(() => { streamRequests = 0; });
    UseMockServer(
        http.get("*/config.json", () => HttpResponse.json({ backendPort: 8080 })),
        http.get("*/api/v1/config/fade", () => HttpResponse.json(groups)),
        http.get("*/api/features", () => HttpResponse.json(features)),
        http.get("*/api/v1/control-mode", () => HttpResponse.json({ browserOnly: true })),
        http.get("*/api/v1/control-state/stream", ({ request }) => { streamRequests++; return stream(request); }),
    );
    function createPage() {
        const ControlPage = Route.options.component!;
        return render(
            <SWRConfig value={{ provider: () => new Map(), dedupingInterval: 0 }}>
                <Suspense fallback="Loading..."><ControlPage /></Suspense>
            </SWRConfig>,
        );
    }
    it("Keeps Browser only out of the control page", async () => {
        const { getByRole, getByText } = await createPage();
        await expect.element(getByText("Control", { exact: true })).toBeVisible();
        expect(getByRole("switch", { name: "Browser only" }).elements()).toHaveLength(0);
        await expect.element(getByRole("switch", { name: "CUT" })).toBeVisible();
        await expect.element(getByText("Browser only: ON")).toBeVisible();
    });
    it("Streams output changes made by another client over one connection", async () => {
        features = ["osc"];
        isMute = true;
        groups = { stage: { name: "Stage", devices: [{ model: "dimmer", channel: 1, max: [255] }] } };
        fadeStates = { stage: { level: 0, state: "idle", isIn: false } };
        const { getByRole, getByText } = await createPage();
        await expect.element(getByText("Idle · Output 0%")).toBeVisible();
        await expect.element(getByText("Muted", { exact: true })).toBeVisible();
        fadeStates = { stage: { level: 0.7, state: "fading", isIn: false } };
        isMute = false;
        await expect.element(getByRole("progressbar", { name: "Stage output level" })).toHaveAttribute("aria-valuenow", "70");
        await expect.element(getByText("Fading out · Output 70%")).toBeVisible();
        await expect.element(getByText("Unmuted", { exact: true })).toBeVisible();
        expect(streamRequests).toBe(1);

    });
    it("Keeps control buttons available when status retrieval fails", async () => {
        groups = { stage: { name: "Stage", devices: [] } };
        fadeStateFailure = true;
        const { getByRole, getByText } = await createPage();
        await expect.element(getByText("State unavailable")).toBeVisible();
        await expect.element(getByRole("button", { name: "Fade In" })).toBeEnabled();
        await expect.element(getByRole("progressbar")).not.toHaveAttribute("aria-valuenow");
    });
});
