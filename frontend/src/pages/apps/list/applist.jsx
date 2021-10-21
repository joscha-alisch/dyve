import ListPage from "../../../components/base/pages/listpage/listpage";
import axios from "axios";
import AppCard from "../../../components/appcard/appcard";


const fetchApps = (perPage, page, setResults) => {
    axios.get("/api/apps?perPage=" + perPage + "&page=" + page)
        .then((res) => {
            if(res.data.result.apps) {
                setResults(res.data.result.apps, res.data.result.totalResults)
            }
        })
}

const itemRender = (app) => <AppCard app={app} />


const AppList = () => <ListPage title="Applications" parent="Platform" fetchItems={fetchApps} itemRender={itemRender}/>

export default AppList