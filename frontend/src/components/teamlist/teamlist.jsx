import {Fragment, useEffect, useState} from "react";
import styles from "./teamlist.module.sass"
import { useQueryParam, NumberParam, withDefault } from 'use-query-params';
import AppCard from "../appcard/appcard";
import ListControl from "../listcontrol/listcontrol";
import Heading from "../heading/heading";
import axios from "axios";
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome'
import {faPlusCircle} from "@fortawesome/free-solid-svg-icons/faPlusCircle";
import {Link} from "react-router-dom";


const TeamList = () => {
    let [teams, setTeams] = useState([])
    let [page] = useQueryParam("page", withDefault(NumberParam, 1))
    let [perPage, setPerPage] = useQueryParam("perPage", withDefault(NumberParam, 20))
    let [totalPages, setTotalPages] = useState(0)
    let [totalResults, setTotalResults] = useState(0)
    let [loading, setLoading] = useState(true)

    useEffect(() => {
        setLoading(true)
        axios.get("/api/teams?perPage=" + perPage + "&page=" + (page-1))
            .then((res) => {
                if(res.data.result.teams) {
                    setTeams(res.data.result.teams)
                    setTotalPages(res.data.result.totalPages)
                    setTotalResults(res.data.result.totalResults)
                    setLoading(false)
                } else {
                    setTeams([])
                    setLoading(false)
                }
            })
    }, [page, perPage])

    let paginationControl = <ListControl totalResults={totalResults} totalPages={totalPages} page={page} perPage={perPage} setPerPage={setPerPage} />

    let cards
    if (loading) {
        cards = <Fragment>
            <AppCard className={styles.AppCard} loading/>
            <AppCard className={styles.AppCard} loading/>
            <AppCard className={styles.AppCard} loading/>
            <AppCard className={styles.AppCard} loading/>
            <AppCard className={styles.AppCard} loading/>
        </Fragment>
    } else {
        cards = teams.map((app) => <AppCard className={styles.AppCard} app={app}/>)
    }

    return <Fragment>
        <Heading title="Teams"/>
        <Link to="/teams/new"><FontAwesomeIcon icon={faPlusCircle}/> New Team</Link>
        {teams.length > perPage ? paginationControl : "" }
        {teams.length > 0 ? cards : <div className={styles.NoContent}>No Teams</div> }
        {teams.length > perPage ? paginationControl : ""}
    </Fragment>
}

export default TeamList