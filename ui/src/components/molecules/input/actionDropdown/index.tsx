import React, { MouseEventHandler, useRef, useState } from "react"
import { Icons } from "../../../atoms/display/icon/icons"
import { useClickAway } from "use-click-away";
import useHotkeys from "@reecelucas/react-use-hotkeys"

import FilterInput from "../../../atoms/input/filterInput";

import { ActionItem, ToolTip, ListGroup, List, Action } from "../../../atoms"

export type ActionDropdownOption = {
    icon?: Icons,
    label: string,
    onClick: MouseEventHandler,
    group: string,
}

type ActionDropdownProps = {
    className?: string,
    label: string,
    icon: Icons,
    options: ActionDropdownOption[]
}

const ActionDropdown = ({
    className = "",
    label,
    icon,
    options,
}: ActionDropdownProps) => {
    const [isOpen, setIsOpen] = useState(false)
    const [filterValue, setFilterValue] = useState("")

    const clickRef = useRef(null);
    
    const close = () => {
        setIsOpen(false)
        setFilterValue("")
    }

    useClickAway(clickRef, close);
    useHotkeys("Escape", close);

    const filtered = options && options.filter(option => option.label.toLowerCase().includes(filterValue.toLowerCase()))

    const groups : {[key: string]: ActionDropdownOption[]} = {}
    filtered.forEach(option => {
        if (!groups[option.group]) {
            groups[option.group] = [option]
        } else {
            groups[option.group] = [...groups[option.group], option]
        }
    })


    return <div className={["", className].join(" ")}>
        <Action className="text-xs" icon={icon} label={label} onClick={() => setIsOpen(!isOpen)} />
        { isOpen && <ToolTip ref={clickRef} className="z-50 bg-gray-50 w-56 absolute max-w-xs rounded flex flex-col mt-2" >
                <FilterInput autoFocus className="mt-2 px-3" value={filterValue} onChange={setFilterValue} />
                <List className="mt-4 max-h-48 overflow-scroll">
                    { Object.keys(groups).map(group => <ListGroup label={group}>
                        {groups[group] && groups[group].map(option => <ActionItem className="text-sm" label={option.label} icon={option.icon} onClick={(ev) => {
                            close()
                            option.onClick(ev)
                        }}/>)}
                    </ListGroup>)}
                </List>
            </ToolTip>}
    </div>
}

export default ActionDropdown