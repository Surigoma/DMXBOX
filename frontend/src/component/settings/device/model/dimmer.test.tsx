import { expect, describe, it } from "vitest";
import { render } from "vitest-browser-react";
import { FormProvider, useForm } from "react-hook-form";
import Dimmer from "./dimmer";
import { user, UserSetup } from "../../../../test/user_helper";

describe("Dimmer Element", async () => {
    UserSetup();
    interface testForm {
        test: {
            max: number[];
        };
    }
    const result: testForm = {
        test: {
            max: [0],
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
                          max: [0],
                      },
                  };
        const configForm = useForm<testForm>({
            defaultValues: d,
        });
        return (
            <FormProvider {...configForm}>
                <form onSubmit={configForm.handleSubmit((v) => f.callback(v))}>
                    <Dimmer name="test" />
                    <input type="submit" value="SUBMIT" />
                </form>
            </FormProvider>
        );
    }
    function CreateTestComponent(defaultValue?: testForm) {
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
    it("Shown", async () => {
        const { getByTestId } = await CreateTestComponent();
        const dimmer = getByTestId("Dimmer");
        await expect.element(dimmer).toBeVisible();
    });
    it("Can change value using Click", async () => {
        const { getByTestId, getByRole } = await CreateTestComponent();

        const dimmer = getByTestId("OpDimmer");
        const slider = getByRole("slider");
        await expect.element(dimmer).toBeVisible();
        const height = dimmer.element().clientHeight;
        const width = dimmer.element().clientWidth;
        await expect.element(slider).toHaveValue("0");
        await user.click(dimmer, {
            position: {
                x: width / 2,
                y: height / 2,
            },
        });
        await expect.element(slider).toHaveValue("128");
        await user.click(dimmer, {
            position: {
                x: width - 0.6,
                y: height / 2,
            },
        });
        await expect.element(slider).toHaveValue("255");
    });
    it("Shown", async () => {
        const { getByTestId, getByRole, getByText } =
            await CreateTestComponent();
        const dimmer = getByTestId("OpDimmer");
        const slider = getByRole("slider");
        const submit = getByText("SUBMIT");

        await expect.element(dimmer).toBeVisible();
        const height = dimmer.element().clientHeight;
        const width = dimmer.element().clientWidth;
        await user.click(dimmer, {
            position: {
                x: width - 0.6,
                y: height / 2,
            },
        });
        await expect.element(slider).toHaveValue("255");
        await user.click(submit);
        await expect(result.test.max).toEqual([255]);
    });
    it("Can cutout data when over length", async () => {
        const { getByTestId, getByRole, getByText } = await CreateTestComponent(
            {
                test: {
                    max: [0, 0, 0],
                },
            },
        );
        const dimmer = getByTestId("OpDimmer");
        const slider = getByRole("slider");
        const submit = getByText("SUBMIT");

        await expect.element(dimmer).toBeVisible();
        const height = dimmer.element().clientHeight;
        const width = dimmer.element().clientWidth;
        await user.click(dimmer, {
            position: {
                x: width - 0.6,
                y: height / 2,
            },
        });
        await expect.element(slider).toHaveValue("255");
        await user.click(submit);
        await expect(result.test.max).toEqual([255]);
    });
});
