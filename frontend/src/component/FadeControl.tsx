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

function FadeControl({ group }: { group: string }) {
    const config = useContext(ConfigContext);
    async function fade(isIn: boolean) {
        const path = genBackendPath(config, "/api/v1/fade/" + group, {"isIn": isIn});
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
                            {group}
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
