import React from "react"
import PageHeading from ".";
import { ComponentMeta, ComponentStory} from "@storybook/react";

export default {
    title: 'Components/Molecules/Display/Page Heading',
    component: PageHeading,
} as ComponentMeta<typeof PageHeading>

export const StoryHeader : ComponentStory<typeof PageHeading> = (args) => <PageHeading {...args}/>

StoryHeader.storyName = "Page Heading"
StoryHeader.args = {
    title: "Apps",
    category: "Platform"
}