import React, { Component, FunctionComponent, MouseEventHandler } from "react"
import * as HeroIcons from '@heroicons/react/solid'
import { Icons } from "./icons"

export interface IconProps {
    className?: string,
    icon: Icons,
    onClick?: MouseEventHandler
}

const Icon: FunctionComponent<IconProps> = ({
    className = "",
    icon,
    onClick = () => { }
}) => {
    const iconProps: React.ComponentProps<"svg"> = {
        onClick: onClick,
        className: className
    }

    switch (icon) {
        case "minus":
            return <HeroIcons.MinusIcon {...iconProps} />
        case "plus":
            return <HeroIcons.PlusIcon {...iconProps} />
        case "mail":
            return <HeroIcons.PaperAirplaneIcon {...iconProps} />
        case "search":
            return <HeroIcons.SearchIcon {...iconProps} />
        case "close":
            return <HeroIcons.XIcon {...iconProps} />
        case "point-left":
            return <HeroIcons.ChevronLeftIcon {...iconProps} />
        case "point-right":
            return <HeroIcons.ChevronRightIcon {...iconProps} />
        case "clipboard":
            return <HeroIcons.ClipboardListIcon {...iconProps} />
        case "chip":
            return <HeroIcons.ChipIcon {...iconProps} />
        case "code":
            return <HeroIcons.CodeIcon {...iconProps} />
        case "server":
            return <HeroIcons.ServerIcon {...iconProps} />
        case "":
            return null
    }
}

export default Icon