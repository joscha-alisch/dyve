import React from "react"
import { Icon } from "../.."

type PageCounterProps = {
    className?: string,
    page: number,
    perPage: number,
    totalItems: number,
    onPageChange: (page: number) => void
}

const PageCounter = ({
    className = "",
    page,
    perPage,
    totalItems,
    onPageChange,
}: PageCounterProps) => {

    let start = (page * perPage) + 1
    let end = Math.min((page * perPage) + perPage, totalItems)

    const onLeft = () => {
        if (page > 0) {
            onPageChange(page-1)
        }
    }
    const onRight = () => {
        if (page < (totalItems / perPage) - 1) {
            onPageChange(page+1)
        }
    }

    return <div className={["flex flex-row items-center text-sm text-gray-500 select-none", className].join(" ")}>
        <Icon icon="point-left" onClick={onLeft} className="w-7 h-7 hover:text-gray-900 cursor-pointer" />
        <span className="text-gray-700">{start}</span>
        <span>-</span>
        <span className="text-gray-700">{end}</span>
        <span className="mx-1 ">of</span>
        <span className="text-gray-700">{totalItems}</span>
        <Icon icon="point-right" onClick={onRight} className="w-7 h-7 hover:text-gray-900 cursor-pointer" />
    </div>
}

export default PageCounter