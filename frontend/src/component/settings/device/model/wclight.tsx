import {
    Box,
    Slider,
    Stack,
    Typography,
    type SxProps,
    type Theme,
} from "@mui/material";
import { useEffect, useMemo, useState } from "react";
import { useFormContext } from "react-hook-form";
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
    const { setValue, getValues } = useFormContext();
    const [colorTemp, setColorTemp] = useState(0.0);
    const [dimmer, setDimmer] = useState(0.0);
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
    let setTimer: number | undefined = undefined;

    function convertDMXtoWCInfo(values: number[]): WCInfo {
        if (values === undefined || values.length < 3) {
            return {
                dimmer: values[0] !== undefined ? values[0] / 255 : 1,
                temp: 0.5,
            };
        }
        const target = values.slice(0, 2);
        const cool = target[0];
        const warm = target[1];
        return {
            dimmer: Math.max(...target) / 255,
            temp: cool / (cool + warm),
        };
    }
    function convertWCInfotoDMX(values: WCInfo): number[] {
        const warm = Math.round(values.dimmer * values.temp * 255);
        const cool = Math.round(255 * values.dimmer - warm);
        return [cool, warm, 0];
    }
    useEffect(() => {
        const values = getValues(prop.name + ".max") as number[];
        const wcinfo = convertDMXtoWCInfo(values);
        setDimmer(wcinfo.dimmer);
        setColorTemp(wcinfo.temp);
    }, []);
    function updateValues() {
        const value: WCInfo = {
            dimmer: dimmer,
            temp: colorTemp,
        };
        const newValue = convertWCInfotoDMX(value);
        clearTimeout(setTimer);
        setTimer = setTimeout(()=>{
            setValue(prop.name + ".max", newValue);
            setTimer = undefined;
        }, 100);
    }
    return (
        <Stack spacing={2}>
            <Stack
                spacing={2}
                direction="row"
                sx={{ alignItems: "center", mb: 1 }}
            >
                <MdLightbulb />
                <Slider
                    aria-label="Dimmer"
                    min={0}
                    max={1}
                    step={0.01}
                    value={dimmer}
                    onChange={(_, v) => {
                        setDimmer(v);
                        updateValues();
                    }}
                />
                <MdLightbulbOutline />
                <Typography
                    variant="caption"
                    noWrap={true}
                    width="48px"
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
                    min={0}
                    max={1}
                    step={0.01}
                    value={colorTemp}
                    onChange={(_, v) => {
                        setColorTemp(v);
                        updateValues();
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
