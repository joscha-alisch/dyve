import React from "react"
import styles from "./logo.module.sass"

const Logo = ({className}) => <img className={styles.Logo + " " + className} alt="dyve logo" src="/img/logo.png"/>

export default Logo