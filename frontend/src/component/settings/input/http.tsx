import {
    FormControl,
    FormGroup,
    Grid,
    IconButton,
    Paper,
    Stack,
    styled,
    TextField,
    Typography,
} from "@mui/material";
import { IoMdTrash } from "react-icons/io";
import { Controller, useFormContext, Watch } from "react-hook-form";
import NumberField from "../../common/numberField";

const Item = styled(Paper)(({ theme }) => ({
    padding: theme.spacing(1),
    textAlign: "left",
    color: (theme.vars ?? theme).palette.text.secondary,
}));

function InputHTTP() {
    const { control, register, setValue } = useFormContext();
    return (
        <Grid direction="column" spacing={1} data-testid="InputHTTP">
            <Typography variant="h5">HTTP</Typography>
            <FormGroup>
                <FormControl margin="normal">
                    <Grid container direction="row" spacing={0}>
                        <Grid size="grow">
                            <TextField
                                fullWidth
                                data-testid="OpIP"
                                label="Address"
                                {...register("http.ip")}
                            />
                        </Grid>
                        <Grid size={3}>
                            <Controller
                                name={"http.port"}
                                control={control}
                                render={({ field }) => (
                                    <NumberField
                                        label="Port"
                                        name="OpPort"
                                        data-testid="OpPort"
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
                    <Stack spacing={1}>
                        <Controller
                            name={"http.accepts"}
                            control={control}
                            render={({ field }) => (
                                <TextField
                                    fullWidth
                                    name="OpAccepts"
                                    data-testid="OpAccepts"
                                    label="Accept addresses"
                                    onKeyDown={(e) => {
                                        if (e.key === "Enter") {
                                            const target =
                                                e.target as HTMLInputElement;
                                            const newValue = (
                                                field.value as string[]
                                            )
                                                .filter(
                                                    (v) => v !== target.value,
                                                )
                                                .concat(target.value);
                                            field.onChange(newValue);
                                            target.value = "";
                                        }
                                    }}
                                ></TextField>
                            )}
                        />
                        <Watch
                            name={"http.accepts"}
                            control={control}
                            render={(field: string[]) => {
                                return field.map((v) => (
                                    <Item key={v}>
                                        <Grid container direction="row">
                                            <Grid size="grow">{v}</Grid>
                                            <Grid size="auto">
                                                <IconButton
                                                    aria-label="delete"
                                                    size="small"
                                                    onClick={() => {
                                                        setValue(
                                                            "http.accepts",
                                                            field.filter(
                                                                (v2) =>
                                                                    v2 !== v,
                                                            ),
                                                        );
                                                    }}
                                                >
                                                    <IoMdTrash />
                                                </IconButton>
                                            </Grid>
                                        </Grid>
                                    </Item>
                                ));
                            }}
                        />
                    </Stack>
                </FormControl>
            </FormGroup>
        </Grid>
    );
}

export default InputHTTP;
