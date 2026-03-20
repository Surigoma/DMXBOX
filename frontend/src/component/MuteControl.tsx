import { Button, Card, CardContent, Grid, Stack } from "@mui/material";
import { FrontConfigContext, genBackendPath } from "../routes/__root";
import { useContext } from "react";

function MuteControl() {
    const config = useContext(FrontConfigContext);
    async function mute(isMute: boolean) {
        const path = genBackendPath(config, "/api/v1/mute", {
            isMute,
        });
        await fetch(path, { method: "POST" });
    }
    return (
        <Card variant="outlined">
            <CardContent
                style={{
                    margin: 0,
                    padding: 0,
                    position: "relative",
                    height: "150px",
                }}
            >
                <Grid
                    container
                    direction="column"
                    alignItems="stretch"
                    justifyContent="center"
                    wrap="nowrap"
                    width="100%"
                    height="100%"
                    zIndex={10}
                >
                    <Grid size="grow" height="100%">
                        <Stack direction="row" spacing={0} height="100%">
                            <Button
                                style={{ width: "100%", height: "100%" }}
                                color="error"
                                size="large"
                                variant="outlined"
                                onClick={async () => {
                                    await mute(true);
                                }}
                            >
                                Mute
                            </Button>
                            <Button
                                style={{ width: "100%", height: "100%" }}
                                color="success"
                                size="large"
                                variant="outlined"
                                onClick={async () => {
                                    await mute(false);
                                }}
                            >
                                Unmute
                            </Button>
                        </Stack>
                    </Grid>
                </Grid>
            </CardContent>
        </Card>
    );
}

export default MuteControl;
