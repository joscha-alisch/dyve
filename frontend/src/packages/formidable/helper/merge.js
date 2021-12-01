export const merge = (...objects) => {
    let target = {};

    let merger = (obj) => {
        for (let prop in obj) {
            if (obj.hasOwnProperty(prop)) {
                if (Object.prototype.toString.call(obj[prop]) === '[object Object]') {
                    target[prop] = merge(target[prop], obj[prop]);
                } else {
                    target[prop] = obj[prop];
                }
            }
        }
    };

    for (let i = 0; i < objects.length; i++) {
        merger(objects[i]);
    }

    return target;
};