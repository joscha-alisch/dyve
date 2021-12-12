import React from "react"
import ActionDropdown from ".";
import { ComponentMeta, ComponentStory } from "@storybook/react";
import {Icon} from "../../../atoms";

export default {
    title: 'Components/Molecules/Input/Action Dropdown',
    component: ActionDropdown,
} as ComponentMeta<typeof ActionDropdown>

export const StoryActionDropdown: ComponentStory<typeof ActionDropdown> = (args) => {

    return <ActionDropdown {...args} />
}

StoryActionDropdown.storyName = "Action Dropdown"
StoryActionDropdown.args = {
    label: "Action",
    icon: "plus",
    options: [
        {label: "Send Mail", group: "Group1", icon: "mail", onClick: () => {}},
        {label: "Add Recipient",group: "Group1", icon: "plus", onClick: () => {}},
        {label: "Remove Recipient", group: "Group2", icon: "minus", onClick: () => {}},
    ]
}