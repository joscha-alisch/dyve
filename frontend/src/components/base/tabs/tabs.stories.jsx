import React from "react"
import Tabs, {Tab} from "./tabs";

export default {
    title: 'Components/Tabs',
    component: Tabs,
}

export const StoryTabs = (args) => <Tabs {...args}>
    <Tab title={"Item 1"}>
        hi1
    </Tab>
    <Tab title={"Item 2"}>
        hi2
    </Tab>
    <Tab title={"Item 3"}>
        hi3
    </Tab>
</Tabs>

StoryTabs.storyName = "Tabs"
StoryTabs.args = {}