import React from "react"
import PropTypes from "prop-types"
import {faNetworkWired} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import Box from "../../base/box/box";
import styles from "./instancescard.module.sass"
import {CircularProgress} from "@mui/material";

const InstancesCard = ({className, instances}) => <Box title="Instances" className={styles.Main}>
    <ul className={styles.Instances}>
        { !instances && <CircularProgress />}
        { instances && instances.map((instance, index) => <li key={index} className={styles.Instance}>
            <div className={ instance.state === "running" ? styles.IconGreen : styles.IconRed}>{index}</div>
            <span className={styles.Info}>
                <span className={styles.State}>{instance.state}</span> since <span className={styles.Date}>{(new Date(instance.since)).toUTCString()}</span>
            </span>
        </li>)}
    </ul>
</Box>


InstancesCard.propTypes = {
    className: PropTypes.string,
}

export default InstancesCard