import React, { useState } from "react"
import { PaginationValue } from "../../../molecules/input/pagination"
import { FilterEditor } from "../../../organisms"
import ListHeader from "../../../organisms/headers/listHeader"

type ListPageProps = {
    className?: string,
    title: string,
    category?: string,
}

const ListPage = ({
    className = "",
    title,
    category,
}: ListPageProps) => {
    let [filters, setFilters] = useState([])
    let [pagination, setPagination] = useState<PaginationValue>({
        page: 0,
        perPage: 10,
        totalItems: 2000,
    })

    return <div className={["", className].join(" ")}>
        <ListHeader filters={filters} onFilterChange={setFilters} pagination={pagination} onPaginationChange={setPagination} title={title} category={category}/>
    </div>
}

export default ListPage