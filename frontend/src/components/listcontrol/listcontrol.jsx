import Pagination from "@material-ui/lab/Pagination";
import {PaginationItem} from "@material-ui/lab";
import {Link} from "react-router-dom";
import {FormControl, FormHelperText, Grid, InputLabel, MenuItem, Select} from "@material-ui/core";
import styles from "./listcontrol.module.sass"

const ListControl = ({totalResults, totalPages, page, perPage, setPerPage}) => <Grid container direction="row"
                                                                       justifyContent="space-between"
                                                                       alignItems="center">
    <FormHelperText>
        {((page-1)*perPage)+1}-{Math.min((page)*perPage, totalResults)} of {totalResults}
    </FormHelperText>
    <FormControl>
        <Pagination page={page} count={totalPages} siblingCount={4} renderItem={(item) => (
            <PaginationItem
                component={Link}
                to={`?page=${item.page}&perPage=${perPage}`}
                {...item}
            />
        )}/>
    </FormControl>
    <FormControl className={styles.PerPageSelect}>
        <InputLabel id="select-per-page-label">Per Page</InputLabel>
        <Select
            labelId="select-per-page-label"
            id="select-per-page"
            value={perPage}
            onChange={(e) => setPerPage(e.target.value)}
        >
            <MenuItem value={20}>20</MenuItem>
            <MenuItem value={50}>50</MenuItem>
            <MenuItem value={100}>100</MenuItem>
            <MenuItem value={200}>200</MenuItem>
            <MenuItem value={500}>500</MenuItem>
            <MenuItem value={20000}>All</MenuItem>
        </Select>
    </FormControl>
</Grid>

export default ListControl