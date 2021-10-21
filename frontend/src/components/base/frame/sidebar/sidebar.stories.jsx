import React from "react"
import SideBar from "./sidebar";
import {userData} from "../../../../../.storybook/data";
import {UserContext} from "../../../../context/auth";
import {menuData} from "../../../../views/MainView";

export default {
    title: 'App/Frame/Side Bar',
    component: SideBar
}

export const StorySideBar = (args) => <UserContext.Provider value={{name: args.userName, picture: args.avatarUrl}}>
    <SideBar {...args}/>
</UserContext.Provider>

StorySideBar.storyName = "Side Bar"
StorySideBar.args = {
    ...userData,
    menuCategories: menuData,
}