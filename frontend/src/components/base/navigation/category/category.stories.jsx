import React from "react"
import Category from "./category";
import {faCoffee, faCog, faHouseUser, faLaptopCode} from "@fortawesome/free-solid-svg-icons";

export default {
    title: 'Components/Menu/Category',
    component: Category,
}

export const StoryCategory= (args) => <Category {...args}/>

StoryCategory.storyName = "Category"
StoryCategory.args = {
    title: "Category",
    items: [
        {
            to: "/item1",
            label: "Menu Item 1",
            icon: faLaptopCode,
        },
        {
            to: "/item2",
            label: "Menu Item 2",
            icon: faCoffee,
            soon: true
        },
        {
            to: "/item3",
            label: "Menu Item 3",
            icon: faHouseUser,
        },
        {
            to: "/item4",
            label: "Menu Item 4",
            icon: faCog,
        }
    ]
}