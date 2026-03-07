import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/test")({
    component: ResponsiveAppBar,
});

import { useContext, useEffect, useState } from "react";
import { fetcher, FrontConfigContext, genBackendPath } from "./__root";
import { type Config } from "../types";
import useSWR from "swr";
import Stack from "@mui/material/Stack";
import Inputs from "../component/settings/Input";
import Outputs from "../component/settings/Output";
import { FormProvider, useForm } from "react-hook-form";
import { Box, Button } from "@mui/material";

// export const BackendConfigContext = createContext<Config>(DefaultConfig());

function ResponsiveAppBar() {
    const config = useContext(FrontConfigContext);
    // const backend = useContext(BackendConfigContext);
    const { data, error, isLoading } = useSWR(
        genBackendPath(config, "/api/v1/config/all"),
        fetcher,
    );
    const [result, setResult] = useState(data as Config);
    const configForm = useForm<Config>({});
    useEffect(() => {
        setResult(data);
        configForm.reset(data as Config, {
            keepDefaultValues: false,
        });
    }, [data, configForm.reset]);

    function onSubmit(data: Config) {
        setResult(data);
    }

    if (isLoading) {
        return <>Please wait</>;
    } else if (error) {
        return <>Error</>;
    }

    return (
        <FormProvider {...configForm}>
            <Box component="form" onSubmit={configForm.handleSubmit(onSubmit)}>
                <Stack margin={2} spacing={2}>
                    <Inputs></Inputs>
                    <Outputs></Outputs>
                    <Button type="submit">Update</Button>
                    <pre>{JSON.stringify(result, undefined, 4)}</pre>
                </Stack>
            </Box>
        </FormProvider>
    );
}
export default ResponsiveAppBar;
