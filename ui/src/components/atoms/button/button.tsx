import * as React from "react"

type ButtonProps = {
    className: string,
    title: string,
}

const Button = (props : ButtonProps)  => <div className={["w-12 bg-red-500 h-12", props.className].join(" ")}>
    {props.title}
</div>

export default Button