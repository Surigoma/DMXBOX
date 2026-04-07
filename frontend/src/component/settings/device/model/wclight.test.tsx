import { expect, describe, it } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";
import { FormProvider, useForm } from "react-hook-form";
import WCLight from "./wclight";
import { user, UserSetup } from "../../../../test/user_helper";

describe("WCLight option", async () => {
    UserSetup();
    interface testForm {
        test: {
            max: number[] | undefined[];
        };
    }
    const result: testForm = {
        test: {
            max: [255, 255, 0],
        },
    };
    function TestForm(f: {
        callback: (v: testForm) => void;
        defaultValue?: testForm;
    }) {
        const d: testForm =
            f.defaultValue !== undefined
                ? JSON.parse(JSON.stringify(f.defaultValue))
                : {
                      test: {
                          max: [255, 255, 0],
                      },
                  };
        const configForm = useForm<testForm>({
            defaultValues: d,
        });
        return (
            <FormProvider {...configForm}>
                <form onSubmit={configForm.handleSubmit((v) => f.callback(v))}>
                    <WCLight name="test" />
                    <input type="submit" value="SUBMIT" />
                </form>
            </FormProvider>
        );
    }
    function CreateTestComponent(
        defaultValue?: testForm,
    ): Promise<RenderResult> {
        return render(
            <TestForm
                callback={(v) => {
                    console.log(v);
                    result.test = v.test;
                }}
                defaultValue={defaultValue}
            ></TestForm>,
        );
    }
    describe("Dimmer", async () => {
        it("Can change value using mouse", async () => {
            const { getByTestId, getByRole } = await CreateTestComponent();
            const dimmer = getByTestId("OpDimmer");
            const DimmerSlider = getByRole("slider").first();
            await expect.element(dimmer).toBeVisible();
            const DimmerHeight = dimmer.element().clientHeight;
            const DimmerWidth = dimmer.element().clientWidth;
            await expect.element(DimmerSlider).toHaveValue("1");
            await user.click(dimmer, {
                position: {
                    x: DimmerWidth / 2,
                    y: DimmerHeight / 2,
                },
            });
            await expect.element(DimmerSlider).toHaveValue("0.5");
            await user.click(dimmer, {
                position: {
                    x: 0.5,
                    y: DimmerHeight / 2,
                },
            });
            await expect.element(DimmerSlider).toHaveValue("0");
        });
    });
    describe("Color Temp", async () => {
        it("Can change value using mouse", async () => {
            const { getByTestId, getByRole } = await CreateTestComponent();
            const temp = getByTestId("OpTemp");
            const TempSlider = getByRole("slider").nth(1);
            await expect.element(temp).toBeVisible();
            const TempHeight = temp.element().clientHeight;
            const TempWidth = temp.element().clientWidth;
            await expect.element(TempSlider).toHaveValue("0.5");
            await user.click(temp, {
                position: {
                    x: 0.5,
                    y: TempHeight / 2,
                },
            });
            await expect.element(TempSlider).toHaveValue("0");
            await user.click(temp, {
                position: {
                    x: TempWidth - 0.5,
                    y: TempHeight / 2,
                },
            });
            await expect.element(TempSlider).toHaveValue("1");
        });
    });
    describe("Calc parameters", async () => {
        interface calc {
            name: string;
            select: {
                dimmer: number;
                temp: number;
            };
            want: number[];
        }
        const tests: calc[] = [
            {
                name: "Dim: min, Temp: cool",
                select: {
                    dimmer: 0,
                    temp: 0,
                },
                want: [0, 0, 0],
            },
            {
                name: "Dim: max, Temp: cool",
                select: {
                    dimmer: 1,
                    temp: 0,
                },
                want: [255, 0, 0],
            },
            {
                name: "Dim: max, Temp: warn",
                select: {
                    dimmer: 1,
                    temp: 1,
                },
                want: [0, 255, 0],
            },
            {
                name: "Dim: max, Temp: even",
                select: {
                    dimmer: 1,
                    temp: 0.5,
                },
                want: [127, 128, 0],
            },
            {
                name: "Dim: half, Temp: cool",
                select: {
                    dimmer: 0.5,
                    temp: 0,
                },
                want: [128, 0, 0],
            },
            {
                name: "Dim: half, Temp: warn",
                select: {
                    dimmer: 0.5,
                    temp: 1,
                },
                want: [0, 128, 0],
            },
        ];
        tests.forEach((v) => {
            it(v.name, async () => {
                const { getByTestId, getByText } = await CreateTestComponent();
                const dimmer = getByTestId("OpDimmer");
                const submit = getByText("SUBMIT");
                await expect.element(dimmer).toBeVisible();
                const DimmerHeight = dimmer.element().clientHeight;
                const DimmerWidth = dimmer.element().clientWidth;
                const temp = getByTestId("OpTemp");
                await expect.element(temp).toBeVisible();
                const TempHeight = temp.element().clientHeight;
                const TempWidth = temp.element().clientWidth;
                await user.click(dimmer, {
                    position: {
                        x: (DimmerWidth - 1) * v.select.dimmer + 0.5,
                        y: DimmerHeight / 2,
                    },
                });
                await user.click(temp, {
                    position: {
                        x: (TempWidth - 1) * v.select.temp + 0.5,
                        y: TempHeight / 2,
                    },
                });
                await user.click(submit);
                await expect(result.test.max).toEqual(v.want);
            });
        });
    });
    it("Can submit", async () => {
        const { getByTestId, getByRole, getByText } =
            await CreateTestComponent();
        const dimmer = getByTestId("OpDimmer");
        const DimmerSlider = getByRole("slider").first();
        const submit = getByText("SUBMIT");
        await expect.element(dimmer).toBeVisible();
        const DimmerHeight = dimmer.element().clientHeight;
        const DimmerWidth = dimmer.element().clientWidth;
        await expect.element(DimmerSlider).toHaveValue("1");
        await user.click(dimmer, {
            position: {
                x: DimmerWidth / 2,
                y: DimmerHeight / 2,
            },
        });
        await expect.element(DimmerSlider).toHaveValue("0.5");

        await user.click(submit);

        await expect(result.test.max).toEqual([64, 64, 0]);
    });

    it("Can convert data when min length", async () => {
        const { getByTestId, getByRole, getByText } = await CreateTestComponent(
            {
                test: {
                    max: [255],
                },
            },
        );
        const dimmer = getByTestId("OpDimmer");
        const DimmerSlider = getByRole("slider").first();
        const submit = getByText("SUBMIT");
        await expect.element(dimmer).toBeVisible();
        const DimmerHeight = dimmer.element().clientHeight;
        const DimmerWidth = dimmer.element().clientWidth;
        await expect.element(DimmerSlider).toHaveValue("1");
        await user.click(dimmer, {
            position: {
                x: DimmerWidth / 2,
                y: DimmerHeight / 2,
            },
        });
        await expect.element(DimmerSlider).toHaveValue("0.5");

        await user.click(submit);

        await expect(result.test.max).toEqual([64, 64, 0]);
    });
    it("Can convert data when undefined", async () => {
        const { getByTestId, getByRole, getByText } = await CreateTestComponent(
            {
                test: {
                    max: [undefined],
                },
            },
        );
        const dimmer = getByTestId("OpDimmer");
        const DimmerSlider = getByRole("slider").first();
        const submit = getByText("SUBMIT");
        await expect.element(dimmer).toBeVisible();
        const DimmerHeight = dimmer.element().clientHeight;
        const DimmerWidth = dimmer.element().clientWidth;
        await expect.element(DimmerSlider).toHaveValue("1");
        await user.click(dimmer, {
            position: {
                x: DimmerWidth / 2,
                y: DimmerHeight / 2,
            },
        });
        await expect.element(DimmerSlider).toHaveValue("0.5");

        await user.click(submit);

        await expect(result.test.max).toEqual([64, 64, 0]);
    });
});
