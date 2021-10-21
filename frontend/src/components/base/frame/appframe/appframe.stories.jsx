import React from "react"
import AppFrame from "./appframe";
import {UserContext} from "../../../../context/auth";
import { userData} from "../../../../../.storybook/data";
import {menuData} from "../../../../views/MainView";
import Page from "../../pages/page/page";

export default {
    title: 'App/Frame/Composed',
    component: AppFrame,
    isFullscreen: true
}

export const StoryAppFrame = (args) => <UserContext.Provider value={{name: args.userName, picture: args.avatarUrl}}>
    <div style={{height: "100vh"}}>
        <AppFrame {...args}>
            <Page title={args.pageTitle} parent={args.pageParent}/>
        </AppFrame>
    </div>
</UserContext.Provider>;

StoryAppFrame.storyName = "Composed"
StoryAppFrame.args = {
    ...userData,
    menuCategories: menuData,
    pageTitle: "Current Page",
    pageParent: "Parent Page"
}