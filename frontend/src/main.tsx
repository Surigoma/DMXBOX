import { createRoot } from "react-dom/client";
import "./index.css";
import { Fragment, StrictMode } from "react";
import { createRouter, RouterProvider } from "@tanstack/react-router";

// Import the generated route tree
import { routeTree } from "./routeTree.gen";
import { ThemeProvider } from "@emotion/react";
import { createTheme, CssBaseline } from "@mui/material";

// Create a new router instance
const router = createRouter({ routeTree });
const darkTheme = createTheme({
    palette: {
        mode: "dark",
    },
});

createRoot(document.getElementById("root")!).render(
    <Fragment>
        <ThemeProvider theme={darkTheme}>
            <CssBaseline enableColorScheme />
            <StrictMode>
                <RouterProvider router={router} basepath="/gui" />
            </StrictMode>
        </ThemeProvider>
    </Fragment>,
);
