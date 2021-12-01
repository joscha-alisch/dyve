import ArrayInput from "../../inputs/arrayInput/arrayinput";
import ConditionBuilder from "../../inputs/conditionBuilder/conditionbuilder";
import {Chip} from "@mui/material";


export const conditionsField = ({common, options, ui}) => ({value, errors, dirty, touched}, {onChange}, data) => {
    return <ArrayInput
        label={common.label}
        helperText={common.hint}
        value={value}
        onChange={onChange}
        newItem={() => {
            return []
        }}
        divider={(args) => <Chip {...args} label={"or"}/>}
        renderField={({value, onChange}) => <ConditionBuilder value={value} onChange={onChange}/>}
    />
}