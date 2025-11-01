import { createFileRoute } from "@tanstack/react-router";
import { ConfigContext, fetcher, genBackendPath } from "./__root";
import useSWR from "swr";
import FadeControl from "../component/FadeControl";
import { Grid, Typography } from "@mui/material";
import ErrorComponent from "../component/Error";
import { useContext } from "react";

export const Route = createFileRoute("/")({
    component: ControlPage,
});

export interface DMXGroupInfo {
    name: string;
    devices: [
        {
            model: string;
            channel: number;
            max: number[];
        },
    ];
}

function ControlPage() {
    const config = useContext(ConfigContext);
    const { data, error, isLoading } = useSWR(
        genBackendPath(config, "/api/v1/config/fade"),
        fetcher,
    );
    const dmxInfo = data as { [group: string]: DMXGroupInfo };
    if (error) {
        return (
            <ErrorComponent>
                Connection Error. Plase check backend config or frontend{" "}
                <a href="/config.json">config.json</a>
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
            <Typography variant="h5" margin={2}>
                Control
            </Typography>
            <Grid container spacing={2} padding={2}>
                {Object.keys(dmxInfo).map((k) => {
                    {
                        return (
                            <Grid size={{ xs: 12, md: 6, lg: 3 }} key={k}>
                                <FadeControl
                                    name={k}
                                    data={dmxInfo[k]}
                                ></FadeControl>
                            </Grid>
                        );
                    }
                })}
            </Grid>
        </>
    );
}
