import React from "react"
import TextField from "./textField";
import { ComponentMeta, ComponentStory} from "@storybook/react";

export default {
    title: 'Components/Atoms/Input/Text Field',
    component: TextField,
} as ComponentMeta<typeof TextField>

export const StoryTextField : ComponentStory<typeof TextField> = (args) => <TextField {...args}/>

StoryTextField.storyName = "Text Field"
StoryTextField.args = {
    multiLine: false,
    lines: 10,
}