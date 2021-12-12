import React, { FunctionComponent } from "react"
import {List} from "../../../atoms"

type AppNavigationProps = {
    className?: string,
}

const AppNavigation : FunctionComponent<AppNavigationProps> = ({
    className = "",
})  => <aside className={["", className].join(" ")}>
    <nav className="sticky top-0 left-0 h-screen group flex-none flex flex-row bg-gray-50 w-12 hover:w-72 transition-all overflow-hidden">
        <List>
        </List>
    </nav>
</aside>

export default AppNavigation