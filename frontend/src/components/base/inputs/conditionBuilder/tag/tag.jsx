import {Chip} from "@mui/material";
import styles from "../conditionbuilder.module.sass";
import React from "react";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faMinusCircle} from "@fortawesome/free-solid-svg-icons";

export const Tag = ({value, onDelete, onClick, component}) => <Chip
    className={styles.Chip}
    color={"primary"}
    clickable
    draggable
    deleteIcon={<FontAwesomeIcon size={"xs"} icon={faMinusCircle}/>}
    onDelete={onDelete}
    onClick={onClick}
    label={component ? component() : <>
        <b>{value.key}: </b>{value.value}
    </>}
/>