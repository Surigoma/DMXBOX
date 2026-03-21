import { expect, describe, it } from "vitest";
import { render } from "vitest-browser-react";
import NumberField from "./numberField";
import { FormProvider, useForm } from "react-hook-form";
import { userEvent } from "vitest/browser";

describe("Number Field", () => {
    it("Shown", async () => {
        interface testForm {
            test: boolean;
        }
        const result: testForm = {
            test: false,
        };
        const user = userEvent.setup();
        function TestForm(f: { callback: (v: testForm) => void }) {
            const configForm = useForm<testForm>({
                defaultValues: {
                    test: false,
                },
            });
            return (
                <FormProvider {...configForm}>
                    <form
                        onSubmit={configForm.handleSubmit((v) => f.callback(v))}
                    >
                        <NumberField
                            min={0}
                            max={100}
                            defaultValue={0}
                            label="TEST"
                            id="test"
                            name="test"
                        ></NumberField>
                        ,
                        <input type="submit" value="SUBMIT" />
                    </form>
                </FormProvider>
            );
        }
        const { getByRole } = await render(
            <TestForm
                callback={(v) => {
                    result.test = v.test;
                }}
            ></TestForm>,
        );
        const numberField = getByRole("textbox", { name: "test" });
        await expect.element(numberField).toHaveValue("0");

        await user.click(numberField);
        await user.keyboard("{ArrowUp}");
        await expect.element(numberField).toHaveValue("1");
        await user.fill(numberField, "100");
        await expect.element(numberField).toHaveValue("100");
        await user.keyboard("{ArrowUp}");
        await expect.element(numberField).toHaveValue("100");
        await user.fill(numberField, "0");
        await expect.element(numberField).toHaveValue("0");
        await user.keyboard("{ArrowDown}");
        await expect.element(numberField).toHaveValue("0");
    });
});
