import { createFileRoute, Link } from "@tanstack/react-router";
import { FrontConfigContext, genBackendPath, typedFetcher } from "./__root";
import useSWR from "swr";
import FadeControl from "../component/FadeControl";
import {
    FormControlLabel,
    FormGroup,
    Grid,
    Switch,
    Typography,
} from "@mui/material";
import ErrorComponent from "../component/Error";
import { useContext, useState } from "react";
import { DMXGroupMap, type TDMXGroupMap } from "../types";
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
    const { data, error, isLoading } = useSWR(
        genBackendPath(config, "/api/v1/config/fade"),
        typedFetcher(DMXGroupMap),
    );
    const [showCutin, setCutin] = useState(false);
    const dmxInfo = data as TDMXGroupMap;
    if (error) {
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
                    {JSON.stringify(error, undefined, 4)}
                </SyntaxHighlighter>
            </ErrorComponent>
        );
    }
    if (isLoading) {
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
                        <FormGroup>
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
        </Grid>
    );
}
