import PropTypes from "prop-types";
import {useEffect, useState} from "react";
import {useFormData} from "../data/formData";
import {buildUI} from "../ui/ui";

export const Form = ({
                         className,
                         schema,
                         ui,
                         state,
                         setState,
                         components,
                         data
                     }) => {
    let {fields, runtime, handlers} = useFormData(schema, state, setState)
    let [renderUi, setRenderUi] = useState(() => (runtime, handlers, data) => <></>)

    useEffect(() => {
        setRenderUi(() => buildUI(ui, fields, components))
    }, [ui, fields])

    return <div>
        {renderUi(runtime, handlers, data)}
    </div>
}

Form.propTypes = {
    className: PropTypes.string,
    schema: PropTypes.object,
    fields: PropTypes.object
}