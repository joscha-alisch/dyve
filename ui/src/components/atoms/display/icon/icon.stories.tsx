import * as React from "react"
import Icon from ".";
import { ComponentMeta, ComponentStory} from "@storybook/react";

export default {
    title: 'Components/Atoms/Display/Icon',
    component: Icon,
} as ComponentMeta<typeof Icon>

export const StoryIcon : ComponentStory<typeof Icon> = (args) => <Icon {...args}/>

StoryIcon.storyName = "Icon"
StoryIcon.args = {
    icon: "plus",
}