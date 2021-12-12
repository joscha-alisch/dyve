import React from "react"
import Chip from ".";
import { ComponentMeta, ComponentStory} from "@storybook/react";

export default {
    title: 'Components/Atoms/Display/Chip',
    component: Chip,
} as ComponentMeta<typeof Chip>

export const StoryChip : ComponentStory<typeof Chip> = (args) => <Chip {...args}/>

StoryChip.storyName = "Chip"
StoryChip.args = {
    specificer: "",
    label: "key",
    value: "value",
}