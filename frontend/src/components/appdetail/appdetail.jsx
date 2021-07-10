import styles from "./appdetail.module.sass"
import {useParams} from "react-router";
import {Fragment, useEffect, useState} from "react";
import {Link} from "react-router-dom";
import Heading from "../heading/heading";

const AppDetail = () => {
    const {id} = useParams()

    const [app, setApp] = useState({})

    useEffect(() => {
        fetch("/api/apps/" + id)
            .then(res => res.json())
            .then((data) => {
                setApp(data.result)
            })
    }, [id])

    return <Fragment>
        <Heading title={app.name} backlink="/apps" backlinkTitle="Apps"/>
    </Fragment>
}

export default AppDetail