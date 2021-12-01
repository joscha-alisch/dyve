import React from "react";
import {faBook, faHouseUser, faLaptopCode, faRocket, faServer, faUsers} from "@fortawesome/free-solid-svg-icons";
import {Route, Switch} from "react-router-dom";
import AppFrame from "../components/base/frame/appframe/appframe";
import Pages from "../pages";

export const MainView = () => <AppFrame menuCategories={menuData}>
    <Switch>
        <Route exact path="/user/logout">
            <Pages.User.Logout/>
        </Route>

        <Route exact path="/apps/">
            <Pages.Apps.List/>
        </Route>
        <Route path="/apps/:id">
            <Pages.Apps.Detail/>
        </Route>

        <Route exact path="/pipelines/">
            <Pages.Pipelines.List/>
        </Route>

        <Route exact path="/teams/">
            <Pages.Teams.List/>
        </Route>
        <Route exact path="/teams/new">
            <Pages.Teams.New/>
        </Route>
        <Route path="/teams/:id">
            <Pages.Teams.Detail/>
        </Route>
    </Switch>
</AppFrame>

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