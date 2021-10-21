import React from "react"
import Logo from "./logo";

export default {
    title: 'Components/Logo',
    component: Logo,
}

export const StoryLogo = (args) => <Logo {...args}/>;

StoryLogo.storyName = "Logo"
StoryLogo.args = {}