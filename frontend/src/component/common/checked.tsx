import { Card, Grid, InputLabel, Switch, Typography } from "@mui/material";

interface CheckedParam {
    title: string;
    check: boolean;
    onCheck: (isChecked: boolean) => void;
}

function Checked(param: CheckedParam) {
    return (
        <Card
            variant="outlined"
            onClick={() => {
                param.onCheck(!param.check);
            }}
        >
            <Grid
                container
                spacing={2}
                justifyContent="space-around"
                alignItems="center"
            >
                <Grid size="grow">
                    <InputLabel>
                        <Typography
                            margin={2}
                            style={{
                                userSelect: "none",
                                wordBreak: "break-all",
                            }}
                        >
                            {param.title}
                        </Typography>
                    </InputLabel>
                </Grid>
                <Grid size="auto" justifyContent="flex-end">
                    <Switch
                        checked={param.check}
                        onChange={(e) => {
                            param.onCheck(e.target.checked);
                        }}
                    ></Switch>
                </Grid>
            </Grid>
        </Card>
    );
}

export default Checked;
