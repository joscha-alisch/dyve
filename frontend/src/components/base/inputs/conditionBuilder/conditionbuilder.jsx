import React from "react"
import styles from "./conditionbuilder.module.sass"
import PropTypes from "prop-types"
import {Chip, FormHelperText, FormLabel, IconButton, List, Paper} from "@mui/material";
import {faPlus} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome"
import {useControlledState} from "../../../../hooks/useControlledState";
import {EditableTag} from "./tag/editable";


const ConditionBuilder = ({
                              className,
                              separator = " and ",
                              options = {},
                              helperText,
                              label,
                              value,
                              onChange
                          }) => {

    const [state, setState] = useControlledState([], value, onChange)

    const add = (item) => {
        setState([...state, item])
    }

    const onDelete = (index) => () => {
        let newState = [...state]
        newState.splice(index, 1)
        setState(newState)
    }

    const update = (index) => (item) => {
        let newState = [...state]
        newState[index] = item
        setState(newState)
    }

    return <>
        {label && <FormLabel>{label}</FormLabel>}
        <Paper
            className={styles.Main + " " + className}
            component="form"
            elevation={1}
            sx={{backgroundColor: "#2a3543", p: '2px 20px', display: 'flex', alignItems: 'center'}}>
            <List className={styles.List} sx={{flex: "1 1 400px"}}>
                {state.map((item, index) => <>
                    {index > 0 && separator !== "" ? <Chip className={styles.Chip} label={separator}/> : ""}
                    <EditableTag options={options} onDelete={onDelete(index)} value={item} onChange={update(index)}/>
                </>)}
                <EditableTag onChange={add}
                             options={options}
                             component={() => <IconButton
                                 sx={{padding: "5px 20px"}}
                                 color={"success"}>
                                 <FontAwesomeIcon icon={faPlus}/>
                             </IconButton>}/>
            </List>
        </Paper>
        {helperText && <FormHelperText>{helperText}</FormHelperText>}
    </>
}

ConditionBuilder.propTypes = {
    className: PropTypes.string,
}

export default ConditionBuilder