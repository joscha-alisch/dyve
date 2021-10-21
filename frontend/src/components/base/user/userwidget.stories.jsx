import React from "react"
import UserWidget from "./userwidget";
import {UserContext} from "../../../context/auth";
import {userData} from "../../../../.storybook/data";

export default {
    title: 'Components/User Widget',
    component: UserWidget,
}

export const StoryUserWidget= (args) => <UserContext.Provider value={{name: args.userName, picture: args.avatarUrl}}><div style={{display: "flex", flexDirection: "column", height: 500, justifyContent: "space-evenly", alignItems: "stretch"}}>
    <div style={{flex: "0 0 100px"}}>
        <UserWidget {...args} variant={"default"}/>;
    </div>
    <div style={{flex: "0 0 100px", display: "flex", flexDirection: "row", height: 500, justifyContent: "center", alignItems: "center"}}>
        <UserWidget  {...args} variant={"small"}/>;
    </div>
</div></UserContext.Provider>

StoryUserWidget.storyName = "User Widget"
StoryUserWidget.args = {
    profileUrl: "/user",
    logoutUrl: "/user/logout",
    smallExpanded: false,
    ...userData
}

//asu28asd