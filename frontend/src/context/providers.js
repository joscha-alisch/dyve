import {AuthProvider} from "./auth";
import {QueryParamProvider} from "use-query-params";
import {BrowserRouter as Router, Route} from "react-router-dom";
import React from "react";
import {ThemeProvider} from "./theme";
import Themes, {defaultTheme} from "../themes/themes";
import {NotificationProvider} from "./notifications";

export const AppProviders = ({children}) => {
    return <Router>
        <NotificationProvider>
            <AuthProvider>
                <QueryParamProvider ReactRouterRoute={Route}>
                    <ThemeProvider themes={Themes} defaultTheme={defaultTheme}>
                        {children}
                    </ThemeProvider>
                </QueryParamProvider>
            </AuthProvider>
        </NotificationProvider>
    </Router>
}