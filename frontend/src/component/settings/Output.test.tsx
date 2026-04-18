import { expect, describe, it, beforeEach } from "vitest";
import { render, type RenderResult } from "vitest-browser-react";
import { FormProvider, useForm } from "react-hook-form";
import { user, UserSetup } from "../../test/user_helper";
import type { TArtnet, TDMXHardware, TOSCServer } from "../../types";
import Outputs from "./Output";
import { UseMockServer } from "../../test/backend_helper";
import { http, HttpResponse } from "msw";

describe("Outputs", async () => {
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
            target?: string[];
            dmx: TDMXHardware;
            artnet: TArtnet;
        };
        osc: TOSCServer;
    }
    const defaultValue: testForm = {
        output: {
            artnet: {
                addr: "localhost",
                net: 0,
                subuni: 1,
                universe: 2,
            },
            dmx: {
                port: "COM1",
            },
            target: [],
        },
        osc: {
            channels: [],
            format: "{}",
            inverse: false,
            ip: "localhost",
            port: 20000,
            type: "int",
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
                    <Outputs />
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
                    result.output = v.output;
                }}
                value={value}
            ></TestForm>,
        );
    }
    it("Shown", async () => {
        const { getByTestId } = await CreateTestComponent();
        const devices = getByTestId("Outputs");
        await expect.element(devices).toBeVisible();
    });
    describe("Components", async () => {
        describe("Module select", async () => {
            it("Show checkbox for targets", async () => {
                const { getByRole } = await CreateTestComponent();
                const ftdi = getByRole("checkbox", { name: "ftdi" });
                const artnet = getByRole("checkbox", { name: "artnet" });
                await expect.element(ftdi).toBeVisible();
                await expect.element(artnet).toBeVisible();
            });
            it("Can show FTDI Options", async () => {
                const { getByRole, getByTestId } = await CreateTestComponent();
                const ftdi = getByRole("checkbox", { name: "ftdi" });
                const ftdiOptions = getByTestId("OutputDMX");
                await user.click(ftdi);
                await expect.element(ftdiOptions).toBeVisible();
            });
            it("Can show Artnet Options", async () => {
                const { getByRole, getByTestId } = await CreateTestComponent();
                const artnet = getByRole("checkbox", { name: "artnet" });
                const artnetOptions = getByTestId("OutputArtnet");
                await user.click(artnet);
                await expect.element(artnetOptions).toBeVisible();
            });
        });
        it("OSC Option", async () => {
            const { getByTestId } = await CreateTestComponent();
            const osc = getByTestId("OutputOSC");
            await expect.element(osc).toBeVisible();
        });
    });
    describe("", async () => {
        it("", async () => {
            const { getByTestId } = await CreateTestComponent({
                output: {
                    artnet: {
                        addr: "localhost",
                        net: 0,
                        subuni: 1,
                        universe: 2,
                    },
                    dmx: {
                        port: "COM1",
                    },
                },
                osc: {
                    channels: [],
                    format: "{}",
                    inverse: false,
                    ip: "localhost",
                    port: 20000,
                    type: "int",
                },
            });
            const devices = getByTestId("Outputs");
            await expect.element(devices).toBeVisible();
        });
    });
    it("Can submit", async () => {
        const { getByRole } = await CreateTestComponent();
        const ftdi = getByRole("checkbox", { name: "ftdi" });
        const submit = getByRole("button", { name: "SUBMIT" });
        await user.click(ftdi);
        await user.click(submit);
        await expect(result).toEqual({
            output: {
                artnet: {
                    addr: "localhost",
                    net: 0,
                    subuni: 1,
                    universe: 2,
                },
                dmx: {
                    port: "COM1",
                },
                target: ["ftdi"],
            },
            osc: {
                channels: [],
                format: "{}",
                inverse: false,
                ip: "localhost",
                port: 20000,
                type: "int",
            },
        } as testForm);
    });
});
