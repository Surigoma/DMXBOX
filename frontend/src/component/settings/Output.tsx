import { Grid, InputLabel, Typography } from "@mui/material";
import type { OutputTargets } from "../../types";
import Checked from "../common/checked";
import { useState } from "react";

interface OutputsParam {
    config: OutputTargets;
}
function Outputs(param: OutputsParam) {
    const [outputTargets, setOutputTargets] = useState<string[]>(
        param.config.target,
    );
    //TODO: Please optimize callbacks.
    const [FTDI, setFTDI] = useState(outputTargets.includes("ftdi"));
    const [Artnet, setArtnet] = useState(outputTargets.includes("artnet"));

    function toggleOutputs(target: string, checked: boolean) {
        let tmp = outputTargets;
        if (checked) {
            if (outputTargets.includes(target)) {
                return tmp;
            }
            tmp.push(target);
        } else {
            if (!outputTargets.includes(target)) {
                return tmp;
            }
            const i = tmp.indexOf(target);
            tmp.splice(i, 1);
        }
        setFTDI(tmp.includes("ftdi"));
        setArtnet(tmp.includes("artnet"));
        return tmp;
    }
    return (
        <Grid>
            <Typography variant="h4">Output</Typography>
            <Grid container alignItems="center" spacing={2} margin={2}>
                <InputLabel id="demo-multiple-chip-label">
                    <Typography variant="h5">Target</Typography>
                </InputLabel>
                <Checked
                    title="FTDI"
                    check={FTDI}
                    onCheck={(c) => setOutputTargets(toggleOutputs("ftdi", c))}
                ></Checked>
                <Checked
                    title="Artnet"
                    check={Artnet}
                    onCheck={(c) =>
                        setOutputTargets(toggleOutputs("artnet", c))
                    }
                ></Checked>
            </Grid>
        </Grid>
    );
}

export default Outputs;
