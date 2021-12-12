import React, { useState } from "react"
import HorizontalSelect from ".";
import { ComponentMeta, ComponentStory } from "@storybook/react";

export default {
    title: 'Components/Atoms/Input/Horizontal Select',
    component: HorizontalSelect,
} as ComponentMeta<typeof HorizontalSelect>

export const StoryHorizontalSelect: ComponentStory<typeof HorizontalSelect> = (args) => {
    let [selected, setSelected] = useState(args.options[0].value)
    
    return <HorizontalSelect {...args} selected={selected} onSelect={setSelected}/>
}

StoryHorizontalSelect.storyName = "Horizontal Select"
StoryHorizontalSelect.args = {
    label: "Cats",
    options: [
        { label: "10", value: 10 },
        { label: "50", value: 50 },
        { label: "100", value: 100 },
        { label: "All", value: -1 },
    ]
}