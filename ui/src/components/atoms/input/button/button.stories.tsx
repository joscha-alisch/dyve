import * as React from "react"
import Button from "./button";
import { ComponentMeta, ComponentStory} from "@storybook/react";

export default {
    title: 'Components/Atoms/Input/Button',
    component: Button,
} as ComponentMeta<typeof Button>

export const StoryButton : ComponentStory<typeof Button> = (args) => <Button {...args}/>

StoryButton.storyName = "Button"
StoryButton.args = {
    title: "Click Me"
}