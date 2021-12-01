import StoryRouter from 'storybook-react-router';
import {addDecorator} from "@storybook/react";
import "../src/index.css"
import React from "react";
import {ThemeProvider} from "../src/context/theme";
import Themes from "../src/themes/themes";
import {withThemes} from "storybook-addon-themes";

export const parameters = {
    layout: "fullscreen",
    themes: {
        clearable: false,
        default: "dark",
        list: [
            {
                name: "dark",
                color: "#222b36"
            },
            {
                name: "light",
                color: "#fff"
            }
        ],
        Decorator: ThemeDecorator
    },
    backgrounds: {disable: true},
    actions: {argTypesRegex: "^on[A-Z].*"},
    controls: {
        matchers: {
            color: /(background|color)$/i,
            date: /Date$/,
        },
    },
}

function ThemeDecorator(props) {
    const {children, themeName} = props;
    return (
        <ThemeProvider themes={Themes} defaultTheme={themeName}>
            {children}
        </ThemeProvider>
    );
}

addDecorator(StoryRouter());
addDecorator(withThemes);
