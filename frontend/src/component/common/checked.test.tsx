import { test, expect } from "vitest";
import { render } from "vitest-browser-react";
import Checked from "./checked";
import { FormProvider, useForm } from "react-hook-form";
import { userEvent } from "vitest/browser";

test("Check box is work", async () => {
    interface testForm {
        test: boolean;
    }
    const result: testForm = {
        test: false,
    };
    function TestForm(f: { callback: (v: testForm) => void }) {
        const configForm = useForm<testForm>({
            defaultValues: {
                test: false,
            },
        });
        return (
            <FormProvider {...configForm}>
                <form onSubmit={configForm.handleSubmit((v) => f.callback(v))}>
                    <Checked target="test" title="test" />
                    <input type="submit" value="SUBMIT" />
                </form>
            </FormProvider>
        );
    }
    const { getByRole, getByText } = await render(
        <TestForm
            callback={(v) => {
                result.test = v.test;
            }}
        ></TestForm>,
    );
    const checked = getByRole("checkbox");
    const submit = getByText("SUBMIT");
    await expect.element(checked).not.toBeChecked();
    await userEvent.click(checked);
    await expect.element(checked).toBeChecked();
    await userEvent.click(submit);
    await expect(result.test).toBe(true);
    await userEvent.click(checked);
    await expect.element(checked).not.toBeChecked();
    await userEvent.click(submit);
    await expect(result.test).toBe(false);
});
