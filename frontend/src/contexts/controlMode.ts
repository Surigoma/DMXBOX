import { useContext } from "react";
import useSWR from "swr";
import { FrontConfigContext, genBackendPath } from "../routes/__root";
import { ControlMode } from "../types";

export function useControlMode() {
    const config = useContext(FrontConfigContext);
    return useSWR(
        genBackendPath(config, "/api/v1/control-mode"),
        async (url: string) => {
            const response = await fetch(url, { credentials: "include", mode: "cors" });
            if (!response.ok) throw new Error(`Control mode request failed: ${response.status}`);
            return ControlMode.parse(await response.json());
        },
        { refreshInterval: 1000, dedupingInterval: 500, errorRetryInterval: 1000 },
    );
}
