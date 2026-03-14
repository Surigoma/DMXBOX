import {
    Button,
    Card,
    CardHeader,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    FormControl,
    FormGroup,
    Grid,
    IconButton,
    List,
    ListItem,
    TextField,
} from "@mui/material";
import { useEffect, useMemo, useState } from "react";
import { Controller, useFormContext } from "react-hook-form";
import type { DMXGroup } from "../../../types";
import { MdDelete, MdEdit } from "react-icons/md";
import Device from "./device";

interface GroupProp {
    name: string;
}
interface AddGroupResult {
    title: string;
    id: string;
}
interface AddGroupProp {
    name?: string;
    open: boolean;
    onClose: (reuslt: AddGroupResult | undefined, isCancel: boolean) => void;
}

export function AddEditGroup(prop: AddGroupProp) {
    const { getValues } = useFormContext();
    const [title, setTitle] = useState("");
    const [id, setId] = useState("");
    useEffect(() => {
        if (!prop.open || prop.name === undefined) {
            return;
        }
        setId(prop.name.split(".").pop() ?? "");
        setTitle(getValues(prop.name + ".name"));
    }, [prop]);
    return (
        <Dialog open={prop.open} aria-hidden={!prop.open}>
            <DialogTitle>
                Group: {title !== "" ? title : "New Group"}
            </DialogTitle>
            <DialogContent>
                <FormGroup>
                    <FormControl fullWidth margin="normal">
                        <TextField
                            aria-label="Title"
                            value={title}
                            onChange={(e) => setTitle(e.target.value)}
                            label="Title"
                        />
                    </FormControl>
                    <FormControl fullWidth margin="normal">
                        <TextField
                            aria-label="ID"
                            value={id}
                            onChange={(e) => setId(e.target.value)}
                            label="ID"
                        />
                    </FormControl>
                </FormGroup>
            </DialogContent>
            <DialogActions>
                <Button onClick={() => prop.onClose(undefined, true)}>
                    Cancel
                </Button>
                <Button onClick={() => prop.onClose({ title, id }, false)}>
                    {prop.name !== undefined ? "Edit" : "Add"}
                </Button>
            </DialogActions>
        </Dialog>
    );
}

function Group(prop: GroupProp) {
    const { control, getValues, setValue } = useFormContext();
    const parent = "dmx.groups";
    const name = useMemo(() => parent + "." + prop.name, [prop]);
    const [openEdit, setOpenEdit] = useState(false);
    const [openDelete, setOpenDelete] = useState(false);
    return (
        <Card variant="outlined">
            <Grid container direction="column">
                <Controller
                    control={control}
                    name={name}
                    render={({ field }) => {
                        const value = field.value as DMXGroup;
                        if (value === undefined) {
                            return <></>;
                        }
                        return (
                            <Grid
                                container
                                direction="row"
                                justifyContent="center"
                                alignItems="center"
                            >
                                <Grid size="grow">
                                    <CardHeader
                                        title={
                                            value.name + " (" + prop.name + ")"
                                        }
                                    ></CardHeader>
                                </Grid>
                                <Grid size="auto">
                                    <IconButton
                                        sx={{ marginRight: "8px" }}
                                        onClick={() => setOpenDelete(true)}
                                    >
                                        <MdDelete />
                                    </IconButton>
                                </Grid>
                                <Grid size="auto">
                                    <IconButton
                                        sx={{ marginRight: "8px" }}
                                        onClick={() => setOpenEdit(true)}
                                    >
                                        <MdEdit />
                                    </IconButton>
                                </Grid>
                            </Grid>
                        );
                    }}
                />
                <Grid
                    container
                    gap={1}
                    margin={2}
                    direction={{ xs: "column", md: "row" }}
                >
                    <Controller
                        control={control}
                        name={name + ".devices"}
                        render={({ field }) => (
                            <>
                                {(field.value as object[]).map((_, i) => (
                                    <Grid key={i} size="grow">
                                        <Device base={name} index={i} />
                                    </Grid>
                                ))}
                            </>
                        )}
                    />
                </Grid>
                <Button>Add Device</Button>
            </Grid>
            <AddEditGroup
                name={name}
                open={openEdit}
                onClose={(r, c) => {
                    if (c || r === undefined) {
                        setOpenEdit(false);
                        return;
                    }
                    const oldId = prop.name;
                    const newId = r.id;
                    const body = getValues(parent);
                    body[oldId].title = r.title;
                    if (oldId !== r.id) {
                        body[newId] = body[oldId];
                        delete body[oldId];
                    }
                    setValue(parent, body);
                    setOpenEdit(false);
                }}
            />
            <Dialog open={openDelete} aria-hidden={!openDelete}>
                <DialogTitle>
                    Are you sure you want to delete this item?
                </DialogTitle>
                <DialogContent>
                    <List>
                        <ListItem>Title: {getValues(name + ".name")}</ListItem>
                        <ListItem>ID: {prop.name}</ListItem>
                    </List>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setOpenDelete(false)}>Cancel</Button>
                    <Button
                        color="error"
                        onClick={() => {
                            const body = getValues(parent);
                            delete body[prop.name];
                            setValue(parent, body);
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

export default Group;
