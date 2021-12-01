import React, {useState} from "react"
import {Form} from "../packages/formidable";
import {FormTeamAccess, FormTeamAssociations, FormTeamGeneral} from "./team";
import {components} from "../components/base/forms/components/components";
import {withComponents} from "../packages/formidable/components";
import Box from "../components/base/box/box";

export default {
    title: 'App/Forms',
    component: Form,
}

const Template = (args) => {
    let [state, setState] = useState({})

    return <>
        <div style={{
            boxSizing: "border-box",
            padding: "30px",
            marginBottom: "50px",
            height: "200px",
            overflow: "scroll",
            borderBottom: "1px solid grey"
        }}>
            <pre>{JSON.stringify(state, null, 2)}</pre>
        </div>
        <Box>
            <Form {...args} components={withComponents(components)} setState={setState}/>
        </Box>
    </>
}

export const StoryTeamGeneral = Template.bind({})
StoryTeamGeneral.storyName = "Team | General"
StoryTeamGeneral.args = FormTeamGeneral

export const StoryTeamAccess = Template.bind({})
StoryTeamAccess.storyName = "Team | Access"
StoryTeamAccess.args = {
    ...FormTeamAccess,
    data: {
        groups: [
            { name: "group a", provider: "github", providerName: "GitHub"},
            { name: "group b", provider: "github", providerName: "GitHub"},
            { name: "group c", provider: "SSO", providerName: "Company"},
            { name: "group d", provider: "SSO", providerName: "Company"},
        ]

    }
}

export const StoryTeamPermissions = Template.bind({})
StoryTeamPermissions.storyName = "Team | Associations"
StoryTeamPermissions.args = FormTeamAssociations