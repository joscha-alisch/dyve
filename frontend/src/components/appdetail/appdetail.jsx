import {useParams} from "react-router";
import {Fragment, useEffect, useState} from "react";
import Heading from "../heading/heading";
import axios from "axios";

const AppDetail = () => {
    const {id} = useParams()

    const [app, setApp] = useState({})

    useEffect(() => {
        axios.get("/api/apps/" + id)
            .then(res => res.data)
            .then((data) => {
                setApp(data.result)
            })
    }, [id])

    return <Fragment>
        <Heading title={app.name} backlink="/apps" backlinkTitle="Apps"/>
    </Fragment>
}

export default AppDetail