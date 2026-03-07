import {
    FormControl,
    FormGroup,
    Grid,
    TextField,
    Typography,
} from "@mui/material";
import { Controller, useFormContext } from "react-hook-form";
import NumberField from "../../common/numberField";

interface UniverseOption {
    name: string;
    title: string;
    min: number;
    max: number;
}

function OutputArtnet() {
    const { register, control } = useFormContext();
    const universeMap: UniverseOption[] = [
        {
            name: "universe",
            title: "Universe",
            min: 0,
            max: 15,
        },
        {
            name: "subuni",
            title: "Sub Universe",
            min: 0,
            max: 15,
        },
        {
            name: "net",
            title: "Net",
            min: 0,
            max: 15,
        },
    ];
    return (
        <Grid spacing={2}>
            <Typography variant="h5">Artnet</Typography>
            <FormGroup>
                <FormControl margin="normal">
                    <TextField
                        label="Address"
                        {...register("output.artnet.addr")}
                    />
                </FormControl>
                <Grid container direction="row" spacing={1}>
                    {universeMap.map((e) => (
                        <Grid size="grow" key={e.name}>
                            <FormControl fullWidth margin="normal">
                                <Controller
                                    name={"output.artnet." + e.name}
                                    control={control}
                                    render={({ field }) => (
                                        <NumberField
                                            label={e.title}
                                            value={field.value}
                                            min={e.min}
                                            max={e.max}
                                            onValueChange={(e) =>
                                                field.onChange(e?.valueOf())
                                            }
                                        ></NumberField>
                                    )}
                                ></Controller>
                            </FormControl>
                        </Grid>
                    ))}
                </Grid>
            </FormGroup>
        </Grid>
    );
}

export default OutputArtnet;
