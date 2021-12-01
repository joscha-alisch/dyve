import React, {useState} from "react"
import styles from "./associationspicker.module.sass"
import PropTypes from "prop-types"
import {
    Autocomplete,
    Button,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    TextField
} from "@mui/material";
import {faMinusCircle, faPlusCircle} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome"

const AssociationsPicker = ({className, fields, value, onChange}) => {
    let [state, setState] = useState([])

    if (!fields) {
        fields = []
    }

    if (!onChange) {
        onChange = setState
    }

    if (!value) {
        value = state
    }

    const add = () => {
        let newItem = {}
        fields.forEach((field) => newItem[field.outputKey] = null)

        onChange([...value, newItem])
    }

    const remove = (index) => {
        let newState = [...value]
        newState.splice(index, 1);
        onChange(newState)
    }

    let onChangeHandler = (index, prop, newValue) => {
        let newState = [...value]
        newState[index][prop] = newValue
        onChange(newState)
    }


    return <div className={styles.Main + " " + className}>
        <TableContainer>
            <Table sx={{minWidth: 650}} aria-label="simple table">
                <TableHead>
                    <TableRow>
                        {fields.map(field => <TableCell>{field.label}</TableCell>)}
                        <TableCell align="right"/>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {value.map((item, index) => <TableRow key={index}>
                        {fields.map(field => {
                            return <TableCell>
                                <Autocomplete
                                    disablePortal
                                    id="combo-box-demo"
                                    options={field.data}
                                    sx={{width: 300}}
                                    value={value[index][field.outputKey]}
                                    getOptionLabel={(option) => field.optionLabel ? option[field.optionLabel] : option}
                                    onChange={(_, value) => onChangeHandler(index, field.outputKey, value)}
                                    renderInput={(params) => <TextField {...params} label={field.label}/>}
                                />
                            </TableCell>
                        })}
                        <TableCell><Button onClick={() => remove(index)}><FontAwesomeIcon size="2x"
                                                                                          icon={faMinusCircle}/></Button></TableCell>
                    </TableRow>)}
                    <TableRow>
                        <TableCell colSpan={3} align={"center"}><Button onClick={() => add()}><FontAwesomeIcon
                            size={"2x"} icon={faPlusCircle}/></Button></TableCell>
                    </TableRow>
                </TableBody>
            </Table>
        </TableContainer>

    </div>
}


AssociationsPicker.propTypes = {
    className: PropTypes.string,
}

export default AssociationsPicker