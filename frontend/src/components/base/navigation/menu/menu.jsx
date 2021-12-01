import React from "react"
import styles from "./menu.module.sass"
import Category from "../category/category";
import PropTypes from "prop-types"

const Menu = ({categories, className}) => <div className={styles.Main + " " + className}>
    {categories.map((category) => <Category key={category.title ? category.title : ""}
                                            className={styles.Category} {...category} />)}
</div>

Category.propTypes = {
    categories: PropTypes.arrayOf(PropTypes.element),
    className: PropTypes.string
}

export default Menu