import React from "react"
import AppCard from "./appcard";

export default {
    title: 'Components/Cards/App',
    component: AppCard,
}

export const StoryAppCard = (args) => <AppCard {...args}/>

StoryAppCard.storyName = "App"
StoryAppCard.args = {
    value: {
        name: "My App"
    }
}