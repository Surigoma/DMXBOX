import { expect, describe, it } from "vitest";
import { render } from "vitest-browser-react";
import { FormProvider, useForm } from "react-hook-form";
import Device from "./device";
import type { TDMXGroup } from "../../../types";
import { user, UserSetup } from "../../../test/helper";

describe("DMXDevice Component", async () => {
    UserSetup();
    interface testForm {
        test: TDMXGroup;
    }
    const result: testForm = {
        test: {
            devices: [
                {
                    model: "dimmer",
                    channel: 1,
                    max: [0],
                },
            ],
            name: "test",
        },
    };
    function TestForm(f: { callback: (v: testForm) => void }) {
        const configForm = useForm<testForm>({
            defaultValues: {
                test: {
                    devices: [
                        {
                            model: "dimmer",
                            channel: 1,
                            max: [0],
                        },
                    ],
                    name: "test",
                },
            },
        });
        return (
            <FormProvider {...configForm}>
                <form onSubmit={configForm.handleSubmit((v) => f.callback(v))}>
                    <Device base="test" index={0} />
                    <input type="submit" value="SUBMIT" />
                </form>
            </FormProvider>
        );
    }
    function CreateTestComponent() {
        return render(
            <TestForm
                callback={(v) => {
                    console.log(v);
                    result.test = v.test;
                }}
            ></TestForm>,
        );
    }
    it("Shown", async () => {
        const { getByTestId } = await CreateTestComponent();
        const device = getByTestId("DMXDevice");
        await expect.element(device).toBeVisible();
    });
    describe("Global config", async () => {
        it("Show Model selector", async () => {
            const { getByTestId } = await CreateTestComponent();
            const device = getByTestId("DMXDevice");
            const dimmer = getByTestId("Dimmer");
            await expect.element(device).toBeVisible();
            await expect.element(dimmer).toBeVisible();
        });
        it("Can change Model selector", async () => {
            const { getByTestId, getByText } = await CreateTestComponent();
            const device = getByTestId("DMXDevice");
            const wclight = getByTestId("WCLight");
            const selector = getByTestId("OpModelSelect");
            const wcText = getByText("White Control Light");
            await expect.element(device).toBeVisible();
            await user.click(selector);
            await user.click(wcText);
            await expect.element(wclight).toBeVisible();
        });
        it("Show Channel Select", async () => {
            const { getByTestId, getByRole } = await CreateTestComponent();
            const device = getByTestId("DMXDevice");
            const channel = getByRole("textbox");
            await expect.element(device).toBeVisible();
            await expect.element(channel).toBeVisible();
        });
        it("Can change channel between 1 to 255", async () => {
            const { getByTestId, getByRole } = await CreateTestComponent();
            const device = getByTestId("DMXDevice");
            const channel = getByRole("textbox");
            await expect.element(device).toBeVisible();
            await expect.element(channel).toBeVisible();

            await user.click(channel);
            await user.keyboard("{BackSpace}");
            await user.fill(channel, "0");
            await user.keyboard("{ArrowDown}");
            await expect.element(channel).toHaveValue("1");
            await user.fill(channel, "256");
            await user.keyboard("{ArrowUp}");
            await expect.element(channel).toHaveValue("255");
        });
    });
    it("Can submit", async () => {
        const { getByTestId, getByText } = await CreateTestComponent();
        const device = getByTestId("DMXDevice");
        const submit = getByText("SUBMIT");
        await expect.element(device).toBeVisible();
        await user.click(submit);
        await expect(result).toEqual({
            test: {
                name: "test",
                devices: [
                    {
                        model: "dimmer",
                        channel: 1,
                        max: [0],
                    },
                ],
            },
        } as testForm);
    });
    it("Can remove", async () => {
        const { getByTestId, getByText } = await CreateTestComponent();
        const device = getByTestId("DMXDevice");
        const deleteButton = getByTestId("DeviceDeleteButton");
        const deleteDialog = getByTestId("DeviceDeleteDialog");
        const Confirm = getByText("Confirm");
        const submit = getByText("SUBMIT");
        await expect.element(device).toBeVisible();
        await user.click(deleteButton);
        await expect.element(deleteDialog).toBeVisible();
        await user.click(Confirm);
        await user.click(submit);
        await expect(result).toEqual({
            test: {
                name: "test",
                devices: [],
            },
        } as testForm);
    });
});
