import {
    FormControl,
    TablePagination
} from "@mui/material";
import styles from "./paginationControl.module.sass"
import PropTypes from "prop-types"

const PaginationControl = ({page, totalResults, perPage, setPerPage, setPage}) => {
    return <TablePagination
        labelRowsPerPage={"Per Page"}
        page={page}
        count={totalResults}
        rowsPerPage={perPage}
        onRowsPerPageChange={(e) => setPerPage(e.target.value)}
        onPageChange={(e, page) => setPage(page)}
    />
}

PaginationControl.propTypes = {
    totalResults: PropTypes.number,
    page: PropTypes.number,
    perPage: PropTypes.number,
    setPerPage: PropTypes.func,
    setPage: PropTypes.func
}

export default PaginationControl