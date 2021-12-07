import * as React from "react"

type ButtonProps = {
    className: string,
}

const Button = ({className} : ButtonProps)  => <div className={["w-12 bg-red-500 h-12", className].join(" ")}>

</div>

export default Button