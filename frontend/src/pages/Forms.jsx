import Page from "../components/base/pages/page/page";
import MultiPageForm from "../components/base/forms/multipage/multiPageForm";
import {FormTeamAccess, FormTeamAssociations, FormTeamGeneral} from "../forms/team";
import API from "../api";
import history from "../helpers/history"
import {flatMap, forEach} from "../helpers/object";

const getTeamData = async () => API.Groups.List((json) => {
    let groups = flatMap(json.result, (provider) => provider.groups.map(group => {
        return {
            id: group.id,
            name: group.name,
            provider: provider.provider,
            providerName: provider.name,
        }
    }))

    return { groups }
})

const submitTeam = (data, notify) => API.Teams.Create(
    {
        ...data,
        access: {
            admin: data.access.admin ? data.access.admin.map(group => group.provider + ":" + group.id) : [],
            member: data.access.member ? data.access.member.map(group => group.provider + ":" + group.id) : [],
            viewer: data.access.viewer ? data.access.viewer.map(group => group.provider + ":" + group.id) : [],
        }
    },
    () => {
        history.push("/teams/" + data.slug)
    },
    (err) => notify(err.response.statusText, "error")
)

export const TeamCreate = () => <Page title="Create New Team" parentRoute={"/teams"} parent="Teams">
    <MultiPageForm forms={[
        {title: "General", form: FormTeamGeneral},
        {title: "Access", form: FormTeamAccess, mappingKey: "access"},
        {title: "Permissions", form: FormTeamAssociations, mappingKey: "permissions"},
    ]} getData={getTeamData} onSubmit={submitTeam}/>
</Page>