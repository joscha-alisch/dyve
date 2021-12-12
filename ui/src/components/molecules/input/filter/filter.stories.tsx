import React, { useState } from "react"
import Filter from "./filter";
import { ComponentMeta, ComponentStory } from "@storybook/react";

export default {
    title: 'Components/Molecules/Input/Filter',
    component: Filter,
} as ComponentMeta<typeof Filter>

export const StoryFilter: ComponentStory<typeof Filter> = (args) => {
    let [state, setState] = useState({
        key: args.filterKey,
        value: args.filterValue,
    })
    let [open, setOpen] = useState(false)

    return <Filter {...args} filterKey={state.key} filterValue={state.value} onChange={(key, value) => {
        setState({
            key: key,
            value: value
        })
    }}  onOpen={() => setOpen(true)} onClose={() => setOpen(false)} open={open}/>
}
StoryFilter.storyName = "Filter"
StoryFilter.args = {
    filterKey: "key",
    filterValue: "valueee"
}