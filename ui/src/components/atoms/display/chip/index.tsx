import React, { FunctionComponent, MouseEventHandler } from "react"

type ChipProps = {
    className?: string,
    specificer?: string,
    label: string,
    value?: string | number,
    onClick?: MouseEventHandler
}

const Chip : FunctionComponent<ChipProps> = ({
    className = "",
    specificer = "",
    label,
    value = "",
    onClick,
    children
}) => <div onClick={onClick} className={["rounded select-none inline-block px-2 py-1 shadow w-auto cursor-pointer bg-indigo-600 hover:bg-indigo-700 group", className].join(" ")}>
    { specificer !== "" && <>
        <span className="text-gray-300 mr-1">{specificer}</span>
    </>}
    <span className="text-gray-100 group-hover:text-gray-50">{label}</span>
    {value !== "" && <>
        <span className="text-gray-300 group-hover:text-gray-200">:</span>
        <span className="text-gray-300 group-hover:text-gray-200 ml-1">{value}</span>
    </>}
    {children}
</div>

export default Chip