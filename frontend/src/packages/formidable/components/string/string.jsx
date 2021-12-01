import {TextField} from "@mui/material";
import React from "react";

export const getStringOptions = ({
                                     multiline = false,
                                     ...rest
                                 }) => {
    return {
        multiline
    }
}

const stringField = ({common, options}) => ({value, errors, dirty, touched}, {onChange}) => {
    return <TextField
        value={value}
        error={errors.length > 0}
        name={common.field}
        label={common.label}
        helperText={errors.length > 0 ? errors[0] : common.hint}
        multiline={options.multiline}
        onChange={(event) => onChange(event.target.value)}
    />
}

export default stringField