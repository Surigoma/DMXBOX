import { expect, describe, it } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";
import { FormProvider, useForm } from "react-hook-form";
import { user, UserSetup } from "../../../test/user_helper";
import type { TArtnet } from "../../../types";
import OutputArtnet from "./artnet";

describe("Output Artnet", async () => {
    UserSetup();
    interface testForm {
        output: {
            artnet: TArtnet;
        };
    }
    const defaultValue: testForm = {
        output: {
            artnet: {
                addr: "127.0.0.1",
                net: 0,
                subuni: 0,
                universe: 0,
            },
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
                    <OutputArtnet />
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
                    result.output = v.output;
                }}
            ></TestForm>,
        );
    }
    it("Shown", async () => {
        const { getByTestId } = await CreateTestComponent();
        const http = getByTestId("OutputArtnet");
        await expect.element(http).toBeVisible();
    });
    describe("Components", async () => {
        describe("Address", async () => {
            it("Show Address", async () => {
                const { getByLabelText } = await CreateTestComponent();
                const address = getByLabelText("Address");
                await expect.element(address).toBeVisible();
            });
        });
        interface Components {
            name: string;
            range: number[];
        }
        const targets: Components[] = [
            {
                name: "Universe",
                range: [0, 15],
            },
            {
                name: "Sub Universe",
                range: [0, 15],
            },
            {
                name: "Net",
                range: [0, 15],
            },
        ];
        targets.forEach((v) => {
            describe(v.name, async () => {
                it("Show " + v.name + " Entry", async () => {
                    const { getByLabelText } = await CreateTestComponent();
                    const universe = getByLabelText(v.name, { exact: true });
                    await expect.element(universe).toBeVisible();
                });
                it(
                    "Universe range is  " + v.range[0] + " to " + v.range[1],
                    async () => {
                        const { getByLabelText } = await CreateTestComponent();
                        const universe = getByLabelText(v.name, {
                            exact: true,
                        });
                        await expect.element(universe).toBeVisible();
                        await user.click(universe);
                        await user.clear(universe);
                        await user.fill(universe, (v.range[0] - 1).toString());
                        await user.keyboard("{ArrowDown}");
                        await expect
                            .element(universe)
                            .toHaveValue(v.range[0].toString());
                        await user.clear(universe);
                        await user.fill(universe, (v.range[1] + 1).toString());
                        await user.keyboard("{ArrowUp}");
                        await expect
                            .element(universe)
                            .toHaveValue(v.range[1].toString());
                    },
                );
            });
        });
    });
    it("Can submit", async () => {
        const { getByText, getByLabelText } = await CreateTestComponent();
        const submit = getByText("SUBMIT");
        const address = getByLabelText("Address");
        await user.clear(address);
        await user.fill(address, "test");
        await user.click(submit);
        await expect(result).toEqual({
            output: {
                artnet: {
                    addr: "test",
                    net: 0,
                    subuni: 0,
                    universe: 0,
                },
            },
        } as testForm);
    });
});
