import {AppList, PipelineList, TeamList} from "./List";
import {AppDetail, TeamDetail} from "./Detail";
import {Logout} from "./Logout"
import {TeamCreate} from "./Forms";

export {TeamList, AppList, PipelineList, AppDetail, Logout, TeamCreate}

const pages = {
    Teams: {
        List: TeamList,
        New: TeamCreate,
        Detail: TeamDetail,
    },
    Apps: {
        List: AppList,
        Detail: AppDetail,
    },
    Pipelines: {
        List: PipelineList
    },
    User: {
        Logout
    }
}

export default pages