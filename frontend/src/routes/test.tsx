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
import {
    Accordion,
    AccordionDetails,
    AccordionSummary,
    Box,
    Button,
    Typography,
} from "@mui/material";
import { MdExpandMore } from "react-icons/md";

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
                    <Accordion defaultExpanded={true} key="input">
                        <AccordionSummary expandIcon={<MdExpandMore />}>
                            <Typography component="span" variant="h4">
                                Input
                            </Typography>
                        </AccordionSummary>
                        <AccordionDetails>
                            <Inputs />
                        </AccordionDetails>
                    </Accordion>
                    <Accordion defaultExpanded={true} key="output">
                        <AccordionSummary expandIcon={<MdExpandMore />}>
                            <Typography component="span" variant="h4">
                                Output
                            </Typography>
                        </AccordionSummary>
                        <AccordionDetails>
                            <Outputs />
                        </AccordionDetails>
                    </Accordion>
                    <Button type="submit">Update</Button>
                    <pre>{JSON.stringify(result, undefined, 4)}</pre>
                </Stack>
            </Box>
        </FormProvider>
    );
}
export default ResponsiveAppBar;
