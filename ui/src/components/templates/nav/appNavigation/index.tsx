import React, { FunctionComponent, MouseEventHandler } from "react"
import internal from "stream"
import { Icon, List } from "../../../atoms"
import { Icons } from "../../../atoms/display/icon/icons"

export type NavElement = {
    label: string,
    icon: Icons,
}
export type NavGroup = NavElement[]
export type Navigation = NavGroup[]

type AppNavigationProps = {
    className?: string,
    nav: Navigation
}

const SingleNav : FunctionComponent<NavElement> = ({
    label,
    icon
}) =>  <li className="py-2 group-scope hover:cursor-pointer rounded-l-md w-full ">
<div className="w-full py-3 group-scope-hover:bg-gray-100 flex flex-col items-center">
    <Icon icon={icon} className="h-6 w-6 transform group-hover:scale-90 translate-y-2 group-hover:-translate-y-1 transition-transform"/>
    <span className="text-tiny opacity-0 group-hover:opacity-100 transition-opacity">{label}</span>
</div>
<div className="w-56 p-6 bg-gray-100 absolute z-10 top-0 left-16 hidden group-scope-hover:block  h-full overflow-y-scroll shadow-inner">
    <h2 className="font-bold uppercase text-xs text-gray-400">Team</h2>
    <h2 className="font-bold text-lg text-gray-900">Projects</h2>

    <ul className="w-full space-y-6 mt-4">
        <li>
            <ul>
                <li>All</li>
                <li>My Team</li>
            </ul>
        </li>
        <li>
            <h3 className="uppercase text-xs tracking-wide text-gray-300 font-bold">Projects</h3>
            <ul>
                <li>CMS</li>
                <li>Production</li>
            </ul>
        </li>
        <li>
            <h3 className="uppercase text-xs tracking-wide text-gray-300 font-bold">Filters</h3>
            <ul>
                <li>CMS</li>
                <li>Production</li>
            </ul>
            <span className="text-gray-400 text-xs hover:text-indigo-400 hover:cursor-pointer">+ Add Filter</span>
        </li>
        <li>
            <h3 className="uppercase text-xs tracking-wide text-gray-300 font-bold">Pinned</h3>
            <ul></ul>
            <span className="text-gray-400 text-xs hover:text-indigo-400 hover:cursor-pointer">+ Pin Current</span>
        </li>
    </ul>
</div>
</li>

const AppNavigation: FunctionComponent<AppNavigationProps> = ({
    className = "",
    nav, 
}) => <aside className={["", className].join(" ")}>
        <nav className="sticky top-[23px] left-0 h-screen group flex-none flex flex-row bg-gray-50 w-12 hover:w-72 transition-all overflow-hidden">
            <ul className="flex flex-col items-center text-gray-700">
               {nav && nav.map(group => group.map(elem => <SingleNav {...elem}/>))}
            </ul>
        </nav>
    </aside>

export default AppNavigation