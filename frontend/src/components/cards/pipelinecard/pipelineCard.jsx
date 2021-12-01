import {Paper, Skeleton} from "@mui/material";
import styles from "./pipelineCard.module.sass"
import {Link} from "react-router-dom";
import SVG from 'react-inlinesvg';

const PipelineCard = ({className, value, svg, loading}) => {
    if (loading) {
        return <Paper className={styles.App + " " + className} elevation={0}>
            <Skeleton animation={"wave"} variant={"text"} width="30%" height={30}/>
            <Skeleton animation={"wave"} variant={"text"} width="15%"/>
        </Paper>
    }

    let pipelineViz
    if (!svg.hasOwnProperty(value.id)) {
        pipelineViz = <div className={styles.Svg}>
            <Skeleton className={styles.PipeSkeleton} animation={"wave"} variant={"rect"}/>
            <Skeleton className={styles.PipeSkeleton} animation={"wave"} variant={"rect"}/>
            <Skeleton className={styles.PipeSkeleton} animation={"wave"} variant={"rect"}/>
            <Skeleton className={styles.PipeSkeleton} animation={"wave"} variant={"rect"}/>
            <Skeleton className={styles.PipeSkeleton} animation={"wave"} variant={"rect"}/>
            <Skeleton className={styles.PipeSkeleton} animation={"wave"} variant={"rect"}/>
        </div>
    } else {
        pipelineViz =
            <SVG viewBox="0 0 3000 600" width={1000} height={200} className={styles.Svg} src={svg[value.id]}/>
    }

    return <Paper className={styles.App + " " + className} elevation={0}>
        <h1 className={styles.Name}>
            <span className={styles.Status}/>
            <Link to={"/pipelines/" + value.id}>{value.name}</Link>
        </h1>
        {pipelineViz}
    </Paper>
}

export default PipelineCard