import React, { useState } from "react"
import FilterInput from ".";
import { ComponentMeta, ComponentStory} from "@storybook/react";

export default {
    title: 'Components/Atoms/Input/Filter Input',
    component: FilterInput,
} as ComponentMeta<typeof FilterInput>

export const StoryFilterInput : ComponentStory<typeof FilterInput> = (args) => {
    let [state, setState] = useState(args.value)

    return <FilterInput {...args} value={state} onChange={setState} />
}

StoryFilterInput.storyName = "Filter Input"
StoryFilterInput.args = {
    value: "",
}