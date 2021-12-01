import ArrayInput from "../../inputs/arrayInput/arrayinput";
import ConditionBuilder from "../../inputs/conditionBuilder/conditionbuilder";
import {Chip, CircularProgress} from "@mui/material";
import TagSelect from "../../inputs/tagSelect/tagselect";


export const tagsField = ({common, options, ui}) => ({value, errors, dirty, touched}, {onChange}, data) => {
    if (!data || !data[ui.dataKey]) {
        return <CircularProgress />
    }

    let sortedOptions =  data[ui.dataKey].sort((a, b) => {
        let labelA = a[ui.labelKey].toUpperCase();
        let labelB = b[ui.labelKey].toUpperCase();
        if (labelA < labelB) {
            return -1;
        }
        if (labelA > labelB) {
            return 1;
        }
        return 0;
    })

    return <TagSelect
        label={common.label}
        helperText={common.hint}
        value={value}
        onChange={onChange}
        options={sortedOptions}
        groupBy={(option) => option[ui.groupKey]}
        getOptionLabel={(option) => option[ui.labelKey]}
        style={ui.style}
    />
}