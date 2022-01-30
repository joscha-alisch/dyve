import React from "react"
import InstancesCard from "./instancescard";

export default {
    title: 'Components/Cards/Instances',
    component: InstancesCard,
}

export const StoryRoutingCard= (args) => <InstancesCard {...args}/>

StoryRoutingCard.storyName = "Routing"
StoryRoutingCard.args = {
    instances: [
    ]
}