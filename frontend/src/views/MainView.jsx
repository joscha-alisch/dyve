import styles from "./MainView.module.sass";
import SideBar from "../components/sidebar/sidebar";
import React from "react";
import {
    faBook,
    faHouseUser,
    faLaptopCode,
    faRocket,
    faServer,
    faUserFriends,
    faUsers
} from "@fortawesome/free-solid-svg-icons";
import {Route, Switch, useParams} from "react-router-dom";
import AppDetail from "../components/appdetail/appdetail";
import AppFrame from "../components/base/frame/appframe/appframe";
import AppList from "../pages/apps/list/applist";

export const MainView = () => {
    return <AppFrame menuCategories={menuData}>
            <Switch>
                <Route exact path="/apps/">
                    <AppList />
                </Route>
                <Route path="/apps/:id">
                    <AppDetail/>
                </Route>
            </Switch>
    </AppFrame>
}

export const menuData = [
    {
        items: [
            {
                label: "Dashboard",
                icon: faHouseUser,
                soon: false,
                to: "/",
                exact: true,
            },
            {
                label: "Projects",
                icon: faBook,
                soon: true,
                to: "/projects",
            },
        ]
    },
    {
        title: "Platform",
        items: [
            {
                label: "Applications",
                icon: faLaptopCode,
                soon: false,
                to: "/apps",
            },
            {
                label: "Pipelines",
                icon: faRocket,
                soon: true,
                to: "/pipelines",
            },
            {
                label: "Services",
                icon: faServer,
                soon: true,
                to: "/services",
            },
        ]
    },
    {
        title: "Manage",
        items: [
            {
                label: "Teams",
                icon: faUsers,
                soon: false,
                to: "/teams",
            },
        ]
    },
]