import {AuthProvider} from "./auth";
import {QueryParamProvider} from "use-query-params";
import {BrowserRouter as Router, Route} from "react-router-dom";
import React from "react";

export const AppProviders = ({children}) => {
    return <Router>
        <AuthProvider>
            <QueryParamProvider ReactRouterRoute={Route}>
                {children}
            </QueryParamProvider>
        </AuthProvider>
    </Router>
}