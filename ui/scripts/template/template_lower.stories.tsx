import React from "react"
import template_upper from ".";
import { ComponentMeta, ComponentStory} from "@storybook/react";

export default {
    title: 'Components/template_upper',
    component: template_upper,
} as ComponentMeta<typeof template_upper>

export const Storytemplate_upper : ComponentStory<typeof template_upper> = (args) => <template_upper {...args}/>

Storytemplate_upper.storyName = "template_upper"
Storytemplate_upper.args = {
}