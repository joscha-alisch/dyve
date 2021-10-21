import Pagination from "@mui/material/Pagination";
import {PaginationItem} from "@mui/material";
import {Link} from "react-router-dom";
import {FormControl, FormHelperText, Grid, InputLabel, MenuItem, Select} from "@mui/material";
import styles from "./paginationControl.module.sass"
import PropTypes from "prop-types"
import {useTheme} from "../../../../context/theme";
import {useQueryParam} from "use-query-params";

const PaginationControl = ({totalResults, page, perPage, setPerPage}) => {
    let totalPages = Math.round(totalResults / perPage)

    // page is 0-based internally
    let displayPage = page + 1

    let startResult = (page * perPage) + 1
    let endResult = Math.min((page+1)*perPage, totalResults)

    console.log(page)

    return <Grid container direction="row"
          justifyContent="space-between"
          alignItems="center">
        <FormHelperText className={styles.PaginationHelper}>
            {startResult}-{endResult} of {totalResults}
        </FormHelperText>
        <FormControl className={styles.PaginationLinks}>
            <Pagination page={displayPage} count={totalPages} siblingCount={4} renderItem={(item) => (
                <PaginationItem
                    component={Link}
                    to={`?page=${item.page - 1}&perPage=${perPage}`}
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
}

PaginationControl.propTypes = {
    totalResults: PropTypes.number,
    page: PropTypes.number,
    perPage: PropTypes.number,
    setPerPage: PropTypes.func,
}

export default PaginationControl