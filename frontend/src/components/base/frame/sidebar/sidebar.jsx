import React from "react"
import styles from "./sidebar.module.sass"
import UserWidget from "../../user/userwidget";
import Menu from "../../navigation/menu/menu";
import PropTypes from "prop-types"

const SideBar = ({menuCategories, className}) => <div className={styles.Main + " " + className}>
    <UserWidget className={styles.User} variant={"default"} />
    <Menu categories={menuCategories} />
</div>

SideBar.propTypes = {
    menuCategories: PropTypes.arrayOf(PropTypes.element),
    className: PropTypes.string
}

export default SideBar