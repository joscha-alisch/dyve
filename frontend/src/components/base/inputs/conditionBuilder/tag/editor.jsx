import React, {useState} from "react";
import {Autocomplete, Chip, Divider, IconButton, TextField} from "@mui/material";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faCheckCircle, faTimesCircle} from "@fortawesome/free-solid-svg-icons";
import {useFocus} from "../../../../../hooks/useFocus";

export const TagEditor = ({
                              onSubmit = (() => {
                              }), onCancel, options = {}, value
                          }) => {
    const initial = value || {key: null, value: null}
    const [state, setState] = useState(initial)
    const [inputRef, setInputFocus] = useFocus()

    const canSubmit = state.key !== null && state.value !== null

    let submit = () => {
        if (!canSubmit) {
            return
        }

        onSubmit(state)
        setState(initial)
        setInputFocus()
    }

    return <Chip color={"primary"} label={<form onSubmit={submit} style={{display: "flex", alignItems: 'center'}}>
        <Autocomplete
            onSubmit={() => onSubmit(state)}
            disablePortal
            id="combo-box-demo"
            options={Object.keys(options)}
            value={state.key}
            disableClearable
            selectOnFocus
            openOnFocus
            clearOnBlur
            handleHomeEndKeys
            autoSelect
            autoComplete
            autoHighlight
            onChange={(e, v) => setState({...state, key: v, value: null})}
            renderInput={(params) =>
                <TextField {...params} variant={"standard"} disableUnderline sx={{width: 100}} inputRef={inputRef}
                           autoFocus placeholder={"Key"}/>
            }
        />
        <Divider sx={{p: '10px', marginRight: "10px"}} orientation={"vertical"}/>
        <Autocomplete
            disablePortal
            id="combo-box-demo"
            value={state.value}
            options={options[state.key] || []}
            disableClearable
            selectOnFocus
            clearOnBlur
            openOnFocus
            handleHomeEndKeys
            autoSelect
            autoComplete
            autoHighlight
            onChange={(e, v) => setState({...state, value: v})}
            renderInput={(params) =>
                <TextField {...params} variant={"standard"} sx={{width: 200}} placeholder={"Value"}/>
            }
        />
        <IconButton onClick={onCancel}><FontAwesomeIcon size={"xs"} icon={faTimesCircle}/></IconButton>
        <IconButton disabled={!canSubmit} onSubmit={submit} onClick={submit}><FontAwesomeIcon size={"xs"}
                                                                                              icon={faCheckCircle}/></IconButton>
    </form>}/>
}
