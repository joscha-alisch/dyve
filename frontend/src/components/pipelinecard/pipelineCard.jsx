import {Chip, Paper} from "@material-ui/core";
import styles from "./pipelineCard.module.sass"
import {Skeleton} from "@material-ui/lab";
import {Link} from "react-router-dom";

const PipelineCard = ({className, pipeline, loading}) => {
    if (loading) {
        return <Paper className={styles.App + " " + className} elevation={0}>
            <Skeleton animation={"wave"} variant={"text"} width="30%" height={30}/>
            <Skeleton animation={"wave"} variant={"text"} width="15%"/>
        </Paper>
    }

    return <Paper className={styles.App + " " + className} elevation={0}>
        <h1 className={styles.Name}>
            <span className={styles.Status}/>
            <Link to={"/pipelines/" + pipeline.id}>{pipeline.name}</Link>
        </h1>
    </Paper>
}

export default PipelineCard