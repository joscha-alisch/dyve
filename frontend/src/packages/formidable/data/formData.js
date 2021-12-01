import {useEffect, useState} from "react";
import {validate} from "../validation/validation";
import {isObject} from "../helper/is";
import {defaulted} from "../helper/defaulted";
import {getNested, setNested} from "../helper/nested";
import {getStringOptions} from "../components/string/string";

const traverseFields = (obj, each, path) => {
    if (!path) {
        path = ""
    }

    for (let key in obj) {
        if (obj.hasOwnProperty(key)) {
            if (isObject(obj[key])) {
                let newPath = key
                if (path !== "") {
                    newPath = path + "." + newPath
                }

                if (!obj[key].type || obj[key].type === "object") {
                    traverseFields(obj[key], each, newPath)
                } else {
                    each(key, newPath, obj[key])
                }
            }
        }
    }
}

const getCommon = (field, info) => {
    return {
        id: field,
        type: info.type,
        label: (info.title && info.title !== "") ? info.title : field,
        hint: info.hint,
    }
}

const getOptionsFor = (type, info) => {
    switch (type) {
        case "string":
            return getStringOptions(info)
        default:
            return {}
    }
}

const buildFormData = (data, initialValues) => {
    let state = {
        values: {},
        fields: {},
        runtime: {},
    }

    traverseFields(data, (field, path, info) => {
            let initial = getNested(initialValues, path)
            let defaultValue = defaulted(info.default, info.type)

            let value = initial ? initial : defaultValue

            state.fields[path] = {
                common: getCommon(field, info),
                options: getOptionsFor(info.type, info)
            }
            state.runtime[path] = {
                value: value,
                errors: [],
                dirty: false,
                touched: false,
            }
            state.values = setNested(state.values, path, value)
        }
    )

    return state
}


export const useFormData = (data, ctrlState, ctrlSetState) => {
    let [state, setState] = useState({
        value: ctrlState
    })

    let change = (path, value) => {
        let newState = setNested(state, "values." + path, value)
        newState["runtime"][path].value = value
        newState["runtime"][path].touched = true
        newState["runtime"][path].dirty = true
        setState(newState)
    }

    useEffect(() => {
        setState(buildFormData(data, ctrlState))
    }, [setState])

    useEffect(() => {
        validate(state)
        ctrlSetState && ctrlSetState(state.values)
    }, [state])

    return {
        values: state.values,
        fields: state.fields,
        runtime: state.runtime,
        handlers: {
            onChange: (path) => (value) => change(path, value),
        },
    }
}

