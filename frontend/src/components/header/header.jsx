import styles from "./header.module.sass"
import Search from "../search/search";

const Header = () => <header className={styles.Header}>
    <Search />
</header>;

export default Header