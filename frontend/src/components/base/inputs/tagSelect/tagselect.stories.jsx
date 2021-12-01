import React from "react"
import TagSelect from "./tagselect";

export default {
    title: 'Components/Inputs/Tag Select',
    component: TagSelect,
}

export const StoryTagSelect= (args) => <TagSelect {...args}/>

StoryTagSelect.storyName = "Tag Select"
StoryTagSelect.args = {
    label: "My Favourite Cities",
    helperText: "Select one or more cities",
    options: [
        { name: "Berlin", country: "Germany"},
        { name: "Stuttgart", country: "Germany" },
        { name: "New York", country: "USA" },
        { name: "Seattle", country: "USA" },
        { name: "Paris", country: "France" },
    ],
    groupBy: (option) => option.country,
    getOptionLabel: (option) => option.name,
}