import React from "react"
import styles from "./box.module.sass"
import PropTypes from "prop-types";

const Box = ({className, children, title}) => <div className={styles.Main + " " + className}>
    {title && title !== "" ? <h2>{title}</h2> : ""}
    {children}
</div>

Box.propTypes = {
    className: PropTypes.string,
    title: ""
}

export default Box