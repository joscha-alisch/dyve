import React from "react"
import {StatusIcon} from ".";

export default {
    title: 'Components/Icons',
    isFullscreen: true,
    subcomponents: {StatusIcon},
}

let liStyle = {
    padding: "40px",
}

let iconStyle = {
    marginLeft: "50px",
    marginRight: "100px",
    display: "inline-block",
    textAlign: "center"
}

let icons = (args) => {
    return {
        Status: {
            Green: <StatusIcon {...args} status={"green"}/>,
            Red: <StatusIcon {...args} status={"red"}/>
        }
    }
}

export const StoryIcons = (args) => {
    let items = []
    let ics = icons(args)

    for (let icon in ics) {
        let variants = []
        for (let variant in ics[icon]) {
            variants.push(<div style={iconStyle}>{ics[icon][variant]}<br/><br/>{variant}</div>)
        }
        items.push(<li style={liStyle}><h2 style={{marginBottom: 50}}>{icon}</h2>{variants}</li>)
    }

    return <ul style={{listStyle: "none", padding: "20px"}}>
        {items}
    </ul>
}


StoryIcons.storyName = "Icons"
StoryIcons.args = {
    scale: 60,
    rotate: 0
}