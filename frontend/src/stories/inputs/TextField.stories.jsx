import {TextField} from "@mui/material";
import React from "react"

export default {
    title: 'Components/inputs/Text Field',
    component: TextField,
}

export const Overview = (args) => <React.Fragment>
    <h2>Outlined</h2>
    <TextField {...args}/>;
</React.Fragment>

const Template = (args) => <TextField {...args}/>;

export const Outlined = Template.bind({});
Outlined.args = {
    label: 'Input',
    placeholder: "Placeholder Value",
    variant: "outlined"
};
