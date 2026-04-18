import { expect, describe, it } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";
import { FormProvider, useForm } from "react-hook-form";
import { user, UserSetup } from "../../test/user_helper";
import type { THttpServer, TInputModules, TTCPServer } from "../../types";
import Inputs from "./Input";

describe("Inputs", async () => {
    UserSetup();
    interface testForm {
        modules: TInputModules;
        http: THttpServer;
        tcp: TTCPServer;
    }
    const defaultValue: testForm = {
        modules: {
            http: false,
            tcp: false,
        },
        http: {
            accepts: [],
            ip: "localhost",
            port: 5000,
        },
        tcp: {
            ip: "localhost",
            port: 8080,
        },
    };
    const result: testForm = JSON.parse(JSON.stringify(defaultValue));
    function TestForm(f: {
        callback: (v: testForm) => void;
        value?: testForm;
    }) {
        const configForm = useForm<testForm>({
            defaultValues: f.value ?? defaultValue,
        });
        return (
            <FormProvider {...configForm}>
                <form onSubmit={configForm.handleSubmit((v) => f.callback(v))}>
                    <Inputs />
                    <input type="submit" value="SUBMIT" />
                </form>
            </FormProvider>
        );
    }
    function CreateTestComponent(value?: testForm): Promise<RenderResult> {
        return render(
            <TestForm
                callback={(v) => {
                    console.log(v);
                    result.modules = v.modules;
                    result.http = v.http;
                    result.tcp = v.tcp;
                }}
                value={value}
            ></TestForm>,
        );
    }
    it("Shown", async () => {
        const { getByTestId } = await CreateTestComponent();
        const devices = getByTestId("Inputs");
        await expect.element(devices).toBeVisible();
    });
    describe("Components", async () => {
        describe("Module select", async () => {
            it("Show checkbox for modules", async () => {
                const { getByRole } = await CreateTestComponent();
                const http = getByRole("checkbox", { name: "http" });
                const tcp = getByRole("checkbox", { name: "tcp" });
                await expect.element(http).toBeVisible();
                await expect.element(tcp).toBeVisible();
            });
            it("Can show No modules", async () => {
                const { getByText } = await CreateTestComponent();
                const message = getByText("Not selected", { exact: true });
                await expect.element(message).toBeVisible();
            });
            it("Can show HTTP", async () => {
                const { getByRole, getByTestId } = await CreateTestComponent();
                const http = getByRole("checkbox", { name: "http" });
                await user.click(http);
                const httpOption = getByTestId("InputHTTP");
                await expect.element(http).toBeVisible();
                await expect.element(httpOption).toBeVisible();
            });
            it("Can show TCP", async () => {
                const { getByRole, getByTestId } = await CreateTestComponent();
                const tcp = getByRole("checkbox", { name: "tcp" });
                await user.click(tcp);
                const tcpOption = getByTestId("InputTCP");
                await expect.element(tcp).toBeVisible();
                await expect.element(tcpOption).toBeVisible();
            });
            it("Can show both", async () => {
                const { getByRole, getByTestId } = await CreateTestComponent();
                const http = getByRole("checkbox", { name: "http" });
                const tcp = getByRole("checkbox", { name: "tcp" });
                await user.click(http);
                await user.click(tcp);
                const httpOption = getByTestId("InputHTTP");
                const tcpOption = getByTestId("InputTCP");
                await expect.element(tcp).toBeVisible();
                await expect.element(httpOption).toBeVisible();
                await expect.element(tcpOption).toBeVisible();
            });
        });
    });
    it("Can submit", async () => {
        const { getByRole } = await CreateTestComponent();
        const http = getByRole("checkbox", { name: "http" });
        const submit = getByRole("button", { name: "SUBMIT" });
        await user.click(http);
        await user.click(submit);
        await expect(result).toEqual({
            modules: {
                http: true,
                tcp: false,
            },
            http: {
                accepts: [],
                ip: "localhost",
                port: 5000,
            },
            tcp: {
                ip: "localhost",
                port: 8080,
            },
        } as testForm);
    });
});
