import React, { ForwardedRef, FunctionComponent, LegacyRef, ReactChildren, ReactNode, RefObject } from "react"

type ToolTipProps = {
    className?: string,
    ref?: LegacyRef<HTMLDivElement>
    children?: ReactNode
}

const ToolTip  = React.forwardRef(({
    className = "",
    children
} : ToolTipProps, ref: ForwardedRef<HTMLDivElement>)  => <div ref={ref} className={["shadow-lg rounded border border-gray-200", className].join(" ")}>
    {children}
</div>)

export default ToolTip