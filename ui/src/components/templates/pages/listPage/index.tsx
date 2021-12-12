import React, { FunctionComponent, useState } from "react"
import { PaginationValue } from "../../../molecules/input/pagination"
import ListHeader, { FilterData } from "../../../organisms/headers/listHeader"

type ListPageProps = {
    className?: string,
    title: string,
    category?: string,
}

const ListPage : FunctionComponent<ListPageProps> = ({
    className = "",
    title,
    category,
}) => {
    let [filters, setFilters] = useState<FilterData[]>([])
    let [pagination, setPagination] = useState<PaginationValue>({
        page: 0,
        perPage: 10,
        totalItems: 2000,
    })

    return <div className={["", className].join(" ")}>
        <ListHeader filters={filters} onFilterChange={setFilters} pagination={pagination} onPaginationChange={setPagination} title={title} category={category} />
    </div>
}

export default ListPage