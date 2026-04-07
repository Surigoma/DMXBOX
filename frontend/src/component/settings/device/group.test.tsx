import { expect, describe, it } from "vitest";
import { render } from "vitest-browser-react";
import { FormProvider, useForm } from "react-hook-form";
import type { TDMXGroup, TDMXServer } from "../../../types";
import Group from "./group";
import { user, UserSetup } from "../../../test/user_helper";
import { useMemo } from "react";

describe("DMXGroup Component", async () => {
    UserSetup();
    interface testForm {
        dmx: TDMXServer;
    }
    const defaultValue: testForm = {
        dmx: {
            delay: 0,
            fadeInterval: 0,
            fps: 30,
            groups: {
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
        },
    };
    const result: testForm = JSON.parse(JSON.stringify(defaultValue));
    function TestForm(f: { callback: (v: testForm) => void }) {
        const configForm = useForm<testForm>({
            defaultValues: defaultValue,
        });
        const groups = configForm.watch("dmx.groups") as {
            [key: string]: TDMXGroup;
        };
        const groupKeys = useMemo(() => Object.keys(groups ?? {}), [groups]);
        return (
            <FormProvider {...configForm}>
                <form onSubmit={configForm.handleSubmit((v) => f.callback(v))}>
                    {groupKeys.map((v) => (
                        <Group name={v} key={v} />
                    ))}
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
                    result.dmx = v.dmx;
                }}
            ></TestForm>,
        );
    }
    it("Shown", async () => {
        const { getByTestId } = await CreateTestComponent();
        const device = getByTestId("DMXDevice");
        await expect.element(device).toBeVisible();
    });
    it("Can submit", async () => {
        const { getByTestId, getByText } = await CreateTestComponent();
        const device = getByTestId("DMXDevice");
        const submit = getByText("SUBMIT");
        await expect.element(device).toBeVisible();
        await user.click(submit);
        await expect(result).toEqual(defaultValue);
    });
    describe("Edit Group information", async () => {
        it("Edit title", async () => {
            const { getByTestId, getByText, getByRole } =
                await CreateTestComponent();
            const groupButton = getByTestId("GroupEditButton");
            const groupDialog = getByTestId("GroupEditDialog");
            const title = getByRole("textbox", { name: "Title" });
            const saveButton = getByText("Edit");
            const submit = getByText("SUBMIT");
            await user.click(groupButton);
            await expect.element(groupDialog).toBeVisible();
            await user.clear(title);
            await user.fill(title, "TEST GROUP");
            await user.click(saveButton);
            await user.click(submit);
            await expect(result).toEqual({
                dmx: {
                    delay: 0,
                    fadeInterval: 0,
                    fps: 30,
                    groups: {
                        test: {
                            devices: [
                                {
                                    model: "dimmer",
                                    channel: 1,
                                    max: [0],
                                },
                            ],
                            name: "TEST GROUP",
                        },
                    },
                },
            } as testForm);
        });
        it("Edit ID", async () => {
            const { getByTestId, getByText, getByRole } =
                await CreateTestComponent();
            const groupButton = getByTestId("GroupEditButton");
            const groupDialog = getByTestId("GroupEditDialog");
            const ID = getByRole("textbox", { name: "ID" });
            const saveButton = getByText("Edit");
            const submit = getByText("SUBMIT");
            await user.click(groupButton);
            await expect.element(groupDialog).toBeVisible();
            await user.clear(ID);
            await user.fill(ID, "TESTID");
            await user.click(saveButton);
            await user.click(submit);
            await expect(result).toEqual({
                dmx: {
                    delay: 0,
                    fadeInterval: 0,
                    fps: 30,
                    groups: {
                        TESTID: {
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
                },
            } as testForm);
        });
        it("Delete Group", async () => {
            const { getByTestId, getByText } = await CreateTestComponent();
            const groupButton = getByTestId("GroupDeleteButton");
            const groupDialog = getByTestId("GroupDeleteDialog");
            const saveButton = getByText("Confirm");
            const submit = getByText("SUBMIT");
            await user.click(groupButton);
            await expect.element(groupDialog).toBeVisible();
            await user.click(saveButton);
            await user.click(submit);
            await expect(result).toEqual({
                dmx: {
                    delay: 0,
                    fadeInterval: 0,
                    fps: 30,
                    groups: {},
                },
            } as testForm);
        });
    });
    describe("Device control", async () => {
        it("Can add device", async () => {
            const { getByText, getByTestId } = await CreateTestComponent();

            const addButton = getByTestId("DeviceAddButton");
            const submit = getByText("SUBMIT");
            await expect.element(addButton).toBeVisible();
            await user.click(addButton);
            await user.click(submit);
            await expect(result).toEqual({
                dmx: {
                    delay: 0,
                    fadeInterval: 0,
                    fps: 30,
                    groups: {
                        test: {
                            name: "test",
                            devices: [
                                {
                                    model: "dimmer",
                                    channel: 1,
                                    max: [0],
                                },
                                {
                                    model: "dimmer",
                                    channel: 1,
                                    max: [255],
                                },
                            ],
                        },
                    },
                },
            } as testForm);
        });
        describe("Delete dialog", async () => {
            it("Can remove device", async () => {
                const { getByTestId, getByText } = await CreateTestComponent();

                const deleteButton = getByTestId("GroupDeleteButton");
                const deleteDialog = getByTestId("GroupDeleteDialog");
                const Confirm = getByText("Confirm");
                const submit = getByText("SUBMIT");
                await user.click(deleteButton);
                await expect.element(deleteDialog).toBeVisible();
                await user.click(Confirm);
                await user.click(submit);
                await expect(result).toEqual({
                    dmx: {
                        delay: 0,
                        fadeInterval: 0,
                        fps: 30,
                        groups: {
                        },
                    },
                } as testForm);
            });

            it("Can cancel", async () => {
                const { getByTestId, getByText } = await CreateTestComponent();

                const deleteButton = getByTestId("GroupDeleteButton");
                const deleteDialog = getByTestId("GroupDeleteDialog");
                const Cancel = getByText("Cancel");
                const submit = getByText("SUBMIT");
                await user.click(deleteButton);
                await expect.element(deleteDialog).toBeVisible();
                await user.click(Cancel);
                await user.click(submit);
                await expect(result).toEqual({
                    dmx: {
                        delay: 0,
                        fadeInterval: 0,
                        fps: 30,
                        groups: {
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
                        },
                    },
                } as testForm);
            });
        });
    });
});
