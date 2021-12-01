import React from "react"
import styles from "./teamCard.module.sass"
import PropTypes from "prop-types"
import Box from "../../base/box/box";
import {StatusIcon} from "../../base/icons";
import {Link} from "react-router-dom";

const TeamCard = ({className, value}) => <Box className={styles.Main + " " + className}>
    <h3><Link to={"/teams/" + value.id }>{value.name}</Link></h3>
</Box>

TeamCard.propTypes = {
    className: PropTypes.string,
    value: PropTypes.element
}

export default TeamCard