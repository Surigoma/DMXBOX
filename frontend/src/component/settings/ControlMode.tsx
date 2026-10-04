import { Alert, Button, FormControlLabel, Stack, Switch, Typography } from "@mui/material";
import { useContext, useState } from "react";
import { FrontConfigContext, genBackendPath } from "../../routes/__root";
import { ControlMode } from "../../types";
import { useControlMode } from "../../contexts/controlMode";

function ControlModeSettings() {
    const config = useContext(FrontConfigContext);
    const { data, error, isLoading, mutate } = useControlMode();
    const [updating, setUpdating] = useState(false);
    const [updateError, setUpdateError] = useState<string>();

    async function setBrowserOnly(browserOnly: boolean) {
        setUpdating(true);
        setUpdateError(undefined);
        try {
            const response = await fetch(genBackendPath(config, "/api/v1/control-mode"), {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "X-DMXBOX-Control": "web-ui",
                },
                body: JSON.stringify({ browserOnly }),
            });
            if (!response.ok) throw new Error(`Request failed: ${response.status}`);
            await mutate(ControlMode.parse(await response.json()), { revalidate: false });
        } catch (error) {
            setUpdateError(`Failed to change control mode. ${error instanceof Error ? error.message : "Please try again."}`);
            // The server may have applied the change before its response was lost.
            await mutate().catch(() => undefined);
        } finally {
            setUpdating(false);
        }
    }

    return (
        <Stack spacing={1}>
            {error && (
                <Alert severity="error" action={<Button onClick={() => { void mutate().catch(() => undefined); }}>Retry</Button>}>
                    Failed to load control mode.
                </Alert>
            )}
            {updateError && <Alert severity="error" onClose={() => setUpdateError(undefined)}>{updateError}</Alert>}
            <FormControlLabel
                label="Browser only"
                control={
                    <Switch
                        checked={data?.browserOnly ?? false}
                        disabled={isLoading || !!error || !data || updating}
                        onChange={(e) => { void setBrowserOnly(e.target.checked); }}
                        slotProps={{ input: { "aria-describedby": "browser-only-description" } }}
                    />
                }
            />
            <Typography id="browser-only-description" variant="body2" color="text.secondary">
                Use browser controls for manual operation and block commands from TCP and other external controllers.
                Changes take effect immediately, without pressing Update. This mode resets when the backend restarts.
            </Typography>
        </Stack>
    );
}

export default ControlModeSettings;
