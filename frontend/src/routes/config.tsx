import { createFileRoute, Link } from "@tanstack/react-router";
import useSWR from "swr";
import { FrontConfigContext, fetcher, genBackendPath } from "./__root";
import { useContext } from "react";
import ErrorComponent from "../component/Error";
import Grid from "@mui/material/Grid";
import type { Config } from "../types";

export const Route = createFileRoute("/config")({
    component: RouteComponent,
});

let configData: Config;
function RouteComponent() {
    const config = useContext(FrontConfigContext);
    const { data, error, isLoading } = useSWR(
        genBackendPath(config, "/api/v1/config/all"),
        fetcher,
    );
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
    configData = data as Config;
    return <pre>{JSON.stringify(configData, undefined, 4)}</pre>;
}
