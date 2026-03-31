import { describe, it, expect, beforeEach } from "vitest";
import { user, UserSetup } from "../../../test/user_helper";
import { useForm, FormProvider } from "react-hook-form";
import { type RenderResult, render } from "vitest-browser-react";
import type { TDMXHardware } from "../../../types";
import OutputDMX from "./dmx";
import { createContext } from "react";
import { SWRConfig } from "swr";
import { UseMockServer } from "../../../test/backend_helper";
import { http, HttpResponse } from "msw";

describe("Output DMX", async () => {
    let APIController: string = "";
    beforeEach(() => {
        APIController = "200";
    });
    UseMockServer(
        http.get("*/api/v1/config/console", () => {
            switch (APIController) {
                case "200":
                    return HttpResponse.json(["COM1", "COM2", "COM3"], {
                        status: 200,
                    });
                case "200NoPort":
                    return HttpResponse.json([], { status: 200 });
                case "404":
                    return HttpResponse.json(undefined, { status: 404 });
            }
            return HttpResponse.json({}, { status: 500 });
        }),
    );
    UserSetup();

    interface testForm {
        output: {
            dmx: TDMXHardware;
        };
    }
    const defaultValue: testForm = {
        output: {
            dmx: {
                port: "COM1",
            },
        },
    };
    const result: testForm = JSON.parse(JSON.stringify(defaultValue));
    const FrontendConfigContext = createContext({
        port: 3030,
    });
    function TestForm(f: { callback: (v: testForm) => void }) {
        const configForm = useForm<testForm>({
            defaultValues: defaultValue,
        });
        return (
            <FrontendConfigContext value={{ port: 3030 }}>
                <FormProvider {...configForm}>
                    <form
                        onSubmit={configForm.handleSubmit((v) => f.callback(v))}
                    >
                        <OutputDMX />
                        <input type="submit" value="SUBMIT" />
                    </form>
                </FormProvider>
            </FrontendConfigContext>
        );
    }
    function CreateTestComponent(): Promise<RenderResult> {
        return render(
            <SWRConfig value={{ provider: () => new Map() }}>
                <TestForm
                    callback={(v) => {
                        console.log(v);
                        result.output = v.output;
                    }}
                ></TestForm>
            </SWRConfig>,
            {
                wrapper: ({ children }) => children,
            },
        );
    }
    it("Show Output DMX", async () => {
        const { getByTestId } = await CreateTestComponent();
        await expect.element(getByTestId("OutputDMX")).toBeVisible();
    });
    describe("Error Message", () => {
        it("Can show error message", async () => {
            APIController = "404";
            const { getByText } = await CreateTestComponent();
            await expect
                .element(getByText("Failed to get Console ports."))
                .toBeVisible();
        });
        it("Can No port", async () => {
            APIController = "200NoPort";
            const { getByText } = await CreateTestComponent();
            await expect
                .element(getByText("DMX Port is not found."))
                .toBeVisible();
        });
    });
    it("Can select COM Ports", async () => {
        const { getByTestId } = await CreateTestComponent();
        const port = getByTestId("OpPort");
        await user.click(port);

        for (let i = 1; i <= 3; i++) {
            await expect.element(getByTestId("OpPortCOM" + i)).toBeVisible();
        }
    });
    it("Can Submit", async () => {
        const { getByText, getByTestId } = await CreateTestComponent();
        const submit = getByText("SUBMIT");
        const port = getByTestId("OpPort");
        const com3 = getByTestId("OpPortCOM3");
        await user.click(port);
        await user.click(com3);
        await user.click(submit);
        await expect(result).toEqual({
            output: {
                dmx: {
                    port: "COM3",
                },
            },
        } as testForm);
    });
});
