import {
    createRootRoute,
    Link,
    linkOptions,
    Outlet,
    useLocation,
    useNavigate,
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
import Container from "@mui/material/Container";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";

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
        "/config": {
            title: "Config",
        },
    };
    const navigate = useNavigate();
    const [anchorElNav, setAnchorElNav] = useState<null | HTMLElement>(null);
    const handleOpenNavMenu = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorElNav(event.currentTarget);
    };
    const handleCloseNavMenu = () => {
        setAnchorElNav(null);
    };
    return (
        <div>
            <AppBar position="static">
                <Container maxWidth="xl">
                    <Toolbar disableGutters>
                        <Box
                            sx={{
                                flexGrow: 1,
                                display: { xs: "flex", md: "none" },
                            }}
                        >
                            <IconButton
                                size="large"
                                edge="start"
                                color="inherit"
                                aria-label="menu"
                                onClick={handleOpenNavMenu}
                                sx={{ mr: 2 }}
                            >
                                <MdMenu />
                            </IconButton>
                            <Menu
                                id="title_menu"
                                anchorEl={anchorElNav}
                                anchorOrigin={{
                                    vertical: "bottom",
                                    horizontal: "left",
                                }}
                                keepMounted
                                transformOrigin={{
                                    vertical: "top",
                                    horizontal: "left",
                                }}
                                open={Boolean(anchorElNav)}
                                onClose={handleCloseNavMenu}
                            >
                                {Object.keys(PageInformation).map((k) => {
                                    return (
                                        <MenuItem
                                            onClick={() =>
                                                navigate(linkOptions({ to: k }))
                                            }
                                        >
                                            {PageInformation[k].title}
                                        </MenuItem>
                                    );
                                })}
                            </Menu>
                        </Box>
                        <Typography
                            variant="h6"
                            noWrap
                            component="a"
                            sx={{
                                mr: 2,
                                flexGrow: { xs: 1, md: "unset" },
                                fontWeight: 700,
                                letterSpacing: ".3rem",
                                color: "inherit",
                                textDecoration: "none",
                            }}
                        >
                            DMXBOX
                        </Typography>
                        <Box
                            sx={{
                                flexGrow: 1,
                                display: { xs: "none", md: "flex" },
                            }}
                        >
                            {Object.keys(PageInformation).map((page) => (
                                <Button
                                    key={page}
                                    onClick={() => {
                                        navigate(linkOptions({ to: page }));
                                    }}
                                    sx={{
                                        my: 2,
                                        color: "white",
                                        display: "block",
                                    }}
                                >
                                    {PageInformation[page].title}
                                </Button>
                            ))}
                        </Box>
                    </Toolbar>
                </Container>
            </AppBar>
            <ConfigContext value={Config.body}>
                <Outlet />
            </ConfigContext>
        </div>
    );
}

export const Route = createRootRoute({ component: RootLayout });
