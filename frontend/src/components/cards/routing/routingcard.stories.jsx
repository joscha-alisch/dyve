import React from "react"
import RoutingCard from "./routingcard";

export default {
    title: 'Components/Cards/Routing',
    component: RoutingCard,
}

export const StoryRoutingCard= (args) => <RoutingCard {...args}/>

StoryRoutingCard.storyName = "Routing"
StoryRoutingCard.args = {
    routes: [
        { host: "some.domain.com", path: "/some/path", appPort: 8080},
        { host: "someother.domain.com", path: "", appPort: 8080}
    ]
}