import { Light as SyntaxHighlighter } from "react-syntax-highlighter";
import json from "react-syntax-highlighter/dist/esm/languages/hljs/json";
import { atomOneDark } from "react-syntax-highlighter/dist/esm/styles/hljs";
import { createFileRoute, useBlocker } from "@tanstack/react-router";
import useSWR from "swr";
import { FrontConfigContext, genBackendPath, typedFetcher } from "./__root";
import { useContext, useEffect, useState, type ReactElement } from "react";
import Grid from "@mui/material/Grid";
import { Config, type TConfig } from "../types";
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
    Stack,
    Typography,
} from "@mui/material";
import { useForm, FormProvider } from "react-hook-form";
import { MdExpandMore } from "react-icons/md";
import Devices from "../component/settings/Device";
import Inputs from "../component/settings/Input";
import Outputs from "../component/settings/Output";
import ControlModeSettings from "../component/settings/ControlMode";
import ControlModeStatus from "../component/ControlModeStatus";

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
        typedFetcher(Config),
    );
    const [submittedResult, setSubmittedResult] = useState<TConfig>();
    const displayResult = submittedResult ?? data;
    const [sendResultShow, setSendResultShow] = useState(false);
    const [restartMessage, setRestartMessage] = useState<string>();
    const [sendResult, setSendResult] = useState<postResult>({
        success: false,
        message: <>Not ready</>,
    });
    const [resultShow, setResultShow] = useState(false);
    const configForm = useForm<TConfig>({});
    const { isDirty, isSubmitting } = configForm.formState;
    const blocker = useBlocker({
        shouldBlockFn: () => isDirty || isSubmitting,
        enableBeforeUnload: isDirty || isSubmitting,
        withResolver: true,
    });
    useEffect(() => {
        if (blocker.status === "blocked" && !isDirty && !isSubmitting) blocker.proceed();
    }, [blocker, isDirty, isSubmitting]);
    useEffect(() => {
        if (data && !configForm.formState.isDirty && !configForm.formState.isSubmitting) {
            configForm.reset(data as TConfig, {
                keepDefaultValues: false,
            });
        }
    }, [data, configForm]);

    async function onSubmit(data: TConfig) {
        setSubmittedResult(data);
        try {
            const result = await fetch(genBackendPath(config, "/api/v1/config/save"), {
                method: "POST",
                body: JSON.stringify(data),
            });
            if (!result.ok) throw new Error(`Request failed: ${result.status}`);
            const saved = await result.json();
            setRestartMessage(saved.restartRequired ? saved.message : undefined);
            configForm.reset(data);
            setSendResult({
                success: true,
                message: <>Success</>,
            });
        } catch (error) {
            setSendResult({
                success: false,
                message: (
                    <>
                        Failed to send configuration.
                        <br />
                        <pre>{error instanceof Error ? error.message : "Please try again."}</pre>
                    </>
                ),
            });
        }
        setSendResultShow(true);
    }

    if (isLoading) {
        return <>Please wait</>;
    } else if (error) {
        return (
            <>
                Error: {error.name}
                <pre>
                    {JSON.stringify(JSON.parse(error.message), undefined, 4)}
                </pre>
            </>
        );
    }

    return (
        <>
            {restartMessage && <Alert severity="warning">{restartMessage}</Alert>}
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
                    <Box component="fieldset" disabled={isSubmitting} sx={{ border: 0, p: 0, m: 0, minWidth: 0, pointerEvents: isSubmitting ? "none" : "auto" }}>
                    <Grid container margin={2} gap={3} direction="column">
                        <Typography variant="h5">Configuration</Typography>
                        {isDirty && <Alert severity="info">Unsaved changes. Press Update to save.</Alert>}
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
                                        disabled={isSubmitting}
                                    >
                                        {isSubmitting ? "Saving..." : "Update"}
                                    </Button>
                                </FormControl>
                            </Grid>
                        </Grid>
                    </Grid>
                    </Box>
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
                            {JSON.stringify(displayResult, undefined, 4)}
                        </SyntaxHighlighter>
                    </DialogContent>
                    <DialogActions>
                        <Button onClick={() => setResultShow(false)}>
                            Close
                        </Button>
                    </DialogActions>
                </Dialog>
            </FormProvider>
            <Box sx={{ m: 2 }}>
                <Accordion defaultExpanded={false} data-testid="ControlModeSettings">
                    <AccordionSummary expandIcon={<MdExpandMore />}>
                        <Stack direction={{ xs: "column", sm: "row" }} alignItems={{ xs: "flex-start", sm: "center" }} spacing={1}>
                            <Typography component="span" variant="h5">Control mode</Typography>
                            <ControlModeStatus />
                        </Stack>
                    </AccordionSummary>
                    <AccordionDetails>
                        <ControlModeSettings />
                    </AccordionDetails>
                </Accordion>
            </Box>
            <Dialog open={blocker.status === "blocked"} onClose={() => blocker.reset?.()} aria-labelledby="unsaved-settings-title">
                <DialogTitle id="unsaved-settings-title">Leave configuration?</DialogTitle>
                <DialogContent>
                    {isSubmitting ? "Configuration is still being saved. Wait for the result before leaving." : "Your unsaved changes will be discarded."}
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => blocker.reset?.()}>Stay</Button>
                    <Button color="error" disabled={isSubmitting} onClick={() => blocker.proceed?.()}>Discard and leave</Button>
                </DialogActions>
            </Dialog>
        </>
    );
}
