import React, {useState} from "react"
import ConditionBuilder from "./conditionbuilder";
import Box from "../../box/box";

export default {
    title: 'Components/Inputs/ConditionBuilder',
    component: ConditionBuilder,
}

export const StoryConditionBuilder = (args) => {
    let [value, setValue] = useState(args.value)

    return <Box>
        <ConditionBuilder {...args} value={value} onChange={setValue}/>
    </Box>
}

StoryConditionBuilder.storyName = "ConditionBuilder"
StoryConditionBuilder.args = {
    label: "Condition Builder",
    helperText: "Build your conditions",
    value: [
        {key: "Food", value: "Pizza"},
        {key: "Drink", value: "Soda"}
    ],
    options: {
        "Food": ["Pizza", "Pasta"],
        "Drink": ["Water", "Soda"]
    }
}