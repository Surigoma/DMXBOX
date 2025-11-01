import {
    Button,
    Card,
    CardContent,
    Grid,
    Stack,
    Typography,
} from "@mui/material";
import { ConfigContext, genBackendPath } from "../routes/__root";
import { useContext } from "react";
import type { DMXGroupInfo } from "../routes";

function FadeControl({ name, data }: { name: string, data: DMXGroupInfo }) {
    const config = useContext(ConfigContext);
    async function fade(isIn: boolean) {
        const path = genBackendPath(config, "/api/v1/fade/" + name, {"isIn": isIn});
        console.log(await fetch(path, {method: "POST"}));
    }
    return (
        <Card variant="outlined">
            <CardContent>
                <Grid
                    container
                    direction="column"
                    spacing={2}
                    alignItems="center"
                    justifyContent="center"
                >
                    <Grid>
                        <Typography variant="h5" component="div">
                            {data.name}
                        </Typography>
                    </Grid>
                    <Grid>
                        <Stack direction="row" spacing={2}>
                            <Button color="primary" size="large" variant="outlined" onClick={async ()=>{await fade(true)}}>
                                Fade In
                            </Button>
                            <Button color="secondary" size="large" variant="outlined" onClick={async ()=>{await fade(false)}}>
                                Fade Out
                            </Button>
                        </Stack>
                    </Grid>
                </Grid>
            </CardContent>
        </Card>
    );
}

export default FadeControl;
