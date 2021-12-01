import {createTheme} from "@mui/material";
import {muiDefaults} from "./mui";

let mui = createTheme(muiDefaults, {
    palette: {
        primary: {
            main: "#008C9E",
            light: "#00B4CC",
            dark: "#005F6B",
            contrastText: "#fff"
        },
        text: {
            primary: "#fff",
            secondary: "#fff",
            disabled: "grey"
        },
    }
})

mui = createTheme(mui, {
        palette: {
            action: {
                active: "#fff",
                hover: mui.palette.primary.dark
            }
        },
        components: {
            MuiPaper: {
                styleOverrides: {
                    root: {
                        backgroundColor: "#222b36",
                    }
                }
            }
        }
    }
)

export const themeDark = {
    className: "themeDark",
    mui: mui
}