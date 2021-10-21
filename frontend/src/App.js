import {MainView} from "./views/MainView";
import styles from "./App.module.sass";

function App() {
  return <div className={styles.App + " nodebug"}>
    <MainView />
  </div>
}

export default App;
