import React from "react"
import styles from "./topbar.module.sass"
import Logo from "../../logo/logo";
import UserWidget from "../../user/userwidget";
import PropTypes from "prop-types";
import SideBar from "../sidebar/sidebar";

const TopBar = ({className}) => <div className={styles.Main + " " + className}>
    <Logo className={styles.Logo} />
    <UserWidget variant={"small"} />
</div>

TopBar.propTypes = {
    className: PropTypes.string
}

export default TopBar