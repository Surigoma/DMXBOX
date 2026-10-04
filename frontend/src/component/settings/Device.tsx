import {
    Button,
    FormControl,
    FormGroup,
    Grid,
    Typography,
} from "@mui/material";
import Group, { AddEditGroup } from "./device/group";
import { Controller, useFormContext, useWatch } from "react-hook-form";
import type { TDMXGroupMap } from "../../types";
import { useState } from "react";
import NumberField from "../common/numberField";

function Devices() {
    const { control, getValues, setValue } = useFormContext();
    const [openAdd, setOpenAdd] = useState(false);
    const parent = "dmx.groups";
    const groups = useWatch({ control, name: parent }) as TDMXGroupMap;
    const groupKeys = Object.keys(groups ?? {});
    return (
        <Grid container spacing={2} direction="column" data-testid="Devices">
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
                    const body = getValues(parent) as TDMXGroupMap;
                    setValue(parent, {
                        ...body,
                        [r.id]: { devices: [], name: r.title },
                    });
                    setOpenAdd(false);
                }}
            />
        </Grid>
    );
}

export default Devices;
