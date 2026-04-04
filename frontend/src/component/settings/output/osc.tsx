import {
    Autocomplete,
    FormControl,
    FormControlLabel,
    FormGroup,
    Grid,
    InputLabel,
    MenuItem,
    Select,
    Switch,
    TextField,
    Typography,
} from "@mui/material";
import { Controller, useFormContext } from "react-hook-form";
import NumberField from "../../common/numberField";

function OutputOSC() {
    const { control, register } = useFormContext();
    const oscChannels = Array.from({ length: 255 }, (_, i) =>
        (i + 1).toString(),
    );
    return (
        <Grid spacing={2} data-testid="OutputOSC">
            <FormGroup>
                <Typography variant="h5">OSC</Typography>
                <FormControl fullWidth margin="normal">
                    <Grid container direction="row" spacing={0}>
                        <Grid size="grow">
                            <TextField
                                fullWidth
                                label="Address"
                                {...register("osc.ip")}
                            />
                        </Grid>
                        <Grid size={3}>
                            <Controller
                                name={"osc.port"}
                                control={control}
                                render={({ field }) => (
                                    <NumberField
                                        label="Port"
                                        value={field.value}
                                        min={1}
                                        max={65535}
                                        format={{ useGrouping: false }}
                                        onValueChange={(e) =>
                                            field.onChange(e?.valueOf())
                                        }
                                    ></NumberField>
                                )}
                            ></Controller>
                        </Grid>
                    </Grid>
                </FormControl>
                <FormControl fullWidth margin="normal">
                    <TextField
                        fullWidth
                        label="OSC Path format"
                        {...register("osc.format")}
                    />
                </FormControl>
                <Grid container spacing={2} alignItems="center">
                    <Grid size="grow">
                        <FormControl fullWidth margin="normal">
                            <InputLabel id="osc-type">
                                Sending data type
                            </InputLabel>
                            <Controller
                                control={control}
                                name="osc.type"
                                render={({ field }) => (
                                    <Select
                                        labelId="osc-type"
                                        label="Sending data type"
                                        data-testid="OpSendingDataType"
                                        value={field.value}
                                        onChange={(e) => {
                                            field.onChange(e.target.value);
                                        }}
                                    >
                                        <MenuItem value="float">Float</MenuItem>
                                        <MenuItem value="int">Int</MenuItem>
                                    </Select>
                                )}
                            />
                        </FormControl>
                    </Grid>
                    <Grid size="auto">
                        <FormControlLabel
                            label="Inverse"
                            style={{ userSelect: "none" }}
                            control={
                                <Switch {...register("osc.inverse")}>
                                    Inverse
                                </Switch>
                            }
                        />
                    </Grid>
                </Grid>
                <FormControl margin="normal">
                    <Controller
                        control={control}
                        name="osc.channels"
                        render={({ field }) => (
                            <Autocomplete
                                multiple
                                value={(field.value as number[])
                                    .sort((a, b) => a - b)
                                    .map((v) => v.toString())}
                                options={oscChannels}
                                onChange={(_, v) => {
                                    field.onChange(
                                        v
                                            .map((v: string) => Number(v))
                                            .sort((a, b) => a - b),
                                    );
                                }}
                                renderInput={(p) => (
                                    <TextField
                                        label="Target OSC Channels"
                                        {...p}
                                    />
                                )}
                            />
                        )}
                    ></Controller>
                </FormControl>
            </FormGroup>
        </Grid>
    );
}

export default OutputOSC;
