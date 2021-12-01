import AssociationsPicker from "../../inputs/associations/associationspicker";


export const associationsField = ({common, options, ui}) => ({value, errors, dirty, touched}, {onChange}, data) => {
    let fieldConfig = ui.fields
    let fields = []

    for (const field of fieldConfig) {
        if (field.dataKey && (!data || !data[field.dataKey])) {
            return "Field is using unknown dataKey '" + field.dataKey + "'"
        }

        fields.push({
            label: field.label,
            data: field.dataKey ? data[field.dataKey] : field.data,
            optionLabel: field.optionLabel,
            outputKey: field.outputKey
        })
    }

    return <AssociationsPicker
        onChange={onChange}
        value={value}
        fields={fields}
    />
}