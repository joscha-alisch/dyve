import React from "react"
import { ComponentMeta, ComponentStory} from "@storybook/react";
import ActionItem from ".";

export default {
    title: 'Components/Atoms/Input/Action Item',
    component: ActionItem,
} as ComponentMeta<typeof ActionItem>

export const StoryContextMenuItem : ComponentStory<typeof ActionItem> = (args) => <ActionItem {...args}/>

StoryContextMenuItem.storyName = "Action Item"
StoryContextMenuItem.args = {
    label: "Action Item",
    icon: "mail"
}