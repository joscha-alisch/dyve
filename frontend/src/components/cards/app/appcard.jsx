import React from "react"
import styles from "./appcard.module.sass"
import PropTypes from "prop-types"
import Box from "../../base/box/box";
import {StatusIcon} from "../../base/icons";
import {Link} from "react-router-dom";

const AppCard = ({className, value}) => <Box className={styles.Main + " " + className}>
    <StatusIcon status="green" scale={20}/>
    <h3><Link to={"/apps/"+value.id}>{value.name}</Link></h3>
</Box>

AppCard.propTypes = {
    className: PropTypes.string,
    value: PropTypes.element
}

export default AppCard