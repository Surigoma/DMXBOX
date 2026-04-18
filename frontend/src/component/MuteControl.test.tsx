import { expect, describe, it, beforeEach } from "vitest";
import { render } from "vitest-browser-react";
import { user, UserSetup } from "../test/user_helper";
import { http, HttpResponse, type DefaultBodyType } from "msw";
import { UseMockServer } from "../test/backend_helper";
import MuteControl from "./MuteControl";

describe("MuteControl", async () => {
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
        http.post("*/api/v1/mute", async (r) => {
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
    function CreateTestComponent() {
        return render(<MuteControl />);
    }
    it("Shown", async () => {
        const { getByTestId } = await CreateTestComponent();
        const ctrl = getByTestId("MuteControl");
        await expect.element(ctrl).toBeVisible();
    });
    describe("Components", async () => {
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
