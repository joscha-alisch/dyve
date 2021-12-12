import React from "react"
import ToolTip from "./toolTip";
import { ComponentMeta, ComponentStory} from "@storybook/react";

export default {
    title: 'Components/Atoms/Display/ToolTip',
    component: ToolTip,
} as ComponentMeta<typeof ToolTip>

export const StoryToolTip : ComponentStory<typeof ToolTip> = (args) => <ToolTip {...args}/>

StoryToolTip.storyName = "ToolTip"
StoryToolTip.args = {
    children: "Lorem ipsum dolor sit amet consectetur adipisicing elit. Qui sequi exercitationem beatae tenetur. Excepturi, quae, labore, eligendi adipisci officiis animi ipsam hic error explicabo ullam magnam placeat eaque repellendus obcaecati!"
}

