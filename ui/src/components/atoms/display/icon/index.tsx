import React, { Component, FunctionComponent, MouseEventHandler } from "react"
import * as HeroSolid from '@heroicons/react/solid'
import * as HeroOutline from '@heroicons/react/outline'

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
            return <HeroSolid.MinusIcon {...iconProps} />
        case "plus":
            return <HeroSolid.PlusIcon {...iconProps} />
        case "mail":
            return <HeroSolid.PaperAirplaneIcon {...iconProps} />
        case "search":
            return <HeroSolid.SearchIcon {...iconProps} />
        case "close":
            return <HeroSolid.XIcon {...iconProps} />
        case "point-left":
            return <HeroSolid.ChevronLeftIcon {...iconProps} />
        case "point-right":
            return <HeroSolid.ChevronRightIcon {...iconProps} />
        case "clipboard":
            return <HeroOutline.ClipboardListIcon {...iconProps} />
        case "chip":
            return <HeroOutline.ChipIcon {...iconProps} />
        case "code":
            return <HeroOutline.CodeIcon {...iconProps} />
        case "server":
            return <HeroOutline.ServerIcon {...iconProps} />
        case "":
            return null
    }
}

export default Icon