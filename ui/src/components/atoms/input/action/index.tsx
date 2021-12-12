import React, { Component, MouseEventHandler } from "react"
import Icon from "../../display/icon"
import { Icons } from "../../display/icon/icons"

export type ActionProps = {
    className?: string,
    icon?: Icons,
    label: string,
    onClick: MouseEventHandler
}

const Action = ({
    className = "",
    icon = "",
    label,
    onClick
} : ActionProps)  => <button onClick={onClick} className={["flex flex-row items-center rounded px-2 py-1  hover:text-indigo-700 group text-gray-500", className].join(" ")}>
    <Icon icon={icon} className="w-4 h-4 mr-1" />
    <span>{label}</span>
</button>

export default Action