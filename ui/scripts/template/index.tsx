import React from "react"

type template_upperProps = {
    className?: string,
}

const template_upper = ({
    className = "",
} : template_upperProps)  => <div className={["", className].join(" ")}>

</div>

export default template_upper