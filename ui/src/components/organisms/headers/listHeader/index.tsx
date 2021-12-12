import React, { useState } from "react"
import Filter from "../../../molecules/input/filter/filter"
import ActionDropdown from "../../../molecules/input/actionDropdown"
import { PageHeading, Pagination } from "../../../molecules"
import { PaginationValue } from "../../../molecules/input/pagination"

type FilterData = {
    key: string,
    value: string
}

type FilterEditorProps = {
    className?: string,
    title: string,
    category?: string,
    filters: FilterData[],
    pagination: PaginationValue,
    onFilterChange: (filters: FilterData[]) => void
    onPaginationChange: (pagination: PaginationValue) => void
}

const ListHeader = ({
    title,
    category,
    className = "",
    filters = [],
    pagination,
    onFilterChange,
    onPaginationChange
}: FilterEditorProps) => {
    const onChangeAtIndex = (index: number, value: FilterData) => {
        const newValue = [...filters]
        newValue[index] = value
        onFilterChange(newValue)
    } 

    const onRemoveAtIndex = (index: number) => {
        const newValue = [...filters]
        newValue.splice(index, 1)
        onFilterChange(newValue)
        closeToolTip(index)
    }

    
    const [toolTipIndex, setToolTipIndex] = useState(-1)
    const closeToolTip = (index: number) => {
        if (toolTipIndex === index) {
            setToolTipIndex(-1)
        }
    }

    const add = (value: FilterData) => {
        onFilterChange([...filters, value])
        if (value.key === "" || value.value === "") {
            setToolTipIndex(filters.length)
        }
    }

    return <div className={["bg-white p-6", className].join(" ")}>
        <PageHeading title={title} category={category} className="mb-4"/>
        <div className="flex flex-row flex-wrap gap-2 mb-4">
            { filters.map((filter, index) => 
                <Filter open={toolTipIndex === index} onClose={() => closeToolTip(index)} onOpen={() => setToolTipIndex(index)} filterKey={filter.key} filterValue={filter.value} onRemove={() => onRemoveAtIndex(index)} onChange={(key, value) => onChangeAtIndex(index, {key, value})} />
            )}
            <ActionDropdown icon="plus" label="Add Filter" options={[
                {label: "Running", group: "App State", onClick: () => add({key: "state", value: "running"}) },
                {label: "Stopped", group: "App State", onClick: () => add({key: "state", value: "stopped"}) },
                {label: "Crashed", group: "App State", onClick: () => add({key: "state", value: "crashed"}) },
                {label: "By Space", group: "Location", onClick: () => add({key: "space", value: ""})},
                {label: "By Org", group: "Location", onClick: () => add({key: "org", value: ""})},
                {label: "key: value", group: "", onClick: () => add({key: "key", value: "value"}) }
            ]}/>
        </div>
        <Pagination value={pagination} onChange={onPaginationChange} />
    </div>
}

export default ListHeader