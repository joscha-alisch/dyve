import {useParams} from "react-router";
import {Fragment, useEffect, useState} from "react";
import Heading from "../heading/heading";

const PipelineDetail = () => {
    const {id} = useParams()

    const [pipeline, setPipeline] = useState({})

    useEffect(() => {
        fetch("/api/pipelines/" + id)
            .then(res => res.json())
            .then((data) => {
                setPipeline(data.result)
            })
    }, [id])

    return <Fragment>
        <Heading title={pipeline.name} backlink="/pipelines" backlinkTitle="Pipelines"/>
    </Fragment>
}

export default PipelineDetail