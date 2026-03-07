import { Grid, InputLabel, Typography } from "@mui/material";
import Checked from "../common/checked";

function Inputs() {
    return (
        <Grid>
            <Typography variant="h4">Input</Typography>
            <Grid container margin={2} spacing={2} alignItems="center">
                <InputLabel>
                    <Typography variant="h5">Modules: </Typography>
                </InputLabel>
                <Checked title="HTTP" target="modules.http"></Checked>
                <Checked title="TCP" target="modules.tcp"></Checked>
            </Grid>
        </Grid>
    );
}

export default Inputs;
