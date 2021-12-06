import React from "react"
import PropTypes from "prop-types"

const Button = ({className}) => <div className={styles.Main + " " + className}>

</div>

Button.propTypes = {
    className: PropTypes.string,
}

export default Button