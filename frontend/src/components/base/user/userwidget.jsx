import React, {useState} from "react"
import styles from "./userwidget.module.sass"
import PropTypes from "prop-types"
import {Link} from "react-router-dom";
import {faCog, faSignOutAlt} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome'
import {useUser} from "../../../context/auth";

const UserWidget = (props) => {
    let user = useUser()

    if (props.variant === "small") {
        return userWidgetSmall(user, props)
    } else {
        return userWidgetDefault(user, props)
    }
}

const userWidgetDefault = ({picture, name}, {className, profileUrl, logoutUrl}) => <Link to={"/user"}><div className={styles.Box + " " + styles.Default + " " + className}>
    <div className={styles.Avatar}><img alt="user avatar" src={picture}/></div>
    <div className={styles.UserName}>{name}</div>
    <div className={styles.Buttons}>
        <Link to={profileUrl}><FontAwesomeIcon icon={faCog}/></Link>
        <Link to={logoutUrl}><FontAwesomeIcon icon={faSignOutAlt}/></Link>
    </div>
</div></Link>

const userWidgetSmall = ({name, picture}, {className, logoutUrl, profileUrl, smallExpanded}) => <div className={styles.Small + " " + className} >
    <div className={styles.Avatar}><img alt="user avatar" src={picture}/></div>
    <div className={styles.HoverContent} style={smallExpanded ? {display: "block"} : {}}>
        <div className={styles.UserName}>{name}</div>
        <div className={styles.Buttons}>
            <Link to={profileUrl}><FontAwesomeIcon icon={faCog}/></Link>
            <Link to={logoutUrl}><FontAwesomeIcon icon={faSignOutAlt} /></Link>
        </div>
    </div>
</div>


UserWidget.propTypes = {
    variant: PropTypes.oneOf(["default", "small"]),
    logoutUrl: PropTypes.string,
    profileUrl: PropTypes.string,
    smallExpanded: PropTypes.bool,
    className: PropTypes.string
}

export default UserWidget