import { expect, describe, it } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";
import { FormProvider, useForm } from "react-hook-form";
import { user, UserSetup } from "../../../test/helper";
import type { TTCPServer } from "../../../types";
import InputTCP from "./tcp";

describe("Input TCP", async () => {
    UserSetup();
    interface testForm {
        tcp: TTCPServer;
    }
    const defaultValue: testForm = {
        tcp: {
            ip: "127.0.0.1",
            port: 8080,
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
                    <InputTCP />
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
                    result.tcp = v.tcp;
                }}
            ></TestForm>,
        );
    }
    it("Shown", async () => {
        const { getByTestId } = await CreateTestComponent();
        const http = getByTestId("InputTCP");
        await expect.element(http).toBeVisible();
    });
    describe("Components", async () => {
        describe("IP Address", async () => {
            it("Show IP Entry", async () => {
                const { getByTestId } = await CreateTestComponent();
                const http = getByTestId("OpIp");
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
    });
    it("Can submit", async () => {
        const { getByText, getByLabelText } = await CreateTestComponent();
        const submit = getByText("SUBMIT");
        const port = getByLabelText("Port");
        await user.clear(port);
        await user.fill(port, "10000");
        await user.click(submit);
        await expect(result).toEqual({
            tcp: {
                ip: "127.0.0.1",
                port: 10000,
            },
        } as testForm);
    });
});
