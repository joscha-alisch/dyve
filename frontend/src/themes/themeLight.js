import {createTheme} from "@mui/material";
import {muiDefaults} from "./mui";

export const themeLight = {
    className: "themeLight",
    mui: createTheme(muiDefaults, {})
}

themeLight.mui = createTheme(themeLight.mui, {})