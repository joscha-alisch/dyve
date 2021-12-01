import React from "react";
import styles from "./unknown.module.sass"

const errorField = ({common, options}) => ({value, errors, dirty, touched}, {onChange}, data) => {
    return <div className={styles.Main}>
        {options.message}<br/>
        <ul>
            {errors && errors.map(err => <li key={err}>{err}</li>)}
        </ul>
    </div>
}

export default errorField