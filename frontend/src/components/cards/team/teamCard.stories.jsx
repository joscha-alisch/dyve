import React from "react"
import TeamCard from "./teamCard";

export default {
    title: 'Components/Cards/App',
    component: TeamCard,
}

export const StoryTeamCard = (args) => <TeamCard {...args}/>

StoryTeamCard.storyName = "Team"
StoryTeamCard.args = {
    value: {
        name: "My Team"
    }
}