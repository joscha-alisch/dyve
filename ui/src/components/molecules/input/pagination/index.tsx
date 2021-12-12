import React, { FunctionComponent } from "react"
import { PageCounter, HorizontalSelect } from "../../../atoms"

export type PaginationValue = {
    totalItems: number,
    page: number,
    perPage: number,
}

type PaginationProps = {
    className?: string,
    value: PaginationValue,
    onChange: (value: PaginationValue) => void
}

const Pagination: FunctionComponent<PaginationProps> = ({
    className = "",
    value,
    onChange,
}) => {
    const onPageChange = (page: number) => {
        onChange({
            ...value,
            page
        })
    } 

    const onPerPageChange = (perPage: string | number) => {
        onChange({
            ...value,
            perPage: perPage as number
        })
    } 

    return <div className={["flex flex-row items-center justify-between", className].join(" ")}>
        <div></div>
        <PageCounter page={value.page} perPage={value.perPage} totalItems={value.totalItems} onPageChange={onPageChange} />
        <HorizontalSelect label="Per Page" options={[
             { label: "10", value: 10 },
             { label: "50", value: 50 },
             { label: "100", value: 100 },
             { label: "All", value: -1 },
        ]} selected={value.perPage} onSelect={onPerPageChange}/>
    </div>
}
export default Pagination