import React from "react"
import TopBar from "./topbar";
import AppFrame from "../appframe/appframe";
import {userData} from "../../../../../.storybook/data";
import {UserContext} from "../../../../context/auth";

export default {
    title: 'App/Frame/Top Bar',
    component: TopBar,
}

export const StoryTopBar = (args) => <UserContext.Provider value={{name: args.userName, picture: args.avatarUrl}}>
    <TopBar {...args}/>
</UserContext.Provider>

StoryTopBar.storyName = "Top Bar"
StoryTopBar.args = {
    ...userData
}
