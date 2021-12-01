import React from "react"
import styles from "./sidebar.module.sass"
import UserWidget from "../../user/userwidget";
import Menu from "../../navigation/menu/menu";
import PropTypes from "prop-types"
import {useNotifications} from "../../../../context/notifications";

const SideBar = ({menuCategories, className}) =>  {
    const {notify} = useNotifications()

    return <div className={styles.Main + " " + className}>
        <UserWidget className={styles.User} variant={"default"} logoutUrl={"/user/logout"} profileUrl={"/user"}/>
        <Menu categories={menuCategories}/>
    </div>
}

SideBar.propTypes = {
    menuCategories: PropTypes.arrayOf(PropTypes.object),
    className: PropTypes.string
}

export default SideBar