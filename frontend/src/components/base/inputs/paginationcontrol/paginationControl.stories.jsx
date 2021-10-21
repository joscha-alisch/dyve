import React, {useEffect} from "react"
import PaginationControl from "./paginationControl";
import {NumberParam, QueryParamProvider, useQueryParam, withDefault} from "use-query-params";
import {Route} from 'react-router-dom';

export default {
    title: 'Components/Inputs/Pagination Control',
    component: PaginationControl,
    isFullscreen: true,
}

export const StoryPaginationControl = (args) => <QueryParamProvider ReactRouterRoute={Route}>
    <Comp {...args}/>
</QueryParamProvider>

const Comp = (args) => {
    let [page, setPage] = useQueryParam("page", withDefault(NumberParam, args.page))
    let [perPage, setPerPage] = useQueryParam("perPage", withDefault(NumberParam, args.perPage))

    return <PaginationControl totalResults={args.totalResults} page={page} perPage={perPage} setPerPage={setPerPage} setPage={setPage}/>
}

StoryPaginationControl.storyName = "Pagination Control"
StoryPaginationControl.args = {
    totalResults: 1000,
    perPage: 25,
    page: 0
}