import React from "react";
import AppNavigation from ".";
import { ComponentMeta, ComponentStory} from "@storybook/react";

export default {
    title: 'Components/AppNavigation',
    component: AppNavigation,
} as ComponentMeta<typeof AppNavigation>

export const StoryAppNavigation : ComponentStory<typeof AppNavigation> = (args) => <AppNavigation {...args}/>

StoryAppNavigation.storyName = "AppNavigation"
StoryAppNavigation.args = {
}