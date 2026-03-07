import useSWR from "swr";
import {
    fetcher,
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

function OutputDMX() {
    const { control } = useFormContext();
    const config = useContext(FrontConfigContext);
    const { data, error, isLoading } = useSWR<string[]>(
        genBackendPath(config, "/api/v1/config/console"),
        fetcher,
    );
    if (isLoading) {
        return <Alert severity="error">Failed to get Console ports.</Alert>;
    }
    if (data === undefined || error) {
        return <Alert severity="error">Failed to get Console ports.</Alert>;
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
