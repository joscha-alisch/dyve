import {Chip, Paper} from "@material-ui/core";
import styles from "./appcard.module.sass"
import {Skeleton} from "@material-ui/lab";
import {Link} from "react-router-dom";

const AppCard = ({className, app, loading}) => {
    if (loading) {
        return <Paper className={styles.App + " " + className} elevation={1}>
            <Skeleton animation={"wave"} variant={"text"} width="30%" height={30}/>
            <Skeleton animation={"wave"} variant={"text"} width="15%"/>
        </Paper>
    }

    return <Paper className={styles.App + " " + className} elevation={1}>
        <h1 className={styles.Name}>
            <span className={styles.Status}/>
            <Link to={"/apps/" + app.id}>{app.name}</Link>
        </h1>

        {Object.keys(app.meta).map((k) => <Chip
            size="small"
            label={k+": "+app.meta[k]}
        />)}
    </Paper>
}

export default AppCard