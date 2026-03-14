import { Slider, Stack } from "@mui/material";
import { useEffect } from "react";
import { Controller, useFormContext } from "react-hook-form";
import { MdLightbulb, MdLightbulbOutline } from "react-icons/md";

interface DimmerProp {
    name: string;
}
function Dimmer(prop: DimmerProp) {
    const { control, getValues, setValue } = useFormContext();
    useEffect(() => {
        const value = getValues(prop.name + ".max") as number[];
        if (value.length > 1) {
            setValue(prop.name + ".max", [value[0]]);
        }
    });
    return (
        <Stack spacing={2} direction="row" sx={{ alignItems: "center", mb: 1 }}>
            <MdLightbulb />
            <Controller
                name={prop.name + ".max[0]"}
                control={control}
                render={({ field }) => (
                    <Slider
                        aria-label="Dimmer"
                        min={0}
                        max={255}
                        value={field.value}
                        onChange={(e) => field.onChange(e?.valueOf())}
                    />
                )}
            ></Controller>
            <MdLightbulbOutline />
        </Stack>
    );
}
export default Dimmer;
