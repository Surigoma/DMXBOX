import z from "zod";

const IPAddress = z.union([z.ipv4(), z.ipv6()]);
const TCPPort = z.number().min(0).max(65535);

export const DMXHardware = z
    .object({
        port: z
            .string()
            .describe("Hardware COM port (ex. COM1 or /dev/ttyUSB0"),
    })
    .describe("Configuration of USB DMX Device");
export type TDMXHardware = z.infer<typeof DMXHardware>;

export const Artnet = z
    .object({
        addr: z.string().describe("Target address for Artnet"),
        universe: z.number().min(0).max(15).describe("Universe of Artnet"),
        subuni: z.number().min(0).max(15).describe("Sub universe of Artnet"),
        net: z.number().min(0).max(127).describe("Net of Artnet"),
    })
    .describe("Artnet Configuration");
export type TArtnet = z.infer<typeof Artnet>;

export const OutputTargets = z.object({
    target: z
        .array(z.enum(["ftdi", "artnet", "console"]))
        .describe("Output target of DMX signal"),
    dmx: DMXHardware,
    artnet: Artnet,
});
export type TOutputTargets = z.infer<typeof OutputTargets>;

export const HttpServer = z.object({
    ip: IPAddress.describe("Listen address"),
    port: TCPPort.describe("Listen port"),
    accepts: z.array(z.string()).describe("CORS"),
});
export type THttpServer = z.infer<typeof HttpServer>;

export const TCPServer = z.object({
    ip: IPAddress.describe("Listen IP"),
    port: TCPPort.describe("Listen port"),
});
export type TTCPServer = z.infer<typeof TCPServer>;

export const DMXDevice = z.object({
    model: z.enum(["dimmer", "wclight"]).describe("DMX Device type"),
    channel: z.number().min(1).max(512).default(1).describe("DMX Channel"),
    max: z
        .array(z.number().min(0).max(255).default(255))
        .describe("Value of max for fade in"),
});
export type TDMXDevice = z.infer<typeof DMXDevice>;

export const DMXGroup = z.object({
    name: z.string().nonempty().describe("Group Name"),
    devices: z.array(DMXDevice),
});
export type TDMXGroup = z.infer<typeof DMXGroup>;

export const DMXGroupMap = z.record(
    z.string().describe("Group ID"),
    DMXGroup.describe("DMX Group Info"),
);
export type TDMXGroupMap = z.infer<typeof DMXGroupMap>;
export const DMXServer = z.object({
    groups: DMXGroupMap,
    fadeInterval: z.number().min(0).describe("Time of between start and end"),
    delay: z.number().min(0).describe("Time of before fade action"),
    fps: z.number().min(0).max(99).describe("Update FPS"),
});
export type TDMXServer = z.infer<typeof DMXServer>;

export const OSCServer = z.object({
    ip: IPAddress.describe("Target Address"),
    port: TCPPort.describe("Target Port"),
    format: z
        .string()
        .includes("{}")
        .describe("Format of OSC message (replace {} to channel)"),
    type: z.enum(["int", "float"]).describe("Format of value"),
    inverse: z.boolean().describe("Inverse value"),
    channels: z.array(z.number().min(0)).describe("Channel"),
});
export type TOSCServer = z.infer<typeof OSCServer>;

export const InputModules = z.record(
    z.enum(["http", "tcp"]).describe("Input module name"),
    z.boolean().describe("Enable of Input module"),
);
export type TInputModules = z.infer<typeof InputModules>;

export const Config = z.object({
    modules: InputModules,
    output: OutputTargets,
    http: HttpServer,
    tcp: TCPServer,
    dmx: DMXServer,
    osc: OSCServer,
});
export type TConfig = z.infer<typeof Config>;

export function DefaultConfig(): TConfig {
    return {
        modules: {
            http: false,
            tcp: false,
        },
        output: {
            target: [],
            dmx: {
                port: "COM1",
            },
            artnet: {
                addr: "127.0.0.1",
                net: 0,
                subuni: 0,
                universe: 0,
            },
        },
        http: {
            ip: "127.0.0.1",
            port: 8080,
            accepts: [],
        },
        tcp: {
            ip: "127.0.0.1",
            port: 50000,
        },
        dmx: {
            delay: 0,
            fadeInterval: 0,
            fps: 30,
            groups: {},
        },
        osc: {
            ip: "127.0.0.1",
            port: 49900,
            format: "",
            type: "float",
            inverse: false,
            channels: [],
        },
    };
}

export const ConsoleAPIResult = z.array(z.string());
