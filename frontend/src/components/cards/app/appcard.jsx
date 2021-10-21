import React from "react"
import styles from "./appcard.module.sass"
import PropTypes from "prop-types"
import Box from "../../base/box/box";
import {StatusIcon} from "../../base/icons";

const AppCard = ({className, value}) => <Box className={styles.Main + " " + className}>
    <StatusIcon status="green" scale={20} />
    <h3>{value.name}</h3>
</Box>

AppCard.propTypes = {
    className: PropTypes.string,
    value: PropTypes.element
}

export default AppCard