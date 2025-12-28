import { createFileRoute, Link } from "@tanstack/react-router";
import { FrontConfigContext, fetcher, genBackendPath } from "./__root";
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
import type { DMXGroupInfo } from "../types";

export const Route = createFileRoute("/")({
    component: ControlPage,
});

function ControlPage() {
    const config = useContext(FrontConfigContext);
    const { data, error, isLoading } = useSWR(
        genBackendPath(config, "/api/v1/config/fade"),
        fetcher,
    );
    const [showCutin, setCutin] = useState(false);
    const dmxInfo = data as { [group: string]: DMXGroupInfo };
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
        <>
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
                <Grid size="auto" justifyContent="center" alignContent="center">
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
        </>
    );
}
