import * as React from "react"

type ButtonProps = {
    className: string,
    title: string
}

const Button = ({className, title}: ButtonProps) => <button className={["w-12 bg-red-100",className].join(' ')}>
    {title}
</button>

export default Button