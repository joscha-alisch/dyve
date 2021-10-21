import {AppList, PipelineList, TeamList} from "./List";
import {AppDetail} from "./Detail";
import {Logout} from "./Logout"

export {TeamList, AppList, PipelineList, AppDetail, Logout}

export default {
    Teams: {
        List: TeamList
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