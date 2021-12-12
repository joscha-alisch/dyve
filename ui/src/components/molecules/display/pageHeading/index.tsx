import React from "react"

type HeaderProps = {
    className?: string,
    title: string,
    category?: string
}

const PageHeading = ({
    className = "",
    category = "",
    title,
} : HeaderProps)  => <div className={["", className].join(" ")}>
    <span className="font-bold tracking-wide uppercase text-gray-400 text-xs">{category}</span>
    <h1 className="font-bold text-3xl">{title}</h1>
</div>

export default PageHeading