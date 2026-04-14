import { expect, describe, it } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";
import { FormProvider, useForm } from "react-hook-form";
import { user, UserSetup } from "../../test/user_helper";
import type { TDMXGroupMap } from "../../types";
import Devices from "./Device";

describe("Devices", async () => {
    UserSetup();
    interface testForm {
        dmx: {
            groups?: TDMXGroupMap;
            delay?: number;
            fps?: number;
            fadeInterval?: number;
        };
    }
    const defaultValue: testForm = {
        dmx: {
            delay: 0,
            fadeInterval: 1,
            fps: 40,
            groups: {},
        },
    };
    const result: testForm = JSON.parse(JSON.stringify(defaultValue));
    function TestForm(f: {
        callback: (v: testForm) => void;
        value?: testForm;
    }) {
        const configForm = useForm<testForm>({
            defaultValues: f.value ?? defaultValue,
        });
        return (
            <FormProvider {...configForm}>
                <form onSubmit={configForm.handleSubmit((v) => f.callback(v))}>
                    <Devices />
                    <input type="submit" value="SUBMIT" />
                </form>
            </FormProvider>
        );
    }
    function CreateTestComponent(value?: testForm): Promise<RenderResult> {
        return render(
            <TestForm
                callback={(v) => {
                    console.log(v);
                    result.dmx = v.dmx;
                }}
                value={value}
            ></TestForm>,
        );
    }
    it("Shown", async () => {
        const { getByTestId } = await CreateTestComponent();
        const devices = getByTestId("Devices");
        await expect.element(devices).toBeVisible();
    });
    describe("Components", async () => {
        describe("Update FPS", async () => {
            it("Show Update FPS", async () => {
                const { getByLabelText } = await CreateTestComponent();
                const updateFPS = getByLabelText("Update FPS");
                await expect.element(updateFPS).toBeVisible();
            });
            it("Can change value", async () => {
                const { getByLabelText } = await CreateTestComponent();
                const updateFPS = getByLabelText("Update FPS");
                await expect.element(updateFPS).toBeVisible();
                await user.click(updateFPS);
                await user.keyboard("{ArrowUp}");
                await expect.element(updateFPS).toHaveValue("41");
            });
            it("Check Default value", async () => {
                const { getByLabelText } = await CreateTestComponent({
                    dmx: {
                        groups: {},
                    },
                });
                const updateFPS = getByLabelText("Update FPS");
                await expect.element(updateFPS).toBeVisible();
                await expect.element(updateFPS).toHaveValue("40");
            });
        });
        describe("Fade Interval", async () => {
            it("Show Fade Interval", async () => {
                const { getByLabelText } = await CreateTestComponent();
                const fadeInterval = getByLabelText("Fade Interval");
                await expect.element(fadeInterval).toBeVisible();
            });
            it("Can change value", async () => {
                const { getByLabelText } = await CreateTestComponent();
                const fadeInterval = getByLabelText("Fade Interval");
                await expect.element(fadeInterval).toBeVisible();
                await user.click(fadeInterval);
                await user.keyboard("{ArrowUp}");
                await expect.element(fadeInterval).toHaveValue("1.1");
            });
            it("Check Default value", async () => {
                const { getByLabelText } = await CreateTestComponent({
                    dmx: {
                        groups: {},
                    },
                });
                const fadeInterval = getByLabelText("Fade Interval");
                await expect.element(fadeInterval).toBeVisible();
                await expect.element(fadeInterval).toHaveValue("0.7");
            });
        });
        describe("Delay", async () => {
            it("Show Delay", async () => {
                const { getByLabelText } = await CreateTestComponent();
                const delay = getByLabelText("Delay");
                await expect.element(delay).toBeVisible();
            });
            it("Can change value", async () => {
                const { getByLabelText } = await CreateTestComponent();
                const delay = getByLabelText("Delay");
                await expect.element(delay).toBeVisible();
                await user.click(delay);
                await user.keyboard("{ArrowUp}");
                await expect.element(delay).toHaveValue("0.1");
            });
            it("Check Default value", async () => {
                const { getByLabelText } = await CreateTestComponent({
                    dmx: {
                        groups: {},
                    },
                });
                const delay = getByLabelText("Delay");
                await expect.element(delay).toBeVisible();
                await expect.element(delay).toHaveValue("0");
            });
        });
        describe("Group", async () => {
            it("Show No Group when no has group", async () => {
                const { getByText } = await CreateTestComponent();
                await expect.element(getByText("No Groups")).toBeVisible();
            });
            it("Show No Group when no has group and no object", async () => {
                const { getByText } = await CreateTestComponent({
                    dmx: {},
                });
                await expect.element(getByText("No Groups")).toBeVisible();
            });
            it("Show Group", async () => {
                const { getByText } = await CreateTestComponent({
                    dmx: {
                        groups: {
                            test: {
                                devices: [],
                                name: "Test",
                            },
                        },
                    },
                });
                await expect.element(getByText("Test")).toBeVisible();
            });
            describe("Group Dialog", async () => {
                it("Can add Group", async () => {
                    const { getByRole, getByTestId, getByText } =
                        await CreateTestComponent({
                            dmx: {
                                groups: {},
                            },
                        });
                    const addGroup = getByRole("button", {
                        name: "Add Group",
                    });
                    await user.click(addGroup);
                    const groupTitle =
                        getByTestId("OpGroupTitle").getByRole("textbox");
                    const groupId =
                        getByTestId("OpGroupId").getByRole("textbox");
                    const saveButton = getByRole("button", { name: "Add" });
                    const submit = getByText("SUBMIT");
                    await user.click(groupTitle);
                    await user.fill(groupTitle, "TEST_NAME");
                    await user.click(groupId);
                    await user.fill(groupId, "TEST_ID");
                    await user.click(saveButton);
                    await user.click(submit);
                });

                it("Can cancel", async () => {
                    const { getByRole, getByText } = await CreateTestComponent({
                        dmx: {
                            groups: {},
                        },
                    });
                    const addGroup = getByRole("button", {
                        name: "Add Group",
                    });
                    await user.click(addGroup);
                    const cancel = getByRole("button", { name: "Cancel" });
                    const submit = getByText("SUBMIT");
                    await user.click(cancel);
                    await user.click(submit);
                });
            });
        });
    });
    it("Can submit", async () => {
        const { getByText, getByLabelText } = await CreateTestComponent();
        const delay = getByLabelText("Delay");
        const submit = getByText("SUBMIT");
        await user.click(delay);
        await user.cleanup();
        await user.fill(delay, "10");
        await user.click(submit);
        await expect(result).toEqual({
            dmx: {
                delay: 10,
                fadeInterval: 1,
                fps: 40,
                groups: {},
            },
        } as testForm);
    });
});
