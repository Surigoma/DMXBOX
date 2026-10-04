import { Button, Box, Card, CardContent, Stack, Typography } from "@mui/material";
import { alpha } from "@mui/material/styles";
import { FrontConfigContext, genBackendPath } from "../routes/__root";
import { useContext } from "react";
import type { TDMXGroup, TFadeState } from "../types";

function FadeControl({ name, data, showCutin, state, onFade }: {
    name: string;
    data: TDMXGroup;
    showCutin: boolean;
    state?: TFadeState;
    onFade?: () => void;
}) {
    const config = useContext(FrontConfigContext);
    const level = state ? Math.round(state.level * 100) : undefined;
    const status = !state ? "State unavailable"
        : state.state === "waiting" ? "Waiting"
        : state.state === "fading" ? (state.isIn ? "Fading in" : "Fading out")
        : "Idle";

    async function fade(isIn: boolean, cutIn: boolean = false) {
        const opts: { [k: string]: string } = { isIn: String(isIn) };
        if (cutIn) {
            opts["interval"] = "0";
            opts["duration"] = "0";
        }
        const path = genBackendPath(config, "/api/v1/fade/" + name, opts);
        const response = await fetch(path, {
            method: "POST",
            headers: { "X-DMXBOX-Control": "web-ui" },
        });
        if (!response.ok) {
            console.error(`Request failed:${response.status}`);
        } else {
            onFade?.();
        }
    }

    return (
        <Card variant="outlined" data-testid="FadeControl">
            <CardContent sx={{ p: 0, "&:last-child": { pb: 0 } }}>
                <Stack alignItems="center" sx={{ px: 2, py: 1, minHeight: 64, justifyContent: "center" }}>
                    <Typography variant="h5" component="div" sx={{ textAlign: "center", overflowWrap: "anywhere", maxWidth: "100%" }}>
                        {data.name}
                    </Typography>
                    <Typography variant="caption">
                        {status}{level === undefined ? "" : ` · Output ${level}%`}
                    </Typography>
                </Stack>
                <Box sx={{ position: "relative", overflow: "hidden" }}>
                    <Box
                        role="progressbar"
                        aria-label={`${data.name} output level`}
                        aria-valuemin={0}
                        aria-valuemax={100}
                        aria-valuenow={level}
                        aria-valuetext={level === undefined ? "State unavailable" : `${status}, Output ${level}%`}
                        sx={{
                            position: "absolute",
                            inset: 0,
                            width: `${level ?? 0}%`,
                            backgroundColor: (theme) => alpha(theme.palette.primary.main, 0.2),
                            transition: "width 250ms linear",
                            pointerEvents: "none",
                            "@media (prefers-reduced-motion: reduce)": { transition: "none" },
                        }}
                    />
                    <Stack direction="row" spacing={0} sx={{ position: "relative", height: showCutin ? 80 : 104 }}>
                        <Button sx={{ width: "50%", borderRadius: 0 }} color="primary" size="large" variant="outlined" onClick={() => fade(true)}>
                            Fade In
                        </Button>
                        <Button sx={{ width: "50%", borderRadius: 0 }} color="secondary" size="large" variant="outlined" onClick={() => fade(false)}>
                            Fade Out
                        </Button>
                    </Stack>
                </Box>
                {showCutin && (
                    <Stack direction="row" spacing={0} height={44}>
                        <Button sx={{ width: "50%" }} color="primary" size="large" variant="text" onClick={() => fade(true, true)}>
                            Cut In
                        </Button>
                        <Button sx={{ width: "50%" }} color="secondary" size="large" variant="text" onClick={() => fade(false, true)}>
                            Cut Out
                        </Button>
                    </Stack>
                )}
            </CardContent>
        </Card>
    );
}

export default FadeControl;
