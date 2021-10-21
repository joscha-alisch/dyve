import React from "react"
import {NavLink} from "react-router-dom";
import styles from "./menuitem.module.sass"
import PropTypes from "prop-types"
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome'

const MenuItem = ({to, label, soon, icon, exact, className}) => <li
    className={styles.Main + " " + className + (soon ? " " + styles.Soon : "")}>
    <NavLink exact={exact} to={to} activeClassName={styles.Active}>
        <span className={styles.Icon}><FontAwesomeIcon icon={icon}/></span>
        <span className={styles.Label}>{label}</span>
        {soon ? <span className={styles.Tag}>soon</span> : ""}
    </NavLink>
</li>

MenuItem.propTypes = {
    to: PropTypes.string,
    label: PropTypes.string,
    soon: PropTypes.bool,
    icon: PropTypes.element,
    className: PropTypes.string,
    exact: PropTypes.bool
}

export default MenuItem