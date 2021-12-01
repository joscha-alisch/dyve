import ListPage from "../components/base/pages/listpage/listpage";
import {fetchList} from "../helpers/fetchList";
import AppCard from "../components/cards/app/appcard";
import TeamCard from "../components/cards/team/teamCard";

export const AppList = () => <ListPage
    title="Applications"
    parent="Platform"
    fetchItems={fetchList("apps")}
    itemRender={AppCard}
/>

export const PipelineList = () => <ListPage
    title="Pipelines"
    parent="Platform"
    fetchItems={fetchList("pipelines")}
    itemRender={AppCard}
/>

export const TeamList = () => <ListPage
    title="Teams"
    parent="Manage"
    fetchItems={fetchList("teams")}
    newItemRoute={"/teams/new"}
    newItemLabel={"Create Team"}
    itemRender={TeamCard}
/>