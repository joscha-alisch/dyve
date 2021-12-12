import React, { useState } from "react"
import PageCounter from ".";
import { ComponentMeta, ComponentStory} from "@storybook/react";

export default {
    title: 'Components/Atoms/Input/Page Counter',
    component: PageCounter,
} as ComponentMeta<typeof PageCounter>

export const StoryPageCounter : ComponentStory<typeof PageCounter> = (args) => {
    let [state, setState] = useState(args.page)
    return <PageCounter {...args} page={state} onPageChange={setState} />
}

StoryPageCounter.storyName = "Page Counter"
StoryPageCounter.args = {
    totalItems:  2000,
    page: 0,
    perPage: 20
}