import React from "react"
import PropTypes from "prop-types"
import {faNetworkWired} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import Box from "../../base/box/box";
import styles from "./routingcard.module.sass"
import {CircularProgress} from "@mui/material";

const RoutingCard = ({className, routing}) => <Box title="Routes">
    <ul className={styles.Routes}>
        { !routing && <CircularProgress />}
        { routing && routing.routes && routing.routes.map((route, index) => <li key={index} className={styles.Route}>
            <div className={styles.Icon}><FontAwesomeIcon icon={faNetworkWired} /></div>
            <span className={styles.Url}>
                <span className={styles.Host}>{route.host}</span>
                <span className={styles.Path}>{route.path}</span>
            </span>
            { route.appPort !== 0 &&
                <span className={styles.Port}><b>Port:</b> {route.appPort}</span>}
        </li>)}
    </ul>
</Box>


RoutingCard.propTypes = {
    className: PropTypes.string,
}

export default RoutingCard