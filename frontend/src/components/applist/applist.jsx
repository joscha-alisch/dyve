import {Fragment, useEffect, useState} from "react";
import styles from "./applist.module.sass"
import { useQueryParam, NumberParam, withDefault } from 'use-query-params';
import {Link} from "react-router-dom";
import Pagination from '@material-ui/lab/Pagination';
import {PaginationItem, Skeleton} from "@material-ui/lab";
import {FormControl, InputLabel, MenuItem, Select} from "@material-ui/core";
import AppCard from "../appcard/appcard";
import ListContext from "@material-ui/core/List/ListContext";
import ListControl from "../listcontrol/listcontrol";

const AppList = () => {
    let [apps, setApps] = useState([])
    let [page, setPage] = useQueryParam("page", withDefault(NumberParam, 1))
    let [perPage, setPerPage] = useQueryParam("perPage", withDefault(NumberParam, 20))
    let [totalPages, setTotalPages] = useState(0)
    let [totalResults, setTotalResults] = useState(0)
    let [loading, setLoading] = useState(true)

    useEffect(() => {
        setLoading(true)
        fetch("/api/apps?perPage=" + perPage + "&page=" + (page-1))
            .then(res => res.json())
            .then((data) => {
                setApps(data.result.apps)
                setTotalPages(data.result.totalPages)
                setTotalResults(data.result.totalResults)
                setLoading(false)
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
        <h1>Apps</h1>
        {paginationControl}
        {cards}
        {paginationControl}
    </Fragment>
}

export default AppList