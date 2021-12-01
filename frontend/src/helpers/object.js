
export const forEach = (obj, each) => {
    for (let key in obj) {
        if (!obj.hasOwnProperty(key)) {
            continue
        }
        each(obj[key])
    }
}

export const flatMap = (obj, flatMap) => {
    let res = []
    for (let key in obj) {
        if (!obj.hasOwnProperty(key)) {
            continue
        }

        res = [
            ...res,
            ...flatMap(obj[key]),
        ]
    }
    return res
}