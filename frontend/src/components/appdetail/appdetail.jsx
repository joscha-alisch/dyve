import styles from "./appdetail.module.sass"
import {useParams} from "react-router";
import {Fragment, useEffect} from "react";

const AppDetail = () => {
    const {id} = useParams()

    useEffect(() => {

    }, [id])

    return <Fragment>
        {id}
    </Fragment>
}

export default AppDetail