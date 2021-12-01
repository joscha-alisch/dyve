import ArrayInput from "../../inputs/arrayInput/arrayinput";


export const arrayField = ({common, options, ui}) => ({value, errors, dirty, touched}, {onChange}, data) => {
    return <ArrayInput label={common.label} helperText={common.hint} value={value} onChange={onChange}/>
}