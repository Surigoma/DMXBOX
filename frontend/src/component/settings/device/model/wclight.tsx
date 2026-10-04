import {
    Box,
    Slider,
    Stack,
    Typography,
    type SxProps,
    type Theme,
} from "@mui/material";
import { useEffect, useMemo, useState } from "react";
import { useFormContext, useWatch } from "react-hook-form";
import { MdLightbulb, MdLightbulbOutline } from "react-icons/md";

interface WCLightProp {
    name: string;
}
interface WCInfo {
    dimmer: number;
    temp: number;
}
function WCLight(prop: WCLightProp) {
    const colorPalette = useMemo(() => {
        return { cool: "#add8e6", warm: "#ffffe0" };
    }, []);
    const style: SxProps<Theme> = {
        width: "1em",
        height: "1em",
        borderRadius: 1,
    };

    function convertDMXtoWCInfo(values: number[]): WCInfo {
        if (values == null) {
            return {
                dimmer: 1,
                temp: 0.5
            }
        }
        if (values.length < 3) {
            return {
                dimmer: values[0] != null ? values[0] / 255 : 1,
                temp: 0.5,
            };
        }
        const target = values.slice(0, 2);
        const cool = target[0];
        const warm = target[1];
        const dimmer = Math.max(...target) / 255;
        const temp = warm / (cool + warm);
        return {
            dimmer: dimmer,
            temp: !isNaN(temp) ? temp : 0.5,
        };
    }
    function convertWCInfoToDMX(values: WCInfo): number[] {
        const scale = 255 * values.dimmer / Math.max(values.temp, 1 - values.temp);
        return [Math.round(scale * (1 - values.temp)), Math.round(scale * values.temp), 0];
    }

    const { setValue, control } = useFormContext();
    const values = useWatch({ control, name: prop.name + ".max" }) as number[];
    const signature = prop.name + JSON.stringify(values);
    const [edited, setEdited] = useState<{ signature: string; info: WCInfo }>();
    const info = edited?.signature === signature ? edited.info : convertDMXtoWCInfo(values);
    const { dimmer, temp: colorTemp } = info;
    function changeInfo(next: WCInfo) {
        const max = convertWCInfoToDMX(next);
        setEdited({ signature: prop.name + JSON.stringify(max), info: next });
        setValue(prop.name + ".max", max, { shouldDirty: true });
    }
    const colorMix = useMemo(
        () =>
            "color-mix(" +
            [
                "in srgb",
                [colorPalette.cool, (1 - colorTemp) * 100 + "%"].join(" "),
                [colorPalette.warm, colorTemp * 100 + "%"].join(" "),
            ].join(",") +
            ")",
        [colorTemp, colorPalette],
    );
    useEffect(() => {
        if (!values || values.length < 3) {
            const level = values?.[0] ?? 255;
            setValue(prop.name + ".max", [level, level, 0]);
        }
    }, [values, prop.name, setValue]);

    return (
        <Stack spacing={2} data-testid="WCLight">
            <Stack
                spacing={2}
                direction="row"
                sx={{ alignItems: "center", mb: 1 }}
            >
                <MdLightbulb />
                <Slider
                    aria-label="Dimmer"
                    data-testid="OpDimmer"
                    min={0}
                    max={1}
                    step={0.01}
                    value={dimmer}
                    onChange={(_, v) => {
                        changeInfo({ dimmer: v as number, temp: colorTemp });
                    }}
                />
                <MdLightbulbOutline />
                <Typography
                    variant="caption"
                    noWrap={true}
                    width="50px"
                    textAlign="right"
                >
                    {(dimmer * 100).toFixed(0)} %
                </Typography>
            </Stack>
            <Stack
                spacing={2}
                direction="row"
                sx={{ alignItems: "center", mb: 1 }}
            >
                <Box
                    sx={{
                        backgroundColor: colorPalette.cool,
                        ...style,
                    }}
                />
                <Slider
                    aria-label="Temp"
                    data-testid="OpTemp"
                    min={0}
                    max={1}
                    step={0.01}
                    value={colorTemp}
                    onChange={(_, v) => {
                        changeInfo({ dimmer, temp: v as number });
                    }}
                />
                <Box
                    sx={{
                        backgroundColor: colorPalette.warm,
                        ...style,
                    }}
                />
                <Box
                    sx={{
                        backgroundColor: colorMix,
                        width: "48px",
                        height: "16px",
                        borderRadius: 1,
                    }}
                />
            </Stack>
        </Stack>
    );
}
export default WCLight;
