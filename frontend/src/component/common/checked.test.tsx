import { expect, describe, it } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";
import Checked from "./checked";
import { FormProvider, useForm } from "react-hook-form";
import { user, UserSetup } from "../../test/user_helper";

describe("Checked", async () => {
    UserSetup();
    interface testForm {
        bool: boolean;
        strArray: string[];
        un?: undefined;
    }
    let result: testForm = {
        bool: false,
        strArray: [],
    };
    function TestForm(f: {
        callback: (v: testForm) => void;
        target: string;
        value?: string;
    }) {
        const configForm = useForm<testForm>({
            defaultValues: {
                bool: false,
                strArray: [],
            },
        });
        return (
            <FormProvider {...configForm}>
                <form onSubmit={configForm.handleSubmit((v) => f.callback(v))}>
                    <Checked target={f.target} title="test" value={f.value} />
                    <input type="submit" value="SUBMIT" />
                </form>
            </FormProvider>
        );
    }
    function CreateTestComponent(
        target: string,
        value?: string,
    ): Promise<RenderResult> {
        return render(
            <TestForm
                target={target}
                value={value}
                callback={(v) => {
                    console.log(v);
                    result = v;
                }}
            ></TestForm>,
        );
    }
    it("Shown", async () => {
        const { getByRole } = await CreateTestComponent("bool");
        const checked = getByRole("checkbox");
        await expect.element(checked).toBeVisible();
    });
    it("Can checkable", async () => {
        const { getByRole } = await CreateTestComponent("bool");
        const checked = getByRole("checkbox");
        await expect.element(checked).not.toBeChecked();
        await user.click(checked);
        await expect.element(checked).toBeChecked();
        await user.click(checked);
        await expect.element(checked).not.toBeChecked();
    });
    describe("Types", async () => {
        it("Can submit for undefined", async () => {
            const { getByRole, getByText } = await CreateTestComponent("un");
            const checked = getByRole("checkbox");
            const submit = getByText("SUBMIT");
            await user.click(checked);
            await user.click(submit);
            await expect(result.bool).toBe(false);
            await user.click(checked);
            await expect.element(checked).not.toBeChecked();
            await user.click(submit);
            await expect(result.bool).toBe(false);
        });
        it("Can submit for boolean", async () => {
            const { getByRole, getByText } = await CreateTestComponent("bool");
            const checked = getByRole("checkbox");
            const submit = getByText("SUBMIT");
            await user.click(checked);
            await user.click(submit);
            await expect(result.bool).toBe(true);
            await user.click(checked);
            await expect.element(checked).not.toBeChecked();
            await user.click(submit);
            await expect(result.bool).toBe(false);
        });
        it("Can submit for list of string", async () => {
            const { getByRole, getByText } = await CreateTestComponent(
                "strArray",
                "test3",
            );
            const checked = getByRole("checkbox");
            const submit = getByText("SUBMIT");
            await expect.element(checked).toBeVisible();
            await user.click(checked);
            await user.click(submit);
            await expect.element(checked).toBeChecked();
            await expect(result.strArray).toEqual(["test3"]);
            await user.click(checked);
            await expect.element(checked).not.toBeChecked();
            await user.click(submit);
            await expect(result.strArray).toEqual([]);
        });
    });
});
