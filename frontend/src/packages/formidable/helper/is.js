export const isObject = (a) => typeof a === 'object' && !Array.isArray(a) && a !== null
export const isString = (a) => typeof a === 'string' || a instanceof String