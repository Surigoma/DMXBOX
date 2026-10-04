import { expect, describe, it, beforeEach, vi, afterEach } from "vitest";
import { render } from "vitest-browser-react";
import { user, UserSetup } from "../test/user_helper";
import { http, HttpResponse } from "msw";
import { UseMockServer } from "../test/backend_helper";
import MuteControl from "./MuteControl";
import { page } from "vitest/browser";
import { ControlStateProvider } from "../contexts/controlState";
import { useStateStreamMock } from "../test/stateStream";

describe("MuteControl", async () => {
    interface postInterface {
        params: { [key: string]: string };
    }
    const postData: postInterface = {
        params: {},
    };
    let isMute: boolean | null = null;
    let stateStatus = 200;
    let postStatus = 200;
    beforeEach(() => {
        postData.params = {};
        isMute = null;
        stateStatus = 200;
        postStatus = 200;
    });
    const stream = useStateStreamMock(() => {
        if (stateStatus !== 200) throw new Error("Disconnected");
        return { fade: {}, mute: { isMute } };
    });
    afterEach(() => vi.useRealTimers());
    UseMockServer(
        http.get("*/api/v1/control-state/stream", ({ request }) => stream(request)),
        http.post("*/api/v1/mute", async (r) => {
            expect(r.request.headers.get("X-DMXBOX-Control")).toBe("web-ui");
            const url = new URL(r.request.url);
            const params: { [key: string]: string } = {};
            url.searchParams.forEach((v, k) => {
                params[k] = v;
            });
            postData.params = params;
            if (postStatus === 200) isMute = params["isMute"] === "true";
            return HttpResponse.json(
                {},
                {
                    status: postStatus,
                },
            );
        }),
    );
    UserSetup();
    function CreateTestComponent() {
        return render(<ControlStateProvider><MuteControl /></ControlStateProvider>);
    }
    it("Shown", async () => {
        const { getByTestId } = await CreateTestComponent();
        const ctrl = getByTestId("MuteControl");
        await expect.element(ctrl).toBeVisible();
    });
    describe("Components", async () => {
        it("Shows the last sent state and highlights one button at a time", async () => {
            const { getByRole } = await CreateTestComponent();
            await expect.element(getByRole("status")).toHaveTextContent("State unavailable");
            const mute = getByRole("button", { name: "Mute", exact: true });
            const unmute = getByRole("button", { name: "Unmute", exact: true });
            await user.click(mute);
            await expect.element(getByRole("status")).toHaveTextContent("Muted");
            await expect.element(mute).toHaveAttribute("aria-pressed", "true");
            await expect.element(unmute).toHaveAttribute("aria-pressed", "false");
            await user.click(unmute);
            await expect.element(getByRole("status")).toHaveTextContent("Unmuted");
            await expect.element(unmute).toHaveAttribute("aria-pressed", "true");
            await expect.element(mute).toHaveAttribute("aria-pressed", "false");
        });
        it("Follows external commands and recovers after status retrieval fails", async () => {
            isMute = true;
            const { getByRole } = await CreateTestComponent();
            await expect.element(getByRole("status")).toHaveTextContent("Muted");
            vi.useFakeTimers({ toFake: ["setTimeout", "clearTimeout"] });
            stateStatus = 503;
            await expect.element(getByRole("status")).toHaveTextContent("State unavailable");
            await expect.element(getByRole("button", { name: "Mute", exact: true })).toBeEnabled();
            isMute = false;
            stateStatus = 200;
            await vi.advanceTimersByTimeAsync(3000);
            vi.useRealTimers();
            await expect.element(getByRole("status")).toHaveTextContent("Unmuted");
        });
        it("Preserves the known state when a command is rejected", async () => {
            isMute = false;
            postStatus = 500;
            const { getByRole } = await CreateTestComponent();
            await expect.element(getByRole("status")).toHaveTextContent("Unmuted");
            await user.click(getByRole("button", { name: "Mute", exact: true }));
            await expect.element(getByRole("alert")).toHaveTextContent("Failed to send mute command");
            await expect.element(getByRole("status")).toHaveTextContent("Unmuted");
        });
        it("Keeps the status above both buttons at mobile width", async () => {
            await page.viewport(360, 500);
            try {
                isMute = true;
                const { getByRole, getByTestId } = await render(<ControlStateProvider><div style={{ width: 320 }}><MuteControl /></div></ControlStateProvider>);
                await expect.element(getByRole("status")).toHaveTextContent("Muted");
                expect(getByRole("status").element().getBoundingClientRect().bottom).toBeLessThan(getByRole("button", { name: "Mute", exact: true }).element().getBoundingClientRect().top);
                await getByTestId("MuteControl").screenshot();
            } finally {
                await page.viewport(800, 600);
            }
        });
        it("Mute buttons", async () => {
            const { getByRole } = await CreateTestComponent();
            const mute = getByRole("button", { name: "Mute", exact: true });
            const unmute = getByRole("button", { name: "Unmute", exact: true });
            await expect.element(mute).toBeVisible();
            await expect.element(unmute).toBeVisible();
        });
    });
    describe("User Action", () => {
        it("Fade In", async () => {
            const { getByRole } = await CreateTestComponent();
            await user.click(
                getByRole("button", { name: "Mute", exact: true }),
            );
            console.log(postData);
            await expect(Object.keys(postData.params)).toStrictEqual([
                "isMute",
            ]);
            await expect(postData.params["isMute"]).toBe("true");
        });
        it("Fade Out", async () => {
            const { getByRole } = await CreateTestComponent();
            await user.click(
                getByRole("button", { name: "Unmute", exact: true }),
            );
            console.log(postData);
            await expect(Object.keys(postData.params)).toStrictEqual([
                "isMute",
            ]);
            await expect(postData.params["isMute"]).toBe("false");
        });
    });
});
