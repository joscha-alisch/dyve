import {Fragment, useEffect, useState} from "react";
import styles from "./applist.module.sass"
import { useQueryParam, NumberParam, withDefault } from 'use-query-params';
import AppCard from "../appcard/appcard";
import ListControl from "../listcontrol/listcontrol";
import Heading from "../heading/heading";
import axios from "axios";
import {useAuth} from "../../context/auth";

const AppList = () => {
    let [apps, setApps] = useState([])
    let [page] = useQueryParam("page", withDefault(NumberParam, 1))
    let [perPage, setPerPage] = useQueryParam("perPage", withDefault(NumberParam, 20))
    let [totalPages, setTotalPages] = useState(0)
    let [totalResults, setTotalResults] = useState(0)
    let [loading, setLoading] = useState(true)

    useEffect(() => {
        setLoading(true)
        axios.get("/api/apps?perPage=" + perPage + "&page=" + (page-1))
            .then((res) => {
                if(res.data.result.apps) {
                    setApps(res.data.result.apps)
                    setTotalPages(res.data.result.totalPages)
                    setTotalResults(res.data.result.totalResults)
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
        cards = apps.map((app) => <AppCard className={styles.AppCard} app={app}/>)
    }

    return <Fragment>
        <Heading title="Apps"/>
        {paginationControl}
        {cards}
        {paginationControl}
    </Fragment>
}

export default AppList