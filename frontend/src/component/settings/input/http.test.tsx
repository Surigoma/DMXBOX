import { expect, describe, it } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";
import { FormProvider, useForm } from "react-hook-form";
import { user, UserSetup } from "../../../test/helper";
import InputHTTP from "./http";
import type { THttpServer } from "../../../types";

describe("Input HTTP", async () => {
    UserSetup();
    interface testForm {
        http: THttpServer;
    }
    const defaultValue: testForm = {
        http: {
            ip: "127.0.0.1",
            port: 8080,
            accepts: ["test"],
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
                    <InputHTTP />
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
                    result.http = v.http;
                }}
            ></TestForm>,
        );
    }
    it("Shown", async () => {
        const { getByTestId } = await CreateTestComponent();
        const http = getByTestId("InputHTTP");
        await expect.element(http).toBeVisible();
    });
    describe("Components", async () => {
        describe("IP Address", async () => {
            it("Show IP Entry", async () => {
                const { getByTestId } = await CreateTestComponent();
                const http = getByTestId("OpIP");
                await expect.element(http).toBeVisible();
            });
        });
        describe("Port", async () => {
            it("Show Port Entry", async () => {
                const { getByLabelText } = await CreateTestComponent();
                const port = getByLabelText("Port");
                await expect.element(port).toBeVisible();
            });
            it("Port range is 1 to 65535", async () => {
                const { getByLabelText } = await CreateTestComponent();
                const port = getByLabelText("Port");
                await expect.element(port).toBeVisible();
                await user.click(port);
                await user.clear(port);
                await user.fill(port, "0");
                await user.keyboard("{ArrowDown}");
                await expect.element(port).toHaveValue("1");
                await user.clear(port);
                await user.fill(port, "65536");
                await user.keyboard("{ArrowUp}");
                await expect.element(port).toHaveValue("65535");
            });
        });
        describe("Accepts", async () => {
            it("Show Accepts Entry", async () => {
                const { getByTestId } = await CreateTestComponent();
                const accepts = getByTestId("OpAccepts");
                await expect.element(accepts).toBeVisible();
            });
            it("Can add accept address", async () => {
                const { getByTestId, getByText, getByRole } =
                    await CreateTestComponent();
                const accepts = getByTestId("OpAccepts");
                const acceptInput = getByRole("textbox", {
                    name: "Accept addresses",
                });
                console.log(accepts.element());
                await expect.element(accepts).toBeVisible();
                await user.click(acceptInput);
                await user.fill(acceptInput, "test.address");
                await user.keyboard("{ArrowRight}");
                await user.keyboard("{Enter}");
                await expect.element(getByText("test.address")).toBeVisible();
            });
            it("Can remote accept address", async () => {
                const { getByText, getByRole } = await CreateTestComponent();
                const accepts = getByText("test");
                const removeButton = getByRole("button", { name: "delete" });
                console.log(accepts.element());
                await expect.element(accepts).toBeVisible();
                await user.click(removeButton);
                await expect(accepts.elements().length).toBe(0);
            });
        });
    });
    it("Can submit", async () => {
        const { getByText, getByLabelText } = await CreateTestComponent();
        const submit = getByText("SUBMIT");
        const port = getByLabelText("Port");
        await user.clear(port);
        await user.fill(port, "10000");
        await user.click(submit);
        await expect(result).toEqual({
            http: {
                ip: "127.0.0.1",
                port: 10000,
                accepts: ["test"],
            },
        } as testForm);
    });
});
