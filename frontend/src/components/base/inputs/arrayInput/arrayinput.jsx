import React from "react"
import styles from "./arrayinput.module.sass"
import PropTypes from "prop-types"
import {useControlledState} from "../../../../hooks/useControlledState";
import {FormHelperText, IconButton, InputLabel, TextField} from "@mui/material";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faMinusCircle, faPlusCircle} from "@fortawesome/free-solid-svg-icons";


const defaultRenderField = ({value, onChange}) => <TextField value={value} onChange={(e) => onChange(e.target.value)}/>

const ArrayInput = ({
                        label,
                        helperText,
                        className,
                        renderField = defaultRenderField,
                        divider,
                        newItem = () => "",
                        value,
                        onChange
                    }) => {
    let [state, setState] = useControlledState([], value, onChange)

    const add = () => {
        setState([...state, newItem()])
    }

    const update = (index) => (value) => {
        let newState = [...state]
        newState[index] = value
        setState(newState)
    }

    const remove = (index) => () => {
        let newState = [...state]
        newState.splice(index, 1)
        setState(newState)
    }

    return <div className={styles.Main + " " + className}>
        {label && <InputLabel>{label}</InputLabel>}
        {helperText && <FormHelperText>{label}</FormHelperText>}
        <ul className={styles.List}>
            {!state || state.length === 0 ? <li className={styles.EmptyText}>Empty</li> : ""}
            {state.map((item, index) => <>{divider && index > 0 ?
                <li key={"divider" + index} className={styles.Divider}>{divider()}</li> : ""}
                <li className={styles.Item} key={index}>
                    <div className={styles.Field}>{renderField({value: item, onChange: update(index)})}</div>
                    <IconButton className={styles.RemoveButton} color={"error"} onClick={remove(index)}><FontAwesomeIcon
                        size={"xs"} icon={faMinusCircle}/></IconButton>
                </li>
            </>)}
            <li className={styles.Item + " " + styles.AddButtonItem}>
                <IconButton color={"primary"} onClick={add}><FontAwesomeIcon size={"1x"}
                                                                             icon={faPlusCircle}/></IconButton>
            </li>
        </ul>

    </div>
}

ArrayInput.propTypes = {
    className: PropTypes.string,
}

export default ArrayInput