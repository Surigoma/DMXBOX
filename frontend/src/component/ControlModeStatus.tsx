import { Chip } from "@mui/material";
import { useControlMode } from "../contexts/controlMode";

export default function ControlModeStatus() {
    const { data, error } = useControlMode();
    const label = error ? "Unavailable" : !data ? "Loading" : data.browserOnly ? "ON" : "OFF";
    return <Chip role="status" aria-label={`Browser only: ${label}`} size="small" variant="outlined" color={!error && data?.browserOnly ? "warning" : "default"} label={`Browser only: ${label}`} />;
}
