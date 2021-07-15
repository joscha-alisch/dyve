import styles from "./sidebar.module.sass"
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome'
import {NavLink} from "react-router-dom";
import {
    faArrowLeft,
    faArrowRight,
    faChevronCircleLeft, faChevronLeft,
    faChevronRight,
    faSearch
} from "@fortawesome/free-solid-svg-icons";
import useLocalStorage from "../../hooks/useLocalStorage";

const SideBar = ({className, menuItems}) => {
    let [collapsed, setCollapsed] = useLocalStorage("sidebarCollapsed", false)

    return <aside className={styles.SideBar + " " + className + (collapsed ? " " + styles.Collapsed : "")}>
        <a className={styles.Toggle} onClick={() => setCollapsed(!collapsed)}>{collapsed ?
            <FontAwesomeIcon icon={faChevronRight}/> : <FontAwesomeIcon icon={faChevronLeft}/>}</a>
        <nav>
            <ul>
                <li>
                    <ul>
                        <li>
                            <NavLink activeClassName={styles.active} to={"/search"}>
                             <span className={styles.MenuIcon}>
                                    <FontAwesomeIcon icon={faSearch}/>
                                </span>
                                <span className={styles.MenuText}>Search</span>
                            </NavLink>
                        </li>
                    </ul>
                </li>
                {menuItems.map((item) => <li>
                    <header>{item.title}</header>
                    <ul>
                        {item.items.map((subItem) => <li><NavLink activeClassName={styles.active}
                                                                  to={subItem.route || "#"}>
                            <span className={styles.MenuIcon}>
                                <FontAwesomeIcon icon={subItem.icon}/>
                            </span>
                            <span className={styles.MenuText}>
                                 {subItem.title}
                            </span>

                        </NavLink>
                        </li>)}
                    </ul>
                </li>)}
            </ul>

        </nav>
        <footer>
            dyve<br/>v0.0.1
        </footer>
    </aside>
}

export default SideBar