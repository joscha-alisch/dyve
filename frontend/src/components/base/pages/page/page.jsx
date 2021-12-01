import React from "react"
import styles from "./page.module.sass"
import PropTypes, {string} from "prop-types";
import {Link} from "react-router-dom";

const Page = ({title, parent, parentRoute, children, className, buttonsRender}) => <div
    className={styles.Main + " " + className}>
    <header className={styles.Header}>
        <section className={styles.TitleSection}>
            <span className={styles.Parent}>{parentRoute ? <Link to={parentRoute}>{parent}</Link> : parent}</span>
            <h1 className={styles.Title}>{title}</h1>
        </section>
        {buttonsRender && <nav className={styles.ButtonList}>
            {buttonsRender()}
        </nav>}
    </header>
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