import {isObject, isString} from "../helper/is";
import styles from "./ui.module.sass"
import errorField from "../components/unknown";
import defaults from "../components";

const buildRenderTree = (ui, fields, components) => {
    let children = []
    for (let elem of ui) {
        if (isString(elem)) {
            let child = renderField(elem, getFieldRenderFunc(components, fields[elem]))
            children.push(child)
        } else if (isObject(elem) && isString(elem.field)) {
            let child = renderField(elem.field, getFieldRenderFunc(components, {...fields[elem.field], ui: elem}))
            children.push(child)
        } else if (isObject(elem) && isString(elem.layout)) {
            let child = renderLayout(elem.layout, buildRenderTree(elem.children, fields, components))
            children.push(child)
        } else {
            console.error("unsupported ui schema")
        }
    }

    return children
}

const getFieldRenderFunc = (components, {common, options = {}, ui = {}}) => {
    let config = {common, options, ui}
    let compVariant = "default"
    if (ui.component) {
        compVariant = ui.component
    }

    let compTypes = components[common.type]
    if (!compTypes) {
        config.options.message = "Unknown component type: " + common.type
        return errorField(config)
    }

    let comp = compTypes[compVariant]
    if (!comp) {
        config.options.message = "Unknown component variant for type " + common.type + ": " + compVariant
        return errorField(config)
    }

    return comp(config)
}

const renderLayout = (layout, children) => (runtime, handlers, data) => {
    let rendered = children.map((child) => child(runtime, handlers, data))
    switch (layout) {
        case "row":
            return <div className={styles.FormLayoutRow}>{rendered}</div>
        case "column":
            return <div className={styles.FormLayoutColumn}>{rendered}</div>
        default:
            return errorField({options: {message: "Unknown layout"}})
    }

}

const renderField = (id, renderFunc) => (runtime, handlers, data) => {
    let rendered = renderFunc(runtime[id], {
        onChange: handlers.onChange(id),
    }, data)

    if (isString(rendered)) {
        return renderRuntimeError([rendered])(runtime, handlers, data)
    }

    return rendered
}

const getDefaultUI = (fields) => [
    {
        layout: "column",
        children: Object.keys(fields).map((key) => key)
    }
]

const renderRuntimeError = (errors) => {
    let field = errorField({options: {message: "Errors during field render"}})
    return (runtime, handlers, data) => {
        runtime.errors = errors

        return field(runtime, handlers, data)
    }
}

export const buildUI = (ui, fields, components) => {
    if (!fields) {
        return (runtime, handlers, data) => []
    }

    if (!components) {
        components = defaults
    }

    if (!ui) {
        ui = getDefaultUI(fields)
    }

    let tree = buildRenderTree(ui, fields, components)
    return (runtime, handlers, data) => tree.map((child) => child(runtime, handlers, data))
}

