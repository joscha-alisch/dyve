import React, {useEffect} from "react"
import PaginationControl from "./paginationControl";
import {NumberParam, QueryParamProvider, useQueryParam, withDefault} from "use-query-params";
import set from "cytoscape/src/set";
import { Route } from 'react-router-dom';
import {withQuery} from "@storybook/addon-queryparams";

export default {
    title: 'Components/Inputs/Pagination Control',
    component: PaginationControl,
    isFullscreen: true,
}

export const StoryPaginationControl = (args) => <QueryParamProvider ReactRouterRoute={Route}>
    <Comp {...args}/>
</QueryParamProvider>

const Comp = (args) => {
    let [page, setPage] = useQueryParam("page", withDefault(NumberParam, 1))
    let [perPage, setPerPage] = useQueryParam("perPage", withDefault(NumberParam, 20))

    useEffect(() => {
        setPage(args.page)
        setPerPage(args.perPage)
    }, [args.page, args.perPage])

    return <PaginationControl totalResults={args.totalResults} page={page} perPage={perPage} setPerPage={() =>{}}/>
}

StoryPaginationControl.storyName = "Pagination Control"
StoryPaginationControl.args = {
    totalResults: 1000,
    perPage: 20,
    page: 0
}