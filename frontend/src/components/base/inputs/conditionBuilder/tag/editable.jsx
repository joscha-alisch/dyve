import React, {useState} from "react";
import {TagEditor} from "./editor";
import {Tag} from "./tag";

export const EditableTag = ({value, onChange, onDelete, component, options}) => {
    const [editMode, setEditMode] = useState(false)

    if (editMode) {
        return <TagEditor options={options} onCancel={() => {
            setEditMode(false)
        }} value={value} onSubmit={(item) => {
            setEditMode(false)
            onChange(item)
        }}/>
    }

    return <Tag onDelete={onDelete} onClick={() => setEditMode(true)} value={value} component={component}/>
}
