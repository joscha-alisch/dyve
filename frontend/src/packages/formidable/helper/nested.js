import {isString} from "./is";

export const setNested = (obj, path, newValue) => {
    if (isString(path)) {
        path = path.split(".")
    }
    if (path.length === 0) {
        return newValue
    }

    let key = path[0]
    if (!obj[key]) {
        obj[key] = {}
    }
    path.shift()

    return {
        ...obj,
        [key]: setNested(obj[key], path, newValue)
    }
}

export const getNested = (obj, path) => {
    if (!obj) {
        return undefined
    }

    if (isString(path)) {
        path = path.split(".")
    }
    if (path.length === 0) {
        return obj
    }
    let key = path[0]
    obj = obj[key]
    if (!obj) {
        return undefined
    }

    path.shift()
    return getNested(obj, path)
}