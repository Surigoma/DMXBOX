import {
    Button,
    FormControl,
    FormGroup,
    Grid,
    Typography,
} from "@mui/material";
import Group, { AddEditGroup } from "./device/group";
import { Controller, useFormContext } from "react-hook-form";
import type { DMXGroup } from "../../types";
import { useMemo, useState } from "react";
import NumberField from "../common/numberField";

function Devices() {
    const { control, getValues, setValue, watch } = useFormContext();
    const [openAdd, setOpenAdd] = useState(false);
    const parent = "dmx.groups";
    const groups = watch(parent) as { [key: string]: DMXGroup };
    const groupKeys = useMemo(() => Object.keys(groups ?? {}), [groups]);
    return (
        <Grid container spacing={2} direction="column">
            <Grid>
                <Typography variant="h6">Global Options</Typography>
                <FormGroup>
                    <Grid
                        container
                        direction={{ xs: "column", md: "row" }}
                        gap={2}
                    >
                        <FormControl margin="normal">
                            <Controller
                                control={control}
                                name="dmx.fps"
                                render={({ field }) => (
                                    <NumberField
                                        label="Update FPS"
                                        value={field.value ?? 40}
                                        min={1}
                                        max={45}
                                        format={{ useGrouping: false }}
                                        onValueChange={(e) =>
                                            field.onChange(e?.valueOf())
                                        }
                                    ></NumberField>
                                )}
                            />
                        </FormControl>
                        <FormControl margin="normal">
                            <Controller
                                control={control}
                                name="dmx.fadeInterval"
                                render={({ field }) => (
                                    <NumberField
                                        label="Fade Interval"
                                        value={field.value ?? 0.7}
                                        step={0.1}
                                        format={{ useGrouping: false }}
                                        help="Interval for fade action"
                                        onValueChange={(e) =>
                                            field.onChange(e?.valueOf())
                                        }
                                    ></NumberField>
                                )}
                            />
                        </FormControl>
                        <FormControl margin="normal">
                            <Controller
                                control={control}
                                name="dmx.delay"
                                render={({ field }) => (
                                    <NumberField
                                        label="Delay"
                                        value={field.value ?? 0}
                                        step={0.1}
                                        format={{ useGrouping: false }}
                                        help="Delay of before fade action"
                                        onValueChange={(e) =>
                                            field.onChange(e?.valueOf())
                                        }
                                    ></NumberField>
                                )}
                            />
                        </FormControl>
                    </Grid>
                </FormGroup>
            </Grid>
            <Grid>
                <Typography variant="h6">Groups</Typography>
                {groupKeys.length > 0 ? (
                    groupKeys.map((v) => <Group key={v} name={v} />)
                ) : (
                    <>No Groups</>
                )}
                <FormControl fullWidth>
                    <Button onClick={() => setOpenAdd(true)}>Add Group</Button>
                </FormControl>
            </Grid>
            <AddEditGroup
                open={openAdd}
                onClose={(r, c) => {
                    if (r === undefined || c) {
                        setOpenAdd(false);
                        return;
                    }
                    const body = getValues(parent) as {
                        [name: string]: DMXGroup;
                    };
                    body[r.id] = {
                        devices: [],
                        name: r.title,
                    };
                    setValue(parent, body);
                    setOpenAdd(false);
                }}
            />
        </Grid>
    );
}

export default Devices;
