import { Card, Grid, InputLabel, Typography } from "@mui/material";
import Checked from "../common/checked";
import { useFormContext, Watch } from "react-hook-form";
import InputHTTP from "./input/http";
import InputTCP from "./input/tcp";

function Inputs() {
    const { control } = useFormContext();
    return (
        <Grid container spacing={2} direction="column" data-testid="Inputs">
            <Grid container margin={2} spacing={2} alignItems="center">
                <InputLabel>
                    <Typography variant="h5">Modules: </Typography>
                </InputLabel>
                <Checked title="HTTP" target="modules.http"></Checked>
                <Checked title="TCP" target="modules.tcp"></Checked>
            </Grid>
            <Card variant="outlined" style={{ margin: "2px" }}>
                <Grid
                    container
                    spacing={2}
                    margin={2}
                    direction={{ xs: "column", md: "row" }}
                >
                    <Watch
                        control={control}
                        name={["modules.http", "modules.tcp"]}
                        render={(v) =>
                            v[0] || v[1] ? (
                                [
                                    v[0] ? (
                                        <Grid size="grow" key="http">
                                            <InputHTTP />
                                        </Grid>
                                    ) : undefined,
                                    v[1] ? (
                                        <Grid size="grow" key="tcp">
                                            <InputTCP />
                                        </Grid>
                                    ) : undefined,
                                ]
                            ) : (
                                <>Not selected</>
                            )
                        }
                    />
                </Grid>
            </Card>
        </Grid>
    );
}

export default Inputs;
