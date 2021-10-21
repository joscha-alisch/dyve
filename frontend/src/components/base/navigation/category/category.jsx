import React from "react"
import styles from "./category.module.sass"
import PropTypes from "prop-types"
import MenuItem from "../item/menuitem";

const Category = ({title, items, className}) => <ul className={styles.Main + " " + className}>
    <span className={styles.Title}>{title}</span>
    {items.map((item) => <MenuItem className={styles.Item} {...item}/>)}
</ul>

Category.propTypes = {
    title: PropTypes.string,
    items: PropTypes.arrayOf(PropTypes.element),
    className: PropTypes.string
}

export default Category