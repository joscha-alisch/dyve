import React from "react"
import Menu from "./menu";
import {menuData} from "../../../../views/MainView";

export default {
    title: 'Components/Menu/Menu',
    component: Menu,
}

export const StoryMenu = (args) => <Menu {...args}/>

StoryMenu.storyName = "Menu"
StoryMenu.args = {
    categories: menuData
}