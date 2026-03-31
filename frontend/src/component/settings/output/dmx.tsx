import useSWR from "swr";
import {
    typedFetcher,
    FrontConfigContext,
    genBackendPath,
} from "../../../routes/__root";
import { useContext } from "react";
import {
    Alert,
    FormControl,
    Grid,
    InputLabel,
    MenuItem,
    Select,
    Typography,
} from "@mui/material";
import { Controller, useFormContext } from "react-hook-form";
import { ConsoleAPIResult } from "../../../types";

function OutputDMX() {
    const { control } = useFormContext();
    const config = useContext(FrontConfigContext);
    const { data, error, isLoading } = useSWR<string[]>(
        genBackendPath(config, "/api/v1/config/console"),
        typedFetcher(ConsoleAPIResult),
    );
    if (isLoading) {
        return "Loading...";
    }
    if (data === undefined || error) {
        return <Alert severity="error">Failed to get Console ports.</Alert>;
    }
    if (data.length <= 0) {
        return <Alert severity="warning">DMX Port is not found.</Alert>;
    }
    return (
        <Grid spacing={2}>
            <Typography variant="h5">DMX</Typography>
            <FormControl fullWidth margin="normal">
                <InputLabel id="output-dmx-port">Port</InputLabel>
                <Controller
                    name="output.dmx.port"
                    control={control}
                    render={({ field }) => (
                        <Select
                            value={field.value}
                            labelId="output-dmx-port"
                            label="Port"
                            onChange={(e) => {
                                field.onChange(e.target.value);
                            }}
                        >
                            {data.map((v) => (
                                <MenuItem key={v} value={v}>
                                    {v}
                                </MenuItem>
                            ))}
                        </Select>
                    )}
                />
            </FormControl>
        </Grid>
    );
}

export default OutputDMX;
