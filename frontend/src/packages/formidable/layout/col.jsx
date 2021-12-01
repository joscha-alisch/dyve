export const renderLayout = (layout, children) => {
    switch (layout) {
        case "row":
            return <div className={styles.FormLayoutRow}>{children}</div>
        case "column":
            return <div className={styles.FormLayoutColumn}>{children}</div>
    }
}