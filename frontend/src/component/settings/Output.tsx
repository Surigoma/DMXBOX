import { Card, FormGroup, Grid, InputLabel, Typography } from "@mui/material";
import Checked from "../common/checked";
import { useFormContext, Watch } from "react-hook-form";
import type { ReactElement } from "react";
import OutputDMX from "./output/dmx";
import OutputArtnet from "./output/artnet";

function Outputs() {
    const { control } = useFormContext();
    return (
        <Grid>
            <Typography variant="h4">Output</Typography>
            <Grid container alignItems="center" spacing={2} margin={2}>
                <InputLabel id="demo-multiple-chip-label">
                    <Typography variant="h5">Target: </Typography>
                </InputLabel>
                <FormGroup row>
                    <Checked
                        value="ftdi"
                        title="FTDI"
                        target="output.target"
                    ></Checked>
                    <Checked
                        value="artnet"
                        title="Artnet"
                        target="output.target"
                    ></Checked>
                </FormGroup>
            </Grid>
            <Watch
                control={control}
                name={"output.target"}
                render={(v) => {
                    if (!(v instanceof Array)) {
                        return <></>;
                    }
                    let result: ReactElement[] = [];
                    if (v.includes("ftdi")) {
                        result.push(<OutputDMX key="ftdi"></OutputDMX>);
                    }
                    if (v.includes("artnet")) {
                        result.push(<OutputArtnet key="artnet"></OutputArtnet>);
                    }
                    return (
                        <Card variant="outlined" style={{ margin: "2px" }}>
                            <Grid
                                container
                                spacing={2}
                                margin={2}
                                direction={{ xs: "column", md: "row" }}
                            >
                                {result.length > 0 ? (
                                    result.map((e) => (
                                        <Grid size="grow" key={e.key}>
                                            {e}
                                        </Grid>
                                    ))
                                ) : (
                                    <>Not selected</>
                                )}
                            </Grid>
                        </Card>
                    );
                }}
            />
        </Grid>
    );
}

export default Outputs;
