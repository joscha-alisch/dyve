import React, { MouseEventHandler, useRef, useState } from "react"
import {Chip, Icon} from "../../../atoms"
import ToolTip from "../../../atoms/display/toolTip/toolTip"
import { useClickAway } from "use-click-away";
import useHotkeys from "@reecelucas/react-use-hotkeys"
import TextField from "../../../atoms/input/textField/textField"

type FilterProps = {
    className?: string,
    filterKey: string,
    filterValue: string,
    onChange: (key: string, value: string) => void
    onRemove: () => void,
    onOpen: () => void,
    onClose: () => void,
    open: boolean
}

const Filter = ({
    className = "",
    filterKey,
    filterValue,
    onChange,
    onRemove,
    onOpen,
    onClose,
    open,
} : FilterProps)  => {
    const clickRef = useRef(null);
    useClickAway(clickRef, onClose);
    useHotkeys("Escape", onClose);

    return <div className={["relative", className].join(" ")}>
            <Chip className="relative group" onClick={onOpen} label={filterKey} value={filterValue}>
                <div className="absolute hidden group-hover:flex inset-y-0 items-center justify-end w-full">
                    <Icon icon="close" onClick={onRemove} className="w-4 h-4 mr-4 bg-gray-400 hover:bg-gray-100 rounded-full" />
                </div>
            </Chip>
        { open && <ToolTip ref={clickRef} className="z-50 bg-gray-50 absolute max-w-xs p-4 flex flex-row space-x-3 mt-2" >
                    <TextField autoFocus={filterKey === "" || filterValue !== ""} onChange={(key: string) => onChange(key, filterValue)} value={filterKey} className="text-xs flex-auto inline-block box-border" />
                    <TextField autoFocus={filterValue === ""} value={filterValue} onChange={(value: string) => onChange(filterKey, value)} className="text-xs flex-auto inline-block box-border" />
            </ToolTip>}
    </div>
}

export default Filter