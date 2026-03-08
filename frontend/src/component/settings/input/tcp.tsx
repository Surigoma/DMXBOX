import {
    FormControl,
    FormGroup,
    Grid,
    TextField,
    Typography,
} from "@mui/material";
import { Controller, useFormContext } from "react-hook-form";
import NumberField from "../../common/numberField";

function InputTCP() {
    const { control, register } = useFormContext();
    return (
        <Grid direction="column" spacing={1}>
            <Typography variant="h5">TCP</Typography>
            <FormGroup>
                <FormControl margin="normal">
                    <Grid container direction="row" spacing={0}>
                        <Grid size="grow">
                            <TextField
                                fullWidth
                                label="Address"
                                {...register("tcp.ip")}
                            />
                        </Grid>
                        <Grid size={3}>
                            <Controller
                                name={"tcp.port"}
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
            </FormGroup>
        </Grid>
    );
}

export default InputTCP;
