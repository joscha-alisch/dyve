import {useState} from "react";


export const useControlledState = (initialValues, value, onChange) => {
    let [state, setState] = useState(initialValues)
    return [value || state, (newState) => {
        setState(newState)
        if (onChange) {
            onChange(newState)
        }
    }]
}
