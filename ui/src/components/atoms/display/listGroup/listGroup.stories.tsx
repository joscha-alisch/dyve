import React from "react"
import ListGroup from ".";
import { ComponentMeta, ComponentStory} from "@storybook/react";
import ActionItem from "../../input/actionItem"

export default {
    title: 'Components/Atoms/Display/List Group',
    component: ListGroup,
} as ComponentMeta<typeof ListGroup>

export const StoryListGroup : ComponentStory<typeof ListGroup> = (args) => <ListGroup {...args}/>

StoryListGroup.storyName = "List Group"
StoryListGroup.args = {
    label: "Group",
    children: <>
        <ActionItem label="List Item 1"/>
        <ActionItem label="List Item 1"/>
        <ActionItem label="List Item 1"/>
    </>
}