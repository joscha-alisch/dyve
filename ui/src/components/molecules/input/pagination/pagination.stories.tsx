import React, { useState } from "react"
import Pagination from ".";
import { ComponentMeta, ComponentStory} from "@storybook/react";

export default {
    title: 'Components/Molecules/Input/Pagination',
    component: Pagination,
} as ComponentMeta<typeof Pagination>

export const StoryPagination : ComponentStory<typeof Pagination> = (args) => {
    let [state, setState] = useState(args.value)
    return <Pagination {...args} value={state} onChange={setState}/>
}
StoryPagination.storyName = "Pagination"
StoryPagination.args = {
    value: {
        page: 0,
        perPage: 10,
        totalItems: 2124
    }
}