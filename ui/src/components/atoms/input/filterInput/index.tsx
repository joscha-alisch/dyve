import React, { ChangeEventHandler, FunctionComponent } from "react"
import {Icon} from "../../../atoms"

type FilterInputProps = {
    className?: string,
    autoFocus?: boolean,
    value: string,
    onChange: (value: string) => void
}

const FilterInput: FunctionComponent<FilterInputProps> = ({
    className = "",
    autoFocus = false,
    value,
    onChange
} : FilterInputProps)  => <div className={["relative w-full", className].join(" ")}>
        <div className="absolute inset-y-0 flex items-center">
            <Icon icon="search" className="w-3 h-3 text-gray-400 ml-2 block" />
        </div>
        <input type="search" value={value} onChange={(ev) => onChange(ev.target.value)} autoFocus={autoFocus} placeholder="Filter" className="w-full min-w-0 rounded-full border-gray-200 text-xs pl-6 py-1" />
    </div>

export default FilterInput