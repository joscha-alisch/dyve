import React from "react"
import ArrayInput from "./arrayinput";

export default {
    title: 'Components/Inputs/Array',
    component: ArrayInput,
}

export const StoryArrayInput = (args) => <ArrayInput {...args}/>

StoryArrayInput.storyName = "Array"
StoryArrayInput.args = {}