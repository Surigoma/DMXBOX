import {
    Button,
    Card,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Grid,
    IconButton,
    List,
    ListItem,
    MenuItem,
    Select,
    Stack,
    Typography,
} from "@mui/material";
import { useMemo, useState } from "react";
import { Controller, useFormContext } from "react-hook-form";
import Dimmer from "./model/dimmer";
import NumberField from "../../common/numberField";
import WCLight from "./model/wclight";
import { MdDelete } from "react-icons/md";
import type { TDMXGroup } from "../../../types";

interface DeviceProp {
    index: number;
    base: string;
}

function ModelSelector(model: string, name: string) {
    switch (model) {
        case "dimmer":
            return <Dimmer name={name} />;
        case "wclight":
            return <WCLight name={name} />;
        default:
            return <div>Undefined</div>;
    }
}

function Device(prop: DeviceProp) {
    const { control, getValues, setValue, watch } = useFormContext();
    const [openDelete, setOpenDelete] = useState(false);
    const name = useMemo(
        () => prop.base + ".devices[" + prop.index + "]",
        [prop],
    );
    const model = watch(name + ".model");

    return (
        <Card style={{ padding: "5px" }} data-testid="DMXDevice">
            <Stack gap={2}>
                <Grid container>
                    <Grid
                        size="grow"
                        justifyContent="center"
                        alignItems="center"
                    >
                        <Typography margin={1} variant="subtitle1">
                            Index: {prop.index + 1}
                        </Typography>
                    </Grid>
                    <Grid>
                        <IconButton
                            onClick={() => setOpenDelete(true)}
                            data-testid="DeviceDeleteButton"
                        >
                            <MdDelete />
                        </IconButton>
                    </Grid>
                </Grid>
                <Controller
                    key="model"
                    control={control}
                    name={name + ".model"}
                    render={({ field }) => (
                        <Select
                            onChange={(e) => {
                                field.onChange(e.target.value);
                            }}
                            value={field.value}
                            data-testid="OpModelSelect"
                        >
                            <MenuItem value="dimmer">Dimmer</MenuItem>
                            <MenuItem value="wclight">
                                White Control Light
                            </MenuItem>
                        </Select>
                    )}
                />
                <Controller
                    key="channel"
                    name={name + ".channel"}
                    control={control}
                    render={({ field }) => (
                        <NumberField
                            label="Start Channel"
                            value={field.value}
                            min={1}
                            max={255}
                            format={{ useGrouping: false }}
                            onValueChange={(e) => field.onChange(e?.valueOf())}
                        ></NumberField>
                    )}
                ></Controller>
                {ModelSelector(model, name)}
            </Stack>
            <Dialog
                open={openDelete}
                aria-hidden={!openDelete}
                data-testid="DeviceDeleteDialog"
            >
                <DialogTitle>
                    Are you sure you want to delete this item?
                </DialogTitle>
                <DialogContent>
                    <List>
                        <ListItem>Index: {prop.index + 1}</ListItem>
                        <ListItem>Model: {getValues(name + ".model")}</ListItem>
                    </List>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setOpenDelete(false)}>Cancel</Button>
                    <Button
                        color="error"
                        onClick={() => {
                            const body = getValues(prop.base) as TDMXGroup;
                            body.devices.splice(prop.index, 1);
                            setValue(prop.base, body);
                            console.log(body);
                            setOpenDelete(false);
                        }}
                    >
                        Confirm
                    </Button>
                </DialogActions>
            </Dialog>
        </Card>
    );
}
export default Device;
