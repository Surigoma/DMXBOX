import {
    Button,
    ButtonGroup,
    Card,
    CardContent,
    Grid,
    Typography,
} from "@mui/material";

function FadeControl({ group }: { group: string }) {
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
                        <ButtonGroup>
                            <Button color="primary" size="large">
                                Fade In
                            </Button>
                            <Button color="secondary" size="large">
                                Fade Out
                            </Button>
                        </ButtonGroup>
                    </Grid>
                </Grid>
            </CardContent>
        </Card>
    );
}

export default FadeControl;
