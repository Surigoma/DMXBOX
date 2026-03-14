import { Card, MenuItem, Select, Stack } from "@mui/material";
import { useMemo } from "react";
import { Controller, useFormContext } from "react-hook-form";
import Dimmer from "./model/dimmer";
import NumberField from "../../common/numberField";
import WCLight from "./model/wclight";

interface DeviceProp {
    index: number;
    base: string;
}

function Device(prop: DeviceProp) {
    const { control } = useFormContext();
    const name = useMemo(
        () => prop.base + ".devices[" + prop.index + "]",
        [prop],
    );
    return (
        <Card style={{ padding: "5px" }}>
            <Stack gap={2}>
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
                            min={0}
                            max={255}
                            format={{ useGrouping: false }}
                            onValueChange={(e) => field.onChange(e?.valueOf())}
                        ></NumberField>
                    )}
                ></Controller>
                <Controller
                    key="type_defined"
                    control={control}
                    name={name + ".model"}
                    render={({ field }) => {
                        switch (field.value) {
                            case "dimmer":
                                return <Dimmer name={name} />;
                            case "wclight":
                                return <WCLight name={name} />;
                            default:
                                return <div>Undefined</div>;
                        }
                    }}
                />
            </Stack>
        </Card>
    );
}
export default Device;
