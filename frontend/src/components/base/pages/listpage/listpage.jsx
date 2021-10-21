import React, {useEffect, useState} from "react"
import styles from "./listpage.module.sass"
import PropTypes, {string} from "prop-types"
import {NumberParam, useQueryParam, withDefault} from "use-query-params";
import {Skeleton} from "@mui/lab";
import Box from "../../box/box";
import PaginationControl from "../../inputs/paginationcontrol/paginationControl";
import Page from "../page/page";

const ListPage = ({className, title, parent, fetchItems, itemRender, skeletonRender}) => {
    if (!skeletonRender) {
        skeletonRender = defaultSkeleton
    }

    let [items, setItems] = useState([])
    let [page] = useQueryParam("page", withDefault(NumberParam, 0))
    let [perPage, setPerPage] = useQueryParam("perPage", withDefault(NumberParam, 20))
    let [totalResults, setTotalResults] = useState(0)
    let [loading, setLoading] = useState(true)

    let setResults = (items, totalResults) => {
        setItems(items)
        setTotalResults(totalResults)
        setLoading(false)
    }

    useEffect(() => {
        setLoading(true)
        fetchItems(perPage, page, setResults)
    }, [page, perPage, fetchItems])

    let content
    if (loading) {
        content = listPageLoading(skeletonRender, perPage)
    } else {
        content = items.map((item) => <li className={styles.ListItem}>
            {itemRender(item)}
        </li>)
    }

    return <Page className={styles.Main + " " + className} title={title} parent={parent}>
        <PaginationControl totalResults={totalResults} page={page} perPage={perPage}
                     setPerPage={setPerPage}/>
        <ul className={styles.List}>
            {content}
        </ul>
    </Page>
}

const listPageLoading = (skeletonRender, perPage) => {
    let content = []
    for (let i = 0; i < perPage; i++) {
        content.push(skeletonRender())
    }
    return content
}

const defaultSkeleton = () => <li className={styles.ListItem}>
    <Box>
        <Skeleton className={styles.Skeleton} animation={"wave"} variant={"text"} width="30%" height={30}/>
        <Skeleton className={styles.Skeleton}  animation={"wave"} variant={"text"} width="15%"/>
    </Box>
</li>

ListPage.propTypes = {
    parent: string,
    title: string,
    className: PropTypes.string,
    fetchItems: PropTypes.func,
    itemRender: PropTypes.func
}


export default ListPage