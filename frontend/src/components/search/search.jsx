import {Fragment} from "react";
import styles from "./search.module.sass"

const Search = () => <Fragment>
    <input type="text" placeholder="Search..." className={styles.Search} />
</Fragment>

export default Search