import React, {useEffect, useState} from "react"
import styles from "./listpage.module.sass"
import PropTypes, {string} from "prop-types"
import {NumberParam, useQueryParam, withDefault} from "use-query-params";
import {Button, Skeleton} from "@mui/material";
import Box from "../../box/box";
import PaginationControl from "../../inputs/paginationcontrol/paginationControl";
import Page from "../page/page";
import {Link} from "react-router-dom";
import {faPlusCircle} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";

const ListPage = ({
                      className,
                      newItemRoute,
                      newItemLabel = "New",
                      title,
                      parent,
                      fetchItems,
                      itemRender,
                      skeletonRender
                  }) => {
    if (!skeletonRender) {
        skeletonRender = defaultSkeleton
    }

    let [items, setItems] = useState([])
    let [page, setPage] = useQueryParam("page", withDefault(NumberParam, 0))
    let [perPage, setPerPage] = useQueryParam("perPage", withDefault(NumberParam, 25))
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
    } else if (!items) {
        content = listPageEmpty()
    } else {
        content = items.map((item) => <li className={styles.ListItem}>
            {itemRender({value: item})}
        </li>)
    }

    let buttonsRender = () => <>
        {newItemRoute &&
        <Button startIcon={<FontAwesomeIcon icon={faPlusCircle}/>} variant={"contained"} component={Link}
                to={newItemRoute}>{newItemLabel}</Button>}
    </>

    return <Page className={styles.Main + " " + className} buttonsRender={buttonsRender} title={title} parent={parent}>
        <PaginationControl totalResults={totalResults} page={page} perPage={perPage}
                           setPerPage={setPerPage} setPage={setPage}/>
        <ul className={styles.List}>
            {content}
        </ul>
    </Page>
}

const listPageEmpty = () => <li className={styles.Empty}>There are no items in this list or you don't have enough
    permissions.</li>

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
        <Skeleton className={styles.Skeleton} animation={"wave"} variant={"text"} width="15%"/>
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