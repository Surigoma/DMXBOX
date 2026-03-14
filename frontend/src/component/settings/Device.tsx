import { Button, Grid } from "@mui/material";
import Group, { AddEditGroup } from "./device/group";
import { Controller, useFormContext } from "react-hook-form";
import type { DMXGroup } from "../../types";
import { useState } from "react";

function Devices() {
    const { control, getValues, setValue } = useFormContext();
    const [openAdd, setOpenAdd] = useState(false);
    const parent = "dmx.groups";
    return (
        <Grid container spacing={2} direction="column">
            <Controller
                control={control}
                name={parent}
                render={({ field }) => (
                    <>
                        {field.value ? (
                            Object.keys(
                                field.value as { [name: string]: DMXGroup },
                            ).map((v) => <Group key={v} name={v} />)
                        ) : (
                            <>No group</>
                        )}
                    </>
                )}
            />
            <Button onClick={() => setOpenAdd(true)}>Add Group</Button>
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
