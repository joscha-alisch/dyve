import React, { FunctionComponent } from "react"

type TextFieldProps = {
    className?: string,
    multiLine?: boolean
    lines?: number,
    placeholder?: string,
    label?: string,
    name?: string,
    value?: string
    autoFocus?: boolean
    onChange?: (value : string) => void
}

const TextField : FunctionComponent<TextFieldProps> = ({
    className,
    multiLine = false,
    autoFocus = false,
    value = "",
    onChange,
})  => {
    if (multiLine) {
        return <textarea />
    } else {
        return <input autoFocus={autoFocus} value={value} onChange={(ev) => onChange && onChange(ev.target.value)} type="text" className={["min-w-0 box-border rounded border-gray-400 hover:border-indigo-600 focus:border-indigo-600 focus:ring-indigo-600", className].join(" ")} />
    }
}
    
export default TextField