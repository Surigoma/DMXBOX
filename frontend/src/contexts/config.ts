export type ConfigBody = {
    backendPort: number;
};

class Configuration {
    body: ConfigBody = {
        backendPort: 8080,
    };
    isLoading: boolean = true;
    isError: boolean = false;
    constructor() {
        fetch("./config.json", {}).then((data) => {
            if (!data.ok) {
                this.isError = true;
                return;
            }
            data.json().then((v) => {
                this.body = v;
            });
            this.isLoading = false;
        });
        return;
    }
}

export default Configuration;
