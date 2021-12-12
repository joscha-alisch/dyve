import React, { FunctionComponent } from "react"

type ListGroupProps = {
    className?: string,
    label: string
}

const ListGroup : FunctionComponent<ListGroupProps> = ({
    className = "",
    label,
    children
})  => <li className={["list-none p-0 m-0", className].join(" ")}>
    <h3 className="text-gray-400 uppercase text-xs pl-2 mt-3 mb-1">{label}</h3>
    <ul className="p-0">
        {children}
    </ul>
</li>

export default ListGroup