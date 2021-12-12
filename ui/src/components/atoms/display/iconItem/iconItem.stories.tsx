import React from "react"
import IconItem from ".";
import { ComponentMeta, ComponentStory} from "@storybook/react";

export default {
    title: 'Components/IconItem',
    component: IconItem,
} as ComponentMeta<typeof IconItem>

export const StoryIconItem : ComponentStory<typeof IconItem> = (args) => <IconItem {...args}/>

StoryIconItem.storyName = "IconItem"
StoryIconItem.args = {
    icon:"server",
    label: "Server",
}