import {Fragment, useEffect, useState} from "react";
import styles from "./pipelinelist.module.sass"
import { useQueryParam, NumberParam, withDefault } from 'use-query-params';
import {Link} from "react-router-dom";
import Pagination from '@material-ui/lab/Pagination';
import {PaginationItem, Skeleton} from "@material-ui/lab";
import {FormControl, InputLabel, MenuItem, Select} from "@material-ui/core";
import AppCard from "../appcard/appcard";
import ListContext from "@material-ui/core/List/ListContext";
import ListControl from "../listcontrol/listcontrol";
import Heading from "../heading/heading";
import PipelineCard from "../pipelinecard/pipelineCard";

const Pipelinelist = () => {
    let [pipelines, setPipelines] = useState([])
    let [page, setPage] = useQueryParam("page", withDefault(NumberParam, 1))
    let [perPage, setPerPage] = useQueryParam("perPage", withDefault(NumberParam, 20))
    let [totalPages, setTotalPages] = useState(0)
    let [totalResults, setTotalResults] = useState(0)
    let [loading, setLoading] = useState(true)
    let [svg, setSvg] = useState({})

    useEffect(() => {
        setLoading(true)
        fetch("/api/pipelines?perPage=" + perPage + "&page=" + (page-1))
            .then(res => res.json())
            .then((data) => {
                setPipelines(data.result.pipelines)
                console.log(data.result)
                setTotalPages(data.result.totalPages)
                setTotalResults(data.result.totalResults)
                setLoading(false)
            })
    }, [page, perPage])

    useEffect(() => {
        pipelines.forEach(p => fetch("/api/pipelines/" + p.id + "/status")
            .then(res => res.json())
            .then((data) => setSvg(prev => ({...prev, [p.id]: data.result.svg})))
        )
    }, [pipelines])

    let paginationControl = <ListControl totalResults={totalResults} totalPages={totalPages} page={page} perPage={perPage} setPerPage={setPerPage} />

    let cards
    if (loading) {
        cards = <Fragment>
            <PipelineCard className={styles.AppCard} loading/>
            <PipelineCard className={styles.AppCard} loading/>
            <PipelineCard className={styles.AppCard} loading/>
            <PipelineCard className={styles.AppCard} loading/>
            <PipelineCard className={styles.AppCard} loading/>
        </Fragment>
    } else {
        cards = pipelines.map((pipeline) => <li key={pipeline.id} className={styles.ListItem}><PipelineCard
            className={styles.AppCard} pipeline={pipeline} svg={svg}/></li>)
    }

    return <Fragment>
        <Heading title="Pipelines"/>
        {paginationControl}
        {cards}
        {paginationControl}
    </Fragment>
}

export default Pipelinelist