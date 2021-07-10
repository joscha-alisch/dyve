import Header from "./components/header/header";
import SideBar from "./components/sidebar/sidebar";
import styles from "./App.module.sass"
import {
    faChartLine,
    faCoffee, faDollarSign, faLaptopCode, faNetworkWired,
    faProjectDiagram,
    faRocket, faSearchPlus,
    faStream,
    faTemperatureHigh, faUserFriends, faWindowRestore
} from "@fortawesome/free-solid-svg-icons";
import {
    Switch,
    Route,
} from "react-router-dom";
import AppList from "./components/applist/applist";
import AppDetail from "./components/appdetail/appdetail";

function App() {
  return (
    <div className={styles.App + " nodebug"}>
        <Header className={styles.Header} />

        <div className={styles.Flex}>
            <SideBar className={styles.SideBar} menuItems={[
                {
                    title: "Platform",
                    items: [
                        { title: "Apps", icon: faLaptopCode, route: "/apps/"},
                        { title: "Pipelines", icon: faRocket, route: "/pipelines/"},
                        { title: "Logging", icon: faStream, route: "/logging/"},
                        { title: "Metrics", icon: faChartLine, route: "/metrics/"},
                        { title: "Error Reporting", icon: faTemperatureHigh, route: "/errors/" },
                        { title: "Graph", icon: faProjectDiagram, route: "/graph/" },
                    ]
                },
                {
                    title: "Tools",
                    items: [
                        { title: "Insights", icon: faSearchPlus },
                        { title: "Network", icon: faNetworkWired },
                        { title: "Teams", icon: faUserFriends },
                        { title: "Costs", icon: faDollarSign},
                    ]
                },
            ]}/>
            <main className={styles.Content}>
                <Switch>
                    <Route exact path="/apps/">
                        <AppList page={0}/>
                    </Route>
                    <Route path="/apps/:id">
                        <AppDetail/>
                    </Route>
                    <Route path="/pipelines/">
                        <h1>Pipelines</h1>
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
                </Switch>
            </main>
        </div>
    </div>
  );
}

export default App;
