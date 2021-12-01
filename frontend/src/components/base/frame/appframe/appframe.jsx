import styles from "./appframe.module.sass"
import TopBar from "../topbar/topbar";
import SideBar from "../sidebar/sidebar";
import PropTypes from "prop-types";

const AppFrame = ({menuCategories, children}) => <div className={styles.Main}>
    <TopBar className={styles.TopBar}/>
    <SideBar className={styles.SideBar} menuCategories={menuCategories}/>
    {children}
</div>

AppFrame.propTypes = {
    menuCategories: PropTypes.arrayOf(PropTypes.object),
    className: PropTypes.string
}

export default AppFrame