import { expect, describe, it, beforeEach } from "vitest";
import { render } from "vitest-browser-react";
import { user, UserSetup } from "../test/user_helper";
import { http, HttpResponse, type DefaultBodyType } from "msw";
import { UseMockServer } from "../test/backend_helper";
import FadeControl from "./FadeControl";
import type { TDMXGroup } from "../types";

describe("FadeControl", async () => {
    interface postInterface {
        params: { [key: string]: string };
        body: DefaultBodyType;
    }
    const postData: postInterface = {
        params: {},
        body: {},
    };
    beforeEach(() => {
        postData.params = {};
        postData.body = {};
    });
    UseMockServer(
        http.post("*/api/v1/fade/*", async (r) => {
            const url = new URL(r.request.url);
            const params: { [key: string]: string } = {};
            url.searchParams.forEach((v, k) => {
                params[k] = v;
            });
            postData.params = params;
            postData.body = await r.request.json();
            return HttpResponse.json(
                {},
                {
                    status: 200,
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
