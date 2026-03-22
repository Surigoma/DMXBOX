import { expect, describe, it } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";
import NumberField from "./numberField";
import { Controller, FormProvider, useForm } from "react-hook-form";
import { userEvent } from "vitest/browser";

describe("Number Field", async () => {
    interface testForm {
        test: number;
    }
    const result: testForm = {
        test: 0,
    };
    const user = userEvent.setup();
    function TestForm(f: {
        callback: (v: testForm) => void;
        min?: number;
        max?: number;
        help?: string;
    }) {
        const configForm = useForm<testForm>({
            defaultValues: {
                test: 0,
            },
        });
        return (
            <FormProvider {...configForm}>
                <form onSubmit={configForm.handleSubmit((v) => f.callback(v))}>
                    <Controller
                        control={configForm.control}
                        name="test"
                        render={({ field }) => (
                            <NumberField
                                min={f.min}
                                max={f.max}
                                help={f.help}
                                defaultValue={field.value}
                                onChange={(e) => {
                                    console.log(
                                        (
                                            e.target as EventTarget &
                                                HTMLInputElement
                                        ).value,
                                    );
                                    field.onChange(
                                        (
                                            e.target as EventTarget &
                                                HTMLInputElement
                                        ).value,
                                    );
                                }}
                                label="TEST"
                                id="test"
                                name="test"
                            ></NumberField>
                        )}
                    />{" "}
                    ,
                    <input type="submit" value="SUBMIT" />
                </form>
            </FormProvider>
        );
    }
    function CreateTestComponent(
        min?: number,
        max?: number,
        help?: string,
    ): Promise<RenderResult> {
        return render(
            <TestForm
                min={min}
                max={max}
                help={help}
                callback={(v) => {
                    result.test = v.test;
                }}
            ></TestForm>,
        );
    }
    it("Shown", async () => {
        const { getByRole } = await CreateTestComponent();
        const numberField = getByRole("textbox", { name: "test" });

        await expect.element(numberField).toBeVisible();
    });
    describe("Help message", async () => {
        it("Not showed help (No set min max)", async () => {
            const { locator, getByRole } = await CreateTestComponent();
            const numberField = getByRole("textbox", { name: "test" });
            await expect.element(numberField).toBeVisible();
            await expect.element(locator).not.toHaveTextContent("Enter value");
        });
        it("Show default message (Set min max)", async () => {
            const { getByText } = await CreateTestComponent(0, 100);
            await expect.element(getByText("Enter value")).toBeVisible();
        });
        it("Show help message (Set help)", async () => {
            const { getByText } = await CreateTestComponent(
                undefined,
                undefined,
                "HELP TEXT",
            );
            await expect.element(getByText("HELP TEXT")).toBeVisible();
        });
        it("Show help message (Set min max help)", async () => {
            const { getByText } = await CreateTestComponent(
                0,
                100,
                "HELP TEXT",
            );
            await expect.element(getByText("HELP TEXT")).toBeVisible();
        });
    });
    it("Can submit value", async () => {
        const { getByRole, getByText } = await CreateTestComponent(0, 100);

        const numberField = getByRole("textbox", { name: "test" });
        const submit = getByText("SUBMIT");
        await expect.element(numberField).toHaveValue("0");
        await user.click(submit);
        await expect(result).toStrictEqual({ test: 0 });
    });
    it("Can change by keyboard", async () => {
        const { getByRole } = await CreateTestComponent(0, 100);

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
        await user.click(numberField);
        await user.keyboard("{ArrowUp}");
        await expect.element(numberField).toHaveValue("1");
    });
    it("Can change by UpDown buttons", async () => {
        const { getByRole } = await CreateTestComponent(0, 100);

        const numberField = getByRole("textbox", { name: "test" });
        const UpButton = getByRole("Button").nth(0);
        const DownButton = getByRole("Button").nth(1);
        console.log(UpButton, DownButton);
        await expect.element(numberField).toHaveValue("0");

        await user.click(UpButton);
        await expect.element(numberField).toHaveValue("1");

        await user.click(DownButton);
        await expect.element(numberField).toHaveValue("0");
    });
});
