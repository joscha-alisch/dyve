import React from "react"
import MenuItem from "./menuitem";
import {faLaptopCode} from "@fortawesome/free-solid-svg-icons";

export default {
    title: 'Components/Menu/Item',
    component: MenuItem,
}

export const StoryMenuItem= (args) => <MenuItem {...args}/>

StoryMenuItem.storyName = "Item"
StoryMenuItem.args = {
    soon: false,
    label: "Menu Item",
    to: "/page",
    icon: faLaptopCode,
}