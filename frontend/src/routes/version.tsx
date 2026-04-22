import { createFileRoute } from "@tanstack/react-router";
import { FrontConfigContext, genBackendPath, typedFetcher } from "./__root";
import useSWR from "swr";
import {
    Grid,
    Paper,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Typography,
} from "@mui/material";
import { useContext } from "react";
import { Light as SyntaxHighlighter } from "react-syntax-highlighter";
import json from "react-syntax-highlighter/dist/esm/languages/hljs/json";
import z from "zod";
import packageJson from "../../package.json";

SyntaxHighlighter.registerLanguage("json", json);

export const Route = createFileRoute("/version")({
    component: ControlPage,
});
const VersionInfo = z.object({
    Version: z.string().describe("Background version"),
});
function ControlPage() {
    const config = useContext(FrontConfigContext);
    const { data, error, isLoading } = useSWR(
        genBackendPath(config, "/api/version"),
        typedFetcher(VersionInfo),
    );
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
        <Grid container justifyContent="center" alignItems="center" direction="column" gap={2} margin="10px">
            <Typography variant="h5">Version information</Typography>
            <TableContainer component={Paper} sx={{ maxWidth: "650px" }}>
                <Table>
                    <TableHead>
                        <TableRow>
                            <TableCell>Module</TableCell>
                            <TableCell>Version</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        <TableRow>
                            <TableCell>Frontend</TableCell>
                            <TableCell>{packageJson.version}</TableCell>
                        </TableRow>
                        <TableRow>
                            <TableCell>Backend</TableCell>
                            <TableCell>
                                {error
                                    ? "Failed to get."
                                    : isLoading
                                      ? "Loading..."
                                      : data!.Version}
                            </TableCell>
                        </TableRow>
                    </TableBody>
                </Table>
            </TableContainer>
        </Grid>
    );
}
