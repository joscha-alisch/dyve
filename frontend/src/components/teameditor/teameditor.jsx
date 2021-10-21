import {Fragment, useEffect, useState} from "react";
import styles from "./teameditor.module.sass"
import AppCard from "../appcard/appcard";
import ListControl from "../listcontrol/listcontrol";
import Heading from "../heading/heading";
import axios from "axios";
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome'
import {faPlusCircle} from "@fortawesome/free-solid-svg-icons/faPlusCircle";
import {Link, useParams} from "react-router-dom";
import {Input, TextField} from "@mui/material";


const TeamEditor = ({newTeam}) => {
    const {id} = useParams()
    let [team, setTeam] = useState({})
    let [loading, setLoading] = useState(true)

    useEffect(() => {
        if (!newTeam) {
            setLoading(true)
            axios.get("/api/teams/" + id)
                .then((res) => {
                    if (res.data.result) {
                        setTeam(res.data.result)
                        setLoading(false)
                    }
                })
        }

    }, [])

    return <Fragment>
        <Heading title={newTeam ? "New Team" : team.name} backlink="/teams" backlinkTitle="Teams"/>

        <TextField id="team_editor_name" required label="Team Name"  helperText="Name people will see displayed in dyve. Should be the canonical, descriptive name of your team."/>
        <TextField id="team_editor_slug" required label="Team Slug"  helperText="Name used in links and internally. Should be concise and without spaces."/>


    </Fragment>
}

export default TeamEditor