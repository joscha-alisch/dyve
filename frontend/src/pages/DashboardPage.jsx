import styles from "./DashboardPage.module.sass";
import SideBar from "../components/sidebar/sidebar";
import React from "react";
import {faLaptopCode, faUserFriends} from "@fortawesome/free-solid-svg-icons";
import Header from "../components/header/header";
import {Route, Switch} from "react-router-dom";
import AppList from "../components/applist/applist";
import AppDetail from "../components/appdetail/appdetail";
import Pipelinelist from "../components/pipelinelist/pipelinelist";
import PipelineDetail from "../components/pipelinedetail/pipelineDetail";
import Pipeline from "../components/pipeline/pipeline";

export const DashboardPage = () => <React.Fragment>
    <SideBar className={styles.SideBar} menuItems={[
        {
            title: "Platform",
            items: [
                {title: "Apps", icon: faLaptopCode, route: "/apps/"},
                /* { title: "Pipelines", icon: faRocket, route: "/pipelines/"},
                 { title: "Logging", icon: faStream, route: "/logging/"},
                 { title: "Metrics", icon: faChartLine, route: "/metrics/"},
                 { title: "Error Reporting", icon: faTemperatureHigh, route: "/errors/" },
                 { title: "Graph", icon: faProjectDiagram, route: "/graph/" },*/
            ]
        },
        /* {
             title: "Tools",
             items: [
                 { title: "Insights", icon: faSearchPlus, route: "/todo/" },
                 { title: "Network", icon: faNetworkWired, route: "/todo/" },
                 { title: "Costs", icon: faDollarSign, route: "/todo/"},
             ]
         },*/
        {
            title: "Manage",
            items: [
                {title: "Teams", icon: faUserFriends, route: "/teams"}
            ]
        },
    ]}/>
    <div className={styles.Flex}>
        <Header className={styles.Header}/>
        <main className={styles.Content}>
            <Switch>
                <Route exact path="/apps/">
                    <AppList page={0}/>
                </Route>
                <Route path="/apps/:id">
                    <AppDetail/>
                </Route>
                <Route exact path="/pipelines/">
                    <Pipelinelist page={0}/>
                </Route>
                <Route path="/pipelines/:id">
                    <PipelineDetail/>
                </Route>
                <Route path="/logging/">
                    <h1>Logging</h1>
                </Route>
                <Route path="/metrics/">
                    <h1>Metrics</h1>
                </Route>
                <Route path="/errors/">
                    <h1>Error Monitoring</h1>
                </Route>
                <Route path="/graph/">
                    <h1>Graph</h1>
                </Route>
                <Route path="/pipeline/">
                    <h1>Pipeline</h1>
                    <Pipeline/>
                </Route>
            </Switch>
        </main>
    </div>
</React.Fragment>