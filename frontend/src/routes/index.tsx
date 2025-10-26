import { createFileRoute } from "@tanstack/react-router";
import { fetcher, genBackendPath } from "./__root";
import useSWR from "swr";
import FadeControl from "../component/FadeControl";
import { Grid } from "@mui/material";
import ErrorComponent from "../component/Error";

export const Route = createFileRoute("/")({
    component: ControlPage,
});

export interface DMXdeviceInfo {
    model: string;
    channel: number;
    max: number[];
}

function ControlPage() {
    //const config = useContext(ConfigContext);
    const { data, error, isLoading } = useSWR(
        genBackendPath("/api/v1/config/fade"),
        fetcher,
    );
    const dmxInfo = data as { [group: string]: DMXdeviceInfo[] };
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
            <h3>Control</h3>
            <Grid container spacing={2} padding={2}>
                {Object.keys(dmxInfo).map((k) => {
                    {
                        return (
                            <Grid size={6}>
                                <FadeControl group={k}></FadeControl>
                            </Grid>
                        );
                    }
                })}
            </Grid>
        </>
    );
}
