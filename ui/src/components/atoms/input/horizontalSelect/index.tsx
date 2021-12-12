import React, { FunctionComponent } from "react"

export type HorizontalSelectOption = {
    label: string,
    value: string | number,
}

type HorizontalSelectProps = {
    className?: string,
    label: string,
    options: HorizontalSelectOption[],
    selected: string | number
    onSelect: (option: string | number) => void
}

const HorizontalSelect: FunctionComponent<HorizontalSelectProps> = ({
    className = "",
    label,
    options,
    selected,
    onSelect,
}) => <div className={["text-xs flex flex-row text-gray-700", className].join(" ")}>
        <label className="mr-2 text-gray-500">
            {label}:
        </label>
        <ul className="flex flex-row space-x-2">
            {options && options.map(option => {
                let classes = "cursor-pointer hover:text-indigo-600"
                if (selected === option.value) {
                    classes += " font-bold text-indigo-600"
                }
                return <li onClick={() => onSelect(option.value)} className={classes}>
                    {option.label}
                </li>
            })}
        </ul>
    </div>

export default HorizontalSelect