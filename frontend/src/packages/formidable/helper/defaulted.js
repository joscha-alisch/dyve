export const defaulted = (value, type) => {
    if (value) {
        return value
    }

    switch (type) {
        case "string":
            return ""
        default:
            return undefined
    }
}