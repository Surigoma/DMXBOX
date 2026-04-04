import { expect, describe, it } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";
import { FormProvider, useForm } from "react-hook-form";
import { user, UserSetup } from "../../../test/user_helper";
import type { TOSCServer } from "../../../types";
import OutputOSC from "./osc";

describe("Output OSC", async () => {
    UserSetup();
    interface testForm {
        osc: TOSCServer;
    }
    const defaultValue: testForm = {
        osc: {
            ip: "127.0.0.1",
            port: 49900,
            format: "/{}/",
            channels: [1],
            inverse: false,
            type: "float",
        },
    };
    const result: testForm = JSON.parse(JSON.stringify(defaultValue));
    function TestForm(f: { callback: (v: testForm) => void }) {
        const configForm = useForm<testForm>({
            defaultValues: defaultValue,
        });
        return (
            <FormProvider {...configForm}>
                <form onSubmit={configForm.handleSubmit((v) => f.callback(v))}>
                    <OutputOSC />
                    <input type="submit" value="SUBMIT" />
                </form>
            </FormProvider>
        );
    }
    function CreateTestComponent(): Promise<RenderResult> {
        return render(
            <TestForm
                callback={(v) => {
                    console.log(v);
                    result.osc = v.osc;
                }}
            ></TestForm>,
        );
    }
    it("Shown", async () => {
        const { getByTestId } = await CreateTestComponent();
        const osc = getByTestId("OutputOSC");
        await expect.element(osc).toBeVisible();
    });
    describe("Components", async () => {
        describe("Address", async () => {
            it("Show Address", async () => {
                const { getByLabelText } = await CreateTestComponent();
                const address = getByLabelText("Address");
                await expect.element(address).toBeVisible();
            });
        });
        describe("Port", async () => {
            it("Show Port", async () => {
                const { getByLabelText } = await CreateTestComponent();
                const port = getByLabelText("Port");
                await expect.element(port).toBeVisible();
            });
            it("Can change 1 to 65535", async () => {
                const { getByLabelText } = await CreateTestComponent();
                const port = getByLabelText("Port");
                await user.click(port);
                await user.clear(port);
                await user.fill(port, "1");
                await user.keyboard("{ArrowDown}");
                await expect.element(port).toHaveValue("1");
                await user.clear(port);
                await user.fill(port, "65535");
                await user.keyboard("{ArrowUp}");
                await expect.element(port).toHaveValue("65535");
            });
        });
        describe("OSC Path format", async () => {
            it("Show OSC Path", async () => {
                const { getByLabelText } = await CreateTestComponent();
                const path = getByLabelText("OSC Path format");
                await expect.element(path).toBeVisible();
            });
        });
        describe("Sending format", async () => {
            it("Show Type", async () => {
                const { getByLabelText } = await CreateTestComponent();
                const type = getByLabelText("Sending data type");
                await expect.element(type).toBeVisible();
            });
            it("Can select float and int", async () => {
                const { getByRole, getByTestId } = await CreateTestComponent();
                const type = getByTestId("OpSendingDataType");
                const float = getByRole("option", { name: "Float" });
                const int = getByRole("option", { name: "Int" });
                await user.click(type);
                await expect.element(float).toBeVisible();
                await expect.element(int).toBeVisible();
                await user.click(int);
                await expect.element(type.getByText("Int")).toBeVisible();
            });
        });
        describe("Inverse", async () => {
            it("Show Inverse", async () => {
                const { getByLabelText } = await CreateTestComponent();
                const inv = getByLabelText("Inverse");
                await expect.element(inv).toBeVisible();
            });
            it("Can change", async () => {
                const { getByLabelText } = await CreateTestComponent();
                const inv = getByLabelText("Inverse");
                await user.click(inv);
                await expect.element(inv).toBeChecked();
            });
        });
        describe("Target OSC Channels", async () => {
            it("Show Target OSC Channels", async () => {
                const { getByLabelText } = await CreateTestComponent();
                const channels = getByLabelText("Target OSC Channels");
                await expect.element(channels).toBeVisible();
            });
            it("Show 1 to 255", async () => {
                const { getByLabelText } = await CreateTestComponent();
                const channels = getByLabelText("Target OSC Channels");
                await user.click(channels);
                Promise.all(
                    Array(256)
                        .map((_, i) => (i + 1).toString())
                        .map(async (v) => {
                            return expect
                                .element(
                                    channels.getByRole("option", {
                                        name: v,
                                        exact: true,
                                    }),
                                )
                                .toBeVisible();
                        }),
                );
            });
            it("Can add channel", async () => {
                const { getByLabelText, getByText } =
                    await CreateTestComponent();
                const channels = getByLabelText("Target OSC Channels");
                const submit = getByText("SUBMIT");
                await user.click(channels);
                await user.click(
                    channels.getByRole("option", { name: "2", exact: true }),
                );
                await user.click(submit);
                await expect(result.osc.channels).toEqual([1, 2]);
            });
            it("Can remove channel", async () => {
                const { getByLabelText, getByText } =
                    await CreateTestComponent();
                const channels = getByLabelText("Target OSC Channels");
                const submit = getByText("SUBMIT");
                await user.click(channels);
                await user.click(
                    channels.getByRole("option", { name: "1", exact: true }),
                );
                await user.click(submit);
                await expect(result.osc.channels).toEqual([]);
            });
        });
    });
    it("Can submit", async () => {
        const { getByText, getByLabelText } = await CreateTestComponent();
        const address = getByLabelText("Address");
        const submit = getByText("SUBMIT");
        await user.click(address);
        await user.clear(address);
        await user.fill(address, "test");
        await user.click(submit);
        await expect(result).toEqual({
            osc: {
                ip: "test",
                port: 49900,
                format: "/{}/",
                channels: [1],
                inverse: false,
                type: "float",
            },
        } as testForm);
    });
});
