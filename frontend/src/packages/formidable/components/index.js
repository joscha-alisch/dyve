import stringField from "./string"
import errorField from "./unknown"
import {merge} from "../helper/merge";


export const defaults = {
    string: {
        default: stringField
    },
    unknown: errorField,
}

export const withComponents = (fields) => {
    return merge(defaults, fields)
}

export default defaults