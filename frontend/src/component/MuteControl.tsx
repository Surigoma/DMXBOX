import { Button, Card, CardContent, Stack, Typography } from "@mui/material";
import { FrontConfigContext, genBackendPath } from "../routes/__root";
import { useContext, useState } from "react";
import useSWR from "swr";
import { MuteState } from "../types";

function MuteControl() {
    const config = useContext(FrontConfigContext);
    const [commandError, setCommandError] = useState(false);
    const { data, error, mutate } = useSWR(
        genBackendPath(config, "/api/v1/mute-state"),
        async (url: string) => {
            const response = await fetch(url);
            if (!response.ok) throw new Error(`Mute state request failed: ${response.status}`);
            return MuteState.parse(await response.json());
        },
        { refreshInterval: 250, dedupingInterval: 0, errorRetryInterval: 250 },
    );
    const isMute = error ? undefined : data?.isMute;
    async function mute(isMute: boolean) {
        try {
            const response = await fetch(genBackendPath(config, "/api/v1/mute", { isMute }), {
                method: "POST",
                headers: { "X-DMXBOX-Control": "web-ui" },
            });
            if (!response.ok) throw new Error(`Request failed: ${response.status}`);
            setCommandError(false);
        } catch (error) {
            setCommandError(true);
            console.error(error);
        }
        await mutate().catch(() => undefined);
    }
    return (
        <Card variant="outlined" data-testid="MuteControl">
            <CardContent sx={{ p: 0, "&:last-child": { pb: 0 } }}>
                <Stack alignItems="center" justifyContent="center" sx={{ py: 1, minHeight: 64 }}>
                    <Typography role="status" color={isMute === true ? "error" : isMute === false ? "success" : "text.secondary"}>
                        {isMute === true ? "Muted" : isMute === false ? "Unmuted" : "State unavailable"}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">Last sent OSC state</Typography>
                    {commandError && <Typography role="alert" variant="caption" color="error">Failed to send mute command</Typography>}
                </Stack>
                <Stack direction="row" spacing={0} height={104}>
                    <Button sx={{ width: "50%", borderRadius: 0 }} color="error" size="large" variant={isMute === true ? "contained" : "outlined"} aria-pressed={isMute === true} onClick={() => mute(true)}>
                        Mute
                    </Button>
                    <Button sx={{ width: "50%", borderRadius: 0 }} color="success" size="large" variant={isMute === false ? "contained" : "outlined"} aria-pressed={isMute === false} onClick={() => mute(false)}>
                        Unmute
                    </Button>
                </Stack>
            </CardContent>
        </Card>
    );
}

export default MuteControl;
