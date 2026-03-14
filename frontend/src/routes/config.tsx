import { Light as SyntaxHighlighter } from "react-syntax-highlighter";
import json from "react-syntax-highlighter/dist/esm/languages/hljs/json";
import { atomOneDark } from "react-syntax-highlighter/dist/esm/styles/hljs";
import { createFileRoute } from "@tanstack/react-router";
import useSWR from "swr";
import { FrontConfigContext, fetcher, genBackendPath } from "./__root";
import { useContext, useEffect, useState, type ReactElement } from "react";
import Grid from "@mui/material/Grid";
import type { Config } from "../types";
import {
    Accordion,
    AccordionDetails,
    AccordionSummary,
    Alert,
    Box,
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    FormControl,
    Snackbar,
    Typography,
} from "@mui/material";
import { useForm, FormProvider } from "react-hook-form";
import { MdExpandMore } from "react-icons/md";
import Devices from "../component/settings/Device";
import Inputs from "../component/settings/Input";
import Outputs from "../component/settings/Output";

SyntaxHighlighter.registerLanguage("json", json);

export const Route = createFileRoute("/config")({
    component: RouteComponent,
});

interface postResult {
    success: boolean;
    message: ReactElement;
}

function RouteComponent() {
    const config = useContext(FrontConfigContext);
    const { data, error, isLoading } = useSWR(
        genBackendPath(config, "/api/v1/config/all"),
        fetcher,
    );
    const [result, setResult] = useState(data as Config);
    const [sendResultShow, setSendResultShow] = useState(false);
    const [sendResult, setSendResult] = useState<postResult>({
        success: false,
        message: <>Not ready</>,
    });
    const [resultShow, setResultShow] = useState(false);
    const configForm = useForm<Config>({});
    useEffect(() => {
        setResult(data);
        configForm.reset(data as Config, {
            keepDefaultValues: false,
        });
    }, [data, configForm]);

    async function onSubmit(data: Config) {
        setResult(data);
        const result = await fetch(genBackendPath(config, "/api/v1/config/save"), {
            method: "POST",
            body: JSON.stringify(data),
        })
        if (result.ok) {
            setSendResult({
                success: true,
                message: <>Success</>,
            });
        } else {
            console.log(result.statusText);
            setSendResult({
                success: false,
                message: (
                    <>
                        Failed to send configuration.
                        <br />
                        <pre>{result.statusText}</pre>
                    </>
                ),
            });
        }
        setSendResultShow(true);
    }

    if (isLoading) {
        return <>Please wait</>;
    } else if (error) {
        return <>Error</>;
    }

    return (
        <>
            <Snackbar
                open={sendResultShow}
                onClose={() => setSendResultShow(false)}
                autoHideDuration={3000}
            >
                <Alert
                    onClose={() => setSendResultShow(false)}
                    severity={sendResult.success ? "success" : "error"}
                    variant="filled"
                    sx={{ width: "100%" }}
                >
                    {sendResult.message}
                </Alert>
            </Snackbar>
            <FormProvider {...configForm}>
                <Box
                    component="form"
                    onSubmit={configForm.handleSubmit(onSubmit)}
                >
                    <Grid container margin={2} gap={3} direction="column">
                        <Typography variant="h4">Configuration</Typography>
                        <Grid size="grow">
                            <Accordion defaultExpanded={false} key="input">
                                <AccordionSummary expandIcon={<MdExpandMore />}>
                                    <Typography component="span" variant="h5">
                                        Input
                                    </Typography>
                                </AccordionSummary>
                                <AccordionDetails>
                                    <Inputs />
                                </AccordionDetails>
                            </Accordion>
                            <Accordion defaultExpanded={false} key="devices">
                                <AccordionSummary expandIcon={<MdExpandMore />}>
                                    <Typography component="span" variant="h5">
                                        Devices
                                    </Typography>
                                </AccordionSummary>
                                <AccordionDetails>
                                    <Devices />
                                </AccordionDetails>
                            </Accordion>
                            <Accordion defaultExpanded={false} key="output">
                                <AccordionSummary expandIcon={<MdExpandMore />}>
                                    <Typography component="span" variant="h5">
                                        Output
                                    </Typography>
                                </AccordionSummary>
                                <AccordionDetails>
                                    <Outputs />
                                </AccordionDetails>
                            </Accordion>
                        </Grid>
                        <Grid
                            container
                            gap={3}
                            direction="row"
                            justifyContent="right"
                        >
                            <Grid>
                                <FormControl fullWidth>
                                    <Button
                                        variant="text"
                                        size="large"
                                        color="secondary"
                                        onClick={() => setResultShow(true)}
                                    >
                                        Show Result
                                    </Button>
                                </FormControl>
                            </Grid>
                            <Grid>
                                <FormControl fullWidth>
                                    <Button
                                        type="submit"
                                        variant="outlined"
                                        size="large"
                                        color="primary"
                                    >
                                        Update
                                    </Button>
                                </FormControl>
                            </Grid>
                        </Grid>
                    </Grid>
                </Box>
                <Dialog
                    open={resultShow}
                    aria-hidden={!resultShow}
                    fullWidth
                    maxWidth="md"
                >
                    <DialogTitle>Current Config</DialogTitle>
                    <DialogContent>
                        <SyntaxHighlighter
                            language="json"
                            style={atomOneDark}
                            wrapLines
                        >
                            {JSON.stringify(result, undefined, 4)}
                        </SyntaxHighlighter>
                    </DialogContent>
                    <DialogActions>
                        <Button onClick={() => setResultShow(false)}>
                            Close
                        </Button>
                    </DialogActions>
                </Dialog>
            </FormProvider>
        </>
    );
}
