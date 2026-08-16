import { createFileRoute } from "@tanstack/react-router";
import { useContext } from "react";
import useSWR from "swr";
import {
    Paper,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Typography,
} from "@mui/material";
import { FrontConfigContext, genBackendPath, typedFetcher } from "./__root";
import { OperationLog } from "../types";

export const Route = createFileRoute("/operations")({
    component: OperationsPage,
});

function OperationsPage() {
    const config = useContext(FrontConfigContext);
    const { data, error, isLoading } = useSWR(
        genBackendPath(config, "/api/v1/operations", { limit: 100 }),
        typedFetcher(OperationLog),
        { refreshInterval: 2000 },
    );
    return (
        <>
            <Typography variant="h5" margin={2}>Operation Log</Typography>
            <TableContainer component={Paper}>
                <Table size="small">
                    <TableHead>
                        <TableRow>
                            <TableCell>Time</TableCell>
                            <TableCell>Source</TableCell>
                            <TableCell>Target</TableCell>
                            <TableCell>Action</TableCell>
                            <TableCell>Arguments</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {data?.slice().reverse().map((entry) => (
                            <TableRow key={`${entry.time}-${entry.target}-${entry.action}`}>
                                <TableCell>{new Date(entry.time).toLocaleString()}</TableCell>
                                <TableCell>{entry.source ?? "unknown"}</TableCell>
                                <TableCell>{entry.target}</TableCell>
                                <TableCell>{entry.action}</TableCell>
                                <TableCell>{JSON.stringify(entry.args)}</TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>
            {isLoading && <Typography margin={2}>Loading...</Typography>}
            {error && <Typography margin={2}>Failed to load operation log. {JSON.stringify(error)} </Typography>}
        </>
    );
}
