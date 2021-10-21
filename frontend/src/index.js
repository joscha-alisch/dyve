import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import {AppProviders} from "./context/providers";
import axios from "axios";

function getCookies() {
    return document.cookie.split("; ").reduce((c, x) => {
        const splitted = x.split("=");
        c[splitted[0]] = splitted[1];
        return c;
    }, {});
}

axios.interceptors.request.use(function (config) {
    const token = getCookies()["XSRF-TOKEN"];
    config.headers["X-XSRF-TOKEN"] = token;

    return config;
});

ReactDOM.render(
    <React.StrictMode>
        <AppProviders>
            <App/>
        </AppProviders>
    </React.StrictMode>,
    document.getElementById('root')
);

