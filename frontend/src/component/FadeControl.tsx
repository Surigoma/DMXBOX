import {
    Button,
    Card,
    CardContent,
    Grid,
    Stack,
    Typography,
} from "@mui/material";
import { ConfigContext, genBackendPath } from "../routes/__root";
import { useContext, useMemo } from "react";
import type { DMXGroupInfo } from "../routes";

function FadeControl({ name, data, showCutin }: { name: string; data: DMXGroupInfo, showCutin: boolean }) {
    const config = useContext(ConfigContext);
    const FadeHeight = useMemo<number>(()=>{return showCutin ? 70 : 100}, [showCutin])
    const CutHeight = useMemo<number>(()=>{return 100 - FadeHeight}, [FadeHeight])
    async function fade(isIn: boolean, cutIn: boolean = false) {
        let opts: {[k: string]:string} = {
            "isIn": String(isIn),
        };
        if (cutIn) {
            opts["interval"] = "0"
            opts["duration"] = "0"
        }
        const path = genBackendPath(config, "/api/v1/fade/" + name, opts);
        console.log(await fetch(path, { method: "POST" }));
    }
    return (
        <Card variant="outlined">
            <CardContent
                style={{ margin: 0, padding: 0, position: "relative", height: "150px"}}
            >
                <div
                    style={{
                        display: "block",
                        position: "absolute",
                        top: 0,
                        left: 0,
                        width: "100%",
                        height: "100%",
                    }}
                >
                    <div
                        style={{
                            display: "flex",
                            width: "100%",
                            height: "100%",
                            alignItems: "center",
                            justifyContent: "center",
                        }}
                    >
                        <Typography variant="h5" component="div">
                            {data.name}
                        </Typography>
                    </div>
                </div>
                <Grid
                    container
                    direction="column"
                    alignItems="stretch"
                    justifyContent="center"
                    width="100%"
                    height="100%"
                    zIndex={10}
                >
                    <Grid size={12} height={FadeHeight + "%"}>
                        <Stack direction="row" spacing={0} height="100%">
                            <Button
                                style={{ width: "100%", height: "100%" }}
                                color="primary"
                                size="large"
                                variant="outlined"
                                onClick={async () => {
                                    await fade(true);
                                }}
                            >
                                Fade In
                            </Button>
                            <Button
                                style={{ width: "100%", height: "100%" }}
                                color="secondary"
                                size="large"
                                variant="outlined"
                                onClick={async () => {
                                    await fade(false);
                                }}
                            >
                                Fade Out
                            </Button>
                        </Stack>
                    </Grid>
                    <Grid size="grow" height={CutHeight + "%"}>
                        <Stack direction="row" spacing={0} height="100%">
                            <Button
                                style={{ width: "100%", height: "100%" }}
                                color="primary"
                                size="large"
                                variant="text"
                                onClick={async () => {
                                    await fade(true, true);
                                }}
                            >
                                Cut In
                            </Button>
                            <Button
                                style={{ width: "100%", height: "100%" }}
                                color="secondary"
                                size="large"
                                variant="text"
                                onClick={async () => {
                                    await fade(false, true);
                                }}
                            >
                                Cut Out
                            </Button>
                        </Stack>
                    </Grid>
                </Grid>
            </CardContent>
        </Card>
    );
}

export default FadeControl;
