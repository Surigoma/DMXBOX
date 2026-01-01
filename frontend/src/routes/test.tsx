import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/test")({
    component: ResponsiveAppBar,
});

import { useContext } from "react";
import { fetcher, FrontConfigContext, genBackendPath } from "./__root";
import type { Config } from "../types";
import useSWR from "swr";
import Stack from "@mui/material/Stack";
import Inputs from "../component/settings/Input";
import Outputs from "../component/settings/Output";

function ResponsiveAppBar() {
    const config = useContext(FrontConfigContext);
    const { data, error, isLoading } = useSWR(
        genBackendPath(config, "/api/v1/config/all"),
        fetcher,
    );
    if (isLoading) {
        return <>Please wait</>;
    }
    if (error) {
        return <>Error</>;
    }

    const configData = data as Config;
    return (
        <Stack margin={2} spacing={2}>
            <Inputs config={configData}></Inputs>
            <Outputs config={configData.output}></Outputs>
        </Stack>
    );
}
export default ResponsiveAppBar;
