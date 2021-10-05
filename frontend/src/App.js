import {DashboardPage} from "./pages/DashboardPage";
import styles from "./App.module.sass";

function App() {
  return <div className={styles.App + " nodebug"}>
    <DashboardPage />
  </div>
}

export default App;
