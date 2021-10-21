import {FontAwesomeIcon} from "@fortawesome/react-fontawesome"
import {faCircle} from "@fortawesome/free-solid-svg-icons";
import React from "react"
import PropTypes from "prop-types";

const StatusIcon = ({status, className, style, scale, rotate}) => {
    let transform = {}

    if (scale) {
        transform.size = scale
    }

    if (rotate) {
        transform.rotate = rotate
    }

    return <FontAwesomeIcon className={className} color={status} style={style} transform={transform} icon={faCircle}/>
}

StatusIcon.propTypes = {
    status: PropTypes.oneOf(["green", "red"]),
    scale: PropTypes.number,
    rotate: PropTypes.number,
    className: PropTypes.string,
    style: PropTypes.element
}

export default StatusIcon