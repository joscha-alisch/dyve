import React from "react"
import styles from "./page.module.sass"
import PropTypes, {string} from "prop-types";

const Page = ({title, parent, children, className}) => <div className={styles.Main + " " + className}>
    <span className={styles.Parent}>{parent}</span>
    <h1 className={styles.Title}>{title}</h1>
    <div className={styles.Children}>
        {children}
    </div>
</div>

Page.propTypes = {
    parent: string,
    title: string,
    className: PropTypes.string,
    children: PropTypes.node,
}


export default Page