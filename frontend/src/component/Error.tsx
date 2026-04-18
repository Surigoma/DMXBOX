import { Alert, Container } from "@mui/material";
import type { ReactNode } from "react";

function ErrorComponent({ children }: { children: ReactNode }) {
    return (
        <Container style={{ padding: "10px" }} data-testid="ErrorComponent">
            <Alert severity="error" variant="filled">
                {children}
            </Alert>
        </Container>
    );
}

export default ErrorComponent;
