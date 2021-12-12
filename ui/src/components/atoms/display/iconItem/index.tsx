import React, { FunctionComponent } from "react"
import {Icon} from "../../../atoms"
import { Icons } from "../icon/icons"

type IconItemProps = {
    className?: string,
    icon: Icons,
    label: string
}

const IconItem: FunctionComponent<IconItemProps> = ({
    className = "",
    icon,
    label
})  => <div className={["w-full py-3 group-scope-hover:bg-gray-100 flex flex-col items-center", className].join(" ")}>
        <Icon className="h-6 w-6 transform group-hover:scale-90 translate-y-2 group-hover:-translate-y-1 transition-transform" icon={icon}/>
        <span className="text-tiny opacity-0 group-hover:opacity-100 transition-opacity">{label}</span>
</div>

export default IconItem