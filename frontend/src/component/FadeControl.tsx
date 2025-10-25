import {
    Button,
    ButtonGroup,
    Card,
    CardContent,
    Typography,
} from "@mui/material";

function FadeControl({ group }: { group: string }) {
    return (
        <Card variant="outlined">
            <CardContent>
                <Typography variant="h5" component="div">
                    {group}
                </Typography>
                <ButtonGroup>
                    <Button color="primary" size="large">
                        Fade In
                    </Button>
                    <Button color="secondary" size="large">
                        Fade Out
                    </Button>
                </ButtonGroup>
            </CardContent>
        </Card>
    );
}

export default FadeControl;
