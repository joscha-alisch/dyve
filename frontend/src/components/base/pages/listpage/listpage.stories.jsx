import React from "react"
import ListPage from "./listpage";
import {QueryParamProvider} from "use-query-params";
import Box from "../../box/box";
import {Route} from "react-router-dom";

export default {
    title: 'App/Pages/ListPage',
    component: ListPage,
}

export const StoryListPage = (args) => <QueryParamProvider ReactRouterRoute={Route}>
    {args.loading ? <ListPage {...args} fetchItems={dontReturn}/> : <ListPage {...args} />}
</QueryParamProvider>

const ItemRender = (item) => <Box title="PipelineList Item">{item.value}</Box>

const returnItems = (perPage, page, setResults) => {
    let results = []
    let max = 1234
    let start = (perPage * page) + 1
    for (let i = start; i < start + perPage; i++) {
        results.push("Item " + i)
        if (i === max) {
            break
        }
    }

    setResults(results, max)
}
const dontReturn = (perPage, page, setResults) => {
}

StoryListPage.storyName = "ListPage"
StoryListPage.args = {
    loading: false,
    title: "List Page",
    parent: "Parent Page",
    fetchItems: returnItems,
    itemRender: ItemRender
}