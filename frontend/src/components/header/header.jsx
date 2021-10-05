import styles from "./header.module.sass"
import {useUser} from "../../context/auth";

const Header = () => {
    let user = useUser()
    return <header className={styles.Header}>
        <div className={styles.UserBox}>
            <img className={styles.Avatar} alt="user avatar" src={user.picture}/>
            <span className={styles.Username}>{user.name}</span>
        </div>
    </header>;
}

export default Header