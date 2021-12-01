import React from "react"
import AssociationsPicker from "./associationspicker";

export default {
    title: 'Components/Inputs/Associations Picker',
    component: AssociationsPicker,
}

export const StoryAssociationsPicker = (args) => <AssociationsPicker {...args}/>

const genOptions = (prefix, n) => {
    let res = []
    for (let i = 0; i < n; i++) {
        res.push(prefix + " " + (i + 1))
    }
    return res
}

StoryAssociationsPicker.storyName = "Associations Picker"
StoryAssociationsPicker.args = {
    fields: [
        {label: "Thing A", data: genOptions("Option ", 100), outputKey: "output1"},
        {label: "Thing B", data: genOptions("Option ", 100), outputKey: "output2"}
    ],
}