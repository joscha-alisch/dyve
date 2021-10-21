import * as React from 'react'
import {useEffect} from 'react'
import useLocalStorage from "../hooks/useLocalStorage";
import {CssBaseline, ThemeProvider as MUIThemeProvider} from "@mui/material"

const ThemeContext = React.createContext([{}, () => {
}])

export const ThemeProvider = ({children, themes, defaultTheme}) => {
    const [themeName, setTheme] = useLocalStorage("theme", defaultTheme)

    useEffect(() => {
        setTheme(defaultTheme)
    }, [defaultTheme])

    let theme = themes[themeName]
    let mui = theme.mui

    if (!theme) {
        setTheme(defaultTheme)
    }

    return <ThemeContext.Provider value={[{name: themeName, theme: theme}, setTheme]}>
        <MUIThemeProvider theme={mui}>
            <CssBaseline/>
            <div id="themeProvider" className={theme.className}>
                {children}
            </div>
        </MUIThemeProvider>
    </ThemeContext.Provider>
}

export const useTheme = () => React.useContext(ThemeContext)
