import React, { FunctionComponent } from "react"

type ListProps = {
    className?: string,
}

const List : FunctionComponent<ListProps> = ({
    className = "",
    children
})  => <ul className={["m-0 p-0 w-full", className].join(" ")}>
    {children}
</ul>

export default List