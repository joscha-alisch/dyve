import React, { useState } from "react"
import ListHeader from ".";
import { ComponentMeta, ComponentStory } from "@storybook/react";

export default {
    title: 'Components/Organisms/Headers/List',
    component: ListHeader,
} as ComponentMeta<typeof ListHeader>

export const StoryListHeader: ComponentStory<typeof ListHeader> = (args) => {
    let [filters, setFilters] = useState(args.filters)
    let [pagination, setPagination] = useState(args.pagination)

    return <ListHeader {...args} filters={filters} onFilterChange={setFilters} pagination={pagination} onPaginationChange={setPagination} />
}

StoryListHeader.storyName = "List"
StoryListHeader.args = {
    title: "Apps",
    category: "Platform",
    filters: [
        { key: "key1", value: "value1" },
        { key: "key2", value: "value2" },
        { key: "key3", value: "value3" },
        { key: "key4", value: "value4" },
        { key: "key5", value: "value5" },
    ],
    pagination: {
        page: 0,
        perPage: 10,
        totalItems: 4183
    }
}