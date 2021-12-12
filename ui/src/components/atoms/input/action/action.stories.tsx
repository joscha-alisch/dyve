import React from "react"
import Action from ".";
import { ComponentMeta, ComponentStory} from "@storybook/react";

export default {
    title: 'Components/Atoms/Input/Action',
    component: Action,
} as ComponentMeta<typeof Action>

export const StoryAction : ComponentStory<typeof Action> = (args) => <Action {...args}/>

StoryAction.storyName = "Action"
StoryAction.args = {
    icon: "plus",
    label: "Add Item",
}