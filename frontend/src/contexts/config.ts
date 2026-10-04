import z from "zod";
const ConfigBody = z.object({
    backendPort: z.number().int().min(1).max(65535)
})
export type ConfigBody = z.infer<typeof ConfigBody>;

class Configuration {
    body: ConfigBody = {
        backendPort: 8080,
    };
    isLoading: boolean = true;
    isError: boolean = false;
    ready: Promise<void>;
    constructor() {
        this.ready = fetch("./config.json", {}).then(async (data) => {
            if (!data.ok) throw new Error(`Configuration request failed: ${data.status}`);
            this.body = ConfigBody.parse(await data.json());
        }).catch(() => {
            this.isError = true;
        }).finally(() => {
            this.isLoading = false;
        });
        return;
    }
}

export default Configuration;
