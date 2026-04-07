import { http, HttpResponse } from "msw";
import { expect, describe, it, beforeEach } from "vitest";
import { UseMockServer } from "../test/backend_helper";
import Configuration from "./config";

const sleep = (t: number) => new Promise(resolve => setTimeout(resolve, t));
const TestData = {
    "success": {
        "backendPort": 8080
    },
}

describe("Config", async () => {
    let APIController: string = "";

    beforeEach(() => {
        APIController = "200";
    });
    UseMockServer(
        http.get("*/config.json", () => {
            switch (APIController) {
                case "200":
                    return HttpResponse.json(TestData["success"], {
                        status: 200,
                    });
                case "404":
                    return HttpResponse.json(undefined, { status: 404 });
            }
            return HttpResponse.json({}, { status: 500 });
        }),
    );
    it("Can load config", async () => {
        const config: Configuration = new Configuration;
        while (config.isLoading) { await sleep(100); }
        expect(config.isError).toBe(false);
        expect(config.isLoading).toBe(false);
    })
    it("Can raise error flag", async () => {
        APIController = "404";
        const config: Configuration = new Configuration;
        while (config.isLoading) { await sleep(100); }
        expect(config.isError).toBe(true);
        expect(config.isLoading).toBe(false);
    })
})