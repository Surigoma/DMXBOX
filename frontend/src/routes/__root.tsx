import {
    createRootRoute,
    Link,
    Outlet,
    useLocation,
} from "@tanstack/react-router";
import { createContext, useContext, useState } from "react";
import AppBar from "@mui/material/AppBar";
import Toolbar from "@mui/material/Toolbar";
import IconButton from "@mui/material/IconButton";
import { MdMenu } from "react-icons/md";
import Typography from "@mui/material/Typography";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import Configuration from "../contexts/config";

const Config = new Configuration();
export const ConfigContext = createContext(Config.body);

export function fetcher(url: string) {
    return fetch(url, { credentials: "include", mode: "cors" }).then((r) =>
        r.json(),
    );
}

export function genBackendPath(path: string): string {
    const config = useContext(ConfigContext);
    return (
        "http://" + window.location.hostname + ":" + config.backendPort + path
    );
}

interface PathInfo {
    title: string;
}

function RootLayout() {
    const PageInformation: { [path: string]: PathInfo } = {
        "/": {
            title: "Control",
        },
    };
    const location = useLocation();
    const [isMenuOpen, setMenuOpen] = useState<boolean>(false);
    return (
        <>
            <AppBar position="static">
                <Toolbar>
                    <IconButton
                        size="large"
                        edge="start"
                        color="inherit"
                        aria-label="menu"
                        sx={{ mr: 2 }}
                    >
                        <MdMenu />
                    </IconButton>
                    <Typography>
                        {location.pathname in PageInformation
                            ? PageInformation[location.pathname].title
                            : "Undefined"}
                    </Typography>
                    <Menu
                        id="title_menu"
                        open={isMenuOpen}
                        onClick={() => {
                            setMenuOpen(!!!isMenuOpen);
                        }}
                    >
                        {Object.keys(PageInformation).map((k) => {
                            return (
                                <Link to={k}>
                                    <MenuItem>
                                        {PageInformation[k].title}
                                    </MenuItem>
                                </Link>
                            );
                        })}
                    </Menu>
                </Toolbar>
            </AppBar>
            <ConfigContext value={Config.body}>
                <Outlet />
            </ConfigContext>
        </>
    );
}

export const Route = createRootRoute({ component: RootLayout });
