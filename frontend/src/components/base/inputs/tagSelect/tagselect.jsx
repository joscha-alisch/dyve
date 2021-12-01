import React from "react"
import PropTypes from "prop-types"
import {Autocomplete, Chip, TextField} from "@mui/material";

const TagSelect = ({
    className,
    label,
    helperText,
    options = [],
    groupBy,
    getOptionLabel = (option) => option,
    value,
    onChange,
    style
}) => {
    let renderInput = (params) => <TextField {...params} style={style} label={label} helperText={helperText}/>

    return <Autocomplete
        multiple
        filterSelectedOptions
        disableClearable
        selectOnFocus
        openOnFocus
        clearOnBlur
        handleHomeEndKeys
        autoComplete
        autoHighlight
        disablePortal
        options={options}
        groupBy={groupBy}
        getOptionLabel={getOptionLabel}
        renderInput={renderInput}
        value={value}
        onChange={(e, v) => onChange(v)}
        renderTags={(value, getTagProps) =>
            value.map((option, index) => (
                <Chip variant="filled" color="primary" label={getOptionLabel(option)} {...getTagProps({ index })} />
            ))
        }
    />
}

TagSelect.propTypes = {
    className: PropTypes.string,
    options: PropTypes.array,
    groupBy: PropTypes.func,
    getOptionLabel: PropTypes.func,
    onChange: PropTypes.func
}

export default TagSelect