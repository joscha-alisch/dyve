import {Fragment} from "react";
import styles from "./sidebar.module.sass"
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import {
    NavLink
} from "react-router-dom";

const SideBar = ({className, menuItems}) => <aside className={styles.SideBar + " " + className}>
    <nav>
        <ul>
            {menuItems.map((item) => <li>
                    <header>{item.title}</header>
                    <ul>
                        {item.items.map((subItem) => <li><NavLink activeClassName={styles.active} to={subItem.route || "#"}>
                            <span className={styles.MenuIcon}>
                                <FontAwesomeIcon icon={subItem.icon} />
                            </span>
                            {subItem.title}
                        </NavLink>
                        </li>)}
                    </ul>
                </li>) }
        </ul>

    </nav>
    <footer>
        dyve<br/>v0.0.1
    </footer>
</aside>

export default SideBar