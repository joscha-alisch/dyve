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
} : ActionProps)  => <button onClick={onClick} className={["flex flex-row items-center rounded px-2 py-1 shadow border border-gray-100 hover:border-indigo-300 hover:text-indigo-700 group text-gray-800", className].join(" ")}>
    <Icon icon={icon} className="w-4 h-4 mr-1" />
    <span className="text-sm">{label}</span>
</button>

export default Action