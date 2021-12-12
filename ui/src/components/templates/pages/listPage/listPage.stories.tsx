import React from "react"
import ListPage from ".";
import { ComponentMeta, ComponentStory} from "@storybook/react";

export default {
    title: 'Components/Templates/Pages/List',
    component: ListPage,
} as ComponentMeta<typeof ListPage>

export const StoryListPage : ComponentStory<typeof ListPage> = (args) => <ListPage {...args}/>

StoryListPage.storyName = "List"
StoryListPage.args = {
    title: "Apps",
    category: "Platform"
}