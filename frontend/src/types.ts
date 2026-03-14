export interface DMXGroupInfo {
    name: string;
    devices: [
        {
            model: string;
            channel: number;
            max: number[];
        },
    ];
}
export interface DMXHardware {
    port: string;
}
export interface Artnet {
    addr: string;
    universe: number;
    subuni: number;
    net: number;
}
export interface OutputTargets {
    target: string[];
    dmx: DMXHardware;
    artnet: Artnet;
}
export interface HttpServer {
    ip: string;
    port: number;
    accepts: string[];
}
export interface TCPServer {
    ip: string;
    port: number;
}
export interface DMXDevice {
    model: string;
    channel: number;
    max: number[];
}
export interface DMXGroup {
    name: string;
    devices: DMXDevice[];
}
export interface DMXServer {
    groups: { [name: string]: DMXGroup };
    fadeInterval: number;
    delay: number;
    fps: number;
}
export interface OSCServer {
    ip: string;
    port: number;
    format: string;
    type: string;
    inverse: boolean;
    channels: number[];
}
export interface InputModules {
    http: boolean;
    tcp: boolean;
}
export interface Config {
    modules: InputModules;
    output: OutputTargets;
    http: HttpServer;
    tcp: TCPServer;
    dmx: DMXServer;
    osc: OSCServer;
}

export function DefaultConfig(): Config {
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
