import { expect, describe, it, beforeEach } from "vitest";
import { render } from "vitest-browser-react";
import { page } from "vitest/browser";
import { user, UserSetup } from "../test/user_helper";
import { http, HttpResponse } from "msw";
import { UseMockServer } from "../test/backend_helper";
import FadeControl from "./FadeControl";
import type { TDMXGroup, TFadeState } from "../types";

describe("FadeControl", async () => {
    interface postInterface {
        params: { [key: string]: string };
    }
    const postData: postInterface = {
        params: {},
    };
    let failure: "http" | "network" | undefined;
    beforeEach(() => {
        postData.params = {};
        failure = undefined;
    });
    UseMockServer(
        http.post("*/api/v1/fade/*", async (r) => {
            expect(r.request.headers.get("X-DMXBOX-Control")).toBe("web-ui");
            const url = new URL(r.request.url);
            const params: { [key: string]: string } = {};
            url.searchParams.forEach((v, k) => {
                params[k] = v;
            });
            postData.params = params;
            if (failure === "network") return HttpResponse.error();
            return HttpResponse.json(
                {},
                {
                    status: failure === "http" ? 503 : 200,
                },
            );
        }),
    );
    UserSetup();
    const defaultOption: TDMXGroup = {
        name: "test",
        devices: [
            {
                channel: 1,
                max: [255],
                model: "dimmer",
            },
        ],
    };
    function CreateTestComponent(data?: TDMXGroup, cutin?: boolean) {
        return render(
            <FadeControl
                data={data ?? defaultOption}
                name={data !== undefined ? data.name : defaultOption.name}
                showCutin={cutin ?? false}
            />,
        );
    }
    it("Shown", async () => {
        const { getByTestId } = await CreateTestComponent({
            name: "test",
            devices: [
                {
                    channel: 1,
                    max: [255],
                    model: "dimmer",
                },
            ],
        });
        const ctrl = getByTestId("FadeControl");
        await expect.element(ctrl).toBeVisible();
    });
    describe("Components", async () => {
        it("Keeps a long group name above the controls on a narrow screen", async () => {
            await page.viewport(360, 500);
            try {
                const name = "舞台照明メインステージ・ウォームホワイト";
                const { getByRole, getByText, getByTestId } = await render(
                    <div style={{ width: 320 }}>
                        <FadeControl name="test" data={{ ...defaultOption, name }} showCutin={true} state={{ level: 0.5, state: "fading", isIn: false }} />
                    </div>,
                );
                const title = getByText(name).element().getBoundingClientRect();
                const status = getByText("Fading out · Output 50%").element().getBoundingClientRect();
                const button = getByRole("button", { name: "Fade In" });
                expect(title.bottom).toBeLessThanOrEqual(status.top);
                expect(status.bottom).toBeLessThanOrEqual(button.element().getBoundingClientRect().top);
                expect(title.right).toBeLessThanOrEqual(getByTestId("FadeControl").element().getBoundingClientRect().right);
                await user.click(button);
                expect(postData.params["isIn"]).toBe("true");
                await user.click(getByRole("button", { name: "Fade Out" }));
                expect(postData.params["isIn"]).toBe("false");
                await getByTestId("FadeControl").screenshot();
            } finally {
                await page.viewport(800, 600);
            }
        });
        it("Shows the live output level and fading direction", async () => {
            const { getByRole, getByText } = await render(<FadeControl name="test" data={defaultOption} showCutin={false} state={{ level: 0.5, state: "fading", isIn: true }} />);
            await expect.element(getByRole("progressbar", { name: "test output level" })).toHaveAttribute("aria-valuenow", "50");
            await expect.element(getByText("Fading in · Output 50%")).toBeVisible();
            const button = getByRole("button", { name: "Fade In" }).element();
            const gauge = getByRole("progressbar").element();
            expect(getByRole("progressbar").elements()).toHaveLength(1);
            expect(getComputedStyle(gauge).backgroundColor).not.toBe("rgba(0, 0, 0, 0)");
            expect(gauge.getBoundingClientRect().width).toBeCloseTo(button.getBoundingClientRect().width, 0);
        });
        for (const state of [
            { level: 1, state: "idle", isIn: true },
            { level: 0, state: "idle", isIn: false },
            { level: 0.6, state: "waiting", isIn: false },
        ] satisfies TFadeState[]) {
            it(`Shows ${state.state} at ${state.level * 100}%`, async () => {
                const { getByRole, getByText } = await render(<FadeControl name="test" data={defaultOption} showCutin={true} state={state} />);
                await expect.element(getByRole("progressbar")).toHaveAttribute("aria-valuenow", String(state.level * 100));
                await expect.element(getByText(`${state.state === "waiting" ? "Waiting" : "Idle"} · Output ${state.level * 100}%`)).toBeVisible();
            });
        }
        it("Shows unavailable status without claiming a zero output", async () => {
            const { getByRole, getByText } = await CreateTestComponent();
            await expect.element(getByText("State unavailable")).toBeVisible();
            await expect.element(getByRole("progressbar")).not.toHaveAttribute("aria-valuenow");
            await expect.element(getByRole("button", { name: "Fade In" })).toBeEnabled();
        });
        it("Show Name", async () => {
            const { getByText } = await CreateTestComponent();
            const name = getByText("test");
            await expect.element(name).toBeVisible();
        });
        it("Show Buttons w/o cut in", async () => {
            const { getByRole } = await CreateTestComponent();
            const fadeIn = getByRole("button", { name: "Fade In" });
            const fadeOut = getByRole("button", { name: "Fade Out" });
            const cutIn = getByRole("button", { name: "Cut In" });
            const cutOut = getByRole("button", { name: "Cut Out" });
            await expect.element(fadeIn).toBeVisible();
            await expect.element(fadeOut).toBeVisible();
            await expect(cutIn.elements().length).toBe(0);
            await expect(cutOut.elements().length).toBe(0);
        });
        it("Show Buttons w/ cut in", async () => {
            const { getByRole } = await CreateTestComponent(undefined, true);
            const fadeIn = getByRole("button", { name: "Fade In" });
            const fadeOut = getByRole("button", { name: "Fade Out" });
            const cutIn = getByRole("button", { name: "Cut In" });
            const cutOut = getByRole("button", { name: "Cut Out" });
            await expect.element(fadeIn).toBeVisible();
            await expect.element(fadeOut).toBeVisible();
            await expect.element(fadeIn).toBeVisible();
            await expect.element(cutIn).toBeVisible();
            await expect.element(cutOut).toBeVisible();
        });
    });
    describe("User Action", () => {
        for (const kind of ["http", "network"] as const) {
            it(`Shows ${kind} failure and retries the same Cut command`, async () => {
                const { getByRole } = await CreateTestComponent(undefined, true);
                failure = kind;
                await user.click(getByRole("button", { name: "Cut Out" }));
                await expect.element(getByRole("alert")).toHaveTextContent("Failed to send Cut Out for test.");
                await expect.element(getByRole("button", { name: "Cut Out" })).toBeEnabled();
                failure = undefined;
                await user.click(getByRole("button", { name: "Retry" }));
                await expect.element(getByRole("alert")).not.toBeInTheDocument();
                expect(postData.params).toEqual({ isIn: "false", interval: "0", duration: "0" });
            });
        }
        it("Fade In", async () => {
            const { getByRole } = await CreateTestComponent();
            await user.click(getByRole("button", { name: "Fade In" }));
            console.log(postData);
            await expect(Object.keys(postData.params)).toStrictEqual(["isIn"]);
            await expect(postData.params["isIn"]).toBe("true");
        });
        it("Fade Out", async () => {
            const { getByRole } = await CreateTestComponent();
            await user.click(getByRole("button", { name: "Fade Out" }));
            console.log(postData);
            await expect(Object.keys(postData.params)).toStrictEqual(["isIn"]);
            await expect(postData.params["isIn"]).toBe("false");
        });
        it("Cut In", async () => {
            const { getByRole } = await CreateTestComponent(undefined, true);
            await user.click(getByRole("button", { name: "Cut In" }));
            console.log(postData);
            await expect(Object.keys(postData.params)).toStrictEqual([
                "isIn",
                "interval",
                "duration",
            ]);
            await expect(postData.params["isIn"]).toBe("true");
            await expect(postData.params["duration"]).toBe("0");
            await expect(postData.params["interval"]).toBe("0");
        });
        it("Cut Out", async () => {
            const { getByRole } = await CreateTestComponent(undefined, true);
            await user.click(getByRole("button", { name: "Cut Out" }));
            console.log(postData);
            await expect(Object.keys(postData.params)).toStrictEqual([
                "isIn",
                "interval",
                "duration",
            ]);
            await expect(postData.params["isIn"]).toBe("false");
            await expect(postData.params["duration"]).toBe("0");
            await expect(postData.params["interval"]).toBe("0");
        });
    });
});
