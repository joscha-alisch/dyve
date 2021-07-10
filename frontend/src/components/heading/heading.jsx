import styles from "./heading.module.sass";
import {Link} from "react-router-dom";
import {Fragment} from "react";
import {Skeleton} from "@material-ui/lab";

const Heading = ({title, backlinkTitle, backlink}) => {
    return <Fragment>
        <h2 className={styles.BackLink}>{(backlink !== null) ? <Link to={backlink}>{backlinkTitle}</Link>:""}</h2>
        {(title !== undefined)? <h1 className={styles.Heading}>{title}</h1> : <Skeleton variant="text" animation={"wave"}  width={200} /> }
    </Fragment>
}

export default Heading