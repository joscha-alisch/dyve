import React, {useEffect, useState} from "react"
import styles from "./detailpage.module.sass"
import PropTypes from "prop-types"
import {Link, useParams} from "react-router-dom";
import Page from "../page/page";
import {Button} from "@mui/material";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faEdit, faSync, faTrash} from "@fortawesome/free-solid-svg-icons";
import {ReadyState} from "react-use-websocket";

const DetailPage = ({
    detailApi,
    useLive = (id, setDetail) => {return {status: "connecting", send: () => {}}},
    render,
    title = (detail) => detail && detail.name,
    className,
    parent,
    parentRoute,
    editRoute,
    editLabel = "Edit",
    deleteEnabled = false
}) => {
    const {id} = useParams()
    const [state, setState] = useState({})
    const {status, send} = useLive(id, (liveData) => setState({
        ...state,
        live: liveData
    }))

    useEffect(() => {
        detailApi(id, (json) => {
            setState({
                ...state,
                detail: json.result
            })
        })
    }, [id])

    let buttons = () => <>
        <span className={styles.Live + " " + (status === "Open" ? styles.Connected : styles.Connecting)}>{status === "Open" ? "Live" : "Connecting..."}
            <FontAwesomeIcon className={styles.Refresh} onClick={() => send("update")} icon={faSync}/>
            </span>

        {editRoute &&
        <Button className={styles.Button} startIcon={<FontAwesomeIcon icon={faEdit}/>} variant={"outlined"}
                component={Link}
                to={editRoute(id)}>{editLabel}</Button>}
        {deleteEnabled &&
        <Button className={styles.Button} startIcon={<FontAwesomeIcon icon={faTrash}/>} variant={"outlined"}
                color={"error"}>Delete</Button>}

    </>

    return <Page className={className} title={title(state.detail)} parentRoute={parentRoute} parent={parent}
                 buttonsRender={buttons}>
        {render && render(state.detail || {}, state.live || {})}
    </Page>
}

DetailPage.propTypes = {
    className: PropTypes.string,
}

export default DetailPage