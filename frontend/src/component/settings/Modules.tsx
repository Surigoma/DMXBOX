import { Grid, InputLabel, Typography } from "@mui/material";
import type { InputModules } from "../../types";
import Checked from "../common/checked";
import { useState } from "react";

interface ModulesParam {
    config: InputModules;
}
function Modules(param: ModulesParam) {
    const [http, setHTTP] = useState(param.config.http);
    const [TCP, setTCP] = useState(param.config.tcp);

    return (
        <Grid>
            <Typography variant="h4">Input</Typography>
            <Grid container margin={2} spacing={2} alignItems="center">
                <InputLabel>
                    <Typography variant="h5">Modules</Typography>
                </InputLabel>
                <Checked title="HTTP" check={http} onCheck={setHTTP}></Checked>
                <Checked title="TCP" check={TCP} onCheck={setTCP}></Checked>
            </Grid>
        </Grid>
    );
}

export default Modules;
