import { createFileRoute, Link } from "@tanstack/react-router";
import { FrontConfigContext, genBackendPath, typedFetcher } from "./__root";
import useSWR from "swr";
import FadeControl from "../component/FadeControl";
import {
    Alert,
    FormControlLabel,
    FormGroup,
    Grid,
    Switch,
    Typography,
} from "@mui/material";
import ErrorComponent from "../component/Error";
import { useContext, useMemo, useState } from "react";
import {
    DMXGroupMap,
    Features,
    ControlMode,
    type TDMXGroupMap,
    type TFeatures,
} from "../types";
import MuteControl from "../component/MuteControl";
import { Light as SyntaxHighlighter } from "react-syntax-highlighter";
import json from "react-syntax-highlighter/dist/esm/languages/hljs/json";
import { atomOneDark } from "react-syntax-highlighter/dist/esm/styles/hljs";
SyntaxHighlighter.registerLanguage("json", json);

export const Route = createFileRoute("/")({
    component: ControlPage,
});

function ControlPage() {
    const config = useContext(FrontConfigContext);
    const {
        data: DMXData,
        error: DMXError,
        isLoading: DMXisLoading,
    } = useSWR(
        genBackendPath(config, "/api/v1/config/fade"),
        typedFetcher(DMXGroupMap),
    );
    const {
        data: FeaturesData,
        error: FeaturesError,
        isLoading: FeaturesLoading,
    } = useSWR(genBackendPath(config, "/api/features"), typedFetcher(Features));
    const {
        data: ControlModeData,
        error: ControlModeError,
        isLoading: ControlModeLoading,
        mutate: mutateControlMode,
    } = useSWR(
        genBackendPath(config, "/api/v1/control-mode"),
        typedFetcher(ControlMode),
    );
    const [showCutin, setCutin] = useState(false);
    const [controlModeUpdating, setControlModeUpdating] = useState(false);
    const [controlModeUpdateError, setControlModeUpdateError] = useState<string>();
    const dmxInfo = DMXData as TDMXGroupMap;
    const features = FeaturesData as TFeatures;
    const showMute = useMemo(
        () => features !== undefined && features.includes("osc"),
        [features],
    );
    const loadError = DMXError ?? FeaturesError ?? ControlModeError;
    async function setBrowserOnly(browserOnly: boolean) {
        setControlModeUpdating(true);
        setControlModeUpdateError(undefined);
        try {
            const response = await fetch(
                genBackendPath(config, "/api/v1/control-mode"),
                {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        "X-DMXBOX-Control": "web-ui",
                    },
                    body: JSON.stringify({ browserOnly }),
                },
            );
            if (!response.ok) {
                throw new Error(`Request failed: ${response.status}`);
            }
            await mutateControlMode(ControlMode.parse(await response.json()), {
                revalidate: false,
            });
        } catch (error) {
            setControlModeUpdateError(
                `Failed to change control mode. ${error instanceof Error ? error.message : "Please try again."}`,
            );
        } finally {
            setControlModeUpdating(false);
        }
    }
    if (loadError) {
        return (
            <ErrorComponent>
                Connection Error. Please check backend config or frontend{" "}
                <Link
                    to="config.json"
                    target="_blank"
                    rel="noopener noreferrer"
                >
                    config.json
                </Link>
                <SyntaxHighlighter
                    language="json"
                    style={atomOneDark}
                    wrapLines
                >
                    {JSON.stringify(loadError, undefined, 4)}
                </SyntaxHighlighter>
            </ErrorComponent>
        );
    }
    if (DMXisLoading || FeaturesLoading || ControlModeLoading) {
        return (
            <Grid
                container
                justifyContent="center"
                alignItems="center"
                padding="10px"
            >
                <a>Loading...</a>
            </Grid>
        );
    }
    return (
        <Grid container direction="column">
            {controlModeUpdateError && (
                <Alert severity="error" onClose={() => setControlModeUpdateError(undefined)}>
                    {controlModeUpdateError}
                </Alert>
            )}
            <Grid size="grow">
                <Grid
                    container
                    direction="row"
                    justifyContent="center"
                    alignItems="center"
                >
                    <Grid size="grow">
                        <Typography variant="h5" margin={2}>
                            Control
                        </Typography>
                    </Grid>
                    <Grid
                        size="auto"
                        justifyContent="center"
                        alignContent="center"
                    >
                        <FormGroup row>
                            <FormControlLabel
                                label="Browser only"
                                control={
                                    <Switch
                                        onChange={(e) => {
                                            void setBrowserOnly(e.target.checked);
                                        }}
                                        checked={ControlModeData?.browserOnly ?? false}
                                        disabled={controlModeUpdating}
                                    />
                                }
                            />
                            <FormControlLabel
                                label="CUT"
                                control={
                                    <Switch
                                        onChange={(e) => {
                                            setCutin(e.target.checked);
                                        }}
                                        checked={showCutin}
                                    />
                                }
                            ></FormControlLabel>
                        </FormGroup>
                    </Grid>
                </Grid>
                <Grid container spacing={3} padding={2}>
                    {Object.keys(dmxInfo).map((k) => {
                        {
                            return (
                                <Grid size={{ xs: 12, md: 6, lg: 4 }} key={k}>
                                    <FadeControl
                                        name={k}
                                        data={dmxInfo[k]}
                                        showCutin={showCutin}
                                    ></FadeControl>
                                </Grid>
                            );
                        }
                    })}
                </Grid>
            </Grid>
            {showMute && (
                <Grid size="grow">
                    <Grid
                        container
                        direction="row"
                        justifyContent="center"
                        alignItems="center"
                    >
                        <Grid size="grow">
                            <Typography variant="h5" margin={2}>
                                Mute
                            </Typography>
                        </Grid>
                    </Grid>
                    <Grid container spacing={3} padding={2}>
                        <Grid size="grow">
                            <MuteControl />
                        </Grid>
                    </Grid>
                </Grid>
            )}
        </Grid>
    );
}
