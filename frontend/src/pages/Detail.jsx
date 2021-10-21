import {useParams} from "react-router-dom";
import {useEffect, useState} from "react";
import axios from "axios";
import Page from "../components/base/pages/page/page";

export const AppDetail = () => {
    const {id} = useParams()
    const [app, setApp] = useState({})

    useEffect(() => {
        axios.get("/api/apps/" + id)
            .then(res => res.data)
            .then((data) => {
                setApp(data.result)
            })
    }, [id])

    return <Page title={app.name} parent={"Applications"}>
    </Page>
}

