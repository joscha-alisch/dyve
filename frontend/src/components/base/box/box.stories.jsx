import React from "react"
import Box from "./box";

export default {
    title: 'Components/Box',
    component: Box,
}

export const StoryBox = (args) => <div style={{padding: "50px"}}>
    <Box {...args}>
        {args.children}
    </Box>
</div>

StoryBox.storyName = "Box"
StoryBox.args = {
    title: "Box Title",
    children: <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit. Aliquid atque debitis illum ipsum iste iusto
        labore magni maxime molestias, nisi porro quam quasi quo, quod rerum similique tempora tenetur unde!</p>
}