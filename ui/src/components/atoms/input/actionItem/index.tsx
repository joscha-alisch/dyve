import React, { MouseEventHandler } from "react"
import { Icon } from "../../"
import { Icons } from "../../display/icon/icons"

type ContextMenuItemProps = {
    className?: string,
    label: string,
    icon?: Icons,
    onClick?: MouseEventHandler
}

const ActionItem = ({
    className = "",
    label,
    icon,
    onClick = () => {}
} : ContextMenuItemProps)  => <li onClick={onClick} className={["w-full list-none flex flex-row items-center gap-2 text-gray-600 hover:text-indigo-600 hover:bg-gray-100 cursor-pointer m-0 p-2", className].join(" ")}>
   <Icon icon={icon || ""} className="w-4 h-4" />{label}
</li>

export default ActionItem