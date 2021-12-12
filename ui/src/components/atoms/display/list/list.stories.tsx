import React from "react"
import List from ".";
import { ComponentMeta, ComponentStory} from "@storybook/react";
import ActionItem from "../../input/actionItem"

export default {
    title: 'Components/Atoms/Display/List',
    component: List,
} as ComponentMeta<typeof List>

export const StoryFilterList : ComponentStory<typeof List> = (args) => <List {...args}/>

StoryFilterList.storyName = "List"
StoryFilterList.args = {
    children: <>
        <ActionItem label="Send Mail" icon="mail" />
        <ActionItem label="Add Recipient" icon="plus" />
    </>
}