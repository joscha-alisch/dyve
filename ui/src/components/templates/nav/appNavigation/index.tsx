import React from "react"
import { IconItem, List } from "../../../atoms"

type AppNavigationProps = {
    className?: string,
}

const AppNavigation = ({
    className = "",
} : AppNavigationProps)  => <aside className={["", className].join(" ")}>
    <nav className="sticky top-0 left-0 h-screen group flex-none flex flex-row bg-gray-50 w-12 hover:w-72 transition-all overflow-hidden">
        <List>
            <IconItem icon="server"></IconItem>
        </List>
    </nav>
</aside>

export default AppNavigation